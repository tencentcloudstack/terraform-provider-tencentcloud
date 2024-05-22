package cvm

import (
	"context"
	"strconv"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func dataSourceTencentCloudReservedInstanceConfigsReadPreRequest0(ctx context.Context, req *cvm.DescribeReservedInstancesOfferingsRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	filters := make(map[string]string)
	if v, ok := d.GetOk("availability_zone"); ok {
		filters["zone"] = v.(string)
	}
	if v, ok := d.GetOk("duration"); ok {
		filters["duration"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		filters["instance-type"] = v.(string)
	}
	if v, ok := d.GetOk("offering_type"); ok {
		filters["offering-type"] = v.(string)
	}
	if v, ok := d.GetOk("product_description"); ok {
		filters["product-description"] = v.(string)
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
