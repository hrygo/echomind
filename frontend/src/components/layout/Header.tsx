"use client";

import { useState, useRef } from "react";
import Link from "next/link";
import { Bell, Settings, LogOut, Globe, Sparkles, Menu } from "lucide-react";
import { ThemeToggle } from "@/components/theme/ThemeToggle";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { useUIStore } from "@/store/ui";
import { CopilotWidget } from "@/components/copilot/CopilotWidget";
import { useCopilotStore } from "@/store/useCopilotStore";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";
import { useOnClickOutside } from "@/hooks/useOnClickOutside";

export function Header() {
    const { language, setLanguage, t } = useLanguage();
    const { openMobileSidebar } = useUIStore();
    const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

    const { setMode, setIsOpen } = useCopilotStore();
    const user = useAuthStore((state) => state.user);
    const logout = useAuthStore((state) => state.logout);
    const router = useRouter();

    const dropdownRef = useRef<HTMLDivElement>(null);
    useOnClickOutside<HTMLDivElement>(dropdownRef, () => setIsUserMenuOpen(false));

    const toggleLanguage = () => {
        setLanguage(language === 'zh' ? 'en' : 'zh');
    };

    const handleOpenChat = () => {
        setMode('chat');
        setIsOpen(true);
    };

    const handleLogout = () => {
        logout();
        setIsUserMenuOpen(false); // Close the dropdown
        router.push('/auth');    // Redirect to the auth page
    };

    return (
        <header className="h-16 md:h-20 bg-white/80 backdrop-blur-xl border-b border-slate-200/60 flex items-center justify-between px-4 md:px-8 sticky top-0 z-30 transition-all duration-200">
            {/* Mobile Header (Simplified for now) */}
            <div className="flex md:hidden items-center w-full justify-between">
                <div className="flex items-center gap-3">
                    <button
                        onClick={openMobileSidebar}
                        className="p-2 -ml-2 text-slate-500 hover:bg-slate-100 rounded-lg"
                    >
                        <Menu className="w-6 h-6" />
                    </button>
                    <span className="font-bold text-lg text-slate-800">EchoMind</span>
                </div>
                <div className="flex items-center gap-2">
                    {/* Mobile Search/Chat Trigger */}
                    <button
                        onClick={handleOpenChat}
                        className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-full"
                    >
                        <Sparkles className="w-5 h-5" />
                    </button>
                    <div ref={dropdownRef} className="relative">
                        <button
                            onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                            className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-xs shadow-sm"
                        >
                            {user?.name?.[0].toUpperCase() || 'D'}
                        </button>
                        {isUserMenuOpen && (
                            // Dropdown menu content remains the same, just for mobile
                            // This is a simplified version for brevity in this example
                            <div className="absolute top-full right-0 mt-2 w-56 bg-white rounded-xl shadow-xl border border-slate-100 z-50">
                                <div className="p-1">
                                    <button
                                        className="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 rounded-lg text-left transition-colors"
                                        onClick={handleLogout}
                                    >
                                        <LogOut className="w-4 h-4" />
                                        {t('common.logout')}
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* Desktop: Standard Layout */}
            <div className="hidden md:flex flex-1 items-center justify-between w-full">
                {/* Copilot Widget (The Omni-Bar) */}
                <div className="flex-1 max-w-2xl relative mr-4">
                    <CopilotWidget />
                </div>

                {/* Right Actions */}
                <div className="flex items-center gap-4 ml-auto shrink-0">
                    <ThemeToggle />
                    <button
                        onClick={toggleLanguage}
                        className="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-colors"
                        title="Switch Language"
                    >
                        <Globe className="w-5 h-5" />
                    </button>

                    <button
                        onClick={handleOpenChat}
                        className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-full transition-colors"
                        title="AI Copilot"
                    >
                        <Sparkles className="w-5 h-5" />
                    </button>

                    <button className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-full transition-colors relative hidden md:block">
                        <Bell className="w-5 h-5" />
                        <span className="absolute top-2 right-2.5 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
                    </button>

                    <div className="h-6 w-px bg-slate-200 mx-1"></div>

                    {/* User Profile Dropdown */}
                    <div ref={dropdownRef} className="relative">
                        <button
                            onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                            className="flex items-center gap-3 p-1.5 pr-3 rounded-full hover:bg-slate-100 transition-all duration-200 focus:outline-none"
                        >
                            <div className="w-9 h-9 rounded-full bg-gradient-to-tr from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-blue-200">
                                {user?.name?.[0].toUpperCase() || 'D'}
                            </div>
                            <div className="hidden lg:block text-left">
                                <p className="text-sm font-semibold text-slate-700 leading-none">{user?.name || '演示用户'}</p>
                                <p className="text-[10px] text-slate-400 font-medium mt-1">{t('sidebar.freePlan')}</p>
                            </div>
                        </button>

                        {/* Dropdown Menu */}
                        {isUserMenuOpen && (
                            <div className="absolute top-full right-0 mt-2 w-56 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                                <div className="p-2 border-b border-slate-50">
                                    <div className="px-3 py-2">
                                        <p className="text-sm font-semibold text-slate-800">{user?.name || 'User Name'}</p>
                                        <p className="text-xs text-slate-500">{user?.email || 'user@example.com'}</p>
                                    </div>
                                </div>
                                <div className="p-1">
                                    <Link
                                        href="/dashboard/settings"
                                        className="flex items-center gap-2 px-3 py-2 text-sm text-slate-600 hover:bg-slate-50 rounded-lg transition-colors"
                                        onClick={() => setIsUserMenuOpen(false)}
                                    >
                                        <Settings className="w-4 h-4" />
                                        {t('common.settings')}
                                    </Link>
                                    <button
                                        className="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 rounded-lg text-left transition-colors"
                                        onClick={handleLogout}
                                    >
                                        <LogOut className="w-4 h-4" />
                                        {t('common.logout')}
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </header>
    );
}