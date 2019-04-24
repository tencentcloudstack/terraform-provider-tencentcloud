---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_instance"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_instance"
description: |-
 Use this data source to get information about a MySQL instance.
---

#tencentcloud_mysql_instance##

Use this data source to get information about a MySQL instance

##Example Usage

```
data "tencentcloud_mysql_instance" "database"{
     mysql_id = "my-test-database" 
}
```


##Argument Reference

The following arguments are supported:

- `mysql_id` - (Required) Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.
- `instance_role` - (Optional) Instance type. Supported values include: master - master instance, dr - disaster recovery instance, and ro - read-only instance.
- `status` - (Optional) Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 – Isolated. 
- `security_group_id` - (Optional) Security groups ID of instance.
- `instance_name` - (Optional) Name of mysql instance.
- `engine_version` - (Optional) The version number of the database engine to use. Supported versions include 5.5/5.6/5.7.
- `init_flag` - (Optional) Initialization mark. Available values: 0 - Uninitialized; 1 – Initialized.
- `with_dr` - (Optional) Indicates whether to query disaster recovery instances.
- `with_ro` - (Optional) Indicates whether to query read-only instances.
- `with_master` - (Optional) Indicates whether to query master instances.
- `offset` - (Optional) Record offset. Default is 0.
- `limit` - (Optional) Number of results returned for a single request. Default is 20, and maximum is 2000.


##Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `instance_list` - List of instances.
- `zone` - Information of available zone.


For detail information of mysql instance, the following information will be included:

- `cpu_core_count`- CPU count.
- `memory_size` - Memory size (in MB). 
- `volume_size` - Disk capacity (in GB).
- `internet_status` - Status of public network.
- `internet_host` - Public network domain name.
- `internet_port` - Public network port.
- `intranet_ip` - Instance IP for internal access.
- `intranet_port` - Transport layer port number for internal purpose.
- `project_id` - Project ID to which the current instance belongs.
- `vpc_id` - ID of Virtual Private Cloud. 
- `device_type` - Supported instance model.HA - high available version; Basic - basic version.
- `subnet_id` - ID of subnet to which the current instance belongs.
- `slave_sync_mode` -  Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.
- `create_time` - the time at which a instance is created.
- `ro_instance_ids` - ID list of read-only type associated with the current instance.
- `dr_instance_ids` - ID list of disaster-recovory type associated with the current instance.