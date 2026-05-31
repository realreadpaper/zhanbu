# 占卜网站 — 开发计划文档

> 版本：v1.0 | 最后更新：2026-05-30
> 基于：[design.md](./design.md)

---

## 目录

1. [里程碑总览](#1-里程碑总览)
2. [Phase 1 详细任务（核心 MVP）](#2-phase-1-详细任务核心-mvp)
3. [Phase 2 详细任务（功能增强）](#3-phase-2-详细任务功能增强)
4. [Phase 3 详细任务（上线准备）](#4-phase-3-详细任务上线准备)
5. [技术风险与应对](#5-技术风险与应对)
6. [开发规范](#6-开发规范)
7. [本地开发环境搭建步骤](#7-本地开发环境搭建步骤)

---

## 1. 里程碑总览

### 1.1 整体时间线

```
Phase 1 (4周)              Phase 2 (4周)              Phase 3 (3周)
基础框架与核心 MVP          功能增强                    上线准备
├─ 项目初始化               ├─ 星座运势                ├─ AI 解读集成
├─ 用户认证                 ├─ 周易六爻                ├─ UI/UX 优化
├─ 塔罗牌占卜               └─ 八字排盘                ├─ 移动端适配
├─ 历史记录                                          ├─ 生产部署
└─ 前后端联调                                          └─ 正式上线
```

| 阶段 | 目标 | 时间估算 | 核心交付物 |
|------|------|---------|-----------|
| **Phase 1** | 核心 MVP：搭建框架 + 塔罗牌占卜全流程 | 4 周 | 可运行的塔罗牌占卜网站（含用户系统、4 种牌阵、抽牌动画、历史记录） |
| **Phase 2** | 功能增强：补齐其余三种占卜类型 | 4 周 | 星座运势 + 周易六爻 + 八字排盘，四种占卜类型全部可用 |
| **Phase 3** | 上线准备：AI 集成 + 优化 + 部署 | 3 周 | AI 大模型解读、移动端适配、生产环境部署、正式上线 |

### 1.2 里程碑节点

| 里程碑 | 时间点 | 验收标准 |
|--------|--------|---------|
| M1 — 项目骨架可运行 | Phase 1 Week 1 结束 | 前后端均可启动，用户可注册/登录 |
| M2 — 塔罗牌占卜可体验 | Phase 1 Week 3 结束 | 塔罗牌占卜完整流程可走通（选牌阵→抽牌→看解读） |
| M3 — MVP 完成 | Phase 1 结束 | 塔罗牌占卜 + 历史记录，可部署测试 |
| M4 — 星座运势上线 | Phase 2 Week 6 结束 | 12 星座日/周/月运势可查 |
| M5 — 全类型占卜完成 | Phase 2 结束 | 六爻 + 八字排盘可用，四种类型全部就绪 |
| M6 — AI 解读可用 | Phase 3 Week 10 结束 | 四种占卜类型均可触发 AI 流式解读 |
| M7 — 正式上线 | Phase 3 结束 | 生产环境部署完成，监控就绪 |

---

## 2. Phase 1 详细任务（核心 MVP）

### Week 1：项目初始化与用户系统

**目标**：搭建前后端项目骨架，完成用户认证全流程。

#### 任务 1.1：初始化 Go 后端项目

- **任务描述**：创建 `server/` 目录，初始化 Go Module，安装核心依赖（Gin、GORM、SQLite 驱动、JWT 库、zerolog、validator），编写 `main.go` 入口文件，实现基础 HTTP 服务启动。
- **预估耗时**：4 小时
- **前置依赖**：无
- **产出物**：
  - `server/go.mod`、`server/go.sum`
  - `server/main.go`（可启动的 HTTP 服务，监听 8080 端口）
  - `server/config/config.go`（配置结构体定义）
  - `server/config/config.yaml`（默认配置文件）

#### 任务 1.2：实现配置管理模块

- **任务描述**：使用 `viper` 或手动解析 `config.yaml`，加载服务器端口、数据库路径、JWT 密钥、AI 配置等。支持环境变量覆盖。
- **预估耗时**：3 小时
- **前置依赖**：任务 1.1
- **产出物**：
  - `server/config/config.go`（完整的配置加载逻辑）

#### 任务 1.3：实现数据库初始化与迁移

- **任务描述**：配置 GORM 连接 SQLite，实现自动迁移（AutoMigrate）创建 `users` 和 `divination_records` 表。编写 `database/db.go` 和 `database/migrate.go`。
- **预估耗时**：3 小时
- **前置依赖**：任务 1.2
- **产出物**：
  - `server/internal/database/db.go`（数据库连接初始化）
  - `server/internal/database/migrate.go`（自动迁移逻辑）
  - `server/data/zhanbu.db`（SQLite 数据库文件）

#### 任务 1.4：实现用户模型与数据访问层

- **任务描述**：定义 `User` 模型（对应 `users` 表），实现 Repository 层的 CRUD 方法（Create、FindByEmail、FindByUsername、FindByID、Update）。
- **预估耗时**：3 小时
- **前置依赖**：任务 1.3
- **产出物**：
  - `server/internal/model/user.go`
  - `server/internal/repository/user.go`

#### 任务 1.5：实现 JWT 认证与密码工具

- **任务描述**：实现 JWT 工具函数（生成 Access Token / Refresh Token、解析验证），实现 bcrypt 密码哈希工具。编写认证中间件，从 Header 提取 Token 并注入用户上下文。
- **预估耗时**：4 小时
- **前置依赖**：任务 1.4
- **产出物**：
  - `server/pkg/utils/jwt.go`
  - `server/pkg/utils/hash.go`
  - `server/internal/middleware/auth.go`

#### 任务 1.6：实现统一响应与错误处理

- **任务描述**：定义统一 JSON 响应格式 `{ code, data, message }`，定义业务错误码常量，实现成功/失败响应工具函数。
- **预估耗时**：2 小时
- **前置依赖**：无
- **产出物**：
  - `server/pkg/response/response.go`
  - `server/pkg/errors/errors.go`

#### 任务 1.7：实现认证 API（注册/登录/刷新/用户信息）

- **任务描述**：编写 `auth` Handler 和 Service，实现以下接口：
  - `POST /api/auth/register` — 注册（校验用户名/邮箱唯一性、密码强度）
  - `POST /api/auth/login` — 登录（返回 access_token + refresh_token）
  - `POST /api/auth/refresh` — 刷新 Token
  - `GET /api/auth/profile` — 获取当前用户信息
  - `PUT /api/auth/profile` — 更新用户信息
- **预估耗时**：5 小时
- **前置依赖**：任务 1.5、1.6
- **产出物**：
  - `server/internal/handler/auth.go`
  - `server/internal/service/auth.go`
  - `server/internal/router/router.go`（路由注册）

#### 任务 1.8：实现 CORS、日志、限流中间件

- **任务描述**：配置 CORS 中间件（允许前端开发域名 `localhost:5173`），接入 zerolog 请求日志中间件，实现简单的内存限流中间件。
- **预估耗时**：3 小时
- **前置依赖**：任务 1.7
- **产出物**：
  - `server/internal/middleware/cors.go`
  - `server/internal/middleware/logger.go`
  - `server/internal/middleware/ratelimit.go`

#### 任务 1.9：初始化 React 前端项目

- **任务描述**：使用 Vite 创建 React + TypeScript 项目，安装 TailwindCSS、React Router、Zustand、Axios、Framer Motion、React Query。配置 Tailwind、路径别名、代理（`/api` → `localhost:8080`）。
- **预估耗时**：3 小时
- **前置依赖**：无
- **产出物**：
  - `client/package.json`、`client/vite.config.ts`、`client/tailwind.config.js`
  - `client/src/main.tsx`、`client/src/App.tsx`
  - 可运行的前端开发服务器（`localhost:5173`）

#### 任务 1.10：实现前端基础布局与路由

- **任务描述**：搭建 `Header`（Logo、导航栏、登录/用户头像）、`Footer` 组件。配置 React Router 路由表（`/`、`/tarot`、`/horoscope`、`/liuyao`、`/bazi`、`/history`、`/profile`、`/login`、`/register`）。实现 `ErrorBoundary`、`Loading` 通用组件。
- **预估耗时**：4 小时
- **前置依赖**：任务 1.9
- **产出物**：
  - `client/src/components/common/Header.tsx`
  - `client/src/components/common/Footer.tsx`
  - `client/src/components/common/Loading.tsx`
  - `client/src/components/common/ErrorBoundary.tsx`
  - `client/src/App.tsx`（路由配置）

#### 任务 1.11：实现前端认证页面与 API 对接

- **任务描述**：实现登录页（`Login.tsx`）和注册页（`Register.tsx`），包含表单验证、错误提示。实现 Axios 实例配置（baseURL、Token 拦截器、自动刷新）。实现 `useAuth` Hook 和 `authStore`（Zustand）管理认证状态。实现路由守卫（未登录跳转登录页）。
- **预估耗时**：5 小时
- **前置依赖**：任务 1.10、1.7
- **产出物**：
  - `client/src/pages/Login.tsx`
  - `client/src/pages/Register.tsx`
  - `client/src/services/api.ts`（Axios 实例）
  - `client/src/services/auth.ts`
  - `client/src/hooks/useAuth.ts`
  - `client/src/stores/authStore.ts`

#### 任务 1.12：后端认证模块单元测试

- **任务描述**：为认证相关的 Service 和工具函数编写单元测试，覆盖注册、登录、Token 刷新、密码哈希等核心逻辑。
- **预估耗时**：3 小时
- **前置依赖**：任务 1.7
- **产出物**：
  - `server/internal/service/auth_test.go`
  - `server/pkg/utils/jwt_test.go`
  - `server/pkg/utils/hash_test.go`

---

### Week 2：塔罗牌后端核心

**目标**：完成塔罗牌数据准备、洗牌算法、牌阵逻辑、全部 API。

#### 任务 2.1：准备塔罗牌数据（78 张）

- **任务描述**：编写完整的 78 张塔罗牌 JSON 数据，包括：22 张大阿尔卡纳 + 56 张小阿尔卡纳。每张牌包含 id、中英文名、类型（major/minor）、花色、编号、正逆位关键词、正逆位含义、牌面描述。
- **预估耗时**：6 小时（数据量大，需要参考经典资料）
- **前置依赖**：无
- **产出物**：
  - `server/data/tarot/major_arcana.json`（22 张大阿尔卡纳）
  - `server/data/tarot/minor_arcana.json`（56 张小阿尔卡纳）

#### 任务 2.2：实现塔罗牌数据模型与种子数据导入

- **任务描述**：定义 `TarotCard` GORM 模型，编写种子数据导入脚本（`seed.go`），将 JSON 数据导入 SQLite。实现 Repository 层（FindByID、FindAll、FindByType、FindBySuit）。
- **预估耗时**：3 小时
- **前置依赖**：任务 2.1、1.3
- **产出物**：
  - `server/internal/model/tarot.go`
  - `server/internal/repository/tarot.go`
  - `server/internal/database/seed.go`（种子数据导入逻辑）

#### 任务 2.3：实现 Fisher-Yates 洗牌与抽牌算法

- **任务描述**：在 Service 层实现核心占卜算法：
  - Fisher-Yates 洗牌（保证均匀随机）
  - 按牌阵类型抽取指定数量的牌
  - 每张牌随机决定正位/逆位
  - 返回完整的抽牌结果（包含牌面数据 + 位置信息 + 正逆位）
- **预估耗时**：4 小时
- **前置依赖**：任务 2.2
- **产出物**：
  - `server/internal/service/tarot.go`

#### 任务 2.4：实现 4 种牌阵逻辑

- **任务描述**：在 Service 层实现 4 种牌阵：
  - 单牌抽取（single）— 1 张
  - 三牌阵（three）— 过去/现在/未来
  - 凯尔特十字阵（celtic）— 10 张，各位置含义
  - 爱情十字阵（love）— 5 张，各位置含义
  - 每种牌阵定义位置名称和含义映射。
- **预估耗时**：4 小时
- **前置依赖**：任务 2.3
- **产出物**：
  - `server/internal/service/tarot.go`（牌阵逻辑部分）

#### 任务 2.5：实现塔罗牌 API Handler 与路由

- **任务描述**：编写 Handler 层，实现以下接口：
  - `GET /api/tarot/cards` — 获取所有牌列表（支持分页）
  - `GET /api/tarot/cards/:id` — 获取单张牌详情
  - `GET /api/tarot/spreads` — 获取可用牌阵列表（返回牌阵名称、张数、位置说明）
  - `POST /api/tarot/draw` — 抽牌（请求体：spread 类型 + 可选 question）
- **预估耗时**：3 小时
- **前置依赖**：任务 2.4、1.7
- **产出物**：
  - `server/internal/handler/tarot.go`
  - 路由注册到 `server/internal/router/router.go`

#### 任务 2.6：塔罗牌模块单元测试

- **任务描述**：编写单元测试覆盖：
  - 洗牌算法的随机性与正确性（所有牌都能被抽到、不重复）
  - 4 种牌阵的抽牌数量和位置映射
  - 正逆位的概率分布（大量抽样统计）
  - API 请求参数校验
- **预估耗时**：4 小时
- **前置依赖**：任务 2.5
- **产出物**：
  - `server/internal/service/tarot_test.go`
  - `server/internal/handler/tarot_test.go`

---

### Week 3：塔罗牌前端

**目标**：完成塔罗牌占卜的完整前端体验，包括首页、占卜页、动画效果。

#### 任务 3.1：实现首页 UI

- **任务描述**：设计并实现首页，展示 4 种占卜类型卡片（图标 + 名称 + 简短描述），每张卡片有悬停动画效果（放大、阴影加深）。点击卡片跳转到对应占卜页面。
- **预估耗时**：4 小时
- **前置依赖**：任务 1.10
- **产出物**：
  - `client/src/pages/Home.tsx`

#### 任务 3.2：实现塔罗牌占卜页 — 输入阶段

- **任务描述**：实现占卜页面的输入阶段 UI：
  - 问题输入框（可选，placeholder 提示用户输入问题）
  - 牌阵选择区域：展示 4 种牌阵的可视化预览（缩略图 + 名称 + 张数说明）
  - 选中牌阵高亮，开始占卜按钮
- **预估耗时**：4 小时
- **前置依赖**：任务 3.1
- **产出物**：
  - `client/src/pages/Tarot.tsx`（输入阶段部分）
  - `client/src/components/tarot/SpreadSelector.tsx`

#### 任务 3.3：实现塔罗牌 TypeScript 类型与 API 服务

- **任务描述**：定义塔罗牌相关的 TypeScript 类型（TarotCard、Spread、DrawResult、CardPosition 等），实现 API 服务函数（fetchCards、fetchSpreads、drawCards）。
- **预估耗时**：2 小时
- **前置依赖**：任务 1.11
- **产出物**：
  - `client/src/types/tarot.ts`
  - `client/src/services/tarot.ts`

#### 任务 3.4：实现抽牌动画（Framer Motion）

- **任务描述**：使用 Framer Motion 实现抽牌动画效果：
  - 牌组展示：牌堆叠的视觉效果
  - 洗牌动画：牌组分裂、交错、合并（可用简化的左右滑动）
  - 翻牌动画：牌面 3D 翻转（先显示牌背 → 翻转显示正面）
  - 逆位动画：逆位牌旋转 180°
  - 动画时序控制（依次翻牌，每张间隔 0.5s）
- **预估耗时**：6 小时（动画是用户体验核心，需精心打磨）
- **前置依赖**：任务 3.3
- **产出物**：
  - `client/src/components/tarot/DrawAnimation.tsx`
  - `client/src/styles/animations.css`（关键帧动画）

#### 任务 3.5：实现牌阵布局组件

- **任务描述**：实现 4 种牌阵的可视化布局：
  - 单牌布局：居中展示
  - 三牌阵：横向三列（过去/现在/未来）
  - 凯尔特十字阵：10 张牌的精确布局（参考 design.md 中的示意图）
  - 爱情十字阵：5 张牌的十字布局
  - 每张牌可点击展开详情。
- **预估耗时**：5 小时
- **前置依赖**：任务 3.4
- **产出物**：
  - `client/src/components/tarot/SpreadLayout.tsx`
  - `client/src/components/tarot/SingleSpread.tsx`
  - `client/src/components/tarot/ThreeSpread.tsx`
  - `client/src/components/tarot/CelticCross.tsx`
  - `client/src/components/tarot/LoveCross.tsx`

#### 任务 3.6：实现牌面解读展示

- **任务描述**：实现占卜结果的解读 UI：
  - 单牌解读：点击某张牌弹出 Modal，显示牌名、正/逆位、关键词、详细含义
  - 综合解读：页面底部展示所有牌的整体分析文字
  - "AI 深度解读"按钮（Phase 3 实现，先预留入口）
- **预估耗时**：4 小时
- **前置依赖**：任务 3.5
- **产出物**：
  - `client/src/components/tarot/CardReading.tsx`
  - `client/src/components/common/Modal.tsx`

#### 任务 3.7：实现 useTarot Hook 与状态管理

- **任务描述**：实现 `useTarot` Hook，封装占卜流程的状态管理（当前阶段：输入/抽牌/结果、选中的牌阵、抽牌结果、加载状态等）。串联整个占卜流程。
- **预估耗时**：3 小时
- **前置依赖**：任务 3.6
- **产出物**：
  - `client/src/hooks/useTarot.ts`

---

### Week 4：历史记录与联调

**目标**：完成历史记录功能，前后端联调，Bug 修复，部署测试。

#### 任务 4.1：实现占卜记录保存（后端）

- **任务描述**：在塔罗牌抽牌完成后，将结果保存到 `divination_records` 表。实现 `History` Repository 和 Service。实现历史记录 API：
  - `GET /api/history` — 分页获取历史列表（支持按类型筛选）
  - `GET /api/history/:id` — 获取单条详情
  - `DELETE /api/history/:id` — 删除记录（仅限本人）
- **预估耗时**：4 小时
- **前置依赖**：任务 2.5、1.4
- **产出物**：
  - `server/internal/model/divination.go`
  - `server/internal/repository/divination.go`
  - `server/internal/service/history.go`
  - `server/internal/handler/history.go`

#### 任务 4.2：实现历史记录页面（前端）

- **任务描述**：实现历史记录页面：
  - 列表展示（卡片形式，显示占卜类型图标、问题摘要、时间）
  - 按类型筛选（Tab 切换：全部/塔罗/星座/六爻/八字）
  - 点击卡片查看详情（跳转到结果页或弹窗展示）
  - 删除操作（确认弹窗）
  - 分页加载（上拉加载更多或分页按钮）
- **预估耗时**：4 小时
- **前置依赖**：任务 4.1、1.11
- **产出物**：
  - `client/src/pages/History.tsx`
  - `client/src/services/history.ts`

#### 任务 4.3：前后端全流程联调

- **任务描述**：端到端联调完整的塔罗牌占卜流程：注册/登录 → 进入首页 → 选择塔罗牌 → 输入问题 → 选择牌阵 → 抽牌动画 → 查看结果 → 查看历史记录。修复联调过程中发现的 API 对接问题、数据格式不一致、CORS 问题等。
- **预估耗时**：5 小时
- **前置依赖**：任务 4.2
- **产出物**：完整可运行的塔罗牌占卜流程

#### 任务 4.4：UI 细节优化与 Bug 修复

- **任务描述**：修复联调发现的 Bug，优化 UI 细节：
  - 加载状态（Loading spinner / skeleton）
  - 错误提示（Toast 通知）
  - 空状态展示（无历史记录时的提示）
  - 表单验证错误的友好提示
  - 页面过渡动画
  - 基础的响应式布局（确保桌面端正常显示）
- **预估耗时**：5 小时
- **前置依赖**：任务 4.3
- **产出物**：优化后的稳定版本

#### 任务 4.5：编写 API 文档与部署测试环境

- **任务描述**：整理 API 文档（可用 Swagger 或 Markdown），配置 Nginx 反向代理，将前端构建产物和后端服务部署到测试服务器（或本地 Docker）。
- **预估耗时**：4 小时
- **前置依赖**：任务 4.4
- **产出物**：
  - `docs/api.md`（API 文档）
  - 可访问的测试环境 URL

---

## 3. Phase 2 详细任务（功能增强）

### Week 5-6：星座运势

**目标**：实现 12 星座的日/周/月运势查询，包含完整的前端体验。

#### 任务 5.1：设计并实现运势生成算法

- **任务描述**：在 Service 层实现基于日期的确定性运势生成算法：
  - 使用 SHA-256(zodiac + date + period) 生成种子
  - 用种子初始化伪随机数生成器
  - 生成各维度评分（overall/love/career/wealth/health，1-5 星）
  - 生成幸运数字（1-9）和幸运颜色
  - 确保同一天同一星座的运势结果一致
- **预估耗时**：4 小时
- **前置依赖**：无
- **产出物**：
  - `server/internal/service/horoscope.go`

#### 任务 5.2：准备运势模板库

- **任务描述**：编写 12 星座 × 3 个运势周期 × 4 个维度 × 3 个评分等级的运势模板文本。模板使用占位符标记星座特质和评分区间，生成时填充具体内容。
  - 每个模板需覆盖：高分（4-5 星）、中等（3 星）、低分（1-2 星）
  - 每个维度还需准备 summary（总述）模板
- **预估耗时**：8 小时（模板数量多，需要丰富的内容）
- **前置依赖**：无
- **产出物**：
  - `server/data/horoscope/daily.json`
  - `server/data/horoscope/weekly.json`
  - `server/data/horoscope/monthly.json`
  - 或数据库表 `horoscope_templates` 的种子数据

#### 任务 5.3：实现星座运势模型、Repository 与 API

- **任务描述**：
  - 定义星座运势相关的数据模型和请求/响应结构体
  - 实现 Repository 层（模板查询）
  - 实现 Handler，提供以下 API：
    - `GET /api/horoscope/:zodiac` — 获取运势（参数：period、date）
    - `GET /api/horoscope/:zodiac/history` — 获取运势历史
  - 将星座运势记录保存到 `divination_records`
- **预估耗时**：4 小时
- **前置依赖**：任务 5.1、5.2
- **产出物**：
  - `server/internal/model/horoscope.go`
  - `server/internal/repository/horoscope.go`
  - `server/internal/handler/horoscope.go`

#### 任务 5.4：实现前端星座运势 TypeScript 类型与 API 服务

- **任务描述**：定义星座运势 TypeScript 类型，实现 API 服务函数。
- **预估耗时**：2 小时
- **前置依赖**：任务 5.3
- **产出物**：
  - `client/src/types/horoscope.ts`
  - `client/src/services/horoscope.ts`

#### 任务 5.5：实现 12 宫轮盘星座选择组件

- **任务描述**：实现 12 星座的圆形轮盘选择 UI：
  - 12 个星座图标按圆形排列
  - 点击选中星座（高亮动画）
  - 支持通过生日自动匹配星座（输入生日 → 计算星座 → 自动选中）
  - 使用 Framer Motion 实现选中动画
- **预估耗时**：5 小时
- **前置依赖**：任务 5.4
- **产出物**：
  - `client/src/components/horoscope/ZodiacWheel.tsx`

#### 任务 5.6：实现运势展示页面与雷达图

- **任务描述**：实现星座运势展示页面：
  - 运势周期 Tab 切换（日/周/月）
  - 日期选择器（日运势可选日期）
  - 雷达图：使用 Recharts 展示 5 个维度评分
  - 运势卡片：各维度详情展示（图文并茂）
  - 幸运元素展示（幸运数字、幸运颜色的视觉设计）
  - 使用 useHoroscope Hook 管理状态
- **预估耗时**：5 小时
- **前置依赖**：任务 5.5
- **产出物**：
  - `client/src/pages/Horoscope.tsx`
  - `client/src/components/horoscope/HoroscopeCard.tsx`
  - `client/src/components/horoscope/RadarChart.tsx`
  - `client/src/components/horoscope/LuckyElements.tsx`
  - `client/src/hooks/useHoroscope.ts`

#### 任务 5.7：星座运势模块联调与测试

- **任务描述**：前后端联调星座运势全流程，修复 Bug，验证运势生成的一致性（同一天结果相同）。
- **预估耗时**：3 小时
- **前置依赖**：任务 5.6
- **产出物**：星座运势功能完整可用

---

### Week 7-8：周易六爻与八字排盘

**目标**：实现六爻占卜和八字排盘的完整功能。

#### 任务 6.1：准备六十四卦数据

- **任务描述**：编写完整的 64 卦 JSON 数据，包括：
  - 卦 id、卦名（全称 + 简称）、上卦、下卦、6 位二进制
  - 卦辞、象辞
  - 六爻爻辞（JSON 数组）
  - 卦的总体描述
  - 同时准备 8 个基本卦（八卦）的数据
- **预估耗时**：6 小时（64 卦数据量大）
- **前置依赖**：无
- **产出物**：
  - `server/data/hexagrams/hexagrams.json`
  - `server/data/hexagrams/trigrams.json`

#### 任务 6.2：实现六爻数据模型与种子导入

- **任务描述**：定义 `Hexagram` 和 `Trigram` GORM 模型，编写种子数据导入逻辑。实现 Repository 层（按 ID 查卦、按上下卦组合查卦）。
- **预估耗时**：3 小时
- **前置依赖**：任务 6.1
- **产出物**：
  - `server/internal/model/hexagram.go`
  - `server/internal/repository/hexagram.go`
  - 种子数据导入更新到 `seed.go`

#### 任务 6.3：实现掷铜钱算法与卦象解析

- **任务描述**：在 Service 层实现六爻核心逻辑：
  - 掷铜钱算法：模拟三枚铜钱投掷，计算老阳/少阴/少阳/老阴
  - 生成本卦：6 次投掷结果组合为 6 爻，查找对应 64 卦
  - 生成变卦：动爻取反，查找变卦
  - 解读逻辑：动爻分析、五行生克、爻位分析
  - 返回完整的占卜结果（六爻详情 + 本卦 + 变卦 + 动爻位置）
- **预估耗时**：5 小时
- **前置依赖**：任务 6.2
- **产出物**：
  - `server/internal/service/liuyao.go`

#### 任务 6.4：实现六爻 API Handler 与路由

- **任务描述**：实现以下接口：
  - `GET /api/liuyao/hexagrams` — 获取 64 卦列表
  - `GET /api/liuyao/hexagrams/:id` — 获取单卦详情
  - `POST /api/liuyao/throw` — 掷铜钱占卜（可选 question）
- **预估耗时**：3 小时
- **前置依赖**：任务 6.3
- **产出物**：
  - `server/internal/handler/liuyao.go`

#### 任务 6.5：实现八字排盘算法

- **任务描述**：在 Service 层实现八字排盘核心算法：
  - 公历转农历（寿星万年历算法或查表法）
  - 年柱计算（以立春为界）
  - 月柱计算（以节气为界 + 年上起月法）
  - 日柱计算（基准日推算）
  - 时柱计算（日上起时法）
  - 五行分析（统计五行数量、判断日主强弱、确定用神/忌神）
  - 十神分析（以日干为基准）
  - 大运排盘（阳年男/阴年女顺排，反之逆排）
- **预估耗时**：8 小时（算法复杂度高，节气计算需精确）
- **前置依赖**：无
- **产出物**：
  - `server/internal/service/bazi.go`

#### 任务 6.6：实现八字 API Handler

- **任务描述**：实现以下接口：
  - `POST /api/bazi/calculate` — 八字排盘（参数：birth_date、birth_time、gender）
- **预估耗时**：3 小时
- **前置依赖**：任务 6.5
- **产出物**：
  - `server/internal/handler/bazi.go`

#### 任务 6.7：实现六爻前端页面

- **任务描述**：实现六爻占卜的完整前端体验：
  - 问题输入
  - 铜钱动画：三枚铜钱的抛掷动画（使用 Framer Motion 或 CSS 3D），点击"掷"按钮触发动画
  - 进度指示（1/6、2/6 ... 6/6）
  - 卦象展示：六爻从下到上依次展示，动爻特殊标记（颜色或闪烁）
  - 本卦 + 变卦对比展示
  - 解读区域：卦名、卦辞、逐爻解读、五行分析
- **预估耗时**：6 小时
- **前置依赖**：任务 6.4
- **产出物**：
  - `client/src/pages/LiuYao.tsx`
  - `client/src/components/liuyao/CoinAnimation.tsx`
  - `client/src/components/liuyao/HexagramDisplay.tsx`
  - `client/src/components/liuyao/LineResult.tsx`
  - `client/src/types/liuyao.ts`
  - `client/src/services/liuyao.ts`
  - `client/src/hooks/useLiuYao.ts`

#### 任务 6.8：实现八字排盘前端页面

- **任务描述**：实现八字排盘的完整前端体验：
  - 出生信息输入：日期选择器（支持公历/农历切换）、时间选择器（24 小时制）、性别选择
  - 四柱表格展示：年柱/月柱/日柱/时柱，每柱显示天干 + 地支 + 纳音 + 藏干
  - 五行分析图表：使用 Recharts 柱状图展示五行数量
  - 十神列表：展示各天干对应的十神
  - 大运时间轴：时间轴布局展示大运，当前大运高亮
- **预估耗时**：6 小时
- **前置依赖**：任务 6.6
- **产出物**：
  - `client/src/pages/BaZi.tsx`
  - `client/src/components/bazi/PillarTable.tsx`
  - `client/src/components/bazi/FiveElementChart.tsx`
  - `client/src/components/bazi/DaYunTimeline.tsx`
  - `client/src/components/bazi/TenGodList.tsx`
  - `client/src/types/bazi.ts`
  - `client/src/services/bazi.ts`
  - `client/src/hooks/useBaZi.ts`

#### 任务 6.9：六爻与八字模块联调测试

- **任务描述**：端到端联调六爻和八字的完整流程，验证算法正确性（用已知结果的案例测试八字排盘），修复 Bug。
- **预估耗时**：5 小时
- **前置依赖**：任务 6.7、6.8
- **产出物**：六爻 + 八字功能完整可用

#### 任务 6.10：更新首页与历史记录

- **任务描述**：更新首页，确保 4 种占卜类型的卡片入口都可正常跳转。更新历史记录页面，确保新增的星座/六爻/八字记录能正确展示。
- **预估耗时**：2 小时
- **前置依赖**：任务 6.9
- **产出物**：四种占卜类型全部可通过首页进入，历史记录正确展示

---

## 4. Phase 3 详细任务（上线准备）

### Week 9-10：AI 解读集成

**目标**：接入 AI 大模型，实现四种占卜类型的流式 AI 解读。

#### 任务 7.1：设计 AI Prompt 模板

- **任务描述**：为四种占卜类型分别设计 Prompt 模板：
  - 塔罗牌 Prompt：包含牌阵信息、每张牌的正逆位和位置含义
  - 星座运势 Prompt：包含星座特质、各维度评分
  - 六爻 Prompt：包含本卦/变卦信息、动爻位置、五行分析
  - 八字 Prompt：包含四柱、五行分析、十神、大运
  - 通用要求：语言温暖亲切、结合用户问题、给出可操作建议、300-500 字
- **预估耗时**：4 小时
- **前置依赖**：无
- **产出物**：
  - `server/data/prompts/`（各类型的 Prompt 模板文件）

#### 任务 7.2：实现 AI 服务抽象层

- **任务描述**：实现 AI 服务的抽象接口，支持多种大模型 Provider（OpenAI、文心一言、通义千问等）。通过配置文件切换 Provider。实现 HTTP 客户端调用大模型 API，处理认证、请求构造、响应解析。
- **预估耗时**：5 小时
- **前置依赖**：任务 7.1
- **产出物**：
  - `server/internal/service/ai.go`（AI 服务接口 + 实现）

#### 任务 7.3：实现流式输出（SSE）

- **任务描述**：实现 Server-Sent Events 流式输出：
  - 后端：将大模型的流式响应转发为 SSE 格式（`data: {"text": "..."}\n\n`）
  - 处理流结束标记（`data: [DONE]`）
  - 错误处理（AI 服务不可用、超时等）
  - 实现 `POST /api/ai/interpret` 接口
- **预估耗时**：4 小时
- **前置依赖**：任务 7.2
- **产出物**：
  - `server/internal/handler/ai.go`
  - SSE 流式响应实现

#### 任务 7.4：实现 AI 限流与错误处理

- **任务描述**：实现 AI 接口的限流策略（每用户每分钟 5 次），实现 AI 服务不可用时的降级方案（显示基础解读 + 提示"AI 解读暂不可用"），实现请求超时处理（30 秒超时）。
- **预估耗时**：3 小时
- **前置依赖**：任务 7.3
- **产出物**：限流逻辑更新到 `middleware/ratelimit.go`，AI 降级逻辑

#### 任务 7.5：实现前端 AI 解读 UI

- **任务描述**：在四种占卜类型的结果页面中，实现"AI 深度解读"按钮和解读展示区域：
  - 点击按钮触发 AI 解读请求
  - 使用 EventSource 或 fetch + ReadableStream 接收 SSE 流
  - 逐字显示 AI 解读内容（打字机效果）
  - 加载状态（等待 AI 响应时的 loading 动画）
  - 错误状态（AI 不可用时的提示 + 重试按钮）
  - 实现 useAI Hook 管理 AI 解读状态
- **预估耗时**：5 小时
- **前置依赖**：任务 7.3
- **产出物**：
  - `client/src/components/common/AIReading.tsx`
  - `client/src/services/ai.ts`
  - `client/src/hooks/useAI.ts`
  - 更新 `Tarot.tsx`、`Horoscope.tsx`、`LiuYao.tsx`、`BaZi.tsx` 集成 AI 解读

#### 任务 7.6：AI 解读模块联调测试

- **任务描述**：联调 AI 解读全流程，测试流式输出效果、错误处理、限流策略。验证 Prompt 模板生成的解读质量。
- **预估耗时**：3 小时
- **前置依赖**：任务 7.5
- **产出物**：AI 解读功能完整可用

---

### Week 11：优化与完善

**目标**：全面优化 UI/UX、性能、移动端适配。

#### 任务 8.1：UI/UX 全面优化

- **任务描述**：
  - 统一视觉风格：配色方案、字体、间距、圆角等
  - 优化页面过渡动画（页面切换的平滑过渡）
  - 优化占卜结果的视觉呈现（卡片阴影、渐变背景等）
  - 添加空状态设计（无历史记录时的插画引导）
  - 优化错误页面（404、500）
  - 添加页面 Title 和 Meta 标签
- **预估耗时**：6 小时
- **前置依赖**：Phase 2 完成
- **产出物**：视觉效果全面提升的 UI

#### 任务 8.2：性能优化

- **任务描述**：
  - 前端路由懒加载（React.lazy + Suspense）
  - 图片懒加载（塔罗牌图片、星座图标等）
  - API 响应缓存（星座运势可缓存当日结果）
  - React Query 缓存策略优化
  - 前端构建产物分析（Bundle Analyzer），优化打包体积
  - 后端 GORM 查询优化（避免 N+1 查询）
- **预估耗时**：5 小时
- **前置依赖**：任务 8.1
- **产出物**：性能优化后的版本

#### 任务 8.3：移动端适配

- **任务描述**：
  - 响应式布局适配（使用 TailwindCSS 的响应式前缀）
  - 移动端导航（汉堡菜单 / 底部 Tab 栏）
  - 触摸手势支持（滑动翻牌等）
  - 移动端占卜页面的布局调整（牌阵缩小适配小屏）
  - 移动端铜钱动画优化（性能考虑）
  - 测试主流手机尺寸（iPhone SE / 12 / 14 Pro Max、Android 主流尺寸）
- **预估耗时**：6 小时
- **前置依赖**：任务 8.1
- **产出物**：移动端适配完成的版本

#### 任务 8.4：安全审计

- **任务描述**：
  - 检查所有 API 的输入验证是否完整
  - 检查 SQL 注入风险（确认全部使用参数化查询）
  - 检查 XSS 风险（前端输出转义）
  - 检查 CORS 配置是否正确
  - 检查 JWT 密钥是否使用环境变量（不硬编码）
  - 检查密码加密是否使用 bcrypt
  - 检查日志中是否泄露敏感信息
  - 更新 CORS 配置为生产域名
- **预估耗时**：4 小时
- **前置依赖**：Phase 2 完成
- **产出物**：安全审计报告 + 修复

#### 任务 8.5：编写用户文档

- **任务描述**：编写用户使用说明文档，包括：
  - 各占卜类型的使用方法和注意事项
  - 常见问题（FAQ）
  - 隐私政策说明
- **预估耗时**：3 小时
- **前置依赖**：无
- **产出物**：
  - `docs/user-guide.md`

---

### Week 12：上线准备与部署

**目标**：完成生产环境部署、监控配置、压力测试、正式上线。

#### 任务 9.1：准备生产环境配置

- **任务描述**：
  - 编写生产环境 `config.production.yaml`（JWT 密钥使用环境变量、关闭 debug 模式、配置生产数据库路径）
  - 编写 Dockerfile（后端 Go 应用 + 前端 Nginx 静态服务）
  - 编写 `docker-compose.yml`（后端 + 前端 + Nginx 反向代理）
  - 配置 HTTPS（Let's Encrypt 或云服务 SSL）
  - 配置域名解析
- **预估耗时**：5 小时
- **前置依赖**：任务 8.4
- **产出物**：
  - `Dockerfile`（后端）
  - `Dockerfile.frontend`（前端）
  - `docker-compose.yml`
  - `nginx/nginx.conf`
  - `server/config/config.production.yaml`

#### 任务 9.2：部署到生产服务器

- **任务描述**：
  - 准备服务器环境（安装 Docker、Docker Compose）
  - 部署应用（docker-compose up）
  - 配置 Nginx 反向代理（前端静态文件 + `/api` 转发后端）
  - 配置 SSL 证书
  - 验证所有功能在生产环境正常工作
- **预估耗时**：4 小时
- **前置依赖**：任务 9.1
- **产出物**：可访问的生产环境 URL

#### 任务 9.3：配置监控与日志

- **任务描述**：
  - 配置后端结构化日志输出（JSON 格式，方便后续接入日志系统）
  - 配置健康检查接口（`GET /api/health`）
  - 配置基础监控（请求量、错误率、响应时间）
  - 配置告警（服务宕机、错误率飙升）
  - 日志轮转配置（避免日志文件过大）
- **预估耗时**：4 小时
- **前置依赖**：任务 9.2
- **产出物**：
  - `server/internal/handler/health.go`
  - 监控与告警配置

#### 任务 9.4：压力测试

- **任务描述**：使用压测工具（如 wrk、hey、k6）对关键 API 进行压力测试：
  - 用户登录 API（并发 100）
  - 塔罗牌抽牌 API（并发 50）
  - 星座运势查询 API（并发 100）
  - AI 解读 API（并发 10，受 AI 服务限制）
  - 记录响应时间、吞吐量、错误率
  - 根据结果优化（数据库连接池、缓存策略等）
- **预估耗时**：4 小时
- **前置依赖**：任务 9.2
- **产出物**：
  - `docs/benchmark.md`（压测报告）

#### 任务 9.5：正式上线与回归测试

- **任务描述**：
  - 全流程回归测试（所有占卜类型 + 用户系统 + 历史记录 + AI 解读）
  - 修复发现的最终 Bug
  - 确认监控和告警正常
  - 确认备份策略就绪
  - 正式开放访问
- **预估耗时**：4 小时
- **前置依赖**：任务 9.4
- **产出物**：🎉 正式上线

---

## 5. 技术风险与应对

### 5.1 高优先级风险

| 风险 | 影响 | 可能性 | 应对方案 |
|------|------|--------|---------|
| **八字排盘算法精度** | 排盘结果错误，用户信任丧失 | 中 | ① 使用成熟的万年历库（如 `6tail/lunar` Go 库）而非自行实现；② 用已知案例进行交叉验证（至少 100 个测试用例）；③ 上线前邀请专业人士审核结果 |
| **节气计算误差** | 月柱切换边界错误 | 中 | ① 使用天文算法计算节气精确时间；② 对于边界日期（节气当天），提示用户确认；③ 维护 2000-2050 年的节气查表作为后备 |
| **AI 大模型服务不稳定** | AI 解读功能不可用 | 高 | ① 实现 Provider 抽象层，支持多 Provider 快速切换；② AI 不可用时降级为基础解读（模板文本）；③ 实现请求队列和重试机制；④ 监控 AI 服务可用性 |
| **塔罗牌图片资源** | 缺少高质量牌面图片影响体验 | 中 | ① Phase 1 先使用文字 + CSS 绘制的简化牌面；② 并行寻找开源塔罗牌图片资源；③ 考虑使用 AI 生图工具生成；④ 预留图片替换接口 |

### 5.2 中优先级风险

| 风险 | 影响 | 可能性 | 应对方案 |
|------|------|--------|---------|
| **SQLite 并发性能** | 高并发时数据库锁争用 | 低 | ① SQLite 读写锁模型适合低并发场景（占卜网站预期并发不高）；② 如果遇到瓶颈，可迁移至 PostgreSQL；③ 实现连接池配置 |
| **动画性能（低端设备）** | 移动端动画卡顿 | 中 | ① 使用 CSS transform/opacity 代替 layout 属性；② 提供"减少动画"开关；③ 铜钱动画在低端设备上简化为静态结果展示 |
| **前端包体积过大** | 首屏加载慢 | 低 | ① 路由懒加载；② 图片懒加载 + WebP 格式；③ 使用 Bundle Analyzer 分析并优化；④ 外部依赖按需引入 |
| **运势模板内容质量** | 模板文本生硬、缺乏吸引力 | 中 | ① 参考专业星座运势网站的文案风格；② 后续用 AI 生成更自然的模板；③ 收集用户反馈持续优化 |

### 5.3 低优先级风险

| 风险 | 影响 | 可能性 | 应对方案 |
|------|------|--------|---------|
| **农历闰月处理** | 八字排盘边界情况 | 低 | ① 使用成熟的农历库处理闰月；② 前端提示用户确认农历日期 |
| **JWT Token 泄露** | 账号安全 | 低 | ① Access Token 有效期 1 小时；② 后续可升级为 HttpOnly Cookie；③ 实现 Token 黑名单（登出时） |
| **多语言支持** | 后续国际化需求 | 低 | 当前阶段仅支持中文，后续如需国际化，前端使用 i18n 库即可 |

---

## 6. 开发规范

### 6.1 Git 分支策略

采用简化版 Git Flow：

```
main ────────────────────────────────────────── 生产环境代码
  │
  ├── develop ───────────────────────────────── 开发主分支
  │     │
  │     ├── feature/tarot-draw ──────────────── 功能分支
  │     ├── feature/horoscope ─────────────────
  │     ├── feature/liuyao ───────────────────
  │     └── feature/bazi ─────────────────────
  │
  ├── hotfix/fix-login-bug ──────────────────── 紧急修复
  └── release/v1.0.0 ────────────────────────── 发布分支
```

**分支命名规范**：
- 功能分支：`feature/<功能名>`（如 `feature/tarot-draw`、`feature/horoscope-api`）
- 修复分支：`fix/<问题描述>`（如 `fix/login-validation`）
- 紧急修复：`hotfix/<问题描述>`
- 发布分支：`release/v<版本号>`

**分支规则**：
- `main` 分支受保护，禁止直接推送，仅通过 PR 合并
- `develop` 分支为日常开发的集成分支
- 功能分支从 `develop` 创建，完成后合并回 `develop`
- 发布时从 `develop` 创建 `release` 分支，测试通过后合并到 `main` 并打 Tag

### 6.2 代码规范

#### Go 代码规范

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- 使用 `gofmt` 格式化代码（提交前自动格式化）
- 使用 `golangci-lint` 进行静态检查
- 包名小写、无下划线、简短
- 导出函数/类型必须有注释
- 错误处理：不丢弃 error，使用 `fmt.Errorf("xxx: %w", err)` 包装
- 日志使用 zerolog，不使用 `fmt.Println`

**推荐 golangci-lint 配置**（`.golangci.yml`）：
```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
```

#### React/TypeScript 代码规范

- 使用 ESLint + Prettier 格式化和检查
- 组件使用函数式组件 + Hooks
- 使用 TypeScript 严格模式（`strict: true`）
- 组件文件使用 PascalCase 命名（如 `TarotCard.tsx`）
- 工具函数使用 camelCase 命名（如 `formatDate.ts`）
- 类型定义集中放在 `types/` 目录
- 避免使用 `any`，必须时使用 `unknown` 并做类型守卫
- CSS 使用 TailwindCSS，避免自定义 CSS（除动画外）

**ESLint 核心规则**：
```json
{
  "extends": [
    "eslint:recommended",
    "plugin:react/recommended",
    "plugin:react-hooks/recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier"
  ],
  "rules": {
    "react/react-in-jsx-scope": "off",
    "@typescript-eslint/no-unused-vars": "error",
    "@typescript-eslint/explicit-function-return-type": "warn"
  }
}
```

### 6.3 提交规范

采用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型**：

| Type | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(tarot): 实现凯尔特十字牌阵` |
| `fix` | Bug 修复 | `fix(auth): 修复 Token 过期未刷新的问题` |
| `docs` | 文档变更 | `docs: 更新 API 文档` |
| `style` | 代码格式（不影响逻辑） | `style: 运行 gofmt 格式化` |
| `refactor` | 重构（不新增功能/不修复 Bug） | `refactor(service): 抽取公共验证逻辑` |
| `perf` | 性能优化 | `perf(frontend): 实现路由懒加载` |
| `test` | 测试相关 | `test(tarot): 补充洗牌算法单元测试` |
| `chore` | 构建/工具/依赖变更 | `chore: 升级 Gin 到 v1.9.1` |
| `ci` | CI/CD 配置 | `ci: 添加 GitHub Actions 配置` |

**Scope 范围**（可选）：

`tarot`、`horoscope`、`liuyao`、`bazi`、`auth`、`ai`、`history`、`frontend`、`backend`、`config`

**示例**：
```
feat(liuyao): 实现掷铜钱算法与卦象解析

- 实现三枚铜钱投掷模拟
- 支持老阳/少阴/少阳/老阴四种结果
- 自动计算本卦和变卦
- 添加动爻标记

Closes #42
```

### 6.4 测试策略

#### 测试金字塔

```
          ┌───────────┐
          │  E2E 测试  │  ← Phase 3 补充（少量关键流程）
          ├───────────┤
          │  集成测试   │  ← API 接口测试
          ├───────────┤
          │  单元测试   │  ← 核心算法、Service 层
          └───────────┘
```

#### 后端测试

| 层级 | 覆盖范围 | 工具 | 目标覆盖率 |
|------|---------|------|-----------|
| 单元测试 | 算法（洗牌、运势生成、排盘）、Service 层逻辑 | `testing` + `testify` | 核心算法 > 90%，Service > 80% |
| 集成测试 | API 接口（Handler → Service → DB） | `httptest` + `testify` | 关键 API 100% |
| 基准测试 | 核心算法性能 | `testing.B` | 排盘算法 < 100ms |

**核心测试用例**：

- **塔罗牌**：
  - 洗牌后 78 张牌不重复
  - 4 种牌阵抽牌数量正确（1/3/5/10）
  - 正逆位概率接近 50:50（大样本统计）
  - 牌面数据完整性（所有字段非空）

- **星座运势**：
  - 同一天同一星座结果一致（确定性）
  - 评分范围 1-5
  - 幸运数字范围 1-9
  - 12 星座都能正确生成运势

- **六爻**：
  - 掷铜钱概率分布正确（老阳 12.5%、少阴 37.5%、少阳 37.5%、老阴 12.5%）
  - 本卦/变卦查找正确
  - 动爻标记正确

- **八字排盘**：
  - 用已知案例验证四柱正确性（至少 50 个案例）
  - 节气边界日期验证
  - 五行统计正确
  - 大运排盘方向正确

#### 前端测试

| 层级 | 覆盖范围 | 工具 |
|------|---------|------|
| 组件测试 | 关键组件的渲染和交互 | Vitest + React Testing Library |
| Hook 测试 | 自定义 Hook 的状态逻辑 | Vitest + @testing-library/react-hooks |

**测试命令**：

```bash
# 后端测试
cd server && go test ./... -v -cover

# 前端测试
cd client && npm test
```

---

## 7. 本地开发环境搭建步骤

### 7.1 前置环境要求

| 工具 | 最低版本 | 推荐版本 | 说明 |
|------|---------|---------|------|
| Go | 1.21 | 1.22+ | 后端语言运行时 |
| Node.js | 18 | 20+ | 前端构建运行时 |
| npm | 9 | 10+ | 前端包管理器 |
| Git | 2.30 | 最新 | 版本控制 |
| SQLite | 3.35 | 最新 | 数据库（macOS/Linux 通常已预装） |

### 7.2 环境安装

#### macOS（Homebrew）

```bash
# 安装 Go
brew install go

# 安装 Node.js（推荐使用 nvm 管理版本）
brew install nvm
nvm install 20
nvm use 20

# 验证安装
go version    # 应输出 go version go1.22+
node -v       # 应输出 v20.x.x
npm -v        # 应输出 10.x.x
```

#### Ubuntu/Debian

```bash
# 安装 Go
wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js（通过 NodeSource）
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# 验证安装
go version
node -v
npm -v
```

### 7.3 项目初始化

```bash
# 1. 克隆项目
git clone <repo-url> zhanbu
cd zhanbu

# 2. 后端初始化
cd server

# 配置 Go 代理（国内用户）
go env -w GOPROXY=https://goproxy.cn,direct

# 安装依赖
go mod tidy

# 3. 前端初始化
cd ../client

# 安装依赖
npm install

# 4. 回到项目根目录
cd ..
```

### 7.4 配置文件

```bash
# 后端配置（复制模板后修改）
cd server
cp config/config.yaml config/config.local.yaml
```

编辑 `config/config.local.yaml`：

```yaml
server:
  port: 8080
  mode: debug

database:
  path: ./data/zhanbu.db

jwt:
  secret: "your-local-dev-secret-key"  # 本地开发随便写
  access_ttl: 1h
  refresh_ttl: 168h

ai:
  provider: openai
  api_key: "sk-your-api-key"           # 填入你的 AI API Key
  model: "gpt-4"
  base_url: "https://api.openai.com/v1"
  max_tokens: 1000
  temperature: 0.7

rate_limit:
  ai_per_minute: 5
  api_per_minute: 60

cors:
  allowed_origins:
    - "http://localhost:5173"
```

### 7.5 启动开发服务

**终端 1 — 后端**：

```bash
cd server

# 初始化数据库（首次运行）
go run main.go --migrate

# 导入种子数据（首次运行）
go run main.go --seed

# 启动开发服务（支持热重载可安装 air）
go run main.go
# 或使用 air 热重载：
# go install github.com/air-verse/air@latest
# air

# ✅ 后端启动在 http://localhost:8080
```

**终端 2 — 前端**：

```bash
cd client

# 启动开发服务器
npm run dev

# ✅ 前端启动在 http://localhost:5173
# ✅ API 请求自动代理到 http://localhost:8080
```

### 7.6 开发工具推荐

| 工具 | 用途 | 安装 |
|------|------|------|
| VS Code | 主力 IDE | [下载](https://code.visualstudio.com/) |
| GoLand | Go 专用 IDE（可选） | [下载](https://www.jetbrains.com/go/) |
| air | Go 热重载 | `go install github.com/air-verse/air@latest` |
| golangci-lint | Go 静态检查 | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| DB Browser for SQLite | 数据库可视化 | [下载](https://sqlitebrowser.org/) |
| Postman / Bruno | API 调试 | [Postman](https://www.postman.com/) / [Bruno](https://www.usebruno.com/) |
| React DevTools | React 调试 | Chrome 扩展 |

### 7.7 常见问题

**Q: `go mod tidy` 报错 "dial tcp: lookup proxy.golang.org: no such host"**
A: 设置国内代理：`go env -w GOPROXY=https://goproxy.cn,direct`

**Q: 前端启动后 API 请求 404**
A: 检查 `vite.config.ts` 中的代理配置是否正确指向 `http://localhost:8080`

**Q: SQLite 数据库文件在哪里**
A: 默认在 `server/data/zhanbu.db`，可通过配置文件修改路径

**Q: 如何重置数据库**
A: 删除 `server/data/zhanbu.db` 文件，重新运行 `go run main.go --migrate && go run main.go --seed`

**Q: AI 解读功能报错**
A: 检查 `config.local.yaml` 中的 AI API Key 是否正确配置，确保网络可以访问对应的 AI 服务

---

> **文档维护说明**：本文档应随项目推进持续更新。任务完成后标记状态，新增任务及时补充。
