import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Ensure AuthGuard doesn't redirect globally if not intended for test
    // AuthPage itself handles redirection based on isAuthenticated state
    await page.route('**/api/v1/auth/me', async (route) => {
      await route.fulfill({
        status: 401, // Initially not authenticated for auth tests
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Unauthorized' }),
      });
    });
    // Go to the auth page
    await page.goto('/auth');
  });

  test('User can switch between login and register modes and see i18n changes', async ({ page }) => {
    // Initial state: Login form
    await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Continue' })).toBeVisible();
    await expect(page.getByLabel('Full Name')).not.toBeVisible();

    // Switch to Register mode
    await page.getByRole('button', { name: 'Create an account' }).click(); // Using button with text from AuthSwitch
    await expect(page.getByRole('heading', { name: 'Create an account' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Create account' })).toBeVisible();
    await expect(page.getByLabel('Full Name')).toBeVisible();

    // Switch language to Chinese
    await page.getByTitle('Switch Language').click();
    await expect(page.getByRole('heading', { name: '创建新账号' })).toBeVisible();
    await expect(page.getByLabel('您的全名')).toBeVisible();
    
    // Switch back to English
    await page.getByTitle('Switch Language').click();
    await expect(page.getByRole('heading', { name: 'Create an account' })).toBeVisible();
  });

  test('Successful registration redirects to onboarding', async ({ page }) => {
    // Mock successful registration
    await page.route('**/api/v1/auth/register', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-jwt-token',
          user: { id: 'user-id-123', email: 'new@example.com', name: 'Test User', role: '', has_account: false },
        }),
      });
    });
    // Mock successful login after registration (AuthGuard check)
    await page.route('**/api/v1/auth/login', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-jwt-token',
          user: { id: 'user-id-123', email: 'new@example.com', name: 'Test User', role: '', has_account: false },
        }),
      });
    });

    await page.goto('/auth?mode=register');

    await page.getByLabel('Full Name').fill('Test User');
    await page.getByLabel('Email Address').fill('new@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByLabel('Confirm Password').fill('password123');

    await page.getByRole('button', { name: 'Create account' }).click();
    
    // Expect redirection to onboarding
    await expect(page).toHaveURL(/\/onboarding/);
  });

  test('Registration with mismatched passwords shows error', async ({ page }) => {
    await page.goto('/auth?mode=register');

    await page.getByLabel('Full Name').fill('Test User');
    await page.getByLabel('Email Address').fill('mismatch@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByLabel('Confirm Password').fill('wrongpassword');

    await page.getByRole('button', { name: 'Create account' }).click();
    
    // Expect error message for password mismatch
    await expect(page.getByText('Passwords do not match')).toBeVisible();
    await expect(page).toHaveURL(/\/auth\\?mode=register/); // Should stay on the same page
  });

  test('Failed registration shows API error', async ({ page }) => {
    // Mock failed registration
    await page.route('**/api/v1/auth/register', async (route) => {
      await route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'User already exists' }),
      });
    });

    await page.goto('/auth?mode=register');

    await page.getByLabel('Full Name').fill('Existing User');
    await page.getByLabel('Email Address').fill('existing@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByLabel('Confirm Password').fill('password123');

    await page.getByRole('button', { name: 'Create account' }).click();
    
    // Expect general error message
    await expect(page.getByText('Registration failed. Please try again.')).toBeVisible(); // This might show general error from AuthForm
    await expect(page).toHaveURL(/\/auth\\?mode=register/); // Should stay on the same page
  });

  test('Successful login redirects to dashboard', async ({ page }) => {
    // Mock successful login
    await page.route('**/api/v1/auth/login', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-jwt-token',
          user: { id: 'user-id-456', email: 'existing@example.com', name: 'Existing User', role: 'manager', has_account: true },
        }),
      });
    });

    await page.goto('/auth?mode=login');

    await page.getByLabel('Email Address').fill('existing@example.com');
    await page.getByLabel('Password').fill('password123');

    await page.getByRole('button', { name: 'Continue' }).click();
    
    // Expect redirection to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('Failed login shows API error', async ({ page }) => {
    // Mock failed login
    await page.route('**/api/v1/auth/login', async (route) => {
      await route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Invalid credentials' }),
      });
    });

    await page.goto('/auth?mode=login');

    await page.getByLabel('Email Address').fill('wrong@example.com');
    await page.getByLabel('Password').fill('wrongpassword');

    await page.getByRole('button', { name: 'Continue' }).click();
    
    // Expect error message for invalid credentials
    await expect(page.getByText('Invalid email or password. Please try again.')).toBeVisible();
    await expect(page).toHaveURL(/\/auth\\?mode=login/); // Should stay on the same page
  });
});
