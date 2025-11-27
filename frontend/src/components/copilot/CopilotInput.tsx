'use client';

import React, { KeyboardEvent } from 'react';
import { Search, Sparkles, X, Settings } from 'lucide-react';
import { useCopilotStore } from '@/store'; // Updated import
import { cn } from '@/lib/utils'; // Assuming this exists, standard in shadcn/ui projects
import { useAuthStore } from '@/store/auth'; // Import useAuthStore
import { useLanguage } from '@/lib/i18n/LanguageContext';
import type { SearchResponse } from '@/types/search';

interface CopilotInputProps {
  showSettings?: boolean;
  onToggleSettings?: () => void;
  onCloseSettings?: () => void;  // 新增：关闭设置面板的回调
}

export function CopilotInput({ showSettings, onToggleSettings, onCloseSettings }: CopilotInputProps = {}) {
  const { t } = useLanguage();
  const { 
    query, 
    setQuery, 
    mode, 
    setMode, 
    setIsSearching, 
    setSearchResults, 
    setIsChatting, 
    addMessage,
    activeContextId,
    // Search enhancement
    enableClustering,
    enableSummary,
    clusterType,
    setClusters,
    setSummary,
  } = useCopilotStore();

  // 判断设置按钮是否应该被禁用
  const isSettingsDisabled = mode !== 'idle';

  const handleKeyDown = async (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && query.trim()) {
      e.preventDefault();
      
      const isExplicitChat = query.startsWith('?') || query.startsWith('？');
      const targetMode = isExplicitChat || mode === 'chat' ? 'chat' : 'search';

      if (targetMode === 'search') {
        handleSearch();
      } else {
        handleChat();
      }
    }
  };

  const handleSearch = async () => {
    setMode('search');
    setIsSearching(true);
    setSearchResults([]); // Clear previous
    setClusters([]); // Clear previous clusters
    setSummary(null); // Clear previous summary

    try {
      const token = useAuthStore.getState().token; // Get token from AuthStore
      const params = new URLSearchParams({
          q: query,
          limit: '10'
      });
      if (activeContextId) {
          params.append('context_id', activeContextId);
      }
      
      // Add search enhancement parameters
      if (enableClustering) {
          params.append('enable_clustering', 'true');
          params.append('cluster_type', clusterType);
      }
      if (enableSummary) {
          params.append('enable_summary', 'true');
      }

      const response = await fetch(`/api/v1/search?${params.toString()}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
      });
      
      if (!response.ok) {
          const errorText = await response.text();
          console.error('Search API Error:', response.status, errorText);
          throw new Error(`Search failed: ${response.status} ${errorText}`);
      }
      
      const data: SearchResponse = await response.json();
      setSearchResults(data.results || []);
      
      // Handle enhanced search data
      if (data.clusters) {
          setClusters(data.clusters);
      }
      if (data.summary) {
          setSummary(data.summary);
      }
    } catch (error) {
      console.error('Search error:', error);
    } finally {
      setIsSearching(false);
    }
  };

  const handleChat = async () => {
    setMode('chat');
    const userMsg = { role: 'user' as const, content: query };
    addMessage(userMsg);
    setQuery('');
    setIsChatting(true);
  };

  return (
    <div className="relative w-full max-w-2xl mx-auto">
      <div className={cn(
        "flex items-center w-full px-4 py-3 bg-white border rounded-xl shadow-sm transition-all duration-200",
        "focus-within:shadow-md focus-within:border-blue-500/50 focus-within:ring-2 focus-within:ring-blue-100",
        mode !== 'idle' ? "rounded-b-none border-b-0" : ""
      )}>
        {mode === 'chat' ? (
          <Sparkles className="w-5 h-5 text-indigo-500 mr-3 animate-pulse" />
        ) : (
          <Search className="w-5 h-5 text-slate-400 mr-3" />
        )}
        
        <input
          type="text"
          className="flex-1 bg-transparent outline-none text-slate-700 placeholder:text-slate-400"
          placeholder={mode === 'chat' ? t('copilot.placeholderChat') : t('copilot.placeholderSearch')}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          onMouseEnter={() => {
            // 当鼠标移入输入框时，关闭设置面板
            if (showSettings && onCloseSettings) {
              onCloseSettings();
            }
          }}
          onFocus={() => {
            // 当输入框获得焦点时，关闭设置面板
            if (showSettings && onCloseSettings) {
              onCloseSettings();
            }
          }}
        />

        {query && (
          <button 
            onClick={() => setQuery('')}
            className="p-1 hover:bg-slate-100 rounded-full text-slate-400 mr-2"
          >
            <X className="w-4 h-4" />
          </button>
        )}

        <div className="flex gap-1">
           {/* Settings Button */}
           {onToggleSettings && (
             <button 
              onClick={onToggleSettings}
              disabled={isSettingsDisabled}
              className={cn(
                  "p-2 rounded-lg transition-colors",
                  isSettingsDisabled 
                    ? "opacity-40 cursor-not-allowed text-slate-300" 
                    : showSettings 
                      ? "bg-blue-50 text-blue-600" 
                      : "hover:bg-slate-50 text-slate-400"
              )}
              title={isSettingsDisabled ? t('copilot.searchEnhancement.disabledHint') || '搜索或聊天时不可用' : t('copilot.searchEnhancement.title')}
             >
               <Settings className="w-4 h-4" />
             </button>
           )}
           
           {/* Chat Mode Toggle */}
           <button 
            onClick={() => setMode('chat')}
            className={cn(
                "p-2 rounded-lg transition-colors",
                mode === 'chat' ? "bg-indigo-50 text-indigo-600" : "hover:bg-slate-50 text-slate-400"
            )}
            title={t('copilot.switchToChat')}
           >
             <Sparkles className="w-4 h-4" />
           </button>
        </div>
      </div>
    </div>
  );
}
