package cvm

import (
	"context"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudEipsReadOutputContent(ctx context.Context) interface{} {
	eipList := ctx.Value("eipList")
	return eipList
}

func dataSourceTencentCloudEipsReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *vpc.DescribeAddressesResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region

	tags := helper.GetTags(d, "tags")
	eipList := make([]map[string]interface{}, 0, len(resp.AddressSet))
	ids := make([]string, 0, len(resp.AddressSet))

EIP_LOOP:
	for _, eip := range resp.AddressSet {
		respTags, err := tagService.DescribeResourceTags(ctx, svcvpc.VPC_SERVICE_TYPE, svcvpc.EIP_RESOURCE_TYPE, region, *eip.AddressId)
		if err != nil {
			log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
			return err
		}

		for k, v := range tags {
			if respTags[k] != v {
				continue EIP_LOOP
			}
		}

		mapping := map[string]interface{}{
			"eip_id":      eip.AddressId,
			"eip_name":    eip.AddressName,
			"eip_type":    eip.AddressType,
			"status":      eip.AddressStatus,
			"public_ip":   eip.AddressIp,
			"instance_id": eip.InstanceId,
			"eni_id":      eip.NetworkInterfaceId,
			"create_time": eip.CreatedTime,
			"tags":        respTags,
		}

		eipList = append(eipList, mapping)
		ids = append(ids, *eip.AddressId)
	}

	context.WithValue(ctx, "eipList", eipList)
	d.SetId(helper.DataResourceIdsHash(ids))
	return nil
}
