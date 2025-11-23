'use client';

import React, { useState, KeyboardEvent, useEffect } from 'react';
import { Search, Sparkles, X, ArrowRight } from 'lucide-react';
import { useCopilotStore } from '@/store';
import { cn } from '@/lib/utils'; // Assuming this exists, standard in shadcn/ui projects
import { useAuthStore } from '@/store/auth'; // Import useAuthStore

export function CopilotInput() {
  const { 
    query, 
    setQuery, 
    mode, 
    setMode, 
    setIsSearching, 
    setSearchResults, 
    setIsChatting, 
    addMessage,
    activeContextId
  } = useCopilotStore();

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

    try {
      const token = useAuthStore.getState().token; // Get token from AuthStore
      const params = new URLSearchParams({
          q: query,
          limit: '10'
      });
      if (activeContextId) {
          params.append('context_id', activeContextId);
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
      
      const data = await response.json();
      setSearchResults(data.results || []);
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
          placeholder={mode === 'chat' ? "Ask Copilot anything..." : "Search emails, tasks, or contacts..."}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
        />

        {query && (
          <button 
            onClick={() => setQuery('')}
            className="p-1 hover:bg-slate-100 rounded-full text-slate-400 mr-2"
          >
            <X className="w-4 h-4" />
          </button>
        )}

        <div className="flex gap-2">
           <button 
            onClick={() => setMode('chat')}
            className={cn(
                "p-2 rounded-lg transition-colors",
                mode === 'chat' ? "bg-indigo-50 text-indigo-600" : "hover:bg-slate-50 text-slate-400"
            )}
            title="Switch to Copilot Chat"
           >
             <Sparkles className="w-4 h-4" />
           </button>
        </div>
      </div>
    </div>
  );
}
