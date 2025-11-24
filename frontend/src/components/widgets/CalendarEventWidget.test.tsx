import { render, screen } from '@testing-library/react';
import { CalendarEventWidget } from './CalendarEventWidget';

describe('CalendarEventWidget', () => {
    it('renders calendar event details correctly', () => {
        const data = {
            title: 'Team Meeting',
            start: '2023-10-27T10:00:00',
            end: '2023-10-27T11:00:00',
            location: 'Conference Room A',
        };
        render(<CalendarEventWidget data={data} />);

        expect(screen.getByText('Team Meeting')).toBeInTheDocument();
        expect(screen.getByText('Conference Room A')).toBeInTheDocument();
        // Check for date/time formatting parts
        expect(screen.getByText('27')).toBeInTheDocument(); // Day
        expect(screen.getByText('Oct')).toBeInTheDocument(); // Month
    });
});
