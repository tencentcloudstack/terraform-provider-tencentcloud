/*
Provides a resource to create a mysql local_binlog_config

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_local_binlog_config" "example" {
  instance_id = tencentcloud_mysql_instance.example.id
  save_hours  = 140
  max_usage   = 50
}
```

Import

mysql local_binlog_config can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_local_binlog_config.local_binlog_config instance_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlLocalBinlogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlLocalBinlogConfigCreate,
		Read:   resourceTencentCloudMysqlLocalBinlogConfigRead,
		Update: resourceTencentCloudMysqlLocalBinlogConfigUpdate,
		Delete: resourceTencentCloudMysqlLocalBinlogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.",
			},

			"save_hours": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Retention period of local binlog. Valid range: 72-168 hours. When there is disaster recovery instance, the valid range will be 120-168 hours.",
			},

			"max_usage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Space utilization of local binlog. Value range: [30,50].",
			},
		},
	}
}

func resourceTencentCloudMysqlLocalBinlogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_local_binlog_config.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_id").(string))

	return resourceTencentCloudMysqlLocalBinlogConfigUpdate(d, meta)
}

func resourceTencentCloudMysqlLocalBinlogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_local_binlog_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	localBinlogConfig, err := service.DescribeMysqlLocalBinlogConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if localBinlogConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_mysql_local_binlog_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil

	}

	_ = d.Set("instance_id", instanceId)

	if localBinlogConfig.SaveHours != nil {
		_ = d.Set("save_hours", localBinlogConfig.SaveHours)
	}

	if localBinlogConfig.MaxUsage != nil {
		_ = d.Set("max_usage", localBinlogConfig.MaxUsage)
	}

	return nil
}

func resourceTencentCloudMysqlLocalBinlogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_local_binlog_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mysql.NewModifyLocalBinlogConfigRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, _ := d.GetOk("save_hours"); v != nil {
		request.SaveHours = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("max_usage"); v != nil {
		request.MaxUsage = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyLocalBinlogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql localBinlogConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlLocalBinlogConfigRead(d, meta)
}

func resourceTencentCloudMysqlLocalBinlogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_local_binlog_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
