/*
Provides a resource to create a CLB instance topic.

Example Usage

```hcl
resource "tencentcloud_clb_instances_topic" "topic" {
    topic_name="clb-topic"
    partition_count=3
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudClbInstanceTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceTopicCreate,
		Read:   resourceTencentCloudClbInstanceTopicRead,
		Update: resourceTencentCloudClbInstanceTopicUpdate,
		Delete: resourceTencentCloudClbInstanceTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic of CLB instance.",
			},
			"partition_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 10),
				Description:  "Topic partition count of CLB instance.(Default 1).",
			},
		},
	}
}

func resourceTencentCloudClbInstanceTopicRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceTencentCloudClbInstanceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceTencentCloudClbInstanceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceTencentCloudClbInstanceTopicCreate(d *schema.ResourceData, meta interface{}) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	params := make(map[string]interface{})
	if topicName, ok := d.GetOk("topic_name"); ok {
		params["topic_name"] = topicName
	}
	if partitionCount, ok := d.GetOk("partition_count"); ok {
		params["partition_count"] = partitionCount
	}
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := clbService.CreateTopic(ctx, params)
	if err != nil {
		log.Printf("[CRITAL]%s create CLB topic failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
