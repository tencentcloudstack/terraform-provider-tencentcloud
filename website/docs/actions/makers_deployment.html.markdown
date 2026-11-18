---
subcategory: "Web & Agents Development Platform(Makers)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_deployment"
sidebar_current: "docs-tencentcloud-action-makers_deployment"
description: |-
  Provides an action to trigger a Makers build via the CreateDeployment API. This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources. Each invocation triggers a brand-new build; destroying the action does not roll back or cancel a build already accepted by the cloud.
---

# tencentcloud_makers_deployment

Provides an action to trigger a Makers build via the CreateDeployment API. This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources. Each invocation triggers a brand-new build; destroying the action does not roll back or cancel a build already accepted by the cloud.

## Example Usage

### ~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

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

### Trigger it explicitly:

```hcl
shell
# Dry-run first: the plan shows the action without executing it
terraform plan - invoke = action.tencentcloud_makers_deployment.example

# Execute only this action, ignoring the rest of the configuration
terraform apply - invoke = action.tencentcloud_makers_deployment.example
```

### :

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

## Argument Reference

The following arguments are supported:

* `env` - (Required, String) Environment slug, e.g. `Production`, `Preview` or a custom environment.
* `project_id` - (Required, String) Project ID.
* `provider` - (Required, String) Code source. Use `Upload` when building from an uploaded artifact.
* `via_meta` - (Required, String) Build trigger source. Valid values: `Github`, `ReDeploy`, `Api`, `Webhook`, `Upload`, `Coding`, `Gitee`, `Gitlab`.
* `zone_id` - (Required, String) Makers site ID.
* `build_from` - (Optional, String) Build source.
* `ci_spec` - (Optional, String) Build specification. Reserved field: the server side currently still reads it from edgeone.json, so it has no effect yet.
* `dist_type` - (Optional, String) Artifact type. Valid values: `Zip`, `Folder`.
* `repo_branch` - (Optional, String) Repository branch to build.
* `repo_commit_hash` - (Optional, String) Commit hash to build.
* `repo_commit_msg` - (Optional, String) Commit message.
* `temp_bucket_path` - (Optional, String) Object path of the uploaded artifact in the temp bucket, used when `provider` is `Upload`.
* `template_id` - (Optional, String) Template ID.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [CreateDeployment](https://www.tencentcloud.com/search/?wd=CreateDeployment&type=document&lang=en)
  - [DescribeDeployments](https://www.tencentcloud.com/search/?wd=DescribeDeployments&type=document&lang=en)

