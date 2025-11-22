"use client";

import { SearchResult } from "@/lib/api";
import { Mail, Calendar, User } from "lucide-react";
import { useRouter } from "next/navigation";

interface SearchResultsProps {
    results: SearchResult[];
    isLoading: boolean;
    query: string;
    onClose: () => void;
}

export function SearchResults({ results, isLoading, query, onClose }: SearchResultsProps) {
    const router = useRouter();

    const handleResultClick = (emailId: string) => {
        router.push(`/dashboard/emails/${emailId}`);
        onClose();
    };

    if (isLoading) {
        return (
            <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                <div className="p-8 text-center">
                    <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                    <p className="mt-3 text-sm text-slate-500">Searching...</p>
                </div>
            </div>
        );
    }

    if (!results || results.length === 0) {
        return (
            <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50">
                <div className="p-8 text-center">
                    <Mail className="w-12 h-12 text-slate-300 mx-auto mb-3" />
                    <p className="text-sm font-medium text-slate-600">No results found</p>
                    <p className="text-xs text-slate-400 mt-1">Try a different search query</p>
                </div>
            </div>
        );
    }

    return (
        <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in slide-in-from-top-2 fade-in duration-200 z-50 max-h-96 overflow-y-auto">
            <div className="p-2 border-b border-slate-100 bg-slate-50">
                <p className="text-xs font-medium text-slate-600 px-3 py-1">
                    Found {results.length} result{results.length > 1 ? 's' : ''} for "{query}"
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
                                    {result.subject || "(No Subject)"}
                                </h4>
                                <p className="text-xs text-slate-500 line-clamp-2 mt-1">
                                    {result.snippet}
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
                                <div className="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-1 rounded">
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
