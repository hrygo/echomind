import React from 'react';
import { AlertTriangle, TrendingUp, Calendar, Loader2 } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { SmartFeed } from './SmartFeed';
import { useExecutiveOverview } from '@/hooks/useExecutiveOverview';
import { useRisks } from '@/hooks/useRisks';
import { useTrends } from '@/hooks/useTrends';

export function ExecutiveView({ contextId }: { contextId?: string | null }) {
    const { t } = useLanguage();

    // Fetch executive data
    const { data: overview, isLoading: overviewLoading } = useExecutiveOverview();
    const { data: risks, isLoading: risksLoading } = useRisks();
    const { data: trends, isLoading: trendsLoading } = useTrends();

    // Calculate display values
    const highRiskCount = contextId ? 1 : (risks?.totalRiskCount || overview?.criticalAlerts || 0);
    const weeklyInteractions = contextId ? 12 : (trends?.weeklyInteraction || 128);
    const interactionChange = trends?.interactionChange || 12;
    const pendingDecisions = contextId ? 0 : (overview?.upcomingDeadlines || 5);

    // Format interaction change
    const formatInteractionChange = (change: number) => {
        const sign = change >= 0 ? '+' : '';
        return `${sign}${change}%`;
    };

    // Get trend icon and color
    const getTrendIcon = (trend: 'upward' | 'downward' | 'stable') => {
        return trend === 'upward' ? '📈' : trend === 'downward' ? '📉' : '➡️';
    };

    return (
        <div className="space-y-6">
            {/* Daily Digest Section */}
            <section>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Risk Card */}
                    <div className="bg-white p-5 rounded-2xl border border-red-100 shadow-sm hover:shadow-md transition-all duration-200 group">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2.5 bg-red-50 rounded-xl group-hover:bg-red-100 transition-colors">
                                {risksLoading ? (
                                    <Loader2 className="w-5 h-5 text-red-600 animate-spin" />
                                ) : (
                                    <AlertTriangle className="w-5 h-5 text-red-600" />
                                )}
                            </div>
                            <span className="text-[11px] font-bold text-red-600 bg-red-50 px-2 py-1 rounded-full uppercase tracking-wide">
                                {t('dashboard.highRisk')}
                            </span>
                        </div>
                        <h3 className="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
                            {t('dashboard.riskWarning')}
                        </h3>
                        <p className="text-3xl font-bold text-slate-800 tracking-tight">
                            {risksLoading ? '--' : highRiskCount}
                        </p>
                        <p className="text-xs text-slate-400 mt-2 font-medium">
                            {risks?.riskTrend === 'increasing' ? '⚠️ ' : ''}{t('dashboard.criticalFeedback')}
                        </p>
                    </div>

                    {/* Trend Card */}
                    <div className="bg-white p-5 rounded-2xl border border-blue-100 shadow-sm hover:shadow-md transition-all duration-200 group">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2.5 bg-blue-50 rounded-xl group-hover:bg-blue-100 transition-colors">
                                {trendsLoading ? (
                                    <Loader2 className="w-5 h-5 text-blue-600 animate-spin" />
                                ) : (
                                    <TrendingUp className="w-5 h-5 text-blue-600" />
                                )}
                            </div>
                            <span className="text-[11px] font-bold text-blue-600 bg-blue-50 px-2 py-1 rounded-full">
                                {formatInteractionChange(interactionChange)}
                            </span>
                        </div>
                        <h3 className="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
                            {t('dashboard.weeklyInteraction')}
                        </h3>
                        <p className="text-3xl font-bold text-slate-800 tracking-tight">
                            {trendsLoading ? '--' : weeklyInteractions.toLocaleString()}
                        </p>
                        <p className="text-xs text-slate-400 mt-2 font-medium">
                            {overview?.productivityTrend && getTrendIcon(overview.productivityTrend)} {t('dashboard.keyStakeholderEngagement')}
                        </p>
                    </div>

                    {/* Schedule Card */}
                    <div className="bg-white p-5 rounded-2xl border border-purple-100 shadow-sm hover:shadow-md transition-all duration-200 group">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2.5 bg-purple-50 rounded-xl group-hover:bg-purple-100 transition-colors">
                                {overviewLoading ? (
                                    <Loader2 className="w-5 h-5 text-purple-600 animate-spin" />
                                ) : (
                                    <Calendar className="w-5 h-5 text-purple-600" />
                                )}
                            </div>
                            <span className="text-[11px] font-bold text-purple-600 bg-purple-50 px-2 py-1 rounded-full">
                                {t('dashboard.today')}
                            </span>
                        </div>
                        <h3 className="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">
                            {t('dashboard.pendingDecisions')}
                        </h3>
                        <p className="text-3xl font-bold text-slate-800 tracking-tight">
                            {overviewLoading ? '--' : pendingDecisions}
                        </p>
                        <p className="text-xs text-slate-400 mt-2 font-medium">
                            {t('dashboard.approvalAndBudgetReview')}
                        </p>
                    </div>
                </div>
            </section>

            {/* Smart Feed Section */}
            <section className="animate-in fade-in slide-in-from-bottom-2 duration-500 delay-100">
                <SmartFeed contextId={contextId} />
            </section>
        </div>
    );
}
