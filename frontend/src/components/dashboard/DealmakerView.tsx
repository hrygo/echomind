import React from 'react';
import { Radar, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, ResponsiveContainer } from 'recharts';
import { ArrowRight, DollarSign, Users, Briefcase, Zap } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";

const radarData = [
  { subject: 'Buying Intent', A: 120, fullMark: 150 },
  { subject: 'Partnership', A: 98, fullMark: 150 },
  { subject: 'Hiring', A: 86, fullMark: 150 },
  { subject: 'Meeting Request', A: 99, fullMark: 150 },
  { subject: 'Referral', A: 85, fullMark: 150 },
  { subject: 'Follow-up', A: 65, fullMark: 150 },
];

const opportunities = [
    { id: 1, title: "Enterprise License Inquiry", company: "TechCorp", value: "$50k", confidence: 92, type: "buying" },
    { id: 2, title: "Strategic Partnership Q1", company: "Innovate Inc", value: "Strategic", confidence: 85, type: "partnership" },
    { id: 3, title: "Team Expansion Discussion", company: "GlobalSoft", value: "10 Seats", confidence: 78, type: "buying" },
];

export function DealmakerView() {
    const { t } = useLanguage();
    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Radar Chart Section */}
            <div className="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm lg:col-span-1 flex flex-col">
                <h3 className="text-lg font-bold text-slate-800 mb-2 flex items-center gap-2">
                    <Zap className="w-5 h-5 text-amber-500" />
                    {t('dashboard.intentRadar')}
                </h3>
                <p className="text-sm text-slate-500 mb-4">{t('dashboard.radarDescription')}</p>
                
                <div className="flex-1 min-h-[300px] -ml-6">
                    <ResponsiveContainer width="100%" height="100%">
                        <RadarChart cx="50%" cy="50%" outerRadius="70%" data={radarData}>
                            <PolarGrid stroke="#e2e8f0" />
                            <PolarAngleAxis dataKey="subject" tick={{ fill: '#64748b', fontSize: 10 }} />
                            <PolarRadiusAxis angle={30} domain={[0, 150]} tick={false} axisLine={false} />
                            <Radar
                                name="Intent"
                                dataKey="A"
                                stroke="#3b82f6"
                                strokeWidth={2}
                                fill="#3b82f6"
                                fillOpacity={0.2}
                            />
                        </RadarChart>
                    </ResponsiveContainer>
                </div>
                <div className="mt-4 text-center">
                    <span className="inline-flex items-center px-3 py-1 rounded-full bg-blue-50 text-blue-700 text-xs font-medium">
                        {t('dashboard.topSignal')}
                    </span>
                </div>
            </div>

            {/* Opportunity List Section */}
            <div className="lg:col-span-2 space-y-6">
                <div className="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
                    <div className="p-5 border-b border-slate-50 flex justify-between items-center">
                        <h3 className="text-lg font-bold text-slate-800">{t('dashboard.detectedOpportunities')}</h3>
                        <button className="text-sm text-blue-600 hover:text-blue-700 font-medium">{t('dashboard.viewAll')}</button>
                    </div>
                    <div className="divide-y divide-slate-50">
                        {opportunities.map((opp) => (
                            <div key={opp.id} className="p-5 hover:bg-slate-50 transition-colors cursor-pointer group flex items-center gap-4">
                                <div className={`w-10 h-10 rounded-xl flex items-center justify-center 
                                    ${opp.type === 'buying' ? 'bg-green-100 text-green-600' : 'bg-purple-100 text-purple-600'}
                                `}>
                                    {opp.type === 'buying' ? <DollarSign className="w-5 h-5" /> : <Briefcase className="w-5 h-5" />}
                                </div>
                                <div className="flex-1">
                                    <h4 className="font-semibold text-slate-800 group-hover:text-blue-600 transition-colors">{opp.title}</h4>
                                    <p className="text-xs text-slate-500 mt-0.5">{opp.company}</p>
                                </div>
                                <div className="text-right">
                                    <div className="text-sm font-bold text-slate-700">{opp.value}</div>
                                    <div className="text-xs font-medium text-green-600 mt-0.5">{opp.confidence}% {t('dashboard.confidence')}</div>
                                </div>
                                <ArrowRight className="w-4 h-4 text-slate-300 group-hover:text-slate-500" />
                            </div>
                        ))}
                    </div>
                </div>

                {/* Recent Connections */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-2xl p-5 text-white shadow-lg">
                        <div className="flex items-center gap-3 mb-4">
                            <div className="p-2 bg-white/10 rounded-lg">
                                <Users className="w-5 h-5 text-white" />
                            </div>
                            <div>
                                <h4 className="font-bold">{t('dashboard.newConnections')}</h4>
                                <p className="text-xs text-slate-400">{t('dashboard.thisWeek')}</p>
                            </div>
                        </div>
                        <div className="text-3xl font-bold mb-2">12</div>
                        <div className="flex -space-x-2 overflow-hidden">
                            {[1,2,3,4].map(i => (
                                <div key={i} className="inline-block h-8 w-8 rounded-full ring-2 ring-slate-900 bg-slate-700 flex items-center justify-center text-xs">
                                    {String.fromCharCode(64 + i)}
                                </div>
                            ))}
                            <div className="inline-block h-8 w-8 rounded-full ring-2 ring-slate-900 bg-slate-800 flex items-center justify-center text-xs text-slate-400">+8</div>
                        </div>
                    </div>
                    
                    <div className="bg-white rounded-2xl p-5 border border-slate-100 shadow-sm flex flex-col justify-center items-center text-center hover:border-blue-200 transition-colors cursor-pointer border-dashed">
                        <div className="w-12 h-12 bg-blue-50 rounded-full flex items-center justify-center mb-3 text-blue-500">
                            <Users className="w-6 h-6" />
                        </div>
                        <h4 className="font-semibold text-slate-800">{t('dashboard.exploreNetwork')}</h4>
                        <p className="text-xs text-slate-500 mt-1">{t('dashboard.viewFullGraph')}</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
