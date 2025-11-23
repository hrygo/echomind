'use client';

import { useState, useRef, useEffect } from 'react';
import { Send, Sparkles, Loader2 } from 'lucide-react';
import { useChatStore } from '@/lib/store/chat';
import { useAuthStore } from '@/store/auth';
import { cn } from '@/lib/utils';
import ReactMarkdown from 'react-markdown';
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/Sheet';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { TaskWidget } from '@/components/widgets/TaskWidget';
import { SearchResultWidget } from '@/components/widgets/SearchResultWidget';


export function ChatSidebar() {
    const { isOpen, setOpen, messages, addMessage, isLoading, setLoading, updateLastMessage, activeContextEmails, clearActiveContextEmails } = useChatStore();
    const [input, setInput] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const { t } = useLanguage();

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    // Handle activeContextEmails when they are set
    useEffect(() => {
        if (activeContextEmails.length > 0) {
            const contextSummary = activeContextEmails.map(email => 
                `Subject: ${email.Subject}\nSender: ${email.Sender}\nSnippet: ${email.Snippet}`
            ).join('\n\n');

            const systemMessageContent = t('chat.contextLoaded').replace('{count}', activeContextEmails.length.toString()) + '\n\n' + contextSummary;
            addMessage({ role: 'system', content: systemMessageContent });
            // Clear the context emails after they have been added to the chat as a system message
            clearActiveContextEmails();
            setOpen(true); // Ensure chat is open when context is loaded
        }
    }, [activeContextEmails, addMessage, clearActiveContextEmails, setOpen, t]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!input.trim() || isLoading) return;

        const userMessageContent = input.trim();
        setInput('');

        const conversationMessages = [...messages, { role: 'user', content: userMessageContent }];
        
        // If there's active context, prepend it as a system message for the API call
        // Note: The context has already been added to `messages` by the useEffect above.
        // So, we just need to ensure the `apiMessages` correctly represents the full conversation history.
        const apiMessages = conversationMessages;

        setLoading(true);

        // Do NOT add a placeholder assistant message here. It will be added when the first content chunk arrives.

        try {
            // Get token from auth store
            const token = useAuthStore.getState().token;
            const baseUrl = process.env.NEXT_PUBLIC_API_URL || '/api/v1';

            const response = await fetch(`${baseUrl}/chat/completions`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'text/event-stream',
                },
                body: JSON.stringify({ messages: apiMessages }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Chat API Error:', response.status, response.statusText, errorText);
                throw new Error(`Failed to send message: ${response.status} ${response.statusText}`);
            }

            const reader = response.body?.getReader();
            const decoder = new TextDecoder();
            let assistantMessage = '';
            let isFirstMessageChunk = true; // Track if this is the first chunk for a new assistant response

            if (reader) {
                let buffer = '';
                while (true) {
                    const { done, value } = await reader.read();
                    if (done) break;

                    buffer += decoder.decode(value, { stream: true });
                    const lines = buffer.split('\n');

                    buffer = lines.pop() || ''; // Keep the last line if it's incomplete

                    for (const line of lines) {
                        if (line.trim() === '') continue;
                        if (line.startsWith('data: ')) {
                            const data = line.slice(6);
                            if (data === '[DONE]') break;

                            try {
                                const parsed = JSON.parse(data);
                                if (parsed.error) {
                                    console.error('Chat API returned error:', parsed.error);
                                    updateLastMessage(`Error: ${parsed.error}`);
                                    break;
                                }
                                if (parsed.choices && parsed.choices[0] && parsed.choices[0].delta) {
                                    const delta = parsed.choices[0].delta;

                                    if (delta.widget) {
                                        // If a widget is received, add it as a new, complete message
                                        addMessage({ role: 'assistant', widget: delta.widget });
                                        assistantMessage = ''; // Reset assistantMessage
                                        isFirstMessageChunk = true; // Next content should start a new message
                                    } else if (delta.content) {
                                        if (isFirstMessageChunk) {
                                            addMessage({ role: 'assistant', content: delta.content });
                                            assistantMessage = delta.content;
                                            isFirstMessageChunk = false;
                                        } else {
                                            assistantMessage += delta.content;
                                            updateLastMessage(assistantMessage);
                                        }
                                    }
                                }
                            } catch {
                                console.warn('JSON parse failed, fallback:', data);
                                // Fallback to raw content if JSON parse fails, append to current assistantMessage
                                if (isFirstMessageChunk) {
                                    addMessage({ role: 'assistant', content: data });
                                    assistantMessage = data;
                                    isFirstMessageChunk = false;
                                } else {
                                    assistantMessage += data;
                                    updateLastMessage(assistantMessage);
                                }
                            }
                        }
                    }
                }

                // Process any remaining buffer at the end of the stream
                if (buffer.trim()) {
                    // This part is for any final content not ending with \n in the last chunk
                    const data = buffer.slice(6); // Assuming it starts with 'data: '
                    if (isFirstMessageChunk) {
                        addMessage({ role: 'assistant', content: data });
                        assistantMessage = data;
                    } else {
                        assistantMessage += data;
                        updateLastMessage(assistantMessage);
                    }
                }

            }
        } catch (error) {
            console.error('Chat error:', error);
            // If an error occurs before any message is added, add it as the first message.
            if (messages.length === apiMessages.length) { // No assistant message added yet
                addMessage({ role: 'assistant', content: 'Sorry, something went wrong. Please try again.' });
            } else {
                updateLastMessage('Sorry, something went wrong. Please try again.');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Sheet open={isOpen} onOpenChange={setOpen}>
            <SheetContent side="right" className="w-full sm:max-w-md flex flex-col p-0">
                {/* Header */}
                <SheetHeader className="p-4 border-b border-gray-200 flex flex-row items-center justify-between bg-gradient-to-r from-indigo-50 to-white">
                    <div className="flex items-center gap-2 text-indigo-700">
                        <Sparkles className="w-5 h-5" />
                        <SheetTitle className="font-semibold text-lg">EchoMind Copilot</SheetTitle>
                    </div>
                </SheetHeader>

                {/* Messages */}
                <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
                    {messages.length === 0 && (
                        <div className="text-center text-gray-500 mt-10">
                            <Sparkles className="w-12 h-12 mx-auto mb-3 text-indigo-200" />
                            <p>How can I help you today?</p>
                            <p className="text-sm mt-2">Try asking about your recent emails.</p>
                        </div>
                    )}

                    {messages.map((msg, idx) => (
                        <div
                            key={idx}
                            className={cn(
                                "flex w-full",
                                msg.role === 'user' ? "justify-end" : "justify-start"
                            )}
                        >
                            <div
                                className={cn(
                                    "max-w-[85%] rounded-2xl px-4 py-2.5 text-sm shadow-sm",
                                    msg.role === 'user'
                                        ? "bg-indigo-600 text-white rounded-br-none"
                                        : "bg-white text-gray-800 border border-gray-100 rounded-bl-none"
                                )}
                            >
                                <div className="prose prose-sm max-w-none dark:prose-invert">
                                    {msg.content && <ReactMarkdown>{msg.content}</ReactMarkdown>}
                                    {msg.widget && (
                                        msg.widget.type === 'task_card' ? <TaskWidget data={msg.widget.data as any} /> :
                                        msg.widget.type === 'search_result_card' ? <SearchResultWidget data={msg.widget.data as any} /> :
                                        null
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                    {isLoading && messages.length > 0 && messages[messages.length - 1].role === 'user' && (
                        <div className="flex justify-start">
                            <div className="bg-white p-3 rounded-2xl rounded-bl-none border border-gray-100 shadow-sm">
                                <Loader2 className="w-4 h-4 animate-spin text-indigo-500" />
                            </div>
                        </div>
                    )}
                    <div ref={messagesEndRef} />
                </div>

                {/* Input */}
                <div className="p-4 bg-white border-t border-gray-200">
                    <form onSubmit={handleSubmit} className="relative">
                        <input
                            type="text"
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder="Ask anything..."
                            className="w-full pl-4 pr-12 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-gray-900 placeholder:text-gray-500"
                            disabled={isLoading}
                        />
                        <button
                            type="submit"
                            disabled={!input.trim() || isLoading}
                            className="absolute right-2 top-1/2 -translate-y-1/2 p-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors shadow-sm"
                        >
                            {isLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
                        </button>
                    </form>
                </div>
            </SheetContent>
        </Sheet>
    );
}
