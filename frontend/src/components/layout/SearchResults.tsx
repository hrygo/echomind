"use client";

import { SearchResult } from "@/lib/api";
import { Mail, Calendar, User } from "lucide-react";
import { useRouter } from "next/navigation";
import { Skeleton } from "@/components/ui/Skeleton";
import { useLanguage } from "@/lib/i18n/LanguageContext";

interface SearchResultsProps {
    results: SearchResult[];
    isLoading: boolean;
    error: string | null;
    query: string;
    onClose: () => void;
}

export function SearchResults({ results, isLoading, error, query, onClose }: SearchResultsProps) {
    const router = useRouter();
    const { t } = useLanguage();

    const handleResultClick = (emailId: string) => {
        router.push(`/dashboard/emails/${emailId}`);
        onClose();
    };

    const highlightText = (text: string, query: string) => {
        if (!query.trim()) return text;
        const words = query.trim().split(/\s+/).filter(word => word.length > 0);
        if (words.length === 0) return text;

        const pattern = new RegExp(`(${words.join('|')})`, 'gi');
        const parts = text.split(pattern);

        return parts.map((part, i) => 
            pattern.test(part) ? <span key={i} className="bg-yellow-200 text-slate-900 font-medium rounded-sm px-0.5">{part}</span> : part
        );
    };

    if (isLoading) {
        return (
            <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                <div className="p-2 border-b border-slate-100 bg-slate-50">
                    <Skeleton className="h-4 w-32" />
                </div>
                <div className="divide-y divide-slate-100">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="px-4 py-3">
                            <div className="flex items-start justify-between gap-3">
                                <div className="flex-1 space-y-2">
                                    <Skeleton className="h-4 w-3/4" />
                                    <Skeleton className="h-3 w-full" />
                                    <div className="flex gap-3 pt-1">
                                        <Skeleton className="h-3 w-20" />
                                        <Skeleton className="h-3 w-20" />
                                    </div>
                                </div>
                                <Skeleton className="h-6 w-10 rounded" />
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-red-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                <div className="p-8 text-center">
                    <div className="text-red-500 mb-2">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mx-auto"><circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/></svg>
                    </div>
                    <p className="text-sm font-medium text-slate-800">{t('common.error')}</p>
                    <p className="text-xs text-slate-500 mt-1">{error}</p>
                </div>
            </div>
        );
    }

    if (!results || results.length === 0) {
        return (
            <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                <div className="p-8 text-center">
                    <Mail className="w-12 h-12 text-slate-300 mx-auto mb-3" />
                    <p className="text-sm font-medium text-slate-600">{t('common.noResults')}</p>
                    <p className="text-xs text-slate-400 mt-1 mb-4">{t('common.noResultsDesc')}</p>
                    
                    <div className="text-left bg-slate-50 rounded-lg p-3">
                        <p className="text-xs font-semibold text-slate-500 mb-2">{t('common.searchTips')}</p>
                        <ul className="text-xs text-slate-500 space-y-1 list-disc pl-4">
                            <li>{t('common.searchTip1')}</li>
                            <li>{t('common.searchTip2')}</li>
                            <li>{t('common.searchTip3')}</li>
                            <li>{t('common.searchTip4')}</li>
                        </ul>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50 max-h-96 overflow-y-auto">
            <div className="p-2 border-b border-slate-100 bg-slate-50">
                <p className="text-xs font-medium text-slate-600 px-3 py-1">
                    {t('common.foundResults')
                        .replace('{count}', results.length.toString())
                        .replace('{plural}', results.length > 1 ? 's' : '')
                        .replace('{query}', query)
                    }
                </p>
            </div>
            <div className="divide-y divide-slate-100">
                {results.map((result) => (
                    <button
                        key={result.email_id}
                        onClick={() => handleResultClick(result.email_id)}
                        className="w-full text-left px-4 py-3 hover:bg-slate-50 transition-colors focus:outline-none focus:bg-slate-50 group"
                    >
                        <div className="flex items-start justify-between gap-3">
                            <div className="flex-1 min-w-0">
                                <h4 className="text-sm font-semibold text-slate-800 truncate group-hover:text-blue-600 transition-colors">
                                    {highlightText(result.subject || "(No Subject)", query)}
                                </h4>
                                <p className="text-xs text-slate-500 line-clamp-2 mt-1">
                                    {highlightText(result.snippet, query)}
                                </p>
                                <div className="flex items-center gap-3 mt-2 text-xs text-slate-400">
                                    <span className="flex items-center gap-1">
                                        <User className="w-3 h-3" />
                                        {result.sender}
                                    </span>
                                    <span className="flex items-center gap-1">
                                        <Calendar className="w-3 h-3" />
                                        {new Date(result.date).toLocaleDateString()}
                                    </span>
                                </div>
                            </div>
                            <div className="flex-shrink-0">
                                <div 
                                    className="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-1 rounded cursor-help"
                                    title={`Relevance Score: ${(result.score * 100).toFixed(1)}%\nCalculated using semantic vector similarity.`}
                                >
                                    {(result.score * 100).toFixed(0)}%
                                </div>
                            </div>
                        </div>
                    </button>
                ))}
            </div>
        </div>
    );
}
