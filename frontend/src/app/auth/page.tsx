'use client';

import React, { useState, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { AuthLayout } from '@/components/auth/AuthLayout';
import { AuthForm } from '@/components/auth/AuthForm';
import { useAuthStore } from '@/store/auth';

export default function AuthPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialMode = searchParams.get('mode') === 'register' ? 'register' : 'login';
  const [mode, setMode] = useState<'login' | 'register'>(initialMode);
  const isAuthenticated = useAuthStore(state => state.isAuthenticated);

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push('/dashboard');
    }
  }, [isAuthenticated, router]);

  useEffect(() => {
    // Update URL query param when mode changes
    const newSearchParams = new URLSearchParams(searchParams.toString());
    newSearchParams.set('mode', mode);
    router.replace(`/auth?${newSearchParams.toString()}`, { shallow: true });
  }, [mode, router, searchParams]);

  const handleModeChange = (newMode: 'login' | 'register') => {
    setMode(newMode);
  };

  if (isAuthenticated) {
    return null; // Render nothing while redirecting
  }

  return (
    <AuthLayout>
      <AuthForm mode={mode} onModeChange={handleModeChange} />
    </AuthLayout>
  );
}
