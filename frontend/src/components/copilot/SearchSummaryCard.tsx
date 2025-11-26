'use client';

import React from 'react';
import { Sparkles, Tag, User, Info } from 'lucide-react';
import { SearchResultsSummary } from '@/types/search';
import { cn } from '@/lib/utils';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface SearchSummaryCardProps {
  summary: SearchResultsSummary;
  resultCount: number;
  className?: string;
}

export function SearchSummaryCard({ summary, resultCount, className }: SearchSummaryCardProps) {
  const { t } = useLanguage();
  const { natural_summary, key_topics, important_people } = summary;

  return (
    <div className={cn(
      "bg-gradient-to-br from-indigo-50/80 to-purple-50/80 rounded-xl p-5 border border-indigo-100/50 shadow-sm",
      className
    )}>
      {/* Header */}
      <div className="flex items-center gap-2 mb-4">
        <div className="p-2 bg-white rounded-lg shadow-sm">
          <Sparkles className="w-4 h-4 text-indigo-600" />
        </div>
        <div>
          <h3 className="text-sm font-semibold text-slate-800">{t('copilot.searchEnhancement.summaryTitle')}</h3>
          <p className="text-xs text-slate-500">{t('copilot.searchEnhancement.emailCount', { count: resultCount })}</p>
        </div>
      </div>

      {/* Natural Summary */}
      {natural_summary && (
        <div className="mb-4 p-4 bg-white/60 rounded-lg">
          <div className="flex items-start gap-2">
            <Info className="w-4 h-4 text-indigo-500 mt-0.5 flex-shrink-0" />
            <p className="text-sm text-slate-700 leading-relaxed">
              {natural_summary}
            </p>
          </div>
        </div>
      )}

      {/* Key Topics and Important People */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Key Topics */}
        {key_topics && key_topics.length > 0 && (
          <div>
            <div className="flex items-center gap-1.5 mb-2">
              <Tag className="w-3.5 h-3.5 text-purple-600" />
              <h4 className="text-xs font-medium text-slate-700">{t('copilot.searchEnhancement.keyTopics')}</h4>
            </div>
            <div className="flex flex-wrap gap-1.5">
              {key_topics.map((topic, index) => (
                <span
                  key={index}
                  className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-700 hover:bg-purple-200 transition-colors cursor-default"
                >
                  {topic}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Important People */}
        {important_people && important_people.length > 0 && (
          <div>
            <div className="flex items-center gap-1.5 mb-2">
              <User className="w-3.5 h-3.5 text-indigo-600" />
              <h4 className="text-xs font-medium text-slate-700">{t('copilot.searchEnhancement.importantPeople')}</h4>
            </div>
            <div className="flex flex-wrap gap-1.5">
              {important_people.map((person, index) => (
                <span
                  key={index}
                  className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-700 hover:bg-indigo-200 transition-colors cursor-default"
                >
                  {person}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Empty State */}
      {(!key_topics || key_topics.length === 0) && (!important_people || important_people.length === 0) && !natural_summary && (
        <div className="text-center py-4">
          <p className="text-sm text-slate-400">{t('copilot.searchEnhancement.noSummary')}</p>
        </div>
      )}
    </div>
  );
}
