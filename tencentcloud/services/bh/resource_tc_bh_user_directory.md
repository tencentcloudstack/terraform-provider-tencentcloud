Provides a resource to create a BH user directory

Example Usage

```hcl
resource "tencentcloud_bh_user_directory" "example" {
  dir_id   = 895784
  dir_name = "tf-example"
  user_org_set {
    org_id        = 1576799
    org_name      = "orgName1"
    org_id_path   = "819729.895784"
    org_name_path = "Root.demo1"
    user_total    = 0
  }

  user_org_set {
    org_id        = 896536
    org_name      = "orgName2"
    org_id_path   = "819729.895784.896536"
    org_name_path = "Root.demo2.demo3"
    user_total    = 1
  }
  source      = 0
  source_name = "sourceName"
}
```

Import

BH user directory can be imported using the id, e.g.

```
terraform import tencentcloud_bh_user_directory.example 32
```
