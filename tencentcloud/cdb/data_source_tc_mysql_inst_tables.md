Use this data source to query detailed information of mysql inst_tables

Example Usage

```hcl
data "tencentcloud_mysql_inst_tables" "inst_tables" {
  instance_id = "cdb-fitq5t9h"
  database = "tf_ci_test"
  # table_regexp = ""
}
```