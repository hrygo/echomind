"use client";

import { useState } from "react";
import Link from "next/link";
import { Search, Bell, Settings, LogOut, Globe, Sparkles, Menu, ArrowLeft, X } from "lucide-react";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { searchEmails, SearchResult } from "@/lib/api";
import { SearchResults } from "./SearchResults";
import { SearchHistory } from "./SearchHistory";
import { useSearchHistory } from "@/hooks/useSearchHistory";
import { useChatStore } from "@/lib/store/chat";
import { useUIStore } from "@/store/ui";

export function Header() {
    const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
    const [isMobileSearchOpen, setIsMobileSearchOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
    const [isSearching, setIsSearching] = useState(false);
    const [searchError, setSearchError] = useState<string | null>(null);
    const [showResults, setShowResults] = useState(false);
    const [isInputFocused, setIsInputFocused] = useState(false);
    const { history, addToHistory, removeFromHistory, clearHistory } = useSearchHistory();
    const { language, setLanguage, t } = useLanguage();
    const { toggleOpen } = useChatStore();
    const { openMobileSidebar } = useUIStore();

    const toggleLanguage = () => {
        setLanguage(language === 'zh' ? 'en' : 'zh');
    };

    const handleSearch = async (queryOverride?: string) => {
        const query = typeof queryOverride === 'string' ? queryOverride : searchQuery;
        if (!query.trim()) {
            setShowResults(false);
            return;
        }

        if (queryOverride) setSearchQuery(queryOverride);

        addToHistory(query);
        setIsSearching(true);
        setSearchError(null);
        setShowResults(true);

        try {
            // filters parameter is removed
            const response = await searchEmails(query, 10);
            setSearchResults(response.results);
        } catch (error) {
            console.error("Search failed:", error);
            setSearchError("Failed to search emails. Please try again.");
            setSearchResults([]);
        } finally {
            setIsSearching(false);
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    const handleCloseResults = () => {
        setShowResults(false);
        setSearchQuery("");
        // On mobile, closing results might also mean closing the search bar, but usually we keep it open until user hits back
    };

    return (
        <header className="h-16 md:h-20 bg-white/80 backdrop-blur-xl border-b border-slate-200/60 flex items-center justify-between px-4 md:px-8 sticky top-0 z-30 transition-all duration-200">
            
            {/* Mobile: Search Overlay Mode */}
            {isMobileSearchOpen ? (
                <div className="absolute inset-0 bg-white z-40 flex items-center px-4 gap-2 md:hidden animate-in fade-in slide-in-from-top-2 duration-200">
                    <button 
                        onClick={() => setIsMobileSearchOpen(false)}
                        className="p-2 -ml-2 text-slate-500 hover:bg-slate-100 rounded-full"
                    >
                        <ArrowLeft className="w-5 h-5" />
                    </button>
                    <div className="flex-1 relative group">
                        <input
                            autoFocus
                            type="text"
                            value={searchQuery}
                            onChange={(e) => {
                                setSearchQuery(e.target.value);
                                if (!e.target.value) setShowResults(false);
                            }}
                            onKeyPress={handleKeyPress}
                            className="block w-full pl-4 pr-10 py-2 bg-slate-100 border-none rounded-full text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:bg-white transition-all text-sm"
                            placeholder={t('common.searchPlaceholder')}
                        />
                        {searchQuery && (
                            <button 
                                onClick={() => setSearchQuery('')}
                                className="absolute right-3 top-2.5 text-slate-400"
                            >
                                <X className="w-4 h-4" />
                            </button>
                        )}
                    </div>
                </div>
            ) : (
                // Mobile: Default Mode
                <div className="flex md:hidden items-center w-full justify-between">
                    <div className="flex items-center gap-3">
                        <button
                            onClick={openMobileSidebar}
                            className="p-2 -ml-2 text-slate-500 hover:bg-slate-100 rounded-lg"
                        >
                            <Menu className="w-6 h-6" />
                        </button>
                        {/* Optional: Mobile Logo Text */}
                        <span className="font-bold text-lg text-slate-800">EchoMind</span>
                    </div>
                    <div className="flex items-center gap-2">
                        <button 
                            onClick={() => setIsMobileSearchOpen(true)}
                            className="p-2 text-slate-500 hover:bg-slate-100 rounded-full"
                        >
                            <Search className="w-5 h-5" />
                        </button>
                        <button
                            onClick={toggleOpen}
                            className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-full"
                        >
                            <Sparkles className="w-5 h-5" />
                        </button>
                        <button
                            onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                            className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-xs shadow-sm"
                        >
                            U
                        </button>
                    </div>
                </div>
            )}

            {/* Desktop: Standard Layout (Hidden on Mobile) */}
            <div className="hidden md:flex flex-1 items-center justify-between w-full">
                {/* Search Bar */}
                <div className="flex-1 max-w-2xl relative flex items-center gap-2 mr-4">
                    <div className="relative group flex-1">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <Search className="h-4 w-4 text-slate-400 group-focus-within:text-blue-500 transition-colors" />
                        </div>
                        <input
                            type="text"
                            value={searchQuery}
                            onChange={(e) => {
                                setSearchQuery(e.target.value);
                                if (!e.target.value) setShowResults(false);
                            }}
                            onKeyPress={handleKeyPress}
                            onFocus={() => {
                                setIsInputFocused(true);
                                if (searchQuery) setShowResults(true);
                            }}
                            onBlur={() => {
                                setTimeout(() => setIsInputFocused(false), 200);
                            }}
                            className="block w-full pl-10 pr-3 py-2.5 border-none rounded-xl bg-slate-100 text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:bg-white transition-all duration-200 text-sm font-medium"
                            placeholder={t('common.searchPlaceholder')}
                        />
                    </div>
                </div>

                {/* Right Actions */}
                <div className="flex items-center gap-4 ml-auto shrink-0">
                    <button
                        onClick={toggleLanguage}
                        className="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-colors"
                        title="Switch Language"
                    >
                        <Globe className="w-5 h-5" />
                    </button>

                    <button
                        onClick={toggleOpen}
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
                    <div className="relative">
                        <button
                            onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                            className="flex items-center gap-3 p-1.5 pr-3 rounded-full hover:bg-slate-100 transition-all duration-200 focus:outline-none"
                        >
                            <div className="w-9 h-9 rounded-full bg-gradient-to-tr from-blue-500 to-cyan-500 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-blue-200">
                                U
                            </div>
                            <div className="hidden lg:block text-left">
                                <p className="text-sm font-semibold text-slate-700 leading-none">{t('sidebar.user')}</p>
                                <p className="text-[10px] text-slate-400 font-medium mt-1">{t('sidebar.freePlan')}</p>
                            </div>
                        </button>

                        {/* Dropdown Menu */}
                        {isUserMenuOpen && (
                            <div className="absolute top-full right-0 mt-2 w-56 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                                <div className="p-2 border-b border-slate-50">
                                    <div className="px-3 py-2">
                                        <p className="text-sm font-semibold text-slate-800">User Name</p>
                                        <p className="text-xs text-slate-500">user@example.com</p>
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
                                        onClick={() => alert("Logout clicked")}
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
            
            {/* Mobile Search Results / History Overlay */}
            {(isMobileSearchOpen && (showResults || (history.length > 0 && !searchQuery))) && (
                <div className="absolute top-16 left-0 right-0 bg-white border-t border-slate-100 shadow-xl z-30 max-h-[80vh] overflow-y-auto md:hidden">
                     {showResults ? (
                        <SearchResults
                            results={searchResults}
                            isLoading={isSearching}
                            error={searchError}
                            query={searchQuery}
                            onClose={handleCloseResults}
                        />
                    ) : (
                        <SearchHistory
                            history={history}
                            onSelect={(query) => handleSearch(query)}
                            onRemove={removeFromHistory}
                            onClear={clearHistory}
                        />
                    )}
                </div>
            )}
        </header>
    );
}
