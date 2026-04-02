# 任务清单

本修复在 `resourceTencentCloudCosBucketRead` 函数中，通过检查 `cosBucketUrl` 是否包含 `.cos-cdc.` 关键字来跳过 CDC 不支持的 Intelligent Tiering 接口调用。

## 1. 修改 Read 函数中的 intelligent tiering 调用逻辑

**文件**: `tencentcloud/services/cos/resource_tc_cos_bucket.go`

- [x] 在调用 `BucketGetIntelligentTiering` 和 `BucketGetIntelligentTieringArchivingRuleList` 之前，添加 `cosBucketUrl` 包含 `.cos-cdc.` 的判断
- [x] 当 `strings.Contains(cosBucketUrl, ".cos-cdc.")` 为 true 时跳过这两个函数的调用
- [x] 清理遗留的调试 `fmt.Println` / `fmt.Printf` 代码

**修改内容**（最终实现）:
```go
//read intelligent tiering
if !strings.Contains(cosBucketUrl, ".cos-cdc.") {
    result, err := cosService.BucketGetIntelligentTiering(ctx, bucket, cdcId)
    // ...

    //read intelligent tiering archiving rule list
    respData, err := cosService.BucketGetIntelligentTieringArchivingRuleList(ctx, bucket, cdcId)
    // ...
}
```

---

## 2. 编译验证

- [x] 运行以下命令确认编译通过：
  ```bash
  go build ./tencentcloud/services/cos/
  ```
- [x] 确认无新增编译错误或警告

---

## 3. 格式化代码

- [x] 对修改的文件执行 `go fmt`：
  ```bash
  go fmt ./tencentcloud/services/cos/
  ```

---

## 总结

- **预计工作量**：极低
- **风险等级**：极低（仅对 CDC 场景短路，非 CDC 行为完全不变）
- **破坏性变更**：无
- **测试需求**：编译通过，CDC 场景验证通过
- **状态**: 所有任务已完成
