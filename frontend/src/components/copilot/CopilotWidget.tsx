'use client';

import React, { useEffect, useRef, useState } from 'react';
import { useCopilotStore } from '@/store/useCopilotStore';
import { CopilotInput } from './CopilotInput';
import { CopilotResults } from './CopilotResults';
import { CopilotChat } from './CopilotChat';
import { SearchEnhancementSettings } from './SearchEnhancementSettings';
import { cn } from '@/lib/utils';

export function CopilotWidget() {
  const { mode, reset } = useCopilotStore();
  const containerRef = useRef<HTMLDivElement>(null);
  const [showSettings, setShowSettings] = useState(false);

  // Close when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        // Close settings panel
        if (showSettings) {
          setShowSettings(false);
        }
        // Close result view if not in a deep interaction state
        if (mode !== 'idle') {
            reset();
        }
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [mode, showSettings, reset]);

  // Auto-close settings when mode changes to search or chat
  useEffect(() => {
    if (mode !== 'idle' && showSettings) {
      setShowSettings(false);
    }
  }, [mode, showSettings]);

  return (
    <div className="relative z-50 w-full max-w-2xl">
      <div ref={containerRef} className="relative">
        {/* The Input Bar with integrated Settings Button */}
        <CopilotInput 
          showSettings={showSettings}
          onToggleSettings={() => setShowSettings(!showSettings)}
          onCloseSettings={() => setShowSettings(false)}
        />

        {/* Settings Panel - positioned below input */}
        {showSettings && (
          <div className="absolute top-full left-0 right-0 mt-2 z-50 animate-in fade-in slide-in-from-top-2 duration-200">
            <SearchEnhancementSettings />
          </div>
        )}

        {/* The Dropdown Area (Search Results or Chat) - Only show when settings is closed */}
        {!showSettings && (
          <div 
            className={cn(
              "absolute top-full left-0 right-0 transition-all duration-200 origin-top mt-1",
              mode === 'idle' ? "opacity-0 scale-95 pointer-events-none" : "opacity-100 scale-100 pointer-events-auto"
            )}
          >
             {mode === 'search' && <CopilotResults />}
             {mode === 'chat' && <CopilotChat />}
          </div>
        )}
      </div>
      
      {/* Backdrop */}
      {(mode !== 'idle' || showSettings) && (
        <div className="fixed inset-0 bg-black/5 -z-10 pointer-events-none" />
      )}
    </div>
  );
}
