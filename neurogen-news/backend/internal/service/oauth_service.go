package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

// OAuthProvider represents an OAuth provider
type OAuthProvider string

const (
	ProviderGoogle   OAuthProvider = "google"
	ProviderVK       OAuthProvider = "vk"
	ProviderTelegram OAuthProvider = "telegram"
	ProviderGithub   OAuthProvider = "github"
)

// OAuthConfig holds OAuth configuration for all providers
type OAuthConfig struct {
	Google   GoogleConfig
	VK       VKConfig
	Telegram TelegramConfig
	Github   GithubConfig
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type VKConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type TelegramConfig struct {
	BotToken string
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// OAuthUser represents user data from OAuth provider
type OAuthUser struct {
	Provider    OAuthProvider
	ProviderID  string
	Email       string
	Username    string
	DisplayName string
	AvatarURL   string
}

// OAuthService handles OAuth authentication
type OAuthService interface {
	// Get authorization URL
	GetAuthURL(provider OAuthProvider, state string) (string, error)

	// Handle OAuth callback
	HandleCallback(ctx context.Context, provider OAuthProvider, code string) (*model.User, *AuthTokens, error)

	// Verify Telegram auth data
	VerifyTelegramAuth(data map[string]string) (*OAuthUser, error)

	// Link OAuth provider to existing user
	LinkProvider(ctx context.Context, userID uuid.UUID, provider OAuthProvider, code string) error

	// Unlink OAuth provider from user
	UnlinkProvider(ctx context.Context, userID uuid.UUID, provider OAuthProvider) error
}

type oauthService struct {
	userRepo    repository.UserRepository
	authService AuthService
	config      OAuthConfig
	httpClient  *http.Client
	logger      *zap.Logger
}

func NewOAuthService(
	userRepo repository.UserRepository,
	authService AuthService,
	config OAuthConfig,
	logger *zap.Logger,
) OAuthService {
	return &oauthService{
		userRepo:    userRepo,
		authService: authService,
		config:      config,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		logger:      logger,
	}
}

// GetAuthURL returns the authorization URL for the given provider
func (s *oauthService) GetAuthURL(provider OAuthProvider, state string) (string, error) {
	switch provider {
	case ProviderGoogle:
		return s.getGoogleAuthURL(state), nil
	case ProviderVK:
		return s.getVKAuthURL(state), nil
	case ProviderGithub:
		return s.getGithubAuthURL(state), nil
	case ProviderTelegram:
		return "", fmt.Errorf("telegram uses widget authentication, no redirect URL needed")
	default:
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
}

// HandleCallback handles OAuth callback from provider
func (s *oauthService) HandleCallback(ctx context.Context, provider OAuthProvider, code string) (*model.User, *AuthTokens, error) {
	var oauthUser *OAuthUser
	var err error

	switch provider {
	case ProviderGoogle:
		oauthUser, err = s.handleGoogleCallback(ctx, code)
	case ProviderVK:
		oauthUser, err = s.handleVKCallback(ctx, code)
	case ProviderGithub:
		oauthUser, err = s.handleGithubCallback(ctx, code)
	default:
		return nil, nil, fmt.Errorf("unknown provider: %s", provider)
	}

	if err != nil {
		return nil, nil, err
	}

	// Find or create user
	user, tokens, err := s.findOrCreateUser(ctx, oauthUser)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// VerifyTelegramAuth verifies Telegram authentication data
func (s *oauthService) VerifyTelegramAuth(data map[string]string) (*OAuthUser, error) {
	// Extract and verify hash
	hash := data["hash"]
	delete(data, "hash")

	// Check auth_date
	authDateStr := data["auth_date"]
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid auth_date")
	}

	// Check if auth is not too old (1 hour)
	if time.Now().Unix()-authDate > 3600 {
		return nil, fmt.Errorf("telegram auth expired")
	}

	// Build data check string
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, data[k]))
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// Calculate hash
	secretKey := sha256.Sum256([]byte(s.config.Telegram.BotToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != hash {
		return nil, fmt.Errorf("invalid telegram hash")
	}

	// Build OAuth user
	displayName := data["first_name"]
	if lastName := data["last_name"]; lastName != "" {
		displayName += " " + lastName
	}

	return &OAuthUser{
		Provider:    ProviderTelegram,
		ProviderID:  data["id"],
		Username:    data["username"],
		DisplayName: displayName,
		AvatarURL:   data["photo_url"],
	}, nil
}

// LinkProvider links OAuth provider to existing user
func (s *oauthService) LinkProvider(ctx context.Context, userID uuid.UUID, provider OAuthProvider, code string) error {
	// TODO: Implement provider linking
	return nil
}

// UnlinkProvider unlinks OAuth provider from user
func (s *oauthService) UnlinkProvider(ctx context.Context, userID uuid.UUID, provider OAuthProvider) error {
	// TODO: Implement provider unlinking
	return nil
}

// ============================================
// Google OAuth
// ============================================

func (s *oauthService) getGoogleAuthURL(state string) string {
	params := url.Values{
		"client_id":     {s.config.Google.ClientID},
		"redirect_uri":  {s.config.Google.RedirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"offline"},
	}
	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

func (s *oauthService) handleGoogleCallback(ctx context.Context, code string) (*OAuthUser, error) {
	// Exchange code for token
	tokenResp, err := s.httpClient.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {s.config.Google.ClientID},
		"client_secret": {s.config.Google.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {s.config.Google.RedirectURL},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer tokenResp.Body.Close()

	var tokenData struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	userResp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer userResp.Body.Close()

	var userData struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userData); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &OAuthUser{
		Provider:    ProviderGoogle,
		ProviderID:  userData.ID,
		Email:       userData.Email,
		DisplayName: userData.Name,
		AvatarURL:   userData.Picture,
	}, nil
}

// ============================================
// VK OAuth
// ============================================

func (s *oauthService) getVKAuthURL(state string) string {
	params := url.Values{
		"client_id":     {s.config.VK.ClientID},
		"redirect_uri":  {s.config.VK.RedirectURL},
		"display":       {"page"},
		"scope":         {"email"},
		"response_type": {"code"},
		"state":         {state},
		"v":             {"5.131"},
	}
	return "https://oauth.vk.com/authorize?" + params.Encode()
}

func (s *oauthService) handleVKCallback(ctx context.Context, code string) (*OAuthUser, error) {
	// Exchange code for token
	tokenResp, err := s.httpClient.Get("https://oauth.vk.com/access_token?" + url.Values{
		"client_id":     {s.config.VK.ClientID},
		"client_secret": {s.config.VK.ClientSecret},
		"redirect_uri":  {s.config.VK.RedirectURL},
		"code":          {code},
	}.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer tokenResp.Body.Close()

	var tokenData struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
		Email       string `json:"email"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	userResp, err := s.httpClient.Get("https://api.vk.com/method/users.get?" + url.Values{
		"user_ids":     {strconv.Itoa(tokenData.UserID)},
		"fields":       {"photo_200,screen_name"},
		"access_token": {tokenData.AccessToken},
		"v":            {"5.131"},
	}.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer userResp.Body.Close()

	var userData struct {
		Response []struct {
			ID         int    `json:"id"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			ScreenName string `json:"screen_name"`
			Photo200   string `json:"photo_200"`
		} `json:"response"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userData); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	if len(userData.Response) == 0 {
		return nil, fmt.Errorf("no user data returned from VK")
	}

	vkUser := userData.Response[0]
	displayName := vkUser.FirstName
	if vkUser.LastName != "" {
		displayName += " " + vkUser.LastName
	}

	return &OAuthUser{
		Provider:    ProviderVK,
		ProviderID:  strconv.Itoa(vkUser.ID),
		Email:       tokenData.Email,
		Username:    vkUser.ScreenName,
		DisplayName: displayName,
		AvatarURL:   vkUser.Photo200,
	}, nil
}

// ============================================
// GitHub OAuth
// ============================================

func (s *oauthService) getGithubAuthURL(state string) string {
	params := url.Values{
		"client_id":    {s.config.Github.ClientID},
		"redirect_uri": {s.config.Github.RedirectURL},
		"scope":        {"read:user user:email"},
		"state":        {state},
	}
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

func (s *oauthService) handleGithubCallback(ctx context.Context, code string) (*OAuthUser, error) {
	// Exchange code for token
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(url.Values{
		"client_id":     {s.config.Github.ClientID},
		"client_secret": {s.config.Github.ClientSecret},
		"code":          {code},
		"redirect_uri":  {s.config.Github.RedirectURL},
	}.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer tokenResp.Body.Close()

	var tokenData struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	userReq, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	userReq.Header.Set("Accept", "application/vnd.github.v3+json")

	userResp, err := s.httpClient.Do(userReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer userResp.Body.Close()

	body, _ := io.ReadAll(userResp.Body)

	var userData struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	displayName := userData.Name
	if displayName == "" {
		displayName = userData.Login
	}

	return &OAuthUser{
		Provider:    ProviderGithub,
		ProviderID:  strconv.Itoa(userData.ID),
		Email:       userData.Email,
		Username:    userData.Login,
		DisplayName: displayName,
		AvatarURL:   userData.AvatarURL,
	}, nil
}

// ============================================
// User management
// ============================================

func (s *oauthService) findOrCreateUser(ctx context.Context, oauthUser *OAuthUser) (*model.User, *AuthTokens, error) {
	// TODO: Check if OAuth provider is already linked to a user
	// For now, try to find user by email

	var user *model.User
	var err error

	if oauthUser.Email != "" {
		user, err = s.userRepo.GetByEmail(ctx, oauthUser.Email)
		if err == nil {
			// User exists, generate tokens
			tokens, err := s.authService.GenerateTokens(user.ID, user.Role)
			if err != nil {
				return nil, nil, err
			}
			return user, tokens, nil
		}
	}

	// Create new user
	username := oauthUser.Username
	if username == "" {
		username = generateUsername(oauthUser.DisplayName)
	}

	// Ensure username is unique
	for i := 0; i < 10; i++ {
		_, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil {
			break // Username is available
		}
		username = username + strconv.Itoa(i+1)
	}

	newUser := &model.User{
		Username:    username,
		Email:       oauthUser.Email,
		DisplayName: oauthUser.DisplayName,
		Role:        model.RoleUser,
	}

	if oauthUser.AvatarURL != "" {
		newUser.AvatarURL = &oauthUser.AvatarURL
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Save OAuth provider link

	tokens, err := s.authService.GenerateTokens(newUser.ID, newUser.Role)
	if err != nil {
		return nil, nil, err
	}

	return newUser, tokens, nil
}

// generateUsername creates a username from display name
func generateUsername(displayName string) string {
	// Convert to lowercase, replace spaces with underscores
	username := strings.ToLower(displayName)
	username = strings.ReplaceAll(username, " ", "_")

	// Remove non-alphanumeric characters
	var result strings.Builder
	for _, r := range username {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}

	username = result.String()
	if username == "" {
		username = "user"
	}

	// Limit length
	if len(username) > 20 {
		username = username[:20]
	}

	return username
}


