import { test, expect } from "@playwright/test";

// Login tests don't use the global auth state
test.use({ storageState: { cookies: [], origins: [] } });

test.describe("Login", () => {
  test("login form renders with required fields", async ({ page }) => {
    await page.goto("/");

    // Should show login form (appears after initial 401 from /api/dashboards)
    const loginForm = page.locator(".login-form");
    await expect(loginForm).toBeVisible({ timeout: 10000 });

    // Should have user ID and password inputs
    await expect(page.locator("#userId")).toBeVisible();
    await expect(page.locator("#password")).toBeVisible();

    // Should have a submit button
    const submitButton = page.locator('button[type="submit"]');
    await expect(submitButton).toBeVisible();
    await expect(submitButton).toHaveText("Log in");
  });

  test("invalid credentials show error message", async ({ page }) => {
    await page.goto("/");

    // Wait for login form to appear
    await expect(page.locator(".login-form")).toBeVisible({ timeout: 10000 });

    await page.fill("#userId", "wrong");
    await page.fill("#password", "wrong");
    await page.click('button[type="submit"]');

    // Should show error message
    const error = page.locator(".login-error");
    await expect(error).toBeVisible();
  });

  test("valid credentials redirect to dashboard", async ({ page }) => {
    await page.goto("/");

    // Wait for login form
    await expect(page.locator(".login-form")).toBeVisible();

    await page.fill("#userId", "admin");
    await page.fill("#password", "admin");
    await page.click('button[type="submit"]');

    // Login form should disappear after successful authentication
    await expect(page.locator(".login-form")).not.toBeVisible({ timeout: 10000 });

    // Reload to trigger dashboard fetch with session cookie
    await page.reload();

    // Should show dashboard with sidebar
    await expect(page.locator(".sidebar")).toBeVisible({ timeout: 15000 });
  });
});
