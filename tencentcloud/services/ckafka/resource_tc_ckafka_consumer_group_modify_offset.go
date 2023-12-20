package ckafka

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaConsumerGroupModifyOffset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaConsumerGroupModifyOffsetCreate,
		Read:   resourceTencentCloudCkafkaConsumerGroupModifyOffsetRead,
		Delete: resourceTencentCloudCkafkaConsumerGroupModifyOffsetDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Kafka instance id.",
			},

			"group": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "kafka group.",
			},

			"strategy": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeInt,
				Description: "Reset the policy of offset.\n" +
					"`0`: Move the offset forward or backward shift bar;\n" +
					"`1`: Alignment reference (by-duration,to-datetime,to-earliest,to-latest), which means moving the offset to the location of the specified timestamp;\n" +
					"`2`: Alignment reference (to-offset), which means to move the offset to the specified offset location.",
			},

			"topics": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Indicates the topics that needs to be reset. Leave it empty means all.",
			},

			"shift": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "This field must be included when strategy is 0. If it is greater than zero, the offset will be moved backward by shift bars, and if it is less than zero, the offset will be traced back to the number of shift entries. After the correct reset, the new offset should be (old_offset + shift). It should be noted that if the new offset is less than partition's earliest, it will be set to earliest, and if the latest greater than partition will be set to latest.",
			},

			"shift_timestamp": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Unit ms. When strategy is 1, you must include this field, where-2 means to reset the offset to the beginning,-1 means to reset to the latest position (equivalent to emptying), and other values represent the specified time. You will get the offset of the specified time in the topic and then reset it. If there is no message at the specified time, get the last offset.",
			},

			"offset": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The offset location that needs to be reset. When strategy is 2, this field must be included.",
			},

			"partitions": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The list of partition that needs to be reset if no Topics parameter is specified. Resets the partition in the corresponding Partition list of all topics. When Topics is specified, the partition of the corresponding topic list of the specified Partitions list is reset.",
			},
		},
	}
}

func resourceTencentCloudCkafkaConsumerGroupModifyOffsetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group_modify_offset.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ckafka.NewModifyGroupOffsetsRequest()
		instanceId string
		group      string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("group"); ok {
		group = v.(string)
		request.Group = helper.String(group)
	}

	if v, _ := d.GetOk("strategy"); v != nil {
		request.Strategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("topics"); ok {
		topicsSet := v.(*schema.Set).List()
		for i := range topicsSet {
			topics := topicsSet[i].(string)
			request.Topics = append(request.Topics, &topics)
		}
	}

	if v, _ := d.GetOk("shift"); v != nil {
		request.Shift = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shift_timestamp"); v != nil {
		request.ShiftTimestamp = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		request.Offset = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("partitions"); ok {
		partitionsSet := v.(*schema.Set).List()
		for i := range partitionsSet {
			partitions := partitionsSet[i].(int)
			request.Partitions = append(request.Partitions, helper.IntInt64(partitions))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().ModifyGroupOffsets(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ckafka consumerGroupModifyOffset failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + group)

	return resourceTencentCloudCkafkaConsumerGroupModifyOffsetRead(d, meta)
}

func resourceTencentCloudCkafkaConsumerGroupModifyOffsetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group_modify_offset.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCkafkaConsumerGroupModifyOffsetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group_modify_offset.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
