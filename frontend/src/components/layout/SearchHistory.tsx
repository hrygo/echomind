import { Clock, X, Trash2 } from "lucide-react";

interface SearchHistoryProps {
  history: string[];
  onSelect: (query: string) => void;
  onRemove: (query: string) => void;
  onClear: () => void;
}

export function SearchHistory({ history, onSelect, onRemove, onClear }: SearchHistoryProps) {
  if (history.length === 0) return null;

  return (
    <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in fade-in slide-in-from-top-2 z-50">
        <div className="flex items-center justify-between px-4 py-2 bg-slate-50 border-b border-slate-100">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Recent Searches</span>
            <button 
                onClick={(e) => {
                    e.stopPropagation();
                    onClear();
                }}
                className="text-xs text-red-500 hover:text-red-700 font-medium flex items-center gap-1 px-2 py-1 rounded hover:bg-red-50 transition-colors"
            >
                <Trash2 className="w-3 h-3" />
                Clear
            </button>
        </div>
        <div className="max-h-[300px] overflow-y-auto py-1">
            {history.map((query, index) => (
                <div 
                    key={index}
                    className="group flex items-center justify-between px-4 py-2 hover:bg-slate-50 cursor-pointer transition-colors"
                    onClick={() => onSelect(query)}
                >
                    <div className="flex items-center gap-3 text-slate-700 overflow-hidden">
                        <Clock className="w-4 h-4 text-slate-400 flex-shrink-0" />
                        <span className="text-sm truncate">{query}</span>
                    </div>
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            onRemove(query);
                        }}
                        className="p-1 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-full opacity-0 group-hover:opacity-100 transition-all"
                        title="Remove from history"
                    >
                        <X className="w-3 h-3" />
                    </button>
                </div>
            ))}
        </div>
    </div>
  );
}
