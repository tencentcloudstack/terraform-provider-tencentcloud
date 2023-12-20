package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudClbTargetGroupAttachments() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetGroupAttachmentsCreate,
		Read:   resourceTencentCloudClbTargetGroupAttachmentsRead,
		Update: resourceTencentCloudClbTargetGroupAttachmentsUpdate,
		Delete: resourceTencentCloudClbTargetGroupAttachmentsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CLB instance ID",
			},
			"associations": {
				Required:    true,
				Type:        schema.TypeSet,
				MaxItems:    20,
				Description: "Association array, the combination cannot exceed 20",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"listener_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Listener ID",
						},
						"target_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target group ID",
						},
						"location_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Forwarding rule ID",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbTargetGroupAttachmentsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = clb.NewAssociateTargetGroupsRequest()
		loadBalancerId string
	)
	if v, ok := d.GetOk("load_balancer_id"); ok {
		loadBalancerId = v.(string)
	}
	if v, ok := d.GetOk("associations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			targetGroupAssociation := clb.TargetGroupAssociation{}
			targetGroupAssociation.LoadBalancerId = helper.String(loadBalancerId)
			if v, ok := dMap["listener_id"]; ok {
				targetGroupAssociation.ListenerId = helper.String(v.(string))
			}
			if v, ok := dMap["target_group_id"]; ok {
				targetGroupAssociation.TargetGroupId = helper.String(v.(string))
			}
			if v, ok := dMap["location_id"]; ok {
				targetGroupAssociation.LocationId = helper.String(v.(string))
			}
			request.Associations = append(request.Associations, &targetGroupAssociation)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().AssociateTargetGroups(request)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), result.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb targetGroupAttachments failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(loadBalancerId)

	return resourceTencentCloudClbTargetGroupAttachmentsRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	loadBalancerId := d.Id()
	associationsSet := make(map[string]struct{}, 0)
	targetGroupList := make([]string, 0)
	if v, ok := d.GetOk("associations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			ids := make([]string, 0)

			ids = append(ids, loadBalancerId)
			if v, ok := dMap["listener_id"]; ok {
				ids = append(ids, v.(string))
			} else {
				ids = append(ids, "null")
			}

			if v, ok := dMap["target_group_id"]; ok {
				ids = append(ids, v.(string))
				targetGroupList = append(targetGroupList, v.(string))
			} else {
				ids = append(ids, "null")
			}

			if v, ok := dMap["location_id"]; ok {
				ids = append(ids, v.(string))
			} else {
				ids = append(ids, "null")
			}

			associationsSet[strings.Join(ids, tccommon.FILED_SP)] = struct{}{}
		}
	}

	targetGroupAttachments, err := service.DescribeClbTargetGroupAttachmentsById(ctx, targetGroupList, associationsSet)
	if err != nil {
		return err
	}

	if len(targetGroupAttachments) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClbTargetGroupAttachments` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	var associationsList []interface{}

	for _, attachment := range targetGroupAttachments {
		info := strings.Split(attachment, tccommon.FILED_SP)
		if len(info) != 4 {
			return fmt.Errorf("id is broken,%s", info)
		}
		associationsMap := map[string]interface{}{}
		_ = d.Set("load_balancer_id", info[0])

		associationsMap["listener_id"] = info[1]

		associationsMap["target_group_id"] = info[2]

		associationsMap["location_id"] = info[3]

		associationsList = append(associationsList, associationsMap)
	}

	_ = d.Set("associations", associationsList)

	return nil
}

func resourceTencentCloudClbTargetGroupAttachmentsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudClbTargetGroupAttachmentsRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := clb.NewDisassociateTargetGroupsRequest()

	loadBalancerId := d.Id()
	if v, ok := d.GetOk("associations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			targetGroupAssociation := clb.TargetGroupAssociation{}
			targetGroupAssociation.LoadBalancerId = helper.String(loadBalancerId)
			if v, ok := dMap["listener_id"]; ok {
				targetGroupAssociation.ListenerId = helper.String(v.(string))
			}
			if v, ok := dMap["target_group_id"]; ok {
				targetGroupAssociation.TargetGroupId = helper.String(v.(string))
			}
			if v, ok := dMap["location_id"]; ok {
				targetGroupAssociation.LocationId = helper.String(v.(string))
			}
			request.Associations = append(request.Associations, &targetGroupAssociation)
		}
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().DisassociateTargetGroups(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), result.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
