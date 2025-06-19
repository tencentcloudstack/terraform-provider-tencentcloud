package clb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func ResourceTencentCloudClbTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetCreate,
		Read:   resourceTencentCloudClbTargetRead,
		Update: resourceTencentCloudClbTargetUpdate,
		Delete: resourceTencentCloudClbTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TF_target_group",
				Description: "Target group name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				ForceNew:    true,
				Description: "VPC ID, default is based on the network.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "The default port of target group, add server after can use it. If `full_listen_switch` is true, setting this parameter is not supported.",
			},
			"target_group_instances": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The backend server of target group bind.",
				Deprecated:  "It has been deprecated from version 1.77.3. Please use `tencentcloud_clb_target_group_instance_attachment` instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateIp,
							Description:  "The internal ip of target group instance.",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "The port of target group instance.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The weight of target group instance.",
						},
						"new_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "The new port of target group instance.",
						},
					},
				},
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"v1", "v2"}),
				Description:  "Target group type, currently supports v1 (old version target group), v2 (new version target group), defaults to v1 (old version target group).",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"TCP", "UDP"}),
				Description:  "Target group backend forwarding protocol. This item is required for the v2 new version target group. Currently supports `TCP`, `UDP`.",
			},
			"tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Label.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
				Description:  "Default weights for backend services. Value range [0, 100]. After setting this value, when adding backend services to the target group, if the backend services do not have separate weights set, the default weights here will be used.",
			},
			"full_listen_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Full listening target group identifier, true indicates full listening target group, false indicates not full listening target group.",
			},
		},
	}
}

func resourceTencentCloudClbTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = clb.NewCreateTargetGroupRequest()
		response = clb.NewCreateTargetGroupResponse()
	)

	if v, ok := d.GetOk("target_group_name"); ok {
		request.TargetGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("port"); ok {
		request.Port = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			clbTags := clb.TagInfo{}
			if v, ok := dMap["tag_key"]; ok {
				clbTags.TagKey = helper.String(v.(string))
			}

			if v, ok := dMap["tag_value"]; ok {
				clbTags.TagValue = helper.String(v.(string))
			}

			request.Tags = append(request.Tags, &clbTags)
		}
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request.Weight = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("full_listen_switch"); ok {
		request.FullListenSwitch = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateTargetGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create target group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create target group failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.TargetGroupId == nil {
		return fmt.Errorf("TargetGroupId is nil.")
	}

	d.SetId(*response.Response.TargetGroupId)
	return resourceTencentCloudClbTargetRead(d, meta)
}

func resourceTencentCloudClbTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService    = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetGroupId = d.Id()
	)

	filters := make(map[string]string)
	targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, targetGroupId, filters)
	if err != nil {
		return err
	}

	if len(targetGroupInfos) < 1 {
		d.SetId("")
		return nil
	}

	if targetGroupInfos[0].TargetGroupName != nil {
		_ = d.Set("target_group_name", targetGroupInfos[0].TargetGroupName)
	}

	if targetGroupInfos[0].VpcId != nil {
		_ = d.Set("vpc_id", targetGroupInfos[0].VpcId)
	}

	if targetGroupInfos[0].Port != nil {
		_ = d.Set("port", targetGroupInfos[0].Port)
	}

	if targetGroupInfos[0].TargetGroupType != nil {
		_ = d.Set("type", targetGroupInfos[0].TargetGroupType)
	}

	if targetGroupInfos[0].Tag != nil {
		tagsList := make([]interface{}, 0, len(targetGroupInfos[0].Tag))
		for _, tags := range targetGroupInfos[0].Tag {
			tagsMap := map[string]interface{}{}
			if tags.TagKey != nil {
				tagsMap["tag_key"] = tags.TagKey
			}

			if tags.TagValue != nil {
				tagsMap["tag_value"] = tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)
	}

	if targetGroupInfos[0].Weight != nil {
		_ = d.Set("weight", targetGroupInfos[0].Weight)
	}

	if targetGroupInfos[0].FullListenSwitch != nil {
		_ = d.Set("full_listen_switch", targetGroupInfos[0].FullListenSwitch)
	}

	return nil
}

func resourceTencentCloudClbTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.update")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		targetGroupId = d.Id()
	)

	immutableArgs := []string{"type", "protocol", "tags", "full_listen_switch"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("target_group_name") || d.HasChange("port") || d.HasChange("weight") {
		request := clb.NewModifyTargetGroupAttributeRequest()
		if v, ok := d.GetOk("target_group_name"); ok {
			request.TargetGroupName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("port"); ok {
			request.Port = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("weight"); ok {
			request.Weight = helper.IntUint64(v.(int))
		}

		request.TargetGroupId = &targetGroupId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyTargetGroupAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify target group failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbTargetRead(d, meta)
}

func resourceTencentCloudClbTargetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.delete")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService    = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetGroupId = d.Id()
	)

	return clbService.DeleteTarget(ctx, targetGroupId)
}
