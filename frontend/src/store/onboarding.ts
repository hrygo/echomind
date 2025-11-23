import { create } from 'zustand';
import { MailProviderConfig } from '@/lib/constants/mail_providers';

interface OnboardingState {
  step: 1 | 2 | 3 | 4; // Step 4 is optional WeChat Bind, for future
  role: 'executive' | 'manager' | 'dealmaker' | null; // User's selected role
  mailbox: {
    email: string;
    password: string;
    imapServer: string;
    imapPort: number;
    smtpServer: string;
    smtpPort: number;
    providerConfig: MailProviderConfig | null; // Auto-detected config
    // manualConfig?: { imap: { host: string; port: number; secure: boolean }; smtp: { host: string; port: number; secure: boolean } }; // User manually entered config
  };
  
  // Actions
  setStep: (step: OnboardingState['step']) => void;
  setRole: (role: OnboardingState['role']) => void;
  setMailboxConfig: (data: Partial<OnboardingState['mailbox']>) => void;
  resetOnboarding: () => void; // Resets the entire onboarding state
}

export const useOnboardingStore = create<OnboardingState>((set) => ({
  step: 1,
  role: null,
  mailbox: {
    email: '',
    password: '',
    imapServer: '',
    imapPort: 0,
    smtpServer: '',
    smtpPort: 0,
    providerConfig: null,
  },

  setStep: (step) => set({ step }),
  setRole: (role) => set({ role }),
  setMailboxConfig: (data) => set(state => ({ mailbox: { ...state.mailbox, ...data } })),
  resetOnboarding: () => set({
    step: 1,
    role: null,
    mailbox: {
      email: '',
      password: '',
      imapServer: '',
      imapPort: 0,
      smtpServer: '',
      smtpPort: 0,
      providerConfig: null,
    },
  }),
}));
