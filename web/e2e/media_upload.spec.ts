import { test, expect, type Page } from '@playwright/test';
import fs from 'fs';
import os from 'os';
import path from 'path';
import { Buffer } from 'buffer';

/**
 * Helper function to log in a user
 * @param page - Playwright page object
 * @param email - User email (defaults to test credentials)
 * @param password - User password (defaults to test credentials)
 */
async function login(page: Page, email?: string, password?: string) {
  await page.goto('/login');
  
  // Use provided credentials or fall back to environment variables or defaults
  const testEmail = email || 'admin@example.com';
  const testPassword = password || 'admin123';
  
  await page.fill('input[name="email"]', testEmail);
  await page.fill('input[name="password"]', testPassword);
  
  await page.click('button[type="submit"]');
  
  // Wait for redirect after login
  await page.waitForURL('/', { timeout: 10000 });
}

/**
 * Helper function to find and add the id to the list of media to be removed after the test
 * @param page - Playwright page object
 * @param uploadedMediaIds - Array to store uploaded media IDs
 */
async function registerUploadedMediaId(page: Page, uploadedMediaIds: string[]) {
  const uploadedLink = page.locator('a[href^="/media/"]').last();

  const href = await uploadedLink.getAttribute('href');
  if (href) {
    const id = href.split('/').pop();
    if (id) {
      uploadedMediaIds.push(id);
    }
  }
}

// Basic tests for media page layout and such
test.describe('Media page', () => {
  // Default selection should be grid
  test('should have Grid layout selected by default on media page', async ({ page }) => {
    await login(page);
    await page.goto('/media');

    // Find buttons
    const gridButton = page.locator('button').filter({ hasText: 'Grid' });
    const listButton = page.locator('button').filter({ hasText: 'List' });

    await expect(gridButton).toBeVisible();
    await expect(listButton).toBeVisible();

    // Check that the Grid button has the active styles (bg-blue-700 text-white)
    await expect(gridButton).toHaveClass(/bg-blue-700/);
    await expect(gridButton).toHaveClass(/text-white/);
  });

  // List layout should be selected when List button is clicked
  test('should switch to List layout when List button is clicked', async ({ page }) => {
    await login(page);
    await page.goto('/media');

    const gridButton = page.locator('button').filter({ hasText: 'Grid' });
    const listButton = page.locator('button').filter({ hasText: 'List' });

    await expect(gridButton).toBeVisible();
    await expect(listButton).toBeVisible();

    await listButton.click();

    // Check that the List button has the active styles (bg-blue-700 text-white)
    await expect(listButton).toHaveClass(/bg-blue-700/);
    await expect(listButton).toHaveClass(/text-white/);

    // Verify that the list table is actually displayed
    const listTable = page.locator('table');
    await expect(listTable).toBeVisible();
  });
});

// More advanced tests for actual media upload functionality
test.describe('Media upload flow', () => {
  // Array to store uploaded media IDs for cleanup after tests
  let uploadedMediaIds: string[] = [];

  test.afterEach(async ({ request }) => {
    for (const id of uploadedMediaIds) {
      const res = await request.delete(`/api/v1/media/${id}`);
      if (!res.ok()) {
        console.warn(`Failed to delete media with id ${id}: ${res.status()} ${res.statusText()}`);
      }
    }
    uploadedMediaIds = [];
  });


  // Upload image test
  test('should upload images', async ({ page }) => {
    await login(page);
    await page.goto('/media');

    const initialMediaCount = await page.locator('.grid img, table img').count();

    const testImagePath = 'e2e/fixtures/test-image.jpg';

    const uploadButton = page.locator('button').filter({ hasText: 'Upload' });
    await expect(uploadButton).toBeVisible();

    const fileChooserPromise = page.waitForEvent('filechooser');
    await uploadButton.click();
    const fileChooser = await fileChooserPromise;
    await fileChooser.setFiles(testImagePath);

    await page.waitForTimeout(3000);
    await page.waitForLoadState('networkidle');

    const finalMediaCount = await page.locator('.grid img, table img').count();
    expect(finalMediaCount).toBeGreaterThan(initialMediaCount);

    await registerUploadedMediaId(page, uploadedMediaIds);

    const uploadedFileSelectors = [
        `img[alt*="test-image"]`,
        `img[alt*="test-image.jpg"]`,
        `img[src*="test-image"]`,
        `a[href*="test-image"]`
    ];

    let fileFound = false;
    for (const selector of uploadedFileSelectors) {
        const elements = page.locator(selector);
        if (await elements.count() > 0) {
            fileFound = true;
            break;
        }
    }

    expect(fileFound).toBe(true);
  });

  // Upload PDF test
  test('should upload pdf', async ({ page }) => {
    await login(page);
    await page.goto('/media');

    const initialMediaCount = await page.locator('.grid img, table img').count();

    const testPdfPath = 'e2e/fixtures/test-document.pdf';

    const uploadButton = page.locator('button').filter({ hasText: 'Upload' });
    await expect(uploadButton).toBeVisible();

    const fileChooserPromise = page.waitForEvent('filechooser');
    await uploadButton.click();

    const fileChooser = await fileChooserPromise;

    await fileChooser.setFiles(testPdfPath);

    await page.waitForTimeout(3000);
    await page.waitForLoadState('networkidle');

    const finalMediaCount = await page.locator('.grid img, table img').count();
    expect(finalMediaCount).toBeGreaterThan(initialMediaCount);

    await registerUploadedMediaId(page, uploadedMediaIds);

    const uploadedFileSelectors = [
        `img[alt*="test-document"]`,
        `img[alt*="test-document.pdf"]`,
        `img[src*="test-document"]`,
        `a[href*="test-document"]`
    ];

    let fileFound = false;
    for (const selector of uploadedFileSelectors) {
        const elements = page.locator(selector);
        if (await elements.count() > 0) {
            fileFound = true;
            break;
        }
    }

    expect(fileFound).toBe(true);
  });

  // Upload oversized file test
  test('should show an error when uploading oversized file', async ({ page }) => {
    const tempDir = os.tmpdir();
    const filePath = path.join(tempDir, 'oversized-file.jpg');

    const fileSize = 256 * 1024 * 1024 + 1;

    await new Promise<void>((resolve, reject) => {
      const fd = fs.openSync(filePath, 'w');
      try {
        fs.writeSync(fd, Buffer.from([0xff, 0xd8, 0xff]), 0, 3, 0);
        fs.writeSync(fd, Buffer.from([0]), 0, 1, fileSize - 1);
        fs.closeSync(fd);
        resolve();
      } catch (err) {
        reject(err);
      }
    });

    // Listen for the alert dialog and verify the message
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('The following files are too large to upload');
      await dialog.dismiss();
    });

    await login(page);
    await page.goto('/media');

    const uploadButton = page.locator('button', { hasText: 'Upload' });
    await expect(uploadButton).toBeVisible();

    const fileChooserPromise = page.waitForEvent('filechooser');
    await uploadButton.click();
    const fileChooser = await fileChooserPromise;

    await fileChooser.setFiles(filePath);
    
    fs.unlinkSync(filePath);
  });
});