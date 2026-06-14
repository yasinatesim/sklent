---
name: playwright-snapshot
description: Capture Playwright snapshots + screenshots for a list of URLs. Saves artifacts to .local-artifacts/screenshots/<branch>/. Lightweight helper invoked by wtf-ux-playwright or directly when verifying a specific surface.
---

# playwright-snapshot

Given a list of URLs, start (or reuse) the web dev server, visit each URL, and save a screenshot +
DOM snapshot to `.local-artifacts/screenshots/<branch>/`. Report any console error or failed
network request per URL. Evidence, not opinion.
