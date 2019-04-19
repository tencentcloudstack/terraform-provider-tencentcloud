---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_readonly_instance"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_readonly_instance"
description: |-
Provides a mysql instance resource to create read-only database instances.

---
#tencentcloud_mysql_readonly_instance

Provides a mysql instance resource to create read-only database instances.


##Example Usage

```
resource " tencentcloud_mysql_readonly_instance " "default" {
  master_instance_id = "cdb-dnqksd9f"
  instance_name =" myTestMysql"
  mem_size = 128000 
  volume_size = 255
  vpc_id = "vpc-12mt3l31"
  subnet_id = "subnet-9uivyb1g"
  intranet_port = 3306
  security_groups = ["sg-ot8eclwz"]
  tags = {
     name ="test"
  }
}

```
##Argument Reference

The following arguments are supported:

- `master_instance_id` - (Required) Indicates the master instance ID of recovery instances. 

- `instance_name` - (Required) The name of a mysql instance.

- `mem_size` – (Required) Memory size (in MB).

- `volume_size` – (Required) Disk size (in GB).

- `vpc_id` – (Required) ID of VPC, which can be modified once every 24 hours and can’t be removed.

- `subnet_id` – (Required) Private network ID. If vpc_id is set, this value is required.

- `intranet_port` – (Optional) Public access port, rang form 1024 to 65535 and default value is 3306.

- `security_groups` – (Optional) Security groups to use.

- `tags` – (Optional) Instance tags.


##Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `gtid` – Indicates whether GTID is enable. 0 - Not enabled; 1 - Enabled.  

- `intranet_ip` – Instance IP.

- `locked` - Indicates whether the instance is locked. 0 - No; 1 - Yes.

- `status` - Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 – Isolated. 

- `task_status` – Indicates which kind of operations is being executed.