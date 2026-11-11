## 1. Preparatory Analysis

- [x] 1.1 Confirm `CreateDBInstanceRequest` and `CreateDBInstanceHourRequest` in vendored SDK (`vendor/.../cdb/v20170320/models.go`) expose an optional `ClientToken *string` field
- [x] 1.2 Review master instance create helpers `mysqlCreateInstancePayByMonth` / `mysqlCreateInstancePayByUse` in `resource_tc_mysql_instance.go` to capture the exact client-token retry pattern to replicate

## 2. Core Implementation

- [x] 2.1 In `mysqlCreateReadonlyInstancePayByMonth` (`resource_tc_mysql_readonly_instance.go`): generate `clientToken := helper.BuildToken()`, set `request.ClientToken = &clientToken`, wrap the `CreateDBInstance` call in `resource.Retry(2*tccommon.WriteRetryTimeout, ...)`, and on `len(r.Response.InstanceIds) < 1 && clientToken != ""` return `resource.RetryableError(...)`; wrap API errors with `tccommon.RetryError(inErr)`; assign the successful response to the outer variable and keep the post-retry `d.SetId(...)` logic
- [x] 2.2 In `mysqlCreateReadonlyInstancePayByUse` (`resource_tc_mysql_readonly_instance.go`): apply the same client-token polling mechanism to the `CreateDBInstanceHour` call, mirroring `mysqlCreateInstancePayByUse`
- [x] 2.3 Verify all returned errors are checked and that the success path still validates `len(response.Response.InstanceIds) == 1` before `d.SetId`

## 3. Unit Tests

- [x] 3.1 In `resource_tc_mysql_readonly_instance_test.go`, add a mock-based (gomonkey) unit test verifying the prepaid create path sets a non-empty `ClientToken` on the `CreateDBInstanceRequest`
- [x] 3.2 Add a mock-based unit test verifying the postpaid create path sets a non-empty `ClientToken` on the `CreateDBInstanceHourRequest`
- [x] 3.3 Add a mock-based unit test simulating an empty-`InstanceIds` response followed by a successful response, verifying the create path retries and ultimately succeeds

## 4. Documentation

- [x] 4.1 Confirm no schema or user-facing argument changes; verify the existing `resource_tc_mysql_readonly_instance.md` example file requires no update (client token is an internal implementation detail, not a schema field)