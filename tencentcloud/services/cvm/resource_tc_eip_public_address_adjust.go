package cvm

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEipPublicAddressAdjust() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipPublicAddressAdjustCreate,
		Read:   resourceTencentCloudEipPublicAddressAdjustRead,
		Delete: resourceTencentCloudEipPublicAddressAdjustDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique ID that identifies the CVM instance. The unique ID of CVM is in the form:`ins-osckfnm7`.",
			},
			"address_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique ID that identifies an EIP instance. The unique ID of EIP is in the form:`eip-erft45fu`.",
			},
		},
	}
}

func resourceTencentCloudEipPublicAddressAdjustCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_public_address_adjust.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		service    = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request    = vpc.NewAdjustPublicAddressRequest()
		instanceId string
		addressId  string
		taskId     uint64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_id"); ok {
		addressId = v.(string)
		request.AddressId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AdjustPublicAddress(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate vpc publicAddressAdjust failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*tccommon.ReadRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(helper.UInt64ToStr(taskId), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId + tccommon.FILED_SP + addressId)
	return resourceTencentCloudEipPublicAddressAdjustRead(d, meta)
}

func resourceTencentCloudEipPublicAddressAdjustRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_public_address_adjust.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEipPublicAddressAdjustDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_public_address_adjust.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
