"use client";

import React from 'react';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { PageHeader } from "@/components/ui/page-header";

export default function InsightsHeader() {
  const { t } = useLanguage();
  return (
    <PageHeader title={t('insights.title')} />
  );
}
