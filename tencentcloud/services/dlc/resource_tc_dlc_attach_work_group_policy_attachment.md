Provides a resource to create a DLC attach work group policy attachment

~> **NOTE:** `policy_id` format: `v1|{SubjectType}|{SubjectId}|{PolicyType}|{Mode}|{Catalog}|{Database}|{Table}|{View}|{Function}|{Column}|{DataEngine}|{Operation}`

Example Usage

If policy_type is ENGINE

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_attachment" "example" {
  work_group_id = 21420
  policy_set {
    policy_type = "ENGINE"
    catalog     = ""
    database    = ""
    table       = ""
    data_engine = "test"
    operation   = "USE,MONITOR,MODIFY"
    source      = "WORKGROUP"
  }
}
```

If policy_type is DATABASE

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_attachment" "example1" {
  work_group_id = 21420
  policy_set {
    policy_type = "DATABASE"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    table       = ""
    operation   = "OWNER"
    source      = "WORKGROUP"
    mode        = "COMMON"
  }
}
```

If policy_type is ROWFILTER

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_attachment" "example2" {
  work_group_id = 21420
  policy_set {
    policy_type = "ROWFILTER"
    catalog     = "DataLakeCatalog"
    database    = "test_database"
    table       = "test_table"
    operation   = "year > 2026 and country == 'US'"
    source      = "WORKGROUP"
    mode        = "SENIOR"
  }
}

```

Import

DLC attach work group policy attachment can be imported using the composite id, e.g. The composite id is `WorkGroupId#PolicyId`.

```
terraform import tencentcloud_dlc_attach_work_group_policy_attachment.example 21420#v1|WORKGROUP|21420|DATABASE|COMMON|DataLakeCatalog|test_database||||||OWNER
```
