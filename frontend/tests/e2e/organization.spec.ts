import { test, expect } from '@playwright/test';

test.describe('Organization Management', () => {
  test.beforeEach(async ({ page }) => {
    // Mock login or use a setup step to authenticate
    // For now, we assume a clean state or mock the auth response
    await page.goto('/login');
    
    // Fill login form (assuming test user exists or we mock the request)
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign in' }).click();
    
    // Wait for navigation to dashboard
    await page.waitForURL('/dashboard');
  });

  test('should create and switch to a new organization', async ({ page }) => {
    // 1. Open Org Switcher
    await page.getByRole('combobox').click();
    
    // 2. Click "Create New Organization"
    await page.getByText('Create New Organization').click();
    
    // 3. Fill Modal
    const orgName = `Test Org ${Date.now()}`;
    await page.getByLabel('Organization Name').fill(orgName);
    await page.getByRole('button', { name: 'Create' }).click();
    
    // 4. Verify Org is selected (Switcher text changes)
    await expect(page.getByRole('combobox')).toHaveText(orgName);
    
    // 5. Verify Org is in the list
    await page.getByRole('combobox').click();
    await expect(page.getByRole('menuitem', { name: orgName })).toBeVisible();
  });
});
