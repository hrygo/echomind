import { create } from 'zustand';

export interface CopilotMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
  widget?: {
    type: string; // e.g., 'task_card', 'search_result_card'
    data: Record<string, unknown>;
  };
}

export interface SearchResult {
  email_id: string;
  subject: string;
  snippet: string;
  sender: string;
  date: string;
  score: number;
}

interface CopilotState {
  // Mode State
  isOpen: boolean;
  mode: 'search' | 'chat' | 'idle';
  
  // Input State
  query: string;
  
  // Search State
  searchResults: SearchResult[];
  isSearching: boolean;
  
  // Chat State
  messages: CopilotMessage[];
  isChatting: boolean;
  
  // Context State
  activeContextId: string | null;

  // Actions
  setIsOpen: (isOpen: boolean) => void;
  setMode: (mode: 'search' | 'chat' | 'idle') => void;
  setQuery: (query: string) => void;
  setSearchResults: (results: SearchResult[]) => void;
  setIsSearching: (isSearching: boolean) => void;
  addMessage: (message: CopilotMessage) => void;
  setMessages: (messages: CopilotMessage[]) => void;
  setIsChatting: (isChatting: boolean) => void;
  setActiveContextId: (id: string | null) => void;
  reset: () => void;
}

export const useCopilotStore = create<CopilotState>((set) => ({
  isOpen: false,
  mode: 'idle',
  query: '',
  searchResults: [],
  isSearching: false,
  messages: [],
  isChatting: false,
  activeContextId: null,

  setIsOpen: (isOpen) => set({ isOpen }),
  setMode: (mode) => set({ mode }),
  setQuery: (query) => set({ query }),
  setSearchResults: (searchResults) => set({ searchResults }),
  setIsSearching: (isSearching) => set({ isSearching }),
  addMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  setMessages: (messages) => set({ messages }),
  setIsChatting: (isChatting) => set({ isChatting }),
  setActiveContextId: (activeContextId) => set({ activeContextId }),
  reset: () => set({
    mode: 'idle',
    query: '',
    searchResults: [],
    isSearching: false,
    messages: [], // Might want to persist chat history in local storage later
    isChatting: false,
  }),
}));
