import { CheckCircle, Calendar } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface TaskWidgetProps {
    data: {
        title: string;
        dueDate?: string;
        priority?: 'High' | 'Medium' | 'Low';
        status?: string;
    };
}

export function TaskWidget({ data }: TaskWidgetProps) {
    return (
        <div className="bg-white border border-gray-200 rounded-xl p-4 shadow-sm my-2">
            <div className="flex items-start gap-3">
                <div className="mt-1 bg-green-100 p-1.5 rounded-full text-green-600">
                    <CheckCircle className="w-4 h-4" />
                </div>
                <div className="flex-1">
                    <h4 className="font-medium text-gray-900 text-sm">{data.title}</h4>
                    {data.dueDate && (
                        <div className="flex items-center gap-1 mt-2 text-xs text-gray-500">
                            <Calendar className="w-3 h-3" />
                            <span>{data.dueDate}</span>
                        </div>
                    )}
                </div>
                {data.priority && (
                    <span className={cn(
                        "text-xs px-2 py-0.5 rounded-full font-medium",
                        data.priority === 'High' ? "bg-red-100 text-red-700" :
                        data.priority === 'Medium' ? "bg-yellow-100 text-yellow-700" :
                        "bg-blue-100 text-blue-700"
                    )}>
                        {data.priority}
                    </span>
                )}
            </div>
        </div>
    );
}
