/*
Provides a resource to create a sqlserver start_xevent

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

resource "tencentcloud_sqlserver_start_xevent" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  event_config {
    event_type = "slow"
    threshold  = 0
  }
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverStartXevent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverStartXeventCreate,
		Read:   resourceTencentCloudSqlserverStartXeventRead,
		Delete: resourceTencentCloudSqlserverStartXeventDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"event_config": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Whether to start or stop an extended event.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event type. Valid values: slow (set threshold for slow SQL ), blocked (set threshold for the blocking and deadlock).",
						},
						"threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Threshold in milliseconds. Valid values: 0(disable), non-zero (enable).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverStartXeventCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_xevent.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = sqlserver.NewStartInstanceXEventRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("event_config"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			eventConfig := sqlserver.EventConfig{}
			if v, ok := dMap["event_type"]; ok {
				eventConfig.EventType = helper.String(v.(string))
			}
			if v, ok := dMap["threshold"]; ok {
				eventConfig.Threshold = helper.IntInt64(v.(int))
			}
			request.EventConfig = append(request.EventConfig, &eventConfig)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().StartInstanceXEvent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver startXevent failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverStartXeventRead(d, meta)
}

func resourceTencentCloudSqlserverStartXeventRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_xevent.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverStartXeventDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_start_xevent.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
