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

  const baseURL = process.env.BASE_URL || "http://localhost:5173";
  await page.goto(`${baseURL}/`);

  // Wait for login form inputs to render (not just the container, which may
  // show a "Loading..." state while /api/auth-info is being fetched)
  await page.waitForSelector("#userId", { timeout: 15000 });

  // Fill login form and submit
  await page.fill("#userId", "admin");
  await page.fill("#password", "admin");

  // Set up response listener before clicking to avoid race condition
  const loginResponse = page.waitForResponse(
    (resp) =>
      resp.url().includes("/api/login") &&
      resp.request().method() === "POST" &&
      resp.ok(),
    { timeout: 10000 }
  );
  await page.click('button[type="submit"]');

  // Wait for login API to respond successfully
  await loginResponse;

  // Reload the page so the app fetches dashboards with the session cookie
  await page.reload();

  // Wait for dashboard to render
  await page.waitForSelector(".sidebar", { timeout: 15000 });

  // Save signed-in state
  await context.storageState({ path: authFile });
  await browser.close();
}

export default globalSetup;
