import { CheckCircle, Calendar, ArrowRight } from 'lucide-react';
import { cn } from '@/lib/utils';

interface TaskItem {
    title: string;
    dueDate?: string;
    priority?: 'High' | 'Medium' | 'Low';
    status?: string;
}

export interface TaskListWidgetProps {
    data: TaskItem[];
}

export function TaskListWidget({ data }: TaskListWidgetProps) {
    if (!Array.isArray(data)) {
        return <div className="text-red-500 text-xs">Invalid task data</div>;
    }

    return (
        <div className="space-y-2">
            {data.map((task, idx) => (
                <div key={idx} className="bg-white border border-slate-200 rounded-lg p-3 shadow-sm hover:border-blue-300 transition-colors group">
                    <div className="flex items-start gap-3">
                        <div className="mt-0.5 text-slate-400 group-hover:text-blue-500 transition-colors">
                            <CheckCircle className="w-4 h-4" />
                        </div>
                        <div className="flex-1 min-w-0">
                            <h4 className="font-medium text-slate-800 text-sm truncate">{task.title}</h4>
                            <div className="flex items-center gap-3 mt-1.5">
                                {task.dueDate && (
                                    <div className="flex items-center gap-1 text-xs text-slate-500">
                                        <Calendar className="w-3 h-3" />
                                        <span>{task.dueDate}</span>
                                    </div>
                                )}
                                {task.priority && (
                                    <span className={cn(
                                        "text-[10px] px-1.5 py-0.5 rounded-full font-medium uppercase tracking-wide",
                                        task.priority === 'High' ? "bg-red-50 text-red-600 border border-red-100" :
                                            task.priority === 'Medium' ? "bg-amber-50 text-amber-600 border border-amber-100" :
                                                "bg-blue-50 text-blue-600 border border-blue-100"
                                    )}>
                                        {task.priority}
                                    </span>
                                )}
                            </div>
                        </div>
                        <button className="opacity-0 group-hover:opacity-100 p-1.5 hover:bg-slate-100 rounded-md transition-all text-slate-400 hover:text-blue-600">
                            <ArrowRight className="w-4 h-4" />
                        </button>
                    </div>
                </div>
            ))}
            <div className="flex justify-end">
                <button className="text-xs font-medium text-blue-600 hover:text-blue-700 hover:underline px-2 py-1">
                    Add All to Tasks
                </button>
            </div>
        </div>
    );
}
