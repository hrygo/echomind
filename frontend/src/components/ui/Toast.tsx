"use client";

import React, { useEffect } from 'react';
import { X } from 'lucide-react';
import { useActionStore } from '@/lib/store/actions';

interface ToastProps {
  id: string;
  message: string;
  action?: () => void; // Optional undo action
  onClose: (id: string) => void;
}

export function Toast({ id, message, action, onClose }: ToastProps) {
  return (
    <div className="flex items-center justify-between gap-4 p-4 pr-5 text-sm font-medium text-white bg-slate-800 rounded-lg shadow-lg animate-in fade-in slide-in-from-bottom-8 duration-300">
      <span>{message}</span>
      <div className="flex items-center gap-2">
        {action && (
          <button 
            onClick={action}
            className="px-3 py-1 -my-1 rounded bg-blue-600 hover:bg-blue-700 text-white transition-colors"
          >
            Undo
          </button>
        )}
        <button onClick={() => onClose(id)} className="text-slate-400 hover:text-slate-200 transition-colors">
          <X className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}

export function ToastContainer() {
  const { toasts, removeToast } = useActionStore();

  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col gap-3 pointer-events-none">
      {toasts.map(toast => (
        <Toast 
          key={toast.id} 
          id={toast.id} 
          message={toast.message} 
          action={() => {
            toast.action();
            removeToast(toast.id);
          }}
          onClose={removeToast}
        />
      ))}
    </div>
  );
}
