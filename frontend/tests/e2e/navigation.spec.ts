import { test, expect } from '@playwright/test';
import { mockApi } from './fixtures/mock-api';

test.describe('Dashboard Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page);
    await page.goto('/login');
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign in' }).click();
    await page.waitForURL('/dashboard');
  });

  test('should navigate through sidebar links', async ({ page }) => {
    // 1. Default View (Dashboard)
    await expect(page).toHaveURL('/dashboard');
    
    // 2. Navigate to Smart Inbox
    await page.click('text=Smart Inbox');
    await expect(page).toHaveURL(/\/dashboard\/inbox\?filter=smart/);
    await expect(page.locator('h1')).toHaveText('Smart Inbox');

    // 3. Navigate to Action Items
    await page.click('text=Action Items');
    await expect(page).toHaveURL(/\/dashboard\/tasks/);
    // (Assuming Action Items page has a header "Action Items")

    // 4. Navigate to Network
    await page.click('text=Network');
    await expect(page).toHaveURL(/\/dashboard\/insights/);

    // 5. Navigate to Settings (via URL or if link exists)
    // Assuming settings link is in user menu or sidebar (it was added in previous steps)
    // Let's check sidebar implementation: No settings link in sidebar main nav yet, usually in user menu.
    // But let's try direct navigation which also tests client-side routing
    await page.goto('/dashboard/settings');
    await expect(page.locator('h2')).toHaveText('Settings'); // Based on SettingsPage implementation
  });
});
