Provides a resource to create a dayu DDoS policy attachment.

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_attachment" "dayu_ddos_policy_attachment_basic" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  resource_id   = "bgpip-00000294"
  policy_id     = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
```