import { test, expect } from '@playwright/test';

test.describe('Omni-Bar Functionality', () => {
  test.beforeEach(async ({ page }) => {
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    // Universal Mock Handler
    await page.route('**/api/v1/**', async (route) => {
        const url = route.request().url();
        
        if (url.includes('/auth/login')) {
            const json = {
                token: 'mock-jwt-token',
                user: { id: 'mock-user-id', email: 'test@example.com', name: 'Test User' },
            };
            return route.fulfill({ json });
        }
        
        if (url.includes('/orgs')) {
             const json = [{ id: 'org-1', name: 'Personal', slug: 'personal', owner_id: 'mock-user-id' }];
              return route.fulfill({ json });
        }

        if (url.includes('/search')) {
             // Mock search results for "project alpha"
             const urlObj = new URL(url);
             const query = urlObj.searchParams.get('q');
             
             if (query === 'project alpha') {
                 const json = {
                    query: 'project alpha',
                    count: 1,
                    results: [
                      {
                        email_id: 'email-alpha',
                        subject: 'Project Alpha Update',
                        snippet: 'Status is green.',
                        sender: 'boss@example.com',
                        date: new Date().toISOString(),
                        score: 0.99,
                      },
                    ],
                  };
                  return route.fulfill({ json });
             }
             
             return route.fulfill({ json: { query, count: 0, results: [] } });
        }
        
        if (url.includes('/emails') || url.includes('/contexts') || url.includes('/tasks') || url.includes('/insights')) {
             return route.fulfill({ json: [] });
        }
        
        return route.fulfill({ status: 200, json: {} });
    });

    await page.goto('/');
    await page.evaluate(() => localStorage.setItem('app-language', 'en'));
    await page.goto('/login');
    await page.getByLabel('Email').fill('test@example.com');
    await page.getByLabel('Password').fill('password123');
    await page.getByRole('button', { name: 'Sign In' }).click();
    await page.waitForURL('/dashboard');
    await page.locator('h1', { hasText: 'EchoMind' }).waitFor({ state: 'visible' });
  });

  test('Smart Routing: Should open chat directly for questions', async ({ page }) => {
    const searchInput = page.locator('header input[type="text"]');
    
    // Type a question
    await searchInput.fill('How do I use this?');
    await searchInput.press('Enter');

    // Expect Chat Sidebar to open
    await expect(page.getByRole('heading', { name: 'EchoMind Copilot' })).toBeVisible();
    
    // Expect the question to be in the chat
    await expect(page.locator('text=How do I use this?')).toBeVisible();
  });

  test('Mixed Mode: Should search and then open chat with context', async ({ page }) => {
    const searchInput = page.locator('header input[type="text"]');
    
    // Type mixed query
    await searchInput.fill('project alpha and summarize status');
    await searchInput.press('Enter');

    // Expect Chat Sidebar to open
    await expect(page.getByRole('heading', { name: 'EchoMind Copilot' })).toBeVisible();

    // Expect Context Loaded message (from Search results)
    // "Loaded 1 emails into context"
    await expect(page.locator('text=Loaded 1 emails into context').first()).toBeVisible();

    // Expect the chat part ("status") to be sent as user message (keyword "and summarize" is stripped)
    await expect(page.locator('text=status').first()).toBeVisible();
  });
});
