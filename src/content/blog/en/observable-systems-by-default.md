---
title: Observable systems by default
description: Why I want to stop bolting on logging after the fact, and start designing endpoints around the traces they'd need to produce.
date: 2026-03-02
tags: [devops, observability, api]
lang: en
draft: true
---

Most observability work happens after something has already gone wrong. A service misbehaves, nobody can explain why, and the fix is a hasty round of `console.log` or a Datadog integration wired in under pressure. It works, but it's reactive — the system was never designed to be watched, it was designed and then *interrogated*.

Designing **API Darwin** is what's changing how I think about this. The plan is to ingest OpenAPI specs and traffic logs to build an endpoint relationship graph — which forces a different question up front: not "how do I debug this later" but "what does this endpoint need to say about itself right now."

## Three habits I'm building around

1. **Structured logs from the first commit.** Not because you'll need them today, but because retrofitting structure into free-text logs later is genuinely miserable.
2. **Name things for the graph, not the file.** An endpoint's name should describe its relationship to the rest of the system, not just its route — the kind of refactor that tends to pay off most on inherited codebases.
3. **Treat deprecation as a first-class state**, not an afterthought comment. If a system can't tell you what's dying, it can't tell you what's safe to build on top of.

None of this is exotic. It's mostly discipline — the kind that's easy to defer under a deadline and expensive to skip once a system has real users.

## The bet

Endpoints designed with observability in mind from the start should be the ones that get consolidated correctly during a refactor — a graph that won't lie about which routes are actually load-bearing versus which ones just have the most historical traffic. That's the distinction I'm building API Darwin to surface, before a real migration has to pay to find it out the hard way.
