"use client";

import React from 'react';
import { useLanguage } from "@/lib/i18n/LanguageContext";

export default function InsightsHeader() {
  const { t } = useLanguage();
  return (
    <>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-800">{t('insights.title')}</h1>
      </div>
      <div className="p-4 border-b border-slate-100 bg-slate-50/30">
          <h2 className="text-sm font-semibold text-slate-600 uppercase tracking-wider">{t('insights.contactNetworkGraph')}</h2>
      </div>
    </>
  );
}
