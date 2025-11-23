'use client';

import { useEffect, useState, useMemo } from 'react';

import { useRouter, usePathname } from 'next/navigation';

import { useAuthStore, useOnboardingStore } from '@/store'; // Import useAuthStore and useOnboardingStore from centralized store



export default function AuthGuard({ children }: { children: React.ReactNode }) {

    const router = useRouter();

    const pathname = usePathname();

    const { isAuthenticated, token, isHydrated, user } = useAuthStore(); // Destructure user

    const { resetOnboarding } = useOnboardingStore();

    const [mounted, setMounted] = useState(false);

    const [isChecking, setIsChecking] = useState(true);



    // Memoize path checks for stable dependencies

    const publicPaths = useMemo(() => ['/auth', '/'], []);

    const isPublicPath = useMemo(() => publicPaths.includes(pathname), [publicPaths, pathname]);

    const isOnboardingPath = useMemo(() => pathname === '/onboarding', [pathname]);



    useEffect(() => {

        setTimeout(() => setMounted(true), 0);

    }, []);



        useEffect(() => {



            if (!mounted || !isHydrated) return;



    



            if (isPublicPath) {



                setTimeout(() => setIsChecking(false), 0); // Defer setState



                if (isAuthenticated && pathname === '/auth') { // Only redirect from /auth if authenticated



                  router.push('/dashboard');



                }



            } else if (!isAuthenticated || !token) {



                router.push(`/auth?mode=login&redirect=${encodeURIComponent(pathname)}`);



            } else {



                // User is authenticated, now check onboarding status



                if (user && !user.has_account && !isOnboardingPath) {



                  router.push('/onboarding');



                } else {



                  setTimeout(() => setIsChecking(false), 0); // Defer setState



                }



            }



        }, [isAuthenticated, token, router, pathname, mounted, isHydrated, user, isPublicPath, isOnboardingPath]);



    useEffect(() => {

      // Reset onboarding state if user navigates away from onboarding or is redirected to login

      // Make sure `isPublicPath` and `isOnboardingPath` are stable via useMemo

      if (!isAuthenticated || (isPublicPath && !isOnboardingPath)) { // Changed condition

        resetOnboarding();

      }

    }, [isAuthenticated, isOnboardingPath, isPublicPath, resetOnboarding]);



    if (!mounted || !isHydrated || isChecking) {

        return (

            <div className="flex min-h-screen items-center justify-center">

                <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>

            </div>

        );

    }



    return <>{children}</>;

}


