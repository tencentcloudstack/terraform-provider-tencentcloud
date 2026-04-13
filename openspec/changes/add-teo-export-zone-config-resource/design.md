## Context

当前项目使用 Terraform Plugin SDK v2 为腾讯云服务提供 Terraform Provider，采用标准的文件组织结构和服务层模式。TEO (Tencent Cloud EdgeOne) 服务已存在基础支持，但缺少站点配置导出功能的 Resource 实现。该功能需要通过 CAPI 接口调用实现配置导出，这是一个标准的 CRUD 资源实现任务。

## Goals / Non-Goals

**Goals:**
- 实现完整的 `tencentcloud_teo_export_zone_config` Resource，包含标准的 CRUD 操作
- 遵循项目现有的代码规范和文件组织结构
- 确保资源与 CAPI 接口的参数定义完全一致
- 提供完整的单元测试和验收测试覆盖
- 生成清晰的使用文档

**Non-Goals:**
- 不涉及现有资源的 Schema 修改（保持向后兼容）
- 不引入新的外部依赖或架构模式
- 不涉及复杂的迁移或数据变更

## Decisions

### 1. 实现模式
采用项目现有的标准资源实现模式：
- 主文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go`
- 测试文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
- 文档文件：`website/docs/r/teo_export_zone_config.html.md`
- 样例文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.md`

**理由**：与项目中已有的 TEO 资源保持一致，便于维护和理解。

### 2. API 调用方式
通过 `tencentcloud-sdk-go` 的 teo 包调用 CAPI 接口，使用项目的标准服务层辅助函数（如 `logId`、`returnCode` 等）。

**理由**：项目已经建立了完善的 SDK 集成和错误处理机制，复用现有代码可以减少维护成本。

### 3. 错误处理和重试
使用 `helper.Retry()` 实现最终一致性重试，在 CRUD 函数中使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 进行错误处理。

**理由**：TEO 服务的异步特性需要重试机制确保最终一致性，项目已有成熟的重试模式。

### 4. Timeouts 支持
在 Schema 中声明 Timeouts 块，支持 Create、Read、Update、Delete 操作的超时配置。

**理由**：异步操作可能需要较长时间，合理的超时配置可以提升用户体验。

### 5. 资源 ID 策略
根据 CAPI 接口的主键设计资源 ID，如果存在复合主键则使用分隔符（如 `zoneId#configId`）。

**理由**：确保资源在 Terraform state 中的唯一性和可识别性。

## Risks / Trade-offs

### 风险 1: CAPI 接口参数变更
**风险**：CAPI 接口参数可能在未来版本中发生变化，导致 Schema 不匹配。
**缓解措施**：严格按照当前版本的 CAPI 接口定义生成 Schema，参数使用 Optional 以兼容未来的扩展。

### 风险 2: 异步操作的最终一致性
**风险**：Create/Update/Delete 操作可能需要较长时间才能完成，状态查询可能延迟。
**缓解措施**：使用项目的重试机制，设置合理的默认超时时间，提供 Timeouts 配置选项。

### 风险 3: 测试环境依赖
**风险**：验收测试需要真实的腾讯云账号和环境变量，可能受到网络或配额限制。
**缓解措施**：单元测试不依赖真实环境，验收测试使用独立测试账号，提供清晰的测试运行说明。
