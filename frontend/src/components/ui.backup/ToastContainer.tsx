'use client';

import { useToastStore, getToastIcon, getToastStyles } from '@/lib/hooks/useToast';
import { X } from 'lucide-react';
import { useEffect, useState } from 'react';

export function ToastContainer() {
  const { toasts, removeToast } = useToastStore();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setMounted(true);
  }, []);

  if (!mounted) return null;

  return (
    <div className="fixed top-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className={`
            pointer-events-auto
            flex items-start gap-3 p-4 rounded-lg border shadow-lg
            min-w-[320px] max-w-md
            animate-in slide-in-from-right duration-300
            ${getToastStyles(toast.type)}
          `}
        >
          <div className="mt-0.5">
            {getToastIcon(toast.type)}
          </div>

          <div className="flex-1 min-w-0">
            {toast.title && (
              <div className="font-semibold mb-1 text-sm">
                {toast.title}
              </div>
            )}
            <div className="text-sm break-words">
              {toast.message}
            </div>
          </div>

          <button
            onClick={() => removeToast(toast.id)}
            className="flex-shrink-0 p-1 rounded-md hover:bg-black/5 transition-colors"
            aria-label="关闭"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
      ))}
    </div>
  );
}
