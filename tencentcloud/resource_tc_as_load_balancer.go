/*
Provides a resource to create a as load_balancer

~> **NOTE:** `load_balancer_ids` A list of traditional load balancer IDs, with a maximum of 20 traditional load balancers bound to each scaling group. Only one LoadBalancerIds and ForwardLoadBalancers can be specified simultaneously.
~> **NOTE:** `forward_load_balancers` List of application type load balancers, with a maximum of 100 bound application type load balancers for each scaling group. Only one LoadBalancerIds and ForwardLoadBalancers can be specified simultaneously.

Example Usage

If use `load_balancer_ids`

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "clb-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_clb_listener" "example" {
  clb_id        = tencentcloud_clb_instance.example.id
  listener_name = "listener-example"
  port          = 80
  protocol      = "HTTP"
}

resource "tencentcloud_clb_listener_rule" "example" {
  listener_id = tencentcloud_clb_listener.example.listener_id
  clb_id      = tencentcloud_clb_instance.example.id
  domain      = "foo.net"
  url         = "/bar"
}

resource "tencentcloud_as_load_balancer" "example" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.example.id
  load_balancer_ids     = [tencentcloud_clb_instance.example.id]
}
```

If use `forward_load_balancers`

```hcl
resource "tencentcloud_as_load_balancer" "example" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.example.id

  forward_load_balancers {
    load_balancer_id = tencentcloud_clb_instance.example.id
    listener_id      = tencentcloud_clb_listener.example.listener_id
    location_id      = tencentcloud_clb_listener_rule.example.rule_id

    target_attributes {
      port   = 8080
      weight = 20
    }
  }
}
```

Import

as load_balancer can be imported using the id, e.g.

```
terraform import tencentcloud_as_load_balancer.load_balancer auto_scaling_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
	defer logElapsed("resource.tencentcloud_as_load_balancer.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	autoScalingGroupId := d.Id()

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

	if loadBalancer.LoadBalancerIdSet != nil {
		_ = d.Set("load_balancer_ids", loadBalancer.LoadBalancerIdSet)
	}

	if loadBalancer.ForwardLoadBalancerSet != nil {
		forwardLoadBalancersList := []interface{}{}
		for _, forwardLoadBalancers := range loadBalancer.ForwardLoadBalancerSet {
			forwardLoadBalancersMap := map[string]interface{}{}

			if forwardLoadBalancers.LoadBalancerId != nil {
				forwardLoadBalancersMap["load_balancer_id"] = forwardLoadBalancers.LoadBalancerId
			}

			if forwardLoadBalancers.ListenerId != nil {
				forwardLoadBalancersMap["listener_id"] = forwardLoadBalancers.ListenerId
			}

			if forwardLoadBalancers.TargetAttributes != nil {
				targetAttributesList := []interface{}{}
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

				forwardLoadBalancersMap["target_attributes"] = []interface{}{targetAttributesList}
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
	defer logElapsed("resource.tencentcloud_as_load_balancer.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := as.NewModifyLoadBalancerTargetAttributesRequest()

	autoScalingGroupId := d.Id()

	request.AutoScalingGroupId = &autoScalingGroupId

	immutableArgs := []string{"load_balancer_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("forward_load_balancers") {
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

	return resourceTencentCloudAsLoadBalancerRead(d, meta)
}

func resourceTencentCloudAsLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_load_balancer.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}
	autoScalingGroupId := d.Id()

	if err := service.DeleteAsLoadBalancerById(ctx, autoScalingGroupId); err != nil {
		return err
	}

	return nil
}
