/*
Provides a resource to create a tsf application_file_config_attachment

Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config_attachment" "application_file_config_attachment" {
  config_id = ""
  group_id = ""
  release_desc = ""
}
```

Import

tsf application_file_config_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_file_config_attachment.application_file_config_attachment application_file_config_attachment_id
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

func resourceTencentCloudTsfApplicationFileConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationFileConfigAttachmentCreate,
		Read:   resourceTencentCloudTsfApplicationFileConfigAttachmentRead,
		Delete: resourceTencentCloudTsfApplicationFileConfigAttachmentDelete,
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

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group id.",
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

			"group_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group name.",
			},

			"namespace_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace id.",
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
		},
	}
}

func resourceTencentCloudTsfApplicationFileConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tsf.NewReleaseFileConfigRequest()
		// response        = tsf.NewReleaseFileConfigResponse()
		configReleaseId string
	)
	if v, ok := d.GetOk("config_id"); ok {
		request.ConfigId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
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
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationFileConfigAttachment failed, reason:%+v", logId, err)
		return err
	}

	// configReleaseId = *response.Response.ConfigReleaseId
	d.SetId(configReleaseId)

	return resourceTencentCloudTsfApplicationFileConfigAttachmentRead(d, meta)
}

func resourceTencentCloudTsfApplicationFileConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	configReleaseId := d.Id()

	applicationFileConfigAttachment, err := service.DescribeTsfApplicationFileConfigAttachmentById(ctx, configReleaseId)
	if err != nil {
		return err
	}

	if applicationFileConfigAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationFileConfigAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationFileConfigAttachment.ConfigId != nil {
		_ = d.Set("config_id", applicationFileConfigAttachment.ConfigId)
	}

	if applicationFileConfigAttachment.GroupId != nil {
		_ = d.Set("group_id", applicationFileConfigAttachment.GroupId)
	}

	if applicationFileConfigAttachment.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationFileConfigAttachment.ReleaseDesc)
	}

	if applicationFileConfigAttachment.ConfigReleaseId != nil {
		_ = d.Set("config_release_id", applicationFileConfigAttachment.ConfigReleaseId)
	}

	if applicationFileConfigAttachment.ConfigName != nil {
		_ = d.Set("config_name", applicationFileConfigAttachment.ConfigName)
	}

	if applicationFileConfigAttachment.ConfigVersion != nil {
		_ = d.Set("config_version", applicationFileConfigAttachment.ConfigVersion)
	}

	if applicationFileConfigAttachment.ReleaseTime != nil {
		_ = d.Set("release_time", applicationFileConfigAttachment.ReleaseTime)
	}

	if applicationFileConfigAttachment.GroupName != nil {
		_ = d.Set("group_name", applicationFileConfigAttachment.GroupName)
	}

	if applicationFileConfigAttachment.NamespaceId != nil {
		_ = d.Set("namespace_id", applicationFileConfigAttachment.NamespaceId)
	}

	if applicationFileConfigAttachment.NamespaceName != nil {
		_ = d.Set("namespace_name", applicationFileConfigAttachment.NamespaceName)
	}

	if applicationFileConfigAttachment.ClusterId != nil {
		_ = d.Set("cluster_id", applicationFileConfigAttachment.ClusterId)
	}

	if applicationFileConfigAttachment.ClusterName != nil {
		_ = d.Set("cluster_name", applicationFileConfigAttachment.ClusterName)
	}

	return nil
}

func resourceTencentCloudTsfApplicationFileConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_application_file_config_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	configReleaseId := d.Id()

	if err := service.DeleteTsfApplicationFileConfigAttachmentById(ctx, configReleaseId); err != nil {
		return err
	}

	return nil
}
