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
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
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
				Description: "Name of cam group.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the cam group.",
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
	request.GroupName = stringToPointer(d.Get("name").(string))
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = stringToPointer(v.(string))
	}

	var response *cam.CreateGroupResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateGroup(request)
		if e != nil {
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
		log.Printf("[CRITAL]%s create cam group failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if response.Response.GroupId == nil {
		return fmt.Errorf("cam group id is nil")
	}
	d.SetId(strconv.Itoa(int(*response.Response.GroupId)))

	return resourceTencentCloudCamGroupRead(d, meta)
}

func resourceTencentCloudCamGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.read")()

	logId := getLogId(contextNil)

	groupId := d.Id()
	request := cam.NewGetGroupRequest()
	groupIdInt, _ := strconv.Atoi(groupId)
	groupIdInt64 := uint64(groupIdInt)
	request.GroupId = &groupIdInt64
	var instance *cam.GetGroupResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().GetGroup(request)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cam group failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.Set("name", *instance.Response.GroupName)
	d.Set("create_time", *instance.Response.CreateTime)
	if instance.Response.Remark != nil {
		d.Set("remark", *instance.Response.Remark)
	}
	return nil
}

func resourceTencentCloudCamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.update")()

	logId := getLogId(contextNil)

	groupId := d.Id()
	groupIdInt, _ := strconv.Atoi(groupId)
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewUpdateGroupRequest()
	request.GroupId = &groupIdInt64
	changeFlag := false

	if d.HasChange("remark") {
		request.Remark = stringToPointer(d.Get("remark").(string))
		changeFlag = true

	}
	if d.HasChange("name") {
		request.GroupName = stringToPointer(d.Get("name").(string))
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
			log.Printf("[CRITAL]%s update cam group description failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	return nil
}

func resourceTencentCloudCamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_group.delete")()

	logId := getLogId(contextNil)

	groupId := d.Id()
	groupIdInt, _ := strconv.Atoi(groupId)
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
		log.Printf("[CRITAL]%s delete cam group failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
