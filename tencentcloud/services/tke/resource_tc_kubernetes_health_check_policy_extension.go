package tke

import (
	"context"
	"fmt"

	tkev20220501 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudKubernetesHealthCheckPolicyCreatePostFillRequest0(ctx context.Context, req *tkev20220501.CreateHealthCheckPolicyRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}
	if healthCheckPolicyMap, ok := helper.InterfacesHeadMap(d, "health_check_policy"); ok {
		healthCheckPolicy := tkev20220501.HealthCheckPolicy{}
		if v, ok := healthCheckPolicyMap["name"]; ok {
			name := v.(string)
			healthCheckPolicy.Name = helper.String(name)
		}
		if v, ok := healthCheckPolicyMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				rulesMap := item.(map[string]interface{})
				healthCheckPolicyRule := tkev20220501.HealthCheckPolicyRule{}
				if v, ok := rulesMap["auto_repair_enabled"]; ok {
					healthCheckPolicyRule.AutoRepairEnabled = helper.Bool(v.(bool))
				}
				if v, ok := rulesMap["enabled"]; ok {
					healthCheckPolicyRule.Enabled = helper.Bool(v.(bool))
				}
				if v, ok := rulesMap["name"]; ok {
					healthCheckPolicyRule.Name = helper.String(v.(string))
				}
				healthCheckPolicy.Rules = append(healthCheckPolicy.Rules, &healthCheckPolicyRule)
			}
		}
		req.HealthCheckPolicy = &healthCheckPolicy
	}

	return nil
}
