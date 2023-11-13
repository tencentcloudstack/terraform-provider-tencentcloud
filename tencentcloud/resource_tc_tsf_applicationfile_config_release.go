/*
Provides a resource to create a tsf applicationfile_config_release

Example Usage

```hcl
resource "tencentcloud_tsf_applicationfile_config_release" "applicationfile_config_release" {
  config_id = "dcfg-f-123456"
  group_id = "group-123456"
  release_desc = "product release"
}
```

Import

tsf applicationfile_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_applicationfile_config_release.applicationfile_config_release applicationfile_config_release_id
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
	"strings"
)

func resourceTencentCloudTsfApplicationfileConfigRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationfileConfigReleaseCreate,
		Read:   resourceTencentCloudTsfApplicationfileConfigReleaseRead,
		Delete: resourceTencentCloudTsfApplicationfileConfigReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File config id.",
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
				Description: "Release Description.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationfileConfigReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_applicationfile_config_release.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewReleaseFileConfigRequest()
		response = tsf.NewReleaseFileConfigResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ReleaseFileConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationfileConfigRelease failed, reason:%+v", logId, err)
		return err
	}

	configId = *response.Response.ConfigId
	d.SetId(strings.Join([]string{configId, groupId}, FILED_SP))

	return resourceTencentCloudTsfApplicationfileConfigReleaseRead(d, meta)
}

func resourceTencentCloudTsfApplicationfileConfigReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_applicationfile_config_release.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]

	applicationfileConfigRelease, err := service.DescribeTsfApplicationfileConfigReleaseById(ctx, configId, groupId)
	if err != nil {
		return err
	}

	if applicationfileConfigRelease == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationfileConfigRelease` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationfileConfigRelease.ConfigId != nil {
		_ = d.Set("config_id", applicationfileConfigRelease.ConfigId)
	}

	if applicationfileConfigRelease.GroupId != nil {
		_ = d.Set("group_id", applicationfileConfigRelease.GroupId)
	}

	if applicationfileConfigRelease.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationfileConfigRelease.ReleaseDesc)
	}

	return nil
}

func resourceTencentCloudTsfApplicationfileConfigReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_applicationfile_config_release.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]

	if err := service.DeleteTsfApplicationfileConfigReleaseById(ctx, configId, groupId); err != nil {
		return err
	}

	return nil
}
