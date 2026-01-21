# Go Shortener - 高性能分布式短链接服务

[![Go](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://golang.org/)
[![Go-Zero](https://img.shields.io/badge/Framework-Go--Zero-green.svg)](https://go-zero.dev/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**Go Shortener** 是一个基于 [go-zero](https://github.com/zeromicro/go-zero) 微服务框架构建的高性能短链接生成系统。它不仅实现了长短链的转换，还针对高并发场景进行了深度优化，集成了 Base62 编码、多级缓存架构及布隆过滤器，确保服务的高可用与低延迟。

## ✨ 核心特性

- **高性能架构**：基于 Go-Zero 框架，原生支持自适应熔断与限流，轻松应对流量突发。
- **多级缓存策略**：
  - **L1**: Redis 热点数据缓存，毫秒级响应。
  - **L2**: 数据库持久化存储。
- **布隆过滤器 (Bloom Filter)**：
  - 集成 Redis 布隆过滤器，有效拦截不存在的短链请求，防止缓存穿透。
  - 支持服务启动时异步预热，自动加载热点数据。
- **并发安全发号器**：优化的数据库发号逻辑，配合 Base62 编码，确保在高并发下短链唯一且短小。
- **有效期管理**：支持设置短链接的生命周期（秒级），过期自动失效，防止无效数据堆积。
- **安全机制**：内置短链黑名单检测，拦截恶意敏感词生成。

## 🛠 技术栈

- **编程语言**: Golang (1.20+)
- **微服务框架**: Go-Zero
- **数据存储**: MySQL 8.0
- **高速缓存**: Redis 7.0 (支持 Bloom Filter)
- **容器化**: Docker & Docker Compose
- **网关**: Nginx

## 📂 目录结构

```text
.
├── etc/                # 配置文件
├── internal/           # 业务逻辑 (Handler, Logic, SVC, Model)
├── model/              # 数据库模型与 SQL 文件
├── pkg/                # 公共工具包 (Base62, MD5, Connect)
├── static/             # 演示用简易前端
├── shortener.api       # API 定义文件
└── shortener.go        # 程序入口
```

## 🚀 快速开始

### 1. 环境准备

确保您的环境已安装：
- Go 1.20+
- MySQL
- Redis

### 2. 数据库初始化

请在 MySQL 中创建一个名为 `go_shortener` 的数据库，并依次执行 `model` 目录下的 SQL 脚本：

1. `model/sequence.sql` —— 初始化发号器表
2. `model/short_url_map.sql` —— 初始化主数据表

### 3. 配置运行

1. 复制配置文件副本：
   ```bash
   cp etc/shortener-api.yaml etc/shortener-api.yaml.local
   ```
2. 编辑 `.local` 文件，填写您的 MySQL 和 Redis 连接信息（**注意：不要将真实密码提交到代码仓库**）。
3. 启动服务：
   ```bash
   go mod tidy
   go run shortener.go -f etc/shortener-api.yaml.local
   ```

服务启动后，监听端口默认为 `8888`。

## 🔌 API 接口文档

### 1. 生成短链接 (Convert)

将长链接转换为短链接，并可指定有效期。

- **URL**: `/convert`
- **Method**: `POST`
- **Content-Type**: `application/json`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `longUrl` | string | 是 | 原始长链接 (需包含 http/https 协议头) |
| `seconds` | int64 | 否 | 有效期(秒)。不传或为0表示永久有效 |

**请求示例**:

```json
{
    "longUrl": "https://www.bilibili.com/video/BV1xx411c7mD",
    "seconds": 86400
}
```

**响应示例**:

```json
{
    "shortUrl": "hp.cn/1u9"
}
```

### 2. 访问/重定向 (Show)

访问生成的短链接，系统会自动重定向到原始长链接。

- **URL**: `/:shortUrl` (例如: `/1u9`)
- **Method**: `GET`

**响应**:
- **302 Found**: 重定向至长链接。
- **404 Not Found**: 链接不存在。
- **400 Bad Request**: 链接已过期。

## 🗺️ 未来规划 (Roadmap)

- [x] 基础短链转换与重定向
- [x] 引入 Redis 缓存与布隆过滤器
- [x] 实现短链有效期控制
- [ ] **自定义短链后缀**：支持用户指定个性化后缀（如 `/mypage`）。
- [ ] **数据统计大屏**：记录 PV/UV、访问来源、地理位置分布。
- [ ] **分布式 ID 生成**：引入 Snowflake 算法替代数据库发号，支持更大规模集群。
- [ ] **用户鉴权系统**：支持多用户登录与短链管理后台。

## 🤝 贡献与支持

如果您觉得这个项目对您有帮助，请给一个 ⭐️ Star！
欢迎提交 Issue 或 Pull Request 参与贡献。

## 📄 License

MIT License
```