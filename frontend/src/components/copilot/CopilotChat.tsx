'use client';

import React, { useEffect, useRef } from 'react';
import { User, Bot, Sparkles } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { useCopilotStore, CopilotMessage } from '@/store';
import { useAuthStore } from '@/store/auth';
import { cn } from '@/lib/utils';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { TaskListWidget, TaskListWidgetProps } from '@/components/widgets/TaskListWidget';
import { EmailDraftWidget, EmailDraftWidgetProps } from '@/components/widgets/EmailDraftWidget';
import { CalendarEventWidget, CalendarEventWidgetProps } from '@/components/widgets/CalendarEventWidget';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';

import { ThinkingIndicator } from '@/components/ui/ThinkingIndicator';

function MessageBubble({ message }: { message: CopilotMessage }) {
  const isUser = message.role === 'user';
  const user = useAuthStore(state => state.user);
  const { t } = useLanguage();

  // Check if this is the thinking message
  const isThinking = message.role === 'assistant' && message.content === t('copilot.thinking');

  return (
    <div className={cn(
      "flex gap-3 mb-4",
      isUser ? "flex-row-reverse" : "flex-row"
    )}>
      {isUser ? (
        <Avatar className="w-8 h-8 flex-shrink-0">
          <AvatarImage src={user?.avatar_url} alt={user?.name || user?.email} />
          <AvatarFallback className="bg-slate-200 text-slate-600">
            {user?.name?.[0]?.toUpperCase() || user?.email?.[0]?.toUpperCase() || 'U'}
          </AvatarFallback>
        </Avatar>
      ) : (
        <div className="w-8 h-8 rounded-full bg-indigo-600 text-white flex items-center justify-center flex-shrink-0">
          <Bot className="w-5 h-5" />
        </div>
      )}

      {isThinking ? (
        <div className="flex items-center py-2 px-1">
          <ThinkingIndicator text={t('copilot.thinking')} />
        </div>
      ) : (
        <div className={cn(
          "max-w-[85%] rounded-2xl px-4 py-2.5 text-sm leading-relaxed",
          isUser
            ? "bg-slate-100 text-slate-800 rounded-tr-sm"
            : "bg-indigo-50/50 text-slate-800 border border-indigo-100 rounded-tl-sm shadow-sm"
        )}>
          {isUser ? (
            message.content
          ) : (
            <div className="prose prose-sm prose-indigo max-w-none">
              <ReactMarkdown>{message.content}</ReactMarkdown>
            </div>
          )}

          {/* Render Widget if present */}
          {message.widget && (
            <div className="mt-3">
              <div className="flex items-center gap-2 text-xs font-semibold text-slate-500 uppercase mb-2">
                <Sparkles className="w-3 h-3 text-indigo-500" />
                {message.widget.type.replace('_', ' ')}
              </div>

              {message.widget.type === 'task_list' && (
                <TaskListWidget data={message.widget.data as unknown as TaskListWidgetProps['data']} />
              )}
              {message.widget.type === 'email_draft' && (
                <EmailDraftWidget data={message.widget.data as unknown as EmailDraftWidgetProps['data']} />
              )}
              {message.widget.type === 'calendar_event' && (
                <CalendarEventWidget data={message.widget.data as unknown as CalendarEventWidgetProps['data']} />
              )}

              {/* Fallback for unknown widgets */}
              {!['task_list', 'email_draft', 'calendar_event'].includes(message.widget.type) && (
                <pre className="text-xs bg-slate-50 p-2 rounded overflow-x-auto text-slate-600">
                  {JSON.stringify(message.widget.data, null, 2)}
                </pre>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export function CopilotChat() {
  const { t } = useLanguage();
  const { messages, isChatting, addMessage, clearMessages, searchResults } = useCopilotStore();
  const bottomRef = useRef<HTMLDivElement>(null);
  const processingRef = useRef(false); // 标记是否正在处理请求
  const pendingTimeoutRef = useRef<NodeJS.Timeout | null>(null); // 延迟处理定时器
  const hasPendingMessagesRef = useRef(false); // 标记是否有待处理的消息

  // Auto-scroll to bottom
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Handle Chat - 合并连续的用户消息
  useEffect(() => {
    const initChat = async () => {
      // 查找最后一条 assistant 消息的位置
      let lastAssistantIndex = -1;
      for (let i = messages.length - 1; i >= 0; i--) {
        if (messages[i].role === 'assistant') {
          lastAssistantIndex = i;
          break;
        }
      }

      // 获取所有未处理的用户消息（在最后一条 assistant 消息之后的）
      const pendingUserMessages = messages.slice(lastAssistantIndex + 1).filter(m => m.role === 'user');

      // 如果没有待处理的用户消息，退出
      if (pendingUserMessages.length === 0) {
        // 清除定时器
        if (pendingTimeoutRef.current) {
          clearTimeout(pendingTimeoutRef.current);
          pendingTimeoutRef.current = null;
        }
        return;
      }

      // 如果正在处理请求，标记有待处理消息，等待处理完成后再次检查
      if (processingRef.current) {
        hasPendingMessagesRef.current = true;
        return;
      }

      // 清除之前的定时器
      if (pendingTimeoutRef.current) {
        clearTimeout(pendingTimeoutRef.current);
      }

      // 设置延迟处理，给用户输入更多问题的时间（1秒）
      pendingTimeoutRef.current = setTimeout(async () => {
        // 再次检查是否正在处理（双重保险）
        if (processingRef.current) {
          return;
        }

        // 标记正在处理
        processingRef.current = true;
        pendingTimeoutRef.current = null;

        // 获取当前最新的消息状态
        const currentMessages = useCopilotStore.getState().messages;

        // 验证最后一条消息必须是用户消息
        if (currentMessages.length > 0 && currentMessages[currentMessages.length - 1].role !== 'user') {
          processingRef.current = false;
          return;
        }

        // Add placeholder for Assistant response
        addMessage({ role: 'assistant', content: t('copilot.thinking') });

        // 记录当前 assistant 消息的索引（刚添加的"思考中..."）
        const assistantMessageIndex = currentMessages.length;

        try {
          const token = useAuthStore.getState().token;

          // 使用最新的消息状态准备对话历史（不包含刚添加的"思考中..."占位消息）
          const conversationHistory = currentMessages.map(m => ({
            role: m.role,
            content: m.content
          }));

          const response = await fetch('/api/v1/chat/completions', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
              messages: conversationHistory,
              context_ref_ids: searchResults.map(r => r.email_id)
            }),
          });

          if (!response.ok) {
            const errorText = await response.text();
            console.error('Chat API Error:', response.status, errorText);
            throw new Error(`Chat failed: ${response.status} ${errorText}`);
          }
          if (!response.body) throw new Error('No response body');

          const reader = response.body.getReader();
          const decoder = new TextDecoder();
          let isFirstChunk = true;
          let fullBuffer = ""; // Buffer to hold full response for parsing

          while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value);
            const lines = chunk.split('\n');

            for (const line of lines) {
              if (line.startsWith('data: ')) {
                const data = line.slice(6);
                if (data.trim() === '[DONE]') break;

                try {
                  const parsed = JSON.parse(data);
                  if (parsed.error) {
                    throw new Error(parsed.error);
                  }

                  if (parsed.choices && parsed.choices[0].delta.content) {
                    const contentChunk = parsed.choices[0].delta.content;
                    fullBuffer += contentChunk;

                    // Widget Parsing Logic
                    let displayContent = fullBuffer;
                    let detectedWidget = undefined;

                    // Regex to find <widget type="...">...</widget>
                    // We use [\s\S]*? to match across newlines non-greedily
                    const widgetRegex = /<widget type="([^"]+)">([\s\S]*?)<\/widget>/;
                    const match = fullBuffer.match(widgetRegex);

                    if (match) {
                      const [fullMatch, type, jsonStr] = match;
                      try {
                        const widgetData = JSON.parse(jsonStr);
                        detectedWidget = {
                          type: type,
                          data: widgetData
                        };
                        // Remove widget tag from display content
                        displayContent = fullBuffer.replace(fullMatch, '').trim();
                      } catch (e) {
                        console.warn("Failed to parse widget JSON:", e);
                        // If JSON is incomplete (streaming), we might fail to parse.
                        // In a real robust impl, we'd wait for the closing tag.
                        // Since we match </code> which implies we have the full block,
                        // failure here means invalid JSON from LLM.
                      }
                    }

                    if (isFirstChunk) {
                      isFirstChunk = false;
                    }

                    // Update the specific assistant message (not the last one)
                    useCopilotStore.setState(state => {
                      const newMessages = [...state.messages];
                      // 确保索引有效且是 assistant 消息
                      if (assistantMessageIndex < newMessages.length &&
                        newMessages[assistantMessageIndex]?.role === 'assistant') {
                        newMessages[assistantMessageIndex] = {
                          role: 'assistant',
                          content: displayContent, // Show content without widget XML
                          widget: detectedWidget
                        };
                      }
                      return { messages: newMessages };
                    });
                  }
                } catch (e) {
                  console.error('Error parsing SSE chunk', e);
                }
              }
            }
          }
        } catch (error) {
          console.error('Chat Error:', error);
          // 更新特定的 assistant 消息为错误信息
          useCopilotStore.setState(state => {
            const newMessages = [...state.messages];
            if (assistantMessageIndex < newMessages.length &&
              newMessages[assistantMessageIndex]?.role === 'assistant') {
              newMessages[assistantMessageIndex] = {
                role: 'assistant',
                content: t('copilot.errorResponse')
              };
            }
            return { messages: newMessages };
          });
        } finally {
          // 处理完成，重置标志
          processingRef.current = false;

          // 检查是否有待处理的消息
          if (hasPendingMessagesRef.current) {
            hasPendingMessagesRef.current = false;
            // 给一个短暂延迟，确保 UI 更新后再检查
            setTimeout(() => {
              // 手动触发重新检查
              initChat();
            }, 50);
          }
        }
      }, 1000); // 1秒延迟，允许用户快速输入多条消息
    };

    initChat();

    // 清理函数：组件卸载时清除定时器
    return () => {
      if (pendingTimeoutRef.current) {
        clearTimeout(pendingTimeoutRef.current);
      }
    };
  }, [messages, searchResults]); // eslint-disable-line react-hooks/exhaustive-deps


  return (
    <div className="w-full max-w-2xl mx-auto bg-white border border-t-0 rounded-b-xl shadow-xl min-h-[300px] max-h-[70vh] flex flex-col">
      <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
        {messages.length === 0 && (
          <div className="h-full flex flex-col items-center justify-center text-slate-400 space-y-2 opacity-50">
            <Bot className="w-12 h-12" />
            <p>{t('copilot.welcomeMessage')}</p>
          </div>
        )}

        {messages.map((msg, idx) => (
          <MessageBubble key={idx} message={msg} />
        ))}

        <div ref={bottomRef} />
      </div>

      {/* Chat Input Area (If we want a persistent input at bottom, 
          but currently we share the top input. This area could be for quick actions.) */}
      <div className="p-3 border-t border-slate-100 bg-slate-50/50 rounded-b-xl flex gap-2 justify-between text-xs text-slate-500">
        <span>{t('copilot.contextAttached').replace('{count}', searchResults.length.toString())}</span>
        {messages.length > 0 && (
          <button
            onClick={clearMessages}
            className="text-red-500 hover:text-red-700 hover:underline transition-colors"
            title={t('copilot.clearHistoryTitle')}
          >
            {t('copilot.clearHistory')}
          </button>
        )}
      </div>
    </div>
  );
}
