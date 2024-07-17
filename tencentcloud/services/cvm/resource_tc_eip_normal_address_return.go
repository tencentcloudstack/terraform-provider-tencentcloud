package cvm

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudEipNormalAddressReturn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipNormalAddressReturnCreate,
		Read:   resourceTencentCloudEipNormalAddressReturnRead,
		Delete: resourceTencentCloudEipNormalAddressReturnDelete,
		Schema: map[string]*schema.Schema{
			"address_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "The IP address of the EIP, example: 101.35.139.183.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudEipNormalAddressReturnCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_normal_address_return.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		addressIps string
	)
	var (
		request  = vpc.NewReturnNormalAddressesRequest()
		response = vpc.NewReturnNormalAddressesResponse()
	)

	if err := resourceTencentCloudEipNormalAddressReturnCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReturnNormalAddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eip normal address return failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(addressIps)

	if err := resourceTencentCloudEipNormalAddressReturnCreateOnExit(ctx); err != nil {
		return err
	}

	return resourceTencentCloudEipNormalAddressReturnRead(d, meta)
}

func resourceTencentCloudEipNormalAddressReturnRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_normal_address_return.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEipNormalAddressReturnDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_normal_address_return.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
