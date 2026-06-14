// Drives a real browser: admin login, create product via API, assert it renders in the store.
import { chromium } from "playwright";

const API = process.env.VERIFY_API ?? "http://localhost:8100";
const WEB = process.env.VERIFY_WEB ?? "http://localhost:3100";

const fail = (msg) => {
  console.error(`✗ ${msg}`);
  process.exit(1);
};

const readCsrf = async (ctx) => {
  await ctx.request.get(`${API}/healthz`);
  const cookies = await ctx.cookies();
  const csrf = cookies.find((c) => c.name === "csrf_token");
  if (!csrf) fail("csrf cookie was not primed by GET");
  return csrf.value;
};

const main = async () => {
  const browser = await chromium.launch();
  const ctx = await browser.newContext();
  const page = await ctx.newPage();

  const csrf = await readCsrf(ctx);
  const login = await ctx.request.post(`${API}/auth/login`, {
    headers: { "X-CSRF-Token": csrf, "Content-Type": "application/json" },
    data: { email: "admin@vela.test", password: "admin12345" },
  });
  if (!login.ok()) fail(`admin login failed: ${login.status()}`);

  const slug = `playwright-${Date.now()}`;
  const csrf2 = await readCsrf(ctx);
  const create = await ctx.request.post(`${API}/admin/products`, {
    headers: { "X-CSRF-Token": csrf2, "Content-Type": "application/json" },
    data: {
      slug,
      titleTr: "Playwright Test Ürünü",
      priceCents: 12345,
      stock: 7,
      published: true,
    },
  });
  if (create.status() !== 201) fail(`product create expected 201, got ${create.status()}`);

  await page.goto(`${WEB}/tr/urun/${slug}`, { waitUntil: "networkidle" });
  const heading = await page.locator("h1").first().textContent();
  if (!heading?.includes("Playwright")) fail("new product is not visible in the store");

  console.log(`✓ verified: "${heading}" visible at ${WEB}/tr/urun/${slug}`);
  await browser.close();
};

main().catch((err) => fail(err.message));
