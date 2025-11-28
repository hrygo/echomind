'use client';

import React from 'react';
import { Settings, Sparkles, GitBranch } from 'lucide-react';
import { useCopilotStore } from '@/store/useCopilotStore';
import { cn } from '@/lib/utils';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface SearchEnhancementSettingsProps {
  className?: string;
}

export function SearchEnhancementSettings({ className }: SearchEnhancementSettingsProps) {
  const { t } = useLanguage();
  const {
    enableClustering,
    enableSummary,
    clusterType,
    setEnableClustering,
    setEnableSummary,
    setClusterType,
  } = useCopilotStore();

  return (
    <div className={cn("bg-white border border-slate-200 rounded-lg p-4 shadow-sm", className)}>
      <div className="flex items-center gap-2 mb-3">
        <Settings className="w-4 h-4 text-slate-600" />
        <h3 className="text-sm font-semibold text-slate-700">{t('copilot.searchEnhancement.title')}</h3>
      </div>

      <div className="space-y-3">
        {/* Enable AI Summary */}
        <label className="flex items-center justify-between cursor-pointer group">
          <div className="flex items-center gap-2">
            <Sparkles className="w-4 h-4 text-indigo-500" />
            <span className="text-sm text-slate-700 group-hover:text-slate-900">{t('copilot.searchEnhancement.aiSummary')}</span>
          </div>
          <button
            onClick={() => setEnableSummary(!enableSummary)}
            className={cn(
              "relative inline-flex h-5 w-9 items-center rounded-full transition-colors",
              enableSummary ? "bg-indigo-600" : "bg-slate-300"
            )}
          >
            <span
              className={cn(
                "inline-block h-4 w-4 transform rounded-full bg-white transition-transform",
                enableSummary ? "translate-x-5" : "translate-x-0.5"
              )}
            />
          </button>
        </label>

        {/* Enable Clustering */}
        <label className="flex items-center justify-between cursor-pointer group">
          <div className="flex items-center gap-2">
            <GitBranch className="w-4 h-4 text-purple-500" />
            <span className="text-sm text-slate-700 group-hover:text-slate-900">{t('copilot.searchEnhancement.clustering')}</span>
          </div>
          <button
            onClick={() => setEnableClustering(!enableClustering)}
            className={cn(
              "relative inline-flex h-5 w-9 items-center rounded-full transition-colors",
              enableClustering ? "bg-purple-600" : "bg-slate-300"
            )}
          >
            <span
              className={cn(
                "inline-block h-4 w-4 transform rounded-full bg-white transition-transform",
                enableClustering ? "translate-x-5" : "translate-x-0.5"
              )}
            />
          </button>
        </label>

        {/* Cluster Type Selection */}
        {enableClustering && (
          <div className="pl-6 pt-2 border-l-2 border-purple-200">
            <p className="text-xs text-slate-500 mb-2">{t('copilot.searchEnhancement.clusterType')}:</p>
            <div className="grid grid-cols-3 gap-2">
              <button
                onClick={() => setClusterType('sender')}
                className={cn(
                  "px-2 py-1.5 text-xs rounded-lg border transition-all",
                  clusterType === 'sender'
                    ? "bg-blue-50 border-blue-300 text-blue-700 font-medium"
                    : "bg-white border-slate-200 text-slate-600 hover:border-slate-300"
                )}
              >
                {t('copilot.searchEnhancement.clusterBySender')}
              </button>
              <button
                onClick={() => setClusterType('time')}
                className={cn(
                  "px-2 py-1.5 text-xs rounded-lg border transition-all",
                  clusterType === 'time'
                    ? "bg-green-50 border-green-300 text-green-700 font-medium"
                    : "bg-white border-slate-200 text-slate-600 hover:border-slate-300"
                )}
              >
                {t('copilot.searchEnhancement.clusterByTime')}
              </button>
              <button
                onClick={() => setClusterType('topic')}
                className={cn(
                  "px-2 py-1.5 text-xs rounded-lg border transition-all",
                  clusterType === 'topic'
                    ? "bg-purple-50 border-purple-300 text-purple-700 font-medium"
                    : "bg-white border-slate-200 text-slate-600 hover:border-slate-300"
                )}
              >
                {t('copilot.searchEnhancement.clusterByTopic')}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
