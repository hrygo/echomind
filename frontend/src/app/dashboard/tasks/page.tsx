import { PageHeader } from "@/components/ui/page-header";

export default function TasksPage() {
    return (
        <div>
            <PageHeader title="任务列表" />
            <div className="rounded-xl border border-slate-200 bg-white shadow-sm p-8 text-center">
                <p className="text-slate-500 text-lg">任务管理功能即将上线。</p>
            </div>
        </div>
    );
}
