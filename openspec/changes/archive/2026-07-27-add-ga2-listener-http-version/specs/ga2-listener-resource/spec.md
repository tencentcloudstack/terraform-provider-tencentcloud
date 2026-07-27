## MODIFIED Requirements

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
- `http_version` (string, optional+computed): HTTP version for HTTPS listeners. Valid values: `HTTP/1.1`, `HTTP/2`. Only applicable to HTTPS listeners. Cannot be modified after creation.

The schema SHALL additionally expose the following Modify-only / Read-only attributes:
- `client_affinity_time` (int, optional+computed): session-stickiness duration. **NOTE:** silently ignored on Create (the SDK `CreateListenerRequest` has no `ClientAffinityTime` field) and forwarded only on Update.

The resource SHALL additionally expose the following read-only attributes hydrated from `DescribeListeners` response:
- `listener_id` (string, computed) — also stored as the second segment of `d.Id()`.
- `create_time` (string, computed)
- `status` (string, computed)
- `endpoint_group_counts` (int, computed)

#### Scenario: All required SDK input fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `ga2v20250115.CreateListenerRequestParams` (GlobalAcceleratorId, Name, PortRanges, Description, ListenerType, Protocol, IdleTimeout, GetRealIpType, ClientAffinity, RequestTimeout, XForwardedForRealIp, CertificationType, CipherPolicyId, ServerCertificates, ClientCaCertificates, HttpVersion) appears in the schema with semantically equivalent typing.

#### Scenario: HttpVersion is Optional+Computed
- **WHEN** a developer inspects the `http_version` schema field
- **THEN** the field is declared as `Optional: true, Computed: true, Type: schema.TypeString` with a description mentioning valid values `HTTP/1.1` and `HTTP/2` and that it only applies to HTTPS listeners.

#### Scenario: HttpVersion defaults to computed value when unset
- **WHEN** a user creates a listener without setting `http_version`
- **THEN** the field is populated from the `DescribeListeners` response after Read, with no diff on subsequent plans.

#### Scenario: HttpVersion change triggers immutable error
- **WHEN** a user changes `http_version` in their Terraform configuration and applies the plan
- **THEN** the Update function returns an error: `field "http_version" cannot be modified after creation; it requires a new resource to be created`.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no fields beyond those listed above; no derived flags or synthetic toggles are introduced.

### Requirement: Async create with task polling
On Create, the resource SHALL invoke `CreateListener`, capture the returned `TaskId` and `ListenerId`, and poll `DescribeTaskResult` via `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` until `Status == "SUCCESS"` or the user-supplied `Timeouts.Create` (default **5 minutes**) elapses.

The Create function SHALL forward `http_version` to `CreateListenerRequest.HttpVersion` when the user has set it, following the standard `if v, ok := d.GetOk("http_version"); ok { request.HttpVersion = helper.String(v.(string)) }` pattern.

#### Scenario: Successful async create with HttpVersion
- **WHEN** `CreateListener` succeeds with `HttpVersion` set and the polled task transitions to `SUCCESS` within the timeout
- **THEN** the resource sets the composite ID, invokes Read, and the state contains the user-specified `http_version` value.

#### Scenario: HttpVersion omitted on Create
- **WHEN** the user does not set `http_version` in their configuration
- **THEN** the Create function does not set `request.HttpVersion`, the API applies its default, and Read populates the computed value from the response.

#### Scenario: Async create timeout
- **WHEN** the task does not reach `SUCCESS` before the configured `Timeouts.Create`
- **THEN** the resource returns an error containing the task ID and last observed status.

#### Scenario: Empty TaskId or ListenerId
- **WHEN** `CreateListener` returns a nil `TaskId` or nil `ListenerId`
- **THEN** the resource returns an explicit error rather than dereferencing the nil pointer.
