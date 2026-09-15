package ioa

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ioav20220601 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIoaCompanyDirectoryConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIoaCompanyDirectoryConfigCreate,
		Read:   resourceTencentCloudIoaCompanyDirectoryConfigRead,
		Update: resourceTencentCloudIoaCompanyDirectoryConfigUpdate,
		Delete: resourceTencentCloudIoaCompanyDirectoryConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enterprise directory type, enum: `WeCom`, `Lark`, `DingTalk`, `MicrosoftEntraID`.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Enterprise directory name.",
			},

			"config": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Configuration data encrypted by SM2 then Hex encoded.",
			},

			"sync_enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether scheduled sync is enabled.",
			},

			"sync_policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Scheduled sync strategy, enum: `4hours`, `daily`, `weekly`.",
			},

			"sync_policy_params": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON string of sync strategy parameters.",
			},

			"create_auth_config": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to synchronously create an auth source.",
			},

			"display_on_login_page": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to display on the login page.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Directory description.",
			},

			"scene": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Usage scene: API created, quick start, normal config. Create-only, not present in Describe output.",
			},

			"name_i18n": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Multi-language names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lang": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Language enum, e.g. `zh-CN`, `en-US`.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Localized name value.",
						},
					},
				},
			},

			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Config source ID.",
			},

			"identify_source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity source config ID (from Create/Modify result).",
			},

			"auth_source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auth source config ID (from Create/Modify result).",
			},

			"auth_config_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Auth config ID (from Create/Modify result).",
			},

			"auth_policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Auth policy ID (from Create/Modify result).",
			},

			"auth_support_platforms": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Auth supported platforms, e.g. PC or Mobile (from Create/Modify result).",
			},

			"auth_methods": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Auth methods, e.g. authorization auth / scan auth (from Create/Modify result).",
			},
		},
	}
}

func resourceTencentCloudIoaCompanyDirectoryConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ioa_company_directory_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = ioav20220601.NewCreateCompanyDirectoryConfigRequest()
		response = ioav20220601.NewCreateCompanyDirectoryConfigResponse()
	)

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config"); ok {
		request.Config = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sync_enable"); ok {
		request.SyncEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("sync_policy"); ok {
		request.SyncPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sync_policy_params"); ok {
		request.SyncPolicyParams = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("create_auth_config"); ok {
		request.CreateAuthConfig = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("display_on_login_page"); ok {
		request.DisplayOnLoginPage = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scene"); ok {
		request.Scene = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name_i18n"); ok {
		for _, item := range v.([]interface{}) {
			nameI18nMap := item.(map[string]interface{})
			i18nString := ioav20220601.I18nString{}
			if v, ok := nameI18nMap["lang"].(string); ok && v != "" {
				i18nString.Lang = helper.String(v)
			}

			if v, ok := nameI18nMap["value"].(string); ok && v != "" {
				i18nString.Value = helper.String(v)
			}

			request.NameI18n = append(request.NameI18n, &i18nString)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIoaV20220601Client().CreateCompanyDirectoryConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ioa_company_directory_config failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ioa_company_directory_config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data == nil || response.Response.Data.Id == nil {
		log.Printf("[CRITAL]%s create ioa_company_directory_config response data or id is nil, d.Id()=%s", logId, d.Id())
		return fmt.Errorf("Create ioa_company_directory_config failed, response Data or Id is nil.")
	}

	d.SetId(strconv.FormatInt(*response.Response.Data.Id, 10))
	setDirectoryConfigResultData(d, response.Response.Data)
	return resourceTencentCloudIoaCompanyDirectoryConfigRead(d, meta)
}

func resourceTencentCloudIoaCompanyDirectoryConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ioa_company_directory_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IoaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCompanyDirectoryConfigById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] ioa_company_directory_config id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Config != nil {
		_ = d.Set("config", respData.Config)
	}

	if respData.SyncEnable != nil {
		_ = d.Set("sync_enable", respData.SyncEnable)
	}

	if respData.SyncPolicy != nil {
		_ = d.Set("sync_policy", respData.SyncPolicy)
	}

	if respData.SyncPolicyParams != nil {
		_ = d.Set("sync_policy_params", respData.SyncPolicyParams)
	}

	if respData.CreateAuthConfig != nil {
		_ = d.Set("create_auth_config", respData.CreateAuthConfig)
	}

	if respData.DisplayOnLoginPage != nil {
		_ = d.Set("display_on_login_page", respData.DisplayOnLoginPage)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.SourceId != nil {
		_ = d.Set("source_id", respData.SourceId)
	}

	if respData.NameI18n != nil && len(respData.NameI18n) > 0 {
		nameI18nList := make([]map[string]interface{}, 0, len(respData.NameI18n))
		for _, i18nString := range respData.NameI18n {
			nameI18nMap := map[string]interface{}{}
			if i18nString.Lang != nil {
				nameI18nMap["lang"] = i18nString.Lang
			}

			if i18nString.Value != nil {
				nameI18nMap["value"] = i18nString.Value
			}

			nameI18nList = append(nameI18nList, nameI18nMap)
		}

		_ = d.Set("name_i18n", nameI18nList)
	}

	return nil
}

func resourceTencentCloudIoaCompanyDirectoryConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ioa_company_directory_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	needChange := false
	mutableArgs := []string{"type", "name", "config", "sync_enable", "sync_policy", "sync_policy_params", "create_auth_config", "display_on_login_page", "description", "name_i18n"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := ioav20220601.NewModifyCompanyDirectoryConfigRequest()
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("config"); ok {
			request.Config = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sync_enable"); ok {
			request.SyncEnable = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("sync_policy"); ok {
			request.SyncPolicy = helper.String(v.(string))
		}

		if v, ok := d.GetOk("sync_policy_params"); ok {
			request.SyncPolicyParams = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("create_auth_config"); ok {
			request.CreateAuthConfig = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("display_on_login_page"); ok {
			request.DisplayOnLoginPage = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("name_i18n"); ok {
			for _, item := range v.([]interface{}) {
				nameI18nMap := item.(map[string]interface{})
				i18nString := ioav20220601.I18nString{}
				if v, ok := nameI18nMap["lang"].(string); ok && v != "" {
					i18nString.Lang = helper.String(v)
				}

				if v, ok := nameI18nMap["value"].(string); ok && v != "" {
					i18nString.Value = helper.String(v)
				}

				request.NameI18n = append(request.NameI18n, &i18nString)
			}
		}

		request.Id = helper.StrToInt64Point(d.Id())
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIoaV20220601Client().ModifyCompanyDirectoryConfigWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify ioa_company_directory_config failed, Response is nil."))
			}

			if result.Response.Data != nil {
				setDirectoryConfigResultData(d, result.Response.Data)
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update ioa_company_directory_config failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudIoaCompanyDirectoryConfigRead(d, meta)
}

func resourceTencentCloudIoaCompanyDirectoryConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ioa_company_directory_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ioav20220601.NewDeleteAccountGroupRequest()
	)

	accountGroupId := helper.StrToUint64Point(d.Id())
	request.AccountGroupId = accountGroupId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIoaV20220601Client().DeleteAccountGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ioa_company_directory_config failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ioa_company_directory_config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func setDirectoryConfigResultData(d *schema.ResourceData, data *ioav20220601.DirectoryConfigResultData) {
	if data == nil {
		return
	}

	if data.IdentifySourceId != nil {
		_ = d.Set("identify_source_id", data.IdentifySourceId)
	}

	if data.AuthSourceId != nil {
		_ = d.Set("auth_source_id", data.AuthSourceId)
	}

	if data.AuthConfigId != nil {
		_ = d.Set("auth_config_id", data.AuthConfigId)
	}

	if data.AuthPolicyId != nil {
		_ = d.Set("auth_policy_id", data.AuthPolicyId)
	}

	if data.AuthSupportPlatforms != nil {
		_ = d.Set("auth_support_platforms", helper.StringsInterfaces(data.AuthSupportPlatforms))
	}

	if data.AuthMethods != nil {
		_ = d.Set("auth_methods", helper.StringsInterfaces(data.AuthMethods))
	}
}
