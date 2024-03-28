package vod

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodEventConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodEventConfigCreate,
		Read:   resourceTencentCloudVodEventConfigRead,
		Update: resourceTencentCloudVodEventConfigUpdate,
		Delete: resourceTencentCloudVodEventConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sub_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sub app id.",
			},

			"mode": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
				Description: "How to receive event notifications. Valid values:\n" +
					"- Push: HTTP callback notification;\n" +
					"- PULL: Reliable notification based on message queuing.",
			},

			"notification_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The address used to receive 3.0 format callbacks when receiving HTTP callback notifications. Note: If you take the NotificationUrl parameter and the value is an empty string, the 3.0 format callback address is cleared.",
			},

			"upload_media_complete_event_switch": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether to receive video upload completion event notification, default `OFF` means to ignore the event notification, `ON` means to receive event notification.",
			},

			"delete_media_complete_event_switch": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether to receive video deletion completion event notification, default `OFF` is to ignore the event notification, `ON` is to receive event notification.",
			},
		},
	}
}

func resourceTencentCloudVodEventConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_event_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var subAppId int
	if v, ok := d.GetOk("sub_app_id"); ok {
		subAppId = v.(int)
	}

	d.SetId(strconv.Itoa(subAppId))

	return resourceTencentCloudVodEventConfigUpdate(d, meta)
}

func resourceTencentCloudVodEventConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_event_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	subAppId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("sub_app_id", subAppId)

	eventConfig, err := service.DescribeVodEventConfig(ctx, uint64(subAppId))
	if err != nil {
		return err
	}

	if eventConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodEventConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if eventConfig.Mode != nil {
		_ = d.Set("mode", eventConfig.Mode)
	}

	if eventConfig.NotificationUrl != nil {
		_ = d.Set("notification_url", eventConfig.NotificationUrl)
	}

	if eventConfig.UploadMediaCompleteEventSwitch != nil {
		_ = d.Set("upload_media_complete_event_switch", eventConfig.UploadMediaCompleteEventSwitch)
	}

	if eventConfig.DeleteMediaCompleteEventSwitch != nil {
		_ = d.Set("delete_media_complete_event_switch", eventConfig.DeleteMediaCompleteEventSwitch)
	}

	return nil
}

func resourceTencentCloudVodEventConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_event_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifyEventConfigRequest()

	subAppId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	request.SubAppId = helper.IntUint64(subAppId)

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notification_url"); ok {
		request.NotificationUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upload_media_complete_event_switch"); ok {
		request.UploadMediaCompleteEventSwitch = helper.String(v.(string))
	}

	if v, ok := d.GetOk("delete_media_complete_event_switch"); ok {
		request.DeleteMediaCompleteEventSwitch = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyEventConfig(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod eventConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodEventConfigRead(d, meta)
}

func resourceTencentCloudVodEventConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_event_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
