/*
Provides a resource to create a tsf application_file_config_release

Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config_release" "application_file_config_release" {
  config_id = "dcfg-f-123456"
  group_id = "group-123456"
  release_desc = "product release"
}
```

Import

tsf applicationfile_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_file_config_release.application_file_config_release application_file_config_release_id
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

func resourceTencentCloudTsfApplicationFileConfigRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationFileConfigReleaseCreate,
		Read:   resourceTencentCloudTsfApplicationFileConfigReleaseRead,
		Delete: resourceTencentCloudTsfApplicationFileConfigReleaseDelete,
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
				Description: "release Description.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationFileConfigReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config_release.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewReleaseFileConfigRequest()
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
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationFileConfigRelease failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(configId + FILED_SP + groupId)

	return resourceTencentCloudTsfApplicationFileConfigReleaseRead(d, meta)
}

func resourceTencentCloudTsfApplicationFileConfigReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config_release.read")()
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

	applicationFileConfigRelease, err := service.DescribeTsfApplicationFileConfigReleaseById(ctx, configId, groupId)
	if err != nil {
		return err
	}

	if applicationFileConfigRelease == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationFileConfigRelease` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationFileConfigRelease.ConfigId != nil {
		_ = d.Set("config_id", applicationFileConfigRelease.ConfigId)
	}

	if applicationFileConfigRelease.GroupId != nil {
		_ = d.Set("group_id", applicationFileConfigRelease.GroupId)
	}

	if applicationFileConfigRelease.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationFileConfigRelease.ReleaseDesc)
	}

	return nil
}

func resourceTencentCloudTsfApplicationFileConfigReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config_release.delete")()
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

	if err := service.DeleteTsfApplicationFileConfigReleaseById(ctx, configId, groupId); err != nil {
		return err
	}

	return nil
}
