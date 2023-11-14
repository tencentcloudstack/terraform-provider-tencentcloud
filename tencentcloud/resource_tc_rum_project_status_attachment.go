/*
Provides a resource to create a rum project_status_attachment

Example Usage

```hcl
resource "tencentcloud_rum_project_status_attachment" "project_status_attachment" {
  project_id = 101
}
```

Import

rum project_status_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_rum_project_status_attachment.project_status_attachment project_status_attachment_id
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

func resourceTencentCloudRumProjectStatusAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumProjectStatusAttachmentCreate,
		Read:   resourceTencentCloudRumProjectStatusAttachmentRead,
		Delete: resourceTencentCloudRumProjectStatusAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},
		},
	}
}

func resourceTencentCloudRumProjectStatusAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = rum.NewResumeProjectRequest()
		response  = rum.NewResumeProjectResponse()
		projectId int
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		projectId = v.(int64)
		request.ProjectId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().ResumeProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum projectStatusAttachment failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.ProjectId
	d.SetId(helper.Int64ToStr(projectId))

	return resourceTencentCloudRumProjectStatusAttachmentRead(d, meta)
}

func resourceTencentCloudRumProjectStatusAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectStatusAttachmentId := d.Id()

	projectStatusAttachment, err := service.DescribeRumProjectStatusAttachmentById(ctx, projectId)
	if err != nil {
		return err
	}

	if projectStatusAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumProjectStatusAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if projectStatusAttachment.ProjectId != nil {
		_ = d.Set("project_id", projectStatusAttachment.ProjectId)
	}

	return nil
}

func resourceTencentCloudRumProjectStatusAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	projectStatusAttachmentId := d.Id()

	if err := service.DeleteRumProjectStatusAttachmentById(ctx, projectId); err != nil {
		return err
	}

	return nil
}
