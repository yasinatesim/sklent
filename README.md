<h3 align="center">
  <br />
  <a href="https://github.com/yasinatesim/vela-commerce"><img src=".github/assets/logo.svg" alt="Sklent" width="200" /></a>
  <br />
  Sklent
  <br />
</h3>

<hr />

<p align="center">A reusable <strong>Claude Code agent bootstrap</strong> — agents, skills, references and hooks, wired around the <strong>BRAID</strong> reasoning model, so a single developer can build and guard production code with discipline. Drop it into any project. Ships with a full e-commerce example.</p>

<p align="center">
· <a href="./README.tr.md">🇹🇷 Türkçe</a>
· <a href="./examples/e-commerce">🛒 E-commerce example</a>
· <a href="./CONTRIBUTING.md">Contributing</a>
· <a href="./LICENSE">License</a>
</p>

## 📖 What is Sklent?

Sklent is **not** an app. It is the **ecosystem that produces and guards apps**: a set of Claude
Code agents, skills, references and hooks that turn project rules into something mechanically
enforced rather than politely hoped for.

The whole point lives under [`examples/e-commerce/.claude/`](examples/e-commerce/.claude). To prove it works on something real, the repo
ships one large worked example — a brandless, full-stack e-commerce platform — under
[`examples/e-commerce/`](examples/e-commerce).

```
sklent/
├── examples/
│   └── e-commerce/
│       ├── .claude/              ← THE POINT
│       │   ├── agents/           reviewers + BRAID solver + constants/issue auditors
│       │   ├── skills/           braid-plan, ship-pr, coverage-gate, security-pentest, …
│       │   ├── references/       coding/backend/frontend/security standards + BRAID model
│       │   └── hooks/            shell gates: no aliased imports, no protected-branch commits, no merges…
│       └── CLAUDE.md             single source of project rules, loaded every session
```

> Full project reference: [Product Requirements Document](examples/e-commerce/docs/PRD.md)

## 🧠 The BRAID reasoning model

The signature of Sklent. An agent's worst failure is **error compounding**: a mistake in one
step feeds the next. A flat to-do list never says *what happens on failure*.

BRAID (*Bounded Reasoning for Autonomous Inference and Decisions*, arXiv 2512.15959) structures a
task as a graph of four node types — **Constraint**, **Fact**, **Step**, **Check** — where every
Check has exactly two edges, **Pass** and **Fail**, and the Fail edge loops back to an earlier
Step. The loop *is* the retry. See [`examples/e-commerce/.claude/references/braid-mental-model.md`](examples/e-commerce/.claude/references/braid-mental-model.md).

## 🤖 The agent ecosystem

All of the following lives under [`examples/e-commerce/.claude/`](examples/e-commerce/.claude):

- **`agents/`** — `wtf-code-reviewer` (dispatcher) routes a diff to `wtf-go`, `wtf-js-react`,
  `wtf-security`, `wtf-ux-playwright` in parallel; plus `braid-solver`, `constants-guard`,
  `issue-auditor`.
- **`skills/`** — `braid-plan`, `spec-driven-development`, `coverage-gate`, `playwright-snapshot`,
  `ship-pr`, `issue-create`, `security-pentest` (web/api/network), `intended-vs-implemented`.
- **`references/`** — language-agnostic coding standards, backend/frontend/security standards,
  git-flow, and the BRAID mental model — the rejection criteria reviewers read.
- **`hooks/`** — shell scripts that make the rules impossible to skip: no aliased Go imports, no
  direct commits to protected branches, no agent-initiated merges, no 2-line comments, CI-mirror
  verify before commit.

## 🛒 The example: Vela Commerce

A real, brandless e-commerce platform that runs with one `docker compose up`: Go 1.25 API
(Gin + GORM + Postgres), Next.js storefront + admin (TR/EN), Iyzico 3D Secure sandbox, ChromaDB
RAG product copy, GIB e-Arşiv invoice proxy, marketplace skeletons. See
[`examples/e-commerce/`](examples/e-commerce).

```bash
cd examples/e-commerce && cp .env.example .env && docker compose up --build
```

## 🔍 Where the ideas come from

Adapted from public work, with attribution kept in each agent/skill. The overall project concept and BRAID integration approach is documented in the [project manifesto](https://gist.github.com/yasinatesim/bd5230ca0cc9b033c16280813c3ce6ff):

- [Claude Code — Subagents](https://docs.claude.com/en/docs/claude-code/sub-agents)
- [Claude Code — Agent Skills](https://docs.claude.com/en/docs/claude-code/skills)
- [Claude Code — Hooks](https://docs.claude.com/en/docs/claude-code/hooks)
- [Anthropic Cookbook](https://github.com/anthropics/anthropic-cookbook)
- [mukul975/Anthropic-Cybersecurity-Skills](https://github.com/mukul975/Anthropic-Cybersecurity-Skills) — inspiration for `security-pentest`
- [phuryn/pm-skills](https://github.com/phuryn/pm-skills) — inspiration for `intended-vs-implemented`
- [BRAID (arXiv 2512.15959)](https://arxiv.org/abs/2512.15959)

## 🔑 License

MIT — see [LICENSE](./LICENSE).
