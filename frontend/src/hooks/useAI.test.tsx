/**
 * Unit Test: AI Hooks
 * 测试自定义 AI Hooks 功能
 */

import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useStreamChat, useAIDraft, useAIReply } from '@/hooks/useAI';
import * as chatClient from '@/lib/ai/chat-client';
import * as draftClient from '@/lib/ai/draft-client';
import type { ReactNode } from 'react';

// Mock AI 客户端
jest.mock('@/lib/ai/chat-client');
jest.mock('@/lib/ai/draft-client');

const mockedChatClient = chatClient as jest.Mocked<typeof chatClient>;
const mockedDraftClient = draftClient as jest.Mocked<typeof draftClient>;

// 创建测试用的 QueryClient wrapper
function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  const Wrapper = ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );

  return Wrapper;
}

describe('useStreamChat Hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should initialize with correct default state', () => {
    const { result } = renderHook(() => useStreamChat(), {
      wrapper: createWrapper(),
    });

    expect(result.current.messages).toEqual([]);
    expect(result.current.isStreaming).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('should handle successful chat streaming', async () => {
    mockedChatClient.streamChat.mockImplementation(async (options) => {
      // 模拟流式响应
      await new Promise(resolve => setTimeout(resolve, 10));
      options.onChunk?.({
        id: '1',
        choices: [{ delta: { content: 'Hello' } }]
      });
      options.onChunk?.({
        id: '1',
        choices: [{ delta: { content: ' world' } }]
      });
      options.onComplete?.();
    });

    const { result } = renderHook(() => useStreamChat(), {
      wrapper: createWrapper(),
    });

    // 发送消息
    await result.current.sendMessage({
      messages: [{ role: 'user', content: 'Test message', id: 'u1' }]
    });

    await waitFor(() => {
      expect(result.current.messages.length).toBeGreaterThan(0);
    });

    // 验证消息内容
    const userMessage = result.current.messages.find((m) => m.role === 'user');
    expect(userMessage?.content).toBe('Test message');

    await waitFor(() => {
      expect(result.current.isStreaming).toBe(false);
    });

    const assistantMessage = result.current.messages.find((m) => m.role === 'assistant');
    expect(assistantMessage?.content).toContain('Hello');
  });

  it('should handle streaming errors', async () => {
    const errorMessage = 'Streaming failed';
    mockedChatClient.streamChat.mockImplementation(async (options) => {
      options.onError?.(new Error(errorMessage));
    });

    const { result } = renderHook(() => useStreamChat(), {
      wrapper: createWrapper(),
    });

    result.current.sendMessage({
      messages: [{ role: 'user', content: 'Test message', id: 'u1' }]
    });

    await waitFor(() => {
      expect(result.current.error).toBeTruthy();
    });

    expect(result.current.error?.message).toBe(errorMessage);
    expect(result.current.isStreaming).toBe(false);
  });

  it('should clear messages', () => {
    const { result } = renderHook(() => useStreamChat(), {
      wrapper: createWrapper(),
    });

    // 清除消息
    result.current.clearMessages();

    expect(result.current.messages).toEqual([]);
  });
});

describe('useAIDraft Hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should initialize mutation correctly', () => {
    const { result } = renderHook(() => useAIDraft(), {
      wrapper: createWrapper(),
    });

    expect(result.current.isPending).toBe(false);
    expect(result.current.data).toBeUndefined();
    expect(result.current.error).toBeNull();
  });

  it('should generate draft successfully', async () => {
    const mockDraft = {
      content: 'Generated draft content',
      subject: 'Meeting Request',
    };

    mockedDraftClient.generateDraft.mockResolvedValue(mockDraft);

    const { result } = renderHook(() => useAIDraft(), {
      wrapper: createWrapper(),
    });

    // 生成草稿
    result.current.mutate({
      emailId: 'email-123',
      context: 'Project discussion',
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockDraft);
    expect(mockedDraftClient.generateDraft).toHaveBeenCalledWith({
      email_id: 'email-123',
      context: 'Project discussion',
    });
  });

  it('should handle draft generation errors', async () => {
    const errorMessage = 'Draft generation failed';
    mockedDraftClient.generateDraft.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useAIDraft(), {
      wrapper: createWrapper(),
    });

    result.current.mutate({
      emailId: 'email-123',
    });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toBeTruthy();
  });

  it('should support different tones', async () => {
    const mockDraft = { content: 'Formal draft', subject: 'Test' };
    mockedDraftClient.generateDraft.mockResolvedValue(mockDraft);

    const { result } = renderHook(() => useAIDraft(), {
      wrapper: createWrapper(),
    });

    result.current.mutate({
      emailId: 'email-123',
      tone: 'formal',
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(mockedDraftClient.generateDraft).toHaveBeenCalledWith({
      email_id: 'email-123',
      tone: 'formal',
    });
  });
});

describe('useAIReply Hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should generate reply successfully', async () => {
    const mockReply = {
      content: 'Thank you for your email. I will review it shortly.',
      subject: 'Re: Meeting Request',
    };

    mockedDraftClient.generateReply.mockResolvedValue(mockReply);

    const { result } = renderHook(() => useAIReply(), {
      wrapper: createWrapper(),
    });

    result.current.mutate({
      emailId: 'email-123',
      tone: 'professional',
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockReply);
    expect(mockedDraftClient.generateReply).toHaveBeenCalledWith(
      'email-123',
      'professional',
      undefined
    );
  });

  it('should pass custom context to reply generation', async () => {
    const mockReply = { content: 'Custom reply', subject: 'Re: Test' };
    mockedDraftClient.generateReply.mockResolvedValue(mockReply);

    const { result } = renderHook(() => useAIReply(), {
      wrapper: createWrapper(),
    });

    const customContext = 'Please mention the meeting scheduled for next week';

    result.current.mutate({
      emailId: 'email-123',
      context: customContext,
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(mockedDraftClient.generateReply).toHaveBeenCalledWith(
      'email-123',
      undefined,
      customContext
    );
  });

  it('should handle reply generation errors', async () => {
    mockedDraftClient.generateReply.mockRejectedValue(new Error('Reply generation failed'));

    const { result } = renderHook(() => useAIReply(), {
      wrapper: createWrapper(),
    });

    result.current.mutate({ emailId: 'email-123' });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toBeTruthy();
  });
});
