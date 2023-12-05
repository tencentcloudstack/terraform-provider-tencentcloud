Provides a resource to create a antiddos scheduling_domain_user_name

Example Usage

```hcl
resource "tencentcloud_antiddos_scheduling_domain_user_name" "scheduling_domain_user_name" {
  domain_name = "test.com"
  domain_user_name = ""
}
```

Import

antiddos scheduling_domain_user_name can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_scheduling_domain_user_name.scheduling_domain_user_name ${domainName}
```