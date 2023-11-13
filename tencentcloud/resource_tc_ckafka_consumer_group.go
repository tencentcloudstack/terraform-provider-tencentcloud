/*
Provides a resource to create a ckafka consumer_group

Example Usage

```hcl
resource "tencentcloud_ckafka_consumer_group" "consumer_group" {
  instance_id = "InstanceId"
  group_name = "GroupName"
  topic_name = "TopicName"
  topic_name_list =
}
```

Import

ckafka consumer_group can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_consumer_group.consumer_group consumer_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCkafkaConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaConsumerGroupCreate,
		Read:   resourceTencentCloudCkafkaConsumerGroupRead,
		Update: resourceTencentCloudCkafkaConsumerGroupUpdate,
		Delete: resourceTencentCloudCkafkaConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "GroupName.",
			},

			"topic_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Topic name, TopicName or TopicNameList need to display and specify an existing topic name.",
			},

			"topic_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of topic names.",
			},
		},
	}
}

func resourceTencentCloudCkafkaConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_consumer_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ckafka.NewCreateConsumerRequest()
		response   = ckafka.NewCreateConsumerResponse()
		instanceId string
		groupName  string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		groupName = v.(string)
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name_list"); ok {
		topicNameListSet := v.(*schema.Set).List()
		for i := range topicNameListSet {
			topicNameList := topicNameListSet[i].(string)
			request.TopicNameList = append(request.TopicNameList, &topicNameList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCkafkaClient().CreateConsumer(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka consumerGroup failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, groupName}, FILED_SP))

	return resourceTencentCloudCkafkaConsumerGroupRead(d, meta)
}

func resourceTencentCloudCkafkaConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_consumer_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupName := idSplit[1]

	consumerGroup, err := service.DescribeCkafkaConsumerGroupById(ctx, instanceId, groupName)
	if err != nil {
		return err
	}

	if consumerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaConsumerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if consumerGroup.InstanceId != nil {
		_ = d.Set("instance_id", consumerGroup.InstanceId)
	}

	if consumerGroup.GroupName != nil {
		_ = d.Set("group_name", consumerGroup.GroupName)
	}

	if consumerGroup.TopicName != nil {
		_ = d.Set("topic_name", consumerGroup.TopicName)
	}

	if consumerGroup.TopicNameList != nil {
		_ = d.Set("topic_name_list", consumerGroup.TopicNameList)
	}

	return nil
}

func resourceTencentCloudCkafkaConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_consumer_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ckafka.NewModifyGroupOffsetsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupName := idSplit[1]

	request.InstanceId = &instanceId
	request.GroupName = &groupName

	immutableArgs := []string{"instance_id", "group_name", "topic_name", "topic_name_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCkafkaClient().ModifyGroupOffsets(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ckafka consumerGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCkafkaConsumerGroupRead(d, meta)
}

func resourceTencentCloudCkafkaConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_consumer_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupName := idSplit[1]

	if err := service.DeleteCkafkaConsumerGroupById(ctx, instanceId, groupName); err != nil {
		return err
	}

	return nil
}
