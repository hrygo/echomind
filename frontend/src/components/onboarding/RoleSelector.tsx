'use client';

import React from 'react';
import { Briefcase, Users, Target, ArrowRight } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useOnboardingStore } from '@/store';
import { cn } from '@/lib/utils';

export function RoleSelector() {
  const { t } = useLanguage();
  const { role, setRole, setStep } = useOnboardingStore();

  const roles = [
    {
      id: 'executive',
      icon: Briefcase,
      title: t('onboarding.step1.roleExecutive'),
      description: t('onboarding.step1.descExecutive'),
    },
    {
      id: 'manager',
      icon: Users,
      title: t('onboarding.step1.roleManager'),
      description: t('onboarding.step1.descManager'),
    },
    {
      id: 'dealmaker',
      icon: Target,
      title: t('onboarding.step1.roleDealmaker'),
      description: t('onboarding.step1.descDealmaker'),
    },
  ];

  const handleNext = () => {
    if (role) {
      setStep(2);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen-auth p-4">
      <h1 className="text-4xl font-extrabold text-slate-900 mb-4 text-center">
        {t('onboarding.step1.title')}
      </h1>
      <p className="text-lg text-slate-600 mb-10 text-center max-w-xl">
        {t('onboarding.step1.subtitle')}
      </p>

      <div className="grid md:grid-cols-3 gap-6 w-full max-w-4xl mb-12">
        {roles.map((r) => (
          <button
            key={r.id}
            onClick={() => setRole(r.id as 'executive' | 'manager' | 'dealmaker')}
            className={cn(
              "flex flex-col items-center p-6 border rounded-xl shadow-sm transition-all duration-200",
              "hover:shadow-md hover:border-blue-500/50",
              role === r.id ? "border-blue-500 ring-2 ring-blue-100 shadow-md" : "border-slate-200 bg-white"
            )}
          >
            <r.icon className={cn("w-12 h-12 mb-4 transition-colors duration-200", role === r.id ? "text-blue-600" : "text-slate-500")} />
            <h3 className="text-xl font-semibold text-slate-800 mb-2">{r.title}</h3>
            <p className="text-sm text-slate-600 text-center">{r.description}</p>
          </button>
        ))}
      </div>

      <button
        onClick={handleNext}
        disabled={!role}
        className={cn(
          "px-8 py-3 rounded-full text-white font-semibold text-lg shadow-lg transition-all duration-200",
          "bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed",
          role ? "opacity-100" : "opacity-0 pointer-events-none"
        )}
      >
        {t('onboarding.step1.next')} <ArrowRight className="w-5 h-5 ml-2 inline-block" />
      </button>
    </div>
  );
}
