package cvm

import (
	"context"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudReservedInstancesReadPreRequest0(ctx context.Context, req *cvm.DescribeReservedInstancesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	filters := make(map[string]string)
	if v, ok := d.GetOk("reserved_instance_id"); ok {
		filters["reserved-instances-id"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		filters["zone"] = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		filters["instance-type"] = v.(string)
	}

	req.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		req.Filters = append(req.Filters, &filter)
	}
	return nil
}
