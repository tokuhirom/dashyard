import { test, expect } from "@playwright/test";

test.describe("Refresh Interval Selector", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.locator(".sidebar-item", { hasText: "overview" }).click();
    await expect(page.locator(".panel").first()).toBeVisible();
  });

  test("selector exists with default Off and correct options", async ({
    page,
  }) => {
    const selector = page.locator(".refresh-interval-selector");
    await expect(selector).toBeVisible();

    // Default should be "Off" (value 0)
    await expect(selector).toHaveValue("0");

    const options = selector.locator("option");
    await expect(options).toHaveCount(5);

    const values = await options.evaluateAll((els) =>
      els.map((el) => (el as HTMLOptionElement).value)
    );
    expect(values).toEqual(["0", "10000", "30000", "60000", "300000"]);

    const labels = await options.evaluateAll((els) =>
      els.map((el) => (el as HTMLOptionElement).textContent)
    );
    expect(labels).toEqual(["Off", "⟳ 10s", "⟳ 30s", "⟳ 1m", "⟳ 5m"]);
  });

  test("auto-refresh triggers periodic API calls for relative time range", async ({
    page,
  }) => {
    // Count query_range requests
    let queryCount = 0;
    page.on("request", (req) => {
      if (req.url().includes("/api/query_range")) {
        queryCount++;
      }
    });

    // Wait for initial load requests to settle
    await page.waitForTimeout(1000);
    queryCount = 0;

    // Select 10s refresh interval (shortest available)
    const selector = page.locator(".refresh-interval-selector");
    await selector.selectOption("10000");

    // Wait long enough for at least one refresh cycle
    await page.waitForTimeout(12000);

    // Should have made at least 1 new batch of API calls
    expect(queryCount).toBeGreaterThan(0);
  });

  test("selecting Off stops periodic requests", async ({ page }) => {
    const selector = page.locator(".refresh-interval-selector");

    // Enable refresh
    await selector.selectOption("10000");
    await page.waitForTimeout(2000);

    // Turn it off
    await selector.selectOption("0");

    // Count requests after turning off
    let queryCount = 0;
    page.on("request", (req) => {
      if (req.url().includes("/api/query_range")) {
        queryCount++;
      }
    });

    // Wait and verify no new requests
    await page.waitForTimeout(12000);
    expect(queryCount).toBe(0);
  });

  test("auto-refresh does not trigger for absolute time range", async ({
    page,
  }) => {
    // Set an absolute time range via URL
    const now = new Date();
    const oneHourAgo = new Date(now.getTime() - 3600 * 1000);
    const fromISO = encodeURIComponent(oneHourAgo.toISOString());
    const toISO = encodeURIComponent(now.toISOString());
    await page.goto(`/d/overview?from=${fromISO}&to=${toISO}`);
    await expect(page.locator(".panel").first()).toBeVisible();

    // Wait for initial requests to settle
    await page.waitForTimeout(1000);

    // Count requests
    let queryCount = 0;
    page.on("request", (req) => {
      if (req.url().includes("/api/query_range")) {
        queryCount++;
      }
    });

    // Enable refresh
    const selector = page.locator(".refresh-interval-selector");
    await selector.selectOption("10000");

    // Wait for a refresh cycle
    await page.waitForTimeout(12000);

    // No new requests should have been made (absolute range returns same ref)
    expect(queryCount).toBe(0);
  });
});
