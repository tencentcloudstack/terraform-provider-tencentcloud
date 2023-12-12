package bh

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbBindDeviceResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbBindDeviceResourceCreate,
		Read:   resourceTencentCloudDasbBindDeviceResourceRead,
		Delete: resourceTencentCloudDasbBindDeviceResourceDelete,

		Schema: map[string]*schema.Schema{
			"device_id_set": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Asset ID collection.",
			},
			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bastion host service ID.",
			},
		},
	}
}

func resourceTencentCloudDasbBindDeviceResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_bind_device_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = dasb.NewBindDeviceResourceRequest()
		resourceId string
	)

	if v, ok := d.GetOk("device_id_set"); ok {
		deviceIdSetSet := v.(*schema.Set).List()
		for i := range deviceIdSetSet {
			deviceIdSet := deviceIdSetSet[i].(int)
			request.DeviceIdSet = append(request.DeviceIdSet, helper.IntUint64(deviceIdSet))
		}
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
		resourceId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().BindDeviceResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dasb bindDeviceResource failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(resourceId)
	return resourceTencentCloudDasbBindDeviceResourceRead(d, meta)
}

func resourceTencentCloudDasbBindDeviceResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_bind_device_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDasbBindDeviceResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_bind_device_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
