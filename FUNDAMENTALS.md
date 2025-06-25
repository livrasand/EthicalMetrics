# 🧱 EthicalMetrics Fundamentals

**EthicalMetrics** is a web analytics tool designed from the ground up to respect people's true privacy. This document establishes the principles that define our development and usage ethics.

## Fundamentals you must NOT compromise

These are the fundamental pillars for maintaining the ethical essence of the project. No feature or implementation must violate these rules:

### 1. No fingerprinting
No type of identification technique based on:
- User-Agent
- IP address
- Screen resolution
- Operating system
- Browser font or language
- Any combination of the above is allowed.

> **Generating hashes or identifiers derived from these attributes is prohibited.

### 2. No cookies or localStorage
EthicalMetrics does not use:
- First-party or third-party cookies
- localStorage, sessionStorage, or IndexedDB

> **The browser should not store data for tracking.**

### 3. No IP storage
IP addresses are not saved or processed for targeting or analytics purposes.

> **Not even partial anonymization like in GA or Matomo. They are not stored. Period.**

---

## Fundamentals you MUST follow

### 4. True anonymity
EthicalMetrics does not assign persistent identifiers to visitors:
- There is no `visitor_id`
- Session events are not associated with each other
- Each event is autonomous and cannot be tracked over time

> **Analytics must work without knowing who you are.**

### 5. Automatic DNT compliance
If the user's browser has **Do Not Track** enabled, EthicalMetrics **automatically respects that choice** and disables all tracking.

> **Full respect for implicit user consent.**

### 6. Minimal and non-invasive JavaScript
JS tracking code should be:
- Lightweight and easy to audit (~1KB if possible)
- No opaque logic or confusing minification
- Easy to disable or adapt for strict environments

> **No dependencies on external libraries or complex JS.**

### 7. No cross-site tracking
Not allowed:
- Correlating activity between domains
- Use of CNAME cloaking or disguised redirects
- Shared identifiers between sites

> **Each domain is independent. The user cannot be tracked between them.**

---

## Conclusion

EthicalMetrics isn't just analytics software. It's a statement of principles.

If you can't implement a metric without compromising one of these pillars, **then that metric shouldn't exist here**.

---

🛡️ *Ethical by design. Anonymous by default. Metrics with ethics. Do EthicalMetrics!*
