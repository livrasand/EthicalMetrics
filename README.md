# 📊 EthicalMetrics

**EthicalMetrics** is a web analytics platform that’s ethical, private, and self-hosted. No cookies, no tracking, no user identification, no third parties. Designed to respect privacy by default.

> 💡 Ideal for sensitive projects, private apps, decentralized platforms, or any product that truly values user privacy.

---

## 🚀 Features

- 🔒 **100% anonymous**: no IP tracking, no fingerprints, no sessions, no cookies.
- 🛡️ **No third parties**: doesn’t rely on external services or CDNs.
- 🔐 **SQLCipher encryption**: AES-256 encrypted database.
- 📈 **Basic analytics**: visits, modules used, errors, load times.
- 🧩 **Universal embeddable script**: pure vanilla JS, zero dependencies.
- 🔧 **Open API**: send custom events with ease.
- 🔐 **Token-based dashboard**: each site has a private, isolated panel.
- 🛠️ **Self-hosted**: run it on Render, VPS, Docker, or locally.

---

## 🧑‍💻 How it works

1. The user registers at `/nuevo.html`
2. Receives a `site_id` and `admin_token`
3. Adds the following `<script>` to their website:

```html
<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js"
        defer data-site-id="YOUR_SITE_ID">
</script>
```

¡Listo! Con eso, se registra automáticamente un evento "visita".

Opción con módulo personalizado:

```html
<body data-modulo="home">
  ...
  <script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js" defer></script>
</body>
```

Así el evento "visita" se atribuye automáticamente al módulo home.

4. Access the dashboard at:

```
https://ethicalmetrics.onrender.com/dashboard.html?site=YOUR_SITE_ID&token=YOUR_ADMIN_TOKEN
```

---

## 🧠 Philosophy

EthicalMetrics was born from a real need: measure without surveillance. It is inspired by ethical tools like Overseer, OnionShare, and Signal. We believe in a web without surveillance.

> ✊ Privacy is not a feature. It’s a right.

---

## 📬 Contact

Project created by [@livrasand](https://github.com/livrasand)
Have questions or ideas? Open an [issue](https://github.com/livrasand/EthicalMetrics/issues) or contribute!

