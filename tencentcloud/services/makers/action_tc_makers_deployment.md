Provides an action to trigger a Makers build via the CreateDeployment API. This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources. Each invocation triggers a brand-new build; destroying the action does not roll back or cancel a build already accepted by the cloud.

Example Usage

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

```hcl
action "tencentcloud_makers_deployment" "example" {
  config {
    zone_id     = "zone-123456"
    project_id  = "project-123456"
    via_meta    = "Api"
    env         = "Production"
    provider    = "Github"
    repo_branch = "main"
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
  # ...
}

action "tencentcloud_makers_deployment" "example" {
  config {
    zone_id     = "zone-123456"
    project_id  = "project-123456"
    via_meta    = "Api"
    env         = "Production"
    provider    = "Github"
    repo_branch = "main"
  }
}

resource "tencentcloud_makers_project" "with_trigger" {
  # ...

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.tencentcloud_makers_deployment.example]
    }
  }
}
```

~> **NOTE:** Repeating the invocation repeats the cloud-side side effect: each run submits a new build, and a build already accepted by the cloud is not rolled back when this action or its enclosing configuration is destroyed. `on_failure` (available in Terraform CLI 1.16+) only controls whether the apply continues (`halt`/`continue`) or the resource is tainted (`taint`, only for `after_create`); it is not a destroy guard, and during a full destroy walk action errors are downgraded to warnings.
