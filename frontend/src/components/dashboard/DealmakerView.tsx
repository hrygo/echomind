import React from 'react';
import { Radar } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";

export function DealmakerView() {
    const { t } = useLanguage();

    return (
        <div className="flex flex-col items-center justify-center h-96 text-slate-400 space-y-4 border-2 border-dashed border-slate-200 rounded-2xl bg-slate-50/50">
            <div className="p-4 bg-white rounded-full shadow-sm">
                <Radar className="w-8 h-8 text-slate-300" />
            </div>
            <div className="text-center">
                <h3 className="text-lg font-semibold text-slate-600">{t('dashboard.dealmakerView')}</h3>
                <p className="text-sm">{t('dashboard.dealmakerDesc')}</p>
            </div>
        </div>
    );
}
