import React, { useState } from 'react';
import { CheckCircle2, Clock, ArrowRight, Calendar, AlertCircle } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";

// Mock Data Types
interface Task {
  id: string;
  title: string;
  sender: string;
  due: string;
  priority: 'High' | 'Medium' | 'Low';
  status: 'pending' | 'completed';
  isOverdue?: boolean;
}

const INITIAL_TASKS: Task[] = [
  { id: '1', title: "Submit Monthly Operations Report", sender: "boss@company.com", due: "Yesterday", priority: "High", status: "pending", isOverdue: true },
  { id: '2', title: "Review Q4 Marketing Budget", sender: "finance@company.com", due: "Today", priority: "High", status: "pending" },
  { id: '3', title: "Team Performance Review - Alice", sender: "hr@company.com", due: "Tomorrow", priority: "Medium", status: "pending" },
  { id: '4', title: "Approve Vacation Request - Bob", sender: "bob@company.com", due: "Fri", priority: "Low", status: "pending" },
];

export function ManagerView() {
  const { t } = useLanguage();
  const [tasks, setTasks] = useState<Task[]>(INITIAL_TASKS);
  const [filter, setFilter] = useState<'all' | 'high'>('all');

  const handleToggleTask = (id: string) => {
    setTasks(prev => prev.map(task => 
      task.id === id 
        ? { ...task, status: task.status === 'pending' ? 'completed' : 'pending' }
        : task
    ));
    
    // In a real app, we would call an API here
    // setTimeout(() => { remove task or move to bottom }, 500)
  };

  const activeTasks = tasks.filter(t => t.status === 'pending');
  const visibleTasks = filter === 'high' ? activeTasks.filter(t => t.priority === 'High') : activeTasks;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Main Column: Action Items */}
      <div className="lg:col-span-2 space-y-6">
        <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold text-slate-800 flex items-center gap-2">
            <CheckCircle2 className="w-5 h-5 text-blue-600" />
            {t('dashboard.actionItems')}
            <span className="text-xs font-normal text-slate-400 bg-slate-100 px-2 py-0.5 rounded-full ml-2">
                {activeTasks.length} Pending
            </span>
            </h2>
            
            <div className="flex gap-2">
                <button 
                    onClick={() => setFilter('all')}
                    className={`text-xs font-medium px-3 py-1.5 rounded-lg transition-colors ${filter === 'all' ? 'bg-slate-800 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}`}
                >
                    All
                </button>
                <button 
                    onClick={() => setFilter('high')}
                    className={`text-xs font-medium px-3 py-1.5 rounded-lg transition-colors ${filter === 'high' ? 'bg-red-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}`}
                >
                    High Priority
                </button>
            </div>
        </div>

        <div className="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
          {visibleTasks.length === 0 ? (
             <div className="p-8 text-center text-slate-400">
                <CheckCircle2 className="w-12 h-12 mx-auto mb-3 opacity-20" />
                <p>No pending tasks! Great job.</p>
             </div>
          ) : (
             <div className="divide-y divide-slate-50">
                {visibleTasks.map((task) => (
                <TaskItem key={task.id} task={task} onToggle={() => handleToggleTask(task.id)} />
                ))}
             </div>
          )}
          
          <div className="p-3 bg-slate-50/50 text-center border-t border-slate-100">
            <Link href="/dashboard/tasks" className="text-sm font-medium text-blue-600 hover:text-blue-700 flex items-center justify-center gap-1 transition-colors">
              {t('dashboard.viewAll')} <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </div>
      </div>

      {/* Right Column: Follow-ups & Stats */}
      <div className="space-y-6">
        {/* Stats Widget */}
        <div className="grid grid-cols-2 gap-3">
             <div className="bg-blue-50/50 p-4 rounded-xl border border-blue-100">
                 <p className="text-xs font-medium text-blue-600 uppercase tracking-wider mb-1">Completed</p>
                 <p className="text-2xl font-bold text-slate-800">14</p>
                 <p className="text-[10px] text-slate-400">This week</p>
             </div>
             <div className="bg-orange-50/50 p-4 rounded-xl border border-orange-100">
                 <p className="text-xs font-medium text-orange-600 uppercase tracking-wider mb-1">Overdue</p>
                 <p className="text-2xl font-bold text-slate-800">{activeTasks.filter(t => t.isOverdue).length}</p>
                 <p className="text-[10px] text-slate-400">Action needed</p>
             </div>
        </div>

        {/* Smart Follow-up */}
        <section>
            <h2 className="text-lg font-bold text-slate-800 mb-3 flex items-center gap-2">
            <Clock className="w-4 h-4 text-orange-500" />
            {t('dashboard.smartFollowUp')}
            </h2>
            <div className="space-y-3">
            {[
                { name: "Alice Smith", subject: "Re: Project Proposal", time: "2d", waiting: true },
                { name: "Bob Jones", subject: "Contract Draft Review", time: "3d", waiting: true },
                { name: "Charlie Day", subject: "Lunch Meeting", time: "5h", waiting: false }, // Not waiting, just recent
            ].filter(i => i.waiting).map((item, i) => (
                <div key={i} className="bg-white p-3 rounded-xl border border-slate-100 shadow-sm flex items-center gap-3 hover:shadow-md transition-shadow cursor-pointer">
                    <div className="w-8 h-8 rounded-full bg-orange-100 flex items-center justify-center text-orange-600 font-bold text-xs">
                        {item.name[0]}
                    </div>
                    <div className="flex-1 min-w-0">
                        <h4 className="font-semibold text-sm text-slate-800 truncate">{item.name}</h4>
                        <p className="text-xs text-slate-500 truncate">{item.subject}</p>
                    </div>
                    <div className="text-right whitespace-nowrap">
                        <span className="text-[10px] font-bold text-orange-500 bg-orange-50 px-1.5 py-0.5 rounded">Waiting</span>
                        <p className="text-[10px] text-slate-400 mt-0.5">{item.time}</p>
                    </div>
                </div>
            ))}
            </div>
        </section>
      </div>
    </div>
  );
}

function TaskItem({ task, onToggle }: { task: Task; onToggle: () => void }) {
    return (
        <div 
            onClick={onToggle}
            className="group p-4 flex items-start gap-4 hover:bg-slate-50 transition-all duration-200 cursor-pointer animate-in fade-in slide-in-from-bottom-1"
        >
            {/* Checkbox */}
            <div className={`mt-0.5 w-5 h-5 rounded-full border-2 flex items-center justify-center transition-colors
                ${task.status === 'completed' ? 'bg-blue-500 border-blue-500' : 'border-slate-300 group-hover:border-blue-500'}
            `}>
                {task.status === 'completed' && <CheckCircle2 className="w-3.5 h-3.5 text-white" />}
            </div>

            {/* Content */}
            <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-2">
                    <h3 className={`font-medium text-sm transition-all ${task.status === 'completed' ? 'text-slate-400 line-through' : 'text-slate-800'}`}>
                        {task.title}
                    </h3>
                    {task.isOverdue && task.status !== 'completed' && (
                        <div className="flex items-center gap-1 text-[10px] font-bold text-red-600 bg-red-50 px-2 py-0.5 rounded-full whitespace-nowrap">
                            <AlertCircle className="w-3 h-3" /> Overdue
                        </div>
                    )}
                </div>
                
                <p className="text-xs text-slate-500 mt-1 truncate">
                    From: <span className="text-slate-600 font-medium">{task.sender}</span>
                </p>

                {/* Footer Tags */}
                <div className="flex items-center gap-2 mt-2.5">
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded border
                        ${task.priority === 'High' ? 'bg-red-50 text-red-700 border-red-100' :
                          task.priority === 'Medium' ? 'bg-orange-50 text-orange-700 border-orange-100' :
                          'bg-slate-50 text-slate-600 border-slate-200'}
                    `}>
                        {task.priority}
                    </span>
                    <span className="flex items-center gap-1 text-[10px] text-slate-500 font-medium bg-slate-50 px-2 py-0.5 rounded border border-slate-100">
                        <Calendar className="w-3 h-3" /> {task.due}
                    </span>
                </div>
            </div>
        </div>
    );
}
import Link from 'next/link';
