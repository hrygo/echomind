import { Calendar, Clock, MapPin } from 'lucide-react';

export interface CalendarEventWidgetProps {
    data: {
        title: string;
        start: string;
        end: string;
        location?: string;
        description?: string;
    };
}

export function CalendarEventWidget({ data }: CalendarEventWidgetProps) {
    return (
        <div className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm my-2 flex gap-4">
            <div className="flex-shrink-0 flex flex-col items-center justify-center w-14 h-14 bg-blue-50 rounded-lg border border-blue-100 text-blue-600">
                <span className="text-[10px] font-bold uppercase tracking-wider">
                    {new Date(data.start).toLocaleString('en-US', { month: 'short' })}
                </span>
                <span className="text-xl font-bold leading-none mt-0.5">
                    {new Date(data.start).getDate()}
                </span>
            </div>

            <div className="flex-1 min-w-0">
                <h4 className="font-semibold text-slate-800 text-sm truncate">{data.title}</h4>

                <div className="mt-2 space-y-1.5">
                    <div className="flex items-center gap-2 text-xs text-slate-600">
                        <Clock className="w-3.5 h-3.5 text-slate-400" />
                        <span>
                            {new Date(data.start).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })} -
                            {new Date(data.end).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </span>
                    </div>

                    {data.location && (
                        <div className="flex items-center gap-2 text-xs text-slate-600">
                            <MapPin className="w-3.5 h-3.5 text-slate-400" />
                            <span className="truncate">{data.location}</span>
                        </div>
                    )}
                </div>
            </div>

            <div className="flex flex-col justify-center">
                <button className="p-2 hover:bg-slate-100 rounded-lg text-blue-600 transition-colors" title="Add to Calendar">
                    <Calendar className="w-5 h-5" />
                </button>
            </div>
        </div>
    );
}
