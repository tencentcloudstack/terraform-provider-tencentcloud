package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsLoadBalancerCreate,
		Read:   resourceTencentCloudAsLoadBalancerRead,
		Update: resourceTencentCloudAsLoadBalancerUpdate,
		Delete: resourceTencentCloudAsLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of a scaling group.",
			},

			"forward_load_balancers": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of application load balancers. The maximum number of application-type load balancers bound to each scaling group is 100.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application load balancer instance ID.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application load balancer listener ID.",
						},
						"target_attributes": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of TargetAttribute.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Target port.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Target weight.",
									},
								},
							},
						},
						"location_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application load balancer location ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Load balancer instance region. Default value is the region of current auto scaling group. The format is the same as the public parameter Region, for example: ap-guangzhou.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAsLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_load_balancer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		request            = as.NewAttachLoadBalancersRequest()
		autoScalingGroupId string
	)

	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		request.AutoScalingGroupId = helper.String(v.(string))
		autoScalingGroupId = v.(string)
	}

	if v, ok := d.GetOk("forward_load_balancers"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			forwardLoadBalancer := as.ForwardLoadBalancer{}
			if v, ok := dMap["load_balancer_id"]; ok {
				forwardLoadBalancer.LoadBalancerId = helper.String(v.(string))
			}

			if v, ok := dMap["listener_id"]; ok {
				forwardLoadBalancer.ListenerId = helper.String(v.(string))
			}

			if v, ok := dMap["target_attributes"]; ok {
				for _, item := range v.([]interface{}) {
					targetAttributesMap := item.(map[string]interface{})
					targetAttribute := as.TargetAttribute{}
					if v, ok := targetAttributesMap["port"]; ok {
						targetAttribute.Port = helper.IntUint64(v.(int))
					}

					if v, ok := targetAttributesMap["weight"]; ok {
						targetAttribute.Weight = helper.IntUint64(v.(int))
					}

					forwardLoadBalancer.TargetAttributes = append(forwardLoadBalancer.TargetAttributes, &targetAttribute)
				}
			}

			if v, ok := dMap["location_id"]; ok {
				forwardLoadBalancer.LocationId = helper.String(v.(string))
			}

			if v, ok := dMap["region"]; ok {
				forwardLoadBalancer.Region = helper.String(v.(string))
			}

			request.ForwardLoadBalancers = append(request.ForwardLoadBalancers, &forwardLoadBalancer)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().AttachLoadBalancers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create as loadBalancer failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(autoScalingGroupId)

	return resourceTencentCloudAsLoadBalancerRead(d, meta)
}

func resourceTencentCloudAsLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_load_balancer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		autoScalingGroupId = d.Id()
	)

	loadBalancer, err := service.DescribeAsLoadBalancerById(ctx, autoScalingGroupId)
	if err != nil {
		return err
	}

	if loadBalancer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AsLoadBalancer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if loadBalancer.AutoScalingGroupId != nil {
		_ = d.Set("auto_scaling_group_id", loadBalancer.AutoScalingGroupId)
	}

	if loadBalancer.ForwardLoadBalancerSet != nil {
		forwardLoadBalancersList := make([]map[string]interface{}, 0, len(loadBalancer.ForwardLoadBalancerSet))
		for _, forwardLoadBalancers := range loadBalancer.ForwardLoadBalancerSet {
			forwardLoadBalancersMap := map[string]interface{}{}
			if forwardLoadBalancers.LoadBalancerId != nil {
				forwardLoadBalancersMap["load_balancer_id"] = forwardLoadBalancers.LoadBalancerId
			}

			if forwardLoadBalancers.ListenerId != nil {
				forwardLoadBalancersMap["listener_id"] = forwardLoadBalancers.ListenerId
			}

			if forwardLoadBalancers.TargetAttributes != nil {
				targetAttributesList := make([]map[string]interface{}, 0, len(forwardLoadBalancers.TargetAttributes))
				for _, targetAttributes := range forwardLoadBalancers.TargetAttributes {
					targetAttributesMap := map[string]interface{}{}
					if targetAttributes.Port != nil {
						targetAttributesMap["port"] = targetAttributes.Port
					}

					if targetAttributes.Weight != nil {
						targetAttributesMap["weight"] = targetAttributes.Weight
					}

					targetAttributesList = append(targetAttributesList, targetAttributesMap)
				}

				forwardLoadBalancersMap["target_attributes"] = targetAttributesList
			}

			if forwardLoadBalancers.LocationId != nil {
				forwardLoadBalancersMap["location_id"] = forwardLoadBalancers.LocationId
			}

			if forwardLoadBalancers.Region != nil {
				forwardLoadBalancersMap["region"] = forwardLoadBalancers.Region
			}

			forwardLoadBalancersList = append(forwardLoadBalancersList, forwardLoadBalancersMap)
		}

		_ = d.Set("forward_load_balancers", forwardLoadBalancersList)
	}

	return nil
}

func resourceTencentCloudAsLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_load_balancer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		request            = as.NewModifyLoadBalancerTargetAttributesRequest()
		autoScalingGroupId = d.Id()
	)

	immutableArgs := []string{"auto_scaling_group_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("forward_load_balancers") {
		request.AutoScalingGroupId = &autoScalingGroupId
		if v, ok := d.GetOk("forward_load_balancers"); ok {
			for _, item := range v.([]interface{}) {
				forwardLoadBalancer := as.ForwardLoadBalancer{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["load_balancer_id"]; ok {
					forwardLoadBalancer.LoadBalancerId = helper.String(v.(string))
				}

				if v, ok := dMap["listener_id"]; ok {
					forwardLoadBalancer.ListenerId = helper.String(v.(string))
				}

				if v, ok := dMap["target_attributes"]; ok {
					for _, item := range v.([]interface{}) {
						targetAttributesMap := item.(map[string]interface{})
						targetAttribute := as.TargetAttribute{}
						if v, ok := targetAttributesMap["port"]; ok {
							targetAttribute.Port = helper.IntUint64(v.(int))
						}

						if v, ok := targetAttributesMap["weight"]; ok {
							targetAttribute.Weight = helper.IntUint64(v.(int))
						}

						forwardLoadBalancer.TargetAttributes = append(forwardLoadBalancer.TargetAttributes, &targetAttribute)
					}
				}

				if v, ok := dMap["location_id"]; ok {
					forwardLoadBalancer.LocationId = helper.String(v.(string))
				}

				if v, ok := dMap["region"]; ok {
					forwardLoadBalancer.Region = helper.String(v.(string))
				}

				request.ForwardLoadBalancers = append(request.ForwardLoadBalancers, &forwardLoadBalancer)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyLoadBalancerTargetAttributes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update as loadBalancer failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAsLoadBalancerRead(d, meta)
}

func resourceTencentCloudAsLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_load_balancer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		autoScalingGroupId = d.Id()
	)

	if err := service.DeleteAsLoadBalancerById(ctx, autoScalingGroupId); err != nil {
		return err
	}

	return nil
}
