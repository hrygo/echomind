"use client";

import { useState } from "react";
import Link from "next/link";
import { Search, Bell, Settings, LogOut, Globe, Filter } from "lucide-react";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { searchEmails, SearchResult } from "@/lib/api";
import { SearchResults } from "./SearchResults";
import { SearchHistory } from "./SearchHistory";
import { SearchFilters } from "./SearchFilters";
import { useSearchHistory } from "@/hooks/useSearchHistory";

export function Header() {
    const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
    const [isSearching, setIsSearching] = useState(false);
    const [searchError, setSearchError] = useState<string | null>(null);
    const [showResults, setShowResults] = useState(false);
    const [isInputFocused, setIsInputFocused] = useState(false);
    const [showFilters, setShowFilters] = useState(false);
    const [filters, setFilters] = useState({ sender: '', startDate: '', endDate: '' });
    const { history, addToHistory, removeFromHistory, clearHistory } = useSearchHistory();
    const { language, setLanguage, t } = useLanguage();

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
            const response = await searchEmails(query, filters, 10);
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
    };

    return (
        <header className="h-20 bg-white/80 backdrop-blur-xl border-b border-slate-200/60 flex items-center justify-between px-8 sticky top-0 z-30 transition-all duration-200">
            {/* Search Bar */}
            <div className="flex-1 max-w-xl relative flex items-center gap-2">
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
                            // Delay to allow clicking on history items
                            setTimeout(() => setIsInputFocused(false), 200);
                        }}
                        className="block w-full pl-10 pr-3 py-2.5 border-none rounded-xl bg-slate-100 text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:bg-white transition-all duration-200 text-sm font-medium"
                        placeholder={t('common.searchPlaceholder')}
                    />
                </div>
                
                <div className="relative">
                    <button
                        onClick={() => setShowFilters(!showFilters)}
                        className={`p-2.5 rounded-xl transition-colors ${
                            showFilters || filters.sender || filters.startDate || filters.endDate
                                ? 'bg-blue-50 text-blue-600 border border-blue-100' 
                                : 'bg-slate-100 text-slate-400 hover:text-slate-600 hover:bg-slate-200'
                        }`}
                        title="Search Filters"
                    >
                        <Filter className="w-4 h-4" />
                        {(filters.sender || filters.startDate || filters.endDate) && (
                            <span className="absolute top-2 right-2 w-1.5 h-1.5 bg-blue-600 rounded-full"></span>
                        )}
                    </button>
                    {showFilters && (
                        <SearchFilters
                            filters={filters}
                            onChange={setFilters}
                            onClose={() => setShowFilters(false)}
                        />
                    )}
                </div>

                {showResults ? (
                    <SearchResults
                        results={searchResults}
                        isLoading={isSearching}
                        error={searchError}
                        query={searchQuery}
                        onClose={handleCloseResults}
                    />
                ) : (isInputFocused && !searchQuery && history.length > 0) && (
                    <SearchHistory
                        history={history}
                        onSelect={(query) => handleSearch(query)}
                        onRemove={removeFromHistory}
                        onClear={clearHistory}
                    />
                )}
            </div>

            {/* Right Actions */}
            <div className="flex items-center gap-4 ml-6">
                <button
                    onClick={toggleLanguage}
                    className="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-colors"
                    title="Switch Language"
                >
                    <Globe className="w-5 h-5" />
                </button>

                <button className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-full transition-colors relative">
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
                        <div className="hidden md:block text-left">
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
        </header>
    );
}
