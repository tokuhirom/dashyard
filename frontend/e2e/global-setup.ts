import { chromium, type FullConfig } from "@playwright/test";
import path from "path";
import fs from "fs";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const authDir = path.join(__dirname, ".auth");
const authFile = path.join(authDir, "session.json");

async function globalSetup(_config: FullConfig) {
  fs.mkdirSync(authDir, { recursive: true });

  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();

  await page.goto("http://localhost:5173/");

  // Wait for login form to render (app starts optimistic, then shows login after 401)
  await page.waitForSelector(".login-form", { timeout: 15000 });

  // Fill login form and submit
  await page.fill("#userId", "admin");
  await page.fill("#password", "admin");
  await page.click('button[type="submit"]');

  // Wait for login API to respond successfully
  await page.waitForResponse(
    (resp) =>
      resp.url().includes("/api/login") &&
      resp.request().method() === "POST" &&
      resp.ok(),
    { timeout: 10000 }
  );

  // Reload the page so the app fetches dashboards with the session cookie
  await page.reload();

  // Wait for dashboard to render
  await page.waitForSelector(".sidebar", { timeout: 15000 });

  // Save signed-in state
  await context.storageState({ path: authFile });
  await browser.close();
}

export default globalSetup;
