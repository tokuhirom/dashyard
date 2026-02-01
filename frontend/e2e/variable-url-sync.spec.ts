import { test, expect } from "@playwright/test";

test.describe("Variable URL Sync", () => {
  test("changing a variable updates the URL with var- parameter", async ({
    page,
  }) => {
    await page.goto("/");

    // Navigate to the network-variable dashboard
    await page
      .locator(".sidebar-item", { hasText: "network-variable" })
      .click();

    // Wait for the variable bar to appear
    const variableBar = page.locator(".variable-bar");
    await expect(variableBar).toBeVisible({ timeout: 10000 });

    const select = variableBar.locator(".variable-select");

    // Default should be eth0 - URL should not have var-device yet (or var-device=eth0)
    const selectedValue = await select.inputValue();
    expect(selectedValue).toBe("eth0");

    // Change dropdown to eth1
    await select.selectOption({ index: 1 });
    const newValue = await select.inputValue();
    expect(newValue).toBe("eth1");

    // URL should now contain var-device=eth1
    const url = new URL(page.url());
    expect(url.searchParams.get("var-device")).toBe("eth1");
  });

  test("variable value is restored from URL on page load", async ({
    page,
  }) => {
    // Navigate directly to dashboard with var-device=eth1 in URL
    await page.goto("/d/network-variable?var-device=eth1");

    // Wait for the variable bar to appear
    const variableBar = page.locator(".variable-bar");
    await expect(variableBar).toBeVisible({ timeout: 10000 });

    const select = variableBar.locator(".variable-select");

    // The dropdown should show eth1 (restored from URL)
    await expect(select).toHaveValue("eth1");

    // Row title should reflect the restored value
    const rowTitle = page.locator(".row-title").first();
    await expect(rowTitle).toContainText("eth1");
  });

  test("navigating to another dashboard clears variable params from URL", async ({
    page,
  }) => {
    // Start on network-variable with a var param
    await page.goto("/d/network-variable?var-device=eth1");

    const variableBar = page.locator(".variable-bar");
    await expect(variableBar).toBeVisible({ timeout: 10000 });

    // Navigate to a different dashboard
    await page.locator(".sidebar-item", { hasText: "overview" }).click();

    // Wait for the new dashboard to load
    await expect(page.locator(".row-title").first()).toBeVisible({
      timeout: 10000,
    });

    // URL should not contain var-device anymore
    const url = new URL(page.url());
    expect(url.searchParams.has("var-device")).toBe(false);
    expect(url.pathname).toBe("/d/overview");
  });

  test("browser back restores previous variable selection", async ({
    page,
  }) => {
    await page.goto("/");

    // Navigate to network-variable dashboard
    await page
      .locator(".sidebar-item", { hasText: "network-variable" })
      .click();

    const variableBar = page.locator(".variable-bar");
    await expect(variableBar).toBeVisible({ timeout: 10000 });

    const select = variableBar.locator(".variable-select");
    expect(await select.inputValue()).toBe("eth0");

    // Navigate to another dashboard (pushes history entry)
    await page.locator(".sidebar-item", { hasText: "overview" }).click();
    await expect(page.locator(".row-title").first()).toBeVisible({
      timeout: 10000,
    });

    // Go back
    await page.goBack();

    // Should be back on network-variable
    await expect(variableBar).toBeVisible({ timeout: 10000 });
    expect(page.url()).toContain("/d/network-variable");
  });

  test("time range and variable params coexist in URL", async ({ page }) => {
    // Load with both time range and variable
    await page.goto("/d/network-variable?t=6h&var-device=eth1");

    const variableBar = page.locator(".variable-bar");
    await expect(variableBar).toBeVisible({ timeout: 10000 });

    const select = variableBar.locator(".variable-select");

    // Variable should be restored
    await expect(select).toHaveValue("eth1");

    // URL should still have both params
    const url = new URL(page.url());
    expect(url.searchParams.get("t")).toBe("6h");
    expect(url.searchParams.get("var-device")).toBe("eth1");
  });
});
