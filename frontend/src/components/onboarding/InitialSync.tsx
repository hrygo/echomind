'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useOnboardingStore, useAuthStore } from '@/store';
import { api } from '@/lib/api';
import { isAxiosError } from 'axios';
import { Loader2, CheckCircle2, AlertCircle } from 'lucide-react';

export function InitialSync() {
  const { t } = useLanguage();
  const router = useRouter();
  const { mailbox, role, setStep } = useOnboardingStore();
  const updateUser = useAuthStore(state => state.updateUser);
  const [syncStatus, setSyncStatus] = useState('pending'); // pending, success, failed
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    const initiateSyncAndSaveProfile = async () => {
      setSyncStatus('pending');
      setErrorMessage(null);

      try {
        // 1. Update user role (from Step 1)
        if (role) {
          await api.patch('/users/me', { role });
        }

        // 2. Save mailbox config (from Step 2) if not already done via SmartMailboxForm's handleSubmit
        // This part is a fallback/double-check; SmartMailboxForm already attempts to save.
        // If SmartMailboxForm successfully saved, `ConnectAndSaveAccount` would just update.
        // So we can just ensure the save happens here if not explicitly saved before.
        // Or more cleanly, ensure the onboarding store has a flag `isMailboxSaved`.
        // For now, let's assume it was saved in Step 2, and just trigger a sync.

        // 3. Trigger initial email sync
        await api.post<{ message: string }>("/sync");

        // 4. Update user state to mark account as connected
        updateUser({ has_account: true });

        setSyncStatus('success');
        setTimeout(() => {
          router.push('/dashboard');
          setStep(1); // Reset onboarding state when navigating away
        }, 2000); // Wait 2 seconds for user to see success message

      } catch (error: unknown) {
        console.error("Initial sync failed:", error);
        setSyncStatus('failed');
        if (isAxiosError(error) && error.response) {
          setErrorMessage(error.response.data?.error || t('onboarding.step3.errors.unknownSync'));
        } else if (error instanceof Error) {
          setErrorMessage(error.message);
        } else {
          setErrorMessage(t('onboarding.step3.errors.unknownError'));
        }
      }
    };

    initiateSyncAndSaveProfile();
  }, [router, mailbox, role, setStep, t, updateUser]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen-auth p-4 text-center">
      <h1 className="text-4xl font-extrabold text-slate-900 mb-4">
        {t('onboarding.step3.title')}
      </h1>
      <p className="text-lg text-slate-600 mb-10 max-w-xl">
        {t('onboarding.step3.subtitle')}
      </p>

      <div className="relative w-24 h-24 mb-8">
        {syncStatus === 'pending' && (
          <Loader2 className="w-full h-full text-blue-600 animate-spin absolute inset-0" />
        )}
        {syncStatus === 'success' && (
          <CheckCircle2 className="w-full h-full text-green-600 absolute inset-0 animate-in fade-in zoom-in-95" />
        )}
        {syncStatus === 'failed' && (
          <AlertCircle className="w-full h-full text-red-600 absolute inset-0 animate-in fade-in zoom-in-95" />
        )}
      </div>

      {errorMessage && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm flex items-center gap-2 mb-4">
          <AlertCircle className="w-5 h-5 flex-shrink-0" />
          {errorMessage}
        </div>
      )}

      {syncStatus === 'failed' && (
        <button
          onClick={() => router.push('/onboarding?step=2')} // Go back to mailbox config
          className="px-8 py-3 rounded-full bg-blue-600 hover:bg-blue-700 text-white font-semibold text-lg shadow-lg transition-all duration-200 mt-6"
        >
          {t('onboarding.step3.retry')} {/* Need to add retry key */}
        </button>
      )}

      {syncStatus === 'success' && (
        <button
          onClick={() => {
            router.push('/dashboard');
            setStep(1); // Reset onboarding state
          }}
          className="px-8 py-3 rounded-full bg-blue-600 hover:bg-blue-700 text-white font-semibold text-lg shadow-lg transition-all duration-200 mt-6"
        >
          {t('onboarding.step3.enterDashboard')}
        </button>
      )}
    </div>
  );
}
