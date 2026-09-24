# Sub2API

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**AI API Gateway Platform for Subscription Quota Distribution**

English | [中文](README_CN.md)

</div>

---

## Demo

Try Sub2API online: **https://v2.pincc.ai/**

Demo credentials (shared demo environment; **not** created automatically for self-hosted installs):

| Email | Password |
|-------|----------|
| admin@sub2api.com | admin123 |

## Overview

Sub2API is an AI API gateway platform designed to distribute and manage API quotas from AI product subscriptions (like Claude Code $200/month). Users can access upstream AI services through platform-generated API Keys, while the platform handles authentication, billing, load balancing, and request forwarding.

## Features

- **Multi-Account Management** - Support multiple upstream account types (OAuth, API Key)
- **API Key Distribution** - Generate and manage API Keys for users
- **Precise Billing** - Token-level usage tracking and cost calculation
- **Smart Scheduling** - Intelligent account selection with sticky sessions
- **Concurrency Control** - Per-user and per-account concurrency limits
- **Rate Limiting** - Configurable request and token rate limits
- **Admin Dashboard** - Web interface for monitoring and management

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25.5, Gin, Ent |
| Frontend | Vue 3.4+, Vite 5+, TailwindCSS |
| Database | PostgreSQL 15+ |
| Cache/Queue | Redis 7+ |

---

## Documentation

- Dependency Security: `docs/dependency-security.md`

---

## Fork features

### Codex tickets and harvesting

- Independent ticket switches and harvest scopes for Codex models, including 292 ticket lifecycle management, validity probes, failure cooldowns, and post-recovery wakeups.
- 780 harvest flow, selected edge-IP routing, configurable harvest parameters, and per-account targeting.
- Admin workbench for harvesting, account quality operations, scheduled-test isolation, and recovery.

### OpenAI Responses and scheduling

- Improved reasoning-content and usage-detail conversion from Chat Completions to Responses.
- Compatibility handling for `function_call_output`, `previous_response_id`, and tool history.
- Account, group, and user model restrictions, including stream-only groups.
- Latest-turn admission checks, sticky scheduling, Cyber session identity isolation, and WebSocket binding.
- Excel Basispoints forwarding for OpenAI OAuth accounts.

### Administration and observability

- Smart account operations, automatic quality recovery, and Pelican test-result presentation.
- Usage views with request phases, time-to-first-token, TPS, and health interpretation.
- User-scoped announcements and `banana*` / `gemini-*image*` image model allowlists.
- Mihomo management, dynamic harvesting egress, the Codex 292 Keeper, and production deployment helpers.

Standard installation, configuration, Docker Compose, source build, and project structure are documented in the [upstream README](https://github.com/Wei-Shaw/sub2api#readme).

## Contributors

Thank you to everyone whose pull request has been merged:

<p>
  <a href="https://github.com/ranxi2001"><img src="https://avatars.githubusercontent.com/u/77790009?v=4" width="56" height="56" alt="Onefly" title="Onefly" /></a>
  <a href="https://github.com/blackdm666"><img src="https://avatars.githubusercontent.com/u/67053678?v=4" width="56" height="56" alt="老黑" title="老黑" /></a>
  <a href="https://github.com/akihitohyh"><img src="https://avatars.githubusercontent.com/u/79531840?v=4" width="56" height="56" alt="akihitohyh" title="akihitohyh" /></a>
  <a href="https://github.com/buluw"><img src="https://avatars.githubusercontent.com/u/45087912?v=4" width="56" height="56" alt="buluw" title="buluw" /></a>
  <a href="https://github.com/spake404"><img src="https://avatars.githubusercontent.com/u/123435269?v=4" width="56" height="56" alt="spake404" title="spake404" /></a>
  <a href="https://github.com/Mickey0811"><img src="https://avatars.githubusercontent.com/u/49522921?v=4" width="56" height="56" alt="Mickey0811" title="Mickey0811" /></a>
  <a href="https://github.com/mracry"><img src="https://avatars.githubusercontent.com/u/112537993?v=4" width="56" height="56" alt="mracry" title="mracry" /></a>
</p>

Merged PRs: [#1](https://github.com/ranxi2001/sub2api/pull/1), [#2](https://github.com/ranxi2001/sub2api/pull/2), [#3](https://github.com/ranxi2001/sub2api/pull/3), [#5](https://github.com/ranxi2001/sub2api/pull/5), [#7](https://github.com/ranxi2001/sub2api/pull/7), [#8](https://github.com/ranxi2001/sub2api/pull/8), [#9](https://github.com/ranxi2001/sub2api/pull/9), [#10](https://github.com/ranxi2001/sub2api/pull/10), [#12](https://github.com/ranxi2001/sub2api/pull/12), [#13](https://github.com/ranxi2001/sub2api/pull/13), [#14](https://github.com/ranxi2001/sub2api/pull/14), [#15](https://github.com/ranxi2001/sub2api/pull/15), [#17](https://github.com/ranxi2001/sub2api/pull/17), [#18](https://github.com/ranxi2001/sub2api/pull/18), [#20](https://github.com/ranxi2001/sub2api/pull/20), [#22](https://github.com/ranxi2001/sub2api/pull/22), [#23](https://github.com/ranxi2001/sub2api/pull/23), [#24](https://github.com/ranxi2001/sub2api/pull/24), [#31](https://github.com/ranxi2001/sub2api/pull/31), [#32](https://github.com/ranxi2001/sub2api/pull/32), [#33](https://github.com/ranxi2001/sub2api/pull/33), [#35](https://github.com/ranxi2001/sub2api/pull/35), [#36](https://github.com/ranxi2001/sub2api/pull/36), [#37](https://github.com/ranxi2001/sub2api/pull/37), [#38](https://github.com/ranxi2001/sub2api/pull/38), [#39](https://github.com/ranxi2001/sub2api/pull/39), [#41](https://github.com/ranxi2001/sub2api/pull/41), [#43](https://github.com/ranxi2001/sub2api/pull/43), [#45](https://github.com/ranxi2001/sub2api/pull/45), [#48](https://github.com/ranxi2001/sub2api/pull/48), [#50](https://github.com/ranxi2001/sub2api/pull/50), [#51](https://github.com/ranxi2001/sub2api/pull/51), [#52](https://github.com/ranxi2001/sub2api/pull/52), [#54](https://github.com/ranxi2001/sub2api/pull/54), [#55](https://github.com/ranxi2001/sub2api/pull/55), [#56](https://github.com/ranxi2001/sub2api/pull/56), [#57](https://github.com/ranxi2001/sub2api/pull/57), [#58](https://github.com/ranxi2001/sub2api/pull/58), [#59](https://github.com/ranxi2001/sub2api/pull/59).

---

## License

MIT License

---

<div align="center">

**If you find this project useful, please give it a star!**

</div>
