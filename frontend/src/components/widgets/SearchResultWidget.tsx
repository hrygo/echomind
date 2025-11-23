import { User, Calendar } from 'lucide-react';

interface SearchResultWidgetProps {
    data: {
        results: Array<{
            subject: string;
            sender: string;
            date: string;
            snippet: string;
        }>;
        count: number;
    };
}

export function SearchResultWidget({ data }: SearchResultWidgetProps) {
    if (!data.results || data.results.length === 0) return null;

    return (
        <div className="bg-white border border-gray-200 rounded-xl overflow-hidden shadow-sm my-2">
            <div className="bg-gray-50 px-4 py-2 border-b border-gray-100 flex justify-between items-center">
                <span className="text-xs font-medium text-gray-600">Found {data.count} emails</span>
            </div>
            <div className="divide-y divide-gray-100">
                {data.results.slice(0, 3).map((email, idx) => (
                    <div key={idx} className="p-3 hover:bg-gray-50 transition-colors cursor-pointer">
                        <h5 className="text-sm font-medium text-gray-900 truncate">{email.subject}</h5>
                        <p className="text-xs text-gray-500 mt-0.5 line-clamp-2">{email.snippet}</p>
                        <div className="flex items-center gap-3 mt-2 text-[10px] text-gray-400">
                            <span className="flex items-center gap-1">
                                <User className="w-3 h-3" /> {email.sender}
                            </span>
                            <span className="flex items-center gap-1">
                                <Calendar className="w-3 h-3" /> {new Date(email.date).toLocaleDateString()}
                            </span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
