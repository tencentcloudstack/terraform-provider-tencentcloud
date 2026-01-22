package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoContentIdentifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoContentIdentifierCreate,
		Read:   resourceTencentCloudTeoContentIdentifierRead,
		Update: resourceTencentCloudTeoContentIdentifierUpdate,
		Delete: resourceTencentCloudTeoContentIdentifierDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the content identifier, length limit of up to 20 characters.",
			},

			"plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target plan id to be bound, available only for the enterprise edition. <li>if there is already a plan under your account, go to [plan management](https://console.cloud.tencent.com/edgeone/package) to get the plan id and directly bind the content identifier to the plan;</li><li>if you do not have a plan to bind, please purchase an enterprise edition plan first.</li>.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags of the content identifier. this parameter is used for authority control. to create tags, go to the [tag console](https://console.cloud.tencent.com/tag/taglist).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag key.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag value.\nNote: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			// computed
			"content_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content identifier ID.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time, which is in Coordinated Universal Time (UTC) and follows the ISO 8601 date and time format..",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of the latest update, in Coordinated Universal Time (UTC), following the ISO 8601 date and time format..",
			},
		},
	}
}

func resourceTencentCloudTeoContentIdentifierCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreateContentIdentifierRequest()
		response = teov20220901.NewCreateContentIdentifierResponse()
	)

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plan_id"); ok {
		request.PlanId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			tag := teov20220901.Tag{}
			if v, ok := tagsMap["tag_key"].(string); ok && v != "" {
				tag.TagKey = helper.String(v)
			}

			if v, ok := tagsMap["tag_value"].(string); ok && v != "" {
				tag.TagValue = helper.String(v)
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateContentIdentifierWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo content identifier failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo content identifier failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ContentId == nil {
		return fmt.Errorf("ContentId is nil.")
	}

	d.SetId(*response.Response.ContentId)
	return resourceTencentCloudTeoContentIdentifierRead(d, meta)
}

func resourceTencentCloudTeoContentIdentifierRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		contentId = d.Id()
	)

	respData, err := service.DescribeTeoContentIdentifierById(ctx, contentId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_content_identifier` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.PlanId != nil {
		_ = d.Set("plan_id", respData.PlanId)
	}

	if respData.Tags != nil {
		tagsList := make([]map[string]interface{}, 0, len(respData.Tags))
		for _, tags := range respData.Tags {
			tagsMap := map[string]interface{}{}
			if tags.TagKey != nil {
				tagsMap["tag_key"] = tags.TagKey
			}

			if tags.TagValue != nil {
				tagsMap["tag_value"] = tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)
	}

	if respData.ContentId != nil {
		_ = d.Set("content_id", respData.ContentId)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoContentIdentifierUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		contentId = d.Id()
	)

	immutableArgs := []string{"plan_id", "tags"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		request := teov20220901.NewModifyContentIdentifierRequest()
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		request.ContentId = &contentId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyContentIdentifierWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo content identifier failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoContentIdentifierRead(d, meta)
}

func resourceTencentCloudTeoContentIdentifierDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teov20220901.NewDeleteContentIdentifierRequest()
		contentId = d.Id()
	)

	request.ContentId = &contentId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteContentIdentifierWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo content identifier failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
