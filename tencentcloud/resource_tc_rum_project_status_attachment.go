/*
Provides a resource to create a rum project_status_attachment

Example Usage

```hcl
resource "tencentcloud_rum_project_status_attachment" "project_status_attachment" {
  project_id = 131407
  operate    = "stop"
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
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
)

func resourceTencentCloudRumProjectStatusAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumProjectStatusAttachmentCreate,
		Read:   resourceTencentCloudRumProjectStatusAttachmentRead,
		Update: resourceTencentCloudRumProjectStatusAttachmentUpdate,
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
			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`resume`, `stop`.",
			},
		},
	}
}

func resourceTencentCloudRumProjectStatusAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.create")()
	defer inconsistentCheck(d, meta)()

	var projectId int
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(int)
	}

	d.SetId(strconv.Itoa(projectId))

	return resourceTencentCloudRumProjectStatusAttachmentRead(d, meta)
}

func resourceTencentCloudRumProjectStatusAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	projectId := d.Id()

	project, err := service.DescribeRumProjectStatusAttachmentById(ctx, projectId)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumProjectStatusAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if project.ID != nil {
		_ = d.Set("project_id", project.ID)
	}

	return nil
}

func resourceTencentCloudRumProjectStatusAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	projectId, _ := strconv.ParseInt(d.Id(), 10, 64)

	var operate string
	if v, ok := d.GetOk("operate"); ok {
		operate = v.(string)
	}

	if operate == "resume" {
		request := rum.NewResumeProjectRequest()
		request.ProjectId = &projectId
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().ResumeProject(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s resume rum project failed, reason:%+v", logId, err)
			return err
		}
	} else if operate == "stop" {
		request := rum.NewStopProjectRequest()
		request.ProjectId = &projectId
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().StopProject(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s stop rum project failed, reason:%+v", logId, err)
			return err
		}
	} else {
		return fmt.Errorf("[CRITAL]%s Operation type error", logId)
	}

	return resourceTencentCloudRumProjectStatusAttachmentRead(d, meta)
}

func resourceTencentCloudRumProjectStatusAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_project_status_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
