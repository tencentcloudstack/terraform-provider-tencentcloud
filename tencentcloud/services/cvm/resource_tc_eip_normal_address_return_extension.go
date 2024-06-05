package cvm

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEipNormalAddressReturnCreatePostFillRequest0(ctx context.Context, req *vpc.ReturnNormalAddressesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	var (
		addressIps string
	)

	if v, ok := d.GetOk("address_ips"); ok {
		addressIpsSet := v.(*schema.Set).List()
		for i := range addressIpsSet {
			addressIp := addressIpsSet[i].(string)
			req.AddressIps = append(req.AddressIps, &addressIp)
			addressIps = addressIp + tccommon.FILED_SP
		}
	}

	_ = addressIps
	return nil
}

func resourceTencentCloudEipNormalAddressReturnCreateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	var (
		addressIps string
	)

	if v, ok := d.GetOk("address_ips"); ok {
		addressIpsSet := v.(*schema.Set).List()
		for i := range addressIpsSet {
			addressIp := addressIpsSet[i].(string)
			addressIps = addressIp + tccommon.FILED_SP
		}
	}

	d.SetId(addressIps)
	return nil
}
