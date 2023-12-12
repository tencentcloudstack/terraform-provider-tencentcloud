Use this data source to query detailed information of cam policy_granting_service_access

Example Usage

```hcl
data "tencentcloud_cam_policy_granting_service_access" "policy_granting_service_access" {
  role_id = 4611686018436805021
  service_type = "cam"
  }
```