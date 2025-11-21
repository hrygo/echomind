import React from 'react';
import { AlertTriangle, TrendingUp, Calendar } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { SmartFeed } from './SmartFeed';

export function ExecutiveView() {
    const { t } = useLanguage();

    return (
        <div className="space-y-6">
            {/* Daily Digest Section */}
            <section>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Risk Card */}
                    <div className="bg-white p-5 rounded-2xl border border-red-100 shadow-sm hover:shadow-md transition-all duration-200 group">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2.5 bg-red-50 rounded-xl group-hover:bg-red-100 transition-colors">
                                <AlertTriangle className="w-5 h-5 text-red-600" />
                            </div>
                            <span className="text-[11px] font-bold text-red-600 bg-red-50 px-2 py-1 rounded-full uppercase tracking-wide">{t('dashboard.highRisk')}</span>
                        </div>
                        <h3 className="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">{t('dashboard.riskWarning')}</h3>
                        <p className="text-3xl font-bold text-slate-800 tracking-tight">3</p>
                        <p className="text-xs text-slate-400 mt-2 font-medium">{t('dashboard.criticalFeedback')}</p>
                    </div>

                    {/* Trend Card */}
                    <div className="bg-white p-5 rounded-2xl border border-blue-100 shadow-sm hover:shadow-md transition-all duration-200 group">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2.5 bg-blue-50 rounded-xl group-hover:bg-blue-100 transition-colors">
                                <TrendingUp className="w-5 h-5 text-blue-600" />
                            </div>
                            <span className="text-[11px] font-bold text-blue-600 bg-blue-50 px-2 py-1 rounded-full">+12%</span>
                        </div>
                        <h3 className="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">{t('dashboard.weeklyInteraction')}</h3>
                        <p className="text-3xl font-bold text-slate-800 tracking-tight">128</p>
                        <p className="text-xs text-slate-400 mt-2 font-medium">{t('dashboard.keyStakeholderEngagement')}</p>
                    </div>

                    {/* Schedule Card */}
                    <div className="bg-white p-5 rounded-2xl border border-purple-100 shadow-sm hover:shadow-md transition-all duration-200 group">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2.5 bg-purple-50 rounded-xl group-hover:bg-purple-100 transition-colors">
                                <Calendar className="w-5 h-5 text-purple-600" />
                            </div>
                            <span className="text-[11px] font-bold text-purple-600 bg-purple-50 px-2 py-1 rounded-full">{t('dashboard.today')}</span>
                        </div>
                        <h3 className="text-slate-500 text-xs font-semibold uppercase tracking-wider mb-1">{t('dashboard.pendingDecisions')}</h3>
                        <p className="text-3xl font-bold text-slate-800 tracking-tight">5</p>
                        <p className="text-xs text-slate-400 mt-2 font-medium">{t('dashboard.approvalAndBudgetReview')}</p>
                    </div>
                </div>
            </section>

            {/* Smart Feed Section */}
            <section className="animate-in fade-in slide-in-from-bottom-2 duration-500 delay-100">
                <SmartFeed />
            </section>
        </div>
    );
}
