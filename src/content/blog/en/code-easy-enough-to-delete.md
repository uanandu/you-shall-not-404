---
title: Code easy enough to delete
description: On writing systems for public sector clients where the real deliverable is the ability to walk away cleanly.
date: 2025-11-08
tags: [devops, consulting, kubernetes]
lang: en
draft: true
---

Consulting on government infrastructure teaches you a lesson that's easy to skip when you're building for yourself: the best system isn't the one that's cleverest to build, it's the one the next team can safely delete a piece of.

Working across engagements with ESDC and the Government of Quebec, the systems that held up under staff turnover, budget cycles, and audits weren't the ones with the most sophisticated abstractions. They were the ones where a single Helm chart, a single ArgoCD app, or a single Terraform module could be removed without anyone needing to first reverse-engineer what it was secretly load-bearing for.

## What "easy to delete" actually means

- **No implicit dependencies between services.** If removing one component requires archaeology in a Slack channel to find out what else breaks, it wasn't decoupled — it just looked decoupled.
- **Config lives with the thing it configures.** Centralized config repos are convenient until the person who understands the mapping leaves.
- **Every deploy is reproducible from source**, not from a runbook someone half-remembers. ArgoCD forces this by construction, which is most of why I reach for it by default.

## The uncomfortable part

This means resisting the urge to build the elegant, general version of something when a client only needs the specific version today. Generality is a bet on a future team's context matching yours — and in public sector work, that team is rarely the one you're sitting across from.

The systems I'm proudest of aren't the ones doing the most. They're the ones where, a year later, deleting the wrong file doesn't take down three unrelated things.
