import { ref, type Ref } from 'vue'

export type MessageType = 
  | 'notification'
  | 'online_count'
  | 'reaction'
  | 'new_comment'
  | 'article_update'
  | 'typing'
  | 'achievement_unlock'
  | 'pong'

export interface WebSocketMessage<T = unknown> {
  type: MessageType
  payload: T
}

export interface NotificationPayload {
  id: string
  type: string
  title: string
  message: string
  link?: string
  actor?: {
    id: string
    username: string
    displayName: string
    avatarUrl?: string
  }
  createdAt: string
}

export interface OnlineCountPayload {
  count: number
}

export interface ReactionPayload {
  articleId: string
  reactions: Array<{
    emoji: string
    count: number
    isReacted: boolean
  }>
}

export interface NewCommentPayload {
  articleId: string
  comment: {
    id: string
    content: string
    author: {
      id: string
      username: string
      displayName: string
      avatarUrl?: string
    }
    createdAt: string
  }
}

export interface TypingPayload {
  articleId: string
  userId: string
}

export interface AchievementPayload {
  id: string
  name: string
  description: string
  icon: string
  points: number
}

type MessageHandler<T> = (payload: T) => void

class WebSocketClient {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private pingInterval: ReturnType<typeof setInterval> | null = null
  private handlers: Map<MessageType, Set<MessageHandler<unknown>>> = new Map()
  
  public isConnected: Ref<boolean> = ref(false)
  public onlineCount: Ref<number> = ref(0)

  connect() {
    const token = localStorage.getItem('accessToken')
    const wsUrl = import.meta.env.VITE_WS_URL || 
      `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`
    
    const url = token ? `${wsUrl}?token=${token}` : wsUrl

    try {
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.isConnected.value = true
        this.reconnectAttempts = 0
        this.startPing()
      }

      this.ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          this.handleMessage(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.isConnected.value = false
        this.stopPing()
        this.attemptReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
      this.attemptReconnect()
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.stopPing()
    this.isConnected.value = false
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
      console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)
      setTimeout(() => this.connect(), delay)
    } else {
      console.error('Max reconnect attempts reached')
    }
  }

  private startPing() {
    this.pingInterval = setInterval(() => {
      this.send('ping', null)
    }, 30000)
  }

  private stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }

  private handleMessage(message: WebSocketMessage) {
    // Handle online count update
    if (message.type === 'online_count') {
      this.onlineCount.value = (message.payload as OnlineCountPayload).count
    }

    // Call registered handlers
    const handlers = this.handlers.get(message.type)
    if (handlers) {
      handlers.forEach(handler => handler(message.payload))
    }
  }

  send(type: string, payload: unknown) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }))
    }
  }

  // Subscribe to article updates
  subscribeArticle(articleId: string) {
    this.send('subscribe_article', { articleId })
  }

  // Unsubscribe from article updates
  unsubscribeArticle(articleId: string) {
    this.send('unsubscribe_article', { articleId })
  }

  // Send typing indicator
  sendTyping(articleId: string) {
    this.send('typing', { articleId })
  }

  // Register message handler
  on<T>(type: MessageType, handler: MessageHandler<T>) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, new Set())
    }
    this.handlers.get(type)!.add(handler as MessageHandler<unknown>)
  }

  // Unregister message handler
  off<T>(type: MessageType, handler: MessageHandler<T>) {
    const handlers = this.handlers.get(type)
    if (handlers) {
      handlers.delete(handler as MessageHandler<unknown>)
    }
  }

  // Clear all handlers for a type
  clearHandlers(type: MessageType) {
    this.handlers.delete(type)
  }
}

// Singleton instance
export const wsClient = new WebSocketClient()

// Composable for Vue components
export function useWebSocket() {
  return {
    isConnected: wsClient.isConnected,
    onlineCount: wsClient.onlineCount,
    connect: () => wsClient.connect(),
    disconnect: () => wsClient.disconnect(),
    send: (type: string, payload: unknown) => wsClient.send(type, payload),
    subscribeArticle: (id: string) => wsClient.subscribeArticle(id),
    unsubscribeArticle: (id: string) => wsClient.unsubscribeArticle(id),
    sendTyping: (id: string) => wsClient.sendTyping(id),
    on: <T>(type: MessageType, handler: MessageHandler<T>) => wsClient.on(type, handler),
    off: <T>(type: MessageType, handler: MessageHandler<T>) => wsClient.off(type, handler),
  }
}

export default wsClient


