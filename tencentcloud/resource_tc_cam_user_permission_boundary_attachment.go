/*
Provides a resource to create a cam user_permission_boundary

Example Usage

```hcl
resource "tencentcloud_cam_user_permission_boundary_attachment" "user_permission_boundary" {
  target_uin = 100032767426
  policy_id = 151113272
}
```

Import

cam user_permission_boundary can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_permission_boundary_attachment.user_permission_boundary user_permission_boundary_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamUserPermissionBoundaryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserPermissionBoundaryAttachmentCreate,
		Read:   resourceTencentCloudCamUserPermissionBoundaryAttachmentRead,
		Delete: resourceTencentCloudCamUserPermissionBoundaryAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Sub account Uin.",
			},

			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID.",
			},
		},
	}
}

func resourceTencentCloudCamUserPermissionBoundaryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_permission_boundary_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cam.NewPutUserPermissionsBoundaryRequest()
		targetUin string
		policyId  string
	)
	if v, ok := d.GetOkExists("target_uin"); ok {
		targetUin = helper.IntToStr(v.(int))
		request.TargetUin = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = helper.IntToStr(v.(int))
		request.PolicyId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().PutUserPermissionsBoundary(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam UserPermissionBoundary failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(targetUin + FILED_SP + policyId)

	return resourceTencentCloudCamUserPermissionBoundaryAttachmentRead(d, meta)
}

func resourceTencentCloudCamUserPermissionBoundaryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_permission_boundary_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	targetUin := idSplit[0]

	UserPermissionBoundary, err := service.DescribeCamUserPermissionBoundaryById(ctx, targetUin)
	if err != nil {
		return err
	}

	if UserPermissionBoundary == nil || UserPermissionBoundary.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamUserPermissionBoundary` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if UserPermissionBoundary.Response.PolicyId != nil {
		_ = d.Set("policy_id", UserPermissionBoundary.Response.PolicyId)
	}
	return nil
}

func resourceTencentCloudCamUserPermissionBoundaryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_permission_boundary_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	targetUin := idSplit[0]

	if err := service.DeleteCamUserPermissionBoundaryById(ctx, targetUin); err != nil {
		return err
	}

	return nil
}
