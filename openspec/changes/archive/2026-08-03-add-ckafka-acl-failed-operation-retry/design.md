## Context

`tencentcloud_ckafka_acl` 资源（文件 `tencentcloud/services/ckafka/resource_tc_ckafka_acl.go`）创建时调用服务层方法 `CkafkaService.CreateAcl`（文件 `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go`，第 331-359 行）。该方法内部通过 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹云 API `CreateAcl` 调用，错误统一交给 `tccommon.RetryError(err)` 判断。

`tccommon.RetryError` 仅对 `retryableErrorCode` 列表内的错误码（如 `ClientError.NetworkError`、`RequestLimitExceeded` 等）返回 `RetryableError`；`FailedOperation` 不在该列表中，因此会被判为 `NonRetryableError`，一旦云 API 返回 `Code=FailedOperation`，创建 ACL 会立即失败。

经核对 vendor 中 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819/client.go` 的 `CreateAclWithContext` 注释，`CreateAcl` 接口可能返回的错误码明确包含 `FAILEDOPERATION = "FailedOperation"`，因此针对该错误码增加重试是可行且有意义的。`extension_ckafka.go` 中已有常量 `CkafkaFailedOperation = "FailedOperation"` 可复用。

## Goals / Non-Goals

**Goals:**
- 在 `CreateAcl` 返回 `Code=FailedOperation` 时，最多重试 3 次，每次间隔固定 5 秒，提高 ACL 创建的最终成功率。
- 保持现有其他错误（网络错误、限流、参数错误等）的处理逻辑不变。
- 保持资源 Schema、ID 格式、`CreateAcl` 方法签名完全向后兼容。
- 通过单元测试（gomonkey mock 云 API）验证重试行为。

**Non-Goals:**
- 不修改 `Read` / `Delete` 等其它 CRUD 方法的重试逻辑。
- 不修改 `tccommon.RetryError` / `retryableErrorCode` 全局配置（避免影响其他服务）。
- 不引入第三方重试库；不修改 vendor 中的云 SDK。
- 不改变 `FailedOperation` 在其它 ckafka 方法中的既有处理方式。

## Decisions

### 决策 1：使用「固定次数 + 固定间隔」的自定义循环，而非 `resource.Retry` 的退避策略

需求明确要求「重试 3 次、每次间隔 5 秒」。`resource.Retry` 只能按「总超时 + 指数退避」重试，无法精确控制次数与固定间隔，因此**不能**仅靠把 `FailedOperation` 加入 `tccommon.RetryError` 的附加可重试错误码来实现（那样会变成超时驱动的退避重试）。

**方案**：在 `CreateAcl` 方法中，用 `for` 循环实现固定次数的重试，重试之间 `time.Sleep(5 * time.Second)`；每次尝试内部保留原有的 `resource.Retry` + `tccommon.RetryError`，从而兼容其他可重试错误。

**替代方���对比：**
- 方案 A（采用）：`for i := 0; i <= RETRY_TIMES; i++` + 每次尝试内嵌 `resource.Retry`。优点：精确满足次数/间隔需求，同时保留原有重试能力。
- 方案 B（不采用）：把 `FailedOperation` 作为 `tccommon.RetryError(e, "FailedOperation")` 的附加可重试错误码。缺点：重试次数与间隔由 `WriteRetryTimeout` 与指数退避决定，无法满足「3 次、5s」的明确约束。
- 方案 C（不采用）：完全移除 `resource.Retry`，改为纯 `for` 循环。缺点：会丢失对网络错误、限流等原有可重试错误的处理，属于行为退化。

### 决策 2：重试逻辑放在服务层 `CkafkaService.CreateAcl`，而非资源层

与项目既有模式一致（该资源的 `DeleteAcl`、`DescribeAclByFilter` 等重试都在服务层），服务层方法可被资源 Create/测试复用，避免在资源层重复实现。

### 决策 3：错误判定方式

使用 `errors.TencentCloudSDKError` 类型断言（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors`，service 文件已 import），判断 `sdkErr.Code == CkafkaFailedOperation`。仅对 `FailedOperation` 进入固定次数重试；断言失败（非 SDK 错误）或 Code 不匹配时直接跳出循环，交给上层返回。

### 决策 4：重试次数与间隔定义为命名常量

在 `extension_ckafka.go` 新增：
- `CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_TIMES = 3`
- `CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_INTERVAL = 5 * time.Second`

便于实现与单测共同引用，避免魔法数字。

### 决策 5：重试语义为「首次调用 + 最多 3 次重试」

循环边界 `i <= RETRY_TIMES`（即 `i` 从 0 到 3，共 4 次尝试：首次 + 3 次重试），在每次重试前 `Sleep(5s)`。若某次成功立即 break；若非 FailedOperation 错误立即 break；若重试次数耗尽仍为 FailedOperation，返回最后一次错误。

**实现示意：**

```go
var response *ckafka.CreateAclResponse
var err error

// 首次调用 + FailedOperation 时最多重试 3 次，每次间隔 5s
for i := 0; i <= CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_TIMES; i++ {
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().CreateAcl(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err == nil {
		break
	}
	sdkErr, ok := err.(*errors.TencentCloudSDKError)
	if !ok || sdkErr.Code != CkafkaFailedOperation {
		break
	}
	if i == CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_TIMES {
		break
	}
	log.Printf("[CRITAL]%s api[%s] fail with Code=FailedOperation, retry %d/%d after %s, reason[%s]",
		logId, request.GetAction(), i+1, CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_TIMES,
		CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_INTERVAL, err.Error())
	time.Sleep(CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_INTERVAL)
}

if err != nil {
	return err
}
if response != nil && response.Response != nil && !me.OperateStatusCheck(ctx, response.Response.Result) {
	return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
}
return nil
```

说明：循环内 `resource.Retry` 的闭包继续捕获外层 `response` / `err`；由于 `tccommon.RetryError` 会把 `FailedOperation` 判为不可重试，`resource.Retry` 会立即返回该错误，交由外层循环按固定次数/间隔重试，行为与需求一致。

## Risks / Trade-offs

- [重试期间阻塞] → 每次重试前 `time.Sleep(5s)` 会阻塞资源 Create 调用，最坏情况（4 次全部失败）总耗时约 15 秒。→ 这是需求明确要求的固定间隔；`terraform apply` 本身是同步操作，可接受。
- [重复创建风险] 云 API 可能在返回 `FailedOperation` 前已实际创建成功，重试可能产生重复 ACL。→ 属于云 API 返回语义问题；需求方明确要求在 `FailedOperation` 时重试，且重试次数有限（3 次），风险可控。
- [测试耗时] 单元测试中「全部失败」路径会真实 sleep 3 次（约 15 秒）。→ 测试用例数量少、耗时有限；测试覆盖重试次数/间隔的语义，收益大于成本。
- [非 FailedOperation 错误行为不变] 循环内对非 FailedOperation 错误直接 break 返回，与原 `resource.Retry` 行为一致，无回归风险。

## Migration Plan

1. 修改 `extension_ckafka.go` 新增重试常量。
2. 修改 `service_tencentcloud_ckafka.go` 的 `CreateAcl` 方法，引入固定次数重试循环。
3. 在 `resource_tc_ckafka_acl_test.go` 中补充 gomonkey 单元测试。
4. 更新 `resource_tc_ckafka_acl.md` 说明重试行为。
5. 回滚策略：本变更为纯逻辑增强，若出现异常可回退 `CreateAcl` 方法到原实现即可，无数据迁移需求。

## Open Questions

- 需求中「重试 3 次」按「首次调用 + 3 次重试」理解（共 4 次尝试）。若业务侧期望「总共 3 次尝试」，只需将循环边界调整为 `i < RETRY_TIMES` 即可，常量与测试同步调整。
