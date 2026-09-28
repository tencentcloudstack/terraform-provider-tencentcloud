package captcha

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	captchaintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCaptchaIpWhiteListInternational() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCaptchaIpWhiteListInternationalCreate,
		Read:   resourceTencentCloudCaptchaIpWhiteListInternationalRead,
		Update: resourceTencentCloudCaptchaIpWhiteListInternationalUpdate,
		Delete: resourceTencentCloudCaptchaIpWhiteListInternationalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP allowlist name.",
			},

			"captcha_appid": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Captcha appid.",
			},

			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCaptchaIpWhiteListInternationalIp,
				Description:  "IP data. Must be a single valid IPv4 or IPv6 address.",
			},

			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark information.",
			},
		},
	}
}

func resourceTencentCloudCaptchaIpWhiteListInternationalCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_ip_white_list_international.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = captchaintl.NewCreateIpWhiteListInternationalRequest()
		response     = captchaintl.NewCreateIpWhiteListInternationalResponse()
		captchaAppid int64
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("captcha_appid"); ok {
		captchaAppid = int64(v.(int))
		request.CaptchaAppid = helper.Int64(captchaAppid)
	}

	if v, ok := d.GetOk("ip"); ok {
		request.Ip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCaptchaClient().CreateIpWhiteListInternationalWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ip white list international failed, IdList is empty."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create captcha ip white list international failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.IdList == nil || len(response.Response.IdList) == 0 {
		return fmt.Errorf("IdList is nil.")
	}

	d.SetId(fmt.Sprintf("%d#%d", captchaAppid, *response.Response.IdList[0]))
	return resourceTencentCloudCaptchaIpWhiteListInternationalRead(d, meta)
}

func resourceTencentCloudCaptchaIpWhiteListInternationalRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_ip_white_list_international.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CaptchaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	parts := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("resource id is broken, id: %s", d.Id())
	}

	captchaAppid := helper.StrToInt64(parts[0])
	id := helper.StrToInt64(parts[1])

	respData, err := service.DescribeCaptchaIpWhiteListInternationalById(ctx, captchaAppid, id)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_captcha_ip_white_list_international` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.CaptchaAppid != nil {
		_ = d.Set("captcha_appid", respData.CaptchaAppid)
	}

	if respData.Ip != nil {
		_ = d.Set("ip", respData.Ip)
	}

	if respData.Comment != nil {
		_ = d.Set("comment", respData.Comment)
	}

	return nil
}

func resourceTencentCloudCaptchaIpWhiteListInternationalUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_ip_white_list_international.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	needChange := false
	mutableArgs := []string{"name", "comment"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		parts := strings.Split(d.Id(), tccommon.FILED_SP)
		if len(parts) != 2 {
			return fmt.Errorf("resource id is broken, id: %s", d.Id())
		}

		request := captchaintl.NewModifyIpWhiteListInternationalRequest()
		request.CaptchaAppid = helper.Int64(helper.StrToInt64(parts[0]))
		request.Id = helper.Int64(helper.StrToInt64(parts[1]))

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCaptchaClient().ModifyIpWhiteListInternationalWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify ip white list international failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update captcha ip white list international failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudCaptchaIpWhiteListInternationalRead(d, meta)
}

func resourceTencentCloudCaptchaIpWhiteListInternationalDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_captcha_ip_white_list_international.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	parts := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("resource id is broken, id: %s", d.Id())
	}

	request := captchaintl.NewDeleteIpWhiteListInternationalRequest()
	request.CaptchaAppid = helper.Int64(helper.StrToInt64(parts[0]))
	request.Id = helper.Int64(helper.StrToInt64(parts[1]))

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCaptchaClient().DeleteIpWhiteListInternationalWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ip white list international failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete captcha ip white list international failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func validateCaptchaIpWhiteListInternationalIp(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if net.ParseIP(value) == nil {
		errors = append(errors, fmt.Errorf("%s must be a single valid IP address (IPv4 or IPv6), got: %s", k, value))
	}
	return warnings, errors
}
