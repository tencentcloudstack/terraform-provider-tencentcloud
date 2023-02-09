/*
Provides a resource to create a tsf application_file_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config" "application_file_config" {
  config_name = "terraform-test"
  config_version = "v1"
  config_file_name = "terraform-config-name"
  config_file_value = "ZWNobyAidGVzdCI="
  application_id = "application-ym9mxmza"
  config_file_path = "/etc/test/"
  config_version_desc = "terraform test version"
  config_file_code = "utf-8"
  config_post_cmd = "echo \"test1\""
  encode_with_base64 = true
  # program_id_list =
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf application_file_config can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_file_config.application_file_config application_file_config_id
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

func resourceTencentCloudTsfApplicationFileConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationFileConfigCreate,
		Read:   resourceTencentCloudTsfApplicationFileConfigRead,
		Update: resourceTencentCloudTsfApplicationFileConfigUpdate,
		Delete: resourceTencentCloudTsfApplicationFileConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config id.",
			},

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
				Description: "Configuration item file content (the original content encoding needs to be in utf-8 format, if the ConfigFileCode is gbk, it will be converted in the background).",
			},

			"application_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration item associated application ID.",
			},

			"config_file_path": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "release path.",
			},

			"config_version_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration item version description.",
			},

			"config_file_code": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration item file encoding, utf-8 or gbk. Note: If you choose gbk, you need the support of a new version of tsf-consul-template (public cloud virtual machines need to use 1.32 tsf-agent, and containers need to obtain the latest tsf-consul-template-docker.tar.gz from the documentation).",
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
				Description: "Base64-encoded configuration items.",
			},

			"creation_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
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
				Description: "Delete flag.",
			},

			"config_version_count": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of configuration item versions.",
			},

			"last_update_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Configuration item last update time.",
			},

			"config_file_value_length": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Configuration item file length.",
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

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationFileConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tsf.NewCreateFileConfigRequest()
		response   = tsf.NewCreateFileConfigResponse()
		configName string
	)
	if v, ok := d.GetOk("config_name"); ok {
		configName = v.(string)
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateFileConfig(request)
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

	if *response.Response.Result {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
		applicationFileConfig, err := service.DescribeTsfApplicationFileConfigById(ctx, "", configName)
		if err != nil {
			return err
		}

		if applicationFileConfig == nil || applicationFileConfig.ConfigId == nil {
			return fmt.Errorf("New file configuration failed [%v]", configName)
		}

		d.SetId(*applicationFileConfig.ConfigId)
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:fileConfigTag/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfApplicationFileConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationFileConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	configId := d.Id()

	applicationFileConfig, err := service.DescribeTsfApplicationFileConfigById(ctx, configId, "")
	if err != nil {
		return err
	}

	if applicationFileConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationFileConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationFileConfig.ConfigId != nil {
		_ = d.Set("config_id", applicationFileConfig.ConfigId)
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

	// if applicationFileConfig.EncodeWithBase64 != nil {
	// 	_ = d.Set("encode_with_base64", applicationFileConfig.EncodeWithBase64)
	// }

	if applicationFileConfig.CreationTime != nil {
		_ = d.Set("creation_time", applicationFileConfig.CreationTime)
	}

	if applicationFileConfig.ApplicationName != nil {
		_ = d.Set("application_name", applicationFileConfig.ApplicationName)
	}

	if applicationFileConfig.DeleteFlag != nil {
		_ = d.Set("delete_flag", applicationFileConfig.DeleteFlag)
	}

	if applicationFileConfig.ConfigVersionCount != nil {
		_ = d.Set("config_version_count", applicationFileConfig.ConfigVersionCount)
	}

	if applicationFileConfig.LastUpdateTime != nil {
		_ = d.Set("last_update_time", applicationFileConfig.LastUpdateTime)
	}

	if applicationFileConfig.ConfigFileValueLength != nil {
		_ = d.Set("config_file_value_length", applicationFileConfig.ConfigFileValueLength)
	}

	// if applicationFileConfig.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", applicationFileConfig.ProgramIdList)
	// }

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "fileConfigTag", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfApplicationFileConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "fileConfigTag", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfGroupRead(d, meta)
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
