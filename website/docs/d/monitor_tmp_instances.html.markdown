---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_instances"
sidebar_current: "docs-tencentcloud-datasource-monitor_tmp_instances"
description: |-
  Use this data source to query detailed information of monitor tmp instances
---

# tencentcloud_monitor_tmp_instances

Use this data source to query detailed information of monitor tmp instances

## Example Usage

```hcl
data "tencentcloud_monitor_tmp_instances" "tmp_instances" {
  instance_ids = ["prom-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_type` - (Optional, Int) Filter according to instance charge type.
	- 2: Prepaid;
	- 3: Postpaid by hour.
* `instance_ids` - (Optional, Set: [`String`]) Query according to one or more instance IDs. The instance ID is like: prom-xxxx. The maximum number of instances requested is 100.
* `instance_name` - (Optional, String) Filter according to instance name.
* `instance_status` - (Optional, Set: [`Int`]) Filter according to instance status.
	- 1: Creating;
	- 2: In operation;
	- 3: Abnormal;
	- 4: Reconstruction;
	- 5: Destruction;
	- 6: Stopped taking;
	- 8: Suspension of service due to arrears;
	- 9: Service has been suspended due to arrears.
* `ipv4_address` - (Optional, Set: [`String`]) Filter according to ipv4 address.
* `result_output_file` - (Optional, String) Used to save results.
* `tag_filters` - (Optional, List) Filter according to tag Key-Value pair. The tag-key is replaced with a specific label key.
* `zones` - (Optional, Set: [`String`]) Filter according to availability area. The availability area is shaped like: ap-Guangzhou-1.

The `tag_filters` object supports the following:

* `key` - (Required, String) The key of the tag.
* `value` - (Required, String) The value of the tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_set` - Instance details list.
  * `alert_rule_limit` - Alert rule limit.
  * `api_root_path` - Prometheus http api root address.
  * `auth_token` - Token required for data writing.
  * `auto_renew_flag` - Automatic renewal flag.
	- 0: No automatic renewal;
	- 1: Enable automatic renewal;
	- 2: Automatic renewal is prohibited;
	- -1: Invalid.
  * `charge_status` - Charge status.
	- 1: Normal;
	- 2: Expires;
	- 3: Destruction;
	- 4: Allocation;
	- 5: Allocation failed.
  * `created_at` - Created_at.
  * `data_retention_time` - Data retention time.
  * `enable_grafana` - Whether to enable grafana.
	- 0: closed;
	- 1: open.
  * `expire_time` - Expires for purchased instances.
  * `grafana_instance_id` - Binding grafana instance id.
  * `grafana_ip_white_list` - Grafana IP whitelist list.
  * `grafana_status` - Grafana status.
	- 1: Creating;
	- 2: In operation;
	- 3: Abnormal;
	- 4: Rebooting;
	- 5: Destruction;
	- 6: Shutdown;
	- 7: Deleted.
  * `grafana_url` - Grafana panel url.
  * `grant` - Authorization information for the instance.
    * `has_agent_manage` - Whether you have permission to manage the agent (1=yes, 2=no).
    * `has_api_operation` - Whether to display API and other information (1=yes, 2=no).
    * `has_charge_operation` - Whether you have charging operation authority (1=yes, 2=no).
    * `has_grafana_status_change` - Whether the status of Grafana can be modified (1=yes, 2=no).
    * `has_tke_manage` - Whether you have permission to manage TKE integration (1=yes, 2=no).
    * `has_vpc_display` - Whether to display VPC information (1=yes, 2=no).
  * `instance_charge_type` - Instance charge type.
	- 2: Prepaid;
	- 3: Postpaid by hour.
  * `instance_id` - Instance id.
  * `instance_name` - Instance name.
  * `instance_status` - Filter according to instance status.
	- 1: Creating;
	- 2: In operation;
	- 3: Abnormal;
	- 4: Reconstruction;
	- 5: Destruction;
	- 6: Stopped taking;
	- 8: Suspension of service due to arrears;
	- 9: Service has been suspended due to arrears.
  * `ipv4_address` - IPV4 address.
  * `is_near_expire` - Whether it is about to expire.
	- 0: No;
	- 1: Expiring soon.
  * `migration_type` - Migration status.
	- 0: Not in migration;
+	- 1: Migrating, original instance;
+	- 2: Migrating, target instance.
  * `proxy_address` - Proxy address.
  * `recording_rule_limit` - Pre-aggregation rule limitations.
  * `region_id` - Region id.
  * `remote_write` - Address of prometheus remote write.
  * `spec_name` - Specification name.
  * `subnet_id` - Subnet id.
  * `tag_specification` - List of tags associated with the instance.
    * `key` - The key of the tag.
    * `value` - The value of the tag.
  * `vpc_id` - VPC id.
  * `zone` - Zone.


