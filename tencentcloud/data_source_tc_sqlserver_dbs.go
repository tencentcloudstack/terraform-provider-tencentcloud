/*
Use this data source to query DB resources for the specific SQLServer instance.

Example Usage

```hcl
variable "availability_zone"{
  default = "ap-guangzhou-2"
}

resource "tencentcloud_vpc" "foo" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.foo.id
  subnet_id         = tencentcloud_subnet.foo.id
  engine_version    = "2008R2"
  project_id        = 0
  memory            = 2
  storage           = 10
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name = "example"
  charset = "Chinese_PRC_BIN"
  remark = "test-remark"
}

data "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description: "SQLServer instance ID which DB belongs to.",
			},
			// Computed
			"db_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of dbs belong to the specific instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of DB.",
						},
						"charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.",
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
	defer logElapsed("data_source.tencent_sqlserver_db.read")()

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
		return fmt.Errorf("[CRITAL]%s SQLServer instance %s dose not exist", logId, instanceId)
	}
	_ = d.Set("instance_id", instanceId)

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
