import { useState, useCallback, useRef, useEffect } from 'react'
import type { ChatSession, ChatMessage, DivinationRecord } from '../services/chat'
import {
  createSession,
  getSession,
  sendMessageStream,
  streamInitialReading,
} from '../services/chat'
import { fetchHistoryDetail } from '../services/history'

export interface UseChatOptions {
  /** The divination record ID to create a session for */
  recordId?: number
  /** Existing session ID to load */
  sessionId?: number
}

export interface UseChatReturn {
  /** Current chat session */
  session: ChatSession | null
  /** Divination record for the current chat session */
  record: DivinationRecord | null
  /** Messages in the session */
  messages: ChatMessage[]
  /** Whether AI is currently streaming a response */
  isStreaming: boolean
  /** Whether the session is loading */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
  /** Create a new session and optionally send first message */
  initSession: (recordId: number, firstMessage?: string) => Promise<void>
  /** Create a chat-mode session from a selected divination type */
  startModeSession: (type: string, question: string) => Promise<void>
  /** Load an existing session */
  loadSession: (sessionId: number) => Promise<void>
  /** Send a message */
  sendMessage: (content: string) => void
  /** Reset the chat state */
  reset: () => void
}

/**
 * Hook for managing chat state and interactions.
 */
export function useChat({ sessionId }: UseChatOptions = {}): UseChatReturn {
  const [session, setSession] = useState<ChatSession | null>(null)
  const [record, setRecord] = useState<DivinationRecord | null>(null)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [isStreaming, setIsStreaming] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const cleanupRef = useRef<(() => void) | null>(null)
  const streamContentRef = useRef('')

  // Refresh session and record (used after stream completes to pick up re-divination updates)
  const refreshSessionAndRecord = useCallback(async (sessionId: number) => {
    try {
      const updatedSession = await getSession(sessionId)
      setSession(updatedSession)
      setMessages(updatedSession.messages || [])
      const updatedRecord = await fetchHistoryDetail(updatedSession.record_id)
      setRecord(updatedRecord)
    } catch {
      // Silently ignore refresh errors — the UI already has the previous state
    }
  }, [])

  // Stream a message
  const streamMessage = useCallback((sessionId: number, content: string) => {
    setIsStreaming(true)
    streamContentRef.current = ''

    // Add placeholder for AI response
    const aiMsgId = Date.now() + 1
    setMessages(prev => [...prev, {
      id: aiMsgId,
      role: 'assistant',
      content: '',
      created_at: new Date().toISOString(),
    }])

    cleanupRef.current = sendMessageStream(
      sessionId,
      content,
      // onChunk
      (chunk: string) => {
        streamContentRef.current += chunk
        setMessages(prev => {
          const last = prev[prev.length - 1]
          if (last && last.role === 'assistant') {
            return [...prev.slice(0, -1), { ...last, content: streamContentRef.current }]
          }
          return prev
        })
      },
      // onDone
      () => {
        setIsStreaming(false)
        cleanupRef.current = null
        // Refresh session/record to pick up potential re-divination updates
        refreshSessionAndRecord(sessionId)
      },
      // onError
      (err: Error) => {
        setIsStreaming(false)
        setError(err.message || 'AI回复失败')
        cleanupRef.current = null

        // Update the placeholder with error
        setMessages(prev => {
          const last = prev[prev.length - 1]
          if (last && last.role === 'assistant' && !last.content) {
            return [...prev.slice(0, -1), { ...last, content: '抱歉，AI回复出现错误，请重试。' }]
          }
          return prev
        })
      }
    )
  }, [refreshSessionAndRecord])

  const streamInitialReadingForSession = useCallback((sessionId: number) => {
    setIsStreaming(true)
    streamContentRef.current = ''

    const aiMsgId = Date.now() + 1
    setMessages(prev => [...prev, {
      id: aiMsgId,
      role: 'assistant',
      content: '',
      created_at: new Date().toISOString(),
    }])

    cleanupRef.current = streamInitialReading(
      sessionId,
      (chunk: string) => {
        streamContentRef.current += chunk
        setMessages(prev => {
          const last = prev[prev.length - 1]
          if (last && last.role === 'assistant') {
            return [...prev.slice(0, -1), { ...last, content: streamContentRef.current }]
          }
          return prev
        })
      },
      () => {
        setIsStreaming(false)
        cleanupRef.current = null
        refreshSessionAndRecord(sessionId)
      },
      (err: Error) => {
        setIsStreaming(false)
        setError(err.message || 'AI解读失败')
        cleanupRef.current = null
        setMessages(prev => {
          const last = prev[prev.length - 1]
          if (last && last.role === 'assistant' && !last.content) {
            return [...prev.slice(0, -1), { ...last, content: '抱歉，AI解读出现错误，请重试。' }]
          }
          return prev
        })
      }
    )
  }, [refreshSessionAndRecord])

  // Initialize session from record
  const initSession = useCallback(async (recordId: number, firstMessage?: string) => {
    setIsLoading(true)
    setError(null)

    try {
      const nextSession = await createSession({ record_id: recordId })
      setSession(nextSession)
      setMessages(nextSession.messages || [])
      const nextRecord = await fetchHistoryDetail(nextSession.record_id)
      setRecord(nextRecord)

      if (firstMessage) {
        const userMsg: ChatMessage = {
          id: Date.now(),
          role: 'user',
          content: firstMessage,
          created_at: new Date().toISOString(),
        }
        setMessages(prev => [...prev, userMsg])
        streamMessage(nextSession.id, firstMessage)
      }
    } catch (err) {
      setError((err as Error).message || '创建会话失败')
    } finally {
      setIsLoading(false)
    }
  }, [streamMessage])

  const startModeSession = useCallback(async (type: string, question: string) => {
    if (isStreaming || !question.trim()) return

    setIsLoading(true)
    setIsStreaming(true)
    setError(null)
    streamContentRef.current = ''

    try {
      const nextSession = await createSession({ type, question: question.trim() })
      setSession(nextSession)
      setMessages(nextSession.messages || [])
      const nextRecord = await fetchHistoryDetail(nextSession.record_id)
      setRecord(nextRecord)
      streamInitialReadingForSession(nextSession.id)
    } catch (err) {
      setError((err as Error).message || '创建会话失败')
    } finally {
      setIsLoading(false)
    }
  }, [isStreaming, streamInitialReadingForSession])

  // Load existing session
  const loadSession = useCallback(async (sessionId: number) => {
    setIsLoading(true)
    setError(null)

    try {
      const nextSession = await getSession(sessionId)
      setSession(nextSession)
      setMessages(nextSession.messages || [])
      const nextRecord = await fetchHistoryDetail(nextSession.record_id)
      setRecord(nextRecord)
    } catch (err) {
      setError((err as Error).message || '加载会话失败')
    } finally {
      setIsLoading(false)
    }
  }, [])

  // Send a message
  const sendMessage = useCallback((content: string) => {
    if (!session || isStreaming || !content.trim()) return

    // Add user message to UI immediately
    const userMsg: ChatMessage = {
      id: Date.now(),
      role: 'user',
      content: content.trim(),
      created_at: new Date().toISOString(),
    }
    setMessages(prev => [...prev, userMsg])

    // Send to backend and stream response
    streamMessage(session.id, content.trim())
  }, [session, isStreaming, streamMessage])

  // Reset state
  const reset = useCallback(() => {
    if (cleanupRef.current) {
      cleanupRef.current()
      cleanupRef.current = null
    }
    setSession(null)
    setRecord(null)
    setMessages([])
    setIsStreaming(false)
    setIsLoading(false)
    setError(null)
    streamContentRef.current = ''
  }, [])

  // Auto-load if sessionId provided
  useEffect(() => {
    if (sessionId) {
      const timeoutId = window.setTimeout(() => {
        void loadSession(sessionId)
      }, 0)
      return () => window.clearTimeout(timeoutId)
    }
  }, [sessionId, loadSession])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (cleanupRef.current) {
        cleanupRef.current()
      }
    }
  }, [])

  return {
    session,
    record,
    messages,
    isStreaming,
    isLoading,
    error,
    initSession,
    startModeSession,
    loadSession,
    sendMessage,
    reset,
  }
}
