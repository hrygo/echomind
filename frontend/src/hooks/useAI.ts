// AI相关的自定义 Hooks

import { useState, useCallback, useRef, useEffect } from 'react'
import { Message, ChatChunk, DraftResponse, AIError } from '@/lib/ai/types'
import { streamChat, ChatOptions } from '@/lib/ai/chat-client'
import { generateDraft, generateReply } from '@/lib/ai/draft-client'
import { useMutation } from '@tanstack/react-query'

/**
 * 流式聊天 Hook
 */
export function useStreamChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [currentMessage, setCurrentMessage] = useState('')
  const [isStreaming, setIsStreaming] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const abortControllerRef = useRef<AbortController | null>(null)

  const sendMessage = useCallback(
    async (options: Omit<ChatOptions, 'onChunk' | 'onError' | 'onComplete' | 'signal'>) => {
      // 取消之前的请求
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }

      // 创建新的 AbortController
      abortControllerRef.current = new AbortController()

      setIsStreaming(true)
      setError(null)
      setCurrentMessage('')

      // 添加用户消息
      const userMessage = options.messages[options.messages.length - 1]
      setMessages(prev => [...prev, userMessage])

      // 创建助手消息占位符
      const assistantMessageId = `msg-${Date.now()}`
      const assistantMessage: Message = {
        id: assistantMessageId,
        role: 'assistant',
        content: '',
        timestamp: new Date(),
      }
      setMessages(prev => [...prev, assistantMessage])

      try {
        await streamChat({
          ...options,
          signal: abortControllerRef.current.signal,
          onChunk: (chunk: ChatChunk) => {
            // 提取内容增量
            const delta = chunk.choices?.[0]?.delta?.content || ''
            if (delta) {
              setCurrentMessage(prev => prev + delta)
              // 更新消息列表中的助手消息
              setMessages(prev => {
                const updated = [...prev]
                const lastMessage = updated[updated.length - 1]
                if (lastMessage.role === 'assistant') {
                  lastMessage.content += delta
                }
                return updated
              })
            }
          },
          onError: (err: Error) => {
            setError(err)
            setIsStreaming(false)
          },
          onComplete: () => {
            setIsStreaming(false)
            setCurrentMessage('')
          },
        })
      } catch (err) {
        if (err instanceof Error && err.name !== 'AbortError') {
          setError(err)
        }
        setIsStreaming(false)
      }
    },
    []
  )

  const stopStreaming = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
      setIsStreaming(false)
    }
  }, [])

  const clearMessages = useCallback(() => {
    setMessages([])
    setCurrentMessage('')
    setError(null)
  }, [])

  // 清理
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [])

  return {
    messages,
    currentMessage,
    isStreaming,
    error,
    sendMessage,
    stopStreaming,
    clearMessages,
  }
}

/**
 * AI 草稿生成 Hook
 */
export function useAIDraft() {
  return useMutation<DraftResponse, AIError, { emailId: string; tone?: string; context?: string }>({
    mutationFn: async ({ emailId, tone, context }) => {
      return await generateDraft({
        email_id: emailId,
        tone,
        context,
      })
    },
  })
}

/**
 * AI 回复生成 Hook
 */
export function useAIReply() {
  return useMutation<DraftResponse, AIError, { emailId: string; tone?: string; context?: string }>({
    mutationFn: async ({ emailId, tone, context }) => {
      return await generateReply(emailId, tone, context)
    },
  })
}
