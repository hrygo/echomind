/**
 * E2E Test: Server Actions
 * 测试 React 19 Server Actions 功能
 */

import { test, expect } from '@playwright/test';

test.describe('Server Actions - Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
  });

  test('should handle login action successfully', async ({ page }) => {
    // 填写登录表单
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');

    // 提交表单（触发 loginAction）
    await page.click('button[type="submit"]');

    // 等待导航完成
    await page.waitForURL('**/dashboard', { timeout: 10000 });

    // 验证登录成功
    expect(page.url()).toContain('/dashboard');
  });

  test('should display error for invalid credentials', async ({ page }) => {
    // 填写错误的登录凭证
    await page.fill('input[name="email"]', 'wrong@example.com');
    await page.fill('input[name="password"]', 'wrongpassword');

    // 提交表单
    await page.click('button[type="submit"]');

    // 等待错误消息显示
    await page.waitForSelector('[role="alert"]', { timeout: 5000 });

    // 验证错误消息
    const errorMessage = await page.textContent('[role="alert"]');
    expect(errorMessage).toBeTruthy();
  });

  test('should handle register action successfully', async ({ page }) => {
    // 导航到注册页面
    await page.goto('/register');

    // 填写注册表单
    await page.fill('input[name="email"]', `test${Date.now()}@example.com`);
    await page.fill('input[name="password"]', 'password123');
    await page.fill('input[name="name"]', 'Test User');

    // 提交表单（触发 registerAction）
    await page.click('button[type="submit"]');

    // 等待成功消息或导航
    await page.waitForURL('**/dashboard', { timeout: 10000 });

    // 验证注册成功
    expect(page.url()).toContain('/dashboard');
  });

  test('should handle logout action', async ({ page, context }) => {
    // 先登录
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });

    // 执行登出
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');

    // 等待导航到登录页
    await page.waitForURL('**/login', { timeout: 5000 });

    // 验证已登出
    expect(page.url()).toContain('/login');

    // 验证认证 Cookie 已清除
    const cookies = await context.cookies();
    const authToken = cookies.find(c => c.name === 'token');
    expect(authToken).toBeUndefined();
  });
});

test.describe('Server Actions - Email Operations', () => {
  test.beforeEach(async ({ page }) => {
    // 登录并导航到邮件页面
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/emails');
  });

  test('should sync emails using server action', async ({ page }) => {
    // 点击同步按钮（触发 syncEmailsAction）
    await page.click('[data-testid="sync-emails-button"]');

    // 等待加载指示器出现
    await page.waitForSelector('[data-testid="loading-indicator"]', { 
      timeout: 3000 
    }).catch(() => {
      // 如果同步太快，加载指示器可能不会出现
    });

    // 等待加载完成
    await page.waitForSelector('[data-testid="loading-indicator"]', { 
      state: 'hidden',
      timeout: 15000 
    }).catch(() => {
      // 已经隐藏
    });

    // 验证邮件列表已更新
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });
    const emailItems = await page.locator('[data-testid="email-item"]').count();
    expect(emailItems).toBeGreaterThan(0);
  });

  test('should delete email using server action', async ({ page }) => {
    // 等待邮件列表加载
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });

    // 获取第一封邮件的ID
    const firstEmail = page.locator('[data-testid="email-item"]').first();
    await firstEmail.click();

    // 点击删除按钮（触发 deleteEmailAction）
    await page.click('[data-testid="delete-email-button"]');

    // 确认删除对话框
    await page.click('[data-testid="confirm-delete-button"]');

    // 等待成功消息
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('删除成功');
  });

  test('should archive email using server action', async ({ page }) => {
    // 等待邮件列表加载
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });

    // 选择第一封邮件
    const firstEmail = page.locator('[data-testid="email-item"]').first();
    await firstEmail.click();

    // 点击归档按钮（触发 archiveEmailAction）
    await page.click('[data-testid="archive-email-button"]');

    // 等待成功消息
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('归档成功');
  });
});

test.describe('Server Actions - Organization Management', () => {
  test.beforeEach(async ({ page }) => {
    // 登录并导航到组织设置页面
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/settings/organization');
  });

  test('should create organization using server action', async ({ page }) => {
    // 点击创建组织按钮
    await page.click('[data-testid="create-org-button"]');

    // 填写组织信息
    await page.fill('input[name="name"]', `Test Org ${Date.now()}`);
    await page.fill('textarea[name="description"]', 'Test organization description');

    // 提交表单（触发 createOrganizationAction）
    await page.click('button[type="submit"]');

    // 等待成功消息
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('创建成功');
  });

  test('should update organization using server action', async ({ page }) => {
    // 等待组织信息加载
    await page.waitForSelector('input[name="name"]', { timeout: 5000 });

    // 修改组织名称
    await page.fill('input[name="name"]', `Updated Org ${Date.now()}`);

    // 提交表单（触发 updateOrganizationAction）
    await page.click('[data-testid="save-org-button"]');

    // 等待成功消息
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('更新成功');
  });

  test('should invite member using server action', async ({ page }) => {
    // 导航到成员管理
    await page.click('[data-testid="members-tab"]');

    // 点击邀请按钮
    await page.click('[data-testid="invite-member-button"]');

    // 填写邀请信息
    await page.fill('input[name="email"]', `newmember${Date.now()}@example.com`);
    await page.selectOption('select[name="role"]', 'member');

    // 提交表单（触发 inviteMemberAction）
    await page.click('button[type="submit"]');

    // 等待成功消息
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('邀请成功');
  });
});

test.describe('Server Actions - Form Validation', () => {
  test('should display validation errors from server actions', async ({ page }) => {
    await page.goto('/login');

    // 提交空表单
    await page.click('button[type="submit"]');

    // 等待验证错误显示
    await page.waitForSelector('[role="alert"]', { timeout: 3000 });

    // 验证错误消息
    const errorMessages = await page.locator('[role="alert"]').allTextContents();
    expect(errorMessages.length).toBeGreaterThan(0);
  });

  test('should handle server-side validation errors', async ({ page }) => {
    await page.goto('/register');

    // 填写无效的邮箱
    await page.fill('input[name="email"]', 'invalid-email');
    await page.fill('input[name="password"]', '123'); // 密码太短

    // 提交表单
    await page.click('button[type="submit"]');

    // 等待验证错误
    await page.waitForSelector('[role="alert"]', { timeout: 3000 });

    // 验证错误消息
    const errorMessage = await page.textContent('[role="alert"]');
    expect(errorMessage).toContain('邮箱格式');
  });
});
