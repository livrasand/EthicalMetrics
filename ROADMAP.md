# Roadmap: Functionalities and Objectives

🛡️ **Ethical by design. Anonymous by default. Metrics with ethics. Do EthicalMetrics!**

---

## 1. Basic metrics (without compromising privacy)

- [x] Visits and views per page
- [x] Referrers
- [x] Load time
- [x] Most visited pages
- [x] Average session duration
- [x] Custom events
- [x] Bounce rate (per page and global)
- [x] Devices, browsers, and operating systems
- [x] Browser languages
- [x] Location (city, region, country) with MaxMind GeoLite2 (no IP address storage)
- [x] UTM tracking (`utm_source`, `utm_medium`, `utm_campaign`, etc.)
- [x] Time-based comparisons (week vs. week, month vs. month)
- [x] Retention (no persistent IDs)
- [x] Basic funnels (event progression or pages)

---

## 2. Advanced Metrics and Ethical Segmentation

- [ ] Filters by country, browser, URL, device, source, etc.
- [ ] Advanced segmentation (multiple criteria without individual identification)
- [ ] Custom dimensions and variables without identifying users
- [ ] Journey visualization (anonymous sequential navigation)
- [ ] Ethical attribution of multi-channel conversions/events
- [ ] User-defined goals, with configurable conditions
- [ ] Time cohorts without unique ID
- [ ] Row evolution (time comparison by row)
- [ ] Rich events: clicks, downloads, scrolling, banners, 404 errors
- [ ] Exact duration per page
- [ ] Entry and exit pages
- [ ] Internal site search
- [ ] CTR per element (text, images, buttons)
- [ ] Form tracking (anonymous, without user identification)

---

## 3. Ecommerce and ethical conversions

- [ ] Tracking Ecommerce: products, carts, purchases, revenue (without identifying users)
- [ ] Multi-channel conversions (ads, organic, direct, referrers)
- [ ] Monetary value of events
- [ ] Funnel visualization
- [ ] Privacy-respecting UTM campaigns and ad tracking

---

## 4. User Experience (UX/UI)

- [ ] SPA (Single Page App)
- [ ] Responsive interface with light/dark mode
- [ ] Default overview view with key metrics
- [ ] Real-time dashboard
- [ ] Custom panels (drag & drop)
- [ ] Intuitive custom segments/filters
- [ ] Multi-site dashboard
- [ ] Visual user flow (anonymous flow)
- [ ] Page transitions
- [ ] Page overlay (site data)
- [ ] Chart annotations
- [ ] Alerts Customizable

---

## 5. Visual Features (Optional)

*(Always with consent and without hidden tracking)*

- [ ] Heatmaps and scrollmaps (without fingerprinting)
- [ ] Session recording (ethical version without identification)
- [ ] Native A/B Testing
- [ ] Form analysis (abandonment, fields)
- [ ] Media analytics (videos/audio)

---

## 6. Backend and Infrastructure

- [ ] Optimized database
- [ ] Asynchronous loading (WebSocket or polling)
- [ ] Efficient processing (stream/batch)
- [ ] Secure data import/export
- [ ] Public API (REST or GraphQL) without exposing sensitive data
- [ ] Multi-site support (roll-up analytics)
- [ ] User control per site (with roles)
- [ ] Secure auditing (without IPs or (personal identifiers)
- [ ] Log tracking (for intranets)
- [ ] Multi-language interface
- [ ] Scalable for high traffic
- [ ] White labeling

---

## 7. Integrations and extensibility

- [ ] Modular plugin/extension system
- [ ] Public administration/tracking API
- [ ] Official SDKs:

- JavaScript
- PHP
- Python
- Android / iOS

- [ ] CMS plugins (WordPress, Joomla, Drupal)
- [ ] Modern frameworks: Vue, React, Nuxt, Next, Astro, etc.
- [ ] Offline tracking (deferred synchronization)
- [ ] Log tracking (Apache, Nginx, IIS)

---

## 8. Privacy and legal compliance (non-negotiable)

- [x] No cookies
- [x] No stored IPs
- [x] Self-hosted
- [x] SQLCipher compatible
- [ ] Automatic DNT (Do Not Track) compliance
- [ ] GDPR, CCPA, PECR certifications
- [ ] Legal tools ready (DPA, policies, terms)
- [ ] Opt-out form (embeddable iFrame)
- [ ] Granular consent by category
- [ ] Referrer and geodata anonymization
- [ ] ARCO panel: access, rectification, cancellation, opposition
- [ ] Data deletion under the right to be forgotten
- [ ] Cookie control (if the user opts in) (use them)
- [ ] User ID replaced with ephemeral pseudonyms

---

## 9. Analyst and Administration Tools

- [ ] Data export: JSON, XML, CSV, Excel
- [ ] Automatic reports (PDF, HTML, PNG)
- [ ] Configurable time zones per site
- [ ] Exclusion of IP addresses or ranges
- [ ] Exclusion of URL parameters
- [ ] Embeddable dashboard
- [ ] Official mobile app (ethical, no tracking)
- [ ] Multi-currency support
- [ ] Management of multiple users and sites
- [ ] Unified configuration panel

---

> _EthicalMetrics doesn't track people, it tracks context._
