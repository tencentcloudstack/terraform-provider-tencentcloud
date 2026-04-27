# Add tencentcloud_teo_security_client_attester Resource

## What

Add a new Terraform resource `tencentcloud_teo_security_client_attester` for managing Tencent Cloud EdgeOne (TEO) client authentication options (Client Attesters). This resource manages a single client attester entry under a site zone, supporting full CRUD lifecycle.

## Why

TEO client attesters allow users to define authentication strategies (TC-RCE, TC-CAPTCHA, TC-EO-CAPTCHA) for client identity verification in API and web security scenarios. Currently no Terraform resource exists to manage these attesters, requiring manual portal operations. This resource enables infrastructure-as-code management of TEO client attesters.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateSecurityClientAttester` | `ClientAttesters` array length limited to **1**; returns `ClientAttesterIds` array |
| Read | `DescribeSecurityClientAttester` | Paginate by zone; match by `Id` field in response |
| Update | `ModifySecurityClientAttester` | Pass `Id` + updated fields in `ClientAttesters[0]` |
| Delete | `DeleteSecurityClientAttester` | Pass `ZoneId` + `ClientAttesterIds=[<id>]` |

## Resource ID

Composite: `<zone_id>#<client_attester_id>` (using `tccommon.FILED_SP`), e.g. `zone-123123322#attest-2184008405`.
`client_attester_id` is taken from `ClientAttesterIds[0]` returned by `CreateSecurityClientAttester`.
