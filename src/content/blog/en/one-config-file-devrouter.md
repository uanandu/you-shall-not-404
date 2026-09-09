---
title: One config file, every local service
description: Design notes on DevRouter — a local dev proxy and self-hosted tunnel in Go I'm building, and why a single YAML file should beat a folder of scripts.
date: 2026-01-14
tags: [go, devrouter, tooling]
lang: en
draft: true
---

Every team I've joined has some version of the same problem: a handful of services that each need their own port, their own `.env`, and a README nobody updates. New hires spend their first day just getting `localhost` to behave.

**DevRouter** is what I'm building to fix that. The plan: one config file starts your services, routes them by hostname, and `devrouter share` exposes any of them publicly through a WebSocket tunnel you control — no third-party tunnel service, no surprise bill.

## Why hostname routing over ports

Ports are an implementation detail leaking into your daily workflow. `api.local` and `web.local` are addresses; `:3001` and `:5173` are trivia you have to remember. Once routing is by hostname, adding a new service is a config diff, not a context switch.

```yaml
services:
  api:
    host: api.local
    cmd: go run ./cmd/api
  web:
    host: web.local
    cmd: pnpm dev
```

## The tunnel is the hard part

Proxying locally is easy. Exposing a local service safely, without handing your traffic to someone else's edge network, is not. The plan for DevRouter's WebSocket tunnel is to keep it deliberately boring — no dynamic DNS magic, no plugin system — because the boring version is the one I'd trust enough to use for client demos, once it's built.

Go is the right call here specifically for speed of iteration: a tunnel client and server small enough to let me rewrite the framing protocol more than once while working out backpressure, without either rewrite eating a full day.
