/*
Provides a resource to create a rum offline_log_config

Example Usage

```hcl
resource "tencentcloud_rum_offline_log_config" "offline_log_config" {
  project_key = &lt;nil&gt;
  unique_i_d = &lt;nil&gt;
  }
```

Import

rum offline_log_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_offline_log_config.offline_log_config offline_log_config_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudRumOfflineLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumOfflineLogConfigCreate,
		Read:   resourceTencentCloudRumOfflineLogConfigRead,
		Delete: resourceTencentCloudRumOfflineLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_key": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unique project key for reporting.",
			},

			"unique_i_d": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of the user to be listened on(aid or uin).",
			},

			"msg": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "API call information.",
			},
		},
	}
}

func resourceTencentCloudRumOfflineLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_offline_log_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = rum.NewCreateOfflineLogConfigRequest()
		response   = rum.NewCreateOfflineLogConfigResponse()
		projectKey string
	)
	if v, ok := d.GetOk("project_key"); ok {
		projectKey = v.(string)
		request.ProjectKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unique_i_d"); ok {
		request.UniqueID = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateOfflineLogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum offlineLogConfig failed, reason:%+v", logId, err)
		return err
	}

	projectKey = *response.Response.ProjectKey
	d.SetId(projectKey)

	return resourceTencentCloudRumOfflineLogConfigRead(d, meta)
}

func resourceTencentCloudRumOfflineLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_offline_log_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	offlineLogConfigId := d.Id()

	offlineLogConfig, err := service.DescribeRumOfflineLogConfigById(ctx, projectKey)
	if err != nil {
		return err
	}

	if offlineLogConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumOfflineLogConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if offlineLogConfig.ProjectKey != nil {
		_ = d.Set("project_key", offlineLogConfig.ProjectKey)
	}

	if offlineLogConfig.UniqueID != nil {
		_ = d.Set("unique_i_d", offlineLogConfig.UniqueID)
	}

	if offlineLogConfig.Msg != nil {
		_ = d.Set("msg", offlineLogConfig.Msg)
	}

	return nil
}

func resourceTencentCloudRumOfflineLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_offline_log_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	offlineLogConfigId := d.Id()

	if err := service.DeleteRumOfflineLogConfigById(ctx, projectKey); err != nil {
		return err
	}

	return nil
}
