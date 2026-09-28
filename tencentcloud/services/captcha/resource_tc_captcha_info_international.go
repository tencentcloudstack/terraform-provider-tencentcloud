package captcha

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	captchaintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCaptchaInfoInternational() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCaptchaInfoInternationalCreate,
		Read:   resourceTencentCloudCaptchaInfoInternationalRead,
		Update: resourceTencentCloudCaptchaInfoInternationalUpdate,
		Delete: resourceTencentCloudCaptchaInfoInternationalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Captcha name.",
			},

			"channel_info": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Client type. Valid values: `web` (web scenario), `android` (android client), `ios` (ios client). Default is `web`.",
			},

			"verify_rank": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Verify rank. Valid values: `1` (experience first), `2` (balanced), `3` (security first). Default is `1`.",
			},

			"user_set_cap_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Verify type. Valid values: `1` (invisible verify, DisableInvisibleSwitch must be `2`), `2` (slider), `8` (image click), `9` (voice).",
			},

			"defend_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defend mode. Valid values: `block` (block mode), `notify` (perceive mode). Default is `notify`.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Resource tags, key&value format.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"disable_invisible_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Verify mechanism. Valid values: `0` (one-click verify), `1` (always verify), `2` (invisible verify, UserSetCapType must be `1`).",
			},

			"verify_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Web domain. Only valid when `channel_info` is `web`.",
			},

			"verify_bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "App BundleId. Only valid when `channel_info` is `ios`.",
			},

			"verify_package": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "App package. Only valid when `channel_info` is `android`.",
			},

			"check_appid_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable captcha encryption. `0`: disable, `1`: enable.",
			},

			"check_iv_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable one-time pad. Valid values: `0` (disable), `1` (enable). Only allowed to be `1` when `check_appid_switch` is `1`.",
			},

			"check_box_style": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Checkbox style. Valid values: `0` (simple), `1` (basic), `2` (invisible).",
			},
		},
	}
}

func resourceTencentCloudCaptchaInfoInternationalCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_info_international.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = captchaintl.NewCreateCaptchaInfoInternationalRequest()
		response = captchaintl.NewCreateCaptchaInfoInternationalResponse()
	)

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("channel_info"); ok {
		request.ChannelInfo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verify_rank"); ok {
		request.VerifyRank = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_set_cap_type"); ok {
		request.UserSetCapType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("defend_mode"); ok {
		request.DefendMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsList := v.([]interface{})
		for i := range tagsList {
			if tag, ok := tagsList[i].(string); ok {
				request.Tags = append(request.Tags, helper.String(tag))
			}
		}
	}

	if v, ok := d.GetOk("disable_invisible_switch"); ok {
		request.DisableInvisibleSwitch = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verify_domain"); ok {
		request.VerifyDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verify_bundle_id"); ok {
		request.VerifyBundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verify_package"); ok {
		request.VerifyPackage = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("check_appid_switch"); ok {
		request.CheckAppidSwitch = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("check_iv_switch"); ok {
		request.CheckIvSwitch = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("check_box_style"); ok {
		request.CheckBoxStyle = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCaptchaClient().CreateCaptchaInfoInternationalWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create captcha info international failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create captcha info international failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data == nil {
		return fmt.Errorf("Data is nil.")
	}

	d.SetId(helper.Int64ToStr(*response.Response.Data))
	return resourceTencentCloudCaptchaInfoInternationalRead(d, meta)
}

func resourceTencentCloudCaptchaInfoInternationalRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_info_international.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = CaptchaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		captchaAppId = d.Id()
	)

	respData, err := service.DescribeCaptchaInfoInternationalById(ctx, captchaAppId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_captcha_info_international` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AppName != nil {
		_ = d.Set("app_name", respData.AppName)
	}

	if respData.ChannelInfo != nil {
		_ = d.Set("channel_info", respData.ChannelInfo)
	}

	if respData.VerifyRank != nil {
		_ = d.Set("verify_rank", helper.Int64ToStr(*respData.VerifyRank))
	}

	if respData.UserSetCapType != nil {
		_ = d.Set("user_set_cap_type", helper.Int64ToStr(*respData.UserSetCapType))
	}

	if respData.DefendMode != nil {
		_ = d.Set("defend_mode", respData.DefendMode)
	}

	if respData.Tags != nil {
		_ = d.Set("tags", respData.Tags)
	}

	if respData.DisableInvisibleSwitch != nil {
		_ = d.Set("disable_invisible_switch", respData.DisableInvisibleSwitch)
	}

	if respData.VerifyDomain != nil {
		_ = d.Set("verify_domain", respData.VerifyDomain)
	}

	if respData.VerifyBundleId != nil {
		_ = d.Set("verify_bundle_id", respData.VerifyBundleId)
	}

	if respData.VerifyPackage != nil {
		_ = d.Set("verify_package", respData.VerifyPackage)
	}

	if respData.CheckAppidSwitch != nil {
		_ = d.Set("check_appid_switch", respData.CheckAppidSwitch)
	}

	if respData.CheckIvSwitch != nil {
		_ = d.Set("check_iv_switch", respData.CheckIvSwitch)
	}

	if respData.CheckBoxStyle != nil {
		_ = d.Set("check_box_style", respData.CheckBoxStyle)
	}

	return nil
}

func resourceTencentCloudCaptchaInfoInternationalUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_info_international.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		captchaAppId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"app_name", "verify_rank", "user_set_cap_type", "defend_mode", "disable_invisible_switch", "verify_domain", "verify_bundle_id", "verify_package", "check_appid_switch", "check_iv_switch", "check_box_style", "tags"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := captchaintl.NewModifyCaptchaInfoInternationalRequest()
		request.CaptchaAppId = helper.String(captchaAppId)

		if v, ok := d.GetOk("app_name"); ok {
			request.AppName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("verify_rank"); ok {
			request.VerifyRank = helper.String(v.(string))
		}

		if v, ok := d.GetOk("user_set_cap_type"); ok {
			request.UserSetCapType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("defend_mode"); ok {
			request.DefendMode = helper.String(v.(string))
		}

		if v, ok := d.GetOk("disable_invisible_switch"); ok {
			request.DisableInvisibleSwitch = helper.String(v.(string))
		}

		if v, ok := d.GetOk("verify_domain"); ok {
			request.VerifyDomain = helper.String(v.(string))
		}

		if v, ok := d.GetOk("verify_bundle_id"); ok {
			request.VerifyBundleId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("verify_package"); ok {
			request.VerifyPackage = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("check_appid_switch"); ok {
			request.CheckAppidSwitch = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("check_iv_switch"); ok {
			request.CheckIvSwitch = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("check_box_style"); ok {
			request.CheckBoxStyle = helper.String(v.(string))
		}

		if v, ok := d.GetOk("tags"); ok {
			tagsList := v.([]interface{})
			for i := range tagsList {
				if tag, ok := tagsList[i].(string); ok {
					request.Tags = append(request.Tags, helper.String(tag))
				}
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCaptchaClient().ModifyCaptchaInfoInternationalWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify captcha info international failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update captcha info international failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudCaptchaInfoInternationalRead(d, meta)
}

func resourceTencentCloudCaptchaInfoInternationalDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_info_international.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = captchaintl.NewRemoveCaptchaInfoInternationalRequest()
		captchaAppId = d.Id()
	)

	request.CaptchaAppId = helper.String(captchaAppId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCaptchaClient().RemoveCaptchaInfoInternationalWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Remove captcha info international failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete captcha info international failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
