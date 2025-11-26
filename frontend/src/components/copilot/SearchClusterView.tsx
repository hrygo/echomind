'use client';

import React, { useState } from 'react';
import { Users, Clock, Hash, ChevronDown, ChevronRight, Mail } from 'lucide-react';
import { SearchCluster, ClusterType, SearchResult } from '@/types/search';
import { cn } from '@/lib/utils';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface SearchClusterViewProps {
  clusters: SearchCluster[];
  clusterType: ClusterType;
  onResultClick?: (result: SearchResult) => void;
  className?: string;
}

export function SearchClusterView({ 
  clusters, 
  clusterType,
  onResultClick, 
  className 
}: SearchClusterViewProps) {
  const { t } = useLanguage();
  const [expandedClusters, setExpandedClusters] = useState<Set<string>>(
    new Set(clusters.slice(0, 2).map(c => c.id)) // 默认展开前2个
  );

  const toggleCluster = (clusterId: string) => {
    setExpandedClusters(prev => {
      const next = new Set(prev);
      if (next.has(clusterId)) {
        next.delete(clusterId);
      } else {
        next.add(clusterId);
      }
      return next;
    });
  };

  const getClusterIcon = (type: ClusterType) => {
    switch (type) {
      case 'sender':
        return Users;
      case 'time':
        return Clock;
      case 'topic':
        return Hash;
      default:
        return Mail;
    }
  };

  const getClusterColor = (type: ClusterType) => {
    switch (type) {
      case 'sender':
        return 'text-blue-600 bg-blue-50 border-blue-200';
      case 'time':
        return 'text-green-600 bg-green-50 border-green-200';
      case 'topic':
        return 'text-purple-600 bg-purple-50 border-purple-200';
      default:
        return 'text-slate-600 bg-slate-50 border-slate-200';
    }
  };

  if (!clusters || clusters.length === 0) {
    return (
      <div className={cn("text-center py-8", className)}>
        <p className="text-sm text-slate-400">{t('copilot.searchEnhancement.noClusterData')}</p>
      </div>
    );
  }

  const Icon = getClusterIcon(clusterType);

  return (
    <div className={cn("space-y-3", className)}>
      {clusters.map((cluster) => {
        const isExpanded = expandedClusters.has(cluster.id);
        const colorClass = getClusterColor(cluster.type);

        return (
          <div key={cluster.id} className="border border-slate-200 rounded-lg overflow-hidden bg-white shadow-sm">
            {/* Cluster Header */}
            <button
              onClick={() => toggleCluster(cluster.id)}
              className={cn(
                "w-full flex items-center justify-between p-4 hover:bg-slate-50 transition-colors",
                isExpanded && "bg-slate-50/50"
              )}
            >
              <div className="flex items-center gap-3">
                <div className={cn("p-2 rounded-lg border", colorClass)}>
                  <Icon className="w-4 h-4" />
                </div>
                <div className="text-left">
                  <h4 className="text-sm font-medium text-slate-800">
                    {cluster.label}
                  </h4>
                  <p className="text-xs text-slate-500">
                    {t('copilot.searchEnhancement.emailCount', { count: cluster.count })}
                  </p>
                </div>
              </div>
              {isExpanded ? (
                <ChevronDown className="w-5 h-5 text-slate-400" />
              ) : (
                <ChevronRight className="w-5 h-5 text-slate-400" />
              )}
            </button>

            {/* Cluster Results */}
            {isExpanded && (
              <div className="border-t border-slate-100 bg-slate-50/30">
                <div className="divide-y divide-slate-100">
                  {cluster.results.map((result, index) => (
                    <div
                      key={`${result.email_id}-${index}`}
                      onClick={() => onResultClick?.(result)}
                      className="p-4 hover:bg-white transition-colors cursor-pointer"
                    >
                      <div className="flex items-start gap-3">
                        <Mail className="w-4 h-4 text-slate-400 mt-1 flex-shrink-0" />
                        <div className="flex-1 min-w-0">
                          <h5 className="text-sm font-medium text-slate-800 truncate mb-1">
                            {result.subject}
                          </h5>
                          <p className="text-xs text-slate-600 mb-2">
                            {t('copilot.searchEnhancement.from')}: <span className="font-medium">{result.sender}</span>
                          </p>
                          <p className="text-xs text-slate-500 line-clamp-2">
                            {result.snippet}
                          </p>
                          <div className="flex items-center gap-2 mt-2">
                            <span className="text-xs text-slate-400">
                              {new Date(result.date).toLocaleDateString('zh-CN')}
                            </span>
                            {result.score && (
                              <span className="text-xs text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-full">
                                {t('copilot.searchEnhancement.matchScore', { score: (result.score * 100).toFixed(0) })}
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
