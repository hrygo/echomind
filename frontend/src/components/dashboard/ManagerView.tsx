import React from 'react';
import { CheckCircle2, Clock, ArrowRight } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";

export function ManagerView() {
  const { t } = useLanguage();

  return (
    <div className="space-y-6">
      {/* Action Items Section */}
      <section>
        <h2 className="text-xl font-bold text-slate-800 mb-4 flex items-center gap-2">
          <CheckCircle2 className="w-5 h-5 text-green-600" />
          {t('dashboard.actionItems')}
        </h2>
        <div className="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
          <div className="divide-y divide-slate-100">
            {[
              { title: "提交月度运营报告", due: "Today", priority: "High" },
              { title: "回复 Client X 的询价邮件", due: "Tomorrow", priority: "Medium" },
              { title: "确认下周团队建设方案", due: "Fri", priority: "Low" },
            ].map((item, i) => (
              <div key={i} className="p-4 flex items-center gap-4 hover:bg-slate-50 transition-colors group cursor-pointer">
                <div className="w-5 h-5 rounded-full border-2 border-slate-300 group-hover:border-blue-500 transition-colors"></div>
                <div className="flex-1">
                  <h3 className="font-medium text-slate-800 group-hover:text-blue-600 transition-colors">{item.title}</h3>
                  <p className="text-xs text-slate-400 mt-0.5">From: email@example.com</p>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`text-xs font-bold px-2 py-1 rounded-full 
                    ${item.priority === 'High' ? 'bg-red-50 text-red-600' :
                      item.priority === 'Medium' ? 'bg-orange-50 text-orange-600' :
                        'bg-slate-100 text-slate-500'}`}>
                    {item.priority}
                  </span>
                  <span className="text-xs font-medium text-slate-500 bg-slate-50 px-2 py-1 rounded border border-slate-100">
                    Due: {item.due}
                  </span>
                </div>
              </div>
            ))}
          </div>
          <div className="p-3 bg-slate-50 text-center border-t border-slate-100">
            <button className="text-sm font-medium text-blue-600 hover:text-blue-700 flex items-center justify-center gap-1">
              {t('dashboard.viewAll')} <ArrowRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      </section>

      {/* Smart Follow-up Section */}
      <section>
        <h2 className="text-xl font-bold text-slate-800 mb-4 flex items-center gap-2">
          <Clock className="w-5 h-5 text-orange-500" />
          {t('dashboard.smartFollowUp')}
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {[
            { name: "Alice Smith", subject: "Re: Project Proposal", time: "2 days ago" },
            { name: "Bob Jones", subject: "Contract Draft Review", time: "3 days ago" },
          ].map((item, i) => (
            <div key={i} className="bg-white p-4 rounded-xl border border-slate-100 shadow-sm flex items-center gap-4">
              <div className="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center text-orange-600 font-bold">
                {item.name[0]}
              </div>
              <div className="flex-1">
                <h4 className="font-semibold text-slate-800">{item.name}</h4>
                <p className="text-sm text-slate-500 truncate">{item.subject}</p>
              </div>
              <div className="text-right">
                <p className="text-xs font-bold text-orange-500">{t('dashboard.waiting')}</p>
                <p className="text-xs text-slate-400">{item.time}</p>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
