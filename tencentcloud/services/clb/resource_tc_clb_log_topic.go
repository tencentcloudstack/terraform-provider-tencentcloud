package clb

import (
	"context"
	"fmt"
	"log"
	"sync"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

var clsActionMu = &sync.Mutex{}

func ResourceTencentCloudClbLogTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceTopicCreate,
		Read:   resourceTencentCloudClbInstanceTopicRead,
		Update: resourceTencentCloudClbInstanceTopicUpdate,
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
			"status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The status of log topic. true: enable; false: disable. Default is true.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags of clb log topic.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			//compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Log topic creation time.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceTopicCreate(d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clsService := svccls.NewClsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

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
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	params := make(map[string]interface{})
	if topicName, ok := d.GetOk("topic_name"); ok {
		params["topic_name"] = topicName
	}
	if partitionCount, ok := d.GetOk("partition_count"); ok {
		params["partition_count"] = partitionCount
	}
	if tags, ok := d.GetOk("tags"); ok {
		params["tags"] = tags.([]interface{})
	}
	resp, err := clbService.CreateTopic(ctx, params)
	if err != nil {
		log.Printf("[CRITAL]%s create clb log topic failed, reason:%+v", logId, err)
		return err
	}
	if resp == nil || resp.Response == nil || resp.Response.TopicId == nil {
		log.Printf("[CRITAL]%s create clb log topic failed, response is nil, logId=%s", logId, logId)
		return fmt.Errorf("create clb log topic failed, response is nil")
	}

	topicId := *resp.Response.TopicId
	d.SetId(topicId)

	if v, ok := d.GetOkExists("status"); ok {
		if !v.(bool) {
			request := cls.NewModifyTopicRequest()
			request.TopicId = &topicId
			request.Status = helper.Bool(false)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyTopic(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudClbInstanceTopicRead(d, meta)
}

func resourceTencentCloudClbInstanceTopicRead(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	clsService := svccls.NewClsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	res, err := clsService.DescribeClsTopicById(ctx, id, nil)
	if err != nil {
		return err
	}
	if res == nil {
		log.Printf("[CRUD] clb log topic id=%s", id)
		d.SetId("")
		return fmt.Errorf("resource `logTopic` %s does not exist", id)
	}
	_ = d.Set("log_set_id", res.LogsetId)
	_ = d.Set("topic_name", res.TopicName)
	_ = d.Set("create_time", res.CreateTime)
	_ = d.Set("status", res.Status)

	if res.Tags != nil {
		tagsList := make([]map[string]interface{}, 0, len(res.Tags))
		for _, tag := range res.Tags {
			if tag == nil {
				continue
			}
			tagMap := make(map[string]interface{})
			if tag.Key != nil {
				tagMap["key"] = *tag.Key
			}
			if tag.Value != nil {
				tagMap["value"] = *tag.Value
			}
			tagsList = append(tagsList, tagMap)
		}
		_ = d.Set("tags", tagsList)
	}

	return nil
}

func resourceTencentCloudClbInstanceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_log_topic.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		topicId = d.Id()
	)

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			request := cls.NewModifyTopicRequest()
			request.TopicId = &topicId
			request.Status = helper.Bool(v.(bool))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyTopic(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("tags") {
		request := cls.NewModifyTopicRequest()
		request.TopicId = &topicId
		if v, ok := d.GetOk("tags"); ok {
			tagsList := v.([]interface{})
			if len(tagsList) > 0 {
				clsTags := make([]*cls.Tag, 0, len(tagsList))
				for _, tag := range tagsList {
					tagMap, ok := tag.(map[string]interface{})
					if !ok {
						continue
					}
					clsTag := &cls.Tag{}
					if key, ok := tagMap["key"].(string); ok {
						clsTag.Key = helper.String(key)
					}
					if value, ok := tagMap["value"].(string); ok {
						clsTag.Value = helper.String(value)
					}
					clsTags = append(clsTags, clsTag)
				}
				request.Tags = clsTags
			}
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyTopic(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return resourceTencentCloudClbInstanceTopicRead(d, meta)
}

func resourceTencentCloudClbInstanceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	clsActionMu.Lock()
	defer clsActionMu.Unlock()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	clsService := svccls.NewClsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := clsService.DeleteClsTopic(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
