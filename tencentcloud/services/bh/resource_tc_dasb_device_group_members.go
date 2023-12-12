package bh

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbDeviceGroupMembers() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_device_group_members.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
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

		memberIdSetStr = strings.Join(tmpList, tccommon.COMMA_SP)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().AddDeviceGroupMembers(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	d.SetId(strings.Join([]string{deviceGroupId, memberIdSetStr}, tccommon.FILED_SP))

	return resourceTencentCloudDasbDeviceGroupMembersRead(d, meta)
}

func resourceTencentCloudDasbDeviceGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_device_group_members.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_device_group_members.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
