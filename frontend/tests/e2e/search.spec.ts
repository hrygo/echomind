import { test, expect } from '@playwright/test';
import { mockApi } from './fixtures/mock-api';

test.describe('Search Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Enable Mocks
    await mockApi(page);

    // Go to login (auth guard will redirect if not logged in)
    await page.goto('/login');
    
    // Perform Mock Login
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign in' }).click();
    
    // Wait for dashboard
    await page.waitForURL(/\/dashboard/);
  });

  test('should perform search and view results', async ({ page }) => {
    // 1. Check Search Bar Visibility
    const searchInput = page.locator('header input[type="text"]');
    await expect(searchInput).toBeVisible();

    // 2. Perform Search
    await searchInput.fill('Test Query');
    await searchInput.press('Enter');

    // 3. Verify Search Results UI (Mocked Response)
    // Should show the result from mock-api
    await expect(page.locator('text=Mock Email Subject')).toBeVisible();
    await expect(page.locator('text=This is a mock email snippet...')).toBeVisible();

    // 4. Test Search History
    // Clear input to show history
    await searchInput.fill('');
    // History should be visible now
    const historyContainer = page.locator('text=Recent Searches');
    await expect(historyContainer).toBeVisible();
    await expect(page.locator('text=Test Query')).toBeVisible();

    // 5. Select from history
    await page.click('text=Test Query');
    // Should trigger search again
    await expect(page.locator('text=Mock Email Subject')).toBeVisible();
    
    // 6. Click Result to navigate (Mock Detail API needed)
    await page.click('text=Mock Email Subject');
    await page.waitForURL(/\/dashboard\/email\/email-1/);
    await expect(page.locator('h1')).toHaveText('Mock Email Subject');
  });
});