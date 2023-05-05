/*
Provides a resource to create a group of AS (Auto scaling) instances.

Example Usage

```hcl
resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-scaling-group"
  configuration_id     = "asc-oqio4yyj"
  max_size             = 1
  min_size             = 0
  vpc_id               = "vpc-3efmz0z"
  subnet_ids           = ["subnet-mc3egos"]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  forward_balancer_ids {
    load_balancer_id = "lb-hk693b1l"
    listener_id      = "lbl-81wr497k"
    rule_id          = "loc-kiodx943"

    target_attribute {
      port   = 80
      weight = 90
    }
  }
}
```

Import

AutoScaling Groups can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_group.scaling_group asg-n32ymck2
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudAsScalingGroup() *schema.Resource {
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
				ValidateFunc: validateStringLengthInRange(1, 55),
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
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Maximum number of CVM instances. Valid value ranges: (0~2000).",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
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
				Type:          schema.TypeList,
				Optional:      true,
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
					ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
						SCALING_GROUP_TERMINATION_POLICY_NEWEST_INSTANCE}),
				},
			},
			"retry_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available values for retry policies. Valid values: IMMEDIATE_RETRY and INCREMENTAL_INTERVALS.",
				Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
					SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS}),
			},
			"scaling_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates scaling mode which creates and terminates instances (classic method), or method first tries to start stopped instances (wake up stopped) to perform scaling operations. Available values: `CLASSIC_SCALING`, `WAKE_UP_STOPPED_SCALING`. Default: `CLASSIC_SCALING`.",
			},
			// Service Settings
			"replace_monitor_unhealthy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables unhealthy instance replacement. If set to `true`, AS will replace instances that are flagged as unhealthy by Cloud Monitor.",
			},
			"replace_load_balancer_unhealthy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable unhealthy instance replacement. If set to `true`, AS will replace instances that are found unhealthy in the CLB health check.",
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
				ValidateFunc: validateAllowedStringValue([]string{MultiZoneSubnetPolicyPriority,
					MultiZoneSubnetPolicyEquality}),
				Description: "Multi zone or subnet strategy, Valid values: PRIORITY and EQUALITY.",
			},
		},
	}
}

func resourceTencentCloudAsScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
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
		forwardBalancers := v.([]interface{})
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

	var (
		scalingMode             = d.Get("scaling_mode").(string)
		replaceMonitorUnhealthy = d.Get("replace_monitor_unhealthy").(bool)
		replaceLBUnhealthy      = d.Get("replace_load_balancer_unhealthy").(bool)
	)

	if scalingMode != "" || replaceMonitorUnhealthy || replaceLBUnhealthy {
		if scalingMode == "" {
			scalingMode = SCALING_MODE_CLASSIC
		}

		request.ServiceSettings = &as.ServiceSettings{
			ScalingMode:                  &scalingMode,
			ReplaceMonitorUnhealthy:      &replaceMonitorUnhealthy,
			ReplaceLoadBalancerUnhealthy: &replaceLBUnhealthy,
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			request.Tags = append(request.Tags, &as.Tag{
				ResourceType: helper.String("auto-scaling-group"),
				Key:          &k,
				Value:        &v,
			})
		}
	}

	var id string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.AutoScalingGroupId == nil {
			err = fmt.Errorf("Auto scaling group id is nil")
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
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, id)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if scalingGroup != nil && *scalingGroup.InActivityStatus == SCALING_GROUP_NOT_IN_ACTIVITY_STATUS {
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
	defer logElapsed("resource.tencentcloud_as_scaling_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var (
		scalingGroup *as.AutoScalingGroup
		e            error
		has          int
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		scalingGroup, has, e = asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
		if e != nil {
			return retryError(e)
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
	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok && v.(string) != "" {
		_ = d.Set("multi_zone_subnet_policy", scalingGroup.MultiZoneSubnetPolicy)
	}

	if v := d.Get("scaling_mode"); v != "" {
		_ = d.Set("scaling_mode", v.(string))
	}

	if v, ok := d.GetOk("replace_monitor_unhealthy"); ok {
		_ = d.Set("replace_monitor_unhealthy", v.(bool))
	}

	if v, ok := d.GetOk("replace_load_balancer_unhealthy"); ok {
		_ = d.Set("replace_load_balancer_unhealthy", v.(bool))
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
				"load_balancer_id":  *v.LoadBalancerId,
				"listener_id":       *v.ListenerId,
				"target_attributes": targetAttributes,
				"rule_id":           *v.LocationId,
			}
			forwardLoadBalancers = append(forwardLoadBalancers, forwardLoadBalancer)
		}
		_ = d.Set("forward_balancer_ids", forwardLoadBalancers)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "as", "auto-scaling-group", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAsScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
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

	if d.HasChange("scaling_mode") ||
		d.HasChange("replace_monitor_unhealthy") ||
		d.HasChange("replace_load_balancer_unhealthy") {
		updateAttrs = append(updateAttrs, "scaling_mode", "replace_monitor_unhealthy", "replace_load_balancer_unhealthy")
		scalingMode := d.Get("scaling_mode").(string)
		if scalingMode == "" {
			scalingMode = SCALING_MODE_CLASSIC
		}
		replaceMonitor := d.Get("replace_monitor_unhealthy").(bool)
		replaceLB := d.Get("replace_load_balancer_unhealthy").(bool)
		request.ServiceSettings = &as.ServiceSettings{
			ScalingMode:                  &scalingMode,
			ReplaceMonitorUnhealthy:      &replaceMonitor,
			ReplaceLoadBalancerUnhealthy: &replaceLB,
		}
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.UseAsClient().ModifyAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
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

		forwardBalancers := d.Get("forward_balancer_ids").([]interface{})
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
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(balancerRequest.GetAction())

			balancerResponse, err := client.UseAsClient().ModifyLoadBalancers(balancerRequest)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), err.Error())
				return retryError(err)
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
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		resourceName := BuildTagResourceName("as", "auto-scaling-group", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_group.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
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
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := asService.ClearScalingGroupInstance(ctx, scalingGroupId)
			if inErr != nil {
				return retryError(inErr)
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
