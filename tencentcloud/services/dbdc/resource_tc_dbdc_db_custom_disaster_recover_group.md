Provides a resource to create a DBDC db custom disaster recover group.

Example Usage

```hcl
resource "tencentcloud_dbdc_db_custom_disaster_recover_group" "example" {
  name     = "tf-example1"
  type     = "HOST"
  strategy = "SPREAD"
  affinity = 1

  tags = {
    createBy = "Terraform"
  }
}
```

Import

DBDC db custom disaster recover group can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_disaster_recover_group.example dbps-sdcm1hfl
```
