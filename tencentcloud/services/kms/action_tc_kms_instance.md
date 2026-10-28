Provides an action to purchase and activate a prepaid KMS (Key Management Service) instance via the CreateKmsInstance API. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_kms_instance" "example" {
  config {
    instance_charge_type = "PREPAID"
    product_version      = "professional"
    client_token         = "tf-token-xxx"

    instance_charge_prepaid {
      period     = 1
      renew_flag = 1
    }

    tags {
      tag_key   = "env"
      tag_value = "prod"
    }
  }
}
```