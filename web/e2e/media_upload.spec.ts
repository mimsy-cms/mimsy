import { test, expect, type Page } from '@playwright/test';

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

test.describe('Media upload flow', () => {
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

    // Verify that the grid layout is actually displayed
    // The grid container should be visible
    const gridContainer = page.locator('.grid.grid-cols-1.gap-4');
    await expect(gridContainer).toBeVisible();
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

    if (!fileFound) {
        console.warn('Could not find uploaded file by name but media count increased');
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

    if (!fileFound) {
        console.warn('Could not find uploaded file by name but media count increased');
    }

    expect(fileFound).toBe(true);
  });
});