Use this data source to query detailed information of oceanus resource_related_job

Example Usage

```hcl
data "tencentcloud_oceanus_resource_related_job" "example" {
  resource_id                    = "resource-8y9lzcuz"
  desc_by_job_config_create_time = 0
  resource_config_version        = 1
  work_space_id                  = "space-2idq8wbr"
}
```