/*
Provides a resource to create a sqlserver config_database_cdc

Example Usage

Turn off database data change capture (CDC)

```hcl
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

resource "tencentcloud_sqlserver_config_database_cdc" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  db_name     = tencentcloud_sqlserver_db.example.name
  modify_type = "disable"
}
```

Enable Database Data Change Capture (CDC)

```hcl
resource "tencentcloud_sqlserver_config_database_cdc" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  db_name     = tencentcloud_sqlserver_db.example.name
  modify_type = "enable"
}
```

Import

sqlserver config_database_cdc can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_cdc.example mssql-i9ma6oy7#tf_example_db
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverConfigDatabaseCDC() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigDatabaseCDCCreate,
		Read:   resourceTencentCloudSqlserverConfigDatabaseCDCRead,
		Update: resourceTencentCloudSqlserverConfigDatabaseCDCUpdate,
		Delete: resourceTencentCloudSqlserverConfigDatabaseCDCDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "database name.",
			},
			"modify_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Enable or disable CDC. Valid values: enable, disable.",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigDatabaseCDCCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_cdc.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
		dbName     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("db_name"); ok {
		dbName = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, dbName}, FILED_SP))

	return resourceTencentCloudSqlserverConfigDatabaseCDCUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseCDCRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_cdc.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		Name       string
		modifyType string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	dbName := idSplit[1]

	configDatabaseCDC, err := service.DescribeSqlserverConfigDatabaseCDCById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configDatabaseCDC == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigDatabaseCDC` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, i := range configDatabaseCDC {
		if *i.Name == dbName {
			if *i.IsCdcEnabled == "0" {
				modifyType = "disable"
			} else {
				modifyType = "enable"
			}
			Name = *i.Name
			break
		}
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("db_name", Name)
	_ = d.Set("modify_type", modifyType)

	return nil
}

func resourceTencentCloudSqlserverConfigDatabaseCDCUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_cdc.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = sqlserver.NewModifyDatabaseCDCRequest()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		flowId      int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	dbName := idSplit[1]

	if v, ok := d.GetOk("modify_type"); ok {
		request.ModifyType = helper.String(v.(string))
	}

	request.InstanceId = &instanceId
	dbNames := make([]*string, 0)
	dbNames = append(dbNames, &dbName)
	request.DBNames = dbNames

	if v, ok := d.GetOk("modify_type"); ok {
		request.ModifyType = helper.String(v.(string))
	}

	request.InstanceId = &instanceId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDatabaseCDC(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver configDatabaseCDC not exists")
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configDatabaseCDC failed, reason:%+v", logId, err)
		return err
	}

	flowRequest.FlowId = &flowId
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return retryError(e)
		}
		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver configDatabaseCDC status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver configDatabaseCDC status is fail"))
		} else {
			e = fmt.Errorf("sqlserver configDatabaseCDC status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver configDatabaseCDC failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigDatabaseCDCRead(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseCDCDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_cdc.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
