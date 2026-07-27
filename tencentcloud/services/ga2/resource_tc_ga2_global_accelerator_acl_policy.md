Provides a resource to create a GA2 (Global Accelerator V2) global accelerator ACL policy

Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_ga2_global_accelerator_acl_policy" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  default_action        = "ACCEPT"
  status                = "OPEN"
}
```

Import

GA2 global accelerator ACL policy can be imported using the composite id `<global_accelerator_id>#<global_accelerator_acl_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_global_accelerator_acl_policy.example ga-jnyfyyss#gapolicy-xxx
```
