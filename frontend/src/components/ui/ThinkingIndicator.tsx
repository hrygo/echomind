import React from 'react';
import { cn } from '@/lib/utils';

interface ThinkingIndicatorProps {
  className?: string;
  text?: string;
}

export function ThinkingIndicator({ className, text = "Thinking" }: ThinkingIndicatorProps) {
  return (
    <div className={cn("flex items-center gap-2 px-1", className)}>
      {/* Gradient Sparkle Icon */}
      <div className="relative flex items-center justify-center">
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          className="w-4 h-4 animate-[spin_3s_linear_infinite]"
        >
          <defs>
            <linearGradient id="sparkle-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#3b82f6" /> {/* blue-500 */}
              <stop offset="50%" stopColor="#6366f1" /> {/* indigo-500 */}
              <stop offset="100%" stopColor="#a855f7" /> {/* purple-500 */}
            </linearGradient>
          </defs>
          <path
            d="M12 2L14.4 9.6L22 12L14.4 14.4L12 22L9.6 14.4L2 12L9.6 9.6L12 2Z"
            fill="url(#sparkle-gradient)"
          />
        </svg>
        {/* Glow effect */}
        <div className="absolute inset-0 bg-indigo-500/20 blur-md rounded-full animate-pulse" />
      </div>

      <span className="text-sm font-medium bg-gradient-to-r from-slate-700 to-slate-500 bg-clip-text text-transparent">
        {text}
      </span>

      {/* Pulsing Dots */}
      <div className="flex gap-1 items-center mt-1">
        <div className="w-1 h-1 bg-slate-400 rounded-full animate-[bounce_1.4s_infinite_-0.3s]" />
        <div className="w-1 h-1 bg-slate-400 rounded-full animate-[bounce_1.4s_infinite_-0.15s]" />
        <div className="w-1 h-1 bg-slate-400 rounded-full animate-[bounce_1.4s_infinite]" />
      </div>
    </div>
  );
}
