'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '@/store/auth';

export default function AuthGuard({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const { isAuthenticated, token } = useAuthStore();
    const [isChecking, setIsChecking] = useState(true);

    useEffect(() => {
        // Define public paths that don't require authentication
        const publicPaths = ['/login', '/register', '/'];
        const isPublicPath = publicPaths.includes(pathname);

        if (isPublicPath) {
            // Allow access to public pages
            setIsChecking(false);
        } else if (!isAuthenticated || !token) {
            // Redirect to login if trying to access protected page without auth
            router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
        } else {
            // User is authenticated and accessing protected page
            setIsChecking(false);
        }
    }, [isAuthenticated, token, router, pathname]);

    if (isChecking) {
        // Show loading spinner while checking auth
        return (
            <div className="flex min-h-screen items-center justify-center">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
            </div>
        );
    }

    return <>{children}</>;
}
