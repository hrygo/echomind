'use client';

import React from 'react';
import { Sparkles } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface AuthLayoutProps {
  children: React.ReactNode;
}

export function AuthLayout({ children }: AuthLayoutProps) {
  const { t } = useLanguage();

  return (
    <div className="flex min-h-screen bg-slate-50">
      {/* Left Side: Visuals (Hidden on mobile) */}
      <div className={cn(
        "hidden lg:flex flex-1 flex-col justify-center items-center p-8 text-white relative overflow-hidden",
        "bg-gradient-to-br from-indigo-800 to-blue-900"
      )}>
        {/* Abstract neural network animation - Placeholder */}
        <div className="absolute inset-0 z-0 opacity-20">
          {/* Replace with actual SVG/particle animation later */}
          <svg className="w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
            <defs>
              <radialGradient id="gradient1" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stopColor="#6366F1" stopOpacity="0.5" />
                <stop offset="100%" stopColor="#1E3A8A" stopOpacity="0" />
              </radialGradient>
            </defs>
            <circle cx="20" cy="30" r="15" fill="url(#gradient1)" opacity="0.7" />
            <circle cx="80" cy="70" r="10" fill="url(#gradient1)" opacity="0.7" />
            <line x1="20" y1="30" x2="80" y2="70" stroke="#818CF8" strokeWidth="0.5" opacity="0.5" />
          </svg>
        </div>

        <div className="relative z-10 text-center">
          <Sparkles className="w-16 h-16 mx-auto text-blue-300 mb-4" />
          <h1 className="text-5xl font-extrabold tracking-tight mb-4 leading-tight">
            {t('common.appName')}
          </h1>
          <p className="text-xl font-light opacity-80 mb-8">
            {t('common.appSlogan')}
          </p>
          {/* Product feature highlights - Placeholder */}
          <ul className="text-lg space-y-3 opacity-90 text-left mx-auto max-w-sm">
            <li className="flex items-center gap-3">
              <Sparkles className="w-5 h-5 flex-shrink-0 text-blue-300" />
              <span>{t('auth.feature1')}</span>
            </li>
            <li className="flex items-center gap-3">
              <Sparkles className="w-5 h-5 flex-shrink-0 text-blue-300" />
              <span>{t('auth.feature2')}</span>
            </li>
            <li className="flex items-center gap-3">
              <Sparkles className="w-5 h-5 flex-shrink-0 text-blue-300" />
              <span>{t('auth.feature3')}</span>
            </li>
          </ul>
        </div>
      </div>

      {/* Right Side: Auth Form */}
      <div className="flex-1 flex items-center justify-center p-4 lg:p-8 overflow-y-auto">
        <div className="w-full max-w-md bg-white rounded-2xl shadow-xl border border-slate-100 p-8 lg:p-10">
          {children}
        </div>
      </div>
    </div>
  );
}
