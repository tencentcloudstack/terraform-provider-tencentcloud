# Design: tencentcloud_teo_security_client_attester Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ tencentcloud/services/teo/resource_tc_teo_security_client_attester.go  (CRUD handlers)
           └─ tencentcloud/services/teo/service_tencentcloud_teo.go (DescribeTeoSecurityClientAttesterById)
                  └─ teo SDK v20220901
```

## Resource ID

Composite: `<zone_id>#<client_attester_id>` (using `tccommon.FILED_SP`), e.g. `zone-123123322#attest-2184008405`.

## Key Constraint

`ClientAttesters` array in Create/Modify is limited to **1 element**. The `client_attesters` field uses `TypeList` with `MaxItems: 1` to align with the SDK `ClientAttester` struct.

## Schema

### Top-level Required

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `zone_id` | String | Yes | Site ID |
| `client_attesters` | List (MaxItems:1) | No | Client attester configuration |

### client_attesters sub-fields (100% mapping to SDK ClientAttester struct)

| Sub-field | Type | Required | SDK Field | Notes |
|---|---|---|---|---|
| `name` | String | Required | `Name` | Attester name |
| `attester_source` | String | Required | `AttesterSource` | `TC-RCE`, `TC-CAPTCHA`, `TC-EO-CAPTCHA` |
| `attester_duration` | String | Optional | `AttesterDuration` | e.g. `300s`, `5m`, `1h` |
| `t_c_r_c_e_option` | List (MaxItems:1) | Optional | `TCRCEOption` | Required when `attester_source=TC-RCE` |
| `t_c_captcha_option` | List (MaxItems:1) | Optional | `TCCaptchaOption` | Required when `attester_source=TC-CAPTCHA` |
| `t_c_e_o_captcha_option` | List (MaxItems:1) | Optional | `TCEOCaptchaOption` | Required when `attester_source=TC-EO-CAPTCHA` |
| `id` | String | Computed | `Id` | Attester ID returned by API |
| `type` | String | Computed | `Type` | `PRESET` or `CUSTOM`, read-only from API |

#### t_c_r_c_e_option sub-fields (TCRCEOption)

| Sub-field | Type | SDK Field | Description |
|---|---|---|---|
| `channel` | String | `Channel` | TC-RCE channel ID |
| `region` | String | `Region` | Channel region |

#### t_c_captcha_option sub-fields (TCCaptchaOption)

| Sub-field | Type | SDK Field | Description |
|---|---|---|---|
| `captcha_app_id` | String | `CaptchaAppId` | CaptchaAppId |
| `app_secret_key` | String | `AppSecretKey` | AppSecretKey |

#### t_c_e_o_captcha_option sub-fields (TCEOCaptchaOption)

| Sub-field | Type | SDK Field | Description |
|---|---|---|---|
| `captcha_mode` | String | `CaptchaMode` | `Invisible` or `Adaptive` |

## Read Logic

Call `DescribeSecurityClientAttester` with `ZoneId`, paginate (Limit=100) until the entry with `ClientAttester.Id == clientAttesterId` is found.
If not found → resource deleted → `d.SetId("")`.

## Update Logic

Call `ModifySecurityClientAttester` with `ZoneId` and `ClientAttesters=[{Id: clientAttesterId, ...}]`.

## Delete Logic

Call `DeleteSecurityClientAttester` with `ZoneId` and `ClientAttesterIds=[clientAttesterId]`.

## Key SDK Types

```go
// teo v20220901
type ClientAttester struct {
    Id                *string
    Name              *string
    Type              *string          // Computed only
    AttesterSource    *string
    AttesterDuration  *string
    TCRCEOption       *TCRCEOption
    TCCaptchaOption   *TCCaptchaOption
    TCEOCaptchaOption *TCEOCaptchaOption
}

type TCRCEOption struct {
    Channel *string
    Region  *string
}

type TCCaptchaOption struct {
    CaptchaAppId  *string
    AppSecretKey  *string
}

type TCEOCaptchaOption struct {
    CaptchaMode *string
}
```
