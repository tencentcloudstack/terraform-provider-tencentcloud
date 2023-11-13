/*
Provides a resource to create a tsf application_release_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id = ""
  group_id = ""
  release_desc = ""
}
```

Import

tsf application_release_config can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_release_config.application_release_config application_release_config_id
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

func resourceTencentCloudTsfApplicationReleaseConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationReleaseConfigCreate,
		Read:   resourceTencentCloudTsfApplicationReleaseConfigRead,
		Delete: resourceTencentCloudTsfApplicationReleaseConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config Id.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group Id.",
			},

			"release_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Release Desc.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationReleaseConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_release_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewReleaseConfigWithDetailRespRequest()
		response = tsf.NewReleaseConfigWithDetailRespResponse()
		configId string
	)
	if v, ok := d.GetOk("config_id"); ok {
		configId = v.(string)
		request.ConfigId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("release_desc"); ok {
		request.ReleaseDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ReleaseConfigWithDetailResp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationReleaseConfig failed, reason:%+v", logId, err)
		return err
	}

	configId = *response.Response.ConfigId
	d.SetId(configId)

	return resourceTencentCloudTsfApplicationReleaseConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationReleaseConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_release_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationReleaseConfigId := d.Id()

	applicationReleaseConfig, err := service.DescribeTsfApplicationReleaseConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if applicationReleaseConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationReleaseConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationReleaseConfig.ConfigId != nil {
		_ = d.Set("config_id", applicationReleaseConfig.ConfigId)
	}

	if applicationReleaseConfig.GroupId != nil {
		_ = d.Set("group_id", applicationReleaseConfig.GroupId)
	}

	if applicationReleaseConfig.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationReleaseConfig.ReleaseDesc)
	}

	return nil
}

func resourceTencentCloudTsfApplicationReleaseConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_release_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	applicationReleaseConfigId := d.Id()

	if err := service.DeleteTsfApplicationReleaseConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
