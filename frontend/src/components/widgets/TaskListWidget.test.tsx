import { render, screen } from '@testing-library/react';
import { TaskListWidget } from './TaskListWidget';

describe('TaskListWidget', () => {
    it('renders list of tasks correctly', () => {
        const data = [
            { title: 'Task 1', dueDate: '2023-10-27', priority: 'High' as const },
            { title: 'Task 2', priority: 'Medium' as const },
        ];
        render(<TaskListWidget data={data} />);

        expect(screen.getByText('Task 1')).toBeInTheDocument();
        expect(screen.getByText('2023-10-27')).toBeInTheDocument();
        expect(screen.getByText('High')).toBeInTheDocument();

        expect(screen.getByText('Task 2')).toBeInTheDocument();
        expect(screen.getByText('Medium')).toBeInTheDocument();
    });

    it('renders error message for invalid data', () => {
        // @ts-expect-error - Testing error handling with invalid data
        render(<TaskListWidget data={null} />);
        expect(screen.getByText('Invalid task data')).toBeInTheDocument();
    });
});
