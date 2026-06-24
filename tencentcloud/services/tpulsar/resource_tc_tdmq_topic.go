package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
)

func ResourceTencentCloudTdmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqTopicCreate,
		Read:   resourceTencentCloudTdmqTopicRead,
		Update: resourceTencentCloudTdmqTopicUpdate,
		Delete: resourceTencentCloudTdmqTopicDelete,

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
				ForceNew:    true,
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
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag description list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
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
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_topic.create")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tdmqService     = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
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

	var tags []*tdmq.Tag
	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagKey := dMap["tag_key"].(string)
			tagValue := dMap["tag_value"].(string)
			tags = append(tags, &tdmq.Tag{
				TagKey:   &tagKey,
				TagValue: &tagValue,
			})
		}
	}

	err := tdmqService.CreateTdmqTopic(ctx, environId, topicName, partitions, topicType, remark, clusterId, pulsarTopicType, tags)
	if err != nil {
		return err
	}

	d.SetId(topicName)
	return resourceTencentCloudTdmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tdmqService = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	topicName := d.Id()
	environId := d.Get("environ_id").(string)
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqTopicById(ctx, environId, topicName, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if !has {
			log.Printf("[WARN] tencentcloud_tdmq_topic id=%s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		if info.Partitions != nil {
			_ = d.Set("partitions", info.Partitions)
		}
		if info.TopicType != nil {
			_ = d.Set("topic_type", info.TopicType)
		}
		if info.PulsarTopicType != nil {
			_ = d.Set("pulsar_topic_type", info.PulsarTopicType)
		}
		if info.Remark != nil {
			_ = d.Set("remark", info.Remark)
		}
		if info.CreateTime != nil {
			_ = d.Set("create_time", info.CreateTime)
		}
		if info.Tags != nil {
			tagsList := []interface{}{}
			for _, t := range info.Tags {
				tagsMap := map[string]interface{}{}
				if t.TagKey != nil {
					tagsMap["tag_key"] = *t.TagKey
				}
				if t.TagValue != nil {
					tagsMap["tag_value"] = *t.TagValue
				}
				tagsList = append(tagsList, tagsMap)
			}
			_ = d.Set("tags", tagsList)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudTdmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_topic.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		partitions uint64
		remark     string
	)

	immutableArgs := []string{"topic_type", "pulsar_topic_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	topicName := d.Id()
	environId := d.Get("environ_id").(string)
	clusterId := d.Get("cluster_id").(string)

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

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		oldRaw, newRaw := d.GetChange("tags")
		oldTagsMap := tagsListToMap(oldRaw.([]interface{}))
		newTagsMap := tagsListToMap(newRaw.([]interface{}))
		replaceTags, deleteTags := svctag.DiffTags(oldTagsMap, newTagsMap)
		resourceName := tccommon.BuildTagResourceName("tdmq", "topic", tcClient.Region, fmt.Sprintf("%s/%s/%s", clusterId, environId, topicName))

		// Untag removed keys
		if len(deleteTags) > 0 {
			unTagRequest := tag.NewUnTagResourcesRequest()
			unTagRequest.ResourceList = []*string{&resourceName}
			tagKeys := make([]*string, 0, len(deleteTags))
			for _, key := range deleteTags {
				k := key
				tagKeys = append(tagKeys, &k)
			}
			unTagRequest.TagKeys = tagKeys

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				log.Printf("[DEBUG]%s api[UnTagResources] request: %s", logId, unTagRequest.ToJsonString())
				response, e := tcClient.UseTagClient().UnTagResources(unTagRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}
				log.Printf("[DEBUG]%s api[UnTagResources] response: %s", logId, response.ToJsonString())
				return nil
			})
			if err != nil {
				return err
			}
		}

		// Tag new/updated keys
		if len(replaceTags) > 0 {
			tagRequest := tag.NewTagResourcesRequest()
			tagRequest.ResourceList = []*string{&resourceName}
			tags := make([]*tag.Tag, 0, len(replaceTags))
			for k, v := range replaceTags {
				tagKey := k
				tagValue := v
				tags = append(tags, &tag.Tag{
					TagKey:   &tagKey,
					TagValue: &tagValue,
				})
			}
			tagRequest.Tags = tags

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				log.Printf("[DEBUG]%s api[TagResources] request: %s", logId, tagRequest.ToJsonString())
				response, e := tcClient.UseTagClient().TagResources(tagRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}
				log.Printf("[DEBUG]%s api[TagResources] response: %s", logId, response.ToJsonString())
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	d.Partial(false)
	return resourceTencentCloudTdmqTopicRead(d, meta)
}

func resourceTencentCloudTdmqTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_instance.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	topicName := d.Id()
	environId := d.Get("environ_id").(string)
	clusterId := d.Get("cluster_id").(string)

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

// tagsListToMap converts a tags list ([]interface{} with tag_key/tag_value) to map[string]interface{}
func tagsListToMap(tagsList []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, item := range tagsList {
		dMap := item.(map[string]interface{})
		if key, ok := dMap["tag_key"].(string); ok {
			if value, ok := dMap["tag_value"].(string); ok {
				result[key] = value
			}
		}
	}
	return result
}
