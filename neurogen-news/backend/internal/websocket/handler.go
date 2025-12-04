package websocket

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

// Handler handles WebSocket connections
type Handler struct {
	hub         *Hub
	authService service.AuthService
	logger      *zap.Logger
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub, authService service.AuthService, logger *zap.Logger) *Handler {
	return &Handler{
		hub:         hub,
		authService: authService,
		logger:      logger,
	}
}

// Upgrade upgrades HTTP connection to WebSocket
func (h *Handler) Upgrade() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Get user ID from locals (set by auth middleware)
		var userID uuid.UUID
		if uid, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID); ok {
			userID = uid
		}

		client := &Client{
			ID:     uuid.New().String(),
			UserID: userID,
			Conn:   c,
			Send:   make(chan []byte, 256),
		}

		h.hub.Register(client)

		// Start goroutines for reading and writing
		go h.writePump(client)
		h.readPump(client)
	})
}

// readPump pumps messages from the WebSocket connection to the hub
func (h *Handler) readPump(client *Client) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		// Handle incoming message
		h.handleMessage(client, message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (h *Handler) writePump(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for message := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			h.logger.Error("WebSocket write error", zap.Error(err))
			return
		}
	}
}

// ClientMessage represents an incoming WebSocket message from client
type ClientMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// handleMessage processes incoming messages from clients
func (h *Handler) handleMessage(client *Client, data []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		h.logger.Error("Failed to unmarshal message", zap.Error(err))
		return
	}

	switch msg.Type {
	case "ping":
		// Respond with pong
		response := Message{Type: "pong", Payload: nil}
		data, _ := json.Marshal(response)
		client.Send <- data

	case "subscribe_article":
		// Subscribe to article updates (comments, reactions)
		var payload struct {
			ArticleID string `json:"articleId"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}
		// TODO: Track article subscriptions

	case "unsubscribe_article":
		// Unsubscribe from article updates
		var payload struct {
			ArticleID string `json:"articleId"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}
		// TODO: Remove article subscription

	case "typing":
		// User is typing (for comments)
		if client.UserID == uuid.Nil {
			return
		}
		var payload struct {
			ArticleID string `json:"articleId"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}
		// Broadcast typing indicator
		h.hub.Broadcast(Message{
			Type: TypeTyping,
			Payload: map[string]interface{}{
				"articleId": payload.ArticleID,
				"userId":    client.UserID.String(),
			},
		})

	default:
		h.logger.Debug("Unknown message type", zap.String("type", msg.Type))
	}
}

// UpgradeCheck middleware to check if request can be upgraded to WebSocket
func UpgradeCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}


