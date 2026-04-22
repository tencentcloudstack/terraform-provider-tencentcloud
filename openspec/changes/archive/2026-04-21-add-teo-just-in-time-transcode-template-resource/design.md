## Context

当前 Terraform Provider for TencentCloud 已经支持 TEO 产品的多个资源，但缺少即时转码模板功能的支持。TEO 的即时转码模板功能允许用户配置视频流和音频流的转码参数，用于边缘节点上的实时转码处理。

现有的 TEO 资源实现模式已经成熟，包括：
- 使用 Terraform Plugin SDK v2 进行资源实现
- 通过 tencentcloud-sdk-go 的 teo/v20220901 包调用云 API
- 使用复合 ID 格式进行资源标识
- 实现异步操作的轮询机制
- 遵循最终一致性重试模式

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_just_in_time_transcode_template` 资源的完整 CRUD 操作
- 支持视频流和音频流的独立配置，包括开关控制和模板参数
- 正确处理异步操作，确保资源状态最终一致
- 提供完整的单元测试和文档
- 保持与现有 TEO 资源的一致性

**Non-Goals:**
- 不实现即时转码模板的更新操作（云 API 不支持 Update 接口，需要通过删除后重新创建实现）
- 不实现批量操作接口
- 不修改现有 TEO 资源的行为

## Decisions

### 1. 资源实现模式

**决策**: 采用标准的 Terraform 资源实现模式，参考现有 TEO 资源（如 `resource_tc_teo_acceleration_domain.go`）。

**理由**:
- 保持代码一致性和可维护性
- 利用已有的错误处理、重试等基础设施
- 降低实现复杂度

**替代方案考虑**:
- 使用新的资源框架（如 Terraform Plugin Framework）：考虑迁移成本和兼容性风险，暂不采用

### 2. 复合 ID 设计

**决策**: 使用 `zone_id#template_id` 作为资源 ID，与现有 TEO 资源保持一致。

**理由**:
- zone_id 是 TEO 的站点标识，template_id 是模板唯一标识
- 这种格式已在多个 TEO 资源中使用，符合用户习惯
- 便于资源的唯一性识别和查询

**实现细节**:
```
func resourceTencentCloudTeoJustInTimeTranscodeTemplateParseId(id string) (string, string, error) {
    items := strings.Split(id, "#")
    if len(items) != 2 {
        return "", "", fmt.Errorf("invalid resource ID format, expected: zone_id#template_id, got: %s", id)
    }
    return items[0], items[1], nil
}
```

### 3. 异步操作处理

**决策**: 对 Create 和 Delete 操作实现轮询机制，确保操作完成后资源状态同步。

**理由**:
- 云 API 是异步操作，创建和删除后模板不会立即生效
- 需要轮询 DescribeJustInTimeTranscodeTemplates 接口确认状态
- Terraform 的 `Timeouts` 配置允许用户自定义等待时间

**实现细节**:
```go
d.SetId(fmt.Sprintf("%s#%s", *response.TemplateId, zoneId))
d.Set("template_id", response.TemplateId)

// 轮询确认创建成功
err := resourceTencentCloudTeoJustInTimeTranscodeTemplateRead(ctx, d, meta)
if err != nil {
    return err
}
```

### 4. 更新操作处理

**决策**: 由于云 API 不提供 Update 接口，Update 操作将强制删除并重新创建资源（ForceNew）。

**理由**:
- 云 API 只支持 Create 和 Delete，不支持 Update
- 通过 ForceNew 确保 Terraform 状态正确更新
- 虽然会增加操作成本，但确保了功能正确性

**实现细节**:
```go
"zone_id": &schema.Schema{
    Type:     schema.TypeString,
    Required: true,
    ForceNew: true,
},
"template_name": &schema.Schema{
    Type:     schema.TypeString,
    Required: true,
    ForceNew: true,
},
// 其他关键字段也标记为 ForceNew
```

### 5. 视频模板和音频模板配置

**决策**: 将 VideoTemplateInfo 和 AudioTemplateInfo 作为嵌套对象在 schema 中定义，使用 TypeList 和 TypeMap 组合。

**理由**:
- VideoTemplateInfo 和 AudioTemplateInfo 是复杂的嵌套结构
- Terraform schema 可以通过嵌套表示复杂数据结构
- 保持与云 API 数据结构的一致性

**实现细节**:
```go
"video_template": &schema.Schema{
    Type:     schema.TypeList,
    MaxItems: 1,
    Optional: true,
    Computed: true,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "video_codec": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
            },
            "fps": &schema.Schema{
                Type:     schema.TypeInt,
                Optional: true,
            },
            // 其他视频参数
        },
    },
},
```

### 6. 错误处理和重试

**决策**: 使用项目已有的 `helper.Retry()` 和 `tccommon.InconsistentCheck()` 模式处理最终一致性和错误。

**理由**:
- TEO 服务是分布式系统，存在最终一致性
- 网络抖动和服务临时不可用需要重试机制
- 利用现有代码避免重复实现

**实现细节**:
```go
defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.create")()

err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    result, err := me.client.UseTeoClient(teoVersion).CreateJustInTimeTranscodeTemplate(request)
    if err != nil {
        if isRetryableError(err) {
            return resource.RetryableError(err)
        }
        return resource.NonRetryableError(err)
    }
    // 处理响应
    return nil
})
```

### 7. 单元测试策略

**决策**: 新资源使用 gomonkey mock 云 API 进行单元测试，不使用 Terraform 集成测试。

**理由**:
- 新资源的测试应该专注于业务逻辑
- 避免依赖真实的云环境（需要凭证和网络）
- gomonkey 可以 mock SDK 接口，测试资源代码逻辑

**实现细节**:
```go
func TestAccTencentCloudTeoJustInTimeTranscodeTemplate_basic(t *testing.T) {
    // 使用 gomonkey mock CreateJustInTimeTranscodeTemplate 接口
    // 测试创建、读取、删除流程
}
```

## Risks / Trade-offs

### Risk 1: API 兼容性变化
**风险**: 云 API 可能在未来版本中变更参数结构，导致代码不兼容。

**缓解措施**:
- 使用 vendor 模式管理依赖，锁定特定版本的 SDK
- 在 README 中注明支持的 SDK 版本
- 关注云 API 变更公告，及时更新

### Risk 2: 异步操作超时
**风险**: 轮询等待资源创建或删除可能超时，导致 Terraform 操作失败。

**缓解措施**:
- 配置合理的默认超时时间（如 10 分钟）
- 允许用户通过 Timeouts 参数自定义超时
- 在文档中说明异步操作特性

### Risk 3: 资源状态不一致
**风险**: 轮询过程中资源状态可能与实际状态不一致。

**缓解措施**:
- 使用 `tccommon.InconsistentCheck()` 检测和重试
- 在读取操作中实现状态一致性校验
- 记录详细日志便于问题排查

### Trade-off 1: 更新操作成本
**权衡**: 由于不支持 Update 接口，任何参数变更都需要删除重建，增加操作成本和时间。

**理由**: 云 API 限制，无法绕过。

### Trade-off 2: 测试覆盖率
**权衡**: 使用 mock 单元测试而非集成测试，可能无法完全覆盖真实场景。

**缓解措施**: 补充文档和使用说明，鼓励用户进行实际验证。

## Migration Plan

### 部署步骤

1. **代码审查和合并**
   - 确保 Pull Request 通过代码审查
   - 验证单元测试全部通过
   - 确认文档完整准确

2. **发布版本**
   - 打包新版本的 Terraform Provider
   - 更新 CHANGELOG
   - 发布到 Terraform Registry

3. **用户文档更新**
   - 在 website/docs/ 添加资源文档
   - 提供使用示例
   - 说明注意事项（如异步操作、ForceNew 字段）

### 回滚策略

- 如果发现严重问题，可以通过快速回滚移除该资源
- 用户状态不受影响，因为这是新增资源
- 保留代码便于问题修复和重新发布

## Open Questions

（无待解决的开放性问题）
