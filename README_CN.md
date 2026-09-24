# Sub2API

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**AI API 网关平台 - 订阅配额分发管理**

[English](README.md) | 中文

</div>

---

## 在线体验

体验地址：**https://v2.pincc.ai/**

演示账号（共享演示环境；自建部署不会自动创建该账号）：

| 邮箱 | 密码 |
|------|------|
| admin@sub2api.com | admin123 |

## 项目概述

Sub2API 是一个 AI API 网关平台，用于分发和管理 AI 产品订阅（如 Claude Code $200/月）的 API 配额。用户通过平台生成的 API Key 调用上游 AI 服务，平台负责鉴权、计费、负载均衡和请求转发。

## 核心功能

- **多账号管理** - 支持多种上游账号类型（OAuth、API Key）
- **API Key 分发** - 为用户生成和管理 API Key
- **精确计费** - Token 级别的用量追踪和成本计算
- **智能调度** - 智能账号选择，支持粘性会话
- **并发控制** - 用户级和账号级并发限制
- **速率限制** - 可配置的请求和 Token 速率限制
- **管理后台** - Web 界面进行监控和管理

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go 1.25.5, Gin, Ent |
| 前端 | Vue 3.4+, Vite 5+, TailwindCSS |
| 数据库 | PostgreSQL 15+ |
| 缓存/队列 | Redis 7+ |

---

## 文档

- 依赖安全：`docs/dependency-security.md`

---

## OpenAI Responses 兼容注意事项

- 当请求包含 `function_call_output` 时，需要携带 `previous_response_id`，或在 `input` 中包含带 `call_id` 的 `tool_call`/`function_call`，或带非空 `id` 且与 `function_call_output.call_id` 匹配的 `item_reference`。
- 若依赖上游历史记录，网关会强制 `store=true` 并需要复用 `previous_response_id`，以避免出现 “No tool call found for function call output” 错误。

---

## 本 fork 二开特性

### Codex ticket 与采票

- 独立控制 Codex 模型 ticket 开关与采集范围，支持 292 ticket 生命周期、有效性探测、失败冷却和账号恢复后的后台唤醒。
- 支持 780 采票链路、指定边缘 IP 直拨、可配置打票参数和按账号定向打票。
- 管理后台提供采票工作台、账号质量运维、定时测试隔离与恢复能力。

### OpenAI Responses 与账号调度

- 完善 Chat Completions 到 Responses 的推理内容和用量明细转换。
- 兼容 `function_call_output`、`previous_response_id` 和工具历史。
- 支持按账号、分组和用户限制模型，以及分组仅允许流式请求。
- 增加最新会话准入校验、粘性调度、Cyber 会话身份隔离和 WebSocket 绑定。
- 支持 Excel Basispoints 协议，使用 OpenAI OAuth 账号转发 Responses 请求。

### 管理与观测

- 支持智能账号操作、质量自动恢复和 Pelican 测智结果展示。
- 使用记录增加请求耗时分段、首字延迟、TPS 和健康状态说明。
- 公告可按用户可见，图片模型白名单支持 `banana*` 与 `gemini-*image*`。
- 保留 Mihomo 管理、动态打票出口、Codex 292 Keeper 和生产部署脚本等配套能力。

基础安装、配置、Docker Compose、源码编译和项目结构请参考[上游 README](https://github.com/Wei-Shaw/sub2api#readme)。

## 贡献者

感谢所有已合并 PR 的贡献者：

<p>
  <a href="https://github.com/ranxi2001"><img src="https://avatars.githubusercontent.com/u/77790009?v=4" width="56" height="56" alt="Onefly" title="Onefly" /></a>
  <a href="https://github.com/blackdm666"><img src="https://avatars.githubusercontent.com/u/67053678?v=4" width="56" height="56" alt="老黑" title="老黑" /></a>
  <a href="https://github.com/akihitohyh"><img src="https://avatars.githubusercontent.com/u/79531840?v=4" width="56" height="56" alt="akihitohyh" title="akihitohyh" /></a>
  <a href="https://github.com/buluw"><img src="https://avatars.githubusercontent.com/u/45087912?v=4" width="56" height="56" alt="buluw" title="buluw" /></a>
  <a href="https://github.com/spake404"><img src="https://avatars.githubusercontent.com/u/123435269?v=4" width="56" height="56" alt="spake404" title="spake404" /></a>
  <a href="https://github.com/Mickey0811"><img src="https://avatars.githubusercontent.com/u/49522921?v=4" width="56" height="56" alt="Mickey0811" title="Mickey0811" /></a>
  <a href="https://github.com/mracry"><img src="https://avatars.githubusercontent.com/u/112537993?v=4" width="56" height="56" alt="mracry" title="mracry" /></a>
</p>

已合并 PR：[#1](https://github.com/ranxi2001/sub2api/pull/1)、[#2](https://github.com/ranxi2001/sub2api/pull/2)、[#3](https://github.com/ranxi2001/sub2api/pull/3)、[#5](https://github.com/ranxi2001/sub2api/pull/5)、[#7](https://github.com/ranxi2001/sub2api/pull/7)、[#8](https://github.com/ranxi2001/sub2api/pull/8)、[#9](https://github.com/ranxi2001/sub2api/pull/9)、[#10](https://github.com/ranxi2001/sub2api/pull/10)、[#12](https://github.com/ranxi2001/sub2api/pull/12)、[#13](https://github.com/ranxi2001/sub2api/pull/13)、[#14](https://github.com/ranxi2001/sub2api/pull/14)、[#15](https://github.com/ranxi2001/sub2api/pull/15)、[#17](https://github.com/ranxi2001/sub2api/pull/17)、[#18](https://github.com/ranxi2001/sub2api/pull/18)、[#20](https://github.com/ranxi2001/sub2api/pull/20)、[#22](https://github.com/ranxi2001/sub2api/pull/22)、[#23](https://github.com/ranxi2001/sub2api/pull/23)、[#24](https://github.com/ranxi2001/sub2api/pull/24)、[#31](https://github.com/ranxi2001/sub2api/pull/31)、[#32](https://github.com/ranxi2001/sub2api/pull/32)、[#33](https://github.com/ranxi2001/sub2api/pull/33)、[#35](https://github.com/ranxi2001/sub2api/pull/35)、[#36](https://github.com/ranxi2001/sub2api/pull/36)、[#37](https://github.com/ranxi2001/sub2api/pull/37)、[#38](https://github.com/ranxi2001/sub2api/pull/38)、[#39](https://github.com/ranxi2001/sub2api/pull/39)、[#41](https://github.com/ranxi2001/sub2api/pull/41)、[#43](https://github.com/ranxi2001/sub2api/pull/43)、[#45](https://github.com/ranxi2001/sub2api/pull/45)、[#48](https://github.com/ranxi2001/sub2api/pull/48)、[#50](https://github.com/ranxi2001/sub2api/pull/50)、[#51](https://github.com/ranxi2001/sub2api/pull/51)、[#52](https://github.com/ranxi2001/sub2api/pull/52)、[#54](https://github.com/ranxi2001/sub2api/pull/54)、[#55](https://github.com/ranxi2001/sub2api/pull/55)、[#56](https://github.com/ranxi2001/sub2api/pull/56)、[#57](https://github.com/ranxi2001/sub2api/pull/57)、[#58](https://github.com/ranxi2001/sub2api/pull/58)、[#59](https://github.com/ranxi2001/sub2api/pull/59)。

---

## 许可证

MIT License

---

<div align="center">

**如果觉得有用，请给个 Star 支持一下！**

</div>
