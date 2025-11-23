import { render, screen } from '@testing-library/react';
import { TaskWidget } from './TaskWidget';

describe('TaskWidget', () => {
    it('renders task details correctly', () => {
        const data = {
            title: 'Test Task',
            dueDate: '2023-12-31',
            priority: 'High' as const,
        };
        render(<TaskWidget data={data} />);
        expect(screen.getByText('Test Task')).toBeInTheDocument();
        expect(screen.getByText('2023-12-31')).toBeInTheDocument();
        expect(screen.getByText('High')).toBeInTheDocument();
    });
});
