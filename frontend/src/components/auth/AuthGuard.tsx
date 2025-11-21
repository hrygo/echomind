'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '@/store/auth';

export default function AuthGuard({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const { isAuthenticated, token, isHydrated } = useAuthStore();
    const [mounted, setMounted] = useState(false);
    const [isChecking, setIsChecking] = useState(true);

    useEffect(() => {
        setTimeout(() => setMounted(true), 0);
    }, []);

    useEffect(() => {
        if (!mounted || !isHydrated) return;

        // Define public paths that don't require authentication
        const publicPaths = ['/login', '/register', '/'];
        const isPublicPath = publicPaths.includes(pathname);

        if (isPublicPath) {
            // Allow access to public pages
            setTimeout(() => setIsChecking(false), 0);
        } else if (!isAuthenticated || !token) {
            // Redirect to login if trying to access protected page without auth
            router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
        } else {
            // User is authenticated and accessing protected page
            setTimeout(() => setIsChecking(false), 0);
        }
    }, [isAuthenticated, token, router, pathname, mounted, isHydrated]);

    if (!mounted || !isHydrated || isChecking) {
        // Show loading spinner while checking auth or hydrating
        return (
            <div className="flex min-h-screen items-center justify-center">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
            </div>
        );
    }

    return <>{children}</>;
}
