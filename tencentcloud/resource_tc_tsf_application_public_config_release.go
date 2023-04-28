/*
Provides a resource to create a tsf application_public_config_release

Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config_release" "application_public_config_release" {
  config_id = "dcfg-p-123456"
  namespace_id = "namespace-123456"
  release_desc = "product version"
}
```

Import

tsf application_public_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_public_config_release.application_public_config_release application_public_config_attachment_id
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

func resourceTencentCloudTsfApplicationPublicConfigRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationPublicConfigReleaseCreate,
		Read:   resourceTencentCloudTsfApplicationPublicConfigReleaseRead,
		Delete: resourceTencentCloudTsfApplicationPublicConfigReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ConfigId.",
			},

			"namespace_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "namespace-id.",
			},

			"release_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Release description.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationPublicConfigReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_release.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tsf.NewReleasePublicConfigRequest()
		configId    string
		namespaceId string
	)
	if v, ok := d.GetOk("config_id"); ok {
		configId = v.(string)
		request.ConfigId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("release_desc"); ok {
		request.ReleaseDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ReleasePublicConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationPublicConfigRelease failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(configId + FILED_SP + namespaceId)

	return resourceTencentCloudTsfApplicationPublicConfigReleaseRead(d, meta)
}

func resourceTencentCloudTsfApplicationPublicConfigReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_release.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	namespaceId := idSplit[1]

	applicationPublicConfigRelease, err := service.DescribeTsfApplicationPublicConfigReleaseById(ctx, configId, namespaceId)
	if err != nil {
		return err
	}

	if applicationPublicConfigRelease == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationPublicConfigRelease` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationPublicConfigRelease.ConfigId != nil {
		_ = d.Set("config_id", applicationPublicConfigRelease.ConfigId)
	}

	if applicationPublicConfigRelease.NamespaceId != nil {
		_ = d.Set("namespace_id", applicationPublicConfigRelease.NamespaceId)
	}

	if applicationPublicConfigRelease.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationPublicConfigRelease.ReleaseDesc)
	}

	return nil
}

func resourceTencentCloudTsfApplicationPublicConfigReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_release.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	namespaceId := idSplit[1]

	if err := service.DeleteTsfApplicationPublicConfigReleaseById(ctx, configId, namespaceId); err != nil {
		return err
	}

	return nil
}
