package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

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
				Description: "Client attestation option list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attestation option name.",
						},
						"attester_source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authentication method. Values: TC-RCE, TC-CAPTCHA, TC-EO-CAPTCHA.",
						},
						"attester_duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication validity duration. Default 60s. Supported units: s (60-43200), m (1-720), h (1-12).",
						},
						"tc_rce_option": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TC-RCE authentication configuration, required when attester_source is TC-RCE.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Channel information.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "RCE Channel region.",
									},
								},
							},
						},
						"tc_captcha_option": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TC-CAPTCHA authentication configuration, required when attester_source is TC-CAPTCHA.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"captcha_app_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "CaptchaAppId information.",
									},
									"app_secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "AppSecretKey information.",
									},
								},
							},
						},
						"tc_eo_captcha_option": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TC-EO-CAPTCHA authentication configuration, required when attester_source is TC-EO-CAPTCHA.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"captcha_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "EdgeOne human-machine verification mode. Values: Invisible, Adaptive.",
									},
								},
							},
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Attestation option ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule type. Values: PRESET, CUSTOM.",
						},
					},
				},
			},

			"client_attester_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of client attestation option IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityClientAttesterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		zoneId            string
		clientAttesterIds []string
	)
	var (
		request  = teov20220901.NewCreateSecurityClientAttesterRequest()
		response = teov20220901.NewCreateSecurityClientAttesterResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	request.ZoneId = helper.String(zoneId)

	if clientAttesters, ok := d.Get("client_attesters").([]interface{}); ok {
		for _, item := range clientAttesters {
			clientAttesterMap := item.(map[string]interface{})
			clientAttester := teov20220901.ClientAttester{}

			if v, ok := clientAttesterMap["name"]; ok {
				clientAttester.Name = helper.String(v.(string))
			}
			if v, ok := clientAttesterMap["attester_source"]; ok {
				clientAttester.AttesterSource = helper.String(v.(string))
			}
			if v, ok := clientAttesterMap["attester_duration"]; ok {
				clientAttester.AttesterDuration = helper.String(v.(string))
			}
			if tcRceOptions, ok := clientAttesterMap["tc_rce_option"].([]interface{}); ok && len(tcRceOptions) > 0 {
				tcRceOptionMap := tcRceOptions[0].(map[string]interface{})
				tcRceOption := teov20220901.TCRCEOption{}
				if v, ok := tcRceOptionMap["channel"]; ok {
					tcRceOption.Channel = helper.String(v.(string))
				}
				if v, ok := tcRceOptionMap["region"]; ok {
					tcRceOption.Region = helper.String(v.(string))
				}
				clientAttester.TCRCEOption = &tcRceOption
			}
			if tcCaptchaOptions, ok := clientAttesterMap["tc_captcha_option"].([]interface{}); ok && len(tcCaptchaOptions) > 0 {
				tcCaptchaOptionMap := tcCaptchaOptions[0].(map[string]interface{})
				tcCaptchaOption := teov20220901.TCCaptchaOption{}
				if v, ok := tcCaptchaOptionMap["captcha_app_id"]; ok {
					tcCaptchaOption.CaptchaAppId = helper.String(v.(string))
				}
				if v, ok := tcCaptchaOptionMap["app_secret_key"]; ok {
					tcCaptchaOption.AppSecretKey = helper.String(v.(string))
				}
				clientAttester.TCCaptchaOption = &tcCaptchaOption
			}
			if tcEoCaptchaOptions, ok := clientAttesterMap["tc_eo_captcha_option"].([]interface{}); ok && len(tcEoCaptchaOptions) > 0 {
				tcEoCaptchaOptionMap := tcEoCaptchaOptions[0].(map[string]interface{})
				tcEoCaptchaOption := teov20220901.TCEOCaptchaOption{}
				if v, ok := tcEoCaptchaOptionMap["captcha_mode"]; ok {
					tcEoCaptchaOption.CaptchaMode = helper.String(v.(string))
				}
				clientAttester.TCEOCaptchaOption = &tcEoCaptchaOption
			}

			request.ClientAttesters = append(request.ClientAttesters, &clientAttester)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityClientAttester(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo security client attester failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil && response.Response.ClientAttesterIds != nil {
		for _, id := range response.Response.ClientAttesterIds {
			clientAttesterIds = append(clientAttesterIds, *id)
		}
	}

	_ = d.Set("client_attester_ids", clientAttesterIds)

	d.SetId(strings.Join([]string{zoneId, strings.Join(clientAttesterIds, ",")}, tccommon.FILED_SP))

	return resourceTencentCloudTeoSecurityClientAttesterRead(d, meta)
}

func resourceTencentCloudTeoSecurityClientAttesterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	clientAttesterIdsStr := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoSecurityClientAttesterById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil || len(respData.ClientAttesters) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_security_client_attester` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// Filter by client_attester_ids from composite ID
	targetIds := strings.Split(clientAttesterIdsStr, ",")
	targetIdSet := make(map[string]bool)
	for _, id := range targetIds {
		targetIdSet[id] = true
	}

	clientAttestersList := make([]map[string]interface{}, 0)
	matchedIds := make([]string, 0)
	for _, clientAttester := range respData.ClientAttesters {
		if clientAttester.Id != nil && targetIdSet[*clientAttester.Id] {
			clientAttesterMap := map[string]interface{}{}

			if clientAttester.Id != nil {
				clientAttesterMap["id"] = *clientAttester.Id
				matchedIds = append(matchedIds, *clientAttester.Id)
			}

			if clientAttester.Name != nil {
				clientAttesterMap["name"] = *clientAttester.Name
			}

			if clientAttester.Type != nil {
				clientAttesterMap["type"] = *clientAttester.Type
			}

			if clientAttester.AttesterSource != nil {
				clientAttesterMap["attester_source"] = *clientAttester.AttesterSource
			}

			if clientAttester.AttesterDuration != nil {
				clientAttesterMap["attester_duration"] = *clientAttester.AttesterDuration
			}

			if clientAttester.TCRCEOption != nil {
				tcRceOptionList := make([]map[string]interface{}, 0, 1)
				tcRceOptionMap := map[string]interface{}{}
				if clientAttester.TCRCEOption.Channel != nil {
					tcRceOptionMap["channel"] = *clientAttester.TCRCEOption.Channel
				}
				if clientAttester.TCRCEOption.Region != nil {
					tcRceOptionMap["region"] = *clientAttester.TCRCEOption.Region
				}
				tcRceOptionList = append(tcRceOptionList, tcRceOptionMap)
				clientAttesterMap["tc_rce_option"] = tcRceOptionList
			} else {
				clientAttesterMap["tc_rce_option"] = []interface{}{}
			}

			if clientAttester.TCCaptchaOption != nil {
				tcCaptchaOptionList := make([]map[string]interface{}, 0, 1)
				tcCaptchaOptionMap := map[string]interface{}{}
				if clientAttester.TCCaptchaOption.CaptchaAppId != nil {
					tcCaptchaOptionMap["captcha_app_id"] = *clientAttester.TCCaptchaOption.CaptchaAppId
				}
				if clientAttester.TCCaptchaOption.AppSecretKey != nil {
					tcCaptchaOptionMap["app_secret_key"] = *clientAttester.TCCaptchaOption.AppSecretKey
				}
				tcCaptchaOptionList = append(tcCaptchaOptionList, tcCaptchaOptionMap)
				clientAttesterMap["tc_captcha_option"] = tcCaptchaOptionList
			} else {
				clientAttesterMap["tc_captcha_option"] = []interface{}{}
			}

			if clientAttester.TCEOCaptchaOption != nil {
				tcEoCaptchaOptionList := make([]map[string]interface{}, 0, 1)
				tcEoCaptchaOptionMap := map[string]interface{}{}
				if clientAttester.TCEOCaptchaOption.CaptchaMode != nil {
					tcEoCaptchaOptionMap["captcha_mode"] = *clientAttester.TCEOCaptchaOption.CaptchaMode
				}
				tcEoCaptchaOptionList = append(tcEoCaptchaOptionList, tcEoCaptchaOptionMap)
				clientAttesterMap["tc_eo_captcha_option"] = tcEoCaptchaOptionList
			} else {
				clientAttesterMap["tc_eo_captcha_option"] = []interface{}{}
			}

			clientAttestersList = append(clientAttestersList, clientAttesterMap)
		}
	}

	if len(clientAttestersList) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_security_client_attester` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("client_attesters", clientAttestersList)
	_ = d.Set("client_attester_ids", matchedIds)

	return nil
}

func resourceTencentCloudTeoSecurityClientAttesterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)

	needChange := false
	mutableArgs := []string{"client_attesters"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifySecurityClientAttesterRequest()

		request.ZoneId = helper.String(zoneId)

		if clientAttesters, ok := d.Get("client_attesters").([]interface{}); ok {
			for _, item := range clientAttesters {
				clientAttesterMap := item.(map[string]interface{})
				clientAttester := teov20220901.ClientAttester{}

				if v, ok := clientAttesterMap["id"]; ok && v.(string) != "" {
					clientAttester.Id = helper.String(v.(string))
				}
				if v, ok := clientAttesterMap["name"]; ok {
					clientAttester.Name = helper.String(v.(string))
				}
				if v, ok := clientAttesterMap["attester_source"]; ok {
					clientAttester.AttesterSource = helper.String(v.(string))
				}
				if v, ok := clientAttesterMap["attester_duration"]; ok {
					clientAttester.AttesterDuration = helper.String(v.(string))
				}
				if tcRceOptions, ok := clientAttesterMap["tc_rce_option"].([]interface{}); ok && len(tcRceOptions) > 0 {
					tcRceOptionMap := tcRceOptions[0].(map[string]interface{})
					tcRceOption := teov20220901.TCRCEOption{}
					if v, ok := tcRceOptionMap["channel"]; ok {
						tcRceOption.Channel = helper.String(v.(string))
					}
					if v, ok := tcRceOptionMap["region"]; ok {
						tcRceOption.Region = helper.String(v.(string))
					}
					clientAttester.TCRCEOption = &tcRceOption
				}
				if tcCaptchaOptions, ok := clientAttesterMap["tc_captcha_option"].([]interface{}); ok && len(tcCaptchaOptions) > 0 {
					tcCaptchaOptionMap := tcCaptchaOptions[0].(map[string]interface{})
					tcCaptchaOption := teov20220901.TCCaptchaOption{}
					if v, ok := tcCaptchaOptionMap["captcha_app_id"]; ok {
						tcCaptchaOption.CaptchaAppId = helper.String(v.(string))
					}
					if v, ok := tcCaptchaOptionMap["app_secret_key"]; ok {
						tcCaptchaOption.AppSecretKey = helper.String(v.(string))
					}
					clientAttester.TCCaptchaOption = &tcCaptchaOption
				}
				if tcEoCaptchaOptions, ok := clientAttesterMap["tc_eo_captcha_option"].([]interface{}); ok && len(tcEoCaptchaOptions) > 0 {
					tcEoCaptchaOptionMap := tcEoCaptchaOptions[0].(map[string]interface{})
					tcEoCaptchaOption := teov20220901.TCEOCaptchaOption{}
					if v, ok := tcEoCaptchaOptionMap["captcha_mode"]; ok {
						tcEoCaptchaOption.CaptchaMode = helper.String(v.(string))
					}
					clientAttester.TCEOCaptchaOption = &tcEoCaptchaOption
				}

				request.ClientAttesters = append(request.ClientAttesters, &clientAttester)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityClientAttester(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo security client attester failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoSecurityClientAttesterRead(d, meta)
}

func resourceTencentCloudTeoSecurityClientAttesterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_client_attester.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	clientAttesterIdsRaw := d.Get("client_attester_ids").([]interface{})

	var (
		request  = teov20220901.NewDeleteSecurityClientAttesterRequest()
		response = teov20220901.NewDeleteSecurityClientAttesterResponse()
	)

	request.ZoneId = helper.String(zoneId)

	for _, idRaw := range clientAttesterIdsRaw {
		request.ClientAttesterIds = append(request.ClientAttesterIds, helper.String(idRaw.(string)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityClientAttester(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo security client attester failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
