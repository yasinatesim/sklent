---
name: evidence-based-debugging
description: >
  Enforces methodical, evidence-based debugging instead of immediate guess-and-patch fixes.
  Use this whenever the user reports a bug, an error, or broken/unexpected behavior —
  phrases like "this isn't working", "it's broken", "throwing an error", "not behaving as
  expected" — even if they don't explicitly ask for debugging. Also trigger when a previous
  fix attempt didn't actually solve the problem, when the root cause isn't 100% obvious from
  the error message/stack trace alone, or when the user describes a symptom without pointing
  to the exact cause. Before writing any fix, this skill requires adding targeted
  logging/instrumentation, asking the user to reproduce the issue, and waiting for real
  console/terminal output — it never patches based on assumption alone.
---

# Evidence-Based Debugging

## Why this skill exists

When a bug report comes in, the easy path is: read the error, guess the most likely cause, change the code, say "fixed it." That feels fast but is often wrong — it can patch the wrong spot, introduce new side effects, or mask the symptom without solving the real problem. A real engineer doesn't do this — they *observe*, form a *hypothesis*, *test* it, and only change code once they have *evidence*.

This skill's job: in any conversation about fixing a bug, require gathering real runtime evidence before touching code whenever possible.

## Core rule

Unless the cause is 100% obvious from the error message/stack trace itself (see the "when to shortcut" section below), *do not jump straight to a code change*. Go through the debug loop first.

## The flow

### 1. Nail down the symptom

Turn what the user said into concrete facts:
- What was the expected behavior, and what actually happened?
- When/under what conditions does it trigger? Always, or intermittently?
- What are the exact steps to reproduce it?

If repro steps, the affected file/route/component/endpoint aren't clear, ask a short question instead of guessing.

### 2. Silently build a hypothesis list

Read the relevant code path and form 2-4 plausible causes (state timing, race conditions, a wrong prop/parameter, an unexpected API response shape, the wrong branch being taken, stale/cached data, etc.). You don't need to narrate this list to the user — its purpose is to tell you *where to put instrumentation* next.

### 3. Add instrumentation — change observability, not logic

Add logs/prints/debugger statements at the suspect points:
- Inputs coming into the function
- Key state/variable updates
- Which branch of an if/switch is actually taken
- Data sent to and received from API calls
- Resolve/reject moments in promises/async code

Use a prefix that's easy to spot and later remove, e.g. `[DEBUG]`:

```js
console.log('[DEBUG] handleSubmit input:', input);
console.log('[DEBUG] API response:', response.status, response.data);
```

At this stage, *don't change the behavior of the code* — only add observation. No "fix" yet.

### 4. Ask the user to reproduce it

Give a clear instruction: which scenario to re-run, and where to copy the output from (browser devtools console, terminal, a server log file, etc.). Example:

> "I added logs at 3 points: (1) the input at form submit, (2) the request body sent to the API, (3) the response coming back. Can you reproduce it again and paste the full [DEBUG] output here?"

Don't propose a fix at this step. Just wait for the data.

### 5. Analyze the real output

Compare the user's actual log output against your hypotheses. Find the *exact point* where the expected value and the real value diverge.

If the logs are inconclusive, or the problem turns out to live somewhere unexpected, say so explicitly, add more targeted logging, and go back to step 4:

> "All three of those look normal, so the issue is deeper than I assumed. I added another log point here — can you try once more?"

Repeat this until the root cause is actually clear. It's fine if it takes more than one round — the point is not to guess.

### 6. Apply a minimal, evidence-backed fix

Make a small, focused change targeting only the root cause you've confirmed with real output. Don't change several unrelated things at once "just in case" — patching five different spots based on "maybe it's this, maybe it's that" is the opposite of evidence-based debugging.

### 7. Clean up debug logs and get confirmation

- Remove the `[DEBUG]` logs you added (or, if persistent logging is genuinely useful, suggest converting them to a production-safe logger).
- Don't say "fixed it." Say something like: "The root cause was X, I changed Y — can you run through the same scenario once more to confirm it's actually resolved?"
- Don't treat the issue as closed until the user gives real confirmation.

## When this flow can be shortened

If the error message/stack trace already points to the root cause unambiguously (a typo, a wrong import, an obvious type error, a clearly-null reference), skip the logging round and fix it directly. Even then, still ask the user to verify after the fix — "this is probably it" and "confirmed" are not the same thing.

A deterministic logic bug that can be pinned down and proven with a unit test (write the failing assertion, watch it fail for the expected reason, fix, watch it pass) satisfies this skill's evidence bar without live console instrumentation — that IS the evidence. Reserve `[DEBUG]` console/log instrumentation for cases where the failure only manifests at runtime (live API responses, race conditions, environment-specific state) and can't be pinned down by reading code + a test alone.

## Don't

- Reply to "I'm getting error X" with "I changed this, try it" without gathering observation first.
- Push multiple untested, unrelated potential fixes into the code at once.
- Leave debug logs sitting in production code.
- Say "this should be fixed now" or "we can close this" before the user shares real output confirming it.
- Skip the user's actual repro steps in favor of a scenario you made up yourself.

## Example dialogue template

*User:* "Sometimes the user info comes back empty after login."

*You:*
1. "Does this happen every time, or intermittently? Does it go empty again after a page refresh?"
2. (Add [DEBUG] logs at the AuthContext, the login handler, and the API response points.)
3. "I added logs at 3 points: (1) the login API response, (2) the moment it's set on context, (3) the context value at render time. Can you reproduce it and paste the full console output?"
4. (Once the user pastes the logs, compare the real data against your hypotheses.)
5. Once the root cause is clear, apply the minimal fix, remove the debug logs, and ask for final confirmation.
