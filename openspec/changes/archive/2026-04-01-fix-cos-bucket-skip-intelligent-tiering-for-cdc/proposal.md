# 变更提案：修复 tencentcloud_cos_bucket 资源 BucketGetIntelligentTiering 在 CDC 场景下的调用问题

## 变更类型

**Bug 修复** — 此变更修复了现有功能的缺陷，不引入新能力，因此不需要 spec deltas。

## Why

### 问题描述

当 `tencentcloud_cos_bucket` 资源的 `cdc_id` 字段不为空（即用户在 CDC 专属云环境中使用 COS）时，`BucketGetIntelligentTiering` 接口不被 CDC 环境支持，会返回错误或异常响应，导致 Read 操作失败。

### 根本原因

在 `resourceTencentCloudCosBucketRead` 函数中，调用 `cosService.BucketGetIntelligentTiering` 时没有对 `cdcId` 字段进行判断。CDC（Cloud Dedicated Cluster，专属集群）环境中的 COS Bucket 不支持智能分层存储（Intelligent Tiering）功能，强行调用该接口会导致错误。

相关代码位置：`tencentcloud/services/cos/service_tencentcloud_cos.go`，`BucketGetIntelligentTiering` 函数（line 1801）。

### 影响

- **CDC 用户受影响**：使用专属集群的用户在执行 `terraform plan` 或 `terraform refresh` 时，Read 操作会因为调用不支持的接口而报错
- **功能不可用**：CDC 场景下 `tencentcloud_cos_bucket` 资源的 Read 操作不能正常完成

## What Changes

### 代码变更

在 `tencentcloud/services/cos/service_tencentcloud_cos.go` 中，`BucketGetIntelligentTiering` 函数添加对 `cdcId` 的判断：

#### 修改 BucketGetIntelligentTiering 函数

**修改前：**
```go
func (me *CosService) BucketGetIntelligentTiering(ctx context.Context, bucket string, cdcId string) (result *cos.BucketGetIntelligentTieringResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	ratelimit.Check("GetIntelligentTiering")
	intelligentTieringResult, response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.GetIntelligentTiering(ctx)
	// ...
}
```

**修改后：**
```go
func (me *CosService) BucketGetIntelligentTiering(ctx context.Context, bucket string, cdcId string) (result *cos.BucketGetIntelligentTieringResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	if cdcId != "" {
		return nil, nil
	}

	ratelimit.Check("GetIntelligentTiering")
	intelligentTieringResult, response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.GetIntelligentTiering(ctx)
	// ...
}
```

**理由**：CDC 专属集群不支持智能分层存储功能，当 `cdcId` 不为空时，直接返回 `nil, nil` 跳过该接口调用，避免错误。

### 影响范围

- **影响的规范**：无 — 这是一个 **bug 修复**，修复了 CDC 场景下的非法调用，不涉及新能力或行为变更
- **影响的文件**：
  - `tencentcloud/services/cos/service_tencentcloud_cos.go` — `BucketGetIntelligentTiering` 函数
- **破坏性变更**：无
  - 非 CDC 场景（`cdcId` 为空）行为完全不变
  - CDC 场景下，Intelligent Tiering 相关字段将保持零值/默认值，符合预期
- **迁移需求**：不适用

### 向后兼容性

✅ **完全向后兼容**：
- 现有的非 CDC Terraform 配置无需任何修改
- 只影响 `cdcId != ""` 的场景，且这些场景原本就无法正常工作
- 不影响 Create/Update/Delete 操作

### 测试建议

1. **非 CDC 场景**：确保非 CDC COS Bucket 的 Read 操作仍然正常调用 `BucketGetIntelligentTiering`
2. **CDC 场景**：在 `cdc_id` 不为空时，Read 操作跳过 `BucketGetIntelligentTiering` 调用，不报错
3. **编译验证**：运行 `go build` 确认代码编译通过

### 参考

- **相关函数**：`tencentcloud/services/cos/service_tencentcloud_cos.go` — `BucketGetIntelligentTiering` (line 1801)
- **调用位置**：`tencentcloud/services/cos/resource_tc_cos_bucket.go` (line 1107)
