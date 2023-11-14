/*
Provides a resource to create a as load_balancer

Example Usage

```hcl
resource "tencentcloud_as_load_balancer" "load_balancer" {
  auto_scaling_group_id = "asg-12wjuh0s"
  load_balancer_ids =
  forward_load_balancers {
		load_balancer_id = "lb-d8u76te5"
		listener_id = "lbl-s8dh4y75"
		target_attributes {
			port = 8080
			weight = 20
		}
		location_id = "loc-fsa87u6d"
		region = "ap-guangzhou"

  }
}
```

Import

as load_balancer can be imported using the id, e.g.

```
terraform import tencentcloud_as_load_balancer.load_balancer load_balancer_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudAsLoadBalancer() *schema.Resource {
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

			"load_balancer_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of traditional load balancer IDs. The maximum number of traditional load balancers bound to each scaling group is 20. Both LoadBalancerIds and ForwardLoadBalancers can specify at most one at the same time.",
			},

			"forward_load_balancers": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of application load balancers. The maximum number of application-type load balancers bound to each scaling group is 100. Both LoadBalancerIds and ForwardLoadBalancers can specify at most one at the same time.",
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
							Description: "Load balancer instance region. Default value is the region of current auto scaling group. The format is the same as the public parameter Region, for example: ap-guangzhou.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAsLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_load_balancer.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = as.NewAttachLoadBalancersRequest()
		response           = as.NewAttachLoadBalancersResponse()
		autoScalingGroupId string
	)
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		autoScalingGroupId = v.(string)
		request.AutoScalingGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIdsSet := v.(*schema.Set).List()
		for i := range loadBalancerIdsSet {
			loadBalancerIds := loadBalancerIdsSet[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerIds)
		}
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().AttachLoadBalancers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create as loadBalancer failed, reason:%+v", logId, err)
		return err
	}

	autoScalingGroupId = *response.Response.AutoScalingGroupId
	d.SetId(autoScalingGroupId)

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESSFUL"}, 2*readRetryTimeout, time.Second, service.AsLoadBalancerStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudAsLoadBalancerRead(d, meta)
}

func resourceTencentCloudAsLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_load_balancer.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	loadBalancerId := d.Id()

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

	if loadBalancer.LoadBalancerIds != nil {
		_ = d.Set("load_balancer_ids", loadBalancer.LoadBalancerIds)
	}

	if loadBalancer.ForwardLoadBalancers != nil {
		forwardLoadBalancersList := []interface{}{}
		for _, forwardLoadBalancers := range loadBalancer.ForwardLoadBalancers {
			forwardLoadBalancersMap := map[string]interface{}{}

			if loadBalancer.ForwardLoadBalancers.LoadBalancerId != nil {
				forwardLoadBalancersMap["load_balancer_id"] = loadBalancer.ForwardLoadBalancers.LoadBalancerId
			}

			if loadBalancer.ForwardLoadBalancers.ListenerId != nil {
				forwardLoadBalancersMap["listener_id"] = loadBalancer.ForwardLoadBalancers.ListenerId
			}

			if loadBalancer.ForwardLoadBalancers.TargetAttributes != nil {
				targetAttributesList := []interface{}{}
				for _, targetAttributes := range loadBalancer.ForwardLoadBalancers.TargetAttributes {
					targetAttributesMap := map[string]interface{}{}

					if targetAttributes.Port != nil {
						targetAttributesMap["port"] = targetAttributes.Port
					}

					if targetAttributes.Weight != nil {
						targetAttributesMap["weight"] = targetAttributes.Weight
					}

					targetAttributesList = append(targetAttributesList, targetAttributesMap)
				}

				forwardLoadBalancersMap["target_attributes"] = []interface{}{targetAttributesList}
			}

			if loadBalancer.ForwardLoadBalancers.LocationId != nil {
				forwardLoadBalancersMap["location_id"] = loadBalancer.ForwardLoadBalancers.LocationId
			}

			if loadBalancer.ForwardLoadBalancers.Region != nil {
				forwardLoadBalancersMap["region"] = loadBalancer.ForwardLoadBalancers.Region
			}

			forwardLoadBalancersList = append(forwardLoadBalancersList, forwardLoadBalancersMap)
		}

		_ = d.Set("forward_load_balancers", forwardLoadBalancersList)

	}

	return nil
}

func resourceTencentCloudAsLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_load_balancer.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := as.NewModifyLoadBalancerTargetAttributesRequest()

	loadBalancerId := d.Id()

	request.AutoScalingGroupId = &autoScalingGroupId

	immutableArgs := []string{"auto_scaling_group_id", "load_balancer_ids", "forward_load_balancers"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("auto_scaling_group_id") {
		if v, ok := d.GetOk("auto_scaling_group_id"); ok {
			request.AutoScalingGroupId = helper.String(v.(string))
		}
	}

	if d.HasChange("forward_load_balancers") {
		if v, ok := d.GetOk("forward_load_balancers"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().ModifyLoadBalancerTargetAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update as loadBalancer failed, reason:%+v", logId, err)
		return err
	}

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESSFUL"}, 2*readRetryTimeout, time.Second, service.AsLoadBalancerStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudAsLoadBalancerRead(d, meta)
}

func resourceTencentCloudAsLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_load_balancer.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}
	loadBalancerId := d.Id()

	if err := service.DeleteAsLoadBalancerById(ctx, autoScalingGroupId); err != nil {
		return err
	}

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESSFUL"}, 2*readRetryTimeout, time.Second, service.AsLoadBalancerStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
