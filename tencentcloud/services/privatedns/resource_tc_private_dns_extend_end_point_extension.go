package privatedns

import (
	"context"
	"fmt"
	"log"

	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudPrivateDnsExtendEndPointReadPreHandleResponse0(ctx context.Context, resp *privatednsIntlv20201028.DescribeExtendEndpointListResponseParams) error {
	logId := tccommon.GetLogId(ctx)
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if resp.OutboundEndpointSet == nil && len(resp.OutboundEndpointSet) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `private_dns_extend_end_point` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, item := range resp.OutboundEndpointSet {
		if item.EndpointName != nil {
			_ = d.Set("end_point_name", item.EndpointName)
		}

		if item.Region != nil {
			_ = d.Set("end_point_region", item.Region)
		}

		if item.EndpointServiceSet != nil {
			tmpList := make([]map[string]interface{}, 0, len(item.EndpointServiceSet))
			for _, v := range item.EndpointServiceSet {
				dMap := make(map[string]interface{}, 0)
				if v.AccessType != nil {
					dMap["access_type"] = v.AccessType
				}

				if v.Pip != nil {
					dMap["host"] = v.Pip
				}

				if v.Pport != nil {
					dMap["port"] = v.Pport
				}

				if v.VpcId != nil {
					dMap["vpc_id"] = v.VpcId
				}

				if v.SubnetId != nil {
					dMap["subnet_id"] = v.SubnetId
				}

				if v.AccessGatewayId != nil {
					dMap["access_gateway_id"] = v.AccessGatewayId
				}

				if v.Vip != nil {
					dMap["vip"] = v.Vip
				}

				if v.Vport != nil {
					dMap["vport"] = v.Vport
				}

				if v.Proto != nil {
					dMap["proto"] = v.Proto
				}

				if v.SnatVipCidr != nil {
					dMap["snat_vip_cidr"] = v.SnatVipCidr
				}

				if v.SnatVipSet != nil {
					dMap["snat_vip_set"] = v.SnatVipSet
				}

				tmpList = append(tmpList, dMap)
			}

			_ = d.Set("forward_ip", tmpList)
		}
	}

	return nil
}
