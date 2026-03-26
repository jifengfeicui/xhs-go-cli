# xhs-go-cli

Go CLI for Xiaohongshu source/query/search/detail/qualify pipeline backed by SQLite.

## 现状

当前已实现首版 5 个命令：

- `import-sources`
- `query-gen`
- `search`
- `fetch-detail`
- `qualify`

主存储使用 SQLite。

---

## 环境准备

### 1. Go

当前环境中 Go 放在：

- `GOROOT=/root/go/go`
- `PATH=/root/go/go/bin:$PATH`

运行前建议先导出：

```bash
export GOROOT=/root/go/go
export PATH=/root/go/go/bin:$PATH
```

### 2. Xiaohongshu MCP

当前使用本地 MCP HTTP 服务，通过 `config.yaml` 配置。

**重要：启动 MCP 时要在 `bin/` 目录里启动，不要在上一级目录启动。**

当前实际路径：

- 运行目录：`/root/.openclaw/workspace/projects/xiaohongshu-ops/runtime/xiaohongshu-mcp/bin`
- 二进制：`/root/.openclaw/workspace/projects/xiaohongshu-ops/runtime/xiaohongshu-mcp/bin/xiaohongshu-mcp-linux-amd64`
- cookies：`/root/.openclaw/workspace/projects/xiaohongshu-ops/runtime/xiaohongshu-mcp/bin/cookies.json`

启动方式：

```bash
cd /root/.openclaw/workspace/projects/xiaohongshu-ops/runtime/xiaohongshu-mcp/bin
./xiaohongshu-mcp-linux-amd64
```

### 3. 配置文件

创建 `config.yaml`：

```yaml
db: ./data/xhs.db
mcp:
  base-url: http://127.0.0.1:18060
```

- `db`: 数据库路径（默认 `./data/xhs.db`）
- `mcp.base-url`: MCP 服务地址

命令行也可覆盖配置：

```bash
xhs-go-cli search --db ./other.db --mcp-url http://other:port
```

---

## 项目结构

```text
projects/xhs-go-cli/
├── go.mod
├── main.go
├── README.md
├── config.yaml
├── .gitignore
└── internal/
    ├── cmd/          # CLI 命令
    ├── db/           # 数据库层 (GORM)
    ├── model/        # 数据模型
    ├── repository/   # 数据仓库
    ├── source/       # 来源导入
    ├── querygen/     # 查询生成
    ├── search/       # 搜索
    ├── detail/       # 详情获取
    ├── qualify/      # 资质校验
    └── mcp/          # MCP 客户端
```

---

## SQLite 数据表

当前会自动初始化以下表（GORM AutoMigrate）：

- `sources`
- `generated_queries`
- `search_results`
- `details`
- `qualifications`

---

## 命令说明

## 1. import-sources

把来源 JSON 导入 SQLite。

### 用法

```bash
go run . import-sources --input <sources_json>
```

### 示例

```bash
cd /root/.openclaw/workspace/projects/xhs-go-cli
export GOROOT=/root/go/go
export PATH=/root/go/go/bin:$PATH

go run . import-sources \
  --input /root/.openclaw/workspace/projects/xiaohongshu-ops/data/xhs_source_records.json
```

### 当前说明

- 当前已实现
- 只迁来源数据
- 不迁历史 `search_results / details / qualifications`

---

## 2. query-gen

从 SQLite 中的来源记录生成 query，并写入 `generated_queries`。

### 用法

```bash
go run . query-gen --limit <n> --per-source <n>
```

### 参数

- `--limit`：处理多少条来源
- `--per-source`：每条来源生成多少个 query

### 示例

```bash
go run . query-gen --limit 5 --per-source 3
```

### 当前行为

- 会自动归类来源类型
- 当前分类包括：
  - `mall`
  - `brand`
  - `official_event`
  - `info_account`
  - `generic`
- 生成结果会同时：
  - 输出到 stdout
  - 写入 `generated_queries`

---

## 3. search

读取 `generated_queries`，调用 MCP 搜索接口，把结果写入 `search_results`。

### 用法

```bash
go run . search --limit <n> --page-size <n>
```

### 参数

- `--limit`：本轮取多少条 query
- `--page-size`：每个 query 请求多少条搜索结果

### 示例

```bash
go run . search --limit 5 --page-size 5
```

### 当前行为

- 从 `generated_queries` 取 query
- 调 `/api/v1/feeds/search`
- 结果写入 `search_results`
- stdout 输出每条 query 的 `stored` 数量或错误信息

### 当前已知问题

- MCP 搜索服务启动目录不对时，可能读不到 cookies
- 当前某些 query 已出现 `stored: 0`，需要继续排查是：
  - MCP 返回空
  - 还是当前 Go 解析路径不对

---

## 4. fetch-detail

从 `search_results` 中取候选，调用 MCP detail 接口，把结果写入 `details`。

### 用法

```bash
go run . fetch-detail --limit <n> --concurrency <n>
```

### 参数

- `--limit`：本轮最多取多少条搜索结果
- `--concurrency`：detail 拉取并发数

### 示例

```bash
go run . fetch-detail --limit 10 --concurrency 3
```

### 当前行为

- 从 `search_results` 读取候选
- 并发调用 `/api/v1/feeds/detail`
- 把 detail JSON 写进 `details`
- stdout 返回每条 feed 的抓取状态

---

## 5. qualify

从 `details` 中读取 detail，做结构化判断，把结果写入 `qualifications`。

### 用法

```bash
go run . qualify --limit <n>
```

### 参数

- `--limit`：本轮处理多少条 detail

### 示例

```bash
go run . qualify --limit 10
```

### 当前判断门槛

必须同时有：

- `title`
- `source_link`
- `claim_rule`
- `location`
- `participation_method`

### 输出

- 写入 `qualifications`
- `status` 为：
  - `accepted`
  - `rejected`
- 若 rejected，会写 `reason`

---

## 端到端示例

```bash
cd /root/.openclaw/workspace/projects/xhs-go-cli
export GOROOT=/root/go/go
export PATH=/root/go/go/bin:$PATH

# 1) 导入来源
go run . import-sources \
  --input /root/.openclaw/workspace/projects/xiaohongshu-ops/data/xhs_source_records.json

# 2) 生成 query
go run . query-gen --limit 5 --per-source 3

# 3) 搜索
go run . search --limit 5 --page-size 5

# 4) 拉详情
go run . fetch-detail --limit 10 --concurrency 3

# 5) 做达标判断
go run . qualify --limit 10
```

---

## 当前真实状态

当前 CLI 已完成首版骨架，且以下模块测试通过：

- `internal/db`
- `internal/source`
- `internal/querygen`
- `internal/search`
- `internal/detail`
- `internal/qualify`

但端到端真实链路目前还有一个未完全收口的问题：

- `search` 这层已经能调用 MCP
- 但当前这轮真实数据里出现了 `stored: 0`
- 下一步需要继续确认：
  - 是 MCP 搜索确实返回空
  - 还是 Go 当前 search 结果解析路径还要调整

---

## Git 提交

当前关键提交：

- `15acc51` — `feat: scaffold go cli with source query and search flow`
- `48505ef` — `feat: add detail fetch and qualification commands`
