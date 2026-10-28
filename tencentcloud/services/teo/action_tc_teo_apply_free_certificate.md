Provides an action to apply for a TEO (EdgeOne) free certificate for a domain via the ApplyFreeCertificate API. After the apply call succeeds, the action polls CheckFreeCertificateVerification until the certificate application succeeds, fails, or times out. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** After invoking this action, you need to configure the verification content (DNS record or verification file) according to the verification information returned by the ApplyFreeCertificate API (available in the Terraform apply logs). Once the verification is configured, the action will poll CheckFreeCertificateVerification to confirm the certificate application result within the timeout period. If the verification is not configured in time, the polling may time out and you can retry the apply after completing the configuration.

Example Usage

```hcl
action "tencentcloud_teo_apply_free_certificate" "example" {
  config {
    zone_id             = "zone-3fkff38fyw8s"
    domain              = "example.com"
    verification_method = "dns_challenge"
  }
}
```