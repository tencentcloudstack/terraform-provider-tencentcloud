/*
Provides a resource to create a CAM group.

Example Usage

```hcl
resource "tencentcloud_cam_group" "foo" {
  name   = "cam-group-test"
  remark = "test"
}
```

Import

CAM group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group.foo 90496
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamGroupCreate,
		Read:   resourceTencentCloudCamGroupRead,
		Update: resourceTencentCloudCamGroupUpdate,
		Delete: resourceTencentCloudCamGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of CAM group.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CAM group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM group.",
			},
		},
	}
}

func resourceTencentCloudCamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.create")()

	logId := getLogId(contextNil)

	request := cam.NewCreateGroupRequest()
	request.GroupName = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	var response *cam.CreateGroupResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateGroup(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "GroupNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.GroupId == nil {
		return fmt.Errorf("CAM group id is nil")
	}
	d.SetId(strconv.Itoa(int(*response.Response.GroupId)))

	//get really instance then read
	ctx := context.WithValue(context.TODO(), "logId", logId)
	groupId := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeGroupById(ctx, groupId)
		if e != nil {
			return retryError(e)
		}
		if instance == nil || instance.Response == nil || instance.Response.GroupId == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamGroupRead(d, meta)
}

func resourceTencentCloudCamGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	groupId := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.GetGroupResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeGroupById(ctx, groupId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil || instance.Response == nil || instance.Response.GroupId == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", *instance.Response.GroupName)
	_ = d.Set("create_time", *instance.Response.CreateTime)
	if instance.Response.Remark != nil {
		_ = d.Set("remark", *instance.Response.Remark)
	}
	return nil
}

func resourceTencentCloudCamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.update")()

	logId := getLogId(contextNil)

	groupId := d.Id()
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		return e
	}
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewUpdateGroupRequest()
	request.GroupId = &groupIdInt64
	changeFlag := false

	if d.HasChange("remark") {
		request.Remark = helper.String(d.Get("remark").(string))
		changeFlag = true

	}
	if d.HasChange("name") {
		request.GroupName = helper.String(d.Get("name").(string))
		changeFlag = true
	}

	if changeFlag {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateGroup(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM group description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCamGroupRead(d, meta)
}

func resourceTencentCloudCamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.delete")()

	logId := getLogId(contextNil)

	groupId := d.Id()
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		return e
	}
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewDeleteGroupRequest()
	request.GroupId = &groupIdInt64
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DeleteGroup(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM group failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
