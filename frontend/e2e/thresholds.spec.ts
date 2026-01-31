import { test, expect } from "@playwright/test";

test.describe("Thresholds", () => {
  test("thresholds dashboard loads and shows graph panels", async ({
    page,
  }) => {
    await page.goto("/");

    // Navigate to the thresholds dashboard
    await page.locator(".sidebar-item", { hasText: "thresholds" }).click();

    // Wait for graph panels to render
    const graphPanel = page.locator(".graph-panel").first();
    await expect(graphPanel).toBeVisible({ timeout: 15000 });

    // Graph panels should contain canvas elements (Chart.js)
    const canvas = graphPanel.locator("canvas");
    await expect(canvas).toBeVisible();
  });

  test("thresholds dashboard has multiple rows", async ({ page }) => {
    await page.goto("/");

    await page.locator(".sidebar-item", { hasText: "thresholds" }).click();

    // Wait for panels to render
    await expect(page.locator(".graph-panel canvas").first()).toBeVisible({
      timeout: 15000,
    });

    // Should have 3 rows (CPU, Memory, Area Chart)
    const rowTitles = page.locator(".row-title");
    await expect(rowTitles.first()).toBeVisible();
    const count = await rowTitles.count();
    expect(count).toBe(3);
  });

  test("thresholds dashboard renders all graph panels with canvas", async ({
    page,
  }) => {
    await page.goto("/");

    await page.locator(".sidebar-item", { hasText: "thresholds" }).click();

    // Wait for panels to render
    await expect(page.locator(".graph-panel canvas").first()).toBeVisible({
      timeout: 15000,
    });

    // Each graph panel should have a canvas element
    const graphPanels = page.locator(".graph-panel");
    const panelCount = await graphPanels.count();
    expect(panelCount).toBe(3);

    for (let i = 0; i < panelCount; i++) {
      const canvas = graphPanels.nth(i).locator("canvas");
      await expect(canvas).toBeVisible();
    }
  });
});
