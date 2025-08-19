import { test, expect } from '@playwright/test';

test.describe('Login Flow', () => {
    test('should allow a user to log in with valid credentials', async ({ page }) => {
        await page.goto('/login');

        await page.fill('input[name="email"]', 'admin@example.com');
        await page.fill('input[name="password"]', 'admin123');
        await page.click('button[type="submit"]');

        await expect(page).toHaveURL('/');
        await expect(page.getByRole('heading', { name: 'Collections' })).toBeVisible();
    });

    test('should show an error for invalid credentials', async ({ page }) => {
        await page.goto('/login');

        await page.fill('input[name="email"]', 'invalid@example.com');
        await page.fill('input[name="password"]', 'invalidpassword');
        await page.click('button[type="submit"]');

        await expect(page).toHaveURL('/login');
        // Check for the error message but give it some time to appear
        await expect(page.getByText("Invalid credentials")).toBeVisible({ timeout: 5000 });


    });

    test('should not allow access to protected routes when not logged in', async ({ page }) => {
        await page.goto('/');

        await expect(page).toHaveURL('/login');
    });
});