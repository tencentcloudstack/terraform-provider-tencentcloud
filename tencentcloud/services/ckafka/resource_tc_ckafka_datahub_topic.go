package ckafka

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaDatahubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaDatahubTopicCreate,
		Read:   resourceTencentCloudCkafkaDatahubTopicRead,
		Update: resourceTencentCloudCkafkaDatahubTopicUpdate,
		Delete: resourceTencentCloudCkafkaDatahubTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name, start with appid, which is a string of no more than 128 characters, must start with a letter, and the rest can contain letters, numbers, and dashes (-).",
			},

			"partition_num": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of Partitions, greater than 0.",
			},

			"retention_ms": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Message retention time, in ms, the current minimum value is 60000 ms.",
			},

			"note": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subject note, which is a string of no more than 64 characters, must start with a letter, and the rest can contain letters, numbers and dashes (-).",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of dataHub topic.",
			},
		},
	}
}

func resourceTencentCloudCkafkaDatahubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_topic.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = ckafka.NewCreateDatahubTopicRequest()
		topicName string
	)
	if v, ok := d.GetOk("name"); ok {
		topicName = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("partition_num"); ok {
		request.PartitionNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("retention_ms"); ok {
		request.RetentionMs = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("note"); ok {
		request.Note = helper.String(v.(string))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tagInfo := ckafka.Tag{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tagInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreateDatahubTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka datahubTopic failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(topicName)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::ckafka:%s:uin/:dipTopic/%s", region, topicName)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCkafkaDatahubTopicRead(d, meta)
}

func resourceTencentCloudCkafkaDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	topicName := d.Id()

	datahubTopic, err := service.DescribeCkafkaDatahubTopicById(ctx, topicName)
	if err != nil {
		return err
	}

	if datahubTopic == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaDatahubTopic` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if datahubTopic.Name != nil {
		_ = d.Set("name", datahubTopic.Name)
	}

	if datahubTopic.PartitionNum != nil {
		_ = d.Set("partition_num", datahubTopic.PartitionNum)
	}

	if datahubTopic.RetentionMs != nil {
		_ = d.Set("retention_ms", datahubTopic.RetentionMs)
	}

	if datahubTopic.Note != nil {
		_ = d.Set("note", datahubTopic.Note)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "ckafka", "dipTopic", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCkafkaDatahubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_topic.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := ckafka.NewModifyDatahubTopicRequest()

	topicName := d.Id()

	request.Name = &topicName

	immutableArgs := []string{"partition_num"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"retention_ms", "note"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOkExists("retention_ms"); ok {
			request.RetentionMs = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("note"); ok {
			request.Note = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().ModifyDatahubTopic(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ckafka datahubTopic failed, reason:%+v", logId, err)
			return err
		}
	}

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("ckafka", "dipTopic", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudCkafkaDatahubTopicRead(d, meta)
}

func resourceTencentCloudCkafkaDatahubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_datahub_topic.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	topicName := d.Id()

	if err := service.DeleteCkafkaDatahubTopicById(ctx, topicName); err != nil {
		return err
	}

	return nil
}
