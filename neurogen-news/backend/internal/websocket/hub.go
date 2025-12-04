package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/repository"
)

// Message types
const (
	TypeNotification      = "notification"
	TypeOnlineCount       = "online_count"
	TypeReaction          = "reaction"
	TypeNewComment        = "new_comment"
	TypeArticleUpdate     = "article_update"
	TypeTyping            = "typing"
	TypeAchievementUnlock = "achievement_unlock"
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID     string
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub manages all WebSocket connections
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// User ID to clients mapping (one user can have multiple connections)
	userClients map[uuid.UUID]map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Redis for pub/sub
	redis *repository.RedisClient

	// Logger
	logger *zap.Logger
}

// NewHub creates a new WebSocket hub
func NewHub(redis *repository.RedisClient, logger *zap.Logger) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[uuid.UUID]map[*Client]bool),
		broadcast:   make(chan []byte, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		redis:       redis,
		logger:      logger,
	}
}

// Run starts the hub's event loop
func (h *Hub) Run(ctx context.Context) {
	// Start Redis subscription for distributed messaging
	go h.subscribeToRedis(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// Add to user clients map
	if client.UserID != uuid.Nil {
		if h.userClients[client.UserID] == nil {
			h.userClients[client.UserID] = make(map[*Client]bool)
		}
		h.userClients[client.UserID][client] = true

		// Update online status in Redis
		h.redis.SAdd(context.Background(), "online:users", client.UserID.String())
	}

	h.logger.Debug("Client registered",
		zap.String("clientID", client.ID),
		zap.String("userID", client.UserID.String()))

	// Broadcast online count update
	h.broadcastOnlineCount()
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// Remove from user clients map
		if client.UserID != uuid.Nil {
			if userClients, ok := h.userClients[client.UserID]; ok {
				delete(userClients, client)
				if len(userClients) == 0 {
					delete(h.userClients, client.UserID)
					// Update online status in Redis
					h.redis.SRem(context.Background(), "online:users", client.UserID.String())
				}
			}
		}

		h.logger.Debug("Client unregistered",
			zap.String("clientID", client.ID),
			zap.String("userID", client.UserID.String()))

		// Broadcast online count update
		h.broadcastOnlineCount()
	}
}

func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendToUser sends a message to all connections of a specific user
func (h *Hub) SendToUser(userID uuid.UUID, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal message", zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.userClients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}

	// Also publish to Redis for distributed messaging
	h.redis.Publish(context.Background(), "ws:user:"+userID.String(), data)
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal message", zap.Error(err))
		return
	}

	h.broadcast <- data

	// Also publish to Redis for distributed messaging
	h.redis.Publish(context.Background(), "ws:broadcast", data)
}

// SendNotification sends a notification to a user
func (h *Hub) SendNotification(userID uuid.UUID, notification interface{}) {
	h.SendToUser(userID, Message{
		Type:    TypeNotification,
		Payload: notification,
	})
}

// SendAchievementUnlock notifies user about new achievement
func (h *Hub) SendAchievementUnlock(userID uuid.UUID, achievement interface{}) {
	h.SendToUser(userID, Message{
		Type:    TypeAchievementUnlock,
		Payload: achievement,
	})
}

// BroadcastNewComment notifies about new comment on article
func (h *Hub) BroadcastNewComment(articleID uuid.UUID, comment interface{}) {
	h.Broadcast(Message{
		Type: TypeNewComment,
		Payload: map[string]interface{}{
			"articleId": articleID,
			"comment":   comment,
		},
	})
}

// BroadcastReaction notifies about reaction on article
func (h *Hub) BroadcastReaction(articleID uuid.UUID, reactions interface{}) {
	h.Broadcast(Message{
		Type: TypeReaction,
		Payload: map[string]interface{}{
			"articleId": articleID,
			"reactions": reactions,
		},
	})
}

func (h *Hub) broadcastOnlineCount() {
	count, _ := h.redis.SCard(context.Background(), "online:users").Result()

	msg := Message{
		Type:    TypeOnlineCount,
		Payload: map[string]int64{"count": count},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	// Don't use the broadcast channel here to avoid deadlock
	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

// GetOnlineCount returns the number of online users
func (h *Hub) GetOnlineCount() int64 {
	count, _ := h.redis.SCard(context.Background(), "online:users").Result()
	return count
}

// IsUserOnline checks if a user is online
func (h *Hub) IsUserOnline(userID uuid.UUID) bool {
	result, _ := h.redis.SIsMember(context.Background(), "online:users", userID.String()).Result()
	return result
}

// subscribeToRedis listens for messages from other server instances
func (h *Hub) subscribeToRedis(ctx context.Context) {
	pubsub := h.redis.Subscribe(ctx, "ws:broadcast")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			// Broadcast to local clients
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- []byte(msg.Payload):
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}


