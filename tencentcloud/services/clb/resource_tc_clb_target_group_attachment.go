package clb

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func ResourceTencentCloudClbTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetGroupAttachmentCreate,
		Read:   resourceTencentCloudClbTargetGroupAttachmentRead,
		Update: resourceTencentCloudClbTargetGroupAttachmentUpdate,
		Delete: resourceTencentCloudClbTargetGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB listener.",
			},
			"targrt_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CLB target group.",
				Deprecated:  "It has been deprecated from version 1.47.1. Use `target_group_id` instead.",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the CLB target group.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the CLB listener rule.",
			},
		},
	}
}
func resourceTencentCloudClbTargetGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachment.create")()

	var (
		clbService = ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		locationId    string
		listenerId    = d.Get("listener_id").(string)
		clbId         = d.Get("clb_id").(string)
		targetGroupId string

		targetInfos []*clb.TargetGroupInfo
		instance    *clb.LoadBalancer
		has         bool
		err         error
	)
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
	}
	vTarget, eHas := d.GetOk("target_group_id")
	vTargrt, rHas := d.GetOk("targrt_group_id")

	if eHas || rHas {
		if rHas {
			targetGroupId = vTargrt.(string)
		}
		if eHas {
			targetGroupId = vTarget.(string)
		}
	} else {
		return fmt.Errorf("'target_group_id' or 'targrt_group_id' at least set one, please use 'target_group_id'")
	}

	//check listenerId
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	//check ruleId
	checkErr = RuleIdCheck(locationId)
	if checkErr != nil {
		return checkErr
	}

	//check target group
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, err = clbService.DescribeLoadBalancerById(ctx, clbId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		targetInfos, err = clbService.DescribeTargetGroups(ctx, targetGroupId, nil)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(targetInfos) > 0 && (*targetInfos[0].VpcId != *instance.TargetRegionInfo.VpcId) {
		return fmt.Errorf("CLB instance needs to be in the same VPC as the backend target group")
	}

	err = clbService.AssociateTargetGroups(ctx, listenerId, clbId, targetGroupId, locationId)
	if err != nil {
		return err
	}

	// wait status
	has, err = clbService.DescribeAssociateTargetGroups(ctx, []string{targetGroupId, listenerId, clbId, locationId})
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("AssociateTargetGroups faild, targetGroupId = %s, listenerId = %s, clbId = %s, ruleId = %s",
			targetGroupId, listenerId, clbId, locationId)
	}

	d.SetId(strings.Join([]string{targetGroupId, listenerId, clbId, locationId}, tccommon.FILED_SP))

	return resourceTencentCloudClbTargetGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		clbService = ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id    = d.Id()
		has   bool
	)

	ids := strings.Split(id, tccommon.FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("CLB target group attachment id is clb_id#listener_id#target_group_id#rule_id(only required for 7 layer CLB)")
	}

	has, err := clbService.DescribeAssociateTargetGroups(ctx, ids)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("target_group_id", ids[0])
	_ = d.Set("listener_id", ids[1])
	_ = d.Set("clb_id", ids[2])
	if ids[3] != "" {
		_ = d.Set("rule_id", ids[3])
	}

	return nil
}

func resourceTencentCloudClbTargetGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachment.update")()
	return resourceTencentCloudClbTargetGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachment.delete")()

	var (
		clbService = ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id          = d.Id()
		targetInfos []*clb.TargetGroupInfo
		err         error
	)

	ids := strings.Split(id, tccommon.FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("CLB target group attachment id is clb_id#listener_id#target_group_id#rule_id(only required for 7 layer CLB)")
	}

	if err := clbService.DisassociateTargetGroups(ctx, ids[0], ids[1], ids[2], ids[3]); err != nil {
		return err
	}

	// check status
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		targetInfos, err = clbService.DescribeTargetGroups(ctx, ids[0], nil)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		for _, info := range targetInfos {
			for _, rule := range info.AssociatedRule {
				var originLocationId string
				originClbId := *rule.LoadBalancerId
				originListenerId := *rule.ListenerId
				if rule.LocationId != nil {
					originLocationId = *rule.LocationId
				}
				if *rule.Protocol == CLB_LISTENER_PROTOCOL_TCP || *rule.Protocol == CLB_LISTENER_PROTOCOL_UDP ||
					*rule.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL || *rule.Protocol == CLB_LISTENER_PROTOCOL_QUIC {
					if originListenerId == ids[1] && originClbId == ids[2] {
						return resource.RetryableError(
							fmt.Errorf("rule association target group instance still exist. [targetGroupId=%s, listenerId=%s, cldId=%s]",
								ids[0], ids[1], ids[2]))
					}
				} else if originListenerId == ids[1] && originClbId == ids[2] && originLocationId == ids[3] {
					return resource.RetryableError(
						fmt.Errorf("rule association target group instance still exist. [targetGroupId=%s, listenerId=%s, cldId=%s, ruleId=%s]",
							ids[0], ids[1], ids[2], ids[3]))
				}

			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
