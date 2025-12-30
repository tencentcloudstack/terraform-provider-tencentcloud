package tke

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

func resourceTencentCloudKubernetesServerlessNodePoolCreatePostFillRequest0(ctx context.Context, req *tke.CreateClusterVirtualNodePoolRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	securityGroupIds := d.Get("security_group_ids").([]interface{})
	sgIds := make([]string, len(securityGroupIds))
	for i := 0; i < len(securityGroupIds); i++ {
		sgIds[i] = securityGroupIds[i].(string)
	}

	req.SecurityGroupIds = common.StringPtrs(sgIds)
	req.Labels = GetTkeLabels(d, "labels")

	return nil
}

func resourceTencentCloudKubernetesServerlessNodePoolReadRequestOnSuccess0(ctx context.Context, resp *tke.VirtualNodePool) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)

	nodePool := resp
	if nodePool.LifeState == nil {
		return resource.NonRetryableError(fmt.Errorf("LifeState is nil."))
	}

	if shouldServerlessNodePoolRetryReading(*nodePool.LifeState) {
		return resource.RetryableError(fmt.Errorf("serverless node pool %s is now %s, retrying", d.Id(), *nodePool.LifeState))
	}
	return nil
}

func resourceTencentCloudKubernetesServerlessNodePoolReadPostHandleResponse0(ctx context.Context, resp *tke.VirtualNodePool) error {
	d := tccommon.ResourceDataFromContext(ctx)

	labels := make(map[string]interface{})
	for i := 0; i < len(resp.Labels); i++ {
		if resp.Labels != nil && resp.Labels[i].Name != nil && resp.Labels[i].Value != nil {
			labels[*resp.Labels[i].Name] = *resp.Labels[i].Value
		}
	}
	_ = d.Set("labels", labels)

	return nil
}

func resourceTencentCloudKubernetesServerlessNodePoolUpdatePostFillRequest0(ctx context.Context, req *tke.ModifyClusterVirtualNodePoolRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if d.HasChange("labels") {
		req.Labels = GetOptimizedTkeLabels(d, "labels")
	}
	return nil
}

func shouldServerlessNodePoolRetryReading(state string) bool {
	return state != "normal"
}

func GetOptimizedTkeLabels(d *schema.ResourceData, k string) []*tke.Label {
	labels := make([]*tke.Label, 0)
	if raw, ok := d.GetOk(k); ok {
		for k, v := range raw.(map[string]interface{}) {
			labels = append(labels, &tke.Label{Name: helper.String(k), Value: common.StringPtr(v.(string))})
		}
	}
	return labels
}
