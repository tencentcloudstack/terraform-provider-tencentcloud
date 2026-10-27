Provides an action to provision (deploy) a permission configuration to a target account of the Tencent Cloud Organization via the ProvisionRoleConfiguration API. The action waits for the asynchronous task returned by the API to finish. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_organization_provision_role_configuration" "example" {
  config {
    zone_id               = "z-1os7c9znogct"
    role_configuration_id = "rc-ihogrs0e6ceg"
    target_type           = "MemberUin"
    target_uin            = 100043****45

    timeouts {
      invoke = "3m"
    }
  }
}
```
