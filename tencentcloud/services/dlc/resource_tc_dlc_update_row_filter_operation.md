Provides a resource to create a DLC update row filter operation

Example Usage

```hcl
resource "tencentcloud_dlc_update_row_filter_operation" "example" {
  policy_id = 103704
  policy {
    database    = "tf_example_db"
    catalog     = "DataLakeCatalog"
    table       = "test_table"
    operation   = "value!=\"0\""
    policy_type = "ROWFILTER"
    source      = "USER"
    mode        = "SENIOR"
    re_auth     = false
  }
}
```