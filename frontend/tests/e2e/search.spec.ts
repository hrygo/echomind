import { test, expect } from '@playwright/test';

const randomEmail = `e2e-${Math.random().toString(36).substring(7)}@example.com`;
const password = 'password123';

test.describe('Search Functionality', () => {
  test('should allow user to register, login, and search', async ({ page }) => {
    // 1. Register
    await page.goto('/register');
    await page.fill('#name', 'E2E User');
    await page.fill('#email', randomEmail);
    await page.fill('#password', password);
    await page.fill('#confirmPassword', password);
    await page.click('button[type="submit"]');

    // 2. Expect redirect to login
    await expect(page).toHaveURL(/login\?registered=true/);

    // 3. Login
    await page.fill('#email', randomEmail);
    await page.fill('#password', password);
    await page.click('button[type="submit"]');

    // 4. Expect redirect to dashboard
    // Wait for navigation to complete
    await page.waitForURL(/\/dashboard/);

    // 5. Check Search Bar Visibility
    const searchInput = page.locator('header input[type="text"]');
    await expect(searchInput).toBeVisible();

    // 6. Perform Search
    await searchInput.fill('Test Query');
    // Trigger search (focus or enter)
    await searchInput.press('Enter');

    // 7. Verify Search Results UI
    // Should show loading or results
    // Since we expect 0 results (empty DB for this user), check for "No results found"
    const resultsContainer = page.locator('text=No results found');
    await expect(resultsContainer).toBeVisible({ timeout: 10000 });
  });
});
