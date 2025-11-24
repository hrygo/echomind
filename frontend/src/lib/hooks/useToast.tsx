import { create } from 'zustand';
import { CheckCircle2, AlertCircle, Info, XCircle } from 'lucide-react';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastMessage {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
  duration?: number;
}

interface ToastStore {
  toasts: ToastMessage[];
  addToast: (toast: Omit<ToastMessage, 'id'>) => void;
  removeToast: (id: string) => void;
}

export const useToastStore = create<ToastStore>((set) => ({
  toasts: [],
  
  addToast: (toast) => {
    const id = `toast-${Date.now()}-${Math.random()}`;
    const duration = toast.duration || 4000;
    
    set((state) => ({
      toasts: [...state.toasts, { ...toast, id }]
    }));

    // Auto remove after duration
    if (duration > 0) {
      setTimeout(() => {
        set((state) => ({
          toasts: state.toasts.filter((t) => t.id !== id)
        }));
      }, duration);
    }
  },

  removeToast: (id) => {
    set((state) => ({
      toasts: state.toasts.filter((t) => t.id !== id)
    }));
  },
}));

export const useToast = () => {
  const { addToast } = useToastStore();

  return {
    success: (message: string, title?: string, duration?: number) => {
      addToast({ type: 'success', message, title, duration });
    },
    error: (message: string, title?: string, duration?: number) => {
      addToast({ type: 'error', message, title, duration });
    },
    info: (message: string, title?: string, duration?: number) => {
      addToast({ type: 'info', message, title, duration });
    },
    warning: (message: string, title?: string, duration?: number) => {
      addToast({ type: 'warning', message, title, duration });
    },
  };
};

export const getToastIcon = (type: ToastType) => {
  const iconProps = { className: 'w-5 h-5 flex-shrink-0' };
  
  switch (type) {
    case 'success':
      return <CheckCircle2 {...iconProps} />;
    case 'error':
      return <XCircle {...iconProps} />;
    case 'warning':
      return <AlertCircle {...iconProps} />;
    case 'info':
      return <Info {...iconProps} />;
  }
};

export const getToastStyles = (type: ToastType) => {
  switch (type) {
    case 'success':
      return 'bg-green-50 text-green-800 border-green-200';
    case 'error':
      return 'bg-red-50 text-red-800 border-red-200';
    case 'warning':
      return 'bg-orange-50 text-orange-800 border-orange-200';
    case 'info':
      return 'bg-blue-50 text-blue-800 border-blue-200';
  }
};
