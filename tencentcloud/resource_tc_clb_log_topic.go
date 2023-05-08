/*
Provides a resource to create a CLB instance topic.

Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  topic_name = "clb-topic"
}
```

Import

CLB log topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_log_topic.topic lb-7a0t6zqb
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var clsActionMu = &sync.Mutex{}

func resourceTencentCloudClbLogTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceTopicCreate,
		Read:   resourceTencentCloudClbInstanceTopicRead,
		Delete: resourceTencentCloudClbInstanceTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"log_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic of CLB instance.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic of CLB instance.",
			},
			//"partition_count": {
			//	Type:         schema.TypeInt,
			//	Optional:     true,
			//	ValidateFunc: validateIntegerInRange(1, 10),
			//	Description:  "Topic partition count of CLB instance.(Default 1).",
			//},
			//compute
			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The status of log topic.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Log topic creation time.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceTopicCreate(d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if v, ok := d.GetOk("log_set_id"); ok {
		info, err := clsService.DescribeClsLogset(ctx, v.(string))
		if err != nil {
			return err
		}
		if info == nil {
			return fmt.Errorf("resource `log_set` %s does not exist", v.(string))
		}
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	params := make(map[string]interface{})
	if topicName, ok := d.GetOk("topic_name"); ok {
		params["topic_name"] = topicName
	}
	if partitionCount, ok := d.GetOk("partition_count"); ok {
		params["partition_count"] = partitionCount
	}
	resp, err := clbService.CreateTopic(ctx, params)
	if err != nil {
		log.Printf("[CRITAL]%s create clb topic failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(*resp.Response.TopicId)
	return resourceTencentCloudClbInstanceTopicRead(d, meta)
}

func resourceTencentCloudClbInstanceTopicRead(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	res, err := clsService.DescribeClsTopicById(ctx, id)
	if err != nil {
		return err
	}
	if res == nil {
		d.SetId("")
		return fmt.Errorf("resource `logTopic` %s does not exist", id)
	}
	_ = d.Set("log_set_id", res.LogsetId)
	_ = d.Set("topic_name", res.TopicName)
	_ = d.Set("create_time", res.CreateTime)
	_ = d.Set("status", res.Status)
	return nil

}

func resourceTencentCloudClbInstanceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := clsService.DeleteClsTopic(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
