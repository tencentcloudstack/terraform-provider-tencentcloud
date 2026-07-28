Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) ACL rule.

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

resource "tencentcloud_ga2_global_accelerator_acl_rule" "example" {
  global_accelerator_id            = tencentcloud_ga2_global_accelerator.example.id
  global_accelerator_acl_policy_id = tencentcloud_ga2_global_accelerator_acl_policy.example.global_accelerator_acl_policy_id
  protocol                         = "TCP"
  port                             = "80"
  source_cidr_block                = "10.0.0.0/24"
  policy                           = "ACCEPT"
  description                      = "tf example acl rule"
}
```

Import

GA2 ACL rule can be imported using the composite id `<global_accelerator_id>#<global_accelerator_acl_policy_id>#<global_accelerator_acl_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_global_accelerator_acl_rule.example ga-5hlomx9m#sp-8oc0ohct#sr-asa55vy5
```