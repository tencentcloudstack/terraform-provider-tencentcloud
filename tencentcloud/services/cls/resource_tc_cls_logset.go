package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsLogset() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudClsLogsetRead,
		Create: resourceTencentCloudClsLogsetCreate,
		Update: resourceTencentCloudClsLogsetUpdate,
		Delete: resourceTencentCloudClsLogsetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset name, which must be unique.",
			},

			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				Deprecated:    "It is recommended to use `tag_list` because the current `tags` field is binding resources by calling the tag API.",
				ConflictsWith: []string{"tag_list"},
				Description:   "Tag description list.",
			},

			"tag_list": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"tags"},
				Description:   "Tag description list. The CLS API supports up to 10 tag key-value pairs, and duplicate keys are not allowed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"topic_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of log topics in logset.",
			},

			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If assumer_uin is not empty, it indicates the service provider who creates the logset.",
			},
		},
	}
}

func resourceTencentCloudClsLogsetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request  = cls.NewCreateLogsetRequest()
		response = cls.NewCreateLogsetResponse()
	)

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_list"); ok {
		tagsList := v.([]interface{})
		for _, item := range tagsList {
			tagMap := item.(map[string]interface{})
			tag := &cls.Tag{}
			if v, ok := tagMap["key"].(string); ok && v != "" {
				tag.Key = helper.String(v)
			}
			if v, ok := tagMap["value"].(string); ok {
				tag.Value = helper.String(v)
			}
			request.Tags = append(request.Tags, tag)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateLogset(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls logset failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls logset failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.LogsetId == nil {
		return fmt.Errorf("LogsetId is nil.")
	}

	logsetId := *response.Response.LogsetId
	d.SetId(logsetId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::cls:%s:uin/:logset/%s", region, logsetId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsLogsetRead(d, meta)
}

func resourceTencentCloudClsLogsetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		logsetId = d.Id()
	)

	logset, err := service.DescribeClsLogset(ctx, logsetId)
	if err != nil {
		return err
	}

	if logset == nil {
		log.Printf("resource `logset` %s does not exist", logsetId)
		d.SetId("")
		return nil
	}

	if logset.LogsetName != nil {
		_ = d.Set("logset_name", logset.LogsetName)
	}

	if logset.CreateTime != nil {
		_ = d.Set("create_time", logset.CreateTime)
	}

	if logset.TopicCount != nil {
		_ = d.Set("topic_count", logset.TopicCount)
	}

	if logset.RoleName != nil {
		_ = d.Set("role_name", logset.RoleName)
	}

	if logset.Tags == nil {
		tagsList := make([]map[string]interface{}, 0, len(logset.Tags))
		for _, tag := range logset.Tags {
			if tag == nil {
				continue
			}
			tagMap := map[string]interface{}{}
			if tag.Key != nil {
				tagMap["key"] = tag.Key
			}
			if tag.Value != nil {
				tagMap["value"] = tag.Value
			}
			tagsList = append(tagsList, tagMap)
		}
		_ = d.Set("tag_list", tagsList)
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		tags, err := tagService.DescribeResourceTags(ctx, "cls", "logset", tcClient.Region, d.Id())
		if err != nil {
			return err
		}

		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudClsLogsetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	if d.HasChange("logset_name") || d.HasChange("tag_list") {
		request := cls.NewModifyLogsetRequest()
		request.LogsetId = helper.String(d.Id())
		if v, ok := d.GetOk("logset_name"); ok {
			request.LogsetName = helper.String(v.(string))
		}

		if d.HasChange("tag_list") {
			if v, ok := d.GetOk("tag_list"); ok {
				tagsList := v.([]interface{})
				for _, item := range tagsList {
					tagMap := item.(map[string]interface{})
					tag := &cls.Tag{}
					if v, ok := tagMap["key"].(string); ok && v != "" {
						tag.Key = helper.String(v)
					}
					if v, ok := tagMap["value"].(string); ok {
						tag.Value = helper.String(v)
					}
					request.Tags = append(request.Tags, tag)
				}
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyLogset(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 && d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("cls", "logset", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsLogsetRead(d, meta)
}

func resourceTencentCloudClsLogsetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		logsetId = d.Id()
	)

	if err := service.DeleteClsLogsetById(ctx, logsetId); err != nil {
		return err
	}

	return nil
}
