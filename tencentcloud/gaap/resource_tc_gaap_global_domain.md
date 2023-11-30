Provides a resource to create a gaap global domain

Example Usage

```hcl
resource "tencentcloud_gaap_global_domain" "global_domain" {
  project_id = 0
  default_value = "xxxxxx.com"
  alias = "demo"
  tags={
		key = "value"
  }
}
```

Import

gaap global_domain can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain.global_domain ${projectId}#${domainId}
```