# 📧 `emailengine`

`emailengine` is a lightweight email framework I built for personal projects.  
It provides both a REST API and an SMTP server for flexible, middleware-driven email handling.

## Features

- **REST API** — Easily send outbound emails via [documented endpoints](API.md).  
- **SMTP Server** — Accept and filter incoming mail through customizable middleware.  
- **Examples Included** — Check out:
  - [Client with Go Templates](examples/client/main.go)
  - [Server with REST API Setup](examples/server/main.go)

Each example is fully commented to help you get started quickly.

## Getting Started
A full setup and deployment guide is available on my [blog](https://panca.kz/goto/emailengine).  
It walks through the entire process of getting `emailengine` running in production.