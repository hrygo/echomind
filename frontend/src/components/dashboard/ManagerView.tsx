import React, { useState } from 'react';
import { CheckCircle2, Clock, ArrowRight, Calendar, AlertCircle } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { TaskWidget } from "./TaskWidget"; // Import TaskWidget
import Link from 'next/link';

export function ManagerView() {
  const { t } = useLanguage();
  const [filter, setFilter] = useState<'all' | 'high'>('all');

  // TODO: Replace with actual counts from useTaskStore or an API call for dashboard stats
  const activeTasksCount = 5; // Placeholder
  const overdueTasksCount = 1; // Placeholder

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Main Column: Action Items */}
      <div className="lg:col-span-2 space-y-6">
        <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold text-slate-800 flex items-center gap-2">
            <CheckCircle2 className="w-5 h-5 text-blue-600" />
            {t('dashboard.actionItems')}
            <span className="text-xs font-normal text-slate-400 bg-slate-100 px-2 py-0.5 rounded-full ml-2">
                {activeTasksCount} {t('dashboard.pending')}
            </span>
            </h2>
            
            <div className="flex gap-2">
                <button 
                    onClick={() => setFilter('all')}
                    className={`text-xs font-medium px-3 py-1.5 rounded-lg transition-colors ${filter === 'all' ? 'bg-slate-800 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}`}
                >
                    {t('dashboard.all')}
                </button>
                <button 
                    onClick={() => setFilter('high')}
                    className={`text-xs font-medium px-3 py-1.5 rounded-lg transition-colors ${filter === 'high' ? 'bg-red-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}`}
                >
                    {t('dashboard.highPriority')}
                </button>
            </div>
        </div>

        <div className="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
          {/* Integrate TaskWidget here */}
          <TaskWidget initialPriority={filter === 'high' ? 'high' : undefined} />
          
          <div className="p-3 bg-slate-50/50 text-center border-t border-slate-100">
            <Link href="/dashboard/tasks" className="text-sm font-medium text-blue-600 hover:text-blue-700 flex items-center justify-center gap-1 transition-colors">
              {t('dashboard.viewAll')} <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </div>
      </div>

      {/* Right Column: Follow-ups & Stats */}
      <div className="space-y-6">
        {/* Stats Widget */}
        <div className="grid grid-cols-2 gap-3">
             <div className="bg-blue-50/50 p-4 rounded-xl border border-blue-100">
                 <p className="text-xs font-medium text-blue-600 uppercase tracking-wider mb-1">{t('dashboard.completed')}</p>
                 <p className="text-2xl font-bold text-slate-800">14</p>
                 <p className="text-[10px] text-slate-400">{t('dashboard.thisWeek')}</p>
             </div>
             <div className="bg-orange-50/50 p-4 rounded-xl border border-orange-100">
                 <p className="text-xs font-medium text-orange-600 uppercase tracking-wider mb-1">{t('dashboard.overdue')}</p>
                 <p className="text-2xl font-bold text-slate-800">{overdueTasksCount}</p>
                 <p className="text-[10px] text-slate-400">{t('dashboard.actionNeeded')}</p>
             </div>
        </div>

        {/* Smart Follow-up */}
        <section>
            <h2 className="text-lg font-bold text-slate-800 mb-3 flex items-center gap-2">
            <Clock className="w-4 h-4 text-orange-500" />
            {t('dashboard.smartFollowUp')}
            </h2>
            <div className="space-y-3">
            {[
                { name: "Alice Smith", subject: "Re: Project Proposal", time: "2d", waiting: true },
                { name: "Bob Jones", subject: "Contract Draft Review", time: "3d", waiting: true },
                { name: "Charlie Day", subject: "Lunch Meeting", time: "5h", waiting: false }, // Not waiting, just recent
            ].filter(i => i.waiting).map((item, i) => (
                <div key={i} className="bg-white p-3 rounded-xl border border-slate-100 shadow-sm flex items-center gap-3 hover:shadow-md transition-shadow cursor-pointer">
                    <div className="w-8 h-8 rounded-full bg-orange-100 flex items-center justify-center text-orange-600 font-bold text-xs">
                        {item.name[0]}
                    </div>
                    <div className="flex-1 min-w-0">
                        <h4 className="font-semibold text-sm text-slate-800 truncate">{item.name}</h4>
                        <p className="text-xs text-slate-500 truncate">{item.subject}</p>
                    </div>
                    <div className="text-right whitespace-nowrap">
                        <span className="text-[10px] font-bold text-orange-500 bg-orange-50 px-1.5 py-0.5 rounded">{t('dashboard.waiting')}</span>
                        <p className="text-[10px] text-slate-400 mt-0.5">{item.time}</p>
                    </div>
                </div>
            ))}
            </div>
        </section>
      </div>
    </div>
  );
}
