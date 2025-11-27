/**
 * E2E Test: AI Streaming Chat
 * 测试 SSE 流式聊天功能
 */

import { test, expect } from '@playwright/test';

test.describe('AI Streaming Chat', () => {
  test.beforeEach(async ({ page }) => {
    // 登录并导航到聊天页面
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/copilot');
  });

  test('should stream AI chat responses', async ({ page }) => {
    // 等待聊天界面加载
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 输入消息
    const testMessage = 'Hello, can you help me with my emails?';
    await page.fill('[data-testid="chat-input"]', testMessage);

    // 发送消息
    await page.click('[data-testid="send-button"]');

    // 验证用户消息已显示
    await page.waitForSelector('[data-testid="user-message"]', { timeout: 3000 });
    const userMessage = await page.locator('[data-testid="user-message"]').last().textContent();
    expect(userMessage).toContain(testMessage);

    // 等待 AI 响应开始流式传输
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });

    // 验证流式响应
    const aiMessage = page.locator('[data-testid="ai-message"]').last();
    
    // 等待流式响应完成（文本长度停止增长）
    let previousLength = 0;
    let stableCount = 0;
    const maxWaitTime = 30000; // 最多等待30秒
    const startTime = Date.now();

    while (Date.now() - startTime < maxWaitTime) {
      await page.waitForTimeout(500);
      const currentText = await aiMessage.textContent();
      const currentLength = currentText?.length || 0;

      if (currentLength > 0) {
        if (currentLength === previousLength) {
          stableCount++;
          if (stableCount >= 3) {
            // 长度连续3次不变，认为流式传输完成
            break;
          }
        } else {
          stableCount = 0;
        }
        previousLength = currentLength;
      }
    }

    // 验证响应内容
    const finalText = await aiMessage.textContent();
    expect(finalText).toBeTruthy();
    expect(finalText!.length).toBeGreaterThan(10);
  });

  test('should display streaming indicator during response', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 发送消息
    await page.fill('[data-testid="chat-input"]', 'Tell me about email management');
    await page.click('[data-testid="send-button"]');

    // 验证流式传输指示器显示
    await page.waitForSelector('[data-testid="streaming-indicator"]', { 
      timeout: 5000 
    });

    // 等待流式传输完成，指示器应该消失
    await page.waitForSelector('[data-testid="streaming-indicator"]', { 
      state: 'hidden',
      timeout: 30000 
    });
  });

  test('should handle multiple streaming messages', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 发送第一条消息
    await page.fill('[data-testid="chat-input"]', 'What is your name?');
    await page.click('[data-testid="send-button"]');

    // 等待第一条响应完成
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });
    await page.waitForTimeout(2000);

    // 发送第二条消息
    await page.fill('[data-testid="chat-input"]', 'What can you do?');
    await page.click('[data-testid="send-button"]');

    // 等待第二条响应完成
    await page.waitForTimeout(5000);

    // 验证有多条消息
    const messageCount = await page.locator('[data-testid="ai-message"]').count();
    expect(messageCount).toBeGreaterThanOrEqual(2);
  });

  test('should handle streaming errors gracefully', async ({ page }) => {
    // 模拟网络错误
    await page.route('**/api/v1/ai/chat/stream', route => {
      route.abort('failed');
    });

    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 发送消息
    await page.fill('[data-testid="chat-input"]', 'Test error handling');
    await page.click('[data-testid="send-button"]');

    // 等待错误消息显示
    await page.waitForSelector('[data-testid="error-message"]', { timeout: 5000 });

    // 验证错误消息
    const errorMessage = await page.textContent('[data-testid="error-message"]');
    expect(errorMessage).toBeTruthy();
  });

  test('should allow canceling streaming response', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 发送消息
    await page.fill('[data-testid="chat-input"]', 'Tell me a very long story');
    await page.click('[data-testid="send-button"]');

    // 等待流式传输开始
    await page.waitForSelector('[data-testid="streaming-indicator"]', { timeout: 5000 });

    // 点击取消按钮
    await page.click('[data-testid="cancel-streaming-button"]');

    // 验证流式传输已停止
    await page.waitForSelector('[data-testid="streaming-indicator"]', { 
      state: 'hidden',
      timeout: 3000 
    });

    // 验证消息显示为已取消
    const statusIndicator = await page.locator('[data-testid="message-status"]').last().textContent();
    expect(statusIndicator).toContain('已取消');
  });

  test('should preserve chat history across streaming', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 发送多条消息
    const messages = ['First message', 'Second message', 'Third message'];
    
    for (const msg of messages) {
      await page.fill('[data-testid="chat-input"]', msg);
      await page.click('[data-testid="send-button"]');
      await page.waitForTimeout(3000); // 等待响应
    }

    // 验证所有用户消息都在历史记录中
    const userMessages = await page.locator('[data-testid="user-message"]').allTextContents();
    expect(userMessages.length).toBeGreaterThanOrEqual(3);
    
    for (const msg of messages) {
      expect(userMessages.some(m => m.includes(msg))).toBeTruthy();
    }
  });

  test('should handle SSE reconnection', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 发送第一条消息
    await page.fill('[data-testid="chat-input"]', 'First test');
    await page.click('[data-testid="send-button"]');
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });

    // 模拟网络中断后恢复
    await page.route('**/api/v1/ai/chat/stream', route => route.continue());

    // 发送第二条消息
    await page.fill('[data-testid="chat-input"]', 'Second test after reconnection');
    await page.click('[data-testid="send-button"]');

    // 验证连接恢复后消息仍能正常发送
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });
    const messageCount = await page.locator('[data-testid="ai-message"]').count();
    expect(messageCount).toBeGreaterThanOrEqual(2);
  });
});

test.describe('AI Streaming - Context Awareness', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/copilot');
  });

  test('should maintain conversation context in streaming', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 第一条消息：建立上下文
    await page.fill('[data-testid="chat-input"]', 'My name is John');
    await page.click('[data-testid="send-button"]');
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });
    await page.waitForTimeout(2000);

    // 第二条消息：测试上下文是否保留
    await page.fill('[data-testid="chat-input"]', 'What is my name?');
    await page.click('[data-testid="send-button"]');
    await page.waitForTimeout(5000);

    // 验证响应包含上下文信息
    const lastResponse = await page.locator('[data-testid="ai-message"]').last().textContent();
    expect(lastResponse?.toLowerCase()).toContain('john');
  });

  test('should include system context in streaming responses', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 询问系统相关问题
    await page.fill('[data-testid="chat-input"]', 'What features does this email system have?');
    await page.click('[data-testid="send-button"]');

    // 等待响应
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });
    await page.waitForTimeout(5000);

    // 验证响应包含系统上下文
    const response = await page.locator('[data-testid="ai-message"]').last().textContent();
    expect(response).toBeTruthy();
    expect(response!.length).toBeGreaterThan(50);
  });
});

test.describe('AI Streaming - Performance', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="email"]', 'test@example.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    await page.goto('/copilot');
  });

  test('should stream responses with acceptable latency', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    const startTime = Date.now();

    // 发送消息
    await page.fill('[data-testid="chat-input"]', 'Quick response test');
    await page.click('[data-testid="send-button"]');

    // 等待第一个 token 出现
    await page.waitForSelector('[data-testid="ai-message"]', { timeout: 10000 });

    const firstTokenTime = Date.now() - startTime;

    // 验证首字节时间在合理范围内（应该小于 5 秒）
    expect(firstTokenTime).toBeLessThan(5000);
  });

  test('should handle rapid consecutive messages', async ({ page }) => {
    await page.waitForSelector('[data-testid="chat-input"]', { timeout: 5000 });

    // 快速发送多条消息
    for (let i = 0; i < 3; i++) {
      await page.fill('[data-testid="chat-input"]', `Rapid message ${i + 1}`);
      await page.click('[data-testid="send-button"]');
      await page.waitForTimeout(100); // 短暂延迟
    }

    // 等待所有响应完成
    await page.waitForTimeout(15000);

    // 验证所有消息都已处理
    const messageCount = await page.locator('[data-testid="ai-message"]').count();
    expect(messageCount).toBeGreaterThanOrEqual(3);
  });
});
