/*
Provides a resource to create a tsf application_release_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id = "dcfg-nalqbqwv"
  group_id = "group-yxmz72gv"
  release_desc = "terraform-test"
}
```

Import

tsf application_release_config can be imported using the configId#groupId#configReleaseId, e.g.

```
terraform import tencentcloud_tsf_application_release_config.application_release_config dcfg-nalqbqwv#group-yxmz72gv#dcfgr-maeeq2ea
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "Configuration ID.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "deployment group ID.",
			},

			"release_desc": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "release description.",
			},

			"config_release_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "configuration item release ID.",
			},

			"config_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "configuration item name.",
			},

			"config_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "configuration item version.",
			},

			"release_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "release time.",
			},

			"group_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "deployment group name.",
			},

			"namespace_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID.",
			},

			"namespace_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},

			"cluster_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "cluster ID.",
			},

			"cluster_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "cluster name.",
			},

			"application_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
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
		groupId  string
	)
	if v, ok := d.GetOk("config_id"); ok {
		configId = v.(string)
		request.ConfigId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
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

	configReleaseId := *response.Response.Result.ConfigReleaseId
	d.SetId(configId + FILED_SP + groupId + FILED_SP + configReleaseId)

	return resourceTencentCloudTsfApplicationReleaseConfigRead(d, meta)
}

func resourceTencentCloudTsfApplicationReleaseConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_release_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]

	applicationReleaseConfig, err := service.DescribeTsfApplicationReleaseConfigById(ctx, configId, groupId)
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

	if applicationReleaseConfig.ConfigReleaseId != nil {
		_ = d.Set("config_release_id", applicationReleaseConfig.ConfigReleaseId)
	}

	if applicationReleaseConfig.ConfigName != nil {
		_ = d.Set("config_name", applicationReleaseConfig.ConfigName)
	}

	if applicationReleaseConfig.ConfigVersion != nil {
		_ = d.Set("config_version", applicationReleaseConfig.ConfigVersion)
	}

	if applicationReleaseConfig.ReleaseTime != nil {
		_ = d.Set("release_time", applicationReleaseConfig.ReleaseTime)
	}

	if applicationReleaseConfig.GroupName != nil {
		_ = d.Set("group_name", applicationReleaseConfig.GroupName)
	}

	if applicationReleaseConfig.NamespaceId != nil {
		_ = d.Set("namespace_id", applicationReleaseConfig.NamespaceId)
	}

	if applicationReleaseConfig.NamespaceName != nil {
		_ = d.Set("namespace_name", applicationReleaseConfig.NamespaceName)
	}

	if applicationReleaseConfig.ClusterId != nil {
		_ = d.Set("cluster_id", applicationReleaseConfig.ClusterId)
	}

	if applicationReleaseConfig.ClusterName != nil {
		_ = d.Set("cluster_name", applicationReleaseConfig.ClusterName)
	}

	if applicationReleaseConfig.ApplicationId != nil {
		_ = d.Set("application_id", applicationReleaseConfig.ApplicationId)
	}

	return nil
}

func resourceTencentCloudTsfApplicationReleaseConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_release_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[2]

	if err := service.DeleteTsfApplicationReleaseConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
