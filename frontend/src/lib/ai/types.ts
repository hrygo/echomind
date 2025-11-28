// AI 服务相关类型定义

export interface Message {
  id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp?: Date;
}

export interface ChatChunk {
  id?: string;
  choices?: Array<{
    delta: {
      content?: string;
    };
    finish_reason?: string | null;
  }>;
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
  metadata?: Record<string, unknown>;
}

export interface ContextSource {
  type: 'email' | 'context' | 'document' | 'search';
  id: string;
  priority?: number;
  metadata?: Record<string, unknown>;
}

export interface ChatRequest {
  messages: Message[];
  context_sources?: ContextSource[];
  max_context_tokens?: number;
  stream?: boolean;
}

export interface DraftRequest {
  email_id: string;
  tone?: string;
  context?: string;
}

export interface DraftResponse {
  subject: string;
  content: string;
  estimated_tokens?: number;
}

export interface AIError {
  code: string;
  message: string;
  details?: unknown;
}