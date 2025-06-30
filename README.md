# EthicalMetrics

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

## Getting Started

1. Register your site at:

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

---

## Get Involved

Made with 💚 by [@livrasand](https://github.com/livrasand).
Want to collaborate, contribute, or ask something?

→ [Open an issue](https://github.com/livrasand/EthicalMetrics/issues)
→ [Submit a pull request](https://github.com/livrasand/EthicalMetrics/pulls)

---

**🛡️ Ethical by design. Anonymous by default. Metrics with ethics. Do EthicalMetrics.**
