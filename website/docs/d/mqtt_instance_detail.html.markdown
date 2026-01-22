---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_instance_detail"
sidebar_current: "docs-tencentcloud-datasource-mqtt_instance_detail"
description: |-
  Use this data source to query detailed information of MQTT instance detail
---

# tencentcloud_mqtt_instance_detail

Use this data source to query detailed information of MQTT instance detail

## Example Usage

```hcl
data "tencentcloud_mqtt_instance_detail" "example" {
  instance_id = "mqtt-kngmpg9p"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `authorization_policy_limit` - Limit on the number of authorization rules.
* `authorization_policy` - Authorization Policy Switch.
* `automatic_activation` - Is it automatically activated when registering device certificates automatically.
* `client_num_limit` - Number of client connections online.
* `created_time` - Creation time, millisecond timestamp.
* `destroy_time` - Pre destruction time, millisecond timestamp.
* `device_certificate_provision_type` - Client certificate registration method: JITP: Automatic Registration; API: Manually register through API.
* `expiry_time` - Expiration time, millisecond level timestamp.
* `instance_name` - Instance type. BASIC- Basic Edition; PRO- professional edition; PLATINUM- Platinum version.
* `instance_status` - Instance status. RUNNING- In operation; MAINTAINING- Under Maintenance; ABNORMAL- abnormal; OVERDUE- Arrears of fees; DESTROYED- Deleted; CREATING- Creating in progress; MODIFYING- In the process of transformation; CREATE_FAILURE- Creation failed; MODIFY_FAILURE- Transformation failed; DELETING- deleting.
* `instance_type` - Instance ID.
* `max_ca_num` - Maximum Ca quota.
* `max_subscription_per_client` - Maximum number of subscriptions per client.
* `max_subscription` - Maximum number of subscriptions in the cluster.
* `pay_mode` - Billing mode, POSTPAID, pay as you go PREPAID, annual and monthly package.
* `registration_code` - Certificate registration code.
* `remark` - Remark.
* `renew_flag` - Whether to renew automatically. Only the annual and monthly package cluster is effective. 1: Automatic renewal; 0: Non automatic renewal.
* `sku_code` - Product specifications.
* `topic_num_limit` - Maximum number of instance topics.
* `topic_num` - Topic num.
* `tps_limit` - Elastic TPS current limit value.
* `x509_mode` - TLS, Unidirectional authentication mTLS, bidirectional authentication BYOC; One machine, one certificate.


