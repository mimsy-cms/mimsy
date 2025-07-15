import { expect, test } from '@playwright/test';

test('login page has expected h1', async ({ page }) => {
	await page.goto('/login');
	await expect(page.locator('h1')).toBeVisible();
});
