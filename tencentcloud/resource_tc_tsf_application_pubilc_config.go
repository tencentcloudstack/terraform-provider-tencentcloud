/*
Provides a resource to create a tsf application_pubilc_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_pubilc_config" "application_pubilc_config" {
  config_name = ""
  config_version = ""
  config_value = ""
  config_version_desc = ""
  config_type = ""
  encode_with_base64 =
  program_id_list =
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf application_pubilc_config can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_pubilc_config.application_pubilc_config application_pubilc_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfApplicationPubilcConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationPubilcConfigCreate,
		Read:   resourceTencentCloudTsfApplicationPubilcConfigRead,
		Delete: resourceTencentCloudTsfApplicationPubilcConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config name.",
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
				Description: "Configuration item value, always receive the content in yaml format.",
			},

			"config_version_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration item version description.",
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
				Description: "Base64-encoded configuration items.",
			},

			"program_id_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"config_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config id.",
			},

			"creation_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"application_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Application id.",
			},

			"application_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Application name.",
			},

			"delete_flag": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Delete flag, true: can be deleted; false: cannot be deleted.",
			},

			"last_update_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "last updated.",
			},

			"config_version_count": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of configuration item versions.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationPubilcConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_pubilc_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tsf.NewCreatePublicConfigRequest()
		// response = tsf.NewCreatePublicConfigResponse()
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreatePublicConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationPubilcConfig failed, reason:%+v", logId, err)
		return err
	}

	// configId = *response.Response.Result
	d.SetId(configId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:publicConfigTag/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfApplicationPubilcConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationPubilcConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_pubilc_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	configId := d.Id()

	applicationPubilcConfig, err := service.DescribeTsfApplicationPubilcConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if applicationPubilcConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationPubilcConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationPubilcConfig.ConfigName != nil {
		_ = d.Set("config_name", applicationPubilcConfig.ConfigName)
	}

	if applicationPubilcConfig.ConfigVersion != nil {
		_ = d.Set("config_version", applicationPubilcConfig.ConfigVersion)
	}

	if applicationPubilcConfig.ConfigValue != nil {
		_ = d.Set("config_value", applicationPubilcConfig.ConfigValue)
	}

	if applicationPubilcConfig.ConfigVersionDesc != nil {
		_ = d.Set("config_version_desc", applicationPubilcConfig.ConfigVersionDesc)
	}

	if applicationPubilcConfig.ConfigType != nil {
		_ = d.Set("config_type", applicationPubilcConfig.ConfigType)
	}

	// if applicationPubilcConfig.EncodeWithBase64 != nil {
	// 	_ = d.Set("encode_with_base64", applicationPubilcConfig.EncodeWithBase64)
	// }

	// if applicationPubilcConfig.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", applicationPubilcConfig.ProgramIdList)
	// }

	if applicationPubilcConfig.ConfigId != nil {
		_ = d.Set("config_id", applicationPubilcConfig.ConfigId)
	}

	if applicationPubilcConfig.CreationTime != nil {
		_ = d.Set("creation_time", applicationPubilcConfig.CreationTime)
	}

	if applicationPubilcConfig.ApplicationId != nil {
		_ = d.Set("application_id", applicationPubilcConfig.ApplicationId)
	}

	if applicationPubilcConfig.ApplicationName != nil {
		_ = d.Set("application_name", applicationPubilcConfig.ApplicationName)
	}

	if applicationPubilcConfig.DeleteFlag != nil {
		_ = d.Set("delete_flag", applicationPubilcConfig.DeleteFlag)
	}

	if applicationPubilcConfig.LastUpdateTime != nil {
		_ = d.Set("last_update_time", applicationPubilcConfig.LastUpdateTime)
	}

	if applicationPubilcConfig.ConfigVersionCount != nil {
		_ = d.Set("config_version_count", applicationPubilcConfig.ConfigVersionCount)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "publicConfigTag", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfApplicationPubilcConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_pubilc_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	configId := d.Id()

	if err := service.DeleteTsfApplicationPubilcConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
