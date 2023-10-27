/*
Provides a resource to create a oceanus work_space

Example Usage

```hcl
resource "tencentcloud_oceanus_work_space" "example" {
  work_space_name = "tf_example"
  description     = "example description."
}
```

Import

oceanus work_space can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_work_space.example space-0dan3yux#tf_example
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
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOceanusWorkSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusWorkSpaceCreate,
		Read:   resourceTencentCloudOceanusWorkSpaceRead,
		Update: resourceTencentCloudOceanusWorkSpaceUpdate,
		Delete: resourceTencentCloudOceanusWorkSpaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"work_space_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workspace name.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace description.",
			},
			// computed
			"app_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "User APPID.",
			},
			"work_space_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Workspace ID.",
			},
			"serial_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Serial ID.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Workspace status.",
			},
			"role_auth_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of workspace members.",
			},
			"jobs_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of Jobs.",
			},
			"creator_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creator UIN.",
			},
			"owner_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Owner UIN.",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create time.",
			},
			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudOceanusWorkSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_work_space.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		request       = oceanus.NewCreateWorkSpaceRequest()
		response      = oceanus.NewCreateWorkSpaceResponse()
		workSpaceId   string
		workSpaceName string
	)

	if v, ok := d.GetOk("work_space_name"); ok {
		request.WorkSpaceName = helper.String(v.(string))
		workSpaceName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().CreateWorkSpace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus WorkSpace not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus WorkSpace failed, reason:%+v", logId, err)
		return err
	}

	workSpaceId = *response.Response.WorkSpaceId
	d.SetId(strings.Join([]string{workSpaceId, workSpaceName}, FILED_SP))

	return resourceTencentCloudOceanusWorkSpaceRead(d, meta)
}

func resourceTencentCloudOceanusWorkSpaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_work_space.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	workSpaceName := idSplit[1]

	WorkSpace, err := service.DescribeOceanusWorkSpaceById(ctx, workSpaceName)
	if err != nil {
		return err
	}

	if WorkSpace == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusWorkSpace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if WorkSpace.WorkSpaceName != nil {
		_ = d.Set("work_space_name", WorkSpace.WorkSpaceName)
	}

	if WorkSpace.Description != nil {
		_ = d.Set("description", WorkSpace.Description)
	}

	if WorkSpace.AppId != nil {
		_ = d.Set("app_id", WorkSpace.AppId)
	}

	if WorkSpace.WorkSpaceId != nil {
		_ = d.Set("work_space_id", WorkSpace.WorkSpaceId)
	}

	if WorkSpace.SerialId != nil {
		_ = d.Set("serial_id", WorkSpace.SerialId)
	}

	if WorkSpace.Status != nil {
		_ = d.Set("status", WorkSpace.Status)
	}

	if WorkSpace.RoleAuthCount != nil {
		_ = d.Set("role_auth_count", WorkSpace.RoleAuthCount)
	}

	if WorkSpace.JobsCount != nil {
		_ = d.Set("jobs_count", WorkSpace.JobsCount)
	}

	if WorkSpace.CreatorUin != nil {
		_ = d.Set("creator_uin", WorkSpace.CreatorUin)
	}

	if WorkSpace.OwnerUin != nil {
		_ = d.Set("owner_uin", WorkSpace.OwnerUin)
	}

	if WorkSpace.CreateTime != nil {
		_ = d.Set("create_time", WorkSpace.CreateTime)
	}

	if WorkSpace.UpdateTime != nil {
		_ = d.Set("update_time", WorkSpace.UpdateTime)
	}

	return nil
}

func resourceTencentCloudOceanusWorkSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_work_space.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = oceanus.NewModifyWorkSpaceRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	workSpaceId := idSplit[0]

	request.WorkSpaceId = &workSpaceId

	if v, ok := d.GetOk("work_space_name"); ok {
		request.WorkSpaceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().ModifyWorkSpace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update oceanus WorkSpace failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOceanusWorkSpaceRead(d, meta)
}

func resourceTencentCloudOceanusWorkSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_work_space.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	workSpaceId := idSplit[0]

	if err := service.DeleteOceanusWorkSpaceById(ctx, workSpaceId); err != nil {
		return err
	}

	return nil
}
