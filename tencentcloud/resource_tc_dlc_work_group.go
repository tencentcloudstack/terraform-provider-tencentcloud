/*
Provides a resource to create a dlc work_group

Example Usage

```hcl
resource "tencentcloud_dlc_work_group" "work_group" {
  work_group_name        = "tf-demo"
  work_group_description = "dlc workgroup test"
}
```

Import

dlc work_group can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_work_group.work_group work_group_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcWorkGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcWorkGroupCreate,
		Read:   resourceTencentCloudDlcWorkGroupRead,
		Update: resourceTencentCloudDlcWorkGroupUpdate,
		Delete: resourceTencentCloudDlcWorkGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"work_group_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of Work Group.",
			},

			"work_group_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of Work Group.",
			},

			"user_ids": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A collection of user IDs that has been bound to the workgroup.",
			},
		},
	}
}

func resourceTencentCloudDlcWorkGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_work_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = dlc.NewCreateWorkGroupRequest()
		response    = dlc.NewCreateWorkGroupResponse()
		workGroupId int64
	)
	if v, ok := d.GetOk("work_group_name"); ok {
		request.WorkGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_group_description"); ok {
		request.WorkGroupDescription = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().CreateWorkGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dlc workGroup failed, reason:%+v", logId, err)
		return err
	}

	workGroupId = *response.Response.WorkGroupId
	d.SetId(helper.Int64ToStr(workGroupId))

	return resourceTencentCloudDlcWorkGroupRead(d, meta)
}

func resourceTencentCloudDlcWorkGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_work_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	workGroupId := d.Id()

	workGroup, err := service.DescribeDlcWorkGroupById(ctx, workGroupId)
	if err != nil {
		return err
	}

	if workGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcWorkGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if workGroup.WorkGroupName != nil {
		_ = d.Set("work_group_name", workGroup.WorkGroupName)
	}

	if workGroup.WorkGroupDescription != nil {
		_ = d.Set("work_group_description", workGroup.WorkGroupDescription)
	}

	if workGroup.UserSet != nil {

		userIds := make([]*string, len(workGroup.UserSet))

		for _, user := range workGroup.UserSet {
			userIds = append(userIds, user.UserId)
		}

		_ = d.Set("user_ids", userIds)
	}

	return nil
}

func resourceTencentCloudDlcWorkGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_work_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dlc.NewModifyWorkGroupRequest()

	workGroupId := d.Id()

	request.WorkGroupId = helper.Int64(helper.StrToInt64(workGroupId))

	if d.HasChange("work_group_description") {
		if v, ok := d.GetOk("work_group_description"); ok {
			request.WorkGroupDescription = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().ModifyWorkGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dlc workGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDlcWorkGroupRead(d, meta)
}

func resourceTencentCloudDlcWorkGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_work_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}
	workGroupId := d.Id()

	if err := service.DeleteDlcWorkGroupById(ctx, workGroupId); err != nil {
		return err
	}

	return nil
}
