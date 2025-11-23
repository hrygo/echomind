import { test, expect } from '@playwright/test';

test.describe('Dashboard Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Universal Mock Handler
    await page.route('**/api/v1/**', async (route) => {
        const url = route.request().url();
        if (url.includes('/auth/login')) {
            return route.fulfill({ json: { token: 'mock', user: { id: '1', email: 'test@test.com' } } });
        }
        if (url.includes('/orgs')) {
             return route.fulfill({ json: [{ id: 'org-1', name: 'Personal', slug: 'personal', owner_id: '1' }] });
        }
        // Return empty lists for other resources
        if (url.includes('/emails') || url.includes('/contexts') || url.includes('/tasks') || url.includes('/insights')) {
             return route.fulfill({ json: [] });
        }
        // Catch-all
        return route.fulfill({ status: 200, json: {} });
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
