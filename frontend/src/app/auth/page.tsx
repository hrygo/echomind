import React, { Suspense } from 'react';
import AuthPageClient from './AuthPageClient';

// Force dynamic rendering to prevent SSR/prerendering
export const dynamic = 'force-dynamic';

export default function AuthPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">Loading...</div>}>
      <AuthPageClient />
    </Suspense>
  );
}
