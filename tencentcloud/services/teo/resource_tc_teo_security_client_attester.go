package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoSecurityClientAttester() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityClientAttesterCreate,
		Read:   resourceTencentCloudTeoSecurityClientAttesterRead,
		Update: resourceTencentCloudTeoSecurityClientAttesterUpdate,
		Delete: resourceTencentCloudTeoSecurityClientAttesterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"client_attesters": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Client attester configuration. Only one attester is allowed per request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attester name.",
						},
						"attester_source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authentication method. Valid values: `TC-RCE` (Tencent Cloud RCE), `TC-CAPTCHA` (Tencent CAPTCHA), `TC-EO-CAPTCHA` (EdgeOne CAPTCHA).",
						},
						"attester_duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication validity duration. Default `60s`. Supported units: `s` (60-43200), `m` (1-720), `h` (1-12). e.g. `300s`, `5m`, `1h`.",
						},
						"tc_rce_option": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TC-RCE authentication configuration. Required when `attester_source` is `TC-RCE`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "TC-RCE channel ID.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TC-RCE channel region. Valid values: `ap-beijing`, `ap-jakarta`, `ap-singapore`, `eu-frankfurt`, `na-siliconvalley`.",
									},
								},
							},
						},
						"tc_captcha_option": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TC-CAPTCHA authentication configuration. Required when `attester_source` is `TC-CAPTCHA`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"captcha_app_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "CaptchaAppId.",
									},
									"app_secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
										Description: "AppSecretKey.",
									},
								},
							},
						},
						"tc_eo_captcha_option": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TC-EO-CAPTCHA authentication configuration. Required when `attester_source` is `TC-EO-CAPTCHA`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"captcha_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "EdgeOne CAPTCHA mode. Valid values: `Invisible`, `Adaptive`.",
									},
								},
							},
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client attester ID, e.g. `attest-xxxxxxxx`.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Attester rule type (read-only). Valid values: `PRESET` (system preset), `CUSTOM` (user defined).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityClientAttesterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request          = teo.NewCreateSecurityClientAttesterRequest()
		response         = teo.NewCreateSecurityClientAttesterResponse()
		zoneId           string
		clientAttesterId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("client_attesters"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		request.ClientAttesters = []*teo.ClientAttester{buildClientAttesterFromMap(m, "")}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityClientAttesterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo security client attester failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo security client attester failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.ClientAttesterIds) == 0 || response.Response.ClientAttesterIds[0] == nil {
		return fmt.Errorf("ClientAttesterIds is empty.")
	}

	clientAttesterId = *response.Response.ClientAttesterIds[0]
	d.SetId(strings.Join([]string{zoneId, clientAttesterId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoSecurityClientAttesterRead(d, meta)
}

func resourceTencentCloudTeoSecurityClientAttesterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	clientAttesterId := idSplit[1]

	respData, err := service.DescribeTeoSecurityClientAttesterById(ctx, zoneId, clientAttesterId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_security_client_attester` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	attesterMap := map[string]interface{}{}

	if respData.Id != nil {
		attesterMap["id"] = *respData.Id
	}

	if respData.Type != nil {
		attesterMap["type"] = *respData.Type
	}

	if respData.Name != nil {
		attesterMap["name"] = *respData.Name
	}

	if respData.AttesterSource != nil {
		attesterMap["attester_source"] = *respData.AttesterSource
	}

	if respData.AttesterDuration != nil {
		attesterMap["attester_duration"] = *respData.AttesterDuration
	}

	if respData.TCRCEOption != nil {
		opt := map[string]interface{}{}
		if respData.TCRCEOption.Channel != nil {
			opt["channel"] = *respData.TCRCEOption.Channel
		}
		if respData.TCRCEOption.Region != nil {
			opt["region"] = *respData.TCRCEOption.Region
		}
		attesterMap["tc_rce_option"] = []interface{}{opt}
	}

	if respData.TCCaptchaOption != nil {
		opt := map[string]interface{}{}
		if respData.TCCaptchaOption.CaptchaAppId != nil {
			opt["captcha_app_id"] = *respData.TCCaptchaOption.CaptchaAppId
		}
		if respData.TCCaptchaOption.AppSecretKey != nil {
			opt["app_secret_key"] = *respData.TCCaptchaOption.AppSecretKey
		}
		attesterMap["tc_captcha_option"] = []interface{}{opt}
	}

	if respData.TCEOCaptchaOption != nil {
		opt := map[string]interface{}{}
		if respData.TCEOCaptchaOption.CaptchaMode != nil {
			opt["captcha_mode"] = *respData.TCEOCaptchaOption.CaptchaMode
		}
		attesterMap["tc_eo_captcha_option"] = []interface{}{opt}
	}

	_ = d.Set("client_attesters", []interface{}{attesterMap})

	return nil
}

func resourceTencentCloudTeoSecurityClientAttesterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewModifySecurityClientAttesterRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	clientAttesterId := idSplit[1]

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("client_attesters"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		request.ClientAttesters = []*teo.ClientAttester{buildClientAttesterFromMap(m, clientAttesterId)}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityClientAttesterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo security client attester failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoSecurityClientAttesterRead(d, meta)
}

func resourceTencentCloudTeoSecurityClientAttesterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteSecurityClientAttesterRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.ClientAttesterIds = []*string{helper.String(idSplit[1])}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityClientAttesterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo security client attester failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildClientAttesterFromMap converts a schema map block to *teo.ClientAttester.
// Pass id="" when creating (API does not accept Id on create).
func buildClientAttesterFromMap(m map[string]interface{}, id string) *teo.ClientAttester {
	attester := &teo.ClientAttester{}

	if id != "" {
		attester.Id = helper.String(id)
	}

	if val, ok := m["name"].(string); ok && val != "" {
		attester.Name = helper.String(val)
	}

	if val, ok := m["attester_source"].(string); ok && val != "" {
		attester.AttesterSource = helper.String(val)
	}

	if val, ok := m["attester_duration"].(string); ok && val != "" {
		attester.AttesterDuration = helper.String(val)
	}

	if val, ok := m["tc_rce_option"].([]interface{}); ok && len(val) > 0 {
		opt := val[0].(map[string]interface{})
		rce := &teo.TCRCEOption{}
		if v, ok := opt["channel"].(string); ok && v != "" {
			rce.Channel = helper.String(v)
		}
		if v, ok := opt["region"].(string); ok && v != "" {
			rce.Region = helper.String(v)
		}
		attester.TCRCEOption = rce
	}

	if val, ok := m["tc_captcha_option"].([]interface{}); ok && len(val) > 0 {
		opt := val[0].(map[string]interface{})
		cap := &teo.TCCaptchaOption{}
		if v, ok := opt["captcha_app_id"].(string); ok && v != "" {
			cap.CaptchaAppId = helper.String(v)
		}
		if v, ok := opt["app_secret_key"].(string); ok && v != "" {
			cap.AppSecretKey = helper.String(v)
		}
		attester.TCCaptchaOption = cap
	}

	if val, ok := m["tc_eo_captcha_option"].([]interface{}); ok && len(val) > 0 {
		opt := val[0].(map[string]interface{})
		eoOpt := &teo.TCEOCaptchaOption{}
		if v, ok := opt["captcha_mode"].(string); ok && v != "" {
			eoOpt.CaptchaMode = helper.String(v)
		}
		attester.TCEOCaptchaOption = eoOpt
	}

	return attester
}
