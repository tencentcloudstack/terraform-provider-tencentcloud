## 1. Service Layer

- [x] 1.1 Append `DescribeTeoSecurityClientAttesterById(ctx, zoneId, clientAttesterId)` to `tencentcloud/services/teo/service_tencentcloud_teo.go` — paginates `DescribeSecurityClientAttester` (Limit=100) until the entry with `Id == clientAttesterId` is found; returns `*teo.ClientAttester` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_security_client_attester.go` with full schema following `tencentcloud_igtm_strategy` style:
  - Top-level fields: `zone_id` (Required, ForceNew)
  - `client_attesters` (Required, TypeList, MaxItems:1) with sub-fields 100% mapping SDK `ClientAttester` struct:
    - `name` (Required), `attester_source` (Required), `attester_duration` (Optional)
    - `t_c_r_c_e_option` (Optional, List MaxItems:1): `channel`, `region`
    - `t_c_captcha_option` (Optional, List MaxItems:1): `captcha_app_id`, `app_secret_key`
    - `t_c_e_o_captcha_option` (Optional, List MaxItems:1): `captcha_mode`
    - `id` (Computed), `type` (Computed)

- [x] 2.2 Implement Create:
  - Build `CreateSecurityClientAttesterRequest` with `ZoneId` and `ClientAttesters` (single-element, no `Id`)
  - Call `CreateSecurityClientAttesterWithContext`; extract `ClientAttesterIds[0]`
  - Set resource ID to `strings.Join([]string{zoneId, clientAttesterId}, tccommon.FILED_SP)`
  - Call Read

- [x] 2.3 Implement Read:
  - Split ID → `zoneId`, `clientAttesterId`
  - Call `DescribeTeoSecurityClientAttesterById`; if nil → `d.SetId("")`
  - Populate `zone_id` and `client_attesters` block (all fields including nested option blocks) from response

- [x] 2.4 Implement Update:
  - Build `ModifySecurityClientAttesterRequest` with `ZoneId` and `ClientAttesters=[{Id, Name, AttesterSource, ...}]`
  - Call `ModifySecurityClientAttesterWithContext`
  - Call Read

- [x] 2.5 Implement Delete:
  - Build `DeleteSecurityClientAttesterRequest` with `ZoneId` and `ClientAttesterIds=[clientAttesterId]`
  - Call `DeleteSecurityClientAttesterWithContext` with Retry

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_security_client_attester` in `tencentcloud/provider.go` ResourcesMap, pointing to `teo.ResourceTencentCloudTeoSecurityClientAttester()`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_security_client_attester.md` — document all arguments, attributes, and import syntax with TC-RCE example HCL
- [x] 4.2 Create `tencentcloud/services/teo/resource_tc_teo_security_client_attester_test.go` — basic acceptance test covering create/update/import/delete following `resource_tc_igtm_strategy_test.go` style
