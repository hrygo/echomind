import { render, screen } from '@testing-library/react';
import { EmailDraftWidget } from './EmailDraftWidget';

describe('EmailDraftWidget', () => {
    it('renders email draft details correctly', () => {
        const data = {
            to: 'alice@example.com',
            subject: 'Project Update',
            body: 'Here is the update...',
        };
        render(<EmailDraftWidget data={data} />);

        expect(screen.getByText('alice@example.com')).toBeInTheDocument();
        expect(screen.getByText('Project Update')).toBeInTheDocument();
        expect(screen.getByText('Here is the update...')).toBeInTheDocument();
        expect(screen.getByText('Draft Email')).toBeInTheDocument();
    });
});
