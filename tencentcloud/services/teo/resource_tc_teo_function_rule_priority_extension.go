package teo

import (
	"context"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudTeoFunctionRulePriorityReadPreHandleResponse0(ctx context.Context, resp *teov20220901.DescribeFunctionRulesResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	rulesList := []string{}
	if resp.FunctionRules != nil {
		for _, functionRules := range resp.FunctionRules {
			if functionRules.RuleId != nil {
				rulesList = append(rulesList, *functionRules.RuleId)
			}
		}

		_ = d.Set("rule_ids", rulesList)
	}

	return nil
}
