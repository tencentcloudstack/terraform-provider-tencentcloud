package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfApplicationFileConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationFileConfigCreate,
		Read:   resourceTencentCloudTsfApplicationFileConfigRead,
		Delete: resourceTencentCloudTsfApplicationFileConfigDelete,

		Schema: map[string]*schema.Schema{
			"config_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config Name.",
			},

			"config_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config version.",
			},

			"config_file_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config file name.",
			},

			"config_file_value": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration file content (the original content encoding needs to be in utf-8 format, if the ConfigFileCode is gbk, it will be converted in the background).",
			},

			"application_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config file associated application ID.",
			},

			"config_file_path": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "config release path.",
			},

			"config_version_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "config version description.",
			},

			"config_file_code": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration file encoding, utf-8 or gbk. Note: If you choose gbk, you need the support of a new version of tsf-consul-template (public cloud virtual machines need to use 1.32 tsf-agent, and containers need to obtain the latest tsf-consul-template-docker.tar.gz from the documentation).",
			},

			"config_post_cmd": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "post command.",
			},

			"encode_with_base64": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "the config value is encoded with base64 or not.",
			},

			"program_id_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "datasource for auth.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationFileConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateFileConfigWithDetailRespRequest()
		response = tsf.NewCreateFileConfigWithDetailRespResponse()
		configId string
	)
	if v, ok := d.GetOk("config_name"); ok {
		request.ConfigName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_version"); ok {
		request.ConfigVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_file_name"); ok {
		request.ConfigFileName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_file_value"); ok {
		request.ConfigFileValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		request.ApplicationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_file_path"); ok {
		request.ConfigFilePath = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_version_desc"); ok {
		request.ConfigVersionDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_file_code"); ok {
		request.ConfigFileCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_post_cmd"); ok {
		request.ConfigPostCmd = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("encode_with_base64"); ok {
		request.EncodeWithBase64 = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateFileConfigWithDetailResp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationFileConfig failed, reason:%+v", logId, err)
		return err
	}

	configId = *response.Response.Result.ConfigId
	d.SetId(configId)

	return resourceTencentCloudTsfApplicationFileConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationFileConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	configId := d.Id()

	applicationFileConfig, err := service.DescribeTsfApplicationFileConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if applicationFileConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationFileConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationFileConfig.ConfigName != nil {
		_ = d.Set("config_name", applicationFileConfig.ConfigName)
	}

	if applicationFileConfig.ConfigVersion != nil {
		_ = d.Set("config_version", applicationFileConfig.ConfigVersion)
	}

	if applicationFileConfig.ConfigFileName != nil {
		_ = d.Set("config_file_name", applicationFileConfig.ConfigFileName)
	}

	if applicationFileConfig.ConfigFileValue != nil {
		_ = d.Set("config_file_value", applicationFileConfig.ConfigFileValue)
	}

	if applicationFileConfig.ApplicationId != nil {
		_ = d.Set("application_id", applicationFileConfig.ApplicationId)
	}

	if applicationFileConfig.ConfigFilePath != nil {
		_ = d.Set("config_file_path", applicationFileConfig.ConfigFilePath)
	}

	if applicationFileConfig.ConfigVersionDesc != nil {
		_ = d.Set("config_version_desc", applicationFileConfig.ConfigVersionDesc)
	}

	if applicationFileConfig.ConfigFileCode != nil {
		_ = d.Set("config_file_code", applicationFileConfig.ConfigFileCode)
	}

	if applicationFileConfig.ConfigPostCmd != nil {
		_ = d.Set("config_post_cmd", applicationFileConfig.ConfigPostCmd)
	}

	return nil
}

func resourceTencentCloudTsfApplicationFileConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	configId := d.Id()

	if err := service.DeleteTsfApplicationFileConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
