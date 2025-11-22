import { Page } from '@playwright/test';

/**
 * Mock API responses for E2E tests.
 * This helps in running tests without a live backend.
 */
export const mockApi = async (page: Page) => {
  // Mock Auth
  await page.route('**/api/v1/auth/login', async (route) => {
    const json = {
      token: 'mock-jwt-token',
      user: {
        id: 'mock-user-id',
        email: 'test@example.com',
        name: 'Test User',
      },
    };
    await route.fulfill({ json });
  });

  // Mock Organizations List
  await page.route('**/api/v1/orgs', async (route) => {
    if (route.request().method() === 'GET') {
      const json = [
        {
          id: 'org-1',
          name: 'Personal Workspace',
          slug: 'personal-workspace',
          owner_id: 'mock-user-id',
        },
        {
          id: 'org-2',
          name: 'Acme Corp',
          slug: 'acme-corp',
          owner_id: 'mock-user-id',
        },
      ];
      await route.fulfill({ json });
    } else {
      await route.continue(); // Let POST/PUT pass through or mock specifically
    }
  });

  // Mock Search
  await page.route('**/api/v1/search*', async (route) => {
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
    await route.fulfill({ json });
  });
};
