Use this data source to query detailed information of dayu DDoS policy attachments

Example Usage

```hcl
data "tencentcloud_dayu_ddos_policy_attachments" "foo_type" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
data "tencentcloud_dayu_ddos_policy_attachments" "foo_resource" {
  resource_id   = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_id
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
data "tencentcloud_dayu_ddos_policy_attachments" "foo_policy" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.policy_id
}
```