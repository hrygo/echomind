'use client';

import React from 'react';
import { Mail, Calendar, User, List, Grid } from 'lucide-react';
import { useCopilotStore, SearchResult } from '@/store/useCopilotStore';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { SearchSummaryCard } from './SearchSummaryCard';
import { SearchClusterView } from './SearchClusterView';
import { cn } from '@/lib/utils';

import { ThinkingIndicator } from '@/components/ui/ThinkingIndicator';

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
  const {
    searchResults,
    isSearching,
    mode,
    // Search enhancement
    summary,
    clusters,
    clusterType,
    searchViewMode,
    setSearchViewMode,
  } = useCopilotStore();

  if (mode !== 'search') return null;

  const hasEnhancedData = summary || (clusters && clusters.length > 0);

  return (
    <div className="w-full max-w-2xl mx-auto bg-white border border-t-0 rounded-b-xl shadow-lg max-h-[70vh] overflow-y-auto">
      {isSearching ? (
        <div className="p-8 flex justify-center">
          <ThinkingIndicator text={t('copilot.searching')} />
        </div>
      ) : searchResults.length > 0 ? (
        <div className="p-2">
          {/* AI Summary Card */}
          {summary && (
            <div className="px-2 py-2">
              <SearchSummaryCard
                summary={summary}
                resultCount={searchResults.length}
              />
            </div>
          )}

          {/* View Mode Toggle (only show if we have clusters) */}
          {clusters && clusters.length > 0 && (
            <div className="flex items-center justify-between px-3 py-2 border-b border-slate-100">
              <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                {searchViewMode === 'all' ? t('copilot.topResults') : t('copilot.searchEnhancement.clusteredResults')}
              </div>
              <div className="flex gap-1 bg-slate-100 rounded-lg p-1">
                <button
                  onClick={() => setSearchViewMode('all')}
                  className={cn(
                    "px-3 py-1 text-xs rounded transition-colors flex items-center gap-1.5",
                    searchViewMode === 'all'
                      ? "bg-white text-slate-700 shadow-sm"
                      : "text-slate-500 hover:text-slate-700"
                  )}
                >
                  <List className="w-3.5 h-3.5" />
                  {t('copilot.searchEnhancement.allResults')}
                </button>
                <button
                  onClick={() => setSearchViewMode('clustered')}
                  className={cn(
                    "px-3 py-1 text-xs rounded transition-colors flex items-center gap-1.5",
                    searchViewMode === 'clustered'
                      ? "bg-white text-slate-700 shadow-sm"
                      : "text-slate-500 hover:text-slate-700"
                  )}
                >
                  <Grid className="w-3.5 h-3.5" />
                  {t('copilot.searchEnhancement.clustered')}
                </button>
              </div>
            </div>
          )}

          {/* Results Display */}
          {searchViewMode === 'clustered' && clusters && clusters.length > 0 ? (
            <div className="px-2 py-2">
              <SearchClusterView clusters={clusters} clusterType={clusterType} />
            </div>
          ) : (
            <div className="space-y-1 mt-2">
              {!hasEnhancedData && (
                <div className="px-3 py-2 text-xs font-semibold text-slate-400 uppercase tracking-wider">
                  {t('copilot.topResults')}
                </div>
              )}
              {searchResults.map((result) => (
                <ResultItem key={result.email_id} result={result} />
              ))}
            </div>
          )}
        </div>
      ) : (
        <div className="p-8 text-center text-slate-400">
          {t('copilot.noResults')}
        </div>
      )}
    </div>
  );
}
