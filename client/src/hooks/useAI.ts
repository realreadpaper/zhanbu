import { useState, useCallback, useRef, useEffect } from 'react'
import { interpretStream } from '../services/ai'

export interface UseAIOptions {
  type: string
  resultId?: number
  result?: string   // Direct result JSON (alternative to resultId)
  question?: string
  autoStart?: boolean
}

export interface UseAIReturn {
  /** The accumulated text from the AI stream */
  text: string
  /** Whether the stream is currently active */
  isStreaming: boolean
  /** Whether the stream has completed */
  isDone: boolean
  /** Any error that occurred */
  error: string | null
  /** Start the streaming interpretation */
  start: () => void
  /** Reset the state */
  reset: () => void
}

/**
 * Hook for streaming AI interpretation.
 */
export function useAI({ type, resultId, result, question, autoStart = false }: UseAIOptions): UseAIReturn {
  const [text, setText] = useState('')
  const [isStreaming, setIsStreaming] = useState(false)
  const [isDone, setIsDone] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const cleanupRef = useRef<(() => void) | null>(null)

  const start = useCallback(() => {
    // Reset state
    setText('')
    setIsStreaming(true)
    setIsDone(false)
    setError(null)

    // Start streaming
    cleanupRef.current = interpretStream(
      { type, result_id: resultId, result, question },
      // onChunk
      (chunk: string) => {
        setText((prev) => prev + chunk)
      },
      // onDone
      () => {
        setIsStreaming(false)
        setIsDone(true)
      },
      // onError
      (err: Error) => {
        setIsStreaming(false)
        setError(err.message || 'AI解读失败，请稍后重试')
      }
    )
  }, [type, resultId, result, question])

  const reset = useCallback(() => {
    // Cleanup any existing stream
    if (cleanupRef.current) {
      cleanupRef.current()
      cleanupRef.current = null
    }

    setText('')
    setIsStreaming(false)
    setIsDone(false)
    setError(null)
  }, [])

  // Auto-start if requested
  useEffect(() => {
    if (autoStart && (resultId || result)) {
      start()
    }

    // Cleanup on unmount
    return () => {
      if (cleanupRef.current) {
        cleanupRef.current()
      }
    }
  }, [autoStart, resultId, result, start])

  return {
    text,
    isStreaming,
    isDone,
    error,
    start,
    reset,
  }
}
