/*
Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.

Example Usage

```hcl
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test-target-keep-1"
}

resource "tencentcloud_clb_target_group_attachment" "group" {
    clb_id          = tencentcloud_clb_instance.clb_basic.id
    listener_id     = tencentcloud_clb_listener.listener_basic.id
    rule_id         = tencentcloud_clb_listener_rule.rule_basic.id
    targrt_group_id = tencentcloud_clb_target_group.test.id
}
```

Import

CLB target group attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_attachment.group lbtg-odareyb2#lbl-bicjmx3i#lb-cv0iz74c#loc-ac6uk7b6
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbTargetGroupAttachment() *schema.Resource {
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
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB listener.",
			},
			"targrt_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the CLB target group.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CLB listener rule.",
			},
		},
	}
}

func resourceTencentCloudClbTargetGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_attachment.create")()

	var (
		clbService = ClbService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		listenerId       = d.Get("listener_id").(string)
		clbId            = d.Get("clb_id").(string)
		targrtGroupId    = d.Get("targrt_group_id").(string)
		locationId       = d.Get("rule_id").(string)
		targetInfos      []*clb.TargetGroupInfo
		instance         *clb.LoadBalancer
		listener         *clb.Listener
		isRuleExist, has bool
		err              error
	)

	//check rule
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		listener, err = clbService.DescribeListenerById(ctx, listenerId, clbId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, rule := range listener.Rules {
		if locationId == *rule.LocationId && (rule.TargetType != nil && *rule.TargetType == CLB_TARGET_TYPE_TARGETGROUP) {
			isRuleExist = true
			break
		}
	}
	if !isRuleExist {
		return fmt.Errorf("rule bound to the listener of the CLB instance does not exist or the rule not in targetgroup mode")
	}

	//check target group
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, err = clbService.DescribeLoadBalancerById(ctx, clbId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targetInfos, err = clbService.DescribeTargetGroups(ctx, targrtGroupId, nil)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(targetInfos) > 0 && (*targetInfos[0].VpcId != *instance.TargetRegionInfo.VpcId) {
		return fmt.Errorf("CLB instance needs to be in the same VPC as the backend target group")
	}

	err = clbService.AssociateTargetGroups(ctx, listenerId, clbId, targrtGroupId, locationId)
	if err != nil {
		return err
	}

	// wait status
	has, err = clbService.DescribeAssociateTargetGroups(ctx, []string{targrtGroupId, listenerId, clbId, locationId})
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("AssociateTargetGroups faild, targrtGroupId = %s, listenerId = %s, clbId = %s, ruleId = %s",
			targrtGroupId, listenerId, clbId, locationId)
	}

	d.SetId(strings.Join([]string{targrtGroupId, listenerId, clbId, locationId}, FILED_SP))

	return resourceTencentCloudClbTargetGroupAttachmentRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		clbService = ClbService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)
		id    = d.Id()
		has   bool
	)

	ids := strings.Split(id, FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("CLB target group attachment id must contains clb_id, listernrt_id, targrt_group_id, rule_id")
	}

	has, err := clbService.DescribeAssociateTargetGroups(ctx, ids)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("targrt_group_id", ids[0])
	_ = d.Set("listener_id", ids[1])
	_ = d.Set("clb_id", ids[2])
	_ = d.Set("rule_id", ids[3])

	return nil
}

func resourceTencentCloudClbTargetGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group_attachment.delete")()

	var (
		clbService = ClbService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		id          = d.Id()
		targetInfos []*clb.TargetGroupInfo
		err         error
	)

	ids := strings.Split(id, FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("CLB target group attachment id must contains clb_id, listernrt_id, targrt_group_id, rule_id")
	}

	if err := clbService.DisassociateTargetGroups(ctx, ids[0], ids[1], ids[2], ids[3]); err != nil {
		return err
	}

	// check status
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targetInfos, err = clbService.DescribeTargetGroups(ctx, ids[0], nil)
		if err != nil {
			return retryError(err, InternalError)
		}
		for _, info := range targetInfos {
			for _, rule := range info.AssociatedRule {
				var originLocationId string
				originClbId := *rule.LoadBalancerId
				originListenerId := *rule.ListenerId
				if rule.LocationId != nil {
					originLocationId = *rule.LocationId
				}

				if originListenerId == ids[1] && originClbId == ids[2] && originLocationId == ids[3] {
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
