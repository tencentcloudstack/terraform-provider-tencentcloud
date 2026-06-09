# 任务清单：fix-tke-attachment-describe-cluster-instances

## 1. 新增 DescribeAllClusterInstances 方法

**文件**: `tencentcloud/services/tke/service_tencentcloud_tke.go`

- [x] 在 `DescribeClusterInstancesByRole` 函数之后新增 `DescribeAllClusterInstances(ctx, id string) (masters, workers []InstanceInfo, errRet error)` 方法
  - 采用 for 循环（参考 `DescribeClusterInstancesByRole` 风格，不使用 goto）
  - **每次分页请求**用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹
  - 分页参数：limit=100（比原来的20更高效），offset 递增
  - 退出条件：`len(response.Response.InstanceSet) < int(limit)` 或 `offset >= total`
  - 内部逻辑（has 去重、InstanceRole 分类）与原 `DescribeClusterInstances` 保持一致
- [x] 执行 `go fmt ./tencentcloud/services/tke/`

---

## 2. 新增 DescribeClusterInstanceById 方法

**文件**: `tencentcloud/services/tke/service_tencentcloud_tke.go`

- [x] 在 `DescribeAllClusterInstances` 函数之后新增 `DescribeClusterInstanceById(ctx, clusterId, instanceId string) (ret *InstanceInfo, errRet error)` 方法
  - 构建 request，设置 `ClusterId` 和 `InstanceIds = []*string{&instanceId}`
  - 用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹单次查询
  - 若 `response.Response.InstanceSet` 为空则返回 `nil, nil`（表示未找到）
  - 找到则构建 `InstanceInfo` 并返回
- [x] 执行 `go fmt ./tencentcloud/services/tke/`

---

## 3. 更新 nodeHasAttachedToCluster

**文件**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_attachment_extension.go`

- [x] 将 `nodeHasAttachedToCluster` 中的 `service.DescribeClusterInstances` 调用改为 `service.DescribeAllClusterInstances`
  - 去除原有的"先调用一次，失败后再 retry"的双重调用逻辑，直接调用 `DescribeAllClusterInstances`（该方法内部已包含 retry）
- [x] 执行 `go fmt ./tencentcloud/services/tke/`

---

## 4. 更新 CreatePostHandleResponse0 等待逻辑

**文件**: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_attachment_extension.go`

- [x] 将 `resourceTencentCloudKubernetesClusterAttachmentCreatePostHandleResponse0` 中等待 TKE init 的 Retry 块内调用由 `tkeService.DescribeClusterInstances(ctx, clusterId)` 改为 `tkeService.DescribeClusterInstanceById(ctx, clusterId, instanceId)`
- [x] 相应调整判断逻辑：不再需要遍历 workers 列表，直接对返回的单个 `*InstanceInfo` 进行判断：
  - 若 `ret == nil`：返回 `resource.NonRetryableError`（实例不在集群中）
  - 若 `ret.InstanceState == "failed"`：返回 `resource.NonRetryableError`（附加 FailedReason）
  - 若 `ret.InstanceState != "running"`：返回 `resource.RetryableError`
  - 否则返回 `nil`
- [x] 执行 `go fmt ./tencentcloud/services/tke/`

---

## 5. 编译验证

- [x] `go build ./tencentcloud/services/tke/` 确认编译通过

---

## 总结

- **预计工作量**：中等（约 1 小时）
- **风险等级**：低（原有逻辑完全保留，只是调用方式重构）
- **破坏性变更**：无
- **状态**: 已完成
