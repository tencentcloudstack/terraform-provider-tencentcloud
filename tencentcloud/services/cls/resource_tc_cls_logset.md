Provides a resource to create a CLS logset

Example Usage

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf-example"
  tags = {
    createdBy = "Terraform"
  }
}
```
Import

CLS logset can be imported using the id, e.g.
```
$ terraform import tencentcloud_cls_logset.example 698902ff-8b5a-4c65-824b-d8956f366351
```