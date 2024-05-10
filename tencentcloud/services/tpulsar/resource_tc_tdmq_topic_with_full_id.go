package tpulsar

import (
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func ResourceTencentCloudTdmqTopicWithFullId() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqTopicWithFullIdCreate,
		Read:   resourceTencentCloudTdmqTopicWithFullIdRead,
		Update: resourceTencentCloudTdmqTopicWithFullIdUpdate,
		Delete: resourceTencentCloudTdmqTopicWithFullIdDelete,
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
				Optional:    true,
				Computed:    true,
				Deprecated:  "This input will be gradually discarded and can be switched to PulsarTopicType parameter 0: Normal message; 1: Global sequential messages; 2: Local sequential messages; 3: Retrying queue; 4: Dead letter queue.",
				Description: "The type of topic.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Dedicated Cluster Id.",
			},
			"pulsar_topic_type": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"topic_type"},
				Description:   "Pulsar Topic Type 0: Non-persistent non-partitioned 1: Non-persistent partitioned 2: Persistent non-partitioned 3: Persistent partitioned.",
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

func resourceTencentCloudTdmqTopicWithFullIdCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_topic_with_full_id.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tdmqService := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var (
		environId       string
		topicName       string
		partitions      uint64
		topicType       int64
		remark          string
		clusterId       string
		pulsarTopicType int64
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
	if temp, ok := d.GetOk("remark"); ok {
		remark = temp.(string)
	}
	if temp, ok := d.GetOk("cluster_id"); ok {
		clusterId = temp.(string)
	}

	if v, ok := d.GetOkExists("pulsar_topic_type"); ok {
		pulsarTopicType = int64(v.(int))
	} else {
		pulsarTopicType = svctdmq.NonePulsarTopicType
		if v, ok := d.GetOkExists("topic_type"); ok {
			topicType = int64(v.(int))
		} else {
			topicType = svctdmq.NoneTopicType
		}
	}

	err := tdmqService.CreateTdmqTopic(ctx, environId, topicName, partitions, topicType, remark, clusterId, pulsarTopicType)
	if err != nil {
		return err
	}
	d.SetId(clusterId + tccommon.FILED_SP + environId + tccommon.FILED_SP + topicName)

	return resourceTencentCloudTdmqTopicWithFullIdRead(d, meta)
}

func resourceTencentCloudTdmqTopicWithFullIdRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_topic_with_full_id.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", d.Id())
	}
	clusterId := items[0]
	environId := items[1]
	topicName := items[2]

	tdmqService := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqTopicById(ctx, environId, topicName, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("cluster_id", clusterId)
		_ = d.Set("environ_id", environId)
		_ = d.Set("topic_name", topicName)
		_ = d.Set("partitions", info.Partitions)
		_ = d.Set("topic_type", info.TopicType)
		_ = d.Set("pulsar_topic_type", info.PulsarTopicType)
		_ = d.Set("remark", info.Remark)
		_ = d.Set("create_time", info.CreateTime)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTdmqTopicWithFullIdUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_topic_with_full_id.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	if d.HasChange("topic_type") {
		return fmt.Errorf("`topic_type` do not support change now.")
	}

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", d.Id())
	}
	clusterId := items[0]
	environId := items[1]
	topicName := items[2]

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

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
	d.Partial(false)
	return resourceTencentCloudTdmqTopicWithFullIdRead(d, meta)
}

func resourceTencentCloudTdmqTopicWithFullIdDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_topic_with_full_id.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", d.Id())
	}
	clusterId := items[0]
	environId := items[1]
	topicName := items[2]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqTopic(ctx, environId, topicName, clusterId); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == svcvpc.VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
