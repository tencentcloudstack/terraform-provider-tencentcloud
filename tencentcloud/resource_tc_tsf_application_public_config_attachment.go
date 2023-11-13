/*
Provides a resource to create a tsf application_public_config_attachment

Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config_attachment" "application_public_config_attachment" {
  config_id = "dcfg-p-123456"
  namespace_id = "namespace-123456"
  release_desc = "product version"
}
```

Import

tsf application_public_config_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_public_config_attachment.application_public_config_attachment application_public_config_attachment_id
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

func resourceTencentCloudTsfApplicationPublicConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationPublicConfigAttachmentCreate,
		Read:   resourceTencentCloudTsfApplicationPublicConfigAttachmentRead,
		Delete: resourceTencentCloudTsfApplicationPublicConfigAttachmentDelete,
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
				Description: "Namespace-id.",
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

func resourceTencentCloudTsfApplicationPublicConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = tsf.NewReleasePublicConfigRequest()
		response        = tsf.NewReleasePublicConfigResponse()
		configReleaseId string
	)
	if v, ok := d.GetOk("config_id"); ok {
		request.ConfigId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationPublicConfigAttachment failed, reason:%+v", logId, err)
		return err
	}

	configReleaseId = *response.Response.ConfigReleaseId
	d.SetId(configReleaseId)

	return resourceTencentCloudTsfApplicationPublicConfigAttachmentRead(d, meta)
}

func resourceTencentCloudTsfApplicationPublicConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationPublicConfigAttachmentId := d.Id()

	applicationPublicConfigAttachment, err := service.DescribeTsfApplicationPublicConfigAttachmentById(ctx, configReleaseId)
	if err != nil {
		return err
	}

	if applicationPublicConfigAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationPublicConfigAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationPublicConfigAttachment.ConfigId != nil {
		_ = d.Set("config_id", applicationPublicConfigAttachment.ConfigId)
	}

	if applicationPublicConfigAttachment.NamespaceId != nil {
		_ = d.Set("namespace_id", applicationPublicConfigAttachment.NamespaceId)
	}

	if applicationPublicConfigAttachment.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationPublicConfigAttachment.ReleaseDesc)
	}

	return nil
}

func resourceTencentCloudTsfApplicationPublicConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	applicationPublicConfigAttachmentId := d.Id()

	if err := service.DeleteTsfApplicationPublicConfigAttachmentById(ctx, configReleaseId); err != nil {
		return err
	}

	return nil
}
