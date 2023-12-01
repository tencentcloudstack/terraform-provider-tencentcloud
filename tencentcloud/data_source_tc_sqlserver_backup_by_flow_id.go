/*
Use this data source to query detailed information of sqlserver datasource_backup_by_flow_id

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

data "tencentcloud_sqlserver_backup_by_flow_id" "example" {
  instance_id = tencentcloud_sqlserver_general_backup.example.instance_id
  flow_id     = tencentcloud_sqlserver_general_backup.example.flow_id
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.id
  backup_name = "tf_example_backup"
  strategy    = 0
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func dataSourceTencentCloudSqlserverBackupByFlowId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverBackupByFlowIdRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Create a backup process ID, which can be obtained through the [CreateBackup](https://cloud.tencent.com/document/product/238/19946) interface.",
			},
			"file_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "File name. For a single-database backup file, only the file name of the first record is returned; for a single-database backup file, the file names of all records need to be obtained through the DescribeBackupFiles interface.",
			},
			"backup_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup task name, customizable.",
			},
			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "backup start time.",
			},
			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "backup end time.",
			},
			"strategy": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup strategy, 0-instance backup; 1-multi-database backup; when the instance status is 0-creating, this field is the default value 0, meaningless.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup file status, 0-creating; 1-success; 2-failure.",
			},
			"backup_way": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup method, 0-scheduled backup; 1-manual temporary backup; instance status is 0-creating, this field is the default value 0, meaningless.",
			},
			"dbs": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "For the DB list, only the library name contained in the first record is returned for a single-database backup file; for a single-database backup file, the library names of all records need to be obtained through the DescribeBackupFiles interface.",
			},
			"internal_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet download address, for a single database backup file, only the intranet download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.",
			},
			"external_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "External network download address, for a single database backup file, only the external network download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.",
			},
			"group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Aggregate Id, this value is not returned for packaged backup files. Use this value to call the DescribeBackupFiles interface to obtain the detailed information of a single database backup file.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverBackupByFlowIdRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_backup_by_flow_id.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		response   *sqlserver.DescribeBackupByFlowIdResponseParams
		instanceId string
		flowId     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("flow_id"); ok {
		flowId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBackupByFlowId(ctx, instanceId, flowId)
		if e != nil {
			return retryError(e)
		}

		response = result.Response
		return nil
	})

	if err != nil {
		return err
	}

	if response.FileName != nil {
		_ = d.Set("file_name", response.FileName)
	}

	if response.BackupName != nil {
		_ = d.Set("backup_name", response.BackupName)
	}

	if response.StartTime != nil {
		_ = d.Set("start_time", response.StartTime)
	}

	if response.EndTime != nil {
		_ = d.Set("end_time", response.EndTime)
	}

	if response.Strategy != nil {
		_ = d.Set("strategy", response.Strategy)
	}

	if response.Status != nil {
		_ = d.Set("status", response.Status)
	}

	if response.BackupWay != nil {
		_ = d.Set("backup_way", response.BackupWay)
	}

	if response.DBs != nil {
		_ = d.Set("dbs", response.DBs)
	}

	if response.InternalAddr != nil {
		_ = d.Set("internal_addr", response.InternalAddr)
	}

	if response.ExternalAddr != nil {
		_ = d.Set("external_addr", response.ExternalAddr)
	}

	if response.GroupId != nil {
		_ = d.Set("group_id", response.GroupId)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
