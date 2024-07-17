package cvm

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEipPublicAddressAdjust() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipPublicAddressAdjustCreate,
		Read:   resourceTencentCloudEipPublicAddressAdjustRead,
		Delete: resourceTencentCloudEipPublicAddressAdjustDelete,
		Schema: map[string]*schema.Schema{
			"address_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A unique ID that identifies an EIP instance. The unique ID of EIP is in the form:`eip-erft45fu`.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A unique ID that identifies the CVM instance. The unique ID of CVM is in the form:`ins-osckfnm7`.",
			},
		},
	}
}

func resourceTencentCloudEipPublicAddressAdjustCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_public_address_adjust.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
		addressId  string
	)
	var (
		request  = vpc.NewAdjustPublicAddressRequest()
		response = vpc.NewAdjustPublicAddressResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("address_id"); ok {
		addressId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_id"); ok {
		request.AddressId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AdjustPublicAddressWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eip public address adjust failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if err := resourceTencentCloudEipPublicAddressAdjustCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{instanceId, addressId}, tccommon.FILED_SP))

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
