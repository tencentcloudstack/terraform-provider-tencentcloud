package clb

import (
	"context"
	"fmt"
	"log"
	"regexp"
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
		Delete: resourceTencentCloudClbTargetGroupAttachmentsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CLB instance ID, (load_balancer_id and target_group_id require at least one).",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Target group ID, (load_balancer_id and target_group_id require at least one).",
			},
			"associations": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				MaxItems:    20,
				Description: "Association array, the combination cannot exceed 20.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "CLB instance ID, when the binding target is target group, load_balancer_id in associations is required.",
						},
						"target_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Target group ID, when the binding target is clb, the target_group_id in associations is required.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "Listener ID.",
						},
						"location_id": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "Forwarding rule ID.",
						},
					},
				},
			},
		},
	}
}

const ResourcePrefixFromClb = "lb"

func checkParam(d *schema.ResourceData) error {
	_, clbIdExists := d.GetOk("load_balancer_id")
	_, groupIdExists := d.GetOk("target_group_id")

	if !clbIdExists && !groupIdExists {
		return fmt.Errorf("one of load_balancer_id and target_group_id must be specified as the binding target")
	}

	if clbIdExists && groupIdExists {
		return fmt.Errorf("load_balancer_id and target_group_id cannot be used as binding targets at the same time")
	}

	if v, ok := d.GetOk("associations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			associationsClbExists := false
			associationsGroupExists := false

			if v, ok := dMap["load_balancer_id"]; ok && v.(string) != "" {
				associationsClbExists = true
			}
			if v, ok := dMap["target_group_id"]; ok && v.(string) != "" {
				associationsGroupExists = true
			}

			if !associationsClbExists && !associationsGroupExists {
				return fmt.Errorf("then in associations, load_balancer_id and target_group_id must be filled in one")
			}

			if clbIdExists && associationsClbExists {
				return fmt.Errorf("if the binding target is clb, then in associations, it is expected that " +
					"target_group_id exists and load_balancer_id does not exist")
			}

			if groupIdExists && associationsGroupExists {
				return fmt.Errorf("if the binding target is targetGroup, then in associations, it is expected that " +
					"load_balancer_id exists and target_group_id does not exist")
			}
		}
	}
	return nil
}
func resourceTencentCloudClbTargetGroupAttachmentsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	if err := checkParam(d); err != nil {
		return err
	}
	var (
		request    = clb.NewAssociateTargetGroupsRequest()
		resourceId string
	)
	if v, ok := d.GetOk("load_balancer_id"); ok && v.(string) != "" {
		resourceId = v.(string)
		request.Associations = parseParamToRequest(d, "load_balancer_id", resourceId)
	}

	if v, ok := d.GetOk("target_group_id"); ok && v.(string) != "" {
		resourceId = v.(string)
		request.Associations = parseParamToRequest(d, "target_group_id", resourceId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().AssociateTargetGroups(request)
		if err != nil {
			if e := processRetryErrMsg(err); e != nil {
				return e
			}
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
		log.Printf("[CRITAL]%s create clb targetGroupAttachments failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(resourceId)

	return resourceTencentCloudClbTargetGroupAttachmentsRead(d, meta)
}

func resourceTencentCloudClbTargetGroupAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	targetGroupList, associationsSet := margeReadRequest(d)
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
		if isBindFromClb(d.Id()) {
			_ = d.Set("load_balancer_id", info[0])
			associationsMap["target_group_id"] = info[1]

		} else {
			_ = d.Set("target_group_id", info[1])
			associationsMap["load_balancer_id"] = info[0]
		}
		if info[2] != "" && info[2] != "null" {
			associationsMap["listener_id"] = info[2]
		}
		if info[3] != "" && info[3] != "null" {
			associationsMap["location_id"] = info[3]
		}
		associationsList = append(associationsList, associationsMap)
	}

	_ = d.Set("associations", associationsList)

	return nil
}

func resourceTencentCloudClbTargetGroupAttachmentsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group_attachments.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := clb.NewDisassociateTargetGroupsRequest()
	id := d.Id()
	if isBindFromClb(id) {
		request.Associations = parseParamToRequest(d, "load_balancer_id", id)
	} else {
		request.Associations = parseParamToRequest(d, "target_group_id", id)
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().DisassociateTargetGroups(request)
		if err != nil {
			if e := processRetryErrMsg(err); e != nil {
				return e
			}
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
func parseParamToRequest(d *schema.ResourceData, param string, id string) (associations []*clb.TargetGroupAssociation) {

	if v, ok := d.GetOk("associations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			targetGroupAssociation := clb.TargetGroupAssociation{}
			dMap[param] = id
			for name := range dMap {
				if dMap[name] != nil && dMap[name].(string) != "" {
					setString(name, dMap[name].(string), &targetGroupAssociation)
				}
			}
			associations = append(associations, &targetGroupAssociation)
		}
	}
	return associations
}
func setString(fieldName string, value string, request *clb.TargetGroupAssociation) {
	switch fieldName {
	case "load_balancer_id":
		request.LoadBalancerId = helper.String(value)
	case "target_group_id":
		request.TargetGroupId = helper.String(value)
	case "listener_id":
		request.ListenerId = helper.String(value)
	case "location_id":
		request.LocationId = helper.String(value)
	default:
		log.Printf("Invalid field name: %s\n", fieldName)
	}
}
func isBindFromClb(id string) bool {
	re := regexp.MustCompile(`^(.*?)-`)
	match := re.FindStringSubmatch(id)

	if len(match) > 1 {
		return match[1] == ResourcePrefixFromClb
	}
	return false
}
func margeReadRequest(d *schema.ResourceData) ([]string, map[string]struct{}) {
	resourceId := d.Id()

	associationsSet := make(map[string]struct{})
	targetGroupList := make([]string, 0)
	if v, ok := d.GetOk("associations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			ids := make([]string, 0)

			processIds(resourceId, dMap, "load_balancer_id", isBindFromClb(resourceId), &ids)
			processIds(resourceId, dMap, "target_group_id", isBindFromClb(resourceId), &ids)
			processIds(resourceId, dMap, "listener_id", isBindFromClb(resourceId), &ids)
			processIds(resourceId, dMap, "location_id", isBindFromClb(resourceId), &ids)

			if groupId, ok := dMap["target_group_id"]; ok && groupId.(string) != "" {
				targetGroupList = append(targetGroupList, groupId.(string))
			}

			associationsSet[strings.Join(ids, tccommon.FILED_SP)] = struct{}{}
		}
	}
	if len(targetGroupList) < 1 && !isBindFromClb(resourceId) {
		targetGroupList = append(targetGroupList, resourceId)
	}
	return targetGroupList, associationsSet
}
func processIds(id string, dMap map[string]interface{}, key string, clbFlag bool, ids *[]string) {
	if clbFlag && key == "load_balancer_id" {
		*ids = append(*ids, id)
	} else if !clbFlag && key == "target_group_id" {
		*ids = append(*ids, id)
	} else {
		if v, ok := dMap[key]; ok && v.(string) != "" {
			*ids = append(*ids, v.(string))
		} else {
			*ids = append(*ids, "null")
		}

	}
}
