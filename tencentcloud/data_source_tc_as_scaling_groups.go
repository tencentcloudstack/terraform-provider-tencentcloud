package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudAsScalingGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingGroupRead,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"subnet_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"default_cooldown": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"desired_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"load_balancer_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"forward_balancer_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"listener_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"location_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_attribute": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"weight": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"termination_policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"retry_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	scalingGroupId := ""
	configurationId := ""
	scalingGroupName := ""
	if v, ok := d.GetOk("scaling_group_id"); ok {
		scalingGroupId = v.(string)
	}
	if v, ok := d.GetOk("configuration_id"); ok {
		configurationId = v.(string)
	}
	if v, ok := d.GetOk("scaling_group_name"); ok {
		scalingGroupName = v.(string)
	}

	scalingGroups, err := asService.DescribeAutoScalingGroupByFilter(ctx, scalingGroupId, configurationId, scalingGroupName)
	if err != nil {
		return err
	}

	scalingGroupList := make([]map[string]interface{}, 0, len(scalingGroups))
	for _, scalingGroup := range scalingGroups {
		mapping := map[string]interface{}{
			"scaling_group_id":     *scalingGroup.AutoScalingGroupId,
			"scaling_group_name":   *scalingGroup.AutoScalingGroupName,
			"configuration_id":     *scalingGroup.LaunchConfigurationId,
			"status":               *scalingGroup.AutoScalingGroupStatus,
			"instance_count":       *scalingGroup.InstanceCount,
			"max_size":             *scalingGroup.MaxSize,
			"min_size":             *scalingGroup.MinSize,
			"vpc_id":               *scalingGroup.VpcId,
			"subnet_ids":           flattenStringList(scalingGroup.SubnetIdSet),
			"zones":                flattenStringList(scalingGroup.ZoneSet),
			"default_cooldown":     *scalingGroup.DefaultCooldown,
			"desired_capacity":     *scalingGroup.DesiredCapacity,
			"load_balancer_ids":    flattenStringList(scalingGroup.LoadBalancerIdSet),
			"termination_policies": flattenStringList(scalingGroup.TerminationPolicySet),
			"retry_policy":         *scalingGroup.RetryPolicy,
			"create_time":          *scalingGroup.CreatedTime,
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
					"location_id":       *v.LocationId,
				}
				forwardLoadBalancers = append(forwardLoadBalancers, forwardLoadBalancer)
			}
			mapping["forward_load_balancers"] = forwardLoadBalancers
		}
		scalingGroupList = append(scalingGroupList, mapping)
	}

	d.SetId("ScalingGroupList" + scalingGroupId + scalingGroupName + configurationId)
	err = d.Set("scaling_group_list", scalingGroupList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set scaling group list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), scalingGroupList); err != nil {
			return err
		}
	}

	return nil
}
