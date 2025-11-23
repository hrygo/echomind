export interface MailProviderConfig {
  id: string;
  domains: string[]; // Matching domain suffixes, e.g., ["gmail.com", "googlemail.com"]
  name: string;      // Display name
  imap: { host: string; port: number; secure: boolean };
  smtp: { host: string; port: number; secure: boolean };
  helpLink?: string; // Link to help article for App Passwords, etc.
  requiresAppPassword?: boolean; // Whether to show a hint about App Passwords
}

export const MAIL_PROVIDERS: MailProviderConfig[] = [
  {
    id: 'gmail',
    domains: ['gmail.com', 'googlemail.com'],
    name: 'Gmail',
    imap: { host: 'imap.gmail.com', port: 993, secure: true },
    smtp: { host: 'smtp.gmail.com', port: 465, secure: true },
    helpLink: 'https://support.google.com/accounts/answer/185833',
    requiresAppPassword: true
  },
  {
    id: 'outlook',
    domains: ['outlook.com', 'hotmail.com', 'live.com', 'msn.com'],
    name: 'Outlook',
    imap: { host: 'outlook.office365.com', port: 993, secure: true },
    smtp: { host: 'smtp.office365.com', port: 587, secure: false }, // STARTTLS
    requiresAppPassword: false
  },
  {
    id: 'qq',
    domains: ['qq.com', 'foxmail.com'],
    name: 'QQ Mail',
    imap: { host: 'imap.qq.com', port: 993, secure: true },
    smtp: { host: 'smtp.qq.com', port: 465, secure: true },
    helpLink: 'https://service.mail.qq.com/detail/0/75',
    requiresAppPassword: true // Authorization code
  },
  {
    id: '163',
    domains: ['163.com'],
    name: 'NetEase 163',
    imap: { host: 'imap.163.com', port: 993, secure: true },
    smtp: { host: 'smtp.163.com', port: 465, secure: true },
    requiresAppPassword: true // Authorization code
  }
];

// Helper function to detect provider based on email address
export function detectProvider(email: string): MailProviderConfig | null {
  if (!email || !email.includes('@')) return null;
  const domain = email.split('@')[1]?.toLowerCase();
  if (!domain) return null;
  return MAIL_PROVIDERS.find(p => p.domains.includes(domain)) || null;
}
