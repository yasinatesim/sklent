<h3 align="center">
  <br />
   <a href="https://github.com/yasinatesim/vela-commerce"><img src=".github/assets/logo.svg" alt="Vela Commerce" width="200" /></a>
  <br />
Vela Commerce
  <br />
</h3>

<hr />

<p align="center">Açık kaynaklı, yapay zeka ajanlarıyla geliştirilen bir e-ticaret referans platformu — gerçek bir prodüksiyon mağazasının markasız bir damıtımı. Tek bir geliştiricinin, disiplinli bir Claude Code ajan, beceri (skill), hook ve BRAID akıl yürütme modeli ekosistemiyle tam kapsamlı bir e-ticaret platformunu nasıl büyütüp sürdürebileceğini göstermek için yapıldı.</p>

<p align="center">
· <a href="./README.md">🇬🇧 English</a>
· <a href="https://github.com/yasinatesim/vela-commerce/issues">Konular (Issues)</a>
· <a href="./CONTRIBUTING.md">Katkıda Bulunma</a>
· <a href="https://github.com/yasinatesim/vela-commerce/blob/master/LICENSE">Lisans</a>
</p>

## 📖 Hakkında

Vela Commerce, **"AI Agent'lar ile Sıfırdan Uçtan Uca E-Ticaret Projesi Geliştirme"** yazısının eşlik eden deposudur. Bilinçli olarak **markasızdır** ve tek bir `docker compose up` komutuyla uçtan uca çalıştırılabilir.

Bu deponun amacı kodun **kendisi değildir**. Asıl amaç, kodu üreten ve koruyan *ajan ekosistemidir*: kısıtlar, gerçekler, adımlar ve kontroller — yanlış olanı yapmak mekanik olarak imkânsız hâle gelecek şekilde birbirine bağlanmıştır.

> **Varlık ≠ tamamlanmışlık.** Pazaryeri senkronizasyonu, ödeme ve fatura referans iskeletleridir. Herhangi bir dış gönderim/içe aktarım yolunun uçtan uca çalıştığını varsaymadan önce kodu inceleyin.

### 📚 Teknoloji Yığını

<table>
<tr>
  <td> <a href="https://go.dev">Go 1.25</a></td>
  <td>Gin, GORM ve PostgreSQL 16 ile backend API.</td>
</tr>
<tr>
  <td> <a href="https://nextjs.org">Next.js</a></td>
  <td>React 19 mağaza + admin paneli, CSS Modules + SCSS, Tailwind yok.</td>
</tr>
<tr>
  <td> <a href="https://next-intl.dev">next-intl</a></td>
  <td>JSON tabanlı çoklu dil desteği, Türkçe + İngilizce.</td>
</tr>
<tr>
  <td> <a href="https://www.postgresql.org">PostgreSQL 16</a></td>
  <td>Katalog, siparişler, stok ve kanal bazlı sayaçlar.</td>
</tr>
<tr>
  <td> <a href="https://www.trychroma.com">ChromaDB</a></td>
  <td>RAG destekli ürün metni üretimi için vektör erişimi.</td>
</tr>
<tr>
  <td> <a href="https://www.iyzico.com">Iyzico</a></td>
  <td>3D Secure ödeme entegrasyonu (sandbox).</td>
</tr>
<tr>
  <td> <a href="https://www.docker.com">Docker Compose</a></td>
  <td>Postgres + ChromaDB + Go API + Next.js web, tek komutla ayağa kalkar.</td>
</tr>
<tr>
  <td> <a href="https://docs.claude.com/en/docs/claude-code">Claude Code</a></td>
  <td>Depoyu geliştiren ve koruyan ajanlar, beceriler, hook'lar ve BRAID akıl yürütme modeli.</td>
</tr>
</table>

## 🧐 İçinde neler var?

| Alan | Tercih |
|---|---|
| Backend | Go 1.25, Gin, GORM, PostgreSQL 16 |
| Frontend | Next.js, React 19, CSS Modules + SCSS, next-intl (TR + EN, JSON tabanlı) |
| Kimlik Doğrulama | JWT 15dk access + dönen 7g refresh, httpOnly cookie, CSRF double-submit, rate limit |
| Sepet | Misafir (oturum cookie) + üye; 15 dakikalık stok rezervasyonu |
| Ödeme | Iyzico 3D Secure (sandbox), `internal/payment/iyzico` altında |
| Promosyonlar | Yüzde / sabit-TL, sepet/ürün/kategori kapsamı, kupon motoru — saf `Evaluate` |
| Pazaryeri | `internal/marketplace/{hb,ty}` altında iskelet istemciler (kategori ağacı, öznitelik, yayınlama) |
| LLM / RAG | Takılabilir sağlayıcı kaydı + ChromaDB destekli ürün metni üretimi (`internal/rag`) |
| Fatura | GIB e-Arşiv proxy konsepti (host izin listeli) |
| E-posta | SMTP + şablonlu sipariş/kargo mailleri, istek yolu dışında gönderilir |
| Dağıtım | `docker compose up`: Postgres + ChromaDB + Go API + Next.js web |

### Ajan ekosistemi

Her şey [`.claude/`](.claude) altında yaşar:

- **`agents/`** — `braid-solver`, `constants-guard`, `wtf-code-reviewer` (yönlendirici) ve dil/alan gözden geçiricileri `wtf-go`, `wtf-js-react`, `wtf-security`, `wtf-ux-playwright`, ayrıca `issue-auditor`.
- **`skills/`** — `braid-plan`, `spec-driven-development`, `coverage-gate`, `playwright-snapshot`, `ship-pr`, `issue-create`, `security-pentest` ve gözden geçirici beceriler.
- **`references/`** — dilden bağımsız kodlama standartları, backend/frontend standartları, güvenlik standartları, git-flow ve BRAID zihinsel modeli.
- **`hooks/`** — kuralları atlamayı imkânsız kılan shell betikleri (aliaslı Go import'u yok, korumalı dallara doğrudan commit yok, ajan kaynaklı merge yok, uzun yorum yok…).

[`CLAUDE.md`](CLAUDE.md) her oturuma yüklenir ve proje kurallarının tek kaynağıdır.

## Başlangıç

### 📦 Ön Koşullar

- [Docker](https://www.docker.com) + Docker Compose
- [Go 1.25+](https://go.dev) ve [Node.js 20+](https://nodejs.org) (yalnızca Docker'sız yerel geliştirme için)

### ⚙️ Nasıl Kullanılır

```bash
cp .env.example .env
docker compose up --build      # Postgres + ChromaDB + API + web
docker compose run --rm api ./bin/seed   # 8 kategori, 25 ürün, admin kullanıcı
open http://localhost:3100/tr  # mağaza (TR/EN, yeşil tema, açık/koyu)
```

Varsayılan admin (seed): `admin@vela.test` / `admin12345`. `/tr/giris` adresinden giriş yapın; adminler `/tr/admin` paneline yönlenir (kontrol paneli, tam oluşturma formlu ürünler, kampanyalar, kuponlar, siparişler).

**Bir ajan gibi doğrulayın:**

```bash
node e2e/verify.mjs            # gerçek bir tarayıcı sürer: giriş yap, ürün ekle, mağazada gör
```

Doğrulama betiği API yanıtına değil, **kullanıcının gördüğü ekrana** güvenir. Burada yeşil görünüyorsa akış gerçekten çalışıyordur.

**Yerel geliştirme (Docker'sız):**

```bash
# Backend
cd api && go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm install && npm run dev
```

## 🔑 Lisans

* Telif Hakkı © 2026 - MIT Lisansı.

Daha fazla bilgi için [LICENSE](https://github.com/yasinatesim/vela-commerce/blob/master/LICENSE) dosyasına bakın.

---

_This README was generated with by [markdown-manager](https://github.com/yasinatesim/markdown-manager)_ 🥲
