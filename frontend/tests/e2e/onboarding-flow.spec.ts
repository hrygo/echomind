import { test, expect } from '@playwright/test';

test.describe('Onboarding Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Mock authentication: Assume user is authenticated but needs onboarding
    await page.route('**/api/v1/auth/me', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'onboarding-user-id', 
          email: 'onboard@example.com', 
          name: 'Onboard User', 
          role: '', 
          has_account: false, // Key for onboarding
        }),
      });
    });

    // Mock patch user role API
    await page.route('**/api/v1/users/me', async (route) => {
        await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ message: 'User role updated' }),
        });
    });

    // Mock successful mailbox connection
    await page.route('**/api/v1/settings/account', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Account connected successfully', account_id: 'acc-id-123' }),
      });
    });

    // Mock successful initial sync
    await page.route('**/api/v1/sync', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Sync initiated successfully' }),
      });
    });

    await page.goto('/onboarding');
  });

  test('New user completes full onboarding flow', async ({ page }) => {
    // Step 1: Role Selection
    await expect(page.getByRole('heading', { name: 'Tell us about your role' })).toBeVisible();
    await page.getByRole('button', { name: 'Manager' }).click();
    await page.getByRole('button', { name: 'Next Step' }).click();

    // Step 2: Smart Mailbox Form
    await expect(page.getByRole('heading', { name: 'Connect your neural network' })).toBeVisible();
    await expect(page.getByLabel('Email Address')).toHaveValue('onboard@example.com'); // Pre-filled
    await page.getByLabel('Password / App Password').fill('app_password_123');
    
    // Assume generic provider detection, or specifically test Gmail
    await page.getByLabel('Email Address').fill('onboard@gmail.com'); // Change to trigger Gmail provider
    await expect(page.getByText('Using Gmail? You likely need an App Password.')).toBeVisible();

    await page.getByRole('button', { name: 'Connect Mailbox' }).click();
    await expect(page.getByText('Verifying credentials...')).toBeVisible();
    await expect(page.getByText('Connected successfully!')).toBeVisible();

    // Step 3: Initial Sync
    await expect(page.getByRole('heading', { name: 'Initializing EchoMind' })).toBeVisible();
    await expect(page.getByText('Syncing recent emails and building your knowledge graph...')).toBeVisible();

    // Should automatically redirect to dashboard after sync
    await page.waitForURL('/dashboard', { timeout: 5000 });
    await expect(page).toHaveURL('/dashboard');
  });

  test('Mailbox connection failure stays on step 2 with error message', async ({ page }) => {
    // Mock failed mailbox connection
    await page.route('**/api/v1/settings/account', async (route) => {
      await route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'IMAP login failed. Check password.' }),
      });
    });

    await page.goto('/onboarding');

    // Step 1: Role Selection (pass quickly)
    await page.getByRole('button', { name: 'Executive' }).click();
    await page.getByRole('button', { name: 'Next Step' }).click();

    // Step 2: Smart Mailbox Form - enter details that will fail
    await page.getByLabel('Email Address').fill('fail@gmail.com');
    await page.getByLabel('Password / App Password').fill('wrong_password');
    await page.getByRole('button', { name: 'Connect Mailbox' }).click();

    // Expect to stay on step 2 and see error
    await expect(page.getByRole('heading', { name: 'Connect your neural network' })).toBeVisible();
    await expect(page.getByText('Failed to connect to mailbox. Please check your credentials and server settings.')).toBeVisible();
    await expect(page.getByText('IMAP login failed. Check password.')).toBeVisible(); // Specific backend error
  });

  test('Unauthenticated user is redirected to login', async ({ page }) => {
    // Mock unauthenticated status for this specific test case
    await page.route('**/api/v1/auth/me', async (route) => {
      await route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Unauthorized' }),
      });
    });
    await page.goto('/onboarding');
    await expect(page).toHaveURL(/\/auth\?mode=login/);
  });
});
