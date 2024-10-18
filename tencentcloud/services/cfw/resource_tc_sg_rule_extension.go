package cfw

import (
	"context"
	"fmt"

	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSgRuleCreatePostFillRequest0(ctx context.Context, req *cfwv20190904.AddEnterpriseSecurityGroupRulesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	// 默认最后添加
	req.Type = helper.IntUint64(0)

	clientToken := helper.BuildToken()
	req.ClientToken = &clientToken

	if len(req.Data) == 1 {
		req.Data[0].OrderIndex = common.StringPtr("-1")
	}

	if v, ok := d.GetOkExists("enable"); ok {
		if v.(int) == 0 {
			return fmt.Errorf("enable must be `1` when creating")
		}
	}

	return nil
}

func resourceTencentCloudSgRuleReadPostFillRequest0(ctx context.Context, req *cfwv20190904.DescribeEnterpriseSecurityGroupRuleListRequest) error {
	req.Limit = common.Int64Ptr(100)
	req.Offset = common.Int64Ptr(0)
	return nil
}

func resourceTencentCloudSgRuleReadPreHandleResponse0(ctx context.Context, resp *cfwv20190904.DescribeEnterpriseSecurityGroupRuleListResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	dataList := make([]map[string]interface{}, 0, len(resp.Data))
	if resp.Data != nil {
		for _, data := range resp.Data {
			dataMap := map[string]interface{}{}

			if data.OrderIndex != nil {
				dataMap["order_index"] = fmt.Sprintf("%d", *data.OrderIndex)
			}

			if data.SourceId != nil {
				dataMap["source_content"] = data.SourceId
			}

			if data.SourceType != nil {
				dataMap["source_type"] = getRuleType(*data.SourceType)
			}

			if data.TargetId != nil {
				dataMap["dest_content"] = data.TargetId
			}

			if data.TargetType != nil {
				dataMap["dest_type"] = getRuleType(*data.TargetType)
			}

			if data.Protocol != nil {
				dataMap["protocol"] = data.Protocol
			}

			if data.Port != nil {
				dataMap["port"] = data.Port
			}

			if data.ServiceTemplateId != nil {
				dataMap["service_template_id"] = data.ServiceTemplateId
			}

			if data.Strategy != nil {
				if *data.Strategy == 2 {
					dataMap["rule_action"] = "accept"
				} else {
					dataMap["rule_action"] = "drop"
				}
			}

			if data.Detail != nil {
				dataMap["description"] = data.Detail
			}

			if data.Status != nil {
				_ = d.Set("enable", data.Status)
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	return nil
}

func resourceTencentCloudSgRuleUpdatePostFillRequest0(ctx context.Context, req *cfwv20190904.ModifyEnterpriseSecurityGroupRuleRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if ok := d.HasChange("enable"); ok {
		if ok := d.HasChange("data"); ok {
			return fmt.Errorf("cannot modify enable and data at the same time")
		}
		req.ModifyType = helper.IntUint64(1)
	} else {
		req.ModifyType = helper.IntUint64(0)
	}

	return nil
}

func resourceTencentCloudSgRuleDeletePostFillRequest0(ctx context.Context, req *cfwv20190904.RemoveEnterpriseSecurityGroupRuleRequest) error {
	req.RemoveType = helper.IntInt64(0)
	return nil
}

func getRuleType(t int64) string {
	switch t {
	case 0:
		return "net"
	case 1, 2, 3, 4, 5, 6:
		return "instance"
	case 7:
		return "template"
	case 8:
		return "tag"
	case 9:
		return "region"
	case 100:
		return "resourcegroup"
	default:
		return fmt.Sprintf("%d", t)
	}
}
