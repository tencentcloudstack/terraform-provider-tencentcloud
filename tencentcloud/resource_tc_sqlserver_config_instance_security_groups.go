/*
Provides a resource to create a sqlserver config_instance_security_groups

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

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_config_instance_security_groups" "config_instance_security_groups" {
  instance_id           = tencentcloud_sqlserver_basic_instance.example.id
  security_group_id_set = [tencentcloud_security_group.security_group.id]
}
```

Import

sqlserver config_instance_security_groups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_security_groups.config_instance_security_groups config_instance_security_groups_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverConfigInstanceSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigInstanceSecurityGroupsCreate,
		Read:   resourceTencentCloudSqlserverConfigInstanceSecurityGroupsRead,
		Update: resourceTencentCloudSqlserverConfigInstanceSecurityGroupsUpdate,
		Delete: resourceTencentCloudSqlserverConfigInstanceSecurityGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"security_group_id_set": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of security group IDs to modify, an array of one or more security group IDs.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigInstanceSecurityGroupsUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	configInstanceSecurityGroups, err := service.DescribeSqlserverConfigInstanceSecurityGroupsById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configInstanceSecurityGroups == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceSecurityGroups` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	sgList := []interface{}{}
	for _, sg := range configInstanceSecurityGroups {
		sgList = append(sgList, sg.SecurityGroupId)
	}
	_ = d.Set("security_group_id_set", sgList)

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = sqlserver.NewModifyDBInstanceSecurityGroupsRequest()
		instanceId = d.Id()
	)

	if v, ok := d.GetOk("security_group_id_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			request.SecurityGroupIdSet = append(request.SecurityGroupIdSet, helper.String(item.(string)))
		}
	}

	request.InstanceId = &instanceId

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configInstanceSecurityGroups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigInstanceSecurityGroupsRead(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
