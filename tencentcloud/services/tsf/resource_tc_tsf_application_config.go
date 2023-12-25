package tsf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfApplicationConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationConfigCreate,
		Read:   resourceTencentCloudTsfApplicationConfigRead,
		Update: resourceTencentCloudTsfApplicationConfigUpdate,
		Delete: resourceTencentCloudTsfApplicationConfigDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"config_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "configuration item name.",
			},

			"config_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "configuration item version.",
			},

			"config_value": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "configuration item value.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"config_version_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "configuration item version description.",
			},

			"config_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "configuration item value type.",
			},

			"encode_with_base64": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Base64 encoded configuration items.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			// "tags": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "Tag description list.",
			// },
		},
	}
}

func resourceTencentCloudTsfApplicationConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_application_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = tsf.NewCreateConfigWithDetailRespRequest()
		response = tsf.NewCreateConfigWithDetailRespResponse()
		configId string
	)
	if v, ok := d.GetOk("config_name"); ok {
		request.ConfigName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_version"); ok {
		request.ConfigVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_value"); ok {
		request.ConfigValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		request.ApplicationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_version_desc"); ok {
		request.ConfigVersionDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_type"); ok {
		request.ConfigType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("encode_with_base64"); v != nil {
		request.EncodeWithBase64 = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().CreateConfigWithDetailResp(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationConfig failed, reason:%+v", logId, err)
		return err
	}

	configId = *response.Response.Result.ConfigId
	d.SetId(configId)

	return resourceTencentCloudTsfApplicationConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_application_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	configId := d.Id()

	applicationConfig, err := service.DescribeTsfApplicationConfigById(ctx, configId, "")
	if err != nil {
		return err
	}

	if applicationConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil

	}
	if applicationConfig.ConfigName != nil {
		_ = d.Set("config_name", applicationConfig.ConfigName)
	}

	if applicationConfig.ConfigVersion != nil {
		_ = d.Set("config_version", applicationConfig.ConfigVersion)
	}

	if applicationConfig.ConfigValue != nil {
		_ = d.Set("config_value", applicationConfig.ConfigValue)
	}

	if applicationConfig.ApplicationId != nil {
		_ = d.Set("application_id", applicationConfig.ApplicationId)
	}

	if applicationConfig.ConfigVersionDesc != nil {
		_ = d.Set("config_version_desc", applicationConfig.ConfigVersionDesc)
	}

	if applicationConfig.ConfigType != nil {
		_ = d.Set("config_type", applicationConfig.ConfigType)
	}

	// if applicationConfig.EncodeWithBase64 != nil {
	// 	_ = d.Set("encode_with_base64", applicationConfig.EncodeWithBase64)
	// }

	// if applicationConfig.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", applicationConfig.ProgramIdList)
	// }

	return nil
}

func resourceTencentCloudTsfApplicationConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_microservice.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{
		"config_name",
		"config_version",
		"config_value",
		"application_id",
		"config_version_desc",
		"config_type",
		"encode_with_base64",
		"program_id_list",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudTsfApplicationConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_application_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	configId := d.Id()

	if err := service.DeleteTsfApplicationConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
