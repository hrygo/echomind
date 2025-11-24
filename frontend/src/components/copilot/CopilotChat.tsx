'use client';

import React, { useEffect, useRef } from 'react';
import { User, Bot, Sparkles } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { useCopilotStore, CopilotMessage } from '@/store';
import { cn } from '@/lib/utils';
import { useAuthStore } from '@/store/auth';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { TaskListWidget, TaskListWidgetProps } from '@/components/widgets/TaskListWidget';
import { EmailDraftWidget, EmailDraftWidgetProps } from '@/components/widgets/EmailDraftWidget';
import { CalendarEventWidget, CalendarEventWidgetProps } from '@/components/widgets/CalendarEventWidget';

function MessageBubble({ message }: { message: CopilotMessage }) {
  const isUser = message.role === 'user';

  return (
    <div className={cn(
      "flex gap-3 mb-4",
      isUser ? "flex-row-reverse" : "flex-row"
    )}>
      <div className={cn(
        "w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0",
        isUser ? "bg-slate-200 text-slate-600" : "bg-indigo-600 text-white"
      )}>
        {isUser ? <User className="w-5 h-5" /> : <Bot className="w-5 h-5" />}
      </div>

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
    </div>
  );
}

export function CopilotChat() {
  const { t } = useLanguage();
  const { messages, isChatting, addMessage, searchResults } = useCopilotStore();
  const bottomRef = useRef<HTMLDivElement>(null);
  const hasInitialized = useRef(false);

  // Auto-scroll to bottom
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Handle Initial Chat Trigger (Real API Call)
  useEffect(() => {
    const initChat = async () => {
      if (isChatting && !hasInitialized.current && messages.length > 0) {
        const lastMsg = messages[messages.length - 1];
        // Only trigger if the last message is from user and we haven't started responding yet
        if (lastMsg.role === 'user') {
          hasInitialized.current = true;

          // Add placeholder for Assistant response
          addMessage({ role: 'assistant', content: t('copilot.thinking') });

          try {
            const token = useAuthStore.getState().token;
            const response = await fetch('/api/v1/chat/completions', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
              },
              body: JSON.stringify({
                messages: messages.map(m => ({ role: m.role, content: m.content })),
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
            let assistantMessage = "";
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

                      // Update the last message (assistant) with new content
                      useCopilotStore.setState(state => {
                        const newMessages = [...state.messages];
                        newMessages[newMessages.length - 1] = {
                          role: 'assistant',
                          content: displayContent, // Show content without widget XML
                          widget: detectedWidget
                        };
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
            useCopilotStore.setState(state => {
              const newMessages = [...state.messages];
              newMessages[newMessages.length - 1] = {
                role: 'assistant',
                content: t('copilot.errorResponse')
              };
              return { messages: newMessages };
            });
          } finally {
            hasInitialized.current = false;
          }
        }
      }
    };
    initChat();
  }, [isChatting, messages, searchResults]); // eslint-disable-line react-hooks/exhaustive-deps


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
      <div className="p-3 border-t border-slate-100 bg-slate-50/50 rounded-b-xl flex gap-2 text-xs text-slate-500">
        <span>{t('copilot.contextAttached').replace('{count}', searchResults.length.toString())}</span>
      </div>
    </div>
  );
}
