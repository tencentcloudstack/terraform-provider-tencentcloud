Provides an action to trigger a Makers build via the CreateDeployment API. This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources. Each invocation triggers a brand-new build; destroying the action does not roll back or cancel a build already accepted by the cloud.

Example Usage

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.16 and later.

~> **NOTE:** Repeating the invocation repeats the cloud-side side effect: each run submits a new build, and a build already accepted by the cloud is not rolled back when this action or its enclosing configuration is destroyed.


```hcl
action "tencentcloud_makers_deployment" "example" {
  config {
    zone_id       = "zone-123456"
    project_id    = "project-123456"
    via_meta      = "Upload"
    env           = "Production"
    code_provider = "Upload"
    local_path    = "./projcet_path"
  }
}
```

Trigger it explicitly:

```shell
# Dry-run first: the plan shows the action without executing it
terraform plan -invoke=action.tencentcloud_makers_deployment.example

# Execute only this action, ignoring the rest of the configuration
terraform apply -invoke=action.tencentcloud_makers_deployment.example
```

Or trigger it automatically during apply through `lifecycle.action_trigger`:

```hcl

resource "tencentcloud_makers_project" "example" {
  zone_id         = "zone-3hzk0ndvnnyp"
  name            = "tf-test-upload"
  channel         = "Custom"
  area            = "global"
  code_provider   = "Upload"
  node_js_version = "22.11.0"
  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.tencentcloud_makers_deployment.zip, action.tencentcloud_makers_deployment.folder]
    }
    action_trigger {
      events  = [after_update]
      actions = [action.tencentcloud_makers_deployment.folder]
    }
  }
}

action "tencentcloud_makers_deployment" "zip" {
  config {
    zone_id       = caller.zone_id
    project_id    = caller.project_id
    via_meta      = "Upload"
    env           = "Production"
    code_provider = "Upload"
    dist_type     = "Zip"
    local_path    = "./makers.zip"
  }
}

action "tencentcloud_makers_deployment" "folder" {
  config {
    zone_id       = caller.zone_id
    project_id    = caller.project_id
    via_meta      = "Upload"
    env           = "Production"
    code_provider = "Upload"
    dist_type     = "Folder"
    local_path    = "./makers"
  }
}
```