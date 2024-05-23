package cvm

import (
	"context"
	"fmt"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"strings"
)

func resourceTencentCloudCvmSecurityGroupAttachmentReadPostHandleResponse0(ctx context.Context, resp *cvm.Instance) error {
	d := tccommon.ResourceDataFromContext(ctx)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	for _, sgId := range resp.SecurityGroupIds {
		if *sgId == securityGroupId {
			_ = d.Set("instance_id", instanceId)
			_ = d.Set("security_group_id", securityGroupId)
			return nil
		}
	}

	_ = d.Set("instance_id", nil)
	_ = d.Set("security_group_id", nil)
	return nil
}
