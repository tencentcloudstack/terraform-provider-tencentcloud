Use this data source to query APM (Application Performance Management) instances.

Example Usage

Query all APM instances

```hcl
data "tencentcloud_apm_instances" "all" {
}

output "instances" {
  value = data.tencentcloud_apm_instances.all.instance_list
}
```

Query APM instances by IDs

```hcl
data "tencentcloud_apm_instances" "by_ids" {
  instance_ids = ["apm-xxxxxxxx", "apm-yyyyyyyy"]
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_ids.instance_list
}
```

Query APM instances by name

```hcl
data "tencentcloud_apm_instances" "by_name" {
  instance_name = "test"
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_name.instance_list
}
```

Query APM instances by tags

```hcl
data "tencentcloud_apm_instances" "by_tags" {
  tags = {
    "Environment" = "Production"
    "Team"        = "DevOps"
  }
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_tags.instance_list
}
```
