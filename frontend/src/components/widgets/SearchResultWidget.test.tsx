import { render, screen } from '@testing-library/react';
import { SearchResultWidget } from './SearchResultWidget';

describe('SearchResultWidget', () => {
    it('renders search results correctly', () => {
        const data = {
            count: 5,
            results: [
                {
                    subject: 'Email 1',
                    sender: 'sender1@test.com',
                    date: '2023-01-01',
                    snippet: 'Snippet 1',
                },
            ],
        };
        render(<SearchResultWidget data={data} />);
        expect(screen.getByText('Found 5 emails')).toBeInTheDocument();
        expect(screen.getByText('Email 1')).toBeInTheDocument();
        expect(screen.getByText('sender1@test.com')).toBeInTheDocument();
    });
});
