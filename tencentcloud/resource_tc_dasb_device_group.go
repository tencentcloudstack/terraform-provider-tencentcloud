/*
Provides a resource to create a dasb device_group

Example Usage

```hcl
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example"
  department_id = "1.2"
}
```

Import

dasb device_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_group.example 36
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
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbDeviceGroupCreate,
		Read:   resourceTencentCloudDasbDeviceGroupRead,
		Update: resourceTencentCloudDasbDeviceGroupUpdate,
		Delete: resourceTencentCloudDasbDeviceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Device group name, the maximum length is 32 characters.",
			},
			"department_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The ID of the department to which the asset group belongs, such as: 1.2.3 name, with a maximum length of 32 characters.",
			},
		},
	}
}

func resourceTencentCloudDasbDeviceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		request       = dasb.NewCreateDeviceGroupRequest()
		response      = dasb.NewCreateDeviceGroupResponse()
		deviceGroupId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().CreateDeviceGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.Id == nil {
			e = fmt.Errorf("dasb DeviceGroup not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb DeviceGroup failed, reason:%+v", logId, err)
		return err
	}

	deviceGroupIdInt := *response.Response.Id
	deviceGroupId = strconv.FormatUint(deviceGroupIdInt, 10)
	d.SetId(deviceGroupId)

	return resourceTencentCloudDasbDeviceGroupRead(d, meta)
}

func resourceTencentCloudDasbDeviceGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceGroupId = d.Id()
	)

	DeviceGroup, err := service.DescribeDasbDeviceGroupById(ctx, deviceGroupId)
	if err != nil {
		return err
	}

	if DeviceGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbDeviceGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DeviceGroup.Name != nil {
		_ = d.Set("name", DeviceGroup.Name)
	}

	if DeviceGroup.Department != nil {
		departmentId := DeviceGroup.Department.Id
		if *departmentId != "1" {
			_ = d.Set("department_id", departmentId)
		}
	}

	return nil
}

func resourceTencentCloudDasbDeviceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		request       = dasb.NewModifyDeviceGroupRequest()
		deviceGroupId = d.Id()
	)

	request.Id = helper.StrToUint64Point(deviceGroupId)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().ModifyDeviceGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dasb DeviceGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDasbDeviceGroupRead(d, meta)
}

func resourceTencentCloudDasbDeviceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceGroupId = d.Id()
	)

	if err := service.DeleteDasbDeviceGroupById(ctx, deviceGroupId); err != nil {
		return err
	}

	return nil
}
