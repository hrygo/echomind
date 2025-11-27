// AI 草稿生成客户端

import { DraftRequest, DraftResponse, AIError } from './types'

export interface GenerateDraftOptions extends DraftRequest {
  signal?: AbortSignal
}

/**
 * 生成邮件草稿
 * @param options 生成选项
 * @returns 草稿响应
 */
export async function generateDraft(
  options: GenerateDraftOptions
): Promise<DraftResponse> {
  const { email_id, tone, context, signal } = options

  try {
    // 获取认证 token
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth-token') : null
    const orgId = typeof window !== 'undefined' ? localStorage.getItem('current-org-id') : null

    // 构建请求体
    const requestBody: DraftRequest = {
      email_id,
      tone,
      context,
    }

    // 发起请求
    const response = await fetch('/api/v1/ai/draft', {
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
      const errorData = await response.json().catch(() => ({
        code: 'UNKNOWN_ERROR',
        message: `HTTP ${response.status}: ${response.statusText}`,
      }))
      throw errorData as AIError
    }

    const data: DraftResponse = await response.json()
    return data
  } catch (error) {
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw {
          code: 'ABORTED',
          message: 'Request was cancelled',
        } as AIError
      }
      throw {
        code: 'NETWORK_ERROR',
        message: error.message,
      } as AIError
    }
    throw error
  }
}

/**
 * 生成邮件回复
 * @param emailId 邮件 ID
 * @param tone 语气
 * @param context 额外上下文
 * @returns 草稿响应
 */
export async function generateReply(
  emailId: string,
  tone?: string,
  context?: string,
  signal?: AbortSignal
): Promise<DraftResponse> {
  try {
    // 获取认证 token
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth-token') : null
    const orgId = typeof window !== 'undefined' ? localStorage.getItem('current-org-id') : null

    // 构建请求体
    const requestBody: DraftRequest = {
      email_id: emailId,
      tone,
      context,
    }

    // 发起请求
    const response = await fetch('/api/v1/ai/reply', {
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
      const errorData = await response.json().catch(() => ({
        code: 'UNKNOWN_ERROR',
        message: `HTTP ${response.status}: ${response.statusText}`,
      }))
      throw errorData as AIError
    }

    const data: DraftResponse = await response.json()
    return data
  } catch (error) {
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw {
          code: 'ABORTED',
          message: 'Request was cancelled',
        } as AIError
      }
      throw {
        code: 'NETWORK_ERROR',
        message: error.message,
      } as AIError
    }
    throw error
  }
}
