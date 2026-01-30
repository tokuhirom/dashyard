import { test, expect } from "@playwright/test";

test.describe("Column Selector", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    // Navigate to a dashboard with panels
    await page.locator(".sidebar-item", { hasText: "overview" }).click();
    // Wait for panels to render
    await expect(page.locator(".panel").first()).toBeVisible();
  });

  test("selector exists with default value and options 2-6", async ({
    page,
  }) => {
    const selector = page.locator(".column-selector");
    await expect(selector).toBeVisible();

    // Default value should be "2"
    await expect(selector).toHaveValue("2");

    // Should have options 2 through 6
    const options = selector.locator("option");
    await expect(options).toHaveCount(5);

    const values = await options.evaluateAll((els) =>
      els.map((el) => (el as HTMLOptionElement).value)
    );
    expect(values).toEqual(["2", "3", "4", "5", "6"]);
  });

  test("changing selector updates gridTemplateColumns style", async ({
    page,
  }) => {
    const selector = page.locator(".column-selector");
    const rowPanels = page.locator(".row-panels").first();

    // Default should be repeat(2, 1fr)
    const defaultStyle = await rowPanels.evaluate(
      (el) => el.style.gridTemplateColumns
    );
    expect(defaultStyle).toContain("repeat(2, 1fr)");

    // Change to 4 columns
    await selector.selectOption("4");

    const newStyle = await rowPanels.evaluate(
      (el) => el.style.gridTemplateColumns
    );
    expect(newStyle).toContain("repeat(4, 1fr)");
  });

  test("panel widths decrease when column count increases", async ({
    page,
  }) => {
    const selector = page.locator(".column-selector");
    const firstPanel = page.locator(".panel").first();

    // Measure width at 2 columns
    await selector.selectOption("2");
    await page.waitForTimeout(200);
    const width2 = await firstPanel.evaluate(
      (el) => el.getBoundingClientRect().width
    );

    // Measure width at 4 columns
    await selector.selectOption("4");
    await page.waitForTimeout(200);
    const width4 = await firstPanel.evaluate(
      (el) => el.getBoundingClientRect().width
    );

    // Panels should be narrower with more columns
    expect(width4).toBeLessThan(width2);
  });

  test("canvas resizes to match container after column change", async ({
    page,
  }) => {
    const selector = page.locator(".column-selector");

    // Wait for chart canvas to be present
    const graphPanel = page.locator(".graph-panel").first();
    await expect(graphPanel).toBeVisible();
    const canvas = graphPanel.locator("canvas");
    await expect(canvas).toBeVisible();

    // Change to different column counts and verify canvas fits
    for (const cols of ["3", "4", "5"]) {
      await selector.selectOption(cols);
      // Allow time for Chart.js resize
      await page.waitForTimeout(500);

      const { containerWidth, canvasWidth } = await graphPanel.evaluate(
        (panel) => {
          const chartContainer = panel.querySelector(".panel-chart");
          const cvs = panel.querySelector("canvas");
          return {
            containerWidth: chartContainer
              ? chartContainer.getBoundingClientRect().width
              : 0,
            canvasWidth: cvs ? cvs.getBoundingClientRect().width : 0,
          };
        }
      );

      // Canvas width should be close to container width (< 5px difference)
      expect(Math.abs(containerWidth - canvasWidth)).toBeLessThan(5);
    }
  });

  test("no canvas overflow after cycling through all column values", async ({
    page,
  }) => {
    const selector = page.locator(".column-selector");

    // Cycle through all column values
    for (const cols of ["2", "3", "4", "5", "6", "2"]) {
      await selector.selectOption(cols);
      await page.waitForTimeout(300);
    }

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
