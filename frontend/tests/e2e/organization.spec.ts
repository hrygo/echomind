import { test, expect } from '@playwright/test';
import { mockApi } from './fixtures/mock-api';

test.describe('Organization Management', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page);
    
    // Login Flow
    await page.goto('/login');
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign in' }).click();
    await page.waitForURL('/dashboard');
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