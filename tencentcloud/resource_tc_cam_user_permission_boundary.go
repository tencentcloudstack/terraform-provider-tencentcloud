/*
Provides a resource to create a cam user_permission_boundary

Example Usage

```hcl
resource "tencentcloud_cam_user_permission_boundary" "user_permission_boundary" {
  target_uin =
  policy_id =
}
```

Import

cam user_permission_boundary can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_permission_boundary.user_permission_boundary user_permission_boundary_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCamUserPermissionBoundary() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserPermissionBoundaryCreate,
		Read:   resourceTencentCloudCamUserPermissionBoundaryRead,
		Delete: resourceTencentCloudCamUserPermissionBoundaryDelete,
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

func resourceTencentCloudCamUserPermissionBoundaryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_permission_boundary.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cam.NewPutUserPermissionsBoundaryRequest()
		response  = cam.NewPutUserPermissionsBoundaryResponse()
		targetUin int
		policyId  int
	)
	if v, ok := d.GetOkExists("target_uin"); ok {
		targetUin = v.(int64)
		request.TargetUin = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int64)
		request.PolicyId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().PutUserPermissionsBoundary(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam UserPermissionBoundary failed, reason:%+v", logId, err)
		return err
	}

	targetUin = *response.Response.TargetUin
	d.SetId(strings.Join([]string{helper.Int64ToStr(targetUin), helper.Int64ToStr(policyId)}, FILED_SP))

	return resourceTencentCloudCamUserPermissionBoundaryRead(d, meta)
}

func resourceTencentCloudCamUserPermissionBoundaryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_permission_boundary.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	targetUin := idSplit[0]
	policyId := idSplit[1]

	UserPermissionBoundary, err := service.DescribeCamUserPermissionBoundaryById(ctx, targetUin, policyId)
	if err != nil {
		return err
	}

	if UserPermissionBoundary == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamUserPermissionBoundary` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if UserPermissionBoundary.TargetUin != nil {
		_ = d.Set("target_uin", UserPermissionBoundary.TargetUin)
	}

	if UserPermissionBoundary.PolicyId != nil {
		_ = d.Set("policy_id", UserPermissionBoundary.PolicyId)
	}

	return nil
}

func resourceTencentCloudCamUserPermissionBoundaryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_user_permission_boundary.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	targetUin := idSplit[0]
	policyId := idSplit[1]

	if err := service.DeleteCamUserPermissionBoundaryById(ctx, targetUin, policyId); err != nil {
		return err
	}

	return nil
}
