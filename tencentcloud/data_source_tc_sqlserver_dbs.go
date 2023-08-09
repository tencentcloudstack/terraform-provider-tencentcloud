/*
Use this data source to query DB resources for the specific SQL Server instance.

Example Usage

```hcl
data "tencentcloud_sqlserver_dbs" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
}

data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
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
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentSqlserverDBs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverDBRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL Server instance ID which DB belongs to.",
			},
			// Computed
			"db_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of dbs belong to the specific instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL Server instance ID which DB belongs to.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of DB.",
						},
						"charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark of the DB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database creation time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database status. Valid values are `creating`, `running`, `modifying`, `dropping`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverDBRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencent_sqlserver_dbs.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	// precheck
	instanceId := d.Get("instance_id").(string)
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s SQL Server instance %s dose not exist", logId, instanceId)
	}

	dbInfos, err := sqlserverService.DescribeDBsOfInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	var dbList []map[string]interface{}
	for _, item := range dbInfos {
		var dbInfo = make(map[string]interface{})
		dbInfo["name"] = item.Name
		dbInfo["charset"] = item.Charset
		dbInfo["remark"] = item.Remark
		dbInfo["create_time"] = item.CreateTime
		dbInfo["status"] = SQLSERVER_DB_STATUS[*item.Status]
		dbList = append(dbList, dbInfo)
	}
	_ = d.Set("db_list", dbList)
	d.SetId(instanceId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), dbList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}
	}
	return nil
}
