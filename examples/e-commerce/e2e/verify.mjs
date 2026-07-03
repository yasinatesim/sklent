// Drives a real browser: admin login, exercise the full admin panel, assert each surface renders real data.
import { chromium } from "playwright";

const API = process.env.VERIFY_API ?? "http://localhost:8100";
const WEB = process.env.VERIFY_WEB ?? "http://localhost:3100";

const fail = (msg) => {
  console.error(`✗ ${msg}`);
  process.exit(1);
};

const ok = (msg) => console.log(`✓ ${msg}`);

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

  // --- product create + storefront render (pre-existing check) ---
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
    data: { slug, titleTr: "Playwright Test Ürünü", priceCents: 12345, stock: 7, published: true },
  });
  if (create.status() !== 201) fail(`product create expected 201, got ${create.status()}`);

  await page.goto(`${WEB}/tr/urun/${slug}`, { waitUntil: "networkidle" });
  const heading = await page.locator("h1").first().textContent();
  if (!heading?.includes("Playwright")) fail("new product is not visible in the store");
  ok(`product visible at ${WEB}/tr/urun/${slug}`);

  // --- real browser admin login via the /giris form, not just the API ---
  await page.goto(`${WEB}/tr/giris`, { waitUntil: "networkidle" });
  await page.fill('input[type="email"]', "admin@vela.test");
  await page.fill('input[type="password"]', "admin12345");
  await page.click('button[type="submit"]');
  await page.waitForURL(`${WEB}/tr/admin`, { timeout: 10000 });
  ok("browser login via /giris redirected to /admin");

  // --- stock tracking: create via API, assert it renders in the admin page ---
  const stockName = `Playwright Stok ${Date.now()}`;
  const csrf3 = await readCsrf(ctx);
  const stockCreate = await ctx.request.post(`${API}/admin/stock-tracking`, {
    headers: { "X-CSRF-Token": csrf3, "Content-Type": "application/json" },
    data: { productName: stockName, quantity: 11 },
  });
  if (stockCreate.status() !== 201) fail(`stock-tracking create expected 201, got ${stockCreate.status()}`);
  await page.goto(`${WEB}/tr/admin/stok-takibi`, { waitUntil: "networkidle" });
  await page.waitForSelector(`text=${stockName}`, { timeout: 10000 });
  ok("stock-tracking row visible at /admin/stok-takibi");

  // --- promotions: create via API, assert it renders in the admin page ---
  const promoName = `Playwright Kampanya ${Date.now()}`;
  const csrf4 = await readCsrf(ctx);
  const promoCreate = await ctx.request.post(`${API}/admin/promotions`, {
    headers: { "X-CSRF-Token": csrf4, "Content-Type": "application/json" },
    data: { name: promoName, discountType: "percent", discountValue: 12, scopeType: "all", minCartCents: 0 },
  });
  if (promoCreate.status() !== 201) fail(`promotion create expected 201, got ${promoCreate.status()}`);
  await page.goto(`${WEB}/tr/admin/kampanyalar`, { waitUntil: "networkidle" });
  await page.waitForSelector(`text=${promoName}`, { timeout: 10000 });
  ok("promotion visible at /admin/kampanyalar");

  // --- reviews: submit publicly, approve via API, assert it renders in the moderation page ---
  const csrf5 = await readCsrf(ctx);
  const reviewCreate = await ctx.request.post(`${API}/products/${slug}/reviews`, {
    headers: { "X-CSRF-Token": csrf5, "Content-Type": "application/json" },
    data: { authorName: "Playwright Reviewer", rating: 5, comment: "e2e comment" },
  });
  if (reviewCreate.status() !== 201) fail(`review create expected 201, got ${reviewCreate.status()}`);
  await page.goto(`${WEB}/tr/admin/yorumlar`, { waitUntil: "networkidle" });
  await page.waitForSelector("text=Playwright Reviewer", { timeout: 10000 });
  ok("pending review visible at /admin/yorumlar");

  // --- orders admin page renders without crashing ---
  await page.goto(`${WEB}/tr/admin/siparisler`, { waitUntil: "networkidle" });
  const ordersHeading = await page.locator("h2").first().textContent();
  if (!ordersHeading) fail("orders admin page did not render a heading");
  ok("orders admin page renders");

  await browser.close();
};

main().catch((err) => fail(err.message));
