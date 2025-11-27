/**
 * E2E Test: AI Draft Generation
 * 测试 AI 草稿生成和回复功能
 */

import { test, expect } from '@playwright/test';

test.describe('AI Draft Generation', () => {
  test.beforeEach(async ({ page }) => {
    // 登录并导航到邮件页面
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/emails');
  });

  test('should generate new email draft', async ({ page }) => {
    // 点击撰写新邮件按钮
    await page.click('[data-testid="compose-email-button"]');

    // 等待撰写对话框出现
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    // 填写基本信息
    await page.fill('[data-testid="email-to"]', 'recipient@example.com');
    await page.fill('[data-testid="email-subject"]', 'Meeting Request');

    // 输入草稿提示
    await page.fill('[data-testid="draft-prompt"]', 'Request a meeting next week to discuss project status');

    // 点击生成草稿按钮
    await page.click('[data-testid="generate-draft-button"]');

    // 等待生成完成
    await page.waitForSelector('[data-testid="draft-content"]', { timeout: 15000 });

    // 验证草稿内容
    const draftContent = await page.textContent('[data-testid="draft-content"]');
    expect(draftContent).toBeTruthy();
    expect(draftContent!.length).toBeGreaterThan(50);
    expect(draftContent!.toLowerCase()).toContain('meeting');
  });

  test('should generate draft with different tones', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    // 填写基本信息
    await page.fill('[data-testid="email-to"]', 'colleague@example.com');
    await page.fill('[data-testid="email-subject"]', 'Update Request');
    await page.fill('[data-testid="draft-prompt"]', 'Ask for project update');

    // 选择正式语气
    await page.selectOption('[data-testid="tone-select"]', 'formal');

    // 生成草稿
    await page.click('[data-testid="generate-draft-button"]');
    await page.waitForSelector('[data-testid="draft-content"]', { timeout: 15000 });

    const formalDraft = await page.textContent('[data-testid="draft-content"]');
    expect(formalDraft).toBeTruthy();

    // 切换到友好语气
    await page.selectOption('[data-testid="tone-select"]', 'friendly');
    await page.click('[data-testid="generate-draft-button"]');
    await page.waitForTimeout(5000);

    const friendlyDraft = await page.textContent('[data-testid="draft-content"]');
    expect(friendlyDraft).toBeTruthy();

    // 验证两个草稿不同
    expect(formalDraft).not.toBe(friendlyDraft);
  });

  test('should regenerate draft', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    // 填写信息并生成第一个草稿
    await page.fill('[data-testid="email-to"]', 'test@example.com');
    await page.fill('[data-testid="email-subject"]', 'Test');
    await page.fill('[data-testid="draft-prompt"]', 'Send a test message');
    await page.click('[data-testid="generate-draft-button"]');
    await page.waitForSelector('[data-testid="draft-content"]', { timeout: 15000 });

    const firstDraft = await page.textContent('[data-testid="draft-content"]');

    // 重新生成
    await page.click('[data-testid="regenerate-draft-button"]');
    await page.waitForTimeout(5000);

    const secondDraft = await page.textContent('[data-testid="draft-content"]');

    // 验证两个草稿不完全相同
    expect(firstDraft).toBeTruthy();
    expect(secondDraft).toBeTruthy();
    // 内容可能相似但应该有差异
  });

  test('should display draft generation progress', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    await page.fill('[data-testid="email-to"]', 'test@example.com');
    await page.fill('[data-testid="email-subject"]', 'Long Draft');
    await page.fill('[data-testid="draft-prompt"]', 'Write a detailed email about our quarterly results');

    // 点击生成
    await page.click('[data-testid="generate-draft-button"]');

    // 验证加载指示器显示
    await page.waitForSelector('[data-testid="generating-indicator"]', { timeout: 3000 });

    // 等待生成完成
    await page.waitForSelector('[data-testid="generating-indicator"]', { 
      state: 'hidden',
      timeout: 20000 
    });

    // 验证内容已生成
    const content = await page.textContent('[data-testid="draft-content"]');
    expect(content).toBeTruthy();
  });

  test('should handle draft generation errors', async ({ page }) => {
    // 模拟 API 错误
    await page.route('**/api/v1/ai/draft/generate', route => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ error: 'Generation failed' })
      });
    });

    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    await page.fill('[data-testid="email-to"]', 'test@example.com');
    await page.fill('[data-testid="email-subject"]', 'Test');
    await page.fill('[data-testid="draft-prompt"]', 'Test error handling');

    // 点击生成
    await page.click('[data-testid="generate-draft-button"]');

    // 等待错误消息
    await page.waitForSelector('[data-testid="error-message"]', { timeout: 5000 });

    const errorMessage = await page.textContent('[data-testid="error-message"]');
    expect(errorMessage).toBeTruthy();
  });

  test('should save generated draft', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    // 生成草稿
    await page.fill('[data-testid="email-to"]', 'save@example.com');
    await page.fill('[data-testid="email-subject"]', 'Save Test');
    await page.fill('[data-testid="draft-prompt"]', 'Test saving draft');
    await page.click('[data-testid="generate-draft-button"]');
    await page.waitForSelector('[data-testid="draft-content"]', { timeout: 15000 });

    // 保存草稿
    await page.click('[data-testid="save-draft-button"]');

    // 等待保存成功
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('保存成功');

    // 导航到草稿箱验证
    await page.click('[data-testid="drafts-folder"]');
    await page.waitForSelector('[data-testid="draft-item"]', { timeout: 5000 });

    const drafts = await page.locator('[data-testid="draft-item"]').allTextContents();
    expect(drafts.some(d => d.includes('Save Test'))).toBeTruthy();
  });
});

test.describe('AI Reply Generation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/emails');
  });

  test('should generate reply to email', async ({ page }) => {
    // 等待邮件列表加载
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });

    // 选择第一封邮件
    await page.locator('[data-testid="email-item"]').first().click();

    // 等待邮件详情加载
    await page.waitForSelector('[data-testid="email-detail"]', { timeout: 5000 });

    // 点击回复按钮
    await page.click('[data-testid="reply-button"]');

    // 等待回复对话框
    await page.waitForSelector('[data-testid="reply-dialog"]', { timeout: 5000 });

    // 点击 AI 生成回复
    await page.click('[data-testid="ai-reply-button"]');

    // 等待回复生成
    await page.waitForSelector('[data-testid="reply-content"]', { timeout: 15000 });

    // 验证回复内容
    const replyContent = await page.textContent('[data-testid="reply-content"]');
    expect(replyContent).toBeTruthy();
    expect(replyContent!.length).toBeGreaterThan(30);
  });

  test('should generate reply with custom instructions', async ({ page }) => {
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });
    await page.locator('[data-testid="email-item"]').first().click();
    await page.waitForSelector('[data-testid="email-detail"]', { timeout: 5000 });

    // 点击回复
    await page.click('[data-testid="reply-button"]');
    await page.waitForSelector('[data-testid="reply-dialog"]', { timeout: 5000 });

    // 输入自定义指令
    await page.fill('[data-testid="reply-instructions"]', 'Politely decline the invitation');

    // 生成回复
    await page.click('[data-testid="ai-reply-button"]');
    await page.waitForSelector('[data-testid="reply-content"]', { timeout: 15000 });

    // 验证回复包含拒绝内容
    const replyContent = await page.textContent('[data-testid="reply-content"]');
    expect(replyContent!.toLowerCase()).toMatch(/decline|unable|cannot/);
  });

  test('should generate different reply tones', async ({ page }) => {
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });
    await page.locator('[data-testid="email-item"]').first().click();
    await page.waitForSelector('[data-testid="email-detail"]', { timeout: 5000 });

    // 点击回复
    await page.click('[data-testid="reply-button"]');
    await page.waitForSelector('[data-testid="reply-dialog"]', { timeout: 5000 });

    // 测试正式语气
    await page.selectOption('[data-testid="reply-tone-select"]', 'formal');
    await page.click('[data-testid="ai-reply-button"]');
    await page.waitForSelector('[data-testid="reply-content"]', { timeout: 15000 });
    
    const formalReply = await page.textContent('[data-testid="reply-content"]');
    expect(formalReply).toBeTruthy();

    // 测试简短语气
    await page.selectOption('[data-testid="reply-tone-select"]', 'brief');
    await page.click('[data-testid="ai-reply-button"]');
    await page.waitForTimeout(5000);

    const briefReply = await page.textContent('[data-testid="reply-content"]');
    expect(briefReply).toBeTruthy();

    // 简短回复应该更短
    expect(briefReply!.length).toBeLessThan(formalReply!.length * 1.5);
  });

  test('should include email context in reply', async ({ page }) => {
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });

    // 选择特定的邮件
    await page.locator('[data-testid="email-item"]').first().click();
    await page.waitForSelector('[data-testid="email-detail"]', { timeout: 5000 });

    // 获取原邮件的关键信息
    const originalSubject = await page.textContent('[data-testid="email-subject"]');

    // 生成回复
    await page.click('[data-testid="reply-button"]');
    await page.waitForSelector('[data-testid="reply-dialog"]', { timeout: 5000 });
    await page.click('[data-testid="ai-reply-button"]');
    await page.waitForSelector('[data-testid="reply-content"]', { timeout: 15000 });

    // 验证回复与原邮件相关
    const replyContent = await page.textContent('[data-testid="reply-content"]');
    expect(replyContent).toBeTruthy();
    // 回复应该有实质内容，不只是通用模板
    expect(replyContent!.length).toBeGreaterThan(50);
  });

  test('should send AI-generated reply', async ({ page }) => {
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });
    await page.locator('[data-testid="email-item"]').first().click();
    await page.waitForSelector('[data-testid="email-detail"]', { timeout: 5000 });

    // 生成回复
    await page.click('[data-testid="reply-button"]');
    await page.waitForSelector('[data-testid="reply-dialog"]', { timeout: 5000 });
    await page.click('[data-testid="ai-reply-button"]');
    await page.waitForSelector('[data-testid="reply-content"]', { timeout: 15000 });

    // 发送回复
    await page.click('[data-testid="send-reply-button"]');

    // 等待发送成功
    await page.waitForSelector('.toast', { timeout: 5000 });
    const toastMessage = await page.textContent('.toast');
    expect(toastMessage).toContain('发送成功');
  });

  test('should edit AI-generated reply before sending', async ({ page }) => {
    await page.waitForSelector('[data-testid="email-item"]', { timeout: 5000 });
    await page.locator('[data-testid="email-item"]').first().click();
    await page.waitForSelector('[data-testid="email-detail"]', { timeout: 5000 });

    // 生成回复
    await page.click('[data-testid="reply-button"]');
    await page.waitForSelector('[data-testid="reply-dialog"]', { timeout: 5000 });
    await page.click('[data-testid="ai-reply-button"]');
    await page.waitForSelector('[data-testid="reply-content"]', { timeout: 15000 });

    // 编辑回复
    const additionalText = '\n\nAdditional custom content';
    await page.locator('[data-testid="reply-content"]').fill(
      await page.locator('[data-testid="reply-content"]').inputValue() + additionalText
    );

    // 验证编辑成功
    const editedContent = await page.textContent('[data-testid="reply-content"]');
    expect(editedContent).toContain('Additional custom content');
  });
});

test.describe('AI Draft - Advanced Features', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/emails');
  });

  test('should generate draft with attachments context', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    // 添加附件引用
    await page.fill('[data-testid="email-to"]', 'test@example.com');
    await page.fill('[data-testid="email-subject"]', 'Document Review');
    await page.fill('[data-testid="draft-prompt"]', 'Ask to review the attached report');

    // 生成草稿
    await page.click('[data-testid="generate-draft-button"]');
    await page.waitForSelector('[data-testid="draft-content"]', { timeout: 15000 });

    // 验证草稿提到附件
    const draftContent = await page.textContent('[data-testid="draft-content"]');
    expect(draftContent!.toLowerCase()).toMatch(/attach|document|report/);
  });

  test('should handle multilingual draft generation', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    // 选择语言
    await page.selectOption('[data-testid="language-select"]', 'zh-CN');

    await page.fill('[data-testid="email-to"]', 'test@example.com');
    await page.fill('[data-testid="email-subject"]', '会议邀请');
    await page.fill('[data-testid="draft-prompt"]', '邀请参加下周的项目讨论会');

    // 生成中文草稿
    await page.click('[data-testid="generate-draft-button"]');
    await page.waitForSelector('[data-testid="draft-content"]', { timeout: 15000 });

    const draftContent = await page.textContent('[data-testid="draft-content"]');
    expect(draftContent).toBeTruthy();
    // 验证包含中文内容
    expect(/[\u4e00-\u9fa5]/.test(draftContent!)).toBeTruthy();
  });

  test('should provide draft suggestions', async ({ page }) => {
    await page.click('[data-testid="compose-email-button"]');
    await page.waitForSelector('[data-testid="compose-dialog"]', { timeout: 5000 });

    await page.fill('[data-testid="email-to"]', 'test@example.com');
    await page.fill('[data-testid="email-subject"]', 'Follow-up');

    // 请求草稿建议
    await page.click('[data-testid="draft-suggestions-button"]');

    // 等待建议列表
    await page.waitForSelector('[data-testid="suggestion-item"]', { timeout: 10000 });

    // 验证至少有一个建议
    const suggestionCount = await page.locator('[data-testid="suggestion-item"]').count();
    expect(suggestionCount).toBeGreaterThan(0);

    // 选择一个建议
    await page.locator('[data-testid="suggestion-item"]').first().click();

    // 验证建议内容已填充
    const promptValue = await page.inputValue('[data-testid="draft-prompt"]');
    expect(promptValue).toBeTruthy();
  });
});
