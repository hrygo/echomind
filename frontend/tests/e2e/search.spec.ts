import { test, expect } from '@playwright/test';


test.describe('Search Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Enable console logging for debugging
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    page.on('pageerror', exception => console.log(`PAGE EXCEPTION: ${exception}`));

    // Universal Mock Handler to prevent 401s
    await page.route('**/api/v1/**', async (route) => {
        const url = route.request().url();
        
        if (url.includes('/auth/login')) {
            console.log('MOCK HIT: Login');
            const json = {
                token: 'mock-jwt-token',
                user: { id: 'mock-user-id', email: 'test@example.com', name: 'Test User' },
            };
            return route.fulfill({ json });
        }
        
        if (url.includes('/orgs')) {
             console.log('MOCK HIT: Orgs');
             const json = [
                {
                  id: 'org-1',
                  name: 'Personal Workspace',
                  slug: 'personal-workspace',
                  owner_id: 'mock-user-id',
                },
              ];
              return route.fulfill({ json });
        }

        if (url.includes('/search')) {
             console.log('MOCK HIT: Search');
             const json = {
                query: 'test',
                count: 1,
                results: [
                  {
                    email_id: 'email-1',
                    subject: 'Mock Email Subject',
                    snippet: 'This is a mock email snippet...',
                    sender: 'sender@example.com',
                    date: new Date().toISOString(),
                    score: 0.95,
                  },
                ],
              };
              return route.fulfill({ json });
        }
        
        // Mock specific endpoints that might cause issues if empty object is returned
        if (url.includes('/emails') || url.includes('/contexts') || url.includes('/tasks') || url.includes('/insights')) {
             return route.fulfill({ json: [] }); // Return empty array for lists
        }
        
        // Default catch-all
        console.log('MOCK CATCH-ALL (200 OK):', url);
        return route.fulfill({ status: 200, json: {} });
    });

    // Force English language
    await page.goto('/');
    await page.evaluate(() => localStorage.setItem('app-language', 'en'));

    // Go to login
    await page.goto('/login');

    // Perform Mock Login
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign In' }).click();

    // Wait for dashboard
    await page.waitForURL('/dashboard');
    console.log('Navigated to /dashboard');

    await page.waitForLoadState('domcontentloaded');
    await page.locator('h1', { hasText: 'EchoMind' }).waitFor({ state: 'visible' }); 
  });

  test('should perform search and view results', async ({ page }) => {
    // 1. Check Search Bar Visibility
    const searchInput = page.locator('header input[type="text"]');
    await expect(searchInput).toBeVisible();

    // 2. Perform Search
    await searchInput.fill('Test Query');
    await searchInput.press('Enter');

    // 3. Verify Search Results UI (Mocked Response)
    await expect(page.locator('text=Mock Email Subject')).toBeVisible();
    await expect(page.locator('text=This is a mock email snippet...')).toBeVisible();
  });

  test('should open chat sidebar with context when Ask Copilot is clicked', async ({ page }) => {
    const searchInput = page.locator('header input[type="text"]');
    await searchInput.fill('Test Query');
    await searchInput.press('Enter');

    // Wait for results
    await expect(page.locator('text=Mock Email Subject')).toBeVisible();

    // Click "Ask Copilot" button
    const askCopilotBtn = page.getByRole('button', { name: 'Ask Copilot' });
    await expect(askCopilotBtn).toBeVisible();
    await askCopilotBtn.evaluate(b => b.click());

    // Wait for animation
    await page.waitForTimeout(1000);

    // Verify Chat Sidebar opens
    await expect(page.locator('text=EchoMind Copilot')).toBeVisible();
    
    // Verify context loaded message
    // "Loaded 1 emails into context"
    await expect(page.locator('text=Loaded 1 emails into context').first()).toBeVisible();
  });
});
