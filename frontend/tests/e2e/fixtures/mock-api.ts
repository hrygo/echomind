import { Page } from '@playwright/test';

/**
 * Mock API responses for E2E tests.
 * This helps in running tests without a live backend.
 */
export const mockApi = async (page: Page) => {
  // Mock Auth Register
  await page.route('**/api/v1/auth/register', async (route) => {
    const json = {
      id: 'mock-user-id',
      email: 'test@example.com',
      name: 'Test User',
    };
    await route.fulfill({ json, status: 201 });
  });

  // Mock Auth Login
  await page.route('**/api/v1/auth/login', async (route) => {
    console.log('MOCK HIT: Login');
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
    } else if (route.request().method() === 'POST') {
      // Mock Create Org
      const requestData = route.request().postDataJSON();
      const json = {
        id: `org-${Date.now()}`,
        name: requestData.name,
        slug: requestData.name.toLowerCase().replace(/\s+/g, '-'),
        owner_id: 'mock-user-id',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
      await route.fulfill({ json, status: 201 });
    } else {
      await route.continue(); 
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
  
  // Mock Email Detail (for Search result click)
  await page.route('**/api/v1/emails/*', async (route) => {
      const json = {
          ID: 'email-1',
          Subject: 'Mock Email Subject',
          Sender: 'sender@example.com',
          Date: new Date().toISOString(),
          BodyText: 'This is a mock email body content.',
          Summary: 'Mock Summary',
          Sentiment: 'Positive',
          Urgency: 'Low',
          Category: 'Work',
          ActionItems: ['Reply']
      };
      await route.fulfill({ json });
  });
};