// SSE 流解析器

export interface SSEChunk {
  data: string
  event?: string
  id?: string
  retry?: number
}

/**
 * 解析 SSE 流数据
 * @param chunk 原始数据块
 * @returns 解析后的 SSE chunk 数组
 */
export function parseSSEChunk(chunk: string): SSEChunk[] {
  const lines = chunk.split('\n')
  const chunks: SSEChunk[] = []
  let currentChunk: Partial<SSEChunk> = {}

  for (const line of lines) {
    if (!line.trim()) {
      // 空行表示一个消息结束
      if (currentChunk.data !== undefined) {
        chunks.push(currentChunk as SSEChunk)
        currentChunk = {}
      }
      continue
    }

    // 解析字段
    const colonIndex = line.indexOf(':')
    if (colonIndex === -1) continue

    const field = line.slice(0, colonIndex).trim()
    const value = line.slice(colonIndex + 1).trim()

    switch (field) {
      case 'data':
        currentChunk.data = value
        break
      case 'event':
        currentChunk.event = value
        break
      case 'id':
        currentChunk.id = value
        break
      case 'retry':
        currentChunk.retry = parseInt(value, 10)
        break
    }
  }

  // 处理最后一个 chunk
  if (currentChunk.data !== undefined) {
    chunks.push(currentChunk as SSEChunk)
  }

  return chunks
}

/**
 * 检查是否为结束标记
 * @param chunk SSE chunk
 * @returns 是否为结束标记
 */
export function isDoneChunk(chunk: SSEChunk): boolean {
  return chunk.data === '[DONE]' || chunk.data === 'DONE'
}

/**
 * 尝试解析 JSON 数据
 * @param data 数据字符串
 * @returns 解析后的对象或 null
 */
export function tryParseJSON<T = unknown>(data: string): T | null {
  try {
    return JSON.parse(data) as T
  } catch {
    return null
  }
}
