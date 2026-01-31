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
const BASE_URL = process.env.BASE_URL || "http://localhost:5173/";
const OUTPUT_DIR = process.env.OUTPUT_DIR || ROOT_DIR;

// Fixed absolute time window for deterministic screenshots.
// Combined with the deterministic dummyprom noise function, this ensures
// identical chart data and axis labels across runs.
const TIME_FROM = "2025-01-01T00:00:00Z";
const TIME_TO = "2025-01-01T01:00:00Z";
const TIME_QS = `from=${encodeURIComponent(TIME_FROM)}&to=${encodeURIComponent(TIME_TO)}`;

function dashboardUrl(dashPath: string): string {
  const base = BASE_URL.replace(/\/$/, "");
  return `${base}/d/${dashPath}?${TIME_QS}`;
}

async function main() {
  fs.mkdirSync(path.join(OUTPUT_DIR, "docs"), { recursive: true });

  const browser = await chromium.launch();
  const context = await browser.newContext({
    viewport: { width: 1280, height: 800 },
  });
  const page = await context.newPage();

  // Go to login page
  await page.goto(BASE_URL);

  // Wait for login form
  await page.waitForSelector(".login-form", { timeout: 15000 });

  // Screenshot: Login page
  await page.screenshot({
    path: path.join(OUTPUT_DIR, "docs", "screenshot-login.png"),
    fullPage: false,
  });
  console.log("Saved: screenshot-login.png");

  // Login
  await page.fill("#userId", "admin");
  await page.fill("#password", "admin");
  const loginResponsePromise = page.waitForResponse(
    (resp) =>
      resp.url().includes("/api/login") &&
      resp.request().method() === "POST" &&
      resp.ok(),
    { timeout: 10000 }
  );
  await page.click('button[type="submit"]');
  await loginResponsePromise;

  // Navigate to overview dashboard with absolute time range
  await page.goto(dashboardUrl("overview"));
  await page.waitForSelector(".graph-panel canvas", { timeout: 15000 });
  // Give Chart.js a moment to finish animations
  await page.waitForTimeout(2000);

  // Screenshot: Main dashboard view (overview)
  await page.screenshot({
    path: path.join(OUTPUT_DIR, "screenshot.png"),
    fullPage: false,
  });
  console.log("Saved: screenshot.png (main)");

  // Screenshot: Full page dashboard
  await page.screenshot({
    path: path.join(OUTPUT_DIR, "docs", "screenshot-dashboard.png"),
    fullPage: true,
  });
  console.log("Saved: screenshot-dashboard.png");

  // Navigate to a dashboard with variables (network-variable)
  await page.goto(dashboardUrl("network-variable"));
  await page.waitForSelector(".graph-panel canvas", { timeout: 15000 });
  await page.waitForTimeout(2000);
  await page.screenshot({
    path: path.join(OUTPUT_DIR, "docs", "screenshot-variables.png"),
    fullPage: false,
  });
  console.log("Saved: screenshot-variables.png");

  // Navigate to a dashboard with repeat rows (network-repeat)
  await page.goto(dashboardUrl("network-repeat"));
  await page.waitForSelector(".graph-panel canvas", { timeout: 15000 });
  await page.waitForTimeout(2000);
  await page.screenshot({
    path: path.join(OUTPUT_DIR, "docs", "screenshot-repeat.png"),
    fullPage: true,
  });
  console.log("Saved: screenshot-repeat.png");

  // Navigate to chart-types dashboard
  await page.goto(dashboardUrl("chart-types"));
  await page.waitForSelector(".panel", { timeout: 15000 });
  await page.waitForTimeout(2000);
  await page.screenshot({
    path: path.join(OUTPUT_DIR, "docs", "screenshot-chart-types.png"),
    fullPage: true,
  });
  console.log("Saved: screenshot-chart-types.png");

  // Navigate to thresholds dashboard using a fresh page so Chart.js
  // annotation plugin initialises cleanly (accumulated chart instances from
  // earlier dashboards can interfere with annotation rendering).
  {
    const freshPage = await context.newPage();
    await freshPage.goto(dashboardUrl("thresholds"));
    await freshPage.waitForSelector(".graph-panel canvas", {
      timeout: 15000,
    });
    // Annotation plugin needs extra time to render after async data loads
    await freshPage.waitForTimeout(10000);
    await freshPage.screenshot({
      path: path.join(OUTPUT_DIR, "docs", "screenshot-thresholds.png"),
      fullPage: true,
    });
    console.log("Saved: screenshot-thresholds.png");
    await freshPage.close();
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
        path: path.join(OUTPUT_DIR, "docs", "screenshot-sidebar-groups.png"),
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
