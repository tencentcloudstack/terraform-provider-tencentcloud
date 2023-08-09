/*
Provides a resource to create a sqlserver renew_postpaid_db_instance

Example Usage

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

resource "tencentcloud_sqlserver_config_terminate_db_instance" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
}

resource "tencentcloud_sqlserver_renew_postpaid_db_instance" "example" {
  instance_id = tencentcloud_sqlserver_config_terminate_db_instance.example.id
}
```

Import

sqlserver renew_postpaid_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_renew_postpaid_db_instance.example mssql-i9ma6oy7
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func resourceTencentCloudSqlserverRenewPostpaidDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRenewPostpaidDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverRenewPostpaidDBInstanceRead,
		Update: resourceTencentCloudSqlserverRenewPostpaidDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRenewPostpaidDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverRenewPostpaidDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	renewPostpaidDBInstance, err := service.DescribeSqlserverRenewPostpaidDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if renewPostpaidDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRenewPostpaidDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if renewPostpaidDBInstance.InstanceId != nil {
		_ = d.Set("instance_id", renewPostpaidDBInstance.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = sqlserver.NewRenewPostpaidDBInstanceRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RenewPostpaidDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver renewPostpaidDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRenewPostpaidDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
