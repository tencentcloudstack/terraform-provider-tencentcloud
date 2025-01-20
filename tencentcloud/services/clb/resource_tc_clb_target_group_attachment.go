package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func ResourceTencentCloudClbTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetGroupAttachmentCreate,
		Read:   resourceTencentCloudClbTargetGroupAttachmentRead,
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
			"target_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB target group.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the CLB listener.",
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
		clbService    = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request       = clb.NewAssociateTargetGroupsRequest()
		targetInfos   []*clb.TargetGroupInfo
		instance      *clb.LoadBalancer
		clbId         string
		listenerId    string
		targetGroupId string
		locationId    string
	)

	targetGroupAssociation := clb.TargetGroupAssociation{}
	if v, ok := d.GetOk("clb_id"); ok {
		targetGroupAssociation.LoadBalancerId = helper.String(v.(string))
		clbId = v.(string)
	}

	if v, ok := d.GetOk("listener_id"); ok {
		targetGroupAssociation.ListenerId = helper.String(v.(string))
		listenerId = v.(string)
	}

	if v, ok := d.GetOk("target_group_id"); ok {
		targetGroupAssociation.TargetGroupId = helper.String(v.(string))
		targetGroupId = v.(string)
	}

	if v, ok := d.GetOk("rule_id"); ok {
		targetGroupAssociation.LocationId = helper.String(v.(string))
		locationId = v.(string)
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
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLoadBalancerById(ctx, clbId)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeLoadBalancers response is nil."))
		}

		instance = result
		return nil
	})

	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeTargetGroups(ctx, targetGroupId, nil)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeTargetGroups response is nil."))
		}

		targetInfos = result
		return nil
	})

	if err != nil {
		return err
	}

	if len(targetInfos) > 0 && (*targetInfos[0].VpcId != *instance.TargetRegionInfo.VpcId) {
		return fmt.Errorf("CLB instance needs to be in the same VPC as the backend target group")
	}

	request.Associations = []*clb.TargetGroupAssociation{&targetGroupAssociation}
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().AssociateTargetGroups(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation.ResourceInOperating" {
					return resource.RetryableError(e)
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("AssociateTargetGroups response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{targetGroupId, listenerId, clbId, locationId}, tccommon.FILED_SP))

	return resourceTencentCloudClbTargetGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id         = d.Id()
		has        bool
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

func resourceTencentCloudClbTargetGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachment.delete")()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService  = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request     = clb.NewDisassociateTargetGroupsRequest()
		id          = d.Id()
		targetInfos []*clb.TargetGroupInfo
		err         error
	)

	ids := strings.Split(id, tccommon.FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("CLB target group attachment id is clb_id#listener_id#target_group_id#rule_id(only required for 7 layer CLB)")
	}

	request.Associations = []*clb.TargetGroupAssociation{
		{
			TargetGroupId:  &ids[0],
			ListenerId:     &ids[1],
			LoadBalancerId: &ids[2],
			LocationId:     &ids[3],
		},
	}
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().DisassociateTargetGroups(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation.ResourceInOperating" {
					return resource.RetryableError(e)
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("DisassociateTargetGroups response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}

		return nil
	})

	if err != nil {
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
