## 1. Schema 定义

- [x] 1.1 在 `resource_tc_tcr_namespace.go` 的资源 schema 中添加 `tags` 字段
  - 类型: `schema.TypeMap`,元素类型为 `schema.TypeString`
  - 设置 `Optional: true` 和 `Computed: true`
  - 添加标签键值对的描述
- [x] 1.2 导入 `svctag` 包: `svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"`

## 2. Create 函数增强

- [x] 2.1 在 `resourceTencentCloudTcrNamespaceCreate` 中从配置读取 `tags`
- [x] 2.2 将标签传递给 `CreateTCRNameSpace` 服务函数(如果 API 不支持则在创建后处理)
- [x] 2.3 如果创建 API 不支持标签,则在命名空间创建后使用 `TagService.ModifyTags`

## 3. Update 函数实现

- [x] 3.1 在 `resourceTencentCloudTcrNamespaceUpdate` 中添加标签变更检测
  ```go
  if d.HasChange("tags") {
      // 标签更新逻辑
  }
  ```
- [x] 3.2 使用 `d.GetChange("tags")` 获取新旧标签值
- [x] 3.3 使用 `svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))` 计算标签差异
- [x] 3.4 创建 TagService 实例: `tagService := svctag.NewTagService(tcClient)`
- [x] 3.5 构建资源名称: `resourceName := tccommon.BuildTagResourceName("tcr", "namespace", tcClient.Region, instanceId+"/"+namespaceName)`
- [x] 3.6 调用 `tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)`
- [x] 3.7 处理错误,如果标签更新失败则返回

## 4. Read 函数增强

- [x] 4.1 检查 TCR namespace 响应是否包含标签
- [x] 4.2 如果响应中有标签,将其扁平化为 map 格式
- [x] 4.3 使用 `d.Set("tags", tagMap)` 在 state 中设置标签
- [x] 4.4 如果 namespace 响应中没有标签,单独使用 TagService 查询标签
  - 如需要使用 `tagService.DescribeResourceTags` 获取标签
  - 将标签列表转换为 map 格式

## 5. 文档

- [x] 5.1 在 `resource_tc_tcr_namespace.md` 的 Argument Reference 部分添加 tags 字段
  - 标注为 `Optional` 和 `Computed`
  - 说明标签键值对格式
- [x] 5.2 在文档中添加包含标签的使用示例
  ```hcl
  resource "tencentcloud_tcr_namespace" "example" {
    instance_id = "tcr-xxx"
    name        = "example-ns"
    is_public   = false
    
    tags = {
      "env"       = "production"
      "createdBy" = "terraform"
    }
  }
  ```
- [x] 5.3 更新现有示例,可选地包含标签

## 6. 测试

- [x] 6.1 验证 schema 编译: `go build -o /dev/null ./tencentcloud/services/tcr/`
- [x] 6.2 运行 linter 检查: 检查 `resource_tc_tcr_namespace.go` 中的警告
- [ ] 6.3 手动测试:
  - 创建带标签的 namespace
  - 更新标签(添加、修改、删除)
  - 验证标签正确应用
  - 导入现有 namespace 并验证标签被正确读取
- [ ] 6.4 验证标签更新不影响 namespace 的其他属性

## 7. Changelog

- [x] 7.1 在 `.changelog/` 目录创建 changelog 条目
- [x] 7.2 将标签支持记录为增强功能:
  ```
  ```release-note:enhancement
  resource/tencentcloud_tcr_namespace: support `tags` parameter for tag management
  ```
  ```

## 验证清单

- [x] 代码编译无错误
- [x] Linter 通过无关键警告
- [ ] 可在 namespace 创建时设置标签
- [ ] 可在 namespace 创建后更新标签
- [ ] 可删除标签
- [x] 文档包含标签示例
- [x] 对现有功能无破坏性变更
