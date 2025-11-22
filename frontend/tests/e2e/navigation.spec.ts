import { test, expect } from '@playwright/test';

test.describe('Dashboard Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Mock Login
    await page.route(url => url.href.includes('/api/v1/auth/login'), async (route) => {
      const json = {
        token: 'mock-jwt-token',
        user: { id: 'mock-user-id', email: 'test@example.com', name: 'Test User' },
      };
      await route.fulfill({ json });
    });

    // Mock Orgs
    await page.route(url => url.href.includes('/api/v1/orgs'), async (route) => {
      const json = [
        { id: 'org-1', name: 'Personal Workspace', slug: 'personal-workspace', owner_id: 'mock-user-id' },
      ];
      await route.fulfill({ json });
    });

    // Mock Inbox (for Smart Inbox nav)
    await page.route(url => url.href.includes('/api/v1/emails'), async (route) => {
        await route.fulfill({ json: [] });
    });

    // Force English
    await page.goto('/');
    await page.evaluate(() => localStorage.setItem('app-language', 'en'));

    // Login Flow
    await page.goto('/login');
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign In' }).click();
    
    await page.waitForURL('/dashboard');
    await page.locator('h1', { hasText: 'EchoMind' }).waitFor({ state: 'visible' });
  });

  test('should navigate through sidebar links', async ({ page }) => {
    // 1. Default View (Dashboard)
    await expect(page).toHaveURL('/dashboard');
    
    // 2. Navigate to Smart Inbox
    await page.click('text=Smart Inbox');
    await expect(page).toHaveURL(/\/dashboard\/inbox\?filter=smart/);
    // Smart Inbox title might be different depending on translation, let's check URL mainly
    
    // 3. Navigate to Action Items
    await page.click('text=Action Items');
    await expect(page).toHaveURL(/\/dashboard\/tasks/);

    // 4. Navigate to Network
    await page.click('text=Network');
    await expect(page).toHaveURL(/\/dashboard\/insights/);

    // 5. Navigate to Settings
    await page.goto('/dashboard/settings');
    // Settings page title is "Settings Center" or similar in en.json
    await expect(page.locator('h2')).toContainText('Settings');
  });
});
