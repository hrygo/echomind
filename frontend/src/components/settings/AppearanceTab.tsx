'use client';

import React from 'react';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { ThemeToggle } from '@/components/theme/ThemeToggle';

export function AppearanceTab() {
  const { t } = useLanguage();

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="border-b border-border pb-4">
        <h3 className="text-2xl font-semibold text-foreground">{t('settings.appearance.title')}</h3>
        <p className="text-muted-foreground text-sm mt-2">{t('settings.appearance.description')}</p>
      </div>

      {/* Theme Selection */}
      <div className="space-y-4">
        <div>
          <label className="text-sm font-medium text-foreground mb-3 block">
            {t('settings.appearance.themeSelector')}
          </label>
          <div className="p-4 bg-card rounded-lg border border-border">
            <ThemeToggle />
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            {t('settings.appearance.description')}
          </p>
        </div>
      </div>

      {/* Theme Preview */}
      <div className="mt-8">
        <h4 className="text-sm font-medium text-foreground mb-3">
          {t('settings.appearance.preview')}
        </h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="p-6 bg-card rounded-lg border border-border space-y-3">
            <h5 className="font-medium text-foreground">
              {t('settings.appearance.previewCard')}
            </h5>
            <p className="text-sm text-muted-foreground">
              {t('settings.appearance.previewDesc')}
            </p>
            <div className="flex gap-2">
              <button className="px-3 py-1.5 bg-primary text-primary-foreground rounded-md text-sm">
                {t('common.primary')}
              </button>
              <button className="px-3 py-1.5 bg-secondary text-secondary-foreground rounded-md text-sm">
                {t('common.secondary')}
              </button>
            </div>
          </div>

          <div className="p-6 bg-muted rounded-lg border border-border space-y-3">
            <h5 className="font-medium text-muted-foreground">
              {t('settings.appearance.previewMuted')}
            </h5>
            <p className="text-sm text-muted-foreground">
              {t('settings.appearance.previewMutedDesc')}
            </p>
            <div className="flex gap-2">
              <button className="px-3 py-1.5 border border-border rounded-md text-sm text-foreground hover:bg-accent">
                {t('common.outline')}
              </button>
              <button className="px-3 py-1.5 text-sm text-foreground hover:bg-accent rounded-md">
                {t('common.ghost')}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}