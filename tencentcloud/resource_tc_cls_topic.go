/*
Provides a resource to create a cls topic.

Example Usage

```hcl
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "topic"
  logset_id            = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
```

Import

cls topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_topic.topic 2f5764c1-c833-44c5-84c7-950979b2a278
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsTopicCreate,
		Read:   resourceTencentCloudClsTopicRead,
		Delete: resourceTencentCloudClsTopicDelete,
		Update: resourceTencentCloudClsTopicUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset ID.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic name.",
			},
			"partition_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of log topic partitions. Default value: 1. Maximum value: 10.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list. Up to 10 tag key-value pairs are supported and must be unique.",
			},
			"auto_split": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable automatic split. Default value: true.",
			},
			"max_split_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Maximum number of partitions to split into for this topic if" +
					" automatic split is enabled. Default value: 50.",
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Log topic storage class. Valid values: hot: real-time storage; cold: offline storage. Default value: hot. If cold is passed in, " +
					"please contact the customer service to add the log topic to the allowlist first..",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Lifecycle in days. Value range: 1~366. Default value: 30.",
			},
		},
	}
}

func resourceTencentCloudClsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_topic.create")()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateTopicRequest()
		response *cls.CreateTopicResponse
	)

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("partition_count"); ok {
		request.PartitionCount = helper.IntInt64(v.(int))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &k,
				Value: &v,
			})
		}
	}

	request.AutoSplit = helper.Bool(d.Get("auto_split").(bool))

	if v, ok := d.GetOk("max_split_partitions"); ok {
		request.MaxSplitPartitions = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls topic failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.TopicId
	d.SetId(id)
	return resourceTencentCloudClsTopicRead(d, meta)
}

func resourceTencentCloudClsTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_topic.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	topic, err := service.DescribeClsTopicById(ctx, id)

	if err != nil {
		return err
	}

	if topic == nil {
		d.SetId("")
		return fmt.Errorf("resource `Topic` %s does not exist", id)
	}

	_ = d.Set("logset_id", topic.LogsetId)
	_ = d.Set("topic_name", topic.TopicName)
	_ = d.Set("partition_count", topic.PartitionCount)

	tags := make(map[string]string, len(topic.Tags))
	for _, tag := range topic.Tags {
		tags[*tag.Key] = *tag.Value
	}
	_ = d.Set("tags", tags)
	_ = d.Set("auto_split", topic.AutoSplit)
	_ = d.Set("max_split_partitions", topic.MaxSplitPartitions)
	_ = d.Set("storage_type", topic.StorageType)
	_ = d.Set("period", topic.Period)

	return nil
}

func resourceTencentCloudClsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_topic.update")()
	logId := getLogId(contextNil)
	request := cls.NewModifyTopicRequest()

	request.TopicId = helper.String(d.Id())

	if d.HasChange("topic_name") {
		request.TopicName = helper.String(d.Get("topic_name").(string))
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})
		request.Tags = make([]*cls.Tag, 0, len(tags))
		for k, v := range tags {
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &k,
				Value: helper.String(v.(string)),
			})
		}
	}

	if d.HasChange("auto_split") {
		request.AutoSplit = helper.Bool(d.Get("auto_split").(bool))
	}

	if d.HasChange("max_split_partitions") {
		request.MaxSplitPartitions = helper.IntInt64(d.Get("max_split_partitions").(int))
	}

	if d.HasChange("period") {
		request.Period = helper.IntInt64(d.Get("period").(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyTopic(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsTopicRead(d, meta)
}

func resourceTencentCloudClsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_topic.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsTopic(ctx, id); err != nil {
		return err
	}

	return nil
}
