import { apiURL } from './api'

export interface ChatSession {
  id: number
  record_id: number
  title: string
  created_at: string
  updated_at: string
  messages?: ChatMessage[]
}

export interface ChatMessage {
  id: number
  role: 'user' | 'assistant'
  content: string
  created_at: string
}

export interface CreateSessionRequest {
  record_id?: number
  type?: string
  question?: string
}

export interface SendMessageRequest {
  content: string
}

export interface DivinationRecord {
  id: number
  user_id: number
  type: string
  question: string
  result: string
  ai_reading: string
  created_at: string
  prompt_profile_id?: string
  prompt_profile_name?: string
  prompt_profile_version?: string
}

/**
 * Create a new chat session for a divination record.
 */
export async function createSession(request: CreateSessionRequest): Promise<ChatSession> {
  const token = localStorage.getItem('access_token')
  const response = await fetch(apiURL('/chat/sessions'), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(request),
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => null)
    throw new Error(errorData?.message || `HTTP error ${response.status}`)
  }

  const data = await response.json()
  return data.data
}

/**
 * Get a chat session with messages.
 */
export async function getSession(sessionId: number): Promise<ChatSession> {
  const token = localStorage.getItem('access_token')
  const response = await fetch(apiURL(`/chat/sessions/${sessionId}`), {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => null)
    throw new Error(errorData?.message || `HTTP error ${response.status}`)
  }

  const data = await response.json()
  return data.data
}

/**
 * List chat sessions for the current user.
 */
export async function listSessions(page = 1, size = 20): Promise<{ sessions: ChatSession[], total: number }> {
  const token = localStorage.getItem('access_token')
  const response = await fetch(apiURL(`/chat/sessions?page=${page}&size=${size}`), {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => null)
    throw new Error(errorData?.message || `HTTP error ${response.status}`)
  }

  const data = await response.json()
  return data.data
}

/**
 * Delete a chat session.
 */
export async function deleteSession(sessionId: number): Promise<void> {
  const token = localStorage.getItem('access_token')
  const response = await fetch(apiURL(`/chat/sessions/${sessionId}`), {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => null)
    throw new Error(errorData?.message || `HTTP error ${response.status}`)
  }
}

/**
 * Send a message and stream the AI response via SSE.
 * Returns a cleanup function to abort the stream.
 */
export function sendMessageStream(
  sessionId: number,
  content: string,
  onChunk: (text: string) => void,
  onDone: () => void,
  onError: (error: Error) => void
): () => void {
  return streamChatEndpoint(
    `/chat/sessions/${sessionId}/messages`,
    { content },
    onChunk,
    onDone,
    onError
  )
}

export function streamInitialReading(
  sessionId: number,
  onChunk: (text: string) => void,
  onDone: () => void,
  onError: (error: Error) => void
): () => void {
  return streamChatEndpoint(
    `/chat/sessions/${sessionId}/initial-reading`,
    undefined,
    onChunk,
    onDone,
    onError
  )
}

function streamChatEndpoint(
  path: string,
  body: unknown,
  onChunk: (text: string) => void,
  onDone: () => void,
  onError: (error: Error) => void
): () => void {
  const token = localStorage.getItem('access_token')
  const controller = new AbortController()

  const fetchData = async () => {
    try {
      const response = await fetch(apiURL(path), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: body === undefined ? undefined : JSON.stringify(body),
        signal: controller.signal,
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => null)
        const message = errorData?.message || `HTTP error ${response.status}`

        if (response.status === 401) {
          throw new Error('登录已过期，请重新登录后再试')
        }
        if (response.status === 503) {
          throw new Error('AI服务暂时不可用，请稍后重试')
        }
        throw new Error(message)
      }

      if (!response.body) {
        throw new Error('ReadableStream not supported')
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()

        if (done) {
          onDone()
          break
        }

        buffer += decoder.decode(value, { stream: true })

        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6).trim()

            if (data === '[DONE]') {
              onDone()
              return
            }

            try {
              const parsed = JSON.parse(data)
              if (parsed.text) {
                onChunk(parsed.text)
              }
            } catch {
              // Ignore parse errors
            }
          }
        }
      }
    } catch (error) {
      if ((error as Error).name !== 'AbortError') {
        onError(error as Error)
      }
    }
  }

  fetchData()

  return () => {
    controller.abort()
  }
}
