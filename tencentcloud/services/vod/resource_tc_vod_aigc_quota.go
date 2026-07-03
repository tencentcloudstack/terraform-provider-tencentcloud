package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodAigcQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodAigcQuotaCreate,
		Read:   resourceTencentCloudVodAigcQuotaRead,
		Update: resourceTencentCloudVodAigcQuotaUpdate,
		Delete: resourceTencentCloudVodAigcQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sub_app_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The VOD [sub application](https://intl.cloud.tencent.com/document/product/266/14574) ID. Users who activated VOD service after December 25, 2023 must specify the application ID.",
			},

			"quota_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quota type. Valid values: `Image` (AIGC image generation), `Video` (AIGC video generation), `Text` (AIGC text generation).",
			},

			"quota_limit": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Quota limit value. Unit: when `quota_type` is `Image`, count by images; when `quota_type` is `Video`, count by seconds; when `quota_type` is `Text`, count by tokens.",
			},

			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "API token for quota restriction. Only meaningful when `quota_type` is `Text`.",
			},

			// computed
			"usage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current usage amount. Unit: when `quota_type` is `Image`, count by images; when `quota_type` is `Video`, count by seconds; when `quota_type` is `Text`, count by tokens.",
			},
		},
	}
}

func resourceTencentCloudVodAigcQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = vod.NewCreateAigcQuotaRequest()
	)

	var (
		subAppId  string
		quotaType string
		apiToken  string
	)

	if v, ok := d.GetOk("sub_app_id"); ok {
		request.SubAppId = helper.IntUint64(v.(int))
		subAppId = helper.IntToStr(v.(int))
	}

	if v, ok := d.GetOk("quota_type"); ok {
		request.QuotaType = helper.String(v.(string))
		quotaType = v.(string)
	}

	if v, ok := d.GetOk("quota_limit"); ok {
		request.QuotaLimit = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("api_token"); ok {
		request.ApiToken = helper.String(v.(string))
		apiToken = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateAigcQuotaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vod aigc quota failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vod aigc quota failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if apiToken != "" {
		d.SetId(strings.Join([]string{subAppId, quotaType, apiToken}, tccommon.FILED_SP))
	} else {
		d.SetId(strings.Join([]string{subAppId, quotaType}, tccommon.FILED_SP))
	}

	return resourceTencentCloudVodAigcQuotaRead(d, meta)
}

func resourceTencentCloudVodAigcQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	subAppId, quotaType, apiToken, err := ParseVodAigcQuotaId(d.Id())
	if err != nil {
		return err
	}

	item, err := service.DescribeVodAigcQuotaById(ctx, subAppId, quotaType, apiToken)
	if err != nil {
		return err
	}

	if item == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vod_aigc_quota` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("sub_app_id", subAppId)
	_ = d.Set("quota_type", quotaType)

	if apiToken != "" {
		_ = d.Set("api_token", apiToken)
	}

	if item.QuotaLimit != nil {
		_ = d.Set("quota_limit", int(*item.QuotaLimit))
	}

	if item.Usage != nil {
		_ = d.Set("usage", int(*item.Usage))
	}

	return nil
}

func resourceTencentCloudVodAigcQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	subAppId, quotaType, apiToken, err := ParseVodAigcQuotaId(d.Id())
	if err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{"quota_limit"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := vod.NewModifyAigcQuotaRequest()
		request.SubAppId = &subAppId
		request.QuotaType = helper.String(quotaType)
		if v, ok := d.GetOk("quota_limit"); ok {
			request.QuotaLimit = helper.IntUint64(v.(int))
		}

		if apiToken != "" {
			request.ApiToken = helper.String(apiToken)
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyAigcQuotaWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify vod aigc quota failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update vod aigc quota failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudVodAigcQuotaRead(d, meta)
}

func resourceTencentCloudVodAigcQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	subAppId, quotaType, apiToken, err := ParseVodAigcQuotaId(d.Id())
	if err != nil {
		return err
	}

	request := vod.NewDeleteAigcQuotaRequest()
	request.SubAppId = &subAppId
	request.QuotaType = helper.String(quotaType)
	if apiToken != "" {
		request.ApiToken = helper.String(apiToken)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().DeleteAigcQuotaWithContext(ctx, request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.HasPrefix(sdkErr.Code, "ResourceNotFound") {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete vod aigc quota failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vod aigc quota failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// ParseVodAigcQuotaId parses the resource id into its constituent parts. The
// resource id has two possible formats depending on whether an api_token is
// present:
//   - sub_app_id#quota_type              (Image / Video, no api_token)
//   - sub_app_id#quota_type#api_token    (Text, with api_token)
func ParseVodAigcQuotaId(id string) (subAppId uint64, quotaType string, apiToken string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	switch len(parts) {
	case 2:
		if parts[0] == "" || parts[1] == "" {
			err = fmt.Errorf("invalid resource id %q, expected format `sub_app_id%squota_type` or `sub_app_id%squota_type%sapi_token`", id, tccommon.FILED_SP, tccommon.FILED_SP, tccommon.FILED_SP)
			return
		}
	case 3:
		if parts[0] == "" || parts[1] == "" {
			err = fmt.Errorf("invalid resource id %q, expected format `sub_app_id%squota_type%sapi_token`", id, tccommon.FILED_SP, tccommon.FILED_SP)
			return
		}
		apiToken = parts[2]
	default:
		err = fmt.Errorf("invalid resource id %q, expected format `sub_app_id%squota_type` or `sub_app_id%squota_type%sapi_token`", id, tccommon.FILED_SP, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	subAppId, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid resource id %q: sub_app_id part %q is not a valid uint64: %s", id, parts[0], err.Error())
		return
	}
	quotaType = parts[1]
	return
}
