import { api } from '@/lib/api';

export interface Task {
    id: string;
    title: string;
    description?: string;
    status: 'todo' | 'in_progress' | 'done';
    priority: 'high' | 'medium' | 'low';
    due_date?: string;
    source_email_id?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateTaskPayload {
    title: string;
    description?: string;
    source_email_id?: string;
    due_date?: string;
}

export interface UpdateTaskPayload {
    title?: string;
    description?: string;
    priority?: 'high' | 'medium' | 'low';
    due_date?: string;
}

export interface UpdateTaskStatusPayload {
    status: 'todo' | 'in_progress' | 'done';
}

export const createTask = async (payload: CreateTaskPayload): Promise<Task> => {
    const response = await api.post<Task>('/tasks', payload);
    return response.data;
};

export const listTasks = async (status?: string, priority?: string): Promise<Task[]> => {
    const params: Record<string, string> = {};
    if (status) params.status = status;
    if (priority) params.priority = priority;
    const response = await api.get<Task[]>('/tasks', { params });
    return response.data;
};

export const updateTask = async (id: string, payload: UpdateTaskPayload): Promise<Task> => {
    const response = await api.patch<Task>(`/tasks/${id}`, payload);
    return response.data;
};

export const updateTaskStatus = async (id: string, status: 'todo' | 'in_progress' | 'done'): Promise<void> => {
    await api.patch(`/tasks/${id}/status`, { status });
};

export const deleteTask = async (id: string): Promise<void> => {
    await api.delete(`/tasks/${id}`);
};
