<h3 align="center">
  <br />
  <a href="https://github.com/yasinatesim/sklent"><img src=".github/assets/logo.svg" alt="Sklent" width="200" /></a>
  <br />
  Sklent
  <br />
</h3>

<hr />

<p align="center">Yeniden kullanılabilir bir <strong>Claude Code agent bootstrap</strong>. Agent'lar, skill'ler, reference'lar ve hook'lar <strong>BRAID</strong> akıl yürütme modeli etrafında kurgulanmış. Tek bir geliştirici production kodunu disiplinle üretip denetleyebilir. İçinde tam bir e-ticaret örneği gelir.</p>

<p align="center">
· <a href="./README.md">🇬🇧 English</a>
· <a href="./examples/e-commerce">🛒 E-ticaret örneği</a>
· <a href="./CONTRIBUTING.md">Katkı</a>
· <a href="./LICENSE">Lisans</a>
</p>

## Sklent nedir?

Sklent bir Claude Code agent bootstrap'ıdır. Agent'lar, skill'ler, reference'lar ve hook'lar proje kurallarını mekanik olarak zorunlu kılar.

Agent ekosistemi [`examples/e-commerce/.claude/`](examples/e-commerce/.claude) altında yaşıyor. Repo [`examples/e-commerce/`](examples/e-commerce) altında bir full-stack e-ticaret platformu örneğiyle birlikte gelir.

```
sklent/
├── examples/
│   └── e-commerce/
│       ├── .claude/              ← agent'lar, skill'ler, reference'lar, hook'lar
│       │   ├── agents/           kod reviewer'ları, BRAID solver, denetçiler
│       │   ├── skills/           braid-plan, ship-pr, coverage-gate, security-pentest
│       │   ├── references/       kodlama, backend, frontend, güvenlik standartları
│       │   └── hooks/            shell gate'ler: aliaslı import yok, korumalı branch'e commit yok
│       └── CLAUDE.md             her seansta yüklenen proje kuralları
```

> Tam proje referansı: [Ürün Gereksinim Dokümanı](examples/e-commerce/docs/PRD.md)

## BRAID akıl yürütme modeli

Bir agent'ın en kötü hatası hata büyümesidir. Bir adımdaki hata sonraki adımı besler. Düz bir yapılacaklar listesi başarısızlık durumunda ne olacağını tanımlamaz.

BRAID (Bounded Reasoning for Autonomous Inference and Decisions, arXiv 2512.15959) bir işi dört düğüm tipinden oluşan bir grafiğe böler: Constraint, Fact, Step, Check. Her Check'in iki çıkışı vardır: Pass ve Fail. Fail önceki bir Step'e döner. Bkz. [`examples/e-commerce/.claude/references/braid-mental-model.md`](examples/e-commerce/.claude/references/braid-mental-model.md).

## Agent ekosistemi

Tümü [`examples/e-commerce/.claude/`](examples/e-commerce/.claude) altında:

- **`agents/`** -- `wtf-code-reviewer` diff'i `wtf-go`, `wtf-js-react`, `wtf-security`, `wtf-ux-playwright`'a paralel gönderir. Ayrıca `braid-solver`, `constants-guard`, `issue-auditor`.
- **`skills/`** -- `braid-plan`, `spec-driven-development`, `coverage-gate`, `playwright-snapshot`, `ship-pr`, `issue-create`, `security-pentest` (web/api/network), `intended-vs-implemented`.
- **`references/`** -- dilden bağımsız kodlama standartları, backend/frontend/güvenlik standartları, git-flow ve BRAID modeli.
- **`hooks/`** -- aliaslı Go import'larını, korumalı branch'lere direkt commit'leri, agent merge'lerini ve 2 satırlık yorumları engelleyen shell script'leri. Commit öncesi CI-aynası doğrulaması çalıştırır.

## Örnek: Vela Commerce

Tek `docker compose up` ile çalışan markasız bir e-ticaret platformu. Go 1.25 API (Gin, GORM, Postgres), Next.js storefront ve admin (TR/EN), Iyzico 3D Secure sandbox, ChromaDB RAG ürün metni, GIB e-Arşiv fatura proxy'si, marketplace iskeletleri. Bkz. [`examples/e-commerce/`](examples/e-commerce).

```bash
cd examples/e-commerce && cp .env.example .env && docker compose up --build
```

## Fikirler nereden geliyor

Açık kaynaktan uyarlandı, atıf her agent ve skill içinde. Proje konsepti ve BRAID entegrasyon yaklaşımı [proje manifestosunda](https://gist.github.com/yasinatesim/bd5230ca0cc9b033c16280813c3ce6ff) belgelenmiştir:

- [Claude Code: Subagents](https://docs.claude.com/en/docs/claude-code/sub-agents)
- [Claude Code: Agent Skills](https://docs.claude.com/en/docs/claude-code/skills)
- [Claude Code: Hooks](https://docs.claude.com/en/docs/claude-code/hooks)
- [Anthropic Cookbook](https://github.com/anthropics/anthropic-cookbook)
- [mukul975/Anthropic-Cybersecurity-Skills](https://github.com/mukul975/Anthropic-Cybersecurity-Skills) (`security-pentest` ilhamı)
- [phuryn/pm-skills](https://github.com/phuryn/pm-skills) (`intended-vs-implemented` ilhamı)
- [BRAID arXiv 2512.15959](https://arxiv.org/abs/2512.15959)

## Lisans

MIT. Bkz. [LICENSE](./LICENSE).
