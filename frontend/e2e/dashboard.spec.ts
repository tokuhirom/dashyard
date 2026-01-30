import { test, expect } from "@playwright/test";

test.describe("Dashboard", () => {
  test("loads and shows panels", async ({ page }) => {
    await page.goto("/");

    // Sidebar should be visible with dashboard links
    const sidebar = page.locator(".sidebar");
    await expect(sidebar).toBeVisible();

    // Should have at least one sidebar item
    const items = sidebar.locator(".sidebar-item");
    await expect(items.first()).toBeVisible();

    // Click the first dashboard link
    await items.first().click();

    // Should display row titles and panels
    await expect(page.locator(".row-title").first()).toBeVisible();
    await expect(page.locator(".panel").first()).toBeVisible();
  });

  test("shows graph panels with canvas elements", async ({ page }) => {
    await page.goto("/");

    // Navigate to the overview dashboard
    await page.locator(".sidebar-item", { hasText: "overview" }).click();

    // Wait for graph panels to render
    const graphPanel = page.locator(".graph-panel").first();
    await expect(graphPanel).toBeVisible();

    // Graph panels should contain a canvas element (Chart.js)
    const canvas = graphPanel.locator("canvas");
    await expect(canvas).toBeVisible();
  });

  test("shows markdown panels with rendered content", async ({ page }) => {
    await page.goto("/");

    await page.locator(".sidebar-item", { hasText: "overview" }).click();

    const markdownPanel = page.locator(".markdown-panel").first();
    await expect(markdownPanel).toBeVisible();

    // Markdown panel should contain rendered HTML content
    await expect(markdownPanel.locator(".panel-content")).toBeVisible();
  });

  test("sidebar navigation switches dashboards", async ({ page }) => {
    await page.goto("/");

    const sidebarItems = page.locator(".sidebar-item");
    const count = await sidebarItems.count();

    if (count >= 2) {
      // Click the first dashboard
      await sidebarItems.nth(0).click();
      const firstTitle = await page.locator("h2").first().textContent();

      // Click a different dashboard
      await sidebarItems.nth(1).click();
      await page.waitForTimeout(500);
      const secondTitle = await page.locator("h2").first().textContent();

      // Dashboard content should have changed
      expect(firstTitle).not.toEqual(secondTitle);
    }
  });

  test("sidebar groups can be expanded and collapsed", async ({ page }) => {
    await page.goto("/");

    const groupHeader = page.locator(".sidebar-group-header").first();
    const hasGroups = (await groupHeader.count()) > 0;

    if (hasGroups) {
      // Click to toggle the group
      await groupHeader.click();

      // The arrow should indicate expanded/collapsed state
      const arrow = groupHeader.locator(".sidebar-arrow");
      await expect(arrow).toBeVisible();
    }
  });
});
