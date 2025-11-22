import { test, expect } from '@playwright/test';
import { mockApi } from './fixtures/mock-api';

test.describe('Search Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Enable console logging for debugging
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    page.on('pageerror', exception => console.log(`PAGE EXCEPTION: ${exception}`));

    // Mock Login
    await page.route(url => url.href.includes('/api/v1/auth/login'), async (route) => {
      console.log('MOCK HIT: Login');
      const json = {
        token: 'mock-jwt-token',
        user: { id: 'mock-user-id', email: 'test@example.com', name: 'Test User' },
      };
      await route.fulfill({ json });
    });

    // Mock Orgs
    await page.route(url => url.href.includes('/api/v1/orgs'), async (route) => {
      console.log('MOCK HIT: Orgs');
      const json = [
        {
          id: 'org-1',
          name: 'Personal Workspace',
          slug: 'personal-workspace',
          owner_id: 'mock-user-id',
        },
      ];
      await route.fulfill({ json });
    });

    // Mock Search
    await page.route(url => url.href.includes('/api/v1/search'), async (route) => {
        const json = {
          query: 'test',
          count: 1,
          results: [
            {
              email_id: 'email-1',
              subject: 'Mock Email Subject',
              snippet: 'This is a mock email snippet...',
              sender: 'sender@example.com',
              date: new Date().toISOString(),
              score: 0.95,
            },
          ],
        };
        await route.fulfill({ json });
    });

    // Force English language
    await page.goto('/');
    await page.evaluate(() => localStorage.setItem('app-language', 'en'));

    // Go to login
    await page.goto('/login');
    
    // Perform Mock Login
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign In' }).click();

    // Wait for dashboard
    await page.waitForURL('/dashboard');
    console.log('Navigated to /dashboard');
    
    await page.waitForLoadState('domcontentloaded');
    await page.locator('h1', { hasText: 'EchoMind' }).waitFor({ state: 'visible' }); // Wait for the main dashboard heading (EchoMind)
  });

  test('should perform search and view results', async ({ page }) => {
    // 1. Check Search Bar Visibility
    const searchInput = page.locator('header input[type="text"]');
    await expect(searchInput).toBeVisible();

    // 2. Perform Search
    await searchInput.fill('Test Query');
    await searchInput.press('Enter');

    // 3. Verify Search Results UI (Mocked Response)
    await expect(page.locator('text=Mock Email Subject')).toBeVisible();
    await expect(page.locator('text=This is a mock email snippet...')).toBeVisible();
  });
});
