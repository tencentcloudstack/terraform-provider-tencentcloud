/*
Provides a resource to create a dasb device_group_members

Example Usage

```hcl
resource "tencentcloud_dasb_device_group_members" "example" {
  device_group_id = 3
  member_id_set   = [1, 2, 3]
}
```

Import

dasb device_group_members can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_group_members.example 3#1,2,3
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbDeviceGroupMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbDeviceGroupMembersCreate,
		Read:   resourceTencentCloudDasbDeviceGroupMembersRead,
		Delete: resourceTencentCloudDasbDeviceGroupMembersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"device_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Device Group ID.",
			},
			"member_id_set": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "A collection of device IDs that need to be added to the device group.",
			},
		},
	}
}

func resourceTencentCloudDasbDeviceGroupMembersCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group_members.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		request        = dasb.NewAddDeviceGroupMembersRequest()
		deviceGroupId  string
		memberIdSetStr string
	)

	if v, ok := d.GetOkExists("device_group_id"); ok {
		request.Id = helper.IntUint64(v.(int))
		deviceGroupIdInt := v.(int)
		deviceGroupId = strconv.Itoa(deviceGroupIdInt)
	}

	if v, ok := d.GetOk("member_id_set"); ok {
		memberIdSetList := v.(*schema.Set).List()
		tmpList := make([]string, 0)
		for i := range memberIdSetList {
			memberIdSet := memberIdSetList[i].(int)
			request.MemberIdSet = append(request.MemberIdSet, helper.IntUint64(memberIdSet))
			tmpList = append(tmpList, strconv.Itoa(memberIdSet))
		}

		memberIdSetStr = strings.Join(tmpList, COMMA_SP)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().AddDeviceGroupMembers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("dasb DeviceGroupMembers not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb DeviceGroupMembers failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{deviceGroupId, memberIdSetStr}, FILED_SP))

	return resourceTencentCloudDasbDeviceGroupMembersRead(d, meta)
}

func resourceTencentCloudDasbDeviceGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group_members.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	deviceGroupId := idSplit[0]

	DeviceGroupMembers, err := service.DescribeDasbDeviceGroupMembersById(ctx, deviceGroupId)
	if err != nil {
		return err
	}

	if DeviceGroupMembers == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbDeviceGroupMembers` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("device_group_id", deviceGroupId)
	_ = d.Set("member_id_set", DeviceGroupMembers)

	return nil
}

func resourceTencentCloudDasbDeviceGroupMembersDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_group_members.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	deviceGroupId := idSplit[0]
	memberIdSetStr := idSplit[1]

	if err := service.DeleteDasbDeviceGroupMembersById(ctx, deviceGroupId, memberIdSetStr); err != nil {
		return err
	}

	return nil
}
