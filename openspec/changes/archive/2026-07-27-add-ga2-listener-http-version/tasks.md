## 1. Schema modification

- [x] 1.1 In `tencentcloud/services/ga2/resource_tc_ga2_listener.go`, change the `http_version` schema field from `Computed: true` only to `Optional: true, Computed: true`. Update the `Description` to mention valid values (`HTTP/1.1`, `HTTP/2`), that it only applies to HTTPS listeners, and that it cannot be modified after creation.
- [x] 1.2 Move the `http_version` schema entry from the `// Computed` section (after `client_ca_certificates`) to the main Optional+Computed section (after `cipher_policy_id` and before `server_certificates`, maintaining declaration order consistent with the `CreateListenerRequest` struct).

## 2. Create function update

- [x] 2.1 In `resourceTencentCloudGa2ListenerCreate`, add a block to forward `http_version` to `CreateListenerRequest.HttpVersion` when the user has set it: `if v, ok := d.GetOk("http_version"); ok { request.HttpVersion = helper.String(v.(string)) }`. Place this block after the existing HTTPS-only fields section (after `client_ca_certificates`), since `HttpVersion` is applicable only to HTTPS listeners per the SDK documentation.

## 3. Update function immutable check

- [x] 3.1 In `resourceTencentCloudGa2ListenerUpdate`, add an `immutableArgs` slice containing `"http_version"` right after parsing the resource ID. Iterate over the slice and call `d.HasChange(field)` for each; if any immutable field has changed, return `fmt.Errorf("field '%s' cannot be modified after creation; it requires a new resource to be created", field)`.

## 4. Unit test update

- [x] 4.1 In `tencentcloud/services/ga2/resource_tc_ga2_listener_test.go`, add or update unit test cases to verify that the `http_version` field is correctly forwarded in the Create request when set, and that it defaults to computed when not set. Use mock (gomonkey) approach consistent with the existing test file pattern.

## 5. Documentation update

- [x] 5.1 In `tencentcloud/services/ga2/resource_tc_ga2_listener.md`, update the HCL example to include the `http_version` parameter in the HTTPS listener example (e.g., `http_version = "HTTP/2"`).

## 6. Build & verification

- [x] 6.1 Run `go build ./tencentcloud/...` and confirm clean compilation (note: this step is for verification only, not to be executed as part of code generation).
- [x] 6.2 Run `go vet ./tencentcloud/services/ga2/...` and confirm no new findings (note: this step is for verification only, not to be executed as part of code generation).
