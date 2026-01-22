package trocket

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTrocketRocketmqConsumerGroup() *schema.Resource {
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

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tag of consumer group.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().CreateConsumerGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create trocket rocketmq consumer group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmq consumer group failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.InstanceId == nil || response.Response.ConsumerGroup == nil {
		return fmt.Errorf("InstanceId or ConsumerGroup is nil.")
	}

	instanceId = *response.Response.InstanceId
	consumerGroup = *response.Response.ConsumerGroup
	d.SetId(strings.Join([]string{instanceId, consumerGroup}, tccommon.FILED_SP))

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::trocket:%s:uin/:consumerGroup/%s/%s", region, instanceId, consumerGroup)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTrocketRocketmqConsumerGroupRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
		log.Printf("[WARN]%s resource `tencentcloud_trocket_rocketmq_consumer_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
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

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tags, err := tagService.DescribeResourceTags(ctx, "trocket", "consumerGroup", tcClient.Region, fmt.Sprintf("%s/%s", instanceId, consumerGroup))
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTrocketRocketmqConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = trocket.NewModifyConsumerGroupRequest()
		needChange bool
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	consumerGroup := idSplit[1]

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

		request.InstanceId = &instanceId
		request.ConsumerGroup = &consumerGroup
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().ModifyConsumerGroup(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update trocket rocketmq consumer group failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("trocket", "consumerGroup", tcClient.Region, fmt.Sprintf("%s/%s", instanceId, consumerGroup))
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTrocketRocketmqConsumerGroupRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_consumer_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
