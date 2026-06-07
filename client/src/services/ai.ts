import { apiURL } from './api'

export interface InterpretRequest {
  type: string
  result_id?: number
  result?: string   // Direct result JSON (alternative to result_id)
  question?: string
  force?: boolean
}

export interface InterpretChunk {
  text: string
}

/**
 * Start streaming AI interpretation via SSE.
 * Returns an EventSource-like interface for streaming.
 */
export function interpretStream(
  request: InterpretRequest,
  onChunk: (text: string) => void,
  onDone: () => void,
  onError: (error: Error) => void
): () => void {
  const token = localStorage.getItem('access_token')

  // Use fetch with streaming for SSE
  const controller = new AbortController()

  const fetchData = async () => {
    try {
      const response = await fetch(apiURL('/ai/interpret'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(request),
        signal: controller.signal,
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => null)
        const message = errorData?.message || `HTTP error ${response.status}`

        if (response.status === 401) {
          throw new Error('登录已过期，请重新登录后再试')
        }
        if (response.status === 503 && message.includes('AI service is not configured')) {
          throw new Error('AI服务未配置，请在服务端设置 AI_API_KEY 后重试')
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

        // Process complete SSE messages
        const lines = buffer.split('\n')
        buffer = lines.pop() || '' // Keep incomplete line in buffer

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6).trim()

            if (data === '[DONE]') {
              onDone()
              return
            }

            try {
              const parsed: InterpretChunk = JSON.parse(data)
              if (parsed.text) {
                onChunk(parsed.text)
              }
            } catch {
              // Ignore parse errors for non-JSON data
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

  // Return cleanup function
  return () => {
    controller.abort()
  }
}
