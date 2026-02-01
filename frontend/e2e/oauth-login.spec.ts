import { test, expect } from "@playwright/test";

// OAuth tests don't use the global auth state
test.use({ storageState: { cookies: [], origins: [] } });

test.describe("OAuth Login", () => {
  test("login page shows OAuth button when provider is configured", async ({
    page,
  }) => {
    await page.goto("/");

    // Wait for login form to render
    await expect(page.locator(".login-form")).toBeVisible({ timeout: 10000 });

    // Should show OAuth button for GitHub
    const oauthButton = page.locator(".oauth-button-github");
    await expect(oauthButton).toBeVisible();
    await expect(oauthButton).toHaveText("Sign in with Github");
  });

  test("login page shows both OAuth and password form", async ({ page }) => {
    await page.goto("/");

    await expect(page.locator(".login-form")).toBeVisible({ timeout: 10000 });

    // Should have OAuth button
    await expect(page.locator(".oauth-button-github")).toBeVisible();

    // Should have the "or" divider
    await expect(page.locator(".login-divider")).toBeVisible();

    // Should have password form
    await expect(page.locator("#userId")).toBeVisible();
    await expect(page.locator("#password")).toBeVisible();
  });

  test("OAuth flow through dummygithub completes login", async ({ page }) => {
    await page.goto("/");

    // Wait for login form
    await expect(page.locator(".login-form")).toBeVisible({ timeout: 10000 });

    // Click "Sign in with Github" — this navigates through the OAuth flow:
    // 1. /auth/github (backend) → redirects to dummygithub /login/oauth/authorize
    // 2. dummygithub shows login page
    const oauthButton = page.locator(".oauth-button-github");
    await oauthButton.click();

    // Should land on the dummygithub login page
    await expect(page.locator("text=Dummy GitHub Login")).toBeVisible({
      timeout: 10000,
    });

    // Click "Sign in as dummyuser" link to complete the OAuth flow.
    // This navigates: dummygithub → callback → redirect to /
    await page.getByRole("link", { name: "Sign in as dummyuser" }).click();

    // Wait until we're back on the app (URL contains localhost:5173)
    await page.waitForURL("http://localhost:5173/**", { timeout: 15000 });

    // The OAuth callback sets a session cookie and redirects to /.
    // Reload to ensure the app picks up the session.
    await page.reload();

    // Should show dashboard with sidebar after successful OAuth login
    await expect(page.locator(".sidebar")).toBeVisible({ timeout: 15000 });
  });
});
