# Comparison with Other Git Worktree Tools

This document compares **lazyworktree** with other Git worktree management tools.
It aims to be **honest and practical**, highlighting both where lazyworktree
excels and where **other tools are objectively better choices** depending on
constraints and workflow.

---

## High-Level Positioning

| Tool | Strength |
|----|----|
| **lazyworktree** | Full interactive environment for humans |
| git-worktree-runner (gtr) | Best CLI + scripting ergonomics |
| worktrunk (wt) | Parallel / AI-agent workflows |
| worktree-plus (wtp) | Minimal, predictable automation |
| branchlet | Lightweight TUI with low cognitive load |
| gwm | Fast fuzzy navigation |
| kosho | Command-centric, agent-oriented |
| wtm | Bare-repo / CI / server environments |

lazyworktree intentionally trades **simplicity and scriptability** for **interactive power**.

---

## Core Worktree Management

| Feature | lazyworktree | gtr | wt | wtp |
|-------|--------------|-----|----|-----|
| Create / delete worktrees | ✅ | ✅ | ✅ | ✅ |
| Rename worktrees | ✅ | ❌ | ❌ | ❌ |
| Absorb into main | ✅ | ❌ | ⚠️ manual | ❌ |
| Prune merged worktrees | ✅ | ⚠️ manual | ⚠️ manual | ⚠️ limited |
| Create from uncommitted changes | ✅ | ❌ | ❌ | ❌ |

### Where other tools win

* **gtr / wtp**: simpler mental model, fewer moving parts
* **wtp**: extremely predictable behaviour suitable for automation
* **wt**: optimized for creating many short-lived worktrees quickly

---

## Interface & Workflow

| Feature | lazyworktree | gtr | wtp | branchlet |
|-------|--------------|-----|-----|-----------|
| Full TUI | ✅ | ❌ | ❌ | ✅ |
| Zero-UI CLI | ❌ | ✅ | ✅ | ❌ |
| Works well over SSH / low latency | ⚠️ | ✅ | ✅ | ⚠️ |
| Easy to script | ❌ | ✅ | ✅ | ❌ |

### Clear advantages of other tools

* **gtr** is superior for:
  * shell pipelines
  * scripting
  * headless usage
* **wtp** is better when:
  * you want “do exactly one thing”
  * no UI dependencies
* **branchlet** is faster to understand for first-time users

lazyworktree is not optimized for scripting or non-interactive environments.

---

## Automation & Hooks

| Feature | lazyworktree | gtr | wt | wtp |
|-------|--------------|-----|----|-----|
| Hooks | ✅ | ✅ | ✅ | ✅ |
| Secure hook execution (TOFU) | ✅ | ❌ | ❌ | ❌ |
| Built-in automation primitives | ✅ | ❌ | ❌ | ❌ |
| Works without config | ❌ | ✅ | ✅ | ✅ |

### Where other tools win

* **gtr**:
  * hooks live in git config
  * easier to reason about in shared repos
* **wtp**:
  * fewer abstractions
  * easier to debug
* **wt**:
  * intentionally avoids policy decisions

lazyworktree’s automation is more powerful but **more complex**.

---

## Forge / PR Integration

| Feature | lazyworktree | others |
|-------|--------------|--------|
| PR/MR status | ✅ | ❌ |
| CI checks | ✅ | ❌ |
| Create worktree from PR | ✅ | ❌ |

### Trade-off

This integration:

* adds dependencies (`gh`, `glab`)
* may be undesirable in minimal or offline setups

Other tools avoid this entirely.

---

## tmux / Shell Integration

| Feature | lazyworktree | wt | gtr |
|-------|--------------|----|-----|
| tmux orchestration | ✅ | ⚠️ basic | ❌ |
| Shell jump | ✅ | ✅ | ⚠️ manual |
| Multi-window sessions | ✅ | ❌ | ❌ |

### Where others win

* **wt**:
  * simpler shell integration
  * easier to reason about
* **gtr**:
  * no tmux dependency
  * fewer assumptions

lazyworktree assumes tmux-heavy workflows.

---

## Configuration & Maintenance

| Aspect | lazyworktree | gtr | wtp |
|-----|--------------|-----|-----|
| Configuration size | Large | Small | Small |
| Learning curve | High | Low | Low |
| Failure modes | More | Fewer | Fewer |
| Upgrade risk | Higher | Lower | Lower |

This is an **explicit trade-off**:
 optimizes for capability, not minimalism.

---

## When NOT to use lazyworktree

lazyworktree is **not the best choice** if you:

* need headless or CI usage
* rely heavily on shell scripting
* want minimal dependencies
* prefer explicit Git commands
* manage worktrees mostly via automation

In these cases:

* use **gtr** or **wtp**

---

## Summary

**lazyworktree** provides the **largest feature surface** and the richest interactive experience for Git worktrees.

However, other tools are objectively better when:

* simplicity matters more than power
* scripting and automation are primary
* environments are constrained
* users prefer explicit Git semantics

lazyworktree is a **workspace manager for humans**.
Other tools remain excellent **worktree utilities for systems**.
