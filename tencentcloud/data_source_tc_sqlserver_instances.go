/*
Use this data source to query SQL Server instances

Example Usage

Filter instance by Id

```hcl
data "tencentcloud_sqlserver_instances" "example_id" {
  id = "mssql-3l3fgqn7"
}
```

Filter instance by project Id

```hcl
data "tencentcloud_sqlserver_instances" "example_project" {
  project_id = 0
}
```

Filter instance by VPC/Subnet

```hcl
data "tencentcloud_sqlserver_instances" "example_vpc" {
  vpc_id    = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the SQL Server instance to be query.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SQL Server instance to be query.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID of the SQL Server instance to be query.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vpc ID of the SQL Server instance to be query.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet ID of the SQL Server instance to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SQL Server instances. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SQL Server instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SQL Server instance.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay type of the SQL Server instance. For now, only `POSTPAID_BY_HOUR` is valid.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.",
						},
						"ha_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type. `DUAL` (dual-server high availability), `CLUSTER` (cluster).",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of subnet.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID, default value is 0.",
						},
						"ro_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Readonly flag. `RO` (read-only instance), `MASTER` (primary instance with read-only instances). If it is left empty, it refers to an instance which is not read-only and has no RO group.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"used_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used storage.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP for private access.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port for private access.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the SQL Server instance.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the SQL Server instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSqlserverInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	service := SqlserverService{client: tcClient}

	id := d.Get("id").(string)

	name := d.Get("name").(string)

	projectId := -1
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(int)
	}

	vpcId := d.Get("vpc_id").(string)

	subnetId := d.Get("subnet_id").(string)

	instanceList, err := service.DescribeSqlserverInstances(ctx, id, name, projectId, vpcId, subnetId, 1)

	if err != nil {
		instanceList, err = service.DescribeSqlserverInstances(ctx, id, name, projectId, vpcId, subnetId, 1)
	}
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	list := make([]map[string]interface{}, 0, len(instanceList))

	for _, v := range instanceList {
		listItem := make(map[string]interface{})
		listItem["id"] = v.InstanceId
		listItem["name"] = v.Name
		listItem["project_id"] = v.ProjectId
		listItem["storage"] = v.Storage
		listItem["memory"] = v.Memory
		listItem["availability_zone"] = v.Zone
		listItem["create_time"] = v.CreateTime
		listItem["vpc_id"] = v.UniqVpcId
		listItem["subnet_id"] = v.UniqSubnetId
		listItem["engine_version"] = v.Version
		listItem["vip"] = v.Vip
		listItem["vport"] = v.Vport
		listItem["used_storage"] = v.UsedStorage
		listItem["status"] = v.Status
		listItem["ro_flag"] = v.ROFlag

		if *v.PayMode == 1 {
			listItem["charge_type"] = COMMON_PAYTYPE_PREPAID
		} else {
			listItem["charge_type"] = COMMON_PAYTYPE_POSTPAID
		}

		tagList, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, *v.InstanceId)
		if err != nil {
			return err
		}

		listItem["tags"] = tagList
		list = append(list, listItem)
		ids = append(ids, *v.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("instance_list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}

	return nil

}
