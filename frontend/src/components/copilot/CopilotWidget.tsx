'use client';

import React, { useEffect, useRef } from 'react';
import { useCopilotStore } from '@/store/useCopilotStore';
import { CopilotInput } from './CopilotInput';
import { CopilotResults } from './CopilotResults';
import { CopilotChat } from './CopilotChat';
import { cn } from '@/lib/utils';

export function CopilotWidget() {
  const { mode, isOpen, setIsOpen, reset } = useCopilotStore();
  const containerRef = useRef<HTMLDivElement>(null);

  // Close when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        // Only close if we are not in a deep interaction state? 
        // For now, close and reset mode to idle, but maybe keep query?
        // Let's just collapse result view.
        if (mode !== 'idle') {
            reset();
        }
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [mode, reset]);

  return (
    <div 
        ref={containerRef}
        className="relative z-50 w-full max-w-2xl"
    >
      {/* The Input Bar (Always Visible) */}
      <CopilotInput />

      {/* The Dropdown Area (Absolute Positioned) */}
      <div className={cn(
          "absolute top-full left-0 right-0 mt-1 transition-all duration-200 origin-top",
          mode === 'idle' ? "opacity-0 scale-95 pointer-events-none" : "opacity-100 scale-100 pointer-events-auto"
      )}>
         {mode === 'search' && <CopilotResults />}
         {mode === 'chat' && <CopilotChat />}
      </div>
      
      {/* Backdrop (Optional, adds focus) */}
      {mode !== 'idle' && (
        <div className="fixed inset-0 bg-black/5 -z-10 pointer-events-none" />
      )}
    </div>
  );
}
