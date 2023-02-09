/*
Provides a resource to create a tsf application_public_config_attachment

Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config_attachment" "application_public_config_attachment" {
  config_id = ""
  namespace_id = ""
  release_desc = ""
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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "Config id.",
			},

			"namespace_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace id.",
			},

			"release_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Release desc.",
			},

			"config_release_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config release id.",
			},

			"config_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config name.",
			},

			"config_version": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Config version.",
			},

			"release_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Release time.",
			},

			"group_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group id.",
			},

			"group_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group name.",
			},

			"namespace_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"cluster_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},

			"cluster_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster name.",
			},

			"application_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Application id.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationPublicConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tsf.NewReleasePublicConfigRequest()
		// response        = tsf.NewReleasePublicConfigResponse()
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
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationPublicConfigAttachment failed, reason:%+v", logId, err)
		return err
	}

	// configReleaseId = *response.Response.ConfigReleaseId
	d.SetId(configReleaseId)

	return resourceTencentCloudTsfApplicationPublicConfigAttachmentRead(d, meta)
}

func resourceTencentCloudTsfApplicationPublicConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	configReleaseId := d.Id()

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

	if applicationPublicConfigAttachment.ConfigReleaseId != nil {
		_ = d.Set("config_release_id", applicationPublicConfigAttachment.ConfigReleaseId)
	}

	if applicationPublicConfigAttachment.ConfigName != nil {
		_ = d.Set("config_name", applicationPublicConfigAttachment.ConfigName)
	}

	if applicationPublicConfigAttachment.ConfigVersion != nil {
		_ = d.Set("config_version", applicationPublicConfigAttachment.ConfigVersion)
	}

	if applicationPublicConfigAttachment.ReleaseTime != nil {
		_ = d.Set("release_time", applicationPublicConfigAttachment.ReleaseTime)
	}

	if applicationPublicConfigAttachment.GroupId != nil {
		_ = d.Set("group_id", applicationPublicConfigAttachment.GroupId)
	}

	if applicationPublicConfigAttachment.GroupName != nil {
		_ = d.Set("group_name", applicationPublicConfigAttachment.GroupName)
	}

	if applicationPublicConfigAttachment.NamespaceName != nil {
		_ = d.Set("namespace_name", applicationPublicConfigAttachment.NamespaceName)
	}

	if applicationPublicConfigAttachment.ClusterId != nil {
		_ = d.Set("cluster_id", applicationPublicConfigAttachment.ClusterId)
	}

	if applicationPublicConfigAttachment.ClusterName != nil {
		_ = d.Set("cluster_name", applicationPublicConfigAttachment.ClusterName)
	}

	if applicationPublicConfigAttachment.ApplicationId != nil {
		_ = d.Set("application_id", applicationPublicConfigAttachment.ApplicationId)
	}

	return nil
}

func resourceTencentCloudTsfApplicationPublicConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_public_config_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	configReleaseId := d.Id()

	if err := service.DeleteTsfApplicationPublicConfigAttachmentById(ctx, configReleaseId); err != nil {
		return err
	}

	return nil
}
