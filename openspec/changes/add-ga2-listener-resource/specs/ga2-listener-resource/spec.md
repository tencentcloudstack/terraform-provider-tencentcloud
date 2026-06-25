## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_ga2_listener` that manages a single Tencent Cloud Global Accelerator V2 Listener per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `ga2` namespace, adjacent to the existing GA2 entries. The resource MUST also appear under the `Global Accelerator(GA2)` Resources section of `tencentcloud/provider.md` so that `make doc` generates a website page for it.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_ga2_listener" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation or vet errors related to the new resource.

#### Scenario: Website doc generated
- **WHEN** an operator runs `make doc`
- **THEN** `website/docs/r/ga2_listener.html.markdown` is generated from the resource Schema/Description and the example markdown file.

### Requirement: Schema mirrors `CreateListener`
The resource schema SHALL expose every input parameter accepted by the `CreateListener` API as a top-level attribute, with no renaming or merging:
- `global_accelerator_id` (string, required, **ForceNew**)
- `name` (string, optional+computed): listener name, ≤60 bytes.
- `port_ranges` (list of one nested block, required, **ForceNew**) — nested fields:
  - `from_port` (int, required): inclusive start port.
  - `to_port` (int, required): inclusive end port.
- `description` (string, optional+computed), ≤100 bytes.
- `listener_type` (string, optional+computed, **ForceNew**): smart-routing type.
- `protocol` (string, optional+computed, **ForceNew**): TCP / UDP / HTTP / HTTPS, default TCP.
- `idle_timeout` (int, optional+computed): connection idle timeout in seconds.
- `get_real_ip_type` (string, optional+computed): `TOA` or `ProxyProtocol`.
- `client_affinity` (string, optional+computed): session-stickiness toggle.
- `request_timeout` (int, optional+computed): request timeout in seconds.
- `x_forwarded_for_real_ip` (bool, optional+computed): L7 real-IP toggle.
- `certification_type` (string, optional+computed): `UNIDIRECTIONAL` / `MUTUAL`.
- `cipher_policy_id` (string, optional+computed): TLS cipher pack ID.
- `server_certificates` (set of string, optional+computed): server certificate IDs.
- `client_ca_certificates` (set of string, optional+computed): client CA certificate IDs.

The schema SHALL additionally expose the following Modify-only / Read-only attributes:
- `client_affinity_time` (int, optional+computed): session-stickiness duration. **NOTE:** silently ignored on Create (the SDK `CreateListenerRequest` has no `ClientAffinityTime` field) and forwarded only on Update.

The resource SHALL additionally expose the following read-only attributes hydrated from `DescribeListeners` response:
- `listener_id` (string, computed) — also stored as the second segment of `d.Id()`.
- `http_version` (string, computed)
- `create_time` (string, computed)
- `status` (string, computed)
- `endpoint_group_counts` (int, computed)

#### Scenario: All required SDK input fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `ga2v20250115.CreateListenerRequestParams` (GlobalAcceleratorId, Name, PortRanges, Description, ListenerType, Protocol, IdleTimeout, GetRealIpType, ClientAffinity, RequestTimeout, XForwardedForRealIp, CertificationType, CipherPolicyId, ServerCertificates, ClientCaCertificates) appears in the schema with semantically equivalent typing.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no fields beyond those listed above; no derived flags or synthetic toggles are introduced.

### Requirement: Resource ID
The resource ID SHALL be the composite `<GlobalAcceleratorId><FILED_SP><ListenerId>`, using the project-standard separator `tccommon.FILED_SP`. The resource SHALL support `terraform import` using the composite ID.

#### Scenario: Create sets the composite ID
- **WHEN** `CreateListener` succeeds and the polled task transitions to `SUCCESS`
- **THEN** the resource calls `d.SetId(strings.Join([]string{gaId, listenerId}, tccommon.FILED_SP))`.

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_ga2_listener.x ga-xxxxxxxx#lsr-yyyyyyy`
- **THEN** the resource state is hydrated from `DescribeListeners` using the parsed `(gaId, listenerId)` pair.

#### Scenario: Malformed import ID
- **WHEN** the import ID does not contain exactly one `tccommon.FILED_SP` separator, or has empty components
- **THEN** the resource returns a descriptive error before any SDK call.

### Requirement: Async create with task polling
On Create, the resource SHALL invoke `CreateListener`, capture the returned `TaskId` and `ListenerId`, and poll `DescribeTaskResult` via `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` until `Status == "SUCCESS"` or the user-supplied `Timeouts.Create` (default **5 minutes**) elapses.

#### Scenario: Successful async create
- **WHEN** `CreateListener` succeeds and the polled task transitions to `SUCCESS` within the timeout
- **THEN** the resource sets the composite ID, invokes Read, and returns no error.

#### Scenario: Async create timeout
- **WHEN** the task does not reach `SUCCESS` before the configured `Timeouts.Create`
- **THEN** the resource returns an error containing the task ID and last observed status.

#### Scenario: Empty TaskId or ListenerId
- **WHEN** `CreateListener` returns a nil `TaskId` or nil `ListenerId`
- **THEN** the resource returns an explicit error rather than dereferencing the nil pointer.

### Requirement: Read with retry and pagination
On Read, the resource SHALL call `Ga2Service.DescribeGa2ListenerById(ctx, gaId, listenerId)`, which:
- Wraps the SDK call `DescribeListenersWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Sets `request.GlobalAcceleratorId` once outside the loop.
- Filters by `Filters=[{Name:"listener-id", Values:[listenerId]}]`.
- Iterates pages with `Limit=100` (the documented maximum), constructing the request object once **outside** the loop and only mutating `Offset` / `Limit` per iteration.
- Strict-equals on both `*item.ListenerId == listenerId` and `*item.GlobalAcceleratorId == gaId` before returning.
- Returns `(nil, nil)` when no matching listener exists.

When the helper returns `(nil, nil)`, the resource SHALL clear `d.SetId("")` and log a `[WARN]` line indicating the listener may have been deleted out of band.

#### Scenario: Listener present
- **WHEN** the helper finds a matching `ListenerSet`
- **THEN** the resource populates all schema fields (input + computed) from the response.

#### Scenario: Listener removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching listener)
- **THEN** the resource calls `d.SetId("")` and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeListenersRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

#### Scenario: PortRanges flattened back into nested block
- **WHEN** the response carries non-nil `PortRanges.FromPort` and `PortRanges.ToPort`
- **THEN** the resource populates `port_ranges = [{ from_port: <FromPort>, to_port: <ToPort> }]`.

### Requirement: Update path
The Update function SHALL:
- Build a `ModifyListenerRequest` populated only with fields the user changed (or, when retaining the existing value is required by the API, with the current value of every Modify-accepted field).
- Always include `GlobalAcceleratorId` and `ListenerId` (mandatory by the API).
- Wrap the SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
- Wait for the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.
- Skip the Modify call entirely if no Modify-supported field changed.
- Refuse to issue a Modify call if any of the ForceNew fields (`port_ranges`, `listener_type`, `protocol`, `global_accelerator_id`) appears to have changed (the SDK plugin layer already triggers a destroy/create in this case; the resource Update function MUST NOT silently no-op those changes).

#### Scenario: Mutable field change
- **WHEN** only `name` or `description` changes
- **THEN** `ModifyListener` is called with `Name` or `Description` populated, awaited via the task helper.

#### Scenario: No-op update
- **WHEN** `terraform apply` runs but no Modify-supported field has changed
- **THEN** the resource skips `ModifyListener` and immediately invokes Read.

### Requirement: Async delete
On Delete, the resource SHALL call `DeleteListener`, capture the returned `TaskId`, and poll the task to completion via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))` (default **5 minutes**).

#### Scenario: Successful async delete
- **WHEN** the delete task transitions to `SUCCESS`
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

### Requirement: Retry coverage
Every SDK call (`CreateListenerWithContext`, `DescribeListenersWithContext`, `ModifyListenerWithContext`, `DeleteListenerWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations.

#### Scenario: Transient SDK error
- **WHEN** any of the four SDK calls returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Nil response defense
- **WHEN** an SDK call returns a nil `Response` (or a nil critical sub-field such as `TaskId` or `ListenerId`)
- **THEN** the wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_ga2_listener.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body (matching the existing GA2 log format).
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure.
- A `[WARN]` line when the resource is detected as deleted out of band during Read.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/ga2/resource_tc_ga2_listener.md` containing a self-contained `terraform { ... } resource "tencentcloud_ga2_listener" "..." { ... }` example and a `terraform import` example using the composite ID. Filename pattern follows `resource_tc_config_compliance_pack.md`.
- An acceptance-test file at `tencentcloud/services/ga2/resource_tc_ga2_listener_test.go` exercising at minimum: create, basic update (e.g. `name` / `description`), import, and destroy. Filename pattern follows `resource_tc_config_compliance_pack_test.go`.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains both an HCL example and an `import` example using the composite ID.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares `TestAccTencentCloudGa2ListenerResource_basic` (or equivalent) using `resource.Test` with a `CheckDestroy`-equivalent and at least two `Steps` (create + update + import).

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** the four Listener APIs plus `DescribeTaskResult` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` before any code is written.
