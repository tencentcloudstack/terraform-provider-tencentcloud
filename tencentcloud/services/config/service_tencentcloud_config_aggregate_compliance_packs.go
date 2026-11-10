// Copyright (c) 2018-2025 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func (me *ConfigService) DescribeConfigListAggregateCompliancePacksByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*configv20220802.ConfigCompliancePack, total uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := configv20220802.NewListAggregateCompliancePacksRequest()
	request.Limit = helper.Uint64(100)
	request.Offset = helper.Uint64(0)
	request.AccountGroupId = helper.String(paramMap["AccountGroupId"].(string))

	if v, ok := paramMap["CompliancePackName"]; ok {
		request.CompliancePackName = helper.String(v.(string))
	}

	if v, ok := paramMap["RiskLevel"]; ok {
		request.RiskLevel = v.([]*uint64)
	}

	if v, ok := paramMap["Status"]; ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := paramMap["ComplianceResult"]; ok {
		request.ComplianceResult = v.([]*string)
	}

	if v, ok := paramMap["OrderType"]; ok {
		request.OrderType = helper.String(v.(string))
	}

	if v, ok := paramMap["IncludeCompliancePackRuleResult"]; ok {
		request.IncludeCompliancePackRuleResult = helper.String(v.(string))
	}

	items = make([]*configv20220802.ConfigCompliancePack, 0)

	for {
		var (
			pageTotal uint64
			pageItems []*configv20220802.ConfigCompliancePack
		)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseConfigV20220802Client().ListAggregateCompliancePacksWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("list aggregate compliance packs failed, Response is nil"))
			}

			if result.Response.Total != nil {
				pageTotal = *result.Response.Total
			}
			pageItems = result.Response.Items

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}

		items = append(items, pageItems...)
		total = pageTotal

		if len(pageItems) < int(*request.Limit) || uint64(len(items)) >= total {
			break
		}
		request.Offset = helper.Uint64(*request.Offset + uint64(len(pageItems)))
	}

	return
}
