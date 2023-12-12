Provides a resource to create a tsf group

Example Usage

```hcl
resource "tencentcloud_tsf_group" "group" {
  application_id = "application-xxx"
  namespace_id = "namespace-aemrxxx"
  group_name = "terraform-test"
  cluster_id = "cluster-vwgjxxxx"
  group_desc = "terraform desc"
  alias = "terraform test"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_group.group group-axxx
```