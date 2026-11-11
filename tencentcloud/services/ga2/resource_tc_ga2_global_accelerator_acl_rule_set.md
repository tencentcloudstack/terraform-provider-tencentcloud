Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) ACL rule set that manages the full collection of ACL rules under one ACL policy.

~> **NOTE:** This resource must exclusive in one acl policy, do not declare additional acl rule resources of this acl policy elsewhere.

~> **NOTE:** The field length for `acl_entries` is subject to account resource quotas.

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

resource "tencentcloud_ga2_global_accelerator_acl_rule_set" "example" {
  global_accelerator_id            = tencentcloud_ga2_global_accelerator.example.id
  global_accelerator_acl_policy_id = tencentcloud_ga2_global_accelerator_acl_policy.example.global_accelerator_acl_policy_id

  acl_entries {
    protocol          = "TCP"
    port              = "80"
    source_cidr_block = "10.0.0.0/24"
    policy            = "ACCEPT"
    description       = "tf example acl rule 1"
  }

  acl_entries {
    protocol          = "UDP"
    port              = "443"
    source_cidr_block = "10.0.1.0/24"
    policy            = "DROP"
    description       = "tf example acl rule 2"
  }
}
```

Import

GA2 ACL rule set can be imported using the composite id `<global_accelerator_id>#<global_accelerator_acl_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_global_accelerator_acl_rule_set.example ga-cmbzp36q#sp-jz94sb2t
```
