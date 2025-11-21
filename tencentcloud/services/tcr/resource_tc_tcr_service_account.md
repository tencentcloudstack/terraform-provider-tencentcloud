Provides a resource to create a TCR service account.

Example Usage

Create custom account with specified duration days

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "standard"
  delete_bucket = true
  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf-example"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
}

resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf-example"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  duration    = 10
  disable     = false
  password    = "Password123"
  tags = {
    createdBy = "Terraform"
  }
}
```

With specified expiration time

```hcl
resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf-example"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  expires_at  = 1676897989000 //time stamp
  disable     = false
  tags = {
    createdBy = "Terraform"
  }
}
```

Import

TCR service account can be imported using the registryId#accountName, e.g.

```
terraform import tencentcloud_tcr_service_account.example tcr-ixgt2l0z#tf-example
```