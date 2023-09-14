/*
Provides a resource to create a trocket rocketmq_consumer_group

Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxxx"
  subnet_id     = "subnet-xxxxx"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_consumer_group" "rocketmq_consumer_group" {
  instance_id             = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  consumer_group          = "test_consumer_group"
  max_retry_times         = 20
  consume_enable          = false
  consume_message_orderly = true
  remark                  = "test for terraform"
}
```

Import

trocket rocketmq_consumer_group can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group  instanceId#consumerGroup
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
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTrocketRocketmqConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTrocketRocketmqConsumerGroupCreate,
		Read:   resourceTencentCloudTrocketRocketmqConsumerGroupRead,
		Update: resourceTencentCloudTrocketRocketmqConsumerGroupUpdate,
		Delete: resourceTencentCloudTrocketRocketmqConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"consumer_group": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Name of consumer group.",
			},

			"max_retry_times": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Max retry times.",
			},

			"consume_enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable consumption.",
			},

			"consume_message_orderly": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "`true`: Sequential delivery, `false`: Concurrent delivery.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "remark.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = trocket.NewCreateConsumerGroupRequest()
		response      = trocket.NewCreateConsumerGroupResponse()
		instanceId    string
		consumerGroup string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("consumer_group"); ok {
		request.ConsumerGroup = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_retry_times"); ok {
		request.MaxRetryTimes = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("consume_enable"); ok {
		request.ConsumeEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("consume_message_orderly"); ok {
		request.ConsumeMessageOrderly = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTrocketClient().CreateConsumerGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmqConsumerGroup failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	consumerGroup = *response.Response.ConsumerGroup
	d.SetId(instanceId + FILED_SP + consumerGroup)

	return resourceTencentCloudTrocketRocketmqConsumerGroupRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	consumerGroup := idSplit[1]

	rocketmqConsumerGroup, err := service.DescribeTrocketRocketmqConsumerGroupById(ctx, instanceId, consumerGroup)
	if err != nil {
		return err
	}

	if rocketmqConsumerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TrocketRocketmqConsumerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("consumer_group", consumerGroup)

	if rocketmqConsumerGroup.MaxRetryTimes != nil {
		_ = d.Set("max_retry_times", rocketmqConsumerGroup.MaxRetryTimes)
	}

	if rocketmqConsumerGroup.ConsumeEnable != nil {
		_ = d.Set("consume_enable", rocketmqConsumerGroup.ConsumeEnable)
	}

	if rocketmqConsumerGroup.ConsumeMessageOrderly != nil {
		_ = d.Set("consume_message_orderly", rocketmqConsumerGroup.ConsumeMessageOrderly)
	}

	if rocketmqConsumerGroup.Remark != nil {
		_ = d.Set("remark", rocketmqConsumerGroup.Remark)
	}

	return nil
}

func resourceTencentCloudTrocketRocketmqConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := trocket.NewModifyConsumerGroupRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	consumerGroup := idSplit[1]

	request.InstanceId = &instanceId
	request.ConsumerGroup = &consumerGroup

	needChange := false

	mutableArgs := []string{"max_retry_times", "consume_enable", "consume_message_orderly", "remark"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
		}
	}

	if needChange {
		if v, ok := d.GetOkExists("max_retry_times"); ok {
			request.MaxRetryTimes = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("consume_enable"); ok {
			request.ConsumeEnable = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("consume_message_orderly"); ok {
			request.ConsumeMessageOrderly = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTrocketClient().ModifyConsumerGroup(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update trocket rocketmqConsumerGroup failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTrocketRocketmqConsumerGroupRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	consumerGroup := idSplit[1]

	if err := service.DeleteTrocketRocketmqConsumerGroupById(ctx, instanceId, consumerGroup); err != nil {
		return err
	}

	return nil
}
