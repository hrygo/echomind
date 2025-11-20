import { render, screen } from '@testing-library/react';
import Home from './page'; // Assuming page.tsx exports default function Home

describe('Home Page', () => {
  it('renders the main heading', () => {
    render(<Home />);
    const heading = screen.getByRole('heading', { name: /EchoMind/i });
    expect(heading).toBeInTheDocument();
  });

  it('renders the slogan', () => {
    render(<Home />);
    const slogan = screen.getByText(/Your external brain for email./i);
    expect(slogan).toBeInTheDocument();
  });

  it('renders a link to the dashboard', () => {
    render(<Home />);
    const link = screen.getByRole('link', { name: /Go to Dashboard/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', '/dashboard');
  });
});
