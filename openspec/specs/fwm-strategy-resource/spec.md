# fwm-strategy-resource Specification

## Purpose
TBD - created by archiving change add-fwm-strategy-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource MUST be registered as `tencentcloud_fwm_strategy`

The provider SHALL register a new general-type resource named `tencentcloud_fwm_strategy` whose CRUD callbacks invoke the `CreateStrategy` / `DescribeStrategy` / `DeleteStrategy` APIs of `tencentcloud-sdk-go/tencentcloud/fwm/v20250611`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `tencentcloud/provider.go` exposes the resource via key `"tencentcloud_fwm_strategy"` mapped to `fwm.ResourceTencentCloudFwmStrategy()`, in a FWM resource section.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the FireWall(FWM) Resource section MUST include `tencentcloud_fwm_strategy` so that `website/docs/r/fwm_strategy.html.markdown` is generated.

#### Scenario: FWM client accessor exists

- **WHEN** a CRUD callback needs the fwm client
- **THEN** `tencentcloud/connectivity/client.go` SHALL provide `UseFwmV20250611Client()` returning `*fwmv20250611.Client`, lazily initialized with `NewClientProfile(300)`.

### Requirement: Schema MUST mirror CreateStrategy API inputs plus Computed StrategyIds output

The resource schema SHALL declare top-level argument keys mapping 1:1 (snake_case) to `CreateStrategyRequest` fields, plus one Computed output field `strategy_ids` mapping to `response.StrategyIds`.

| HCL key | SDK field | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `product` | `Product` | TypeString | Yes | Yes | No |
| `receive_account` | `ReceiveAccount` | TypeList of TypeString | No | No | No |
| `pre_strategy` | `PreStrategy` | TypeList of (group_id, sequence) | No | No | No |
| `post_strategy` | `PostStrategy` | TypeList of (group_id, sequence) | No | No | No |
| `receive_group` | `ReceiveGroup` | TypeList of TypeString | No | No | No |
| `strategy_ids` | `response.StrategyIds` | TypeList of TypeString | No | No | **Yes** |

`strategy_ids` is the sole new parameter introduced by this change; it captures the `CreateStrategy` response output `response.StrategyIds`.

#### Scenario: strategy_ids populated after create

- **WHEN** `CreateStrategy` returns `response.StrategyIds` containing `["s1","s2"]`
- **THEN** the resource state SHALL have `strategy_ids = ["s1","s2"]` and the resource id SHALL equal `"s1#s2"` (joined by `tccommon.FILED_SP`).

#### Scenario: Empty StrategyIds response is rejected

- **WHEN** `CreateStrategy` returns `response.StrategyIds` that is nil or empty
- **THEN** Create SHALL return a non-retryable error and MUST NOT write an empty id into state.

#### Scenario: pre_strategy / post_strategy element mapping

- **WHEN** the user sets `pre_strategy { group_id = "g1" sequence = 1 }`
- **THEN** the corresponding `StrategyReq{ GroupId: "g1", Sequence: 1 }` SHALL be sent in `CreateStrategyRequest.PreStrategy`.

### Requirement: Resource id MUST be the joined StrategyIds list

The resource id SHALL be `strings.Join(strategyIds, tccommon.FILED_SP)`. Read / Update / Delete SHALL recover the list via `strings.Split(d.Id(), tccommon.FILED_SP)`.

#### Scenario: Delete uses recovered StrategyIds

- **WHEN** the resource is destroyed and its id is `"s1#s2"`
- **THEN** `DeleteStrategyRequest.StrategyIds` SHALL be set to `["s1","s2"]`.

#### Scenario: Import requires joined id

- **WHEN** the user imports via `terraform import tencentcloud_fwm_strategy.x s1#s2`
- **THEN** the resource SHALL be imported and Read SHALL query each id via `DescribeStrategy`.

### Requirement: Update MUST enforce immutability (CRD-only resource)

Because the fwm SDK exposes no `ModifyStrategy` API, the resource SHALL treat `CreateStrategy` inputs as immutable except via recreate. `product` is `ForceNew`; the other top-level inputs (`receive_account`, `pre_strategy`, `post_strategy`, `receive_group`) SHALL be checked in `immutableArgs` and any `d.HasChange` on them SHALL return an error. `strategy_ids` is Computed and excluded from immutability checks.

#### Scenario: Update rejects changed immutable field

- **WHEN** the user changes `pre_strategy` in-place and runs `terraform apply`
- **THEN** Update SHALL return an error instructing the user to recreate the resource instead.

### Requirement: CRUD calls MUST use retry and nil-guard

All `CreateStrategy` / `DescribeStrategy` / `DeleteStrategy` SDK calls SHALL be wrapped in `resource.Retry` with `tccommon.RetryError` and nil-pointer guards on `result`/`result.Response`. SetId / Set operations SHALL happen outside the retry block. Errors SHALL be reported with the resource name `fwm_strategy`.

#### Scenario: Read empty does not silently clear id

- **WHEN** `DescribeStrategy` returns empty for all ids
- **THEN** Read SHALL print `log.Printf("[CRUD] fwm_strategy id=%s", d.Id())` before `d.SetId("")`.

