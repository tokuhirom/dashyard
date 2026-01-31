// take-screenshots.ts -- Capture README screenshots using Playwright.
// Prerequisites: start dev services first (make dev-dummyprom, make dev-backend, make dev-frontend).
// Usage: cd frontend && npx tsx take-screenshots.ts
import { chromium } from "@playwright/test";
import path from "path";
import fs from "fs";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const ROOT_DIR = path.resolve(__dirname, "..");

async function main() {
  fs.mkdirSync(path.join(ROOT_DIR, "docs"), { recursive: true });

  const browser = await chromium.launch();
  const context = await browser.newContext({
    viewport: { width: 1280, height: 800 },
  });
  const page = await context.newPage();

  // Go to login page
  await page.goto("http://localhost:5173/");

  // Wait for login form
  await page.waitForSelector(".login-form", { timeout: 15000 });

  // Screenshot: Login page
  await page.screenshot({
    path: path.join(ROOT_DIR, "docs", "screenshot-login.png"),
    fullPage: false,
  });
  console.log("Saved: screenshot-login.png");

  // Login
  await page.fill("#userId", "admin");
  await page.fill("#password", "admin");
  await page.click('button[type="submit"]');
  await page.waitForResponse(
    (resp) =>
      resp.url().includes("/api/login") &&
      resp.request().method() === "POST" &&
      resp.ok(),
    { timeout: 10000 }
  );
  await page.reload();
  await page.waitForSelector(".sidebar", { timeout: 15000 });

  // Navigate to overview dashboard
  await page.locator(".sidebar-item", { hasText: "overview" }).click();
  // Wait for graphs to render
  await page.waitForSelector(".graph-panel canvas", { timeout: 15000 });
  // Give Chart.js a moment to finish animations
  await page.waitForTimeout(2000);

  // Screenshot: Main dashboard view (overview)
  await page.screenshot({
    path: path.join(ROOT_DIR, "screenshot.png"),
    fullPage: false,
  });
  console.log("Saved: screenshot.png (main)");

  // Screenshot: Full page dashboard
  await page.screenshot({
    path: path.join(ROOT_DIR, "docs", "screenshot-dashboard.png"),
    fullPage: true,
  });
  console.log("Saved: screenshot-dashboard.png");

  // Navigate to a dashboard with variables (network-variable)
  const varItem = page.locator(".sidebar-item", {
    hasText: "network-variable",
  });
  if ((await varItem.count()) > 0) {
    await varItem.click();
    await page.waitForSelector(".graph-panel canvas", { timeout: 15000 });
    await page.waitForTimeout(2000);
    await page.screenshot({
      path: path.join(ROOT_DIR, "docs", "screenshot-variables.png"),
      fullPage: false,
    });
    console.log("Saved: screenshot-variables.png");
  }

  // Navigate to a dashboard with repeat rows (network-repeat)
  const repeatItem = page.locator(".sidebar-item", {
    hasText: "network-repeat",
  });
  if ((await repeatItem.count()) > 0) {
    await repeatItem.click();
    await page.waitForSelector(".graph-panel canvas", { timeout: 15000 });
    await page.waitForTimeout(2000);
    await page.screenshot({
      path: path.join(ROOT_DIR, "docs", "screenshot-repeat.png"),
      fullPage: true,
    });
    console.log("Saved: screenshot-repeat.png");
  }

  // Navigate to chart-types dashboard
  const chartItem = page.locator(".sidebar-item", {
    hasText: "chart-types",
  });
  if ((await chartItem.count()) > 0) {
    await chartItem.click();
    await page.waitForSelector(".panel", { timeout: 15000 });
    await page.waitForTimeout(2000);
    await page.screenshot({
      path: path.join(ROOT_DIR, "docs", "screenshot-chart-types.png"),
      fullPage: true,
    });
    console.log("Saved: screenshot-chart-types.png");
  }

  // Navigate to sidebar groups (infra/)
  const groupHeader = page.locator(".sidebar-group-header").first();
  if ((await groupHeader.count()) > 0) {
    await groupHeader.click();
    await page.waitForTimeout(500);
    const groupItems = page.locator(".sidebar-group-items .sidebar-item");
    if ((await groupItems.count()) > 0) {
      await groupItems.first().click();
      await page.waitForSelector(".panel", { timeout: 15000 });
      await page.waitForTimeout(2000);
      await page.screenshot({
        path: path.join(ROOT_DIR, "docs", "screenshot-sidebar-groups.png"),
        fullPage: false,
      });
      console.log("Saved: screenshot-sidebar-groups.png");
    }
  }

  await browser.close();
  console.log("Done!");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
