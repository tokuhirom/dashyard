import { test, expect } from "@playwright/test";

test.describe("Template Variables", () => {
  test("variable dropdown renders and selects values on variable dashboard", async ({
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

    // Should have a select dropdown
    const select = variableBar.locator(".variable-select");
    await expect(select).toBeVisible();

    // Should have a label
    const label = variableBar.locator(".variable-label");
    await expect(label).toHaveText("Network Device");

    // Dropdown should have options loaded from dummyprom
    const options = select.locator("option");
    const count = await options.count();
    expect(count).toBeGreaterThanOrEqual(2);

    // First option should be auto-selected (eth0)
    const selectedValue = await select.inputValue();
    expect(selectedValue).toBe("eth0");

    // Row title should contain the substituted variable value
    const rowTitle = page.locator(".row-title").first();
    await expect(rowTitle).toContainText("eth0");

    // Graph panels should be visible
    const panels = page.locator(".panel");
    await expect(panels.first()).toBeVisible();

    // Change dropdown to second value
    await select.selectOption({ index: 1 });
    const newValue = await select.inputValue();
    expect(newValue).toBe("eth1");

    // Row title should update with new variable value
    await expect(rowTitle).toContainText("eth1");
  });

  test("panel titles reflect variable substitution", async ({ page }) => {
    await page.goto("/");

    await page
      .locator(".sidebar-item", { hasText: "network-variable" })
      .click();

    // Wait for variable bar
    await expect(page.locator(".variable-bar")).toBeVisible({ timeout: 10000 });

    // Panel titles should contain the substituted variable (first selected value)
    const panelTitles = page.locator(".panel-title");
    await expect(panelTitles.first()).toBeVisible();

    const firstTitle = await panelTitles.first().textContent();
    // The title should contain the selected variable value, not the raw $device
    expect(firstTitle).not.toContain("$device");
    expect(firstTitle).toContain("eth0");
  });
});

test.describe("Repeat Rows", () => {
  test("repeat row creates one row per variable value", async ({ page }) => {
    await page.goto("/");

    // Navigate to the network-repeat dashboard
    await page
      .locator(".sidebar-item", { hasText: "network-repeat" })
      .click();

    // Variable bar is hidden for repeat-only variables, wait for rows instead
    await expect(page.locator(".row-title").first()).toBeVisible({ timeout: 10000 });

    // Should have multiple rows (one per device value from dummyprom: eth0, eth1)
    const rowTitles = page.locator(".row-title");
    await expect(rowTitles.first()).toBeVisible();

    const rowCount = await rowTitles.count();
    expect(rowCount).toBeGreaterThanOrEqual(2);

    // Each row should have a different device in its title
    const titles: string[] = [];
    for (let i = 0; i < rowCount; i++) {
      const text = await rowTitles.nth(i).textContent();
      titles.push(text || "");
    }

    // Verify rows contain different device names
    const hasEth0 = titles.some((t) => t.includes("eth0"));
    const hasEth1 = titles.some((t) => t.includes("eth1"));
    expect(hasEth0).toBe(true);
    expect(hasEth1).toBe(true);
  });

  test("repeat rows each have graph panels", async ({ page }) => {
    await page.goto("/");

    await page
      .locator(".sidebar-item", { hasText: "network-repeat" })
      .click();

    // Variable bar is hidden for repeat-only variables, wait for panels instead
    await expect(page.locator(".panel").first()).toBeVisible({ timeout: 10000 });

    // Wait for panels to render
    await expect(page.locator(".panel").first()).toBeVisible();

    // Each repeated row should have panels
    const rows = page.locator(".row");
    const rowCount = await rows.count();
    expect(rowCount).toBeGreaterThanOrEqual(2);

    // Each row should contain graph panels with canvas elements
    for (let i = 0; i < rowCount; i++) {
      const panels = rows.nth(i).locator(".panel");
      const panelCount = await panels.count();
      expect(panelCount).toBeGreaterThanOrEqual(1);
    }
  });

  test("repeat row panel titles do not contain raw variable syntax", async ({
    page,
  }) => {
    await page.goto("/");

    await page
      .locator(".sidebar-item", { hasText: "network-repeat" })
      .click();

    // Variable bar is hidden for repeat-only variables, wait for panel titles instead
    await expect(page.locator(".panel-title").first()).toBeVisible({ timeout: 10000 });

    // No panel title should contain raw $device
    const panelTitles = page.locator(".panel-title");
    const count = await panelTitles.count();
    for (let i = 0; i < count; i++) {
      const text = await panelTitles.nth(i).textContent();
      expect(text).not.toContain("$device");
    }
  });
});
