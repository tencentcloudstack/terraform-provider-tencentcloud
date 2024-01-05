package tpulsar

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqSubscriptionCreate,
		Read:   resourceTencentCloudTdmqSubscriptionRead,
		Delete: resourceTencentCloudTdmqSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Pulsar cluster ID.",
			},
			"environment_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Environment (namespace) name.",
			},
			"topic_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Topic name.",
			},
			"subscription_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subscriber name, which can contain up to 128 characters.",
			},
			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Remarks (up to 128 characters).",
			},
			"auto_create_policy_topic": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Default:     false,
				Description: "Whether to automatically create a dead letter topic and a retry letter topic. true: yes; false: no(default value).",
			},
			"auto_delete_policy_topic": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Default:     false,
				Description: "Whether to automatically delete a dead letter topic and a retry letter topic. Setting is only allowed when `auto_create_policy_topic` is true. Default is false.",
			},
		},
	}
}

func resourceTencentCloudTdmqSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_subscription.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                 = tccommon.GetLogId(tccommon.ContextNil)
		request               = tdmq.NewCreateSubscriptionRequest()
		clusterId             string
		environmentId         string
		topicName             string
		subscriptionName      string
		autoCreatePolicyTopic bool
		autoDeletePolicyTopic string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("environment_id"); ok {
		request.EnvironmentId = helper.String(v.(string))
		environmentId = v.(string)
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
		topicName = v.(string)
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		request.SubscriptionName = helper.String(v.(string))
		subscriptionName = v.(string)
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_create_policy_topic"); ok {
		request.AutoCreatePolicyTopic = helper.Bool(v.(bool))
		autoCreatePolicyTopic = v.(bool)

		if v, ok = d.GetOkExists("auto_delete_policy_topic"); ok {
			if !autoCreatePolicyTopic && v.(bool) {
				return errors.New("If `auto_create_policy_topic` is false, Can't set `auto_delete_policy_topic` param.")
			} else {
				if v.(bool) {
					autoDeletePolicyTopic = "true"
				} else {
					autoDeletePolicyTopic = "false"
				}
			}
		}
	}

	request.IsIdempotent = helper.Bool(false)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateSubscription(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || !*result.Response.Result {
			e = fmt.Errorf("create tdmq subscription failed")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq subscription failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{clusterId, environmentId, topicName, subscriptionName, autoDeletePolicyTopic}, tccommon.FILED_SP))

	return resourceTencentCloudTdmqSubscriptionRead(d, meta)
}

func resourceTencentCloudTdmqSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_subscription.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	environmentId := idSplit[1]
	topicName := idSplit[2]
	subscriptionName := idSplit[3]
	autoDeletePolicyTopic := idSplit[4]

	subscription, err := service.DescribeTdmqSubscriptionById(ctx, clusterId, environmentId, topicName, subscriptionName)
	if err != nil {
		return err
	}

	if subscription == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqSubscription` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("environment_id", subscription.EnvironmentId)
	_ = d.Set("topic_name", subscription.TopicName)
	_ = d.Set("subscription_name", subscriptionName)
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("remark", subscription.Remark)

	if autoDeletePolicyTopic == "true" {
		_ = d.Set("auto_delete_policy_topic", true)
	} else {
		_ = d.Set("auto_delete_policy_topic", false)
	}

	// Get Topics Status For auto_create_policy_topic
	has, err := service.GetTdmqTopicsAttachmentById(ctx, environmentId, topicName, subscriptionName, clusterId)
	if err != nil {
		return err
	}

	_ = d.Set("auto_create_policy_topic", has)
	return nil
}

func resourceTencentCloudTdmqSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_subscription.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	environmentId := idSplit[1]
	topicName := idSplit[2]
	subscriptionName := idSplit[3]
	autoDeletePolicyTopic := idSplit[4]

	if err := service.DeleteTdmqSubscriptionById(ctx, clusterId, environmentId, topicName, subscriptionName); err != nil {
		return err
	}

	if autoDeletePolicyTopic == "true" {
		// Delete Topics
		if err := service.DeleteTdmqTopicsAttachmentById(ctx, environmentId, topicName, subscriptionName, clusterId); err != nil {
			return err
		}
	}

	return nil
}
