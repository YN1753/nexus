# Nexus

##  项目简介

Nexus 是一个面向开发者（Developer）和 HomeLab 用户的现代化 Linux 服务器管理平台。

项目通过 **SSH** 对 Linux 服务器进行统一管理，无需在目标服务器安装 Agent。

Nexus 不追求功能堆砌，而是专注于开发者日常真正会使用的功能，提供简洁、高效的服务器管理体验。

---

## 项目特点

-  基于 Go + Vue3 构建
-  基于 SSH 管理 Linux 服务器
-  Docker 容器管理
-  在线终端（Web Terminal）
-  文件管理（SFTP）
-  系统监控（CPU、内存、磁盘等）
-  常用脚本管理
-  Local First
-  面向开发者设计

---

## Roadmap

### V0.1

- [x] 用户登录
- [ ] 节点（Node）管理
- [ ] SSH 连接
- [ ] 系统信息
- [ ] Dashboard
- [ ] 在线终端
- [ ] 文件管理
- [ ] Docker 管理
- [ ] Script 管理

### 未来

- [ ] Docker Compose
- [ ] Cron 管理
- [ ] 一键部署
- [ ] 对象存储接入
- [ ] 多节点管理优化
- [ ] Cloud Sync

---

##  技术栈

### Backend

- Go
- Gin
- GORM
- MySQL
- Redis
- JWT
- Zap
- Viper

### Frontend

- Vue3
- TypeScript
- Pinia
- Vue Router

### Infrastructure

- Docker
- Docker Compose
- SSH
- SFTP
- WebSocket

---

## 项目结构

```
server/
web/
README.md
```

---

