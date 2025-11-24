import { Mail, Edit2, Send } from 'lucide-react';

export interface EmailDraftWidgetProps {
    data: {
        to: string;
        subject: string;
        body: string;
    };
}

export function EmailDraftWidget({ data }: EmailDraftWidgetProps) {
    return (
        <div className="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm my-2">
            <div className="bg-slate-50 px-4 py-2 border-b border-slate-200 flex items-center justify-between">
                <div className="flex items-center gap-2 text-xs font-medium text-slate-600">
                    <Mail className="w-3.5 h-3.5" />
                    <span>Draft Email</span>
                </div>
                <div className="flex gap-1">
                    <button className="p-1.5 hover:bg-white rounded-md text-slate-500 hover:text-blue-600 transition-colors" title="Edit">
                        <Edit2 className="w-3.5 h-3.5" />
                    </button>
                </div>
            </div>
            <div className="p-4 space-y-3">
                <div className="space-y-1">
                    <div className="text-xs text-slate-500 font-medium">To:</div>
                    <div className="text-sm text-slate-800 bg-slate-50 px-2 py-1 rounded border border-slate-100">{data.to}</div>
                </div>
                <div className="space-y-1">
                    <div className="text-xs text-slate-500 font-medium">Subject:</div>
                    <div className="text-sm font-medium text-slate-800">{data.subject}</div>
                </div>
                <div className="pt-2 border-t border-slate-100">
                    <div className="text-sm text-slate-600 whitespace-pre-wrap leading-relaxed font-mono text-[13px] bg-slate-50 p-3 rounded-lg border border-slate-100">
                        {data.body}
                    </div>
                </div>
            </div>
            <div className="bg-slate-50 px-4 py-3 border-t border-slate-200 flex justify-end gap-2">
                <button className="px-3 py-1.5 text-xs font-medium text-slate-600 hover:text-slate-800 hover:bg-slate-200 rounded-lg transition-colors">
                    Discard
                </button>
                <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium rounded-lg shadow-sm flex items-center gap-1.5 transition-colors">
                    <Send className="w-3 h-3" />
                    Send Now
                </button>
            </div>
        </div>
    );
}
