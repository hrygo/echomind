'use client';

import { create } from 'zustand';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { AlertTriangle } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface ConfirmDialogState {
  isOpen: boolean;
  title: string;
  message: string;
  confirmText: string;
  cancelText: string;
  onConfirm: () => void;
  onCancel?: () => void;
}

interface ConfirmDialogStore extends ConfirmDialogState {
  openConfirm: (options: Omit<ConfirmDialogState, 'isOpen'>) => void;
  closeConfirm: () => void;
}

const useConfirmDialogStore = create<ConfirmDialogStore>((set) => ({
  isOpen: false,
  title: '',
  message: '',
  confirmText: '',
  cancelText: '',
  onConfirm: () => {},
  onCancel: undefined,
  
  openConfirm: (options) => set({ ...options, isOpen: true }),
  closeConfirm: () => set({ isOpen: false }),
}));

export function useConfirm() {
  const { openConfirm, closeConfirm } = useConfirmDialogStore();
  const { t } = useLanguage();

  return (message: string, onConfirm: () => void, options?: {
    title?: string;
    confirmText?: string;
    cancelText?: string;
    onCancel?: () => void;
  }) => {
    openConfirm({
      title: options?.title || t('common.confirmAction'),
      message,
      confirmText: options?.confirmText || t('common.confirm'),
      cancelText: options?.cancelText || t('common.cancel'),
      onConfirm: () => {
        onConfirm();
        closeConfirm();
      },
      onCancel: options?.onCancel,
    });
  };
}

export function ConfirmDialog() {
  const { isOpen, title, message, confirmText, cancelText, onConfirm, onCancel, closeConfirm } = useConfirmDialogStore();

  const handleCancel = () => {
    if (onCancel) onCancel();
    closeConfirm();
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && closeConfirm()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center flex-shrink-0">
              <AlertTriangle className="w-5 h-5 text-orange-600" />
            </div>
            <DialogTitle className="text-lg">{title}</DialogTitle>
          </div>
          <DialogDescription className="pt-2 text-base">
            {message}
          </DialogDescription>
        </DialogHeader>
        
        <div className="flex gap-3 justify-end mt-4">
          <Button
            onClick={handleCancel}
            variant="outline"
            className="px-4 py-2"
          >
            {cancelText}
          </Button>
          <Button
            onClick={onConfirm}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white"
          >
            {confirmText}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
