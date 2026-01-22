package gwlb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gwlbv20240906 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGwlbInstanceAssociateTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGwlbInstanceAssociateTargetGroupCreate,
		Read:   resourceTencentCloudGwlbInstanceAssociateTargetGroupRead,
		Delete: resourceTencentCloudGwlbInstanceAssociateTargetGroupDelete,
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GWLB instance ID.",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target group ID.",
			},
		},
	}
}

func resourceTencentCloudGwlbInstanceAssociateTargetGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance_associate_target_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	var (
		request  = gwlbv20240906.NewAssociateTargetGroupsRequest()
		response = gwlbv20240906.NewAssociateTargetGroupsResponse()
	)
	loadBalancerId := d.Get("load_balancer_id").(string)
	targetGroupId := d.Get("target_group_id").(string)
	request.Associations = []*gwlbv20240906.TargetGroupAssociation{
		{
			LoadBalancerId: helper.String(loadBalancerId),
			TargetGroupId:  helper.String(targetGroupId),
		},
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().AssociateTargetGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gwlb instance associate target groups failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(strings.Join([]string{loadBalancerId, targetGroupId}, tccommon.FILED_SP))

	return resourceTencentCloudGwlbInstanceAssociateTargetGroupRead(d, meta)
}

func resourceTencentCloudGwlbInstanceAssociateTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance_associate_target_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	_ = d.Set("load_balancer_id", items[0])
	_ = d.Set("target_group_id", items[1])
	_ = ctx
	return nil
}

func resourceTencentCloudGwlbInstanceAssociateTargetGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance_associate_target_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	var (
		request  = gwlbv20240906.NewDisassociateTargetGroupsRequest()
		response = gwlbv20240906.NewDisassociateTargetGroupsResponse()
	)

	request.Associations = []*gwlbv20240906.TargetGroupAssociation{
		{
			LoadBalancerId: helper.String(items[0]),
			TargetGroupId:  helper.String(items[1]),
		},
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().DisassociateTargetGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete gwlb instance associate target groups failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
