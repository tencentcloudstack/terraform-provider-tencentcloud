---
subcategory: "Cwp"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cwp_machines"
sidebar_current: "docs-tencentcloud-datasource-cwp_machines"
description: |-
  Use this data source to query detailed information of CWP machines
---

# tencentcloud_cwp_machines

Use this data source to query detailed information of CWP machines

## Example Usage

```hcl
data "tencentcloud_cwp_machines" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
}
```

### Query by Keyword filter

```hcl
data "tencentcloud_cwp_machines" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"

  filters {
    name        = "Keywords"
    values      = ["tf_example"]
    exact_match = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `machine_region` - (Required, String) Machine region. For example, ap-guangzhou and ap-shanghai.
* `machine_type` - (Required, String) Type of the machine's zone
CVM: Cloud Virtual Machine
BM: BMECM: Edge Computing Machine
LH: Lighthouse
Other: Hybrid Cloud Zone.
* `filters` - (Optional, List) Filter criteria
<li>Ips - String - required: no - query by IP</li>
<li>Names - String - required: no - query by instance name</li>
<li>InstanceIds - String - required: no - instance ID for query </li>
<li>Status - String - required: no - client online status (OFFLINE: offline/shut down | ONLINE: online | UNINSTALLED: not installed | AGENT_OFFLINE: agent offline | AGENT_SHUTDOWN: agent shut down)</li>
<li>Version - String required: no - current edition ( PRO_VERSION: Pro Edition | BASIC_VERSION: Basic Edition | Flagship: Ultimate Edition | ProtectedMachines: Pro + Ultimate Editions)</li>
<li>Risk - String - required: no - risky host (yes)</li>
<li>Os - String - required: no - operating system (value of DescribeMachineOsList)</li>
Each filter criterion supports only one value.
<li>Quuid - String - required: no - CVM instance UUID. Maximum value: 100.</li>
<li>AddedOnTheFifteen - String required: no - whether to query only hosts added within the last 15 days (1: yes) </li>
<li> TagId - String required: no - query the list of hosts associated with the specified tag </li>.
* `project_ids` - (Optional, Set: [`Int`]) ID List of Businesses to which machines belong.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Name of filter key.
* `values` - (Required, Set) One or more filter values.
* `exact_match` - (Optional, Bool) Fuzzy search.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `machines` - List of hosts.
  * `agent_status` - ONLINE: Protected; OFFLINE: Offline; UNINSTALLED: Not installed.
  * `agent_version` - Host security agent version.
  * `baseline_num` - Number of baseline risks.
  * `cloud_tags` - Cloud Tag Information
Note: This field may return null, indicating that no valid values can be obtained.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `cyber_attack_num` - Number of network risks.
  * `has_asset_scan` - Whether there is an available asset scanning API: 0 - no; 1 - yes.
  * `instance_id` - Instance ID.
  * `instance_state` - Instance status: TERMINATED_PRO_VERSION - terminated.
  * `instance_status` - RUNNING; STOPPED; EXPIRED (awaiting recycling).
  * `invasion_num` - Number of intrusion events.
  * `ip_list` - Host IP List
Note: This field may return null, indicating that no valid values can be obtained.
  * `is_added_on_the_fifteen` - Whether a host added within the last 15 days: 0: no; 1: yes
Note: This field may return null, indicating that no valid values can be obtained.
  * `is_pro_version` - Whether the edition is Pro Edition
<li>true: yes</li>
<li>false: no</li>.
  * `kernel_version` - Kernel version.
  * `license_status` - Tamper-proof; authorization status: 1 - authorized; 0 - unauthorized.
  * `machine_extra_info` - Additional information
Note: This field may return null, indicating that no valid values can be obtained.
    * `host_name` - Host name
Note: This field may return null, indicating that no valid values can be obtained.
    * `instance_id` - Instance ID
Note: This field may return null, indicating that no valid values can be obtained.
    * `network_name` - Network Name, returns vpc_id in the case of a VPC network
Note: This field may return null, indicating that no valid values can be obtained.
    * `network_type` - Network Type. 1: VPC network; 2: Basic Network; 3: Non-Tencent Cloud Network
Note: This field may return null, indicating that no valid values can be obtained.
    * `private_ip` - Private IP address
Note: This field may return null, indicating that no valid values can be obtained.
    * `wan_ip` - Public IP address
Note: This field may return null, indicating that no valid values can be obtained.
  * `machine_ip` - Host IP.
  * `machine_name` - Host name.
  * `machine_os` - Host System.
  * `machine_status` - Host status
<li>OFFLINE: Offline</li>
<li>ONLINE: Online</li>
<li>SHUTDOWN: Shut down</li>
<li>UNINSTALLED: Unprotected</li>.
  * `machine_type` - Machine Zone Type. CVM - Cloud Virtual Machine; BM: Bare Metal; ECM: Edge Computing Machine; LH: Lightweight Application Server; Other: Hybrid Cloud Zone.
  * `machine_wan_ip` - Public IP address of a host.
  * `malware_num` - Number of Trojans.
  * `pay_mode` - Host status
<li>POSTPAY: postpaid, indicating pay-as-you-go mode  </li>
<li>PREPAY: prepaid, indicating monthly subscription mode</li>.
  * `project_id` - Project ID.
  * `protect_type` - Protection version: BASIC_VERSION - Basic Edition; PRO_VERSION - Professional Edition; Flagship - Ultimate Edition; GENERAL_DISCOUNT - Inclusive Edition.
  * `quuid` - CVM or BM Machine Unique UUID.
  * `region_info` - Region information.
    * `region_code` - Region code, such as gz, sh, and bj.
    * `region_id` - Region ID.
    * `region_name_en` - English name of the region.
    * `region_name` - Chinese name of a region, such as South China (Guangzhou), East China (Shanghai Finance), and North China (Beijing).
    * `region` - Region identifiers, such as ap-guangzhou, ap-shanghai, and ap-beijing.
  * `remark` - Remarks
Note: This field may return null, indicating that no valid values can be obtained.
  * `security_status` - Risk status
<li>SAFE: Safe</li>
<li>RISK: Risk</li>
<li>UNKNOWN: Unknown</li>.
  * `tag` - Tag information.
    * `name` - Tag name.
    * `rid` - Associated tag ID.
    * `tag_id` - Tag ID.
  * `uuid` - Yunjing client UUID. If the client is offline for a long time, an empty string is returned.
  * `vpc_id` - Network
Note: This field may return null, indicating that no valid values can be obtained.
  * `vul_num` - Number of vulnerabilities.


