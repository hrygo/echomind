import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { SearchCluster, SearchResultsSummary, ClusterType, SearchViewMode } from '@/types/search';

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

  // Search Enhancement State (NEW)
  clusters: SearchCluster[];
  clusterType: ClusterType;
  summary: SearchResultsSummary | null;
  searchViewMode: SearchViewMode;
  enableClustering: boolean;
  enableSummary: boolean;

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

  // Search Enhancement Actions (NEW)
  setClusters: (clusters: SearchCluster[]) => void;
  setClusterType: (clusterType: ClusterType) => void;
  setSummary: (summary: SearchResultsSummary | null) => void;
  setSearchViewMode: (mode: SearchViewMode) => void;
  setEnableClustering: (enable: boolean) => void;
  setEnableSummary: (enable: boolean) => void;

  addMessage: (message: CopilotMessage) => void;
  setMessages: (messages: CopilotMessage[]) => void;
  setIsChatting: (isChatting: boolean) => void;
  setActiveContextId: (id: string | null) => void;
  clearMessages: () => void; // 清空会话历史
  reset: () => void;
}

export const useCopilotStore = create<CopilotState>()(persist(
  (set) => ({
    // Mode State
    isOpen: false,
    mode: 'idle' as const,

    // Input State
    query: '',

    // Search State
    searchResults: [],
    isSearching: false,

    // Search Enhancement State
    clusters: [],
    clusterType: 'sender' as ClusterType,
    summary: null,
    searchViewMode: 'all' as SearchViewMode,
    enableClustering: false,
    enableSummary: false,

    // Chat State
    messages: [],
    isChatting: false,

    // Context State
    activeContextId: null,

    // Basic Actions
    setIsOpen: (isOpen) => set({ isOpen }),
    setMode: (mode) => set({ mode }),
    setQuery: (query) => set({ query }),
    setSearchResults: (searchResults) => set({ searchResults }),
    setIsSearching: (isSearching) => set({ isSearching }),

    // Search Enhancement Actions
    setClusters: (clusters) => set({ clusters }),
    setClusterType: (clusterType) => set({ clusterType }),
    setSummary: (summary) => set({ summary }),
    setSearchViewMode: (searchViewMode) => set({ searchViewMode }),
    setEnableClustering: (enableClustering) => set({ enableClustering }),
    setEnableSummary: (enableSummary) => set({ enableSummary }),

    // Chat Actions
    addMessage: (message) => set((state) => {
      const newMessages = [...state.messages, message];
      // 限制为最近 50 轮对话（100 条消息）
      const MAX_MESSAGES = 100;
      if (newMessages.length > MAX_MESSAGES) {
        return { messages: newMessages.slice(newMessages.length - MAX_MESSAGES) };
      }
      return { messages: newMessages };
    }),
    setMessages: (messages) => set({ messages }),
    setIsChatting: (isChatting) => set({ isChatting }),

    // Context Actions
    setActiveContextId: (activeContextId) => set({ activeContextId }),

    // Clear Messages
    clearMessages: () => set({ messages: [] }),

    // Reset (保留历史会话，但清理其他状态)
    reset: () => set((state) => ({
      mode: 'idle',
      query: '',
      searchResults: [],
      isSearching: false,
      clusters: [],
      clusterType: 'sender',
      summary: null,
      searchViewMode: 'all',
      enableClustering: false,
      enableSummary: false,
      // 过滤掉"思考中"等未完成的 AI 消息
      messages: state.messages.filter(m => {
        if (m.role === 'assistant') {
          const thinkingTexts = ['思考中', 'Thinking'];
          return !thinkingTexts.some(text => m.content === text);
        }
        return true;
      }),
      isChatting: false,
    })),
  }),
  {
    name: 'copilot-storage',
    partialize: (state) => ({
      // 过滤掉"思考中"等未完成的 AI 消息
      messages: state.messages.filter(m => {
        if (m.role === 'assistant') {
          // 排除包含特定占位文本的消息
          const thinkingTexts = ['思考中', 'Thinking', '...'];
          return !thinkingTexts.some(text => m.content.includes(text) && m.content.length < 20);
        }
        return true;
      }),
      // 持久化搜索增强设置
      enableClustering: state.enableClustering,
      enableSummary: state.enableSummary,
      clusterType: state.clusterType,
    }),
  }
));
