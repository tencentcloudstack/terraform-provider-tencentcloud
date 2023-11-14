/*
Provides a resource to create a tsf application_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_config" "application_config" {
  config_name = "my_config"
  config_version = "1.0"
  config_value = "test: 1"
  application_id = "app-123456"
  config_version_desc = "product version"
  config_type = "A"
  encode_with_base64 = true
  program_id_list =
}
```

Import

tsf application_config can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_config.application_config application_config_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfApplicationConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationConfigCreate,
		Read:   resourceTencentCloudTsfApplicationConfigRead,
		Delete: resourceTencentCloudTsfApplicationConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"config_value": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config value, yaml or properties file.",
			},

			"application_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"config_version_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config version Description.",
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
				Description: "The config value is encoded with base64 or not.",
			},

			"program_id_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Datasource for auth.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateConfigWithDetailResp(request)
		if e != nil {
			return retryError(e)
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

	configId = *response.Response.ConfigId
	d.SetId(configId)

	return resourceTencentCloudTsfApplicationConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationConfigId := d.Id()

	applicationConfig, err := service.DescribeTsfApplicationConfigById(ctx, configId)
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

	if applicationConfig.EncodeWithBase64 != nil {
		_ = d.Set("encode_with_base64", applicationConfig.EncodeWithBase64)
	}

	if applicationConfig.ProgramIdList != nil {
		_ = d.Set("program_id_list", applicationConfig.ProgramIdList)
	}

	return nil
}

func resourceTencentCloudTsfApplicationConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	applicationConfigId := d.Id()

	if err := service.DeleteTsfApplicationConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
