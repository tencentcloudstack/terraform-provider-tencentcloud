package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudAsScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingGroupCreate,
		Read:   resourceTencentCloudAsScalingGroupRead,
		Update: resourceTencentCloudAsScalingGroupUpdate,
		Delete: resourceTencentCloudAsScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 55),
				Description:  "Name of a scaling group.",
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An available ID for a launch configuration.",
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
				Description:  "Maximum number of CVM instances. Valid value ranges: (0~2000).",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
				Description:  "Minimum number of CVM instances. Valid value ranges: (0~2000).",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of VPC network.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies to which project the scaling group belongs.",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID list of subnet, and for VPC it is required.",
			},
			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of available zones, for Basic network it is required.",
			},
			"default_cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Default cooldown time in second, and default value is `300`.",
			},
			"desired_capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Desired volume of CVM instances, which is between `max_size` and `min_size`.",
			},
			"load_balancer_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"forward_balancer_ids"},
				Description:   "ID list of traditional load balancers.",
			},
			"forward_balancer_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"load_balancer_ids"},
				Description:   "List of application load balancers, which can't be specified with `load_balancer_ids` together.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of available load balancers.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Listener ID for application load balancers.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of forwarding rules.",
						},
						"target_attribute": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Attribute list of target rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Port number.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Weight.",
									},
								},
							},
						},
					},
				},
			},
			"termination_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Available values for termination policies. Valid values: OLDEST_INSTANCE and NEWEST_INSTANCE.",
				Elem: &schema.Schema{
					Type:    schema.TypeString,
					Default: SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
					ValidateFunc: tccommon.ValidateAllowedStringValue([]string{SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
						SCALING_GROUP_TERMINATION_POLICY_NEWEST_INSTANCE}),
				},
			},
			"retry_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available values for retry policies. Valid values: IMMEDIATE_RETRY and INCREMENTAL_INTERVALS.",
				Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
					SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS}),
			},
			// Service Settings
			"replace_monitor_unhealthy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Enables unhealthy instance replacement. If set to `true`, AS will replace instances that are flagged as unhealthy by Cloud Monitor.",
			},
			"scaling_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Indicates scaling mode which creates and terminates instances (classic method), or method first tries to start stopped instances (wake up stopped) to perform scaling operations. Available values: `CLASSIC_SCALING`, `WAKE_UP_STOPPED_SCALING`. Default: `CLASSIC_SCALING`.",
			},
			"replace_load_balancer_unhealthy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Enable unhealthy instance replacement. If set to `true`, AS will replace instances that are found unhealthy in the CLB health check.",
			},
			"replace_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Replace mode of unhealthy replacement service. Valid values: RECREATE: Rebuild an instance to replace the original unhealthy instance. RESET: Performing a system reinstallation on unhealthy instances to keep information such as data disks, private IP addresses, and instance IDs unchanged. The instance login settings, HostName, enhanced services, and UserData will remain consistent with the current launch configuration. Default value: RECREATE. Note: This field may return null, indicating that no valid values can be obtained.",
			},
			"desired_capacity_sync_with_max_min_size": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "The expected number of instances is synchronized with the maximum and minimum values. The default value is `False`. This parameter is effective only in the scenario where the expected number is not passed in when modifying the scaling group interface. True: When modifying the maximum or minimum value, if there is a conflict with the current expected number, the expected number is adjusted synchronously. For example, when modifying, if the minimum value 2 is passed in and the current expected number is 1, the expected number is adjusted synchronously to 2; False: When modifying the maximum or minimum value, if there is a conflict with the current expected number, an error message is displayed indicating that the modification is not allowed.",
			},
			"priority_scale_in_unhealthy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable priority for unhealthy instances during scale-in operations. If set to `true`, unhealthy instances will be removed first when scaling in.",
			},
			"health_check_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Health check type of instances in a scaling group.<br><li>CVM: confirm whether an instance is healthy based on the network status. If the pinged instance is unreachable, the instance will be considered unhealthy. For more information, see [Instance Health Check](https://intl.cloud.tencent.com/document/product/377/8553?from_cn_redirect=1)<br><li>CLB: confirm whether an instance is healthy based on the CLB health check status. For more information, see [Health Check Overview](https://intl.cloud.tencent.com/document/product/214/6097?from_cn_redirect=1).<br>If the parameter is set to `CLB`, the scaling group will check both the network status and the CLB health check status. If the network check indicates unhealthy, the `HealthStatus` field will return `UNHEALTHY`. If the CLB health check indicates unhealthy, the `HealthStatus` field will return `CLB_UNHEALTHY`. If both checks indicate unhealthy, the `HealthStatus` field will return `UNHEALTHY|CLB_UNHEALTHY`. Default value: `CLB`.",
			},
			"lb_health_check_grace_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Grace period of the CLB health check during which the `IN_SERVICE` instances added will not be marked as `CLB_UNHEALTHY`.<br>Valid range: 0-7200, in seconds. Default value: `0`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of a scaling group.",
			},

			// computed value
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of a scaling group.",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance number of a scaling group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the AS group was created.",
			},
			"multi_zone_subnet_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{MultiZoneSubnetPolicyPriority,
					MultiZoneSubnetPolicyEquality}),
				Description: "Multi zone or subnet strategy, Valid values: PRIORITY and EQUALITY.",
			},
		},
	}
}

func resourceTencentCloudAsScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	request := as.NewCreateAutoScalingGroupRequest()

	request.AutoScalingGroupName = helper.String(d.Get("scaling_group_name").(string))
	request.LaunchConfigurationId = helper.String(d.Get("configuration_id").(string))
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	if v, ok := d.GetOk("default_cooldown"); ok {
		request.DefaultCooldown = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("retry_policy"); ok {
		request.RetryPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := v.([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}

	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIds := v.([]interface{})
		request.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerId)
		}
	}

	if v, ok := d.GetOk("forward_balancer_ids"); ok {
		forwardBalancers := v.(*schema.Set).List()
		request.ForwardLoadBalancers = make([]*as.ForwardLoadBalancer, 0, len(forwardBalancers))
		for _, v := range forwardBalancers {
			vv := v.(map[string]interface{})
			targets := vv["target_attribute"].([]interface{})
			forwardBalancer := as.ForwardLoadBalancer{
				LoadBalancerId: helper.String(vv["load_balancer_id"].(string)),
				ListenerId:     helper.String(vv["listener_id"].(string)),
				LocationId:     helper.String(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   helper.IntUint64(t["port"].(int)),
					Weight: helper.IntUint64(t["weight"].(int)),
				}
				forwardBalancer.TargetAttributes = append(forwardBalancer.TargetAttributes, &targetAttribute)
			}

			request.ForwardLoadBalancers = append(request.ForwardLoadBalancers, &forwardBalancer)
		}
	}

	if v, ok := d.GetOk("termination_policies"); ok {
		terminationPolicies := v.([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok {
		request.MultiZoneSubnetPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		request.HealthCheckType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("lb_health_check_grace_period"); ok {
		request.LoadBalancerHealthCheckGracePeriod = helper.IntUint64(v.(int))
	}

	var (
		replaceMonitorUnhealthy           = d.Get("replace_monitor_unhealthy").(bool)
		scalingMode                       = d.Get("scaling_mode").(string)
		replaceLBUnhealthy                = d.Get("replace_load_balancer_unhealthy").(bool)
		replaceMode                       = d.Get("replace_mode").(string)
		desiredCapacitySyncWithMaxMinSize = d.Get("desired_capacity_sync_with_max_min_size").(bool)
		priorityScaleInUnhealthy          = d.Get("priority_scale_in_unhealthy").(bool)
	)

	if replaceMonitorUnhealthy || scalingMode != "" || replaceLBUnhealthy || replaceMode != "" || desiredCapacitySyncWithMaxMinSize || priorityScaleInUnhealthy {
		if scalingMode == "" {
			scalingMode = SCALING_MODE_CLASSIC
		}

		if replaceMode == "" {
			replaceMode = REPLACE_MODE_RECREATE
		}

		request.ServiceSettings = &as.ServiceSettings{
			ReplaceMonitorUnhealthy:           &replaceMonitorUnhealthy,
			ScalingMode:                       &scalingMode,
			ReplaceLoadBalancerUnhealthy:      &replaceLBUnhealthy,
			ReplaceMode:                       &replaceMode,
			DesiredCapacitySyncWithMaxMinSize: &desiredCapacitySyncWithMaxMinSize,
			PriorityScaleInUnhealthy:          &priorityScaleInUnhealthy,
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := as.Tag{
				ResourceType: helper.String("auto-scaling-group"),
				Key:          helper.String(tagKey),
				Value:        helper.String(tagValue),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	var id string
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.AutoScalingGroupId == nil {
			err = fmt.Errorf("Create auto scaling group failed, Auto scaling group id is nil.")
			return resource.NonRetryableError(err)
		}

		id = *response.Response.AutoScalingGroupId

		return nil
	}); err != nil {
		return err
	}
	d.SetId(id)

	// wait for status
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, id)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if scalingGroup != nil && scalingGroup.InActivityStatus != nil && *scalingGroup.InActivityStatus == SCALING_GROUP_NOT_IN_ACTIVITY_STATUS {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("scaling group status is %s, retry...", *scalingGroup.InActivityStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var (
		scalingGroup *as.AutoScalingGroup
		e            error
		has          int
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scalingGroup, has, e = asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if has == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("scaling_group_name", scalingGroup.AutoScalingGroupName)
	_ = d.Set("configuration_id", scalingGroup.LaunchConfigurationId)
	_ = d.Set("status", scalingGroup.AutoScalingGroupStatus)
	_ = d.Set("instance_count", scalingGroup.InstanceCount)
	_ = d.Set("max_size", scalingGroup.MaxSize)
	_ = d.Set("min_size", scalingGroup.MinSize)
	_ = d.Set("vpc_id", scalingGroup.VpcId)
	_ = d.Set("project_id", scalingGroup.ProjectId)
	_ = d.Set("subnet_ids", helper.StringsInterfaces(scalingGroup.SubnetIdSet))
	_ = d.Set("zones", helper.StringsInterfaces(scalingGroup.ZoneSet))
	_ = d.Set("default_cooldown", scalingGroup.DefaultCooldown)
	_ = d.Set("desired_capacity", scalingGroup.DesiredCapacity)
	_ = d.Set("load_balancer_ids", helper.StringsInterfaces(scalingGroup.LoadBalancerIdSet))
	_ = d.Set("termination_policies", helper.StringsInterfaces(scalingGroup.TerminationPolicySet))
	_ = d.Set("retry_policy", scalingGroup.RetryPolicy)
	_ = d.Set("create_time", scalingGroup.CreatedTime)
	_ = d.Set("retry_policy", scalingGroup.RetryPolicy)
	_ = d.Set("health_check_type", scalingGroup.HealthCheckType)
	_ = d.Set("lb_health_check_grace_period", scalingGroup.LoadBalancerHealthCheckGracePeriod)
	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok && v.(string) != "" {
		_ = d.Set("multi_zone_subnet_policy", scalingGroup.MultiZoneSubnetPolicy)
	}

	if scalingGroup.ServiceSettings != nil {
		_ = d.Set("replace_monitor_unhealthy", scalingGroup.ServiceSettings.ReplaceMonitorUnhealthy)
		_ = d.Set("scaling_mode", scalingGroup.ServiceSettings.ScalingMode)
		_ = d.Set("replace_load_balancer_unhealthy", scalingGroup.ServiceSettings.ReplaceLoadBalancerUnhealthy)
		_ = d.Set("replace_mode", scalingGroup.ServiceSettings.ReplaceMode)
		_ = d.Set("desired_capacity_sync_with_max_min_size", scalingGroup.ServiceSettings.DesiredCapacitySyncWithMaxMinSize)
		_ = d.Set("priority_scale_in_unhealthy", scalingGroup.ServiceSettings.PriorityScaleInUnhealthy)
	}

	if scalingGroup.ForwardLoadBalancerSet != nil && len(scalingGroup.ForwardLoadBalancerSet) > 0 {
		forwardLoadBalancers := make([]map[string]interface{}, 0, len(scalingGroup.ForwardLoadBalancerSet))
		for _, v := range scalingGroup.ForwardLoadBalancerSet {
			targetAttributes := make([]map[string]interface{}, 0, len(v.TargetAttributes))
			for _, vv := range v.TargetAttributes {
				targetAttribute := map[string]interface{}{
					"port":   *vv.Port,
					"weight": *vv.Weight,
				}
				targetAttributes = append(targetAttributes, targetAttribute)
			}
			forwardLoadBalancer := map[string]interface{}{
				"load_balancer_id": *v.LoadBalancerId,
				"listener_id":      *v.ListenerId,
				"target_attribute": targetAttributes,
				"rule_id":          *v.LocationId,
			}
			forwardLoadBalancers = append(forwardLoadBalancers, forwardLoadBalancer)
		}
		_ = d.Set("forward_balancer_ids", forwardLoadBalancers)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "as", "auto-scaling-group", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAsScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region

	request := as.NewModifyAutoScalingGroupRequest()
	scalingGroupId := d.Id()

	d.Partial(true)

	var updateAttrs []string

	request.AutoScalingGroupId = &scalingGroupId
	if d.HasChange("scaling_group_name") {
		updateAttrs = append(updateAttrs, "scaling_group_name")
		request.AutoScalingGroupName = helper.String(d.Get("scaling_group_name").(string))
	}
	if d.HasChange("configuration_id") {
		updateAttrs = append(updateAttrs, "configuration_id")
		request.LaunchConfigurationId = helper.String(d.Get("configuration_id").(string))
	}
	if d.HasChange("max_size") {
		updateAttrs = append(updateAttrs, "max_size")
		request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	}
	if d.HasChange("min_size") {
		updateAttrs = append(updateAttrs, "max_size")
		request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	}
	if d.HasChange("vpc_id") {
		updateAttrs = append(updateAttrs, "vpc_id")
		request.VpcId = helper.String(d.Get("vpc_id").(string))
	}
	if d.HasChange("project_id") {
		updateAttrs = append(updateAttrs, "project_id")
		request.ProjectId = helper.IntUint64(d.Get("project_id").(int))
	}
	if d.HasChange("default_cooldown") {
		updateAttrs = append(updateAttrs, "default_cooldown")
		request.DefaultCooldown = helper.IntUint64(d.Get("default_cooldown").(int))
	}
	if d.HasChange("desired_capacity") {
		updateAttrs = append(updateAttrs, "desired_capacity")
		request.DesiredCapacity = helper.IntUint64(d.Get("desired_capacity").(int))
	}
	if d.HasChange("retry_policy") {
		updateAttrs = append(updateAttrs, "retry_policy")
		request.RetryPolicy = helper.String(d.Get("retry_policy").(string))
	}
	if d.HasChange("subnet_ids") {
		updateAttrs = append(updateAttrs, "subnet_ids")
		subnetIds := d.Get("subnet_ids").([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}
	if d.HasChange("zones") {
		updateAttrs = append(updateAttrs, "zones")
		zones := d.Get("zones").([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}
	if d.HasChange("termination_policies") {
		updateAttrs = append(updateAttrs, "termination_policies")
		terminationPolicies := d.Get("termination_policies").([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if d.HasChange("multi_zone_subnet_policy") {
		updateAttrs = append(updateAttrs, "multi_zone_subnet_policy")
		request.MultiZoneSubnetPolicy = helper.String(d.Get("multi_zone_subnet_policy").(string))
	}

	if d.HasChange("replace_monitor_unhealthy") ||
		d.HasChange("scaling_mode") ||
		d.HasChange("replace_load_balancer_unhealthy") ||
		d.HasChange("replace_mode") ||
		d.HasChange("desired_capacity_sync_with_max_min_size") ||
		d.HasChange("priority_scale_in_unhealthy") {
		updateAttrs = append(updateAttrs, "replace_monitor_unhealthy", "scaling_mode", "replace_load_balancer_unhealthy", "replace_mode", "desired_capacity_sync_with_max_min_size", "priority_scale_in_unhealthy")
		scalingMode := d.Get("scaling_mode").(string)
		replaceMode := d.Get("replace_mode").(string)
		if scalingMode == "" {
			scalingMode = SCALING_MODE_CLASSIC
		}
		if replaceMode == "" {
			replaceMode = REPLACE_MODE_RECREATE
		}
		replaceMonitor := d.Get("replace_monitor_unhealthy").(bool)
		replaceLB := d.Get("replace_load_balancer_unhealthy").(bool)
		desiredCapacitySyncWithMaxMinSize := d.Get("desired_capacity_sync_with_max_min_size").(bool)
		priorityScaleInUnhealthy := d.Get("priority_scale_in_unhealthy").(bool)
		request.ServiceSettings = &as.ServiceSettings{
			ReplaceMonitorUnhealthy:           &replaceMonitor,
			ScalingMode:                       &scalingMode,
			ReplaceLoadBalancerUnhealthy:      &replaceLB,
			ReplaceMode:                       &replaceMode,
			DesiredCapacitySyncWithMaxMinSize: &desiredCapacitySyncWithMaxMinSize,
			PriorityScaleInUnhealthy:          &priorityScaleInUnhealthy,
		}
	}

	if d.HasChange("health_check_type") || d.HasChange("lb_health_check_grace_period") {
		request.HealthCheckType = helper.String(d.Get("health_check_type").(string))
		if v, ok := d.GetOkExists("lb_health_check_grace_period"); ok {
			request.LoadBalancerHealthCheckGracePeriod = helper.IntUint64(v.(int))
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.UseAsClient().ModifyAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		return nil
	}); err != nil {
		return err
	}

	updateAttrs = updateAttrs[:0]

	balancerRequest := as.NewModifyLoadBalancersRequest()
	balancerRequest.AutoScalingGroupId = &scalingGroupId
	if d.HasChange("load_balancer_ids") {
		updateAttrs = append(updateAttrs, "load_balancer_ids")

		loadBalancerIds := d.Get("load_balancer_ids").([]interface{})
		balancerRequest.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			balancerRequest.LoadBalancerIds = append(balancerRequest.LoadBalancerIds, &loadBalancerId)
		}
	}

	if d.HasChange("forward_balancer_ids") {
		updateAttrs = append(updateAttrs, "forward_balancer_ids")

		forwardBalancers := d.Get("forward_balancer_ids").(*schema.Set).List()
		balancerRequest.ForwardLoadBalancers = make([]*as.ForwardLoadBalancer, 0, len(forwardBalancers))
		for _, v := range forwardBalancers {
			vv := v.(map[string]interface{})
			targets := vv["target_attribute"].([]interface{})
			forwardBalancer := as.ForwardLoadBalancer{
				LoadBalancerId: helper.String(vv["load_balancer_id"].(string)),
				ListenerId:     helper.String(vv["listener_id"].(string)),
				LocationId:     helper.String(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   helper.IntUint64(t["port"].(int)),
					Weight: helper.IntUint64(t["weight"].(int)),
				}
				forwardBalancer.TargetAttributes = append(forwardBalancer.TargetAttributes, &targetAttribute)
			}

			balancerRequest.ForwardLoadBalancers = append(balancerRequest.ForwardLoadBalancers, &forwardBalancer)
		}
	}

	if len(updateAttrs) > 0 {
		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(balancerRequest.GetAction())

			balancerResponse, err := client.UseAsClient().ModifyLoadBalancers(balancerRequest)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), err.Error())
				return tccommon.RetryError(err)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), balancerResponse.ToJsonString())

			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("as", "auto-scaling-group", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	// We need read the scaling group in order to check if there are instances.
	// If so, we need to remove those first.
	scalingGroup, has, err := asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	if *scalingGroup.InstanceCount > 0 || *scalingGroup.DesiredCapacity > 0 {
		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := asService.ClearScalingGroupInstance(ctx, scalingGroupId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		if errRet := asService.DeleteScalingGroup(ctx, scalingGroupId); errRet != nil {
			if sdkErr, ok := errRet.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == AsScalingGroupNotFound {
					return nil
				} else if sdkErr.Code == AsScalingGroupInProgress || sdkErr.Code == AsScalingGroupInstanceInGroup {
					return resource.RetryableError(sdkErr)
				}
			}
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
