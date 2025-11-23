"use client";

import { useEffect } from "react";
import { Checkbox } from "@/components/ui/Checkbox"; // Assuming you have a Checkbox component
import { useTaskStore } from "@/store/task";
import { listTasks, updateTaskStatus, Task } from "@/lib/api/tasks";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { cn } from "@/lib/utils";
import { Loader2, Calendar } from "lucide-react";

interface TaskWidgetProps {
  // Optional: filters for the tasks to display in this widget
  initialStatus?: 'todo' | 'in_progress' | 'done';
  initialPriority?: 'high' | 'medium' | 'low';
}

export function TaskWidget({ initialStatus = 'todo', initialPriority }: TaskWidgetProps) {
  const { tasks, isLoading, error, setTasks, updateTask, setError, setLoading } = useTaskStore();
  const { t } = useLanguage();

  useEffect(() => {
    const fetchTasks = async () => {
      setLoading(true);
      try {
        const fetchedTasks = await listTasks(initialStatus, initialPriority);
        // Sort tasks: pending/in_progress first, then due date, then creation date
        const sortedTasks = fetchedTasks.sort((a, b) => {
          // Prioritize non-done tasks
          if (a.status !== 'done' && b.status === 'done') return -1;
          if (a.status === 'done' && b.status !== 'done') return 1;

          // Then by due date
          const dateA = a.due_date ? new Date(a.due_date).getTime() : Infinity;
          const dateB = b.due_date ? new Date(b.due_date).getTime() : Infinity;
          if (dateA !== dateB) return dateA - dateB;

          // Finally by creation date (newest first for same due date)
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        });
        setTasks(sortedTasks);
      } catch (err: unknown) {
        console.error("Failed to fetch tasks:", err);
        setError(t('common.error')); // Use i18n
      } finally {
        setLoading(false);
      }
    };
    fetchTasks();
  }, [initialStatus, initialPriority, setTasks, setError, setLoading, t]);

  const handleToggleTaskStatus = async (task: Task) => {
    const newStatus = task.status === 'done' ? 'todo' : 'done';
    // Optimistic UI update
    updateTask({ ...task, status: newStatus });

    try {
      await updateTaskStatus(task.id, newStatus);
      // If successful, UI is already updated
    } catch (err: unknown) {
      // If API fails, revert UI and show error
      updateTask(task); // Revert to original status
      console.error("Failed to update task status:", err);
      setError(t('common.error')); // Use i18n
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-6 text-slate-500">
        <Loader2 className="w-5 h-5 animate-spin mr-2" /> {t('common.loading')}
      </div>
    );
  }

  if (error) {
    return <div className="p-6 text-red-500">{t('common.error')}: {error}</div>;
  }

  if (tasks.length === 0) {
    return (
      <div className="p-6 text-slate-500 text-center">
        {t('dashboard.noPendingTasks')}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {tasks.map((task) => (
        <div
          key={task.id}
          className={cn(
            "flex items-center p-3 bg-white rounded-lg border border-slate-100 shadow-sm",
            task.status === 'done' && "opacity-70 line-through text-slate-500"
          )}
        >
          <Checkbox
            id={task.id}
            checked={task.status === 'done'}
            onCheckedChange={() => handleToggleTaskStatus(task)}
            className="mr-3"
          />
          <label htmlFor={task.id} className="flex-1 font-medium text-slate-800 cursor-pointer">
            {task.title}
            {task.due_date && (
              <span className={cn(
                "ml-2 text-xs flex items-center gap-1 text-slate-500",
                new Date(task.due_date) < new Date() && task.status !== 'done' && "text-red-500 font-semibold"
              )}>
                <Calendar className="w-3 h-3" />
                {new Date(task.due_date).toLocaleDateString()}
              </span>
            )}
          </label>
          {/* Priority or other actions can go here */}
        </div>
      ))}
    </div>
  );
}
