import { useState, useEffect } from "react";
import { Filter, Calendar, User, X } from "lucide-react";

interface SearchFiltersProps {
    filters: {
        sender: string;
        startDate: string;
        endDate: string;
    };
    onChange: (filters: { sender: string; startDate: string; endDate: string }) => void;
    onClose: () => void;
}

export function SearchFilters({ filters, onChange, onClose }: SearchFiltersProps) {
    const [localFilters, setLocalFilters] = useState(filters);

    useEffect(() => {
        setLocalFilters(filters);
    }, [filters]);

    const handleChange = (key: string, value: string) => {
        const newFilters = { ...localFilters, [key]: value };
        setLocalFilters(newFilters);
        onChange(newFilters);
    };

    return (
        <div className="absolute top-full right-0 mt-2 w-72 bg-white rounded-xl shadow-xl border border-slate-100 p-4 animate-in slide-in-from-top-2 fade-in duration-200 z-50">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-slate-800 flex items-center gap-2">
                    <Filter className="w-4 h-4 text-slate-500" />
                    Search Filters
                </h3>
                <button 
                    onClick={onClose}
                    className="text-slate-400 hover:text-slate-600 p-1 hover:bg-slate-100 rounded-full transition-colors"
                >
                    <X className="w-4 h-4" />
                </button>
            </div>

            <div className="space-y-4">
                {/* Sender Filter */}
                <div>
                    <label className="block text-xs font-medium text-slate-500 mb-1.5">Sender</label>
                    <div className="relative">
                        <User className="absolute left-2.5 top-2.5 w-4 h-4 text-slate-400" />
                        <input
                            type="text"
                            value={localFilters.sender}
                            onChange={(e) => handleChange('sender', e.target.value)}
                            className="w-full pl-9 pr-3 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                            placeholder="e.g. alice@example.com"
                        />
                    </div>
                </div>

                {/* Date Range Filter */}
                <div>
                    <label className="block text-xs font-medium text-slate-500 mb-1.5">Date Range</label>
                    <div className="space-y-2">
                        <div className="relative">
                            <Calendar className="absolute left-2.5 top-2.5 w-4 h-4 text-slate-400" />
                            <input
                                type="date"
                                value={localFilters.startDate}
                                onChange={(e) => handleChange('startDate', e.target.value)}
                                className="w-full pl-9 pr-3 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-slate-600"
                            />
                        </div>
                        <div className="relative">
                            <Calendar className="absolute left-2.5 top-2.5 w-4 h-4 text-slate-400" />
                            <input
                                type="date"
                                value={localFilters.endDate}
                                onChange={(e) => handleChange('endDate', e.target.value)}
                                className="w-full pl-9 pr-3 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-slate-600"
                            />
                        </div>
                    </div>
                </div>

                {/* Active Filter Indicator */}
                {(localFilters.sender || localFilters.startDate || localFilters.endDate) && (
                    <div className="pt-2 border-t border-slate-50">
                        <button
                            onClick={() => {
                                const reset = { sender: '', startDate: '', endDate: '' };
                                setLocalFilters(reset);
                                onChange(reset);
                            }}
                            className="text-xs text-blue-600 hover:text-blue-700 font-medium w-full text-center hover:underline"
                        >
                            Clear All Filters
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
