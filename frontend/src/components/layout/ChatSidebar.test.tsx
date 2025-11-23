import { render, screen, waitFor } from '@testing-library/react';
import { ChatSidebar } from './ChatSidebar';
import { useChatStore } from '@/lib/store/chat';
import { useAuthStore } from '@/store/auth';

// Mock dependencies
jest.mock('@/lib/store/chat');
jest.mock('@/store/auth', () => ({
    useAuthStore: {
        getState: jest.fn(),
    },
}));
jest.mock('react-markdown', () => ({ children }: { children: React.ReactNode }) => <div>{children}</div>);
jest.mock('@/components/ui/Sheet', () => ({
  Sheet: ({ children, open }: { children: React.ReactNode; open: boolean }) => open ? <div>{children}</div> : null,
  SheetContent: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SheetHeader: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SheetTitle: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));
jest.mock('@/lib/i18n/LanguageContext', () => ({
  useLanguage: () => ({ 
    t: (key: string) => {
        if (key === 'chat.contextLoaded') return 'Context loaded {count}';
        return key;
    } 
  }),
}));

describe('ChatSidebar', () => {
  const mockSetOpen = jest.fn();
  const mockAddMessage = jest.fn();
  const mockClearActiveContextEmails = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    window.HTMLElement.prototype.scrollIntoView = jest.fn();
    (useChatStore as unknown as jest.Mock).mockReturnValue({
      isOpen: true,
      setOpen: mockSetOpen,
      messages: [],
      addMessage: mockAddMessage,
      isLoading: false,
      setLoading: jest.fn(),
      updateLastMessage: jest.fn(),
      activeContextEmails: [],
      clearActiveContextEmails: mockClearActiveContextEmails,
    });
    (useAuthStore.getState as jest.Mock).mockReturnValue({
      token: 'mock-token',
    });
  });

  it('renders default message when empty', () => {
    render(<ChatSidebar />);
    expect(screen.getByText('How can I help you today?')).toBeInTheDocument();
  });

  it('adds system message when activeContextEmails are present', async () => {
    const mockEmails = [
      {
        ID: '1',
        Subject: 'Test Subject',
        Sender: 'sender@example.com',
        Snippet: 'Test Snippet',
        BodyText: 'Body',
        Date: '2023-01-01',
        Summary: 'Summary',
        Category: 'Work',
        Sentiment: 'Neutral',
        Urgency: 'Low',
        IsRead: true,
        ActionItems: [],
        SmartActions: {},
      },
    ];

    (useChatStore as unknown as jest.Mock).mockReturnValue({
      isOpen: true,
      setOpen: mockSetOpen,
      messages: [],
      addMessage: mockAddMessage,
      isLoading: false,
      setLoading: jest.fn(),
      updateLastMessage: jest.fn(),
      activeContextEmails: mockEmails,
      clearActiveContextEmails: mockClearActiveContextEmails,
    });

    render(<ChatSidebar />);

    await waitFor(() => {
        // Expect the system message to be added
        // The mock t returns 'Context loaded {count}' -> 'Context loaded 1'
        // The logic appends contextSummary
      expect(mockAddMessage).toHaveBeenCalledWith(expect.objectContaining({
        role: 'system',
        content: expect.stringContaining('Subject: Test Subject'),
      }));
    });
    
    expect(mockAddMessage).toHaveBeenCalledWith(expect.objectContaining({
        content: expect.stringContaining('Context loaded 1'),
    }));

    expect(mockClearActiveContextEmails).toHaveBeenCalled();
    expect(mockSetOpen).toHaveBeenCalledWith(true);
  });
});
