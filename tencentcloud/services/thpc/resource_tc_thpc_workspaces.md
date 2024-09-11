Provides a resource to create a thpc thpc_workspaces

Example Usage

```hcl
resource "tencentcloud_thpc_workspaces" "thpc_workspaces" {
  placement = {
  }
  space_charge_prepaid = {
  }
  system_disk = {
  }
  data_disks = {
  }
  virtual_private_cloud = {
  }
  internet_accessible = {
  }
  login_settings = {
  }
  enhanced_service = {
    security_service = {
    }
    monitor_service = {
    }
    automation_service = {
    }
  }
  tag_specification = {
    tags = {
    }
  }
}
```

Import

thpc thpc_workspaces can be imported using the id, e.g.

```
terraform import tencentcloud_thpc_workspaces.thpc_workspaces thpc_workspaces_id
```
