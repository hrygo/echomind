'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useAuthStore, useOnboardingStore } from '@/store';
import { AuthLayout } from '@/components/auth/AuthLayout';
import { RoleSelector } from '@/components/onboarding/RoleSelector';
import { SmartMailboxForm } from '@/components/onboarding/SmartMailboxForm';
import { InitialSync } from '@/components/onboarding/InitialSync';

export default function OnboardingPageClient() {
    const router = useRouter();
    const { t } = useLanguage();
    const isAuthenticated = useAuthStore(state => state.isAuthenticated);
    const { step, resetOnboarding } = useOnboardingStore();

    // Redirect if not authenticated
    useEffect(() => {
        if (!isAuthenticated) {
            router.push('/auth?mode=login');
        }
    }, [isAuthenticated, router]);

    // Clear onboarding state on unmount or if leaving page early
    useEffect(() => {
        return () => {
            resetOnboarding();
        };
    }, [resetOnboarding]);

    if (!isAuthenticated) {
        return null; // Render nothing while redirecting
    }

    let content;
    switch (step) {
        case 1:
            content = <RoleSelector />;
            break;
        case 2:
            content = <SmartMailboxForm />;
            break;
        case 3:
            content = <InitialSync />;
            break;
        case 4:
            // Future WeChat binding step
            content = (
                <div className="text-center p-8">
                    <h2 className="text-2xl font-bold mb-4">{t('onboarding.step4.title')}</h2>
                    <p className="text-slate-600">{t('onboarding.step4.subtitle')}</p>
                </div>
            );
            break;
        default:
            content = (
                <div className="text-center p-8">
                    <h2 className="text-2xl font-bold text-red-500">Error: Unknown Onboarding Step</h2>
                    <p className="text-slate-600">Please restart the onboarding process.</p>
                </div>
            );
    }

    return (
        <AuthLayout>
            {content}
        </AuthLayout>
    );
}
