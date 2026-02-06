Provides a list of Lighthouse blueprints (images).

Use this data source to query available blueprints for Lighthouse instances.

Example Usage

Query all blueprints:

```hcl
data "tencentcloud_lighthouse_blueprints" "all" {
}

output "blueprints" {
  value = data.tencentcloud_lighthouse_blueprints.all.blueprint_set
}
```

Filter by platform type:

```hcl
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
```

Filter by blueprint type:

```hcl
data "tencentcloud_lighthouse_blueprints" "app_os" {
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
```

Query specific blueprints by ID:

```hcl
data "tencentcloud_lighthouse_blueprints" "specific" {
  blueprint_ids = ["lhbp-xxx", "lhbp-yyy"]
}
```
