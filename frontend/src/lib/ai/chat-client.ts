// SSE 流式聊天客户端

import { Message, ChatChunk, ContextSource } from './types'
import { parseSSEChunk, isDoneChunk, tryParseJSON } from './stream-parser'

export interface ChatOptions {
  messages: Message[]
  contextSources?: ContextSource[]
  maxContextTokens?: number
  signal?: AbortSignal
  onChunk?: (chunk: ChatChunk) => void
  onError?: (error: Error) => void
  onComplete?: () => void
}

/**
 * 流式聊天客户端
 * 使用 SSE (Server-Sent Events) 接收后端流式响应
 */
export async function streamChat(options: ChatOptions): Promise<void> {
  const {
    messages,
    contextSources = [],
    maxContextTokens,
    signal,
    onChunk,
    onError,
    onComplete,
  } = options

  try {
    // 获取认证 token (从 localStorage 或其他地方)
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth-token') : null
    const orgId = typeof window !== 'undefined' ? localStorage.getItem('current-org-id') : null

    // 构建请求体
    const requestBody: any = {
      messages,
      stream: true,
    }

    if (contextSources.length > 0) {
      requestBody.context_sources = contextSources
    }

    if (maxContextTokens) {
      requestBody.max_context_tokens = maxContextTokens
    }

    // 发起请求
    const response = await fetch('/api/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...(orgId && { 'X-Organization-ID': orgId }),
      },
      body: JSON.stringify(requestBody),
      signal,
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || `HTTP ${response.status}: ${response.statusText}`)
    }

    if (!response.body) {
      throw new Error('Response body is null')
    }

    // 读取流
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()

      if (done) {
        break
      }

      // 解码数据块
      buffer += decoder.decode(value, { stream: true })

      // 按 \n\n 分割多个 SSE 消息
      const parts = buffer.split('\n\n')
      buffer = parts.pop() || '' // 保留最后一个未完成的部分

      for (const part of parts) {
        if (!part.trim()) continue

        const chunks = parseSSEChunk(part)

        for (const sseChunk of chunks) {
          // 检查是否为结束标记
          if (isDoneChunk(sseChunk)) {
            if (onComplete) onComplete()
            return
          }

          // 解析 JSON 数据
          const chatChunk = tryParseJSON<ChatChunk>(sseChunk.data)
          if (chatChunk && onChunk) {
            onChunk(chatChunk)
          }
        }
      }
    }

    // 处理剩余的 buffer
    if (buffer.trim()) {
      const chunks = parseSSEChunk(buffer)
      for (const sseChunk of chunks) {
        if (!isDoneChunk(sseChunk)) {
          const chatChunk = tryParseJSON<ChatChunk>(sseChunk.data)
          if (chatChunk && onChunk) {
            onChunk(chatChunk)
          }
        }
      }
    }

    if (onComplete) onComplete()
  } catch (error) {
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        // 请求被取消,不视为错误
        return
      }
      if (onError) onError(error)
    } else {
      if (onError) onError(new Error('Unknown error occurred'))
    }
  }
}

/**
 * 重连机制的流式聊天
 * 支持自动重连
 */
export async function streamChatWithRetry(
  options: ChatOptions,
  maxRetries = 3,
  retryDelay = 1000
): Promise<void> {
  let attempt = 0

  while (attempt < maxRetries) {
    try {
      await streamChat(options)
      return // 成功则返回
    } catch (error) {
      attempt++
      if (attempt >= maxRetries) {
        throw error // 达到最大重试次数,抛出错误
      }

      // 等待后重试
      await new Promise(resolve => setTimeout(resolve, retryDelay * attempt))
    }
  }
}
