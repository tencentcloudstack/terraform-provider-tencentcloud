package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
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
			"usage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current usage amount. Unit: when `quota_type` is `Image`, count by images; when `quota_type` is `Video`, count by seconds; when `quota_type` is `Text`, count by tokens.",
			},
		},
	}
}

func ParseVodAigcQuotaId(id string) (subAppId uint64, quotaType string, apiToken string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		err = fmt.Errorf("invalid resource id %q, expected format `sub_app_id%squota_type%sapi_token`", id, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	subAppId, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid resource id %q: sub_app_id part %q is not a valid uint64: %s", id, parts[0], err.Error())
		return
	}
	quotaType = parts[1]
	apiToken = parts[2]
	return
}

func resourceTencentCloudVodAigcQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	subAppId := uint64(d.Get("sub_app_id").(int))
	quotaType := d.Get("quota_type").(string)
	quotaLimit := uint64(d.Get("quota_limit").(int))
	apiToken := ""
	if v, ok := d.GetOk("api_token"); ok {
		apiToken = v.(string)
	}

	if err := service.CreateVodAigcQuota(ctx, subAppId, quotaType, quotaLimit, apiToken); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d%s%s%s%s", subAppId, tccommon.FILED_SP, quotaType, tccommon.FILED_SP, apiToken))

	// Poll until the quota becomes visible via DescribeAigcQuotas (cloud side has sync delay).
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		item, e := service.DescribeVodAigcQuotaById(ctx, subAppId, quotaType, apiToken)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if item == nil {
			return resource.RetryableError(fmt.Errorf("newly created vod_aigc_quota is not yet visible for sub_app_id [%d], quota_type [%s]; retrying", subAppId, quotaType))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudVodAigcQuotaRead(d, meta)
}

func resourceTencentCloudVodAigcQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	subAppId, quotaType, apiToken, err := ParseVodAigcQuotaId(d.Id())
	if err != nil {
		return err
	}

	item, err := service.DescribeVodAigcQuotaById(ctx, subAppId, quotaType, apiToken)
	if err != nil {
		return err
	}
	if item == nil {
		log.Printf("[CRUD] vod_aigc_quota id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("sub_app_id", int(subAppId))
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	subAppId, quotaType, apiToken, err := ParseVodAigcQuotaId(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("quota_limit") {
		quotaLimit := uint64(d.Get("quota_limit").(int))
		if err := service.ModifyVodAigcQuota(ctx, subAppId, quotaType, quotaLimit, apiToken); err != nil {
			return err
		}
	}

	return resourceTencentCloudVodAigcQuotaRead(d, meta)
}

func resourceTencentCloudVodAigcQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_aigc_quota.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	subAppId, quotaType, apiToken, err := ParseVodAigcQuotaId(d.Id())
	if err != nil {
		return err
	}

	if err := service.DeleteVodAigcQuota(ctx, subAppId, quotaType, apiToken); err != nil {
		return err
	}

	// Poll until the quota disappears from the list (cloud side has sync delay).
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		item, e := service.DescribeVodAigcQuotaById(ctx, subAppId, quotaType, apiToken)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if item != nil {
			return resource.RetryableError(fmt.Errorf("vod_aigc_quota still exists for sub_app_id [%d], quota_type [%s]; waiting for sync", subAppId, quotaType))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
