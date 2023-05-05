/*
Provides a resource to create a tsf application_public_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config" "application_public_config" {
  config_name = "my_config"
  config_version = "1.0"
  config_value = "test: 1"
  config_version_desc = "product version"
  config_type = "P"
  encode_with_base64 = true
  # program_id_list =
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfApplicationPublicConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationPublicConfigCreate,
		Read:   resourceTencentCloudTsfApplicationPublicConfigRead,
		Delete: resourceTencentCloudTsfApplicationPublicConfigDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
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
				Description: "config version.",
			},

			"config_value": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "config value, only yaml file allowed.",
			},

			"config_version_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config version description.",
			},

			"config_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config type.",
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

func resourceTencentCloudTsfApplicationPublicConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreatePublicConfigWithDetailRespRequest()
		response = tsf.NewCreatePublicConfigWithDetailRespResponse()
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

	if v, ok := d.GetOk("config_version_desc"); ok {
		request.ConfigVersionDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_type"); ok {
		request.ConfigType = helper.String(v.(string))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreatePublicConfigWithDetailResp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationPublicConfig failed, reason:%+v", logId, err)
		return err
	}

	configId = *response.Response.Result.ConfigId
	d.SetId(configId)

	return resourceTencentCloudTsfApplicationPublicConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationPublicConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	configId := d.Id()

	applicationPublicConfig, err := service.DescribeTsfApplicationPublicConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if applicationPublicConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationPublicConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationPublicConfig.ConfigName != nil {
		_ = d.Set("config_name", applicationPublicConfig.ConfigName)
	}

	if applicationPublicConfig.ConfigVersion != nil {
		_ = d.Set("config_version", applicationPublicConfig.ConfigVersion)
	}

	if applicationPublicConfig.ConfigValue != nil {
		_ = d.Set("config_value", applicationPublicConfig.ConfigValue)
	}

	if applicationPublicConfig.ConfigVersionDesc != nil {
		_ = d.Set("config_version_desc", applicationPublicConfig.ConfigVersionDesc)
	}

	if applicationPublicConfig.ConfigType != nil {
		_ = d.Set("config_type", applicationPublicConfig.ConfigType)
	}

	// if applicationPublicConfig.EncodeWithBase64 != nil {
	// 	_ = d.Set("encode_with_base64", applicationPublicConfig.EncodeWithBase64)
	// }

	// if applicationPublicConfig.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", applicationPublicConfig.ProgramIdList)
	// }

	return nil
}

func resourceTencentCloudTsfApplicationPublicConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	configId := d.Id()

	if err := service.DeleteTsfApplicationPublicConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
