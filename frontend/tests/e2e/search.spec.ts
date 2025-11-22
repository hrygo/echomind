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

    // 8. Test Search History
    // Clear input to show history
    await searchInput.fill('');
    // History should be visible now
    const historyContainer = page.locator('text=Recent Searches');
    await expect(historyContainer).toBeVisible();
    await expect(page.locator('text=Test Query')).toBeVisible();

    // 9. Select from history
    await page.click('text=Test Query');
    // Should trigger search again
    await expect(resultsContainer).toBeVisible();
    
    // 10. Clear History
    await searchInput.fill('');
    await expect(historyContainer).toBeVisible();
    await page.click('text=Clear');
    // History should disappear (or be empty)
    await expect(historyContainer).not.toBeVisible();

    // 11. Test Filters
    // Open Filter Menu
    await page.click('button[title="Search Filters"]');
    const filterMenu = page.locator('text=Search Filters');
    await expect(filterMenu).toBeVisible();

    // Set Sender Filter
    await page.fill('input[placeholder="e.g. alice@example.com"]', 'test@example.com');
    
    // Close Filter Menu
    await page.click('button[title="Search Filters"]'); // Toggle off or click X
    // (Note: The X button is inside the menu, the toggle button is outside)
    
    // Check if filter indicator is active (dot)
    const activeIndicator = page.locator('button[title="Search Filters"] span.bg-blue-600');
    await expect(activeIndicator).toBeVisible();

    // Perform search with filter
    await searchInput.fill('Filtered Query');
    await searchInput.press('Enter');
    await expect(resultsContainer).toBeVisible();
  });
});
