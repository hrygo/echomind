import { create } from 'zustand';
import { Email } from '@/lib/api/emails';

export interface WidgetData {
  type: string; // e.g., 'task_card', 'search_result_card'
  data: Record<string, any>; // Arbitrary data for the widget
}

interface Message {
  role: 'user' | 'assistant' | 'system';
  content?: string; // Make content optional if widget is present
  widget?: WidgetData; // Optional widget data
}

interface ChatState {
  isOpen: boolean;
  messages: Message[];
  isLoading: boolean;
  activeContextEmails: Email[]; // New state for emails passed as context
  toggleOpen: () => void;
  setOpen: (open: boolean) => void;
  addMessage: (message: Message) => void;
  setLoading: (loading: boolean) => void;
  clearMessages: () => void;
  updateLastMessage: (content: string) => void;
  setActiveContextEmails: (emails: Email[]) => void;
  clearActiveContextEmails: () => void;
}

export const useChatStore = create<ChatState>((set) => ({
  isOpen: false,
  messages: [],
  isLoading: false,
  activeContextEmails: [], // Initialize as empty array
  toggleOpen: () => set((state) => ({ isOpen: !state.isOpen })),
  setOpen: (open) => set({ isOpen: open }),
  addMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  setLoading: (loading) => set({ isLoading: loading }),
  clearMessages: () => set({ messages: [] }),
  updateLastMessage: (content) =>
    set((state) => {
      const messages = [...state.messages];
      if (messages.length > 0) {
        messages[messages.length - 1] = {
          ...messages[messages.length - 1],
          content: content,
        };
      }
      return { messages };
    }),
  setActiveContextEmails: (emails) => set({ activeContextEmails: emails }),
  clearActiveContextEmails: () => set({ activeContextEmails: [] }),
}));
