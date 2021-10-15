/*
Provide a resource to create a TDMQ topic.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark = "this is description."
}

resource "tencentcloud_tdmq_namespace" "bar" {
  environ_name = "example"
  msg_ttl = 300
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark = "this is description."
}

resource "tencentcloud_tdmq_topic" "bar" {
  environ_id = "${tencentcloud_tdmq_namespace.bar.id}"
  topic_name = "example"
  partitions = 6
  topic_type = 0
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark = "this is description."
}
```

Import

Tdmq Topic can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_topic.test topic_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudTdmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqTopicCreate,
		Read:   resourceTencentCloudTdmqTopicRead,
		Update: resourceTencentCloudTdmqTopicUpdate,
		Delete: resourceTencentCloudTdmqTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environ_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of tdmq namespace.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of topic to be created.",
			},
			"partitions": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The partitions of topic.",
			},
			"topic_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The type of topic.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Dedicated Cluster Id.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the namespace.",
			},

			//compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
		},
	}
}

func resourceTencentCloudTdmqTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_topic.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		environId  string
		topicName  string
		partitions uint64
		topicType  uint64
		remark     string
		clusterId  string
	)
	if temp, ok := d.GetOk("environ_id"); ok {
		environId = temp.(string)
		if len(environId) < 1 {
			return fmt.Errorf("environ_id should be not empty string")
		}
	}
	if temp, ok := d.GetOk("topic_name"); ok {
		topicName = temp.(string)
		if len(topicName) < 1 {
			return fmt.Errorf("topic_name should be not empty string")
		}
	}
	partitions = uint64(d.Get("partitions").(int))
	topicType = uint64(d.Get("topic_type").(int))
	if temp, ok := d.GetOk("remark"); ok {
		remark = temp.(string)
	}
	if temp, ok := d.GetOk("cluster_id"); ok {
		clusterId = temp.(string)
	}

	err := tdmqService.CreateTdmqTopic(ctx, environId, topicName, partitions, topicType, remark, clusterId)
	if err != nil {
		return err
	}
	d.SetId(topicName)

	return resourceTencentCloudTdmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	topicName := d.Id()
	environId := d.Get("environ_id").(string)
	clusterId := d.Get("cluster_id").(string)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqTopicById(ctx, environId, topicName, clusterId)
		if e != nil {
			return retryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("partitions", info.Partitions)
		_ = d.Set("topic_type", info.TopicType)
		_ = d.Set("remark", info.Remark)
		_ = d.Set("create_time", info.CreateTime)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTdmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_topic.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	topicName := d.Id()
	environId := d.Get("environ_id").(string)
	clusterId := d.Get("cluster_id").(string)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		partitions uint64
		remark     string
	)
	old, now := d.GetChange("partitions")
	if d.HasChange("partitions") {
		partitions = uint64(now.(int))
	} else {
		partitions = uint64(old.(int))
	}

	old, now = d.GetChange("remark")
	if d.HasChange("remark") {
		remark = now.(string)
	} else {
		remark = old.(string)
	}

	d.Partial(true)

	if err := service.ModifyTdmqTopicAttribute(ctx, environId, topicName,
		partitions, remark, clusterId); err != nil {
		return err
	}
	d.SetPartial("partitions")
	d.SetPartial("remark")

	d.Partial(false)
	return resourceTencentCloudTdmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	topicName := d.Id()
	environId := d.Get("environ_id").(string)
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqTopic(ctx, environId, topicName, clusterId); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
