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
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var clsActionMu = &sync.Mutex{}

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
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Fetch topic info pagination limit.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Fetch topic info pagination offset.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceTopicRead(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var topicName string
	if value, ok := d.GetOk("topic_name"); ok {
		topicName = value.(string)
	}
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	res, ok := clsService.DescribeTopicsByTopicName(ctx, topicName)
	if ok != nil {
		return ok
	}
	_ = d.Set("logset_id", res.LogsetId)
	_ = d.Set("topic_id", res.TopicId)
	_ = d.Set("topic_name", res.TopicName)
	_ = d.Set("partition_count", res.PartitionCount)
	_ = d.Set("index", res.Index)
	_ = d.Set("create_time", res.CreateTime)
	_ = d.Set("status", res.Status)
	_ = d.Set("tags", res.Tags)
	_ = d.Set("auto_split", res.AutoSplit)
	_ = d.Set("max_split_partitions", res.MaxSplitPartitions)
	_ = d.Set("storage_type", res.StorageType)
	_ = d.Set("period", res.Period)
	return nil

}

func resourceTencentCloudClbInstanceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTencentCloudClbInstanceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var topicName string
	if value, ok := d.GetOk("topic_name"); ok {
		topicName = value.(string)
	}
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	_, ok := clsService.DeleteTopicsByTopicName(ctx, topicName)
	if ok != nil {
		return ok
	}
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
	resp, err := clbService.CreateTopic(ctx, params)
	d.SetId(*resp.Response.TopicId)
	if err != nil {
		log.Printf("[CRITAL]%s create CLB topic failed, reason:%+v", logId, err)
		return err
	}
	return resourceTencentCloudClbInstanceTopicRead(d, meta)
}
