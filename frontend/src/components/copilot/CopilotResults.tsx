'use client';

import React from 'react';
import { Mail, Calendar, User } from 'lucide-react';
import { useCopilotStore, SearchResult } from '@/store/useCopilotStore';
import { useLanguage } from '@/lib/i18n/LanguageContext';

function ResultItem({ result }: { result: SearchResult }) {
  const { t } = useLanguage();
  return (
    <div className="group flex items-start p-3 hover:bg-slate-50 rounded-lg cursor-pointer transition-colors border border-transparent hover:border-slate-100">
      <div className="mt-1 mr-3 p-2 bg-blue-50 text-blue-600 rounded-full group-hover:bg-blue-100">
        <Mail className="w-4 h-4" />
      </div>
      <div className="flex-1 min-w-0">
        <h4 className="font-medium text-slate-800 truncate">{result.subject}</h4>
        <div className="flex items-center text-xs text-slate-500 mt-0.5 space-x-2">
            <span className="flex items-center">
                <User className="w-3 h-3 mr-1" />
                {result.sender}
            </span>
            <span className="w-1 h-1 bg-slate-300 rounded-full" />
            <span className="flex items-center">
                 {/* Simple date formatting fallback */}
                <Calendar className="w-3 h-3 mr-1" />
                {new Date(result.date).toLocaleDateString()}
            </span>
        </div>
        <p className="text-sm text-slate-600 mt-1 line-clamp-2">
          {result.snippet}
        </p>
      </div>
      <div className="text-xs text-slate-400 self-center opacity-0 group-hover:opacity-100 transition-opacity">
        {Math.round(result.score * 100)}% {t('copilot.match')}
      </div>
    </div>
  );
}

export function CopilotResults() {
  const { t } = useLanguage();
  const { searchResults, isSearching, mode } = useCopilotStore();

  if (mode !== 'search') return null;

  return (
    <div className="w-full max-w-2xl mx-auto bg-white border border-t-0 rounded-b-xl shadow-lg max-h-[60vh] overflow-y-auto">
      {isSearching ? (
        <div className="p-8 text-center text-slate-500">
          <div className="animate-spin w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full mx-auto mb-2"></div>
          {t('copilot.searching')}
        </div>
      ) : searchResults.length > 0 ? (
        <div className="p-2 space-y-1">
            <div className="px-3 py-2 text-xs font-semibold text-slate-400 uppercase tracking-wider">
                {t('copilot.topResults')}
            </div>
          {searchResults.map((result) => (
            <ResultItem key={result.email_id} result={result} />
          ))}
        </div>
      ) : (
        <div className="p-8 text-center text-slate-400">
            {t('copilot.noResults')}
        </div>
      )}
    </div>
  );
}
