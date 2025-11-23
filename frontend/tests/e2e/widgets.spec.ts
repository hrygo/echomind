import { test, expect } from '@playwright/test';

test.describe('Generative Widgets', () => {
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

        if (url.includes('/chat/completions')) {
            // Mock a streaming response with a widget
            // Data must be SSE formatted
            const widgetData = {
                type: 'task_card',
                data: { title: 'Generated Task', priority: 'High' }
            };
            
            const chunk = {
                id: '1',
                choices: [{ index: 0, delta: { widget: widgetData } }]
            };
            
            const body = `data: ${JSON.stringify(chunk)}

data: [DONE]

`;
            
            return route.fulfill({
                status: 200,
                contentType: 'text/event-stream',
                body: body
            });
        }
        
        // Default catch-all to avoid 401s
        if (url.includes('/emails') || url.includes('/contexts') || url.includes('/tasks') || url.includes('/insights') || url.includes('/search')) {
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

  test('should render task widget in chat', async ({ page }) => {
    // 1. Open Chat Sidebar using the Sparkles icon in header
    await page.locator('header button[title="AI Copilot"]').click();
    
    // Wait for sidebar to open
    await expect(page.getByRole('heading', { name: 'EchoMind Copilot' })).toBeVisible();

    // 2. Type in Chat Input
    const chatInput = page.locator('input[placeholder="Ask anything..."]');
    await chatInput.fill('Create a task');
    await chatInput.press('Enter');
    
    // Expect the Widget to be rendered within the Chat Dialog
    const chatDialog = page.getByRole('dialog');
    await expect(chatDialog.locator('text=Generated Task')).toBeVisible();
    await expect(chatDialog.locator('text=High')).toBeVisible();
  });
});
