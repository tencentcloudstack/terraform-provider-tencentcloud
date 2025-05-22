package privatedns

import (
	"context"
	"fmt"
	"log"

	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudPrivateDnsEndPointReadPreHandleResponse0(ctx context.Context, resp *privatednsIntlv20201028.DescribeEndPointListResponseParams) error {
	logId := tccommon.GetLogId(ctx)
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if resp.EndPointSet == nil && len(resp.EndPointSet) < 1 {
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_end_point` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	for _, item := range resp.EndPointSet {
		if item.EndPointName != nil {
			_ = d.Set("end_point_name", item.EndPointName)
		}

		if item.EndPointServiceId != nil {
			_ = d.Set("end_point_service_id", item.EndPointServiceId)
		}

		if item.RegionCode != nil {
			_ = d.Set("end_point_region", item.RegionCode)
		}

		if item.EndPointVipSet != nil {
			endPointVipSetLen := len(item.EndPointVipSet)
			tmpList := make([]string, 0, endPointVipSetLen)
			for _, v := range item.EndPointVipSet {
				tmpList = append(tmpList, *v)
			}

			_ = d.Set("ip_num", endPointVipSetLen)
			_ = d.Set("end_point_vip_set", tmpList)
		}
	}

	return nil
}
