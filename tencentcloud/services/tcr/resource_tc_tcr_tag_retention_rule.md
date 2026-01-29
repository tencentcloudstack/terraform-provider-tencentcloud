Provides a resource to create a TCR tag retention rule.

Example Usage

Create and enable a tcr tag retention rule instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "standard"
  delete_bucket = true
  tags = {
    "createdBy" = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example"
  severity    = "medium"
}

resource "tencentcloud_tcr_tag_retention_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  advanced_rule_items {
    repository_filter {
      decoration = "repoMatches"
      pattern    = "**"
    }

    retention_policy {
      key   = "nDaysSinceLastPush"
      value = 2
    }
    
    tag_filter {
      decoration = "matches"
      pattern    = "**"
    }
  }

  cron_setting = "daily"
}
```

Import

TCR tag retention rule can be imported using the registryId#namespaceName#retentionId, e.g.

```
terraform import tencentcloud_tcr_tag_retention_rule.example tcr-s1jud21h#tf_example#3
```
