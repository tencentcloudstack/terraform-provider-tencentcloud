# 项目上下文

## 项目目的
Terraform Provider for TencentCloud - 腾讯云官方 Terraform Provider，用于实现腾讯云服务的基础设施自动化和管理。该 Provider 允许用户使用基础设施即代码（IaC）的原则来定义、配置和管理腾讯云资源。

## 技术栈
- **编程语言**: Go 1.17+
- **框架**: Terraform Plugin SDK v2 (github.com/hashicorp/terraform-plugin-sdk/v2)
- **云服务 SDK**: TencentCloud SDK Go (tencentcloud-sdk-go) - 多个服务专用包
- **测试**: Go 标准测试包 + Terraform 验收测试
- **构建工具**: GNU Make, GolangCI-Lint, tfproviderlint
- **版本控制**: Git
- **包管理**: Go Modules (vendor 模式)

## 项目规范

### 代码风格
- **代码格式化**: 使用 `gofmt` 和 `goimports` 进行自动代码格式化
- **代码检查**: 通过 golangci-lint 强制执行自定义配置 (.golangci.yml)
  - 启用的检查器: errcheck, gofmt, ineffassign, misspell, unconvert, unused, vet
  - 通过 tfproviderlint 进行 Terraform 特定检查
- **命名规范**:
  - 资源文件: `resource_tc_<service>_<resource_name>.go`
  - 数据源文件: `data_source_tc_<service>_<resource_name>.go`
  - 服务文件: `service_tencentcloud_<service>.go`
  - 函数命名: 导出函数使用大驼峰，内部函数使用小驼峰
  - Terraform 中的资源名: `tencentcloud_<service>_<resource_name>`
- **文件组织**: 服务按子目录组织在 `tencentcloud/services/` 下
- **导入别名**: 公共包使用 `tccommon`，服务包使用 `svc<service>` 格式

### 架构模式
- **Provider 架构**: 遵循 Terraform Plugin SDK v2 模式
  - 资源 CRUD 操作: Create（创建）, Read（读取）, Update（更新）, Delete（删除）
  - 基于 Schema 的资源和数据源定义
  - 通过 Terraform SDK 进行状态管理
- **服务层模式**: 每个腾讯云服务在 `tencentcloud/services/` 下有自己的包
  - 资源和数据源按服务分组
  - 共享服务客户端在 `service_tencentcloud_<service>.go` 文件中
- **错误处理**: 
  - 使用 `defer tccommon.LogElapsed()` 记录操作耗时
  - 使用 `defer tccommon.InconsistentCheck()` 进行状态一致性检查
  - 通过 `helper.Retry()` 实现最终一致性的重试逻辑
- **常用模式**:
  - 资源 ID 格式: 使用分隔符的复合 ID（例如：`instanceId#userId`）
  - 使用 `helper` 包提供通用工具函数
  - 敏感字段在 schema 中标记 `Sensitive: true`

### 测试策略
- **验收测试**: 使用 `TF_ACC=1` 环境变量启用
  - 测试文件命名: `*_test.go`
  - 测试函数命名: `TestAccTencentCloud<Service>_<scenario>`
  - 测试超时: 默认 120 分钟
- **测试要求**:
  - 设置 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY` 环境变量
  - 设置 `TENCENTCLOUD_APPID` 用于 COS 存储桶测试
- **清理测试**: 基础设施清理测试 (`make sweep`)
- **单元测试**: 标准 Go 测试，30 秒超时，4 路并行
- **测试命令**:
  - `make test`: 运行单元测试
  - `make testacc`: 运行验收测试
  - `make sweep`: 清理测试资源

### Git 工作流
- **分支命名**: 功能分支命名（例如：`001-rabbitmq-user-permissions`）
- **提交规范**: 标准 Git 提交信息
- **变更日志**: 在 `.changelog/` 目录中维护变更日志条目
  - 格式: `<PR_NUMBER>.txt` 文件
- **开发流程**:
  1. 创建功能分支
  2. 实现变更
  3. 运行 `make fmt` 进行代码格式化
  4. 运行 `make lint` 进行代码检查
  5. 运行 `make doc` 生成文档
  6. 运行测试
  7. 提交并创建 PR

## 领域上下文
- **腾讯云服务**: Provider 覆盖 80+ 个腾讯云服务，包括:
  - 计算: CVM（云服务器）, AS（弹性伸缩）, Lighthouse（轻量应用服务器）
  - 存储: CBS（云硬盘）, CFS（文件存储）, COS（对象存储）
  - 数据库: CDB（云数据库 MySQL）, PostgreSQL, MongoDB, Redis, TencentDB 系列
  - 网络: VPC（私有网络）, CLB（负载均衡）, CCN（云联网）, DC（专线接入）, VPN
  - 容器: TKE（容器服务）, TCR（容器镜像服务）, TCM（服务网格）
  - 中间件: TDMQ (Pulsar/RabbitMQ/RocketMQ), CKafka（消息队列）
  - 安全: CAM（访问管理）, CFW（云防火墙）, WAF（Web 应用防火墙）, SSL（证书）
  - 监控: Monitor（云监控）, RUM（前端性能监控）, CLS（日志服务）
  - CDN（内容分发）, API Gateway（API 网关）, Serverless（无服务器）等更多服务
- **资源生命周期**: 资源遵循 Terraform 生命周期（创建 → 读取 → 更新 → 删除 → 导入）
- **API 集成**: 通过 SDK 直接集成腾讯云 OpenAPI
- **重试与一致性**: 由于云 API 的最终一致性，许多操作需要重试逻辑

## 重要约束
- **Go 版本**: 必须使用 Go 1.17+（在 go.mod 中指定）
- **Terraform 版本**: 兼容 Terraform 0.13.x 及以上版本
- **Vendor 模式**: 依赖项通过 vendor 目录管理
- **API 速率限制**: 腾讯云 API 速率限制可能影响操作
- **向后兼容性**: 必须保持与现有 Terraform 配置的兼容性
- **资源状态**: 资源必须正确处理 Terraform 状态管理
- **SDK 版本**: 必须保持腾讯云 SDK 包版本同步
- **文档**: 所有资源和数据源必须在 `website/docs/` 中有相应文档

## 外部依赖
- **腾讯云 OpenAPI**: 所有云操作的主要依赖
  - 通过 Secret ID/Key 或 STS Token 认证
  - 区域端点服务访问
  - API 版本控制（例如：tdmq/v20200217）
- **Terraform Registry**: Provider 发布到 Terraform Registry `registry.terraform.io/tencentcloudstack/tencentcloud`
- **开发工具**:
  - golangci-lint: 代码质量和检查
  - tfproviderlint: Terraform 特定检查
  - misspell: 拼写检查
  - terrafmt: 文档中的 HCL 格式化
- **文档生成**: 自定义 `gendoc` 工具用于生成 Provider 文档
- **CI/CD**: Travis CI 配置用于自动化测试和发布
