/*
Provides a resource to create a ckafka renew_ckafka_instance

Example Usage

```hcl
resource "tencentcloud_ckafka_renew_ckafka_instance" "renew_ckafka_instance" {
  instance_id = "InstanceId"
  time_span = 1
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

ckafka renew_ckafka_instance can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_renew_ckafka_instance.renew_ckafka_instance renew_ckafka_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCkafkaRenewCkafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaRenewCkafkaInstanceCreate,
		Read:   resourceTencentCloudCkafkaRenewCkafkaInstanceRead,
		Delete: resourceTencentCloudCkafkaRenewCkafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"time_span": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration, the default is 1, and the unit is month.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudCkafkaRenewCkafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_renew_ckafka_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ckafka.NewRenewCkafkaInstanceRequest()
		response   = ckafka.NewRenewCkafkaInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("time_span"); v != nil {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCkafkaClient().RenewCkafkaInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ckafka renewCkafkaInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCkafkaRenewCkafkaInstanceRead(d, meta)
}

func resourceTencentCloudCkafkaRenewCkafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_renew_ckafka_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCkafkaRenewCkafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_renew_ckafka_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
