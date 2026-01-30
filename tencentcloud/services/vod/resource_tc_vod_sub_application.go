package vod

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodSubApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSubApplicationCreate,
		Read:   resourceTencentCloudVodSubApplicationRead,
		Update: resourceTencentCloudVodSubApplicationUpdate,
		Delete: resourceTencentCloudVodSubApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 40),
				Description:  "Sub application name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.",
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Sub appliaction status.",
				ValidateFunc: tccommon.ValidateAllowedStringValue(VOD_SUB_APPLICATION_STATUS),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sub application description.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tag key-value pairs for resource management. Maximum 10 tags.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the sub application was created.",
			},
		},
	}
}

func resourceTencentCloudVodSubApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sub_application.create")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewCreateSubAppIdRequest()
		subAppId   *uint64
		subAppName *string
	)

	if v, ok := d.GetOk("name"); ok {
		subAppName = helper.String(v.(string))
		request.Name = subAppName
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		for key, value := range tags {
			tag := vod.ResourceTag{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateSubAppId(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		subAppId = response.Response.SubAppId
		return nil
	}); err != nil {
		return err
	}

	d.SetId(*subAppName + tccommon.FILED_SP + helper.UInt64ToStr(*subAppId))

	if v, ok := d.GetOk("status"); ok {
		statusResquest := vod.NewModifySubAppIdStatusRequest()
		statusResquest.SubAppId = subAppId
		statusResquest.Status = helper.String(v.(string))

		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(statusResquest.GetAction())
			_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySubAppIdStatus(statusResquest)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
						return resource.RetryableError(err)
					}
				}
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentCloudVodSubApplicationRead(d, meta)
}

func resourceTencentCloudVodSubApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sub_application.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		//logId   = tccommon.GetLogId(tccommon.ContextNil)
		//ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client  = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		request = vod.NewDescribeSubAppIdsRequest()
		appInfo = vod.SubAppIdInfo{}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sub application id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[1]

	// Paginated query with retry logic
	var offset uint64 = 0
	var limit uint64 = 200
	found := false

	for !found {
		request.Offset = &offset
		request.Limit = &limit

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, err := client.UseVodClient().DescribeSubAppIds(request)
			if err != nil {
				return tccommon.RetryError(err)
			}

			infoSet := response.Response.SubAppIdInfoSet
			if len(infoSet) == 0 && offset == 0 {
				// No sub applications found at all
				d.SetId("")
				found = true
				return nil
			}

			// Search for the target sub application in current page
			for _, info := range infoSet {
				if helper.UInt64ToStr(helper.PUint64(info.SubAppId)) == subAppId {
					appInfo = *info
					found = true
					return nil
				}
			}

			// Check if we need to continue pagination
			totalCount := response.Response.TotalCount
			if totalCount == nil || offset+uint64(len(infoSet)) >= *totalCount {
				// Reached the end without finding the sub application
				d.SetId("")
				found = true
			} else {
				// Continue to next page
				offset += limit
			}

			return nil
		})

		if err != nil {
			return err
		}

		// If ID was set to empty (not found), return early
		if d.Id() == "" {
			return nil
		}
	}
	_ = d.Set("name", appInfo.Name)
	_ = d.Set("description", appInfo.Description)
	// there may hide a bug, appInfo do not return status. So use the user input
	//_ = d.Set("status", appInfo.Status)
	_ = d.Set("status", helper.String(d.Get("status").(string)))
	_ = d.Set("create_time", appInfo.CreateTime)

	// Set tags if returned by API
	if appInfo.Tags != nil {
		tags := make(map[string]string)
		for _, tag := range appInfo.Tags {
			if tag.TagKey != nil && tag.TagValue != nil {
				tags[*tag.TagKey] = *tag.TagValue
			}
		}
		_ = d.Set("tags", tags)
	}

	// Note: DescribeSubAppIds API does not return Type field
	// Type is only set during creation and cannot be queried via API
	// It is preserved in Terraform state (ForceNew field doesn't need to be updated)

	return nil
}

func resourceTencentCloudVodSubApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sub_application.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = vod.NewModifySubAppIdInfoRequest()
		changeFlag = false
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sub application id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[1]

	if d.HasChange("name") {
		changeFlag = true
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}
	if d.HasChange("description") {
		changeFlag = true
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}
	request.SubAppId = helper.Uint64(helper.StrToUInt64(subAppId))

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySubAppIdInfo(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if d.HasChange("status") {
		var statusRequest = vod.NewModifySubAppIdStatusRequest()
		if v, ok := d.GetOk("status"); ok {
			statusRequest.Status = helper.String(v.(string))
		}
		statusRequest.SubAppId = helper.Uint64(helper.StrToUInt64(subAppId))
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(statusRequest.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySubAppIdStatus(statusRequest)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
						return resource.RetryableError(err)
					}
				}
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, statusRequest.GetAction(), err.Error())
				return resource.NonRetryableError(err)
			}
			return nil
		})
	}

	// Handle tags update using unified tag service
	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		region := tcClient.Region
		resourceName := fmt.Sprintf("qcs::vod:%s:uin/:subAppId/%s", region, subAppId)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudVodSubApplicationRead(d, meta)
}

func resourceTencentCloudVodSubApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sub_application.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sub application id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[1]
	var statusRequest = vod.NewModifySubAppIdStatusRequest()
	// first turn off
	statusRequest.Status = helper.String("Off")
	statusRequest.SubAppId = helper.Uint64(helper.StrToUInt64(subAppId))
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(statusRequest.GetAction())
		if _, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySubAppIdStatus(statusRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, statusRequest.GetAction(), err.Error())
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	// then destroy
	statusRequest.Status = helper.String("Destroyed")
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(statusRequest.GetAction())
		if _, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySubAppIdStatus(statusRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, statusRequest.GetAction(), err.Error())
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
