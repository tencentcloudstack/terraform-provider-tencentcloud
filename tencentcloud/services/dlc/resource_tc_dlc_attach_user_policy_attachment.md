Provides a resource to create a DLC attach user policy attachment

~> **NOTE:** `policy_id` format: `v1|{SubjectType}|{SubjectId}|{PolicyType}|{Mode}|{Catalog}|{Database}|{Table}|{View}|{Function}|{Column}|{DataEngine}|{Operation}`

Example Usage

If policy_type is ENGINE

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example" {
  user_id      = "100010109702"
  account_type = "TencentAccount"
  policy_set {
    policy_type = "ENGINE"
    data_engine = "test_engine"
    operation   = "USE,MONITOR"
    source      = "USER"
  }
}
```

If policy_type is DATABASE

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example1" {
  user_id      = "100010109702"
  account_type = "TencentAccount"
  policy_set {
    policy_type = "DATABASE"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    mode        = "COMMON"
    operation   = "ASSAYER"
    source      = "USER"
  }
}
```

If policy_type is ROWFILTER

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example2" {
  user_id      = "100010109702"
  account_type = "TencentAccount"
  policy_set {
    policy_type = "ROWFILTER"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    table       = "test_table"
    mode        = "SENIOR"
    operation   = "year > 2026 and country == 'US'"
    source      = "USER"
  }
}
```

Import

DLC attach user policy attachment can be imported using the composite id (`user_id#policy_id`), e.g.

```
terraform import tencentcloud_dlc_attach_user_policy_attachment.example 100010109702#v1|USER|100010109702|DATABASE|COMMON|DataLakeCatalog|test_database||||||ASSAYER
```
