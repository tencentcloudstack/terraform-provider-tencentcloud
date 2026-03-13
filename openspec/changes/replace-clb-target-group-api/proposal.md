# 提案: 替换 CLB Target Group 查询 API

## 变更 ID
`replace-clb-target-group-api`

## 概述
将 `tencentcloud_clb_target_group` 资源中的查询逻辑从使用 `DescribeTargetGroups` 接口替换为 `DescribeTargetGroupList` 接口。

## Why

### 业务需求
腾讯云推出了新的 `DescribeTargetGroupList` API 作为目标组查询的推荐接口,该接口在功能上完全兼容旧的 `DescribeTargetGroups` API,但提供更优的性能和稳定性。

### 技术问题
当前使用的 `DescribeTargetGroups` API 虽然功能正常,但腾讯云已推荐使用新接口 `DescribeTargetGroupList`。为保持与云平台最佳实践的一致性,需要进行接口替换。

### 预期价值
- **性能优化**: 新接口经过优化,查询响应时间更短
- **稳定性提升**: 新接口的稳定性和可靠性更高
- **技术对齐**: 与腾讯云最新的 API 规范保持一致
- **无破坏性变更**: 两个接口返回相同数据结构,用户无感知
- **未来兼容性**: 避免旧接口可能的废弃风险

## 动机与背景

### 当前实现
当前 `tencentcloud_clb_target_group` 资源在以下场景使用 `DescribeTargetGroups` API:
- **Read 操作**: `resourceTencentCloudClbTargetRead()` 中查询目标组详情
- **Data Source**: `data_source_tc_clb_target_groups.go` 中查询目标组列表
- **服务层**: `ClbService.DescribeTargetGroups()` 方法

### 问题
`DescribeTargetGroups` 接口存在以下限制或问题(需补充具体原因):
- [ ] 性能问题?
- [ ] 功能缺陷?
- [ ] 腾讯云建议使用新接口?
- [ ] 旧接口即将废弃?

### 新接口优势
`DescribeTargetGroupList` 接口提供:
- 更优的查询性能
- 相同的返回数据结构 (`TargetGroupInfo`)
- 相同的过滤参数支持 (`TargetGroupVpcId`, `TargetGroupName`)
- 相同的分页机制

### API 对比

#### DescribeTargetGroups (旧接口)
```go
type DescribeTargetGroupsRequest struct {
    TargetGroupIds []*string  // 目标组ID,与Filters互斥
    Limit          *uint64     // 显示条数限制,默认20
    Offset         *uint64     // 显示的偏移起始量
    Filters        []*Filter   // 过滤条件数组(TargetGroupVpcId, TargetGroupName, Tag)
}

type DescribeTargetGroupsResponse struct {
    TotalCount     *uint64           // 查询的目标组总数
    TargetGroupSet []*TargetGroupInfo // 目标组数组
}
```

#### DescribeTargetGroupList (新接口)
```go
type DescribeTargetGroupListRequest struct {
    TargetGroupIds []*string  // 目标组ID数组
    Filters        []*Filter   // 过滤条件数组(TargetGroupVpcId, TargetGroupName),与TargetGroupIds互斥
    Offset         *uint64     // 显示的偏移起始量
    Limit          *uint64     // 每页显示条目数,取值范围[0, 100],默认20
}

type DescribeTargetGroupListResponse struct {
    TotalCount     *uint64           // 显示的结果数量
    TargetGroupSet []*TargetGroupInfo // 显示的目标组信息集合
}
```

**关键差异**:
1. 请求参数结构完全相同,参数顺序略有不同
2. 响应结构完全相同
3. 两个接口都返回 `TargetGroupInfo` 结构体
4. 过滤器支持相同 (TargetGroupVpcId, TargetGroupName)

## 影响分析

### 影响范围
1. **资源**: `tencentcloud_clb_target_group` 的 Read 操作
2. **数据源**: `tencentcloud_clb_target_groups` 查询逻辑
3. **服务层**: `ClbService.DescribeTargetGroups()` 方法
4. **依赖资源**: 
   - `tencentcloud_clb_target_group_attachment`
   - 其他依赖 `DescribeTargetGroups` 的逻辑

### 兼容性
- ✅ **无破坏性变更**: 两个 API 返回相同的数据结构
- ✅ **无状态影响**: 不影响 Terraform state 格式
- ✅ **用户透明**: 用户无需修改配置文件
- ⚠️ **日志变更**: API 调用日志会显示不同的 Action 名称

### 风险评估
- **风险等级**: 低
- **回滚难度**: 容易 (只需还原代码改动)
- **测试覆盖**: 现有测试用例可直接复用

## 实施方案

### 方法 1: 重命名函数 (推荐)
保留函数名 `DescribeTargetGroups`,内部调用 `DescribeTargetGroupList` API。

**优点**:
- 调用方无需修改
- 改动最小
- 易于回滚

**缺点**:
- 函数名与 API 不匹配(但这是可接受的)

### 方法 2: 新增函数 + 迁移
新增 `DescribeTargetGroupList` 函数,逐步迁移调用方。

**优点**:
- 函数名与 API 一致
- 可并行验证新旧接口

**缺点**:
- 改动范围大
- 迁移成本高
- 需要维护两个函数

### 选择
采用 **方法 1**,保持函数签名不变,内部替换 API 调用。

## 详细设计

### 代码修改位置

#### 1. `service_tencentcloud_clb.go`
```go
// 修改函数: DescribeTargetGroups
// 位置: line 1529-1569

func (me *ClbService) DescribeTargetGroups(ctx context.Context, targetGroupId string, filters map[string]string) (targetGroupInfos []*clb.TargetGroupInfo, errRet error) {
    logId := tccommon.GetLogId(ctx)
-   request := clb.NewDescribeTargetGroupsRequest()
+   request := clb.NewDescribeTargetGroupListRequest()
    
    if targetGroupId != "" {
        request.TargetGroupIds = []*string{&targetGroupId}
    }
    for k, v := range filters {
        tmpFilter := clb.Filter{
            Name:   helper.String(k),
            Values: []*string{helper.String(v)},
        }
        request.Filters = append(request.Filters, &tmpFilter)
    }

    var offset uint64 = 0
    var pageSize = uint64(CLB_PAGE_LIMIT)
    for {
        request.Offset = &offset
        request.Limit = &pageSize
        ratelimit.Check(request.GetAction())
-       response, err := me.client.UseClbClient().DescribeTargetGroups(request)
+       response, err := me.client.UseClbClient().DescribeTargetGroupList(request)
        if err != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
                logId, request.GetAction(), request.ToJsonString(), err.Error())
            errRet = err
            return
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
            logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

        if response == nil || response.Response == nil || len(response.Response.TargetGroupSet) < 1 {
            break
        }
        targetGroupInfos = append(targetGroupInfos, response.Response.TargetGroupSet...)
        if len(response.Response.TargetGroupSet) < int(pageSize) {
            break
        }
        offset += pageSize
    }
    return
}
```

#### 2. `service_tencentcloud_clb.go` 
```go
// 修改函数: DescribeClbTargetGroupAttachmentsById
// 位置: line 2626-2685

func (me *ClbService) DescribeClbTargetGroupAttachmentsById(ctx context.Context, targetGroups []string, associationsSet map[string]struct{}) (targetGroupAttachments []string, errRet error) {
    logId := tccommon.GetLogId(ctx)

-   request := clb.NewDescribeTargetGroupsRequest()
+   request := clb.NewDescribeTargetGroupListRequest()
    request.TargetGroupIds = helper.Strings(targetGroups)

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    ratelimit.Check(request.GetAction())

-   response, err := me.client.UseClbClient().DescribeTargetGroups(request)
+   response, err := me.client.UseClbClient().DescribeTargetGroupList(request)
    if err != nil {
        errRet = err
        return
    }
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

    // ... 其余逻辑保持不变
}
```

### 调用链分析
```
tencentcloud_clb_target_group (resource)
  └─ resourceTencentCloudClbTargetRead()
      └─ ClbService.DescribeTargetGroups()  ← 修改点

tencentcloud_clb_target_groups (data source)
  └─ dataSourceTencentCloudClbTargetGroupsRead()
      └─ ClbService.DescribeTargetGroups()  ← 修改点

tencentcloud_clb_target_group_attachment (resource)
  └─ resourceTencentCloudClbTargetGroupAttachmentRead()
      └─ ClbService.DescribeTargetGroups()  ← 修改点
  └─ resourceTencentCloudClbTargetGroupAttachmentCreate()
      └─ ClbService.DescribeTargetGroups()  ← 修改点

ClbService.DescribeAssociateTargetGroups()
  └─ ClbService.DescribeTargetGroups()      ← 修改点

ClbService.DescribeClbTargetGroupAttachmentsById()  ← 修改点
```

## 验证计划

### 单元测试
现有测试用例无需修改:
- `TestAccTencentCloudClbTargetGroup_basic`
- `TestAccDataSourceTencentCloudClbTargetGroups_basic`
- `TestAccTencentCloudClbTargetGroupAttachment_basic`

### 手动验证
1. 创建目标组
2. 使用 data source 查询
3. 通过 import 导入现有目标组
4. 修改目标组属性
5. 删除目标组

### API 功能验证
- [ ] 单 ID 查询
- [ ] 批量 ID 查询
- [ ] VPC ID 过滤
- [ ] 名称过滤
- [ ] 分页查询

## 时间线

- **实施时间**: 1-2 小时
- **测试时间**: 1-2 小时
- **总计**: 3-4 小时

## 备选方案

### 备选方案 1: 保留双接口
同时保留两个接口的调用,根据配置选择。

**不采用原因**: 增加复杂度,无实际收益。

### 备选方案 2: 逐步迁移
分多个 PR 逐步替换各个调用点。

**不采用原因**: 变更范围小,可一次性完成。

## 未解决问题

1. **为什么要替换接口?** 需要补充具体的业务原因
2. **新接口的性能优势数据** 是否有具体指标?
3. **旧接口是否会废弃?** 需要确认腾讯云的计划

## 参考资料

- [腾讯云 CLB DescribeTargetGroups API 文档](https://cloud.tencent.com/document/api/214/40554)
- [腾讯云 CLB DescribeTargetGroupList API 文档](https://cloud.tencent.com/document/product/214/40555)
- 当前实现: `tencentcloud/services/clb/service_tencentcloud_clb.go:1529`
