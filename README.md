# 🔮 占卜网 (ZhanBu)

一个现代化的在线占卜平台，支持塔罗牌、星座运势、六爻、八字等多种占卜方式，集成 AI 解读功能。

## ✨ 功能特性

- **塔罗牌占卜** - 支持单牌、三牌、凯尔特十字、爱情牌阵等多种牌阵
- **星座运势** - 12 星座日/周/月运势查询
- **六爻占卜** - 传统铜钱六爻起卦，自动解析卦象
- **八字排盘** - 根据出生日期时间自动排盘，分析五行、十神
- **AI 解读** - 集成 OpenAI 兼容 API，提供智能解读（支持流式输出）
- **用户系统** - 注册、登录、个人资料管理
- **历史记录** - 保存占卜记录，随时回顾
- **邮箱验证** - 支持邮箱验证码注册

## 🏗️ 技术栈

- **后端**: Go + Gin + GORM + PostgreSQL
- **前端**: React + TypeScript + Vite + Tailwind CSS
- **数据库**: PostgreSQL 16
- **部署**: Docker + Docker Compose + Nginx

## 🚀 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose v2.0+

### 1. 克隆项目

```bash
git clone <repo-url>
cd zhanbu
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，修改 JWT_SECRET 和数据库密码
```

### 3. 启动服务

```bash
docker-compose up -d
```

### 4. 访问应用

- **前端**: http://localhost
- **API**: http://localhost:18080
- **健康检查**: http://localhost:18080/api/health

## 📁 项目结构

```
zhanbu/
├── docker-compose.yml          # Docker Compose 编排
├── .env.example                # 环境变量模板
├── nginx/
│   └── nginx.conf              # Nginx 配置（顶层备用）
├── client/                     # 前端项目
│   ├── Dockerfile
│   ├── nginx.conf              # 前端 Nginx 配置
│   └── src/
└── server/                     # 后端项目
    ├── Dockerfile
    ├── main.go
    ├── config/
    │   ├── config.go           # 配置加载
    │   └── config.yaml         # 默认配置
    ├── data/                   # 静态数据（塔罗牌、卦象等）
    ├── internal/
    │   ├── database/           # 数据库初始化、迁移、种子
    │   ├── handler/            # HTTP 处理器
    │   ├── middleware/          # 中间件（CORS、Auth、RateLimit）
    │   ├── model/              # 数据模型
    │   ├── repository/         # 数据访问层
    │   ├── router/             # 路由配置
    │   └── service/            # 业务逻辑层
    └── pkg/                    # 公共工具包
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_USER` | PostgreSQL 用户名 | `zhanbu` |
| `DB_PASSWORD` | PostgreSQL 密码 | `zhanbu_secret` |
| `DB_NAME` | 数据库名 | `zhanbu` |
| `JWT_SECRET` | JWT 签名密钥 | `change-this-to-a-secure-secret` |
| `AI_API_KEY` | OpenAI API Key | 空 |
| `AI_MODEL` | AI 模型名称 | `gpt-4` |
| `AI_BASE_URL` | AI API 地址 | `https://api.openai.com/v1` |
| `PORT` | Render 注入的后端监听端口 | `18080` |
| `ZHANBU_SERVER_PORT` | 后端监听端口，优先级高于 `PORT` | `18080` |
| `ZHANBU_BCRYPT_COST` | 注册/登录密码哈希成本；Render 低配实例可设为 `8` 或 `10` | `10` |
| `ZHANBU_CORS_ALLOWED_ORIGINS` | 允许访问后端的前端域名，逗号分隔 | `http://localhost:5173` |
| `VITE_API_BASE_URL` | Vercel 前端调用的 API 前缀，例如 `https://your-api.onrender.com/api` | `/api` |

### Render + Vercel 部署要点

- Render 后端使用 `server/Dockerfile`，数据库请使用 Render PostgreSQL 或 Neon/Supabase 等外部 PostgreSQL，不要把 Postgres 放进同一个 Web Service 容器。
- Render 后端环境变量至少设置：`ZHANBU_SERVER_MODE=release`、`ZHANBU_JWT_SECRET`、`ZHANBU_DB_HOST`、`ZHANBU_DB_PORT`、`ZHANBU_DB_USER`、`ZHANBU_DB_PASSWORD`、`ZHANBU_DB_NAME`、`ZHANBU_DB_SSLMODE=require`、`ZHANBU_CORS_ALLOWED_ORIGINS=https://your-app.vercel.app`、`ZHANBU_AI_API_KEY`、`ZHANBU_AI_MODEL`、`ZHANBU_AI_BASE_URL`。
- Vercel 前端项目根目录选 `client`，构建命令 `npm run build`，输出目录 `dist`，环境变量设置 `VITE_API_BASE_URL=https://your-render-service.onrender.com/api`。
- 本地开发不设置 `VITE_API_BASE_URL` 时，前端继续请求 `/api`，由 Vite 代理到 `http://localhost:18080`；Docker Compose 版继续由 Nginx 代理到 `api:8080`。

### 数据持久化

- **PostgreSQL 数据**: Docker Volume `postgres-data`，存储在 `/var/lib/postgresql/data`
- **静态资源**: Docker Volume `api-data`，存储塔罗牌数据等

## 🛠️ 开发模式

### 本地开发（不使用 Docker）

```bash
# 1. 启动 PostgreSQL（可使用 Docker）
docker run -d --name zhanbu-pg \
  -e POSTGRES_USER=zhanbu \
  -e POSTGRES_PASSWORD=zhanbu_secret \
  -e POSTGRES_DB=zhanbu \
  -p 5432:5432 \
  postgres:16-alpine

# 2. 启动后端
cd server
go run main.go

# 3. 启动前端
cd client
npm install
npm run dev
```

### 运行测试

```bash
cd server
go test ./...
```

> 测试使用 SQLite 内存数据库，无需 PostgreSQL 实例。

## 📝 API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/health` | 健康检查 | ❌ |
| POST | `/api/auth/register` | 用户注册 | ❌ |
| POST | `/api/auth/login` | 用户登录 | ❌ |
| GET | `/api/auth/profile` | 获取个人信息 | ✅ |
| PUT | `/api/auth/profile` | 更新个人信息 | ✅ |
| GET | `/api/tarot/cards` | 获取所有塔罗牌 | ❌ |
| GET | `/api/tarot/spreads` | 获取牌阵列表 | ❌ |
| POST | `/api/tarot/draw` | 抽牌 | ✅ |
| GET | `/api/horoscope/:zodiac` | 星座运势 | ❌ |
| POST | `/api/liuyao/throw` | 六爻起卦 | ✅ |
| POST | `/api/bazi/calculate` | 八字排盘 | ✅ |
| GET | `/api/history` | 历史记录列表 | ✅ |
| GET | `/api/history/:id` | 历史记录详情 | ✅ |
| DELETE | `/api/history/:id` | 删除历史记录 | ✅ |
| POST | `/api/ai/interpret` | AI 解读 | ✅ |

## 🔒 安全说明

- **JWT 密钥**: 生产环境必须通过 `JWT_SECRET` 环境变量设置强密码
- **数据库密码**: 生产环境必须修改 `DB_PASSWORD`
- **HTTPS**: 建议在生产环境配置 SSL/TLS（可在 `nginx.conf` 中取消注释 HTTPS 配置）

## 📄 License

MIT
