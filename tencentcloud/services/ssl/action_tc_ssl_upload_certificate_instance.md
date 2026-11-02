Provides an action to upload new certificate content and update the associated cloud resources (such as clb) without changing the certificate ID via the UploadUpdateCertificateInstance API. The new certificate is given either by an existing `certificate_id` or by a pair of `certificate_public_key`/`certificate_private_key`. The action waits for the asynchronous deploy task to complete. This is a one-time operation; no cloud-side state is persisted.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** Currently, this product only supports updating cloud resources of the `clb` type.

~> **NOTE:** If a deploy task for the same old certificate is already in progress, the API does not create a new task; the action waits for the existing task to complete.

~> **NOTE:** The `invoke` timeout (default `15m`) only limits how long the action polls the deploy record status; it does not bound the task-creation request itself, which uses the provider's built-in write retry timeout.

~> **NOTE:** `certificate_public_key` and `certificate_private_key` are shown in plain text in the Terraform plan output, because action schemas do not support write-only or sensitive attributes yet. Prefer `certificate_id` when possible, and take care to protect the configuration and plan files when key content is used.

Example Usage

```hcl
action "tencentcloud_ssl_upload_certificate_instance" "example" {
  config {
    old_certificate_id      = "f8xnynbd"
    resource_types          = ["clb"]
    certificate_public_key  = "-----BEGIN CERTIFICATE-----\nxxx\n-----END CERTIFICATE-----"
    certificate_private_key = "-----BEGIN RSA PRIVATE KEY-----\nxxx\n-----END RSA PRIVATE KEY-----"

    resource_types_regions {
      resource_type = "clb"
      regions       = ["ap-guangzhou"]
    }

    timeouts {
      invoke = "30m"
    }
  }
}
```
