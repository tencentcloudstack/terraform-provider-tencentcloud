Provides a resource to create a tcr service account.

Example Usage

Create custom account with specified duration days

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr-instance"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_test_tcr_namespace"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "tf_example_cve_id"
  }
}

resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  duration    = 10
  disable     = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

With specified expiration time

```hcl
resource "tencentcloud_tcr_service_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  expires_at  = 1676897989000 //time stamp
  disable     = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr service_account can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_service_account.service_account registry_id#account_name
```