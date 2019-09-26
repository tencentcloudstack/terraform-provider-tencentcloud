---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_instance"
sidebar_current: "docs-tencentcloud-datasource-mysql_instance"
description: |-
  Use this data source to get information about a MySQL instance.
---

# tencentcloud_mysql_instance

Use this data source to get information about a MySQL instance.

## Example Usage

```hcl
data "tencentcloud_mysql_instance" "database" {
  mysql_id           = "my-test-database"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Optional) The version number of the database engine to use. Supported versions include 5.5/5.6/5.7.
* `init_flag` - (Optional) Initialization mark. Available values: 0 - Uninitialized; 1 - Initialized.
* `instance_name` - (Optional) Name of mysql instance.
* `instance_role` - (Optional) Instance type. Supported values include: master - master instance, dr - disaster recovery instance, and ro - read-only instance.
* `limit` - (Optional) Number of results returned for a single request. Default is 20, and maximum is 2000.
* `mysql_id` - (Optional) Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.
* `offset` - (Optional) Record offset. Default is 0.
* `result_output_file` - (Optional) Used to store results.
* `security_group_id` - (Optional) Security groups ID of instance.
* `status` - (Optional) Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 - Isolated.
* `with_dr` - (Optional) Indicates whether to query disaster recovery instances.
* `with_master` - (Optional) Indicates whether to query master instances.
* `with_ro` - (Optional) Indicates whether to query read-only instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of instances. Each element contains the following attributes:
  * `cpu_core_count` - CPU count.
  * `create_time` - The time at which a instance is created.
  * `device_type` - Supported instance model.HA - high available version; Basic - basic version.
  * `dr_instance_ids` - ID list of disaster-recovory type associated with the current instance.
  * `engine_version` - The version number of the database engine to use. Supported versions include 5.5/5.6/5.7.
  * `init_flag` - Initialization mark. Available values: 0 - Uninitialized; 1 - Initialized.
  * `instance_name` - Name of mysql instance.
  * `instance_role` - Instance type. Supported values include: master - master instance, dr - disaster recovery instance, and ro - read-only instance.
  * `internet_host` - Public network domain name.
  * `internet_port` - Public network port.
  * `internet_status` - Status of public network.
  * `intranet_ip` - Instance IP for internal access.
  * `intranet_port` - Transport layer port number for internal purpose.
  * `memory_size` - Memory size (in MB).
  * `mysql_id` - Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.
  * `project_id` - Project ID to which the current instance belongs.
  * `ro_instance_ids` - ID list of read-only type associated with the current instance.
  * `slave_sync_mode` - Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.
  * `status` - Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 - Isolated.
  * `subnet_id` - ID of subnet to which the current instance belongs.
  * `volume_size` - Disk capacity (in GB).
  * `vpc_id` - ID of Virtual Private Cloud.
  * `zone` - Information of available zone.


