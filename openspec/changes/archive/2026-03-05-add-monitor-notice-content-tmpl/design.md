# Design: Monitor Notice Content Template Resource

## Context
腾讯云监控服务提供了通知内容模板功能，允许用户自定义告警通知的内容格式。该功能通过 4 个 API 接口实现完整的 CRUD 操作。本设计需要将这些 API 映射为 Terraform Provider 的标准资源类型。

**背景约束:**
- 创建接口返回 TmplID，但查询和修改需要同时提供 TmplID 和 TmplName
- 删除接口仅需要 TmplID
- 模板内容结构复杂，包含多种通知渠道的嵌套配置
- 需要遵循项目现有的架构模式和编码规范

## Goals / Non-Goals

**Goals:**
- 实现完整的 CRUD 操作，支持用户通过 Terraform 管理通知内容模板
- 正确处理复合 ID，确保资源的唯一标识和可追溯性
- 支持所有主要通知渠道的配置（邮件、短信、企业微信、钉钉、飞书等）
- 提供清晰的错误处理和日志记录
- 遵循项目的代码风格和架构模式

**Non-Goals:**
- 不实现模板内容的自动验证（由云 API 负责）
- 不支持批量操作（单个资源管理）
- 不处理模板绑定告警策略的逻辑（属于其他资源的职责）

## Decisions

### Decision 1: ID 格式 (已优化)
**选择:** 使用单一 `tmplID` 作为资源 ID 格式

**理由:**
- CreateNoticeContentTmpl 返回的 TmplID 已经是唯一标识符
- DescribeNoticeContentTmpl API 接受 TmplIDs 列表参数进行查询,TmplName 仅用于过滤
- 单一 ID 更简洁,符合 Terraform 资源 ID 最佳实践
- 避免了复合 ID 解析的复杂性和潜在错误

**实现:**
```go
// Create 操作
d.SetId(*response.Response.TmplID)

// Read/Update/Delete 操作
tmplID := d.Id()
request.TmplIDs = []*string{&tmplID}
```

**变更原因:**
- 初始设计过度复杂,TmplID 本身足以唯一标识资源
- 简化了 Update 和 Delete 操作的 ID 处理逻辑
- 提高了代码可维护性

### Decision 2: tmpl_contents Schema 结构设计
**选择:** 使用 TypeMap + Elem TypeString 的嵌套结构表示复杂的模板内容

**理由:**
- 腾讯云 API 接受 JSON 格式的 NoticeContentTmplItem 对象
- 不同渠道的配置结构差异较大，过度细化 Schema 会导致维护成本高
- 使用 Map 结构保持灵活性，同时在文档中提供清晰的结构说明
- 可以直接序列化为 API 需要的 JSON 格式

**备选方案考虑:**
- **方案 A:** 完全定义每个渠道的嵌套 Schema
  - 优点: 类型安全，IDE 提示友好
  - 缺点: 维护成本高，API 结构变化时需要同步更新
- **方案 B (选中):** 使用 Map 结构，在文档中说明格式
  - 优点: 灵活性高，易于维护
  - 缺点: 类型检查较弱，依赖用户按文档配置

### Decision 3: DescribeNoticeContentTmpl 参数映射
**选择:** 将 tmplID 转换为字符串列表传递给 TmplIDs 参数

**理由:**
- API 设计上 TmplIDs 是数组类型，支持批量查询
- 单个资源读取时传递单元素列表
- TmplID 本身已经是唯一标识，足以精确查询

**实现:**
```go
tmplID := d.Id()
request.TmplIDs = []*string{&tmplID}
// TmplName 可选，不强制传递
```

**实现:**
```go
request.TmplIDs = []*string{&tmplID}
request.TmplName = &tmplName
```

### Decision 4: 重试机制实现
**选择:** 参考项目现有模式（igtm service），将 ratelimit.Check、判空逻辑和成功日志放在重试函数内部

**理由:**
- 符合项目统一的重试模式
- ratelimit.Check 在每次重试前执行，避免超出 API 限流
- 响应为空时返回 NonRetryableError，避免无效重试
- 成功日志在重试函数内部记录，确保日志完整

**实现模式:**
```go
err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    ratelimit.Check(request.GetAction())
    result, e := client.DescribeNoticeContentTmpl(request)
    if e != nil {
        return tccommon.RetryError(e)
    }
    log.Printf("[DEBUG]%s api[%s] success...", logId, request.GetAction())
    
    if result == nil || result.Response == nil {
        return resource.NonRetryableError(fmt.Errorf("Response is nil."))
    }
    
    response = result
    return nil
})
```

### Decision 5: 文件组织
**选择:** 在 `tencentcloud/services/monitor/` 目录下创建新文件

**文件列表:**
- `resource_tc_monitor_notice_content_tmpl.go` - 资源实现
- `resource_tc_monitor_notice_content_tmpl_test.go` - 验收测试
- `service_tencentcloud_monitor.go` - 扩展 service 方法（已存在文件）

**理由:**
- 符合项目按服务组织的架构模式
- monitor 服务目录已存在，直接添加新文件
- 保持资源、测试、service 的清晰分离

## Risks / Trade-offs

### Risk 1: tmpl_contents 结构复杂度
**风险:** 用户配置复杂嵌套结构时容易出错

**缓解措施:**
- 在文档中提供详细的结构说明和完整示例
- 在 examples 目录提供多个典型场景的配置
- 错误信息中包含 API 返回的详细错误，帮助用户定位问题

### Risk 2: API 字段变化
**风险:** 腾讯云 API 可能新增或修改通知渠道配置字段

**缓解措施:**
- 使用完整的 Schema 定义，新增字段时只需扩展 Schema
- 定期检查 API 文档更新，及时更新资源文档
- 验收测试覆盖主要渠道，确保向后兼容性

### Risk 3: 资源漂移问题 (已解决)
**风险:** flatten 函数设置空对象导致 Terraform 检测到资源漂移

**解决方案:**
- 优化 flattenNoticeContentTmplItem 函数，只在有值时设置 template 子字段
- 所有 template 对象（email、sms、qywx 等）先收集到 map，检查 `len(map) > 0` 后才添加
- 避免空对象被写入 state，消除不必要的 diff
- 经过实际测试验证，资源创建后不再报告漂移

### Trade-off: Schema 灵活性 vs 类型安全
**选择:** 优先考虑灵活性，使用 Map 结构

**理由:**
- Terraform Provider 主要用于基础设施管理，配置稳定性高
- 用户通常参考文档和示例进行配置，不会频繁修改
- API 变化时无需修改代码，降低维护成本
- 可以通过文档和示例降低配置错误风险

## Migration Plan
本次为新增资源，无需迁移计划。

**部署步骤:**
1. 合并代码到 main 分支
2. 打包新版本 Provider
3. 发布到 Terraform Registry
4. 更新官方文档站点

**回滚方案:**
- 如果发现严重问题，可以在 Provider 中移除资源注册
- 已创建的资源可以通过云控制台或 API 手动管理
- 不会影响其他现有资源的使用

## Open Questions
- Q: 是否需要支持 Import 操作？
  - A: 建议支持，通过 `{tmplID}#{tmplName}` 格式导入现有资源
  - 实现: 添加 Importer 配置，复用 Read 逻辑

- Q: tmpl_contents 中 Base64 编码字段是否需要自动处理？
  - A: 建议在文档中说明需要用户手动进行 Base64 编码
  - 理由: 保持与 API 行为一致，避免引入额外复杂度

- Q: 是否需要支持查询所有模板的 DataSource？
  - A: 当前 change 聚焦于 Resource 实现
  - 建议后续单独创建 DataSource，使用 DescribeNoticeContentTmpl 的列表查询能力
