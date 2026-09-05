package main

func seedFS() *node {
	return dir("/",
		dir("home",
			dir("anandu",
				file("README.md", `
# anandu — fullstack & devops engineer

  Location  Montréal, QC, Canada
  Role      Consultant @ Deloitte Canada
  Contact   anandu_dev@icloud.com
  Web       https://anandu.dev

── Quick navigation ──────────────────────────────────────
  ls projects/          three interconnected Go tools
  cat skills.txt        full technical stack
  cat experience/       career history
  neofetch              system overview

Type 'help' for all available commands.
`),
				file("skills.txt", `
── Languages ─────────────────────────────────────────────
  Go (Golang)      Primary backend language
  Rust             Systems / performance-critical work
  TypeScript / JS  Full frontend and Node tooling
  Node.js          Server-side JS runtimes

── Infrastructure ────────────────────────────────────────
  Docker           Container build and runtime
  Kubernetes        Cluster orchestration
  Helm             K8s package management
  ArgoCD           GitOps continuous delivery
  Azure            Cloud platform (AZ-900 certified)

── Frontend ──────────────────────────────────────────────
  Svelte           Reactive UI (this site)
  Astro            Static site generation
  Shopify Liquid   Storefront templating

── Data ──────────────────────────────────────────────────
  PostgreSQL       Primary relational database

── Certifications ────────────────────────────────────────
  AZ-900           Microsoft Azure Fundamentals (2024)
`),
				file("contact.txt", `
  Email     anandu_dev@icloud.com
  GitHub    github.com/anandu-dev
  LinkedIn  linkedin.com/in/anandu-dev
  Location  Montréal, QC, Canada
`),
				dir("projects",
					file("README.md", `
Three tools built to work together:

  api-darwin/     API observability + architecture intelligence
  devrouter/      Local dev proxy + self-hosted WebSocket tunnel
  infra-mcp/      MCP server — LLM ↔ infrastructure bridge

Run 'cat <project>/README.md' for details.
`),
					dir("api-darwin",
						file("README.md", `
# API Darwin

API observability and architecture intelligence platform.

Ingests OpenAPI specs and live traffic event logs. Builds a directed
endpoint relationship graph. Scores API health across four axes:

  - Naming consistency   (conventions, casing, plurality)
  - Usage patterns       (hit rates, dead endpoints)
  - Error rates          (4xx/5xx concentration)
  - Duplication          (overlapping resource paths)

Surfaces AI-powered recommendations:
  → Consolidation candidates
  → Naming normalisation fixes
  → Deprecation targets

Advises only. Never modifies production resources.

Stack    Go, Next.js, PostgreSQL, Redis, Docker
Repo     github.com/anandu-dev/api-darwin
Status   In progress
`),
						file("main.go", `
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/anandu-dev/api-darwin/internal/ingestion"
	"github.com/anandu-dev/api-darwin/internal/graph"
	"github.com/anandu-dev/api-darwin/internal/scoring"
)

func main() {
	ctx := context.Background()

	ingestor := ingestion.New()
	g := graph.Build(ctx, ingestor)
	scores := scoring.Run(g)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: newRouter(g, scores),
	}
	log.Fatal(srv.ListenAndServe())
}
`),
					),
					dir("devrouter",
						file("README.md", `
# DevRouter

Local development proxy and self-hosted tunnel. Written in Go.

  devrouter start       reads devrouter.yml, starts all services,
                        routes by hostname on :80
  devrouter share web   exposes 'web' service publicly via WebSocket
                        tunnel — your own ngrok

No /etc/hosts edits. No cloud dependency. You own the tunnel server.

Stack    Go, WebSocket, Cobra CLI
Repo     github.com/anandu-dev/devrouter
Status   In progress
`),
						file("devrouter.yml", `
version: 1
tunnel:
  server: tunnel.anandu.dev
  token:  ${DEVROUTER_TOKEN}

services:
  web:
    command: pnpm dev
    port:    3000
    host:    web.localhost
  api:
    command: go run ./cmd/api
    port:    8080
    host:    api.localhost
  db-ui:
    command: npx drizzle-studio
    port:    4983
    host:    db.localhost
`),
					),
					dir("infra-mcp",
						file("README.md", `
# Infra MCP Server

A Model Context Protocol (MCP) server in Go that exposes live
infrastructure tooling to LLMs — safely, without uncontrolled access.

Capabilities exposed to Claude / any MCP-compatible client:

  query_cluster       Kubernetes cluster state snapshot
  list_containers     Running containers + health status
  api_darwin_score    Pull health scores for any API endpoint
  tunnel_expose       Trigger a DevRouter share session

All operations are read-heavy by design. Mutations require explicit
confirmation through the MCP tool-call surface.

Stack    Go, MCP protocol, Kubernetes client-go, Docker SDK
Repo     github.com/anandu-dev/infra-mcp
Status   Planning
`),
					),
				),
				dir("experience",
					file("README.md", `
  deloitte.txt        Jan 2023 – Present
  off-paper.txt       Oct 2022 – Jan 2023
  ypl.txt             Mar 2022 – Oct 2022
  techni-mek.txt      Oct 2021 – Mar 2022
`),
					file("deloitte.txt", `
Deloitte Canada                                Jan 2023 – Present
Stack: Go · JavaScript · Docker · Kubernetes · Azure
─────────────────────────────────────────────────────────────────

Solutions Specialist                           Jun 2025 – Present
  Leading technical delivery on client engagements — fullstack
  development, cloud infrastructure, and DevOps practice across
  enterprise-scale systems.

Solutions Analyst                              Jan 2023 – Jun 2025
  Fullstack and DevOps engineering across public sector mandates.

  ▸ ESDC (Employment and Social Development Canada)
    System-to-system integration connecting legacy federal services.
    REST and SOAP bridges, data transformation pipelines, and
    event-driven messaging.

  ▸ Government of Quebec
    API integration and middleware development. Designed service
    contracts, led technical workshops for client teams, produced
    architecture documentation and runbooks.
`),
					file("off-paper.txt", `
Off-Paper Creative Agency                      Oct 2022 – Jan 2023
Stack: Shopify · Liquid · HTML · CSS · JavaScript
─────────────────────────────────────────────────────────────────

Shopify Web Developer
  Built editorial storefronts for a boutique creative agency.
  Converted visual designs into functional web pages using Shopify
  Liquid. Clients included Basma Beauty and Sol Botanicals.
  Collaborated across design, development, and management stages.
`),
					file("ypl.txt", `
YPL Consulting Inc.                            Mar 2022 – Oct 2022
Stack: HTML · CSS · JavaScript · Linux · DevOps
─────────────────────────────────────────────────────────────────

Contract Web Developer                         Aug 2022 – Oct 2022
  Built client-facing web projects on contract.

  ▸ Client Portfolio
    Portfolio website designed and developed from scratch.

  ▸ CRM Application
    Customer Relationship Management app with payment and
    appointment systems. Delivered as a web app and cross-platform.

Mentee                                         Mar 2022 – Jul 2022
  Mentorship program shadowing the team across DevOps and Linux
  system administration. Automated tooling, tested and deployed
  applications, built with HTML and CSS.
`),
					file("techni-mek.txt", `
techni-mek inc                                 Oct 2021 – Mar 2022
Stack: HVAC
─────────────────────────────────────────────────────────────────

Heating & Air Conditioning Technician Apprentice
  Apprenticeship in HVAC systems installation and maintenance.
  (Career pivot to software engineering followed immediately after.)
`),
				),
			),
		),
		dir("etc",
			file("os-release", `NAME="PortfolioOS"
VERSION="1.0.0"
ID=portfolio
PRETTY_NAME="PortfolioOS 1.0.0 (WebAssembly)"
HOME_URL="https://anandu.dev"
SUPPORT_URL="mailto:anandu_dev@icloud.com"
`),
			file("hostname", "portfolio\n"),
		),
		dir("usr",
			dir("local",
				dir("bin",
					file("api-darwin", `api-darwin — API observability platform

USAGE
  api-darwin [command] [flags]

COMMANDS
  ingest   --spec <openapi.yaml>    ingest an OpenAPI specification
  score    --endpoint <path>        score a single endpoint
  report   --format (json|text)     generate full API health report
  serve    --port <n>               start the web dashboard

FLAGS
  --config  path to config file (default: .api-darwin.yml)
  --db      postgres connection string
  --help    show this help

REPO
  github.com/anandu-dev/api-darwin
`),
					file("devrouter", `devrouter — local dev proxy + self-hosted tunnel

USAGE
  devrouter [command] [flags]

COMMANDS
  start            start all services defined in devrouter.yml
  stop             stop all managed services
  share <service>  expose a service publicly via WebSocket tunnel
  status           show status of all services
  logs <service>   tail service logs

FLAGS
  --config  path to devrouter.yml (default: ./devrouter.yml)
  --help    show this help

REPO
  github.com/anandu-dev/devrouter
`),
				),
			),
		),
	)
}
