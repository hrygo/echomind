import { create } from 'zustand';
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
  reset: () => void;
}

export const useCopilotStore = create<CopilotState>((set) => ({
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
  addMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  setMessages: (messages) => set({ messages }),
  setIsChatting: (isChatting) => set({ isChatting }),
  
  // Context Actions
  setActiveContextId: (activeContextId) => set({ activeContextId }),
  
  // Reset
  reset: () => set({
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
    messages: [],
    isChatting: false,
  }),
}));
