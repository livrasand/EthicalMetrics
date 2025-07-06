# EthicalMetrics

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/livrasand/EthicalMetrics) [![Render Deploy](https://img.shields.io/badge/render-live-brightgreen?logo=render)](https://ethicalmetrics.onrender.com) [![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg?logo=go)](https://golang.org/dl/) [![Docker](https://img.shields.io/badge/docker-supported-blue?logo=docker)](https://hub.docker.com) [![Go Report Card](https://goreportcard.com/badge/github.com/livrasand/EthicalMetrics)](https://goreportcard.com/report/github.com/livrasand/EthicalMetrics) [![codecov](https://codecov.io/gh/livrasand/EthicalMetrics/branch/main/graph/badge.svg)](https://codecov.io/gh/livrasand/EthicalMetrics)

**EthicalMetrics** is a next-generation web analytics platform built with one radical goal:  
**protecting your users' privacy while delivering meaningful insights.**

No cookies. No fingerprinting. No personal data. No compromises.

> **Perfect for privacy-first projects, sensitive applications, decentralized platforms, or any product that truly values user trust.**

---

## Why EthicalMetrics?

- **True Anonymity by Design**  
  No IP collection. No device fingerprinting. No persistent IDs.  
  100% anonymous, with no way to track individual users — by default.

- **Self-Hosted, Zero External Dependencies**  
  No third-party scripts, CDNs, or vendors. You own your data.

- **End-to-End Encryption**  
  Data is protected at rest using Redis. Built-in, not bolted on.

- **Core, Actionable Metrics**  
  Track visits, module usage, performance, events, and engagement — not people.

- **Lightweight & Easy to Use**  
  A single vanilla JS snippet. No frameworks. No bloat.  
  Ethical by default — even without configuration.

- **Flexible Open API**  
  Send custom events and modules programmatically with minimal overhead.

- **Private Dashboards per Site**  
  Each site has its own dashboard secured with token-based access.

- **Portable by Nature**  
  Works on VPS, Docker, Render, or local environments. No lock-in.

EthicalMetrics exists to prove that **analytics can be useful without being creepy.**  
👉 [Read the Ethical Analytics Manifesto in `FUNDAMENTALS.md`](./FUNDAMENTALS.md)

EthicalMetrics is growing fast — without compromising ethics.  
👉 [Check out the roadmap here](./ROADMAP.md)

---

## Architecture Overview

Understanding how EthicalMetrics works internally helps you make the most of its privacy-first approach. The following diagrams illustrate the system's core components and data flow.

### System Architecture

The complete EthicalMetrics system consists of four main layers working together:

```mermaid
flowchart TD
    %% Clusters/Groups
    subgraph Client_Layer
        Website["Website Integration"]
        ethicalmetrics_js["ethicalmetrics.js"]
    end

    subgraph Web_Interface
        nuevo_html["pricing.html"]
        dashboard_html["dashboard.html"]
        index_html["index.html"]
    end

    subgraph Server_Layer
        main_go["cmd/server/main.go"]
        handlers_go["internal/api/handlers.go"]
    end

    subgraph Data_Layer
        database_go["internal/db/database.go"]
        event_go["internal/models/event.go"]
        metrics_db["Redis"]
        sites_table["sites"]
        events_table["events"]
    end

    %% Connections
    Website --> ethicalmetrics_js
    ethicalmetrics_js -.-> handlers_go
    dashboard_html -.-> handlers_go
    index_html -.-> handlers_go

    main_go --> handlers_go
    main_go --> nuevo_html
    main_go --> dashboard_html
    main_go --> index_html

    handlers_go --> database_go
    handlers_go --> event_go

    database_go --> metrics_db
    metrics_db --> sites_table
    metrics_db --> events_table

    %% Style for dotted lines
    linkStyle 2,3,4 stroke:#999,stroke-dasharray:3
```

### User Workflow

EthicalMetrics follows a simple three-stage process from registration to analytics:

```mermaid
flowchart TD
    %% Main Nodes
    START["Website Owner Starts"]
    REGISTER["Site Registration"]
    INTEGRATE["Client Integration"]
    ANALYTICS["View Analytics"]
    NUEVO_HANDLER["/nuevo API Endpoint"]
    TRACK_HANDLER["/track API Endpoint"]
    STATS_HANDLER["/stats API Endpoint"]
    REDIS_SITES["Redis Sites Storage"]
    REDIS_EVENTS["Redis Events Storage"]
    QUERY_DATA["Query aggregated data"]
    SITE_ID["site_id + admin_token"]
    EVENT_STORAGE["Anonymous Event Storage"]
    DASHBOARD_DATA["por_modulo + por_dia stats"]

    %% Connections
    START --> REGISTER
    REGISTER --> INTEGRATE
    REGISTER --> NUEVO_HANDLER
    INTEGRATE --> ANALYTICS
    INTEGRATE --> TRACK_HANDLER
    ANALYTICS --> STATS_HANDLER
    NUEVO_HANDLER --> REDIS_SITES
    TRACK_HANDLER --> REDIS_EVENTS
    STATS_HANDLER --> QUERY_DATA
    REDIS_SITES --> SITE_ID
    REDIS_EVENTS --> EVENT_STORAGE
    QUERY_DATA --> DASHBOARD_DATA

    %% Styling for better visualization
    classDef process fill:#e1f5fe,stroke:#039be5
    classDef storage fill:#e8f5e9,stroke:#43a047
    classDef endpoint fill:#fff3e0,stroke:#fb8c00
    classDef data fill:#f3e5f5,stroke:#8e24aa

    class START,REGISTER,INTEGRATE,ANALYTICS process
    class NUEVO_HANDLER,TRACK_HANDLER,STATS_HANDLER endpoint
    class REDIS_SITES,REDIS_EVENTS,QUERY_DATA storage
    class SITE_ID,EVENT_STORAGE,DASHBOARD_DATA data
```

### API Endpoints

The system exposes three main API endpoints that handle all core functionality:

```mermaid
flowchart TD
    %% Clusters/Groups
    subgraph HTTP_Endpoints
        track_endpoint["/track"]
        nuevo_endpoint["/nuevo"]
        stats_endpoint["/stats"]
        static_endpoint["/static/*"]
    end

    subgraph Handler_Functions
        TrackHandler["TrackHandler"]
        NuevoHandler["NuevoHandler"]
        StatsHandler["StatsHandler"]
        StaticHandler["http.FileServer"]
    end

    subgraph Database_Operations
        InsertEvent["INSERT INTO events"]
        InsertSite["INSERT INTO sites"]
        SelectStats["SELECT aggregated data"]
    end

    subgraph Response_Data
        EventModel["models.Event"]
        SiteResponse["site_id + admin_token"]
        StatsResponse["por_modulo + por_dia"]
    end

    %% Connections
    track_endpoint --> TrackHandler
    nuevo_endpoint --> NuevoHandler
    stats_endpoint --> StatsHandler
    static_endpoint --> StaticHandler

    TrackHandler --> InsertEvent
    NuevoHandler --> InsertSite
    StatsHandler --> SelectStats

    InsertEvent --> EventModel
    InsertSite --> SiteResponse
    SelectStats --> StatsResponse

    %% Styling
    classDef endpoint fill:#e3f2fd,stroke:#1976d2
    classDef handler fill:#e8f5e9,stroke:#388e3c
    classDef db_operation fill:#fff3e0,stroke:#ffa000
    classDef response fill:#f3e5f5,stroke:#8e24aa

    class HTTP_Endpoints endpoint
    class Handler_Functions handler
    class Database_Operations db_operation
    class Response_Data response
```

---

## Getting Started

1. Register your site with one clic at:

```
https://ethicalmetrics.onrender.com
````

2. Get your `site_id` and `admin_token`.

3. Add the tracking script to your HTML:

```html
<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=YOUR_SITE_ID"></script>
```

**Ultra-simple banner customization:**  
You can easily customize the consent banner by adding `data-*` attributes to the script tag:

```html
<script
  src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=YOUR_SITE_ID"
  data-banner-bg="#222"
  data-banner-color="#fff"
  data-banner-btn-bg="#4a90e2"
  data-banner-btn-color="#fff"
  data-banner-text="We respect your privacy. DNT is enabled. Only if you accept, we collect anonymous analytics."
  data-banner-btn-text="Accept"
></script>
```

Change colors and text as you wish — no coding required!

4. Optionally, define a custom module for richer context:

```html
<body data-modulo="home">
  ...
  <script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=YOUR_SITE_ID"></script>
</body>
```

5. Access your private dashboard:

```
https://ethicalmetrics.onrender.com/dashboard.html?site=YOUR_SITE_ID&token=YOUR_ADMIN_TOKEN
```

---

## Experimental & Evolving

This is a living project under active development.
We **welcome contributors**, audits, and ideas from the community.

* 🌱 Open issues for feedback or ideas
* 🔍 Fork or inspect the code — transparency is key
* 🛠️ Help us build the most **trustworthy** analytics platform on the web

### Polyrepo

EthicalMetrics is organized into multiple specialized repositories under a decoupled architecture. This allows contributors to contribute in a segmented manner, based on their interests or skills.

#### Backend
Main repository containing the server logic, API, and database:

→ [`github.com/livrasand/EthicalMetrics`](https://github.com/livrasand/EthicalMetrics)

#### Frontend
**Modular, decoupled static frontend with dynamic loading from CDN**, based on reusable HTML components and lightweight JavaScript for SPA-like navigation:

→ [`github.com/EthicalMetrics/website`](https://github.com/EthicalMetrics/website)

This frontend is distributed through [jsDelivr](https://www.jsdelivr.com/) for instant loading and real-time integration without external dependencies. You can contribute new components, visual improvements, or interactive features.

---

## Get Involved

Made with 💚 by [@livrasand](https://github.com/livrasand).
Want to collaborate, contribute, or ask something?

→ [Open an issue](https://github.com/livrasand/EthicalMetrics/issues)
→ [Submit a pull request](https://github.com/livrasand/EthicalMetrics/pulls)

---

**🛡️ Ethical by design. Anonymous by default. Metrics with ethics. Do EthicalMetrics!**

## License

This project is available under a dual license model:

- **GNU GPL v3** for personal, educational, or compatible open source use.
- **EthicalMetrics Commercial License** for commercial, SaaS, or proprietary products.

For a commercial license, see [`LICENSE-COMMERCIAL.md`](./LICENSE-COMMERCIAL.md).