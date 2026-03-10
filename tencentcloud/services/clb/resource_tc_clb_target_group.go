package clb

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ForceNew:    true,
				Computed:    true,
				Description: "VPC ID, default is based on the network.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "The default port of target group, add server after can use it.",
			},
			"target_group_instances": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The backend server of target group bind.",
				Deprecated: "It has been deprecated from version 1.77.3. " +
					"please use `tencentcloud_clb_target_group_instance_attachment` instead.",
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Target group type, currently supported v1 (legacy version target group) and v2 (new version target group), defaults to v1 (legacy version target group).",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backend forwarding protocol of the target group. this field is required for the new version (v2) target group. currently supports TCP, UDP, HTTP, HTTPS, GRPC.",
			},
			"health_check": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Health check configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_switch": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to enable health check. true: enable, false: disable.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS", "PING", "CUSTOM", "GRPC"}, false),
							Description:  "Health check protocol. Valid values: TCP, HTTP, HTTPS, PING, CUSTOM, GRPC. Valid for v2 target groups.",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "Health check port. If not specified, the backend server port is used by default.",
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 60),
							Description:  "Health check response timeout in seconds. Range: [2, 60]. Default: 2.",
						},
						"gap_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 300),
							Description:  "Health check interval in seconds. Range: [2, 300]. Default: 5.",
						},
						"good_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 10),
							Description:  "Healthy threshold. Number of consecutive successful health checks required before marking the backend as healthy. Range: [2, 10]. Default: 3.",
						},
						"bad_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 10),
							Description:  "Unhealthy threshold. Number of consecutive failed health checks required before marking the backend as unhealthy. Range: [2, 10]. Default: 3.",
						},
						"http_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "HTTP status codes indicating health. For HTTP/HTTPS protocol. Example: 1 (1xx), 2 (2xx), 4 (3xx), 8 (4xx), 16 (5xx). Multiple values can be combined, e.g., 7 (1xx, 2xx, 3xx).",
						},
						"http_check_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Health check path. For HTTP/HTTPS protocol. Must start with /. If not specified, / is used by default.",
						},
						"http_check_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Health check domain. For HTTP/HTTPS protocol.",
						},
						"http_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"HEAD", "GET"}, false),
							Description:  "Health check HTTP method. For HTTP/HTTPS protocol. Valid values: HEAD, GET. Default: HEAD.",
						},
						"http_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP/1.0", "HTTP/1.1"}, false),
							Description:  "HTTP version for health check. Required when health check protocol is HTTP. Valid values: HTTP/1.0, HTTP/1.1. Only valid for TCP target groups.",
						},
						"extended_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Extended status code for health check.",
						},
					},
				},
			},
			"schedule_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"WRR", "LEAST_CONN", "IP_HASH"}, false),
				Description:  "Scheduling algorithm. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Valid values: WRR (weighted round robin), LEAST_CONN (least connections), IP_HASH (IP hash). Default: WRR.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource tags for the target group.",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
				Description:  "Default backend server weight. Range: [0, 100]. Only valid for v2 target groups. When set, backend servers added to the target group will use this default weight if not specified.",
			},
			"full_listen_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this is a full listener target group. Only valid for v2 target groups. true: full listener target group, false: normal target group.",
			},
			"keepalive_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable keep-alive connections. Only valid for HTTP/HTTPS target groups. true: enable, false: disable. Default: false.",
			},
			"session_expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					validation.IntBetween(30, 3600),
					validation.IntInSlice([]int{0}),
				),
				Description: "Session persistence time in seconds. Only valid for v2 target groups with HTTP/HTTPS/GRPC protocols. Range: 30-3600 or 0 (disabled). Default: 0 (disabled).",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP version type. Common values: IPv4, IPv6, IPv6FullChain.",
			},
		},
	}
}

func resourceTencentCloudClbTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.create")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService      = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		vpcId           = d.Get("vpc_id").(string)
		targetGroupName = d.Get("target_group_name").(string)
		port            = uint64(d.Get("port").(int))
		targetGroupType = d.Get("type").(string)
		protocol        = d.Get("protocol").(string)
		insAttachments  = make([]*clb.TargetGroupInstance, 0)
		targetGroupId   string
		err             error
	)

	if v, ok := d.GetOk("target_group_instances"); ok {
		targetGroupInstances := v.([]interface{})
		for _, v1 := range targetGroupInstances {
			value := v1.(map[string]interface{})
			bindIP := value["bind_ip"].(string)
			port := uint64(value["port"].(int))
			weight := uint64(value["weight"].(int))
			newPort := uint64(value["new_port"].(int))
			tgtGrp := &clb.TargetGroupInstance{
				BindIP:  &bindIP,
				Port:    &port,
				Weight:  &weight,
				NewPort: &newPort,
			}
			insAttachments = append(insAttachments, tgtGrp)
		}
	}

	// Extract new parameters
	var healthCheck *clb.TargetGroupHealthCheck
	if v, ok := d.GetOk("health_check"); ok && len(v.([]interface{})) > 0 {
		healthCheck = expandHealthCheck(v.([]interface{}))
	}

	scheduleAlgorithm := d.Get("schedule_algorithm").(string)

	var tags []*clb.TagInfo
	if v, ok := d.GetOk("tags"); ok {
		tags = expandTags(v.(map[string]interface{}))
	}

	var weight *uint64
	if v, ok := d.GetOk("weight"); ok {
		w := uint64(v.(int))
		weight = &w
	}

	var fullListenSwitch *bool
	if v, ok := d.GetOkExists("full_listen_switch"); ok {
		fullListenSwitch = helper.Bool(v.(bool))
	}

	var keepaliveEnable *bool
	if v, ok := d.GetOkExists("keepalive_enable"); ok {
		keepaliveEnable = helper.Bool(v.(bool))
	}

	var sessionExpireTime *uint64
	if v, ok := d.GetOk("session_expire_time"); ok {
		s := uint64(v.(int))
		sessionExpireTime = &s
	}

	ipVersion := d.Get("ip_version").(string)

	targetGroupId, err = clbService.CreateTargetGroup(ctx, targetGroupName, vpcId, port, insAttachments, targetGroupType, protocol,
		healthCheck, scheduleAlgorithm, tags, weight, fullListenSwitch, keepaliveEnable, sessionExpireTime, ipVersion)
	if err != nil {
		return err
	}
	d.SetId(targetGroupId)

	return resourceTencentCloudClbTargetRead(d, meta)

}

func resourceTencentCloudClbTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id         = d.Id()
	)
	filters := make(map[string]string)
	targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, id, filters)
	if err != nil {
		return err
	}
	if len(targetGroupInfos) < 1 {
		d.SetId("")
		return nil
	}

	targetGroup := targetGroupInfos[0]

	_ = d.Set("target_group_name", targetGroup.TargetGroupName)
	_ = d.Set("vpc_id", targetGroup.VpcId)
	_ = d.Set("port", targetGroup.Port)
	_ = d.Set("type", targetGroup.TargetGroupType)
	_ = d.Set("protocol", targetGroup.Protocol)

	// Set new parameters
	if targetGroup.HealthCheck != nil {
		_ = d.Set("health_check", flattenHealthCheck(targetGroup.HealthCheck))
	}

	if targetGroup.ScheduleAlgorithm != nil {
		_ = d.Set("schedule_algorithm", targetGroup.ScheduleAlgorithm)
	}

	if targetGroup.Tag != nil && len(targetGroup.Tag) > 0 {
		_ = d.Set("tags", flattenTags(targetGroup.Tag))
	}

	if targetGroup.Weight != nil {
		_ = d.Set("weight", targetGroup.Weight)
	}

	if targetGroup.FullListenSwitch != nil {
		_ = d.Set("full_listen_switch", targetGroup.FullListenSwitch)
	}

	if targetGroup.KeepaliveEnable != nil {
		_ = d.Set("keepalive_enable", targetGroup.KeepaliveEnable)
	}

	if targetGroup.SessionExpireTime != nil {
		_ = d.Set("session_expire_time", targetGroup.SessionExpireTime)
	}

	if targetGroup.IpVersion != nil {
		_ = d.Set("ip_version", targetGroup.IpVersion)
	}

	return nil
}

func resourceTencentCloudClbTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.update")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService    = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetGroupId = d.Id()
	)

	immutableFields := []string{"full_listen_switch", "ip_version", "vpc_id"}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			return fmt.Errorf("field %s cannot be modified after creation", field)
		}
	}

	isChanged := false
	request := clb.NewModifyTargetGroupAttributeRequest()
	request.TargetGroupId = &targetGroupId

	if d.HasChange("target_group_name") {
		request.TargetGroupName = helper.String(d.Get("target_group_name").(string))
		isChanged = true
	}

	if d.HasChange("port") {
		port := uint64(d.Get("port").(int))
		request.Port = &port
		isChanged = true
	}

	if d.HasChange("schedule_algorithm") {
		if v := d.Get("schedule_algorithm").(string); v != "" {
			request.ScheduleAlgorithm = helper.String(v)
			isChanged = true
		}
	}

	if d.HasChange("health_check") {
		if v, ok := d.GetOk("health_check"); ok && len(v.([]interface{})) > 0 {
			request.HealthCheck = expandHealthCheck(v.([]interface{}))
			isChanged = true
		}
	}

	if d.HasChange("weight") {
		if v, ok := d.GetOk("weight"); ok {
			w := uint64(v.(int))
			request.Weight = &w
			isChanged = true
		}
	}

	if d.HasChange("keepalive_enable") {
		request.KeepaliveEnable = helper.Bool(d.Get("keepalive_enable").(bool))
		isChanged = true
	}

	if d.HasChange("session_expire_time") {
		if v, ok := d.GetOk("session_expire_time"); ok {
			s := uint64(v.(int))
			request.SessionExpireTime = &s
			isChanged = true
		}
	}

	if isChanged {
		err := clbService.ModifyTargetGroupAttribute(ctx, request)
		if err != nil {
			return err
		}
	}

	// Handle tags separately if changed
	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("clb", "targetgroup", tcClient.Region, targetGroupId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
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

	err := clbService.DeleteTarget(ctx, targetGroupId)

	if err != nil {
		return err
	}
	return nil
}

// expandHealthCheck converts schema health_check to SDK TargetGroupHealthCheck
func expandHealthCheck(l []interface{}) *clb.TargetGroupHealthCheck {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	hcMap := l[0].(map[string]interface{})
	hc := &clb.TargetGroupHealthCheck{}

	if v, ok := hcMap["health_switch"].(bool); ok {
		hc.HealthSwitch = helper.Bool(v)
	}

	if v, ok := hcMap["protocol"].(string); ok && v != "" {
		hc.Protocol = helper.String(v)
	}

	if v, ok := hcMap["port"].(int); ok && v > 0 {
		hc.Port = helper.IntInt64(v)
	}

	if v, ok := hcMap["timeout"].(int); ok {
		hc.Timeout = helper.IntInt64(v)
	}

	if v, ok := hcMap["gap_time"].(int); ok {
		hc.GapTime = helper.IntInt64(v)
	}

	if v, ok := hcMap["good_limit"].(int); ok {
		hc.GoodLimit = helper.IntInt64(v)
	}

	if v, ok := hcMap["bad_limit"].(int); ok {
		hc.BadLimit = helper.IntInt64(v)
	}

	if v, ok := hcMap["http_code"].(int); ok && v > 0 {
		hc.HttpCode = helper.IntInt64(v)
	}

	if v, ok := hcMap["http_check_path"].(string); ok && v != "" {
		hc.HttpCheckPath = helper.String(v)
	}

	if v, ok := hcMap["http_check_domain"].(string); ok && v != "" {
		hc.HttpCheckDomain = helper.String(v)
	}

	if v, ok := hcMap["http_check_method"].(string); ok && v != "" {
		hc.HttpCheckMethod = helper.String(v)
	}

	if v, ok := hcMap["http_version"].(string); ok && v != "" {
		hc.HttpVersion = helper.String(v)
	}

	if v, ok := hcMap["extended_code"].(string); ok && v != "" {
		hc.ExtendedCode = helper.String(v)
	}

	return hc
}

// expandTags converts map[string]interface{} to []*clb.TagInfo
func expandTags(tags map[string]interface{}) []*clb.TagInfo {
	if len(tags) == 0 {
		return nil
	}

	tagInfos := make([]*clb.TagInfo, 0, len(tags))
	for k, v := range tags {
		tagInfo := &clb.TagInfo{
			TagKey:   helper.String(k),
			TagValue: helper.String(v.(string)),
		}
		tagInfos = append(tagInfos, tagInfo)
	}

	return tagInfos
}

// flattenHealthCheck converts SDK TargetGroupHealthCheck to schema health_check
func flattenHealthCheck(hc *clb.TargetGroupHealthCheck) []interface{} {
	if hc == nil {
		return nil
	}

	result := make(map[string]interface{})

	if hc.HealthSwitch != nil {
		result["health_switch"] = *hc.HealthSwitch
	}

	if hc.Protocol != nil {
		result["protocol"] = *hc.Protocol
	}

	if hc.Port != nil {
		result["port"] = *hc.Port
	}

	if hc.Timeout != nil {
		result["timeout"] = *hc.Timeout
	}

	if hc.GapTime != nil {
		result["gap_time"] = *hc.GapTime
	}

	if hc.GoodLimit != nil {
		result["good_limit"] = *hc.GoodLimit
	}

	if hc.BadLimit != nil {
		result["bad_limit"] = *hc.BadLimit
	}

	if hc.HttpCode != nil {
		result["http_code"] = *hc.HttpCode
	}

	if hc.HttpCheckPath != nil {
		result["http_check_path"] = *hc.HttpCheckPath
	}

	if hc.HttpCheckDomain != nil {
		result["http_check_domain"] = *hc.HttpCheckDomain
	}

	if hc.HttpCheckMethod != nil {
		result["http_check_method"] = *hc.HttpCheckMethod
	}

	if hc.HttpVersion != nil {
		result["http_version"] = *hc.HttpVersion
	}

	if hc.ExtendedCode != nil {
		result["extended_code"] = *hc.ExtendedCode
	}

	return []interface{}{result}
}

// flattenTags converts []*clb.TagInfo to map[string]string
func flattenTags(tagInfos []*clb.TagInfo) map[string]string {
	if len(tagInfos) == 0 {
		return nil
	}

	tags := make(map[string]string, len(tagInfos))
	for _, tagInfo := range tagInfos {
		if tagInfo.TagKey != nil && tagInfo.TagValue != nil {
			tags[*tagInfo.TagKey] = *tagInfo.TagValue
		}
	}

	return tags
}
