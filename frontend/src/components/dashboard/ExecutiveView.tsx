import React from 'react';
import { AlertTriangle, TrendingUp, Calendar } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";

export function ExecutiveView() {
    const { t } = useLanguage();

    return (
        <div className="space-y-6">
            {/* Daily Digest Section */}
            <section>
                <h2 className="text-xl font-bold text-slate-800 mb-4 flex items-center gap-2">
                    <NewspaperIcon className="w-5 h-5 text-blue-600" />
                    {t('dashboard.dailyDigest')}
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Risk Card */}
                    <div className="bg-white p-5 rounded-2xl border border-red-100 shadow-sm hover:shadow-md transition-shadow">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2 bg-red-50 rounded-lg">
                                <AlertTriangle className="w-5 h-5 text-red-600" />
                            </div>
                            <span className="text-xs font-bold text-red-600 bg-red-50 px-2 py-1 rounded-full">High Risk</span>
                        </div>
                        <h3 className="text-slate-600 text-sm font-medium mb-1">{t('dashboard.riskWarning')}</h3>
                        <p className="text-2xl font-bold text-slate-800">3 封</p>
                        <p className="text-xs text-slate-400 mt-2">来自关键客户的负面反馈</p>
                    </div>

                    {/* Trend Card */}
                    <div className="bg-white p-5 rounded-2xl border border-blue-100 shadow-sm hover:shadow-md transition-shadow">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2 bg-blue-50 rounded-lg">
                                <TrendingUp className="w-5 h-5 text-blue-600" />
                            </div>
                            <span className="text-xs font-bold text-blue-600 bg-blue-50 px-2 py-1 rounded-full">+12%</span>
                        </div>
                        <h3 className="text-slate-600 text-sm font-medium mb-1">{t('dashboard.weeklyInteraction')}</h3>
                        <p className="text-2xl font-bold text-slate-800">128 次</p>
                        <p className="text-xs text-slate-400 mt-2">与核心团队及客户</p>
                    </div>

                    {/* Schedule Card */}
                    <div className="bg-white p-5 rounded-2xl border border-purple-100 shadow-sm hover:shadow-md transition-shadow">
                        <div className="flex items-start justify-between mb-3">
                            <div className="p-2 bg-purple-50 rounded-lg">
                                <Calendar className="w-5 h-5 text-purple-600" />
                            </div>
                            <span className="text-xs font-bold text-purple-600 bg-purple-50 px-2 py-1 rounded-full">Today</span>
                        </div>
                        <h3 className="text-slate-600 text-sm font-medium mb-1">{t('dashboard.pendingDecisions')}</h3>
                        <p className="text-2xl font-bold text-slate-800">5 项</p>
                        <p className="text-xs text-slate-400 mt-2">合同审批 & 预算确认</p>
                    </div>
                </div>
            </section>

            {/* Briefing List */}
            <section className="bg-white rounded-2xl border border-slate-100 shadow-sm p-6">
                <h3 className="text-lg font-bold text-slate-800 mb-4">{t('dashboard.briefing')}</h3>
                <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="flex gap-4 p-4 rounded-xl hover:bg-slate-50 transition-colors border border-transparent hover:border-slate-100">
                            <div className="w-1 h-12 bg-blue-500 rounded-full flex-shrink-0"></div>
                            <div>
                                <h4 className="font-semibold text-slate-800">Q4 季度财务预算审批申请</h4>
                                <p className="text-sm text-slate-500 mt-1 line-clamp-2">
                                    财务部提交了 Q4 预算方案，总额比 Q3 增长 15%，主要用于市场推广。需要您在周五前确认。
                                </p>
                                <div className="flex gap-2 mt-2">
                                    <span className="text-xs font-medium text-slate-400 bg-slate-100 px-2 py-0.5 rounded">财务</span>
                                    <span className="text-xs font-medium text-slate-400 bg-slate-100 px-2 py-0.5 rounded">审批</span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </section>
        </div>
    );
}

function NewspaperIcon(props: any) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a2 2 0 0 1-2 2Zm0 0a2 2 0 0 1-2-2v-9c0-1.1.9-2 2-2h2" />
            <path d="M18 14h-8" />
            <path d="M15 18h-5" />
            <path d="M10 6h8v4h-8V6Z" />
        </svg>
    );
}
