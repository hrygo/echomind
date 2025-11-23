import { create } from 'zustand';
import { ActionAPI } from '../api/actions';
import { Email } from '../api/emails'; // To temporarily store emails for undo
import { useEmailStore } from './emails'; // To trigger a re-fetch or direct update

interface UndoToast {
  id: string;
  message: string;
  action: () => void; // Function to revert the action
  timeoutId: NodeJS.Timeout; // For auto-dismissal
}

interface ActionState {
  // State for undo toasts
  toasts: UndoToast[];
  addToast: (message: string, onUndo: () => void, timeout?: number) => void;
  removeToast: (id: string) => void;

  // Email actions
  approveEmail: (emailId: string) => Promise<void>;
  snoozeEmail: (emailId: string, duration: string) => Promise<void>;
  dismissEmail: (emailId: string) => Promise<void>;
}

export const useActionStore = create<ActionState>((set, get) => ({
  toasts: [],

  addToast: (message, onUndo, timeout = 5000) => {
    const id = Date.now().toString();
    const timeoutId = setTimeout(() => {
      get().removeToast(id);
    }, timeout);

    set(state => ({
      toasts: [...state.toasts, { id, message, action: onUndo, timeoutId }]
    }));
  },

  removeToast: (id: string) => {
    set(state => ({
      toasts: state.toasts.filter(toast => {
        if (toast.id === id) {
          clearTimeout(toast.timeoutId);
          return false;
        }
        return true;
      })
    }));
  },

  approveEmail: async (emailId: string) => {
    // Optimistic UI update
    const { emails, fetchEmails } = useEmailStore.getState();
    const originalEmails = [...emails];
    
    useEmailStore.setState({
      emails: emails.filter(e => e.ID !== emailId)
    });

    try {
      await ActionAPI.approve(emailId);
      get().addToast(
        'Email approved. Undo?',
        async () => {
          // Revert: add email back to store (simple re-fetch for now)
          // In a more complex setup, you might have an ActionAPI.unapprove()
          await fetchEmails();
          get().removeToast(emailId); // Remove specific toast
        }
      );
    } catch (error) {
      console.error("Failed to approve email:", error);
      // Revert optimistic update on error
      useEmailStore.setState({ emails: originalEmails });
      get().addToast('Failed to approve email. Retry?', () => get().approveEmail(emailId));
      throw error;
    }
  },

  snoozeEmail: async (emailId: string, duration: string) => {
    const { emails, fetchEmails } = useEmailStore.getState();
    const originalEmails = [...emails];

    // Optimistic UI: remove from current view
    useEmailStore.setState({
      emails: emails.filter(e => e.ID !== emailId)
    });

    try {
      const result = await ActionAPI.snooze(emailId, duration);
      get().addToast(
        `Email snoozed until ${new Date(result.until).toLocaleString()}. Undo?`,
        async () => {
          // Revert: re-fetch or directly update the snoozed_until to null
          // For simplicity, re-fetch emails
          await fetchEmails();
          get().removeToast(emailId);
        }
      );
    } catch (error) {
      console.error("Failed to snooze email:", error);
      // Revert optimistic update on error
      useEmailStore.setState({ emails: originalEmails });
      get().addToast('Failed to snooze email. Retry?', () => get().snoozeEmail(emailId, duration));
      throw error;
    }
  },

  dismissEmail: async (emailId: string) => {
    const { emails, fetchEmails } = useEmailStore.getState();
    const originalEmails = [...emails];

    // Optimistic UI: update urgency to Low locally
    useEmailStore.setState({
      emails: emails.map(e => e.ID === emailId ? { ...e, Urgency: 'Low' } : e)
    });

    try {
      await ActionAPI.dismiss(emailId);
      get().addToast(
        'Email dismissed. Undo?',
        async () => {
          // Revert: re-fetch or directly update urgency to original
          await fetchEmails();
          get().removeToast(emailId);
        }
      );
    } catch (error) {
      console.error("Failed to dismiss email:", error);
      // Revert optimistic update on error
      useEmailStore.setState({ emails: originalEmails });
      get().addToast('Failed to dismiss email. Retry?', () => get().dismissEmail(emailId));
      throw error;
    }
  },
}));
