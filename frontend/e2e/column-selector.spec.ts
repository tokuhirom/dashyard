import { test, expect } from "@playwright/test";

test.describe("Auto Panel Layout", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    // Navigate to a dashboard with panels
    await page.locator(".sidebar-item", { hasText: "overview" }).click();
    // Wait for panels to render
    await expect(page.locator(".panel").first()).toBeVisible();
  });

  test("column selector is removed from header", async ({ page }) => {
    const selector = page.locator(".column-selector");
    await expect(selector).toHaveCount(0);
  });

  test("row panels use grid layout with auto-calculated columns", async ({
    page,
  }) => {
    const rowPanels = page.locator(".row-panels").first();
    const style = await rowPanels.evaluate(
      (el) => el.style.gridTemplateColumns
    );
    // Should have a repeat(N, 1fr) pattern based on panel count
    expect(style).toMatch(/repeat\(\d+, 1fr\)/);
  });

  test("panels render at correct widths within their row", async ({
    page,
  }) => {
    // Each row should have panels that fill the available width
    const firstPanel = page.locator(".panel").first();
    const panelWidth = await firstPanel.evaluate(
      (el) => el.getBoundingClientRect().width
    );
    // Panel should have a reasonable width (> min-width of 200px)
    expect(panelWidth).toBeGreaterThan(100);
  });

  test("no canvas overflow in graph panels", async ({ page }) => {
    // Wait for chart canvas to be present
    const graphPanel = page.locator(".graph-panel").first();
    await expect(graphPanel).toBeVisible();
    const canvas = graphPanel.locator("canvas");
    await expect(canvas).toBeVisible();

    // Allow time for Chart.js render
    await page.waitForTimeout(500);

    // Check that no canvas overflows its container
    const overflows = await page.evaluate(() => {
      const panels = document.querySelectorAll(".graph-panel");
      const results: boolean[] = [];
      panels.forEach((panel) => {
        const container = panel.querySelector(".panel-chart");
        const cvs = panel.querySelector("canvas");
        if (container && cvs) {
          const containerRect = container.getBoundingClientRect();
          const canvasRect = cvs.getBoundingClientRect();
          results.push(canvasRect.width > containerRect.width + 1);
        }
      });
      return results;
    });

    // No panel should have canvas overflow
    overflows.forEach((overflow) => {
      expect(overflow).toBe(false);
    });
  });
});
