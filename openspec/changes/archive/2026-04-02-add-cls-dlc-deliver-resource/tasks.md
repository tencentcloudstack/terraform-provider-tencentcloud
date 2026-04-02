# 任务清单：add-cls-dlc-deliver-resource

## 1. 新增 service 层方法

**文件**: `tencentcloud/services/cls/service_tencentcloud_cls.go`

- [x] 在文件末尾新增 `DescribeClsDlcDeliverById(ctx, topicId, taskId string) (*cls.DlcDeliverInfo, error)` 方法
- [x] 执行 `go fmt ./tencentcloud/services/cls/`

---

## 2. 新增资源主文件

**文件**: `tencentcloud/services/cls/resource_tc_cls_dlc_deliver.go`

- [x] 实现 `ResourceTencentCloudClsDlcDeliver()` 函数及完整 Schema（所有字段 ForceNew）
- [x] 实现 `resourceTencentCloudClsDlcDeliverCreate`
- [x] 实现 `resourceTencentCloudClsDlcDeliverRead`
- [x] 实现 `resourceTencentCloudClsDlcDeliverDelete`
- [x] 执行 `go fmt ./tencentcloud/services/cls/`

---

## 3. 新增资源文档示例

**文件**: `tencentcloud/services/cls/resource_tc_cls_dlc_deliver.md`

- [x] 编写资源示例 HCL，包含完整的 dlc_info 嵌套结构

---

## 4. 新增单元测试文件

**文件**: `tencentcloud/services/cls/resource_tc_cls_dlc_deliver_test.go`

- [x] 实现 `TestAccTencentCloudClsDlcDeliverResource_basic`
- [x] 编写 `testAccClsDlcDeliver` HCL 常量

---

## 5. 注册资源

**文件**: `tencentcloud/provider.go`

- [x] 在 `ResourcesMap` 中 `tencentcloud_cls_scheduled_sql` 后追加 `"tencentcloud_cls_dlc_deliver": cls.ResourceTencentCloudClsDlcDeliver()`

---

## 6. 编译验证

- [x] `go build ./tencentcloud/services/cls/` 确认编译通过
- [x] `go build ./tencentcloud/` 确认 provider 编译通过

---

## 7. 补充 Update 模块（ModifyDlcDeliver）

**文件**: `tencentcloud/services/cls/resource_tc_cls_dlc_deliver.go`

- [x] `ResourceTencentCloudClsDlcDeliver()` 中添加 `Update` 函数引用
- [x] 所有可变字段去掉 `ForceNew`（`topic_id` 保留）
- [x] 新增 `status` 字段（Optional, Computed）：任务状态，1=运行，2=停止
- [x] 实现 `resourceTencentCloudClsDlcDeliverUpdate`，调用 `ModifyDlcDeliverWithContext`
- [x] Read 函数中补充 `status` 字段的读取
- [x] 执行 `go fmt` 及编译验证通过

---

## 总结

- **状态**: 🎉 所有任务已完成
- **风险等级**：低（纯新增，不影响现有资源）
- **破坏性变更**：无
