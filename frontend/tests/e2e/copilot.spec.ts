import { test, expect } from '@playwright/test';

test.describe('Smart Copilot', () => {
  test.beforeEach(async ({ page }) => {
    // Mock authentication
    await page.route('**/api/v1/auth/me', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ id: 'test-user-id', email: 'test@example.com' }),
      });
    });

    // Mock search API
    await page.route('**/api/v1/search**', async (route) => {
        const url = new URL(route.request().url());
        const query = url.searchParams.get('q');
        
        if (query === 'budget') {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({
                    results: [
                        {
                            email_id: 'email-1',
                            subject: 'Q4 Budget Plan',
                            snippet: 'Here is the draft for Q4 budget...', 
                            sender: 'finance@example.com',
                            date: new Date().toISOString(),
                            score: 0.95
                        }
                    ]
                })
            });
        } else {
             await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ results: [] })
            });
        }
    });

     // Mock Chat Stream API
     await page.route('**/api/v1/chat/completions', async (route) => {
         await route.fulfill({
             status: 200,
             contentType: 'text/event-stream',
             body: 'data: {"id":"1","choices":[{"delta":{"content":"Based on your budget emails..."}}]}\n\ndata: [DONE]\n\n'
         });
     });

    await page.goto('/dashboard');
  });

  test('Omni-Bar switches between Search and Chat modes', async ({ page }) => {
    // 1. Initial State: Omni-Bar is visible
    const omniBar = page.getByPlaceholder('Search emails, tasks, or contacts...');
    await expect(omniBar).toBeVisible();

    // 2. Perform Search
    await omniBar.fill('budget');
    await omniBar.press('Enter');

    // 3. Verify Search Results
    const resultItem = page.getByText('Q4 Budget Plan');
    await expect(resultItem).toBeVisible();
    await expect(page.getByText('finance@example.com')).toBeVisible();

    // 4. Switch to Chat (Click Sparkles Icon)
    const chatButton = page.getByTitle('Switch to Copilot Chat');
    await chatButton.click();

    // 5. Verify Chat Mode UI
    await expect(page.getByPlaceholder('Ask Copilot anything...')).toBeVisible();
    
    // 6. Send Chat Message
    await omniBar.fill('Summarize this');
    await omniBar.press('Enter');

    // 7. Verify Chat Response (Streaming Mock)
    await expect(page.getByText('Based on your budget emails...')).toBeVisible();
    
    // 8. Verify Context Attachment Indicator
    await expect(page.getByText('Context: 1 items attached')).toBeVisible();
  });

  test('Auto-Chat trigger with question mark', async ({ page }) => {
      const omniBar = page.getByPlaceholder('Search emails, tasks, or contacts...');
      
      // Type a question
      await omniBar.fill('How much is the budget?');
      await omniBar.press('Enter');

      // Should automatically switch to chat mode
      await expect(page.getByPlaceholder('Ask Copilot anything...')).toBeVisible();
      
      // And show the answer
      await expect(page.getByText('Based on your budget emails...')).toBeVisible();
  });
});
