import { test, expect } from '@playwright/test';
import { mockApi } from './fixtures/mock-api';

test.describe.skip('Organization Management', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page);

    // Catch all other API requests to avoid unmocked errors
    await page.route('**/api/v1/**', route => route.fulfill({ status: 404, body: 'Not Found' }));
    
    // Login Flow - Directly go to dashboard, as mocks handle login state
    await page.goto('/dashboard');
    await page.waitForURL('/dashboard');
    await page.locator('h1').waitFor({ state: 'visible' }); // Wait for the main dashboard heading
  });

  test('should create and switch to a new organization', async ({ page }) => {
    // 1. Open Org Switcher
    // Note: The mock returns 2 orgs initially: "Personal Workspace" and "Acme Corp".
    // We assume the first one is selected by default or we select one.
    
    await page.getByRole('combobox').click();
    
    // Verify initial list from mock
    await expect(page.getByRole('menuitem', { name: 'Personal Workspace' })).toBeVisible();
    await expect(page.getByRole('menuitem', { name: 'Acme Corp' })).toBeVisible();

    // 2. Click "Create New Organization"
    await page.getByText('Create New Organization').click();
    
    // 3. Fill Modal
    const orgName = 'New E2E Org';
    await page.getByLabel('Organization Name').fill(orgName);
    await page.getByRole('button', { name: 'Create' }).click();
    
    // 4. Verify Org is selected (Switcher text changes)
    // The mock API for POST /orgs returns the new org, and our frontend code automatically switches to it.
    await expect(page.getByRole('combobox')).toHaveText(orgName);
    
    // 5. Verify Org is in the list
    await page.getByRole('combobox').click();
    await expect(page.getByRole('menuitem', { name: orgName })).toBeVisible();
  });
});
