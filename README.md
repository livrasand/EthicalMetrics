# EthicalMetrics (Experimental)

**EthicalMetrics** is a next-generation web analytics platform built with one goal in mind: *protecting your users' privacy while delivering meaningful insights.*

No cookies. No trackers. No personal data. No compromises.

> **Perfect for privacy-first projects, sensitive applications, decentralized platforms, or any product that truly values user trust.**

---

## Why EthicalMetrics?

* **Absolute Anonymity**
  Understand your audience without ever collecting IPs, fingerprints, sessions, or cookies. Your users stay private — guaranteed.

* **Zero Third-Party Dependencies**
  No external services. No CDNs. No surprises. Full control, fully self-hosted.

* **End-to-End Encryption**
  Your data is protected with AES-256 encryption powered by SQLCipher. Privacy isn’t just a promise — it’s built in.

* **Essential, Actionable Insights**
  Track visits, module usage, errors, and load times — the core metrics you need to optimize your experience.

* **Effortless Integration**
  Lightweight, pure vanilla JavaScript snippet — zero dependencies, zero bloat.

* **Flexible Open API**
  Easily send custom events tailored to your needs.

* **Secure, Token-Based Dashboard**
  Each site gets its own private, isolated control panel with token authentication.

* **Run It Anywhere**
  Deploy on Render, VPS, Docker, or your local machine. You decide where your data lives.

---

## Getting Started

1. Sign up at `/nuevo.html` to create your site.
2. Receive your unique `site_id` and `admin_token`.
3. Add this simple script tag to your website’s HTML:

```html
<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=YOUR_SITE_ID"></script>
```

And just like that, every visit is automatically and anonymously tracked.

Want to attribute visits by page or feature? Use a custom module:

```html
<body data-modulo="home">
  ...
  <script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=YOUR_SITE_ID"></script>
</body>
```

4. Access your secure dashboard here:

```
https://ethicalmetrics.onrender.com/dashboard.html?site=YOUR_SITE_ID&token=YOUR_ADMIN_TOKEN
```

---

## Our Philosophy

EthicalMetrics is born from a simple truth: **analytics should empower, not surveil.** Inspired by tools like Overseer, OnionShare, and Signal, we believe privacy is not a feature — it’s a fundamental right.

Your users deserve transparency, respect, and control over their data. EthicalMetrics lets you measure impact *without compromise.*

---

## Experimental & Evolving

This project is currently **experimental and under active development**. We welcome contributions, audits, and feedback from the community to help us build the most trustworthy analytics platform possible.

---

## Get Involved

Created by [@livrasand](https://github.com/livrasand).
Questions, ideas, or want to contribute? Open an [issue](https://github.com/livrasand/EthicalMetrics/issues) or submit a pull request.

Together, let’s redefine what analytics can be.
