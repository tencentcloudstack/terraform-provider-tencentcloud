/*
Use this data source to query target group information.

Example Usage

```hcl
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test-target-keep-1"
}

resource "tencentcloud_clb_target_group_attachment" "group" {
    clb_id          = tencentcloud_clb_instance.clb_basic.id
    listener_id     = tencentcloud_clb_listener.listener_basic.listener_id
    rule_id         = tencentcloud_clb_listener_rule.rule_basic.rule_id
    targrt_group_id = tencentcloud_clb_target_group.test.id
}

data "tencentcloud_clb_target_groups" "target_group_info_id" {
  target_group_id = tencentcloud_clb_target_group.test.id
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbTargetGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbTargetGroupRead,

		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_id", "target_group_name"},
				AtLeastOneOf:  []string{"vpc_id", "target_group_name"},
				Description:   "ID of Target group. Mutually exclusive with `vpc_id` and `target_group_name`. `target_group_id` is preferred.",
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"target_group_id", "target_group_name"},
				Description:  "Target group VPC ID. Mutually exclusive with `target_group_id`. `target_group_id` is preferred.",
			},
			"target_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"target_group_id", "vpc_id"},
				Description:  "Name of target group. Mutually exclusive with `target_group_id`. `target_group_id` is preferred.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Target group info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of Target group.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of target group.",
						},
						"target_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group VPC ID.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of target group.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the target group.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time of the target group.",
						},
						"associated_rule_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of associated rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Load balance ID.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener ID.",
									},
									"location_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Forwarding rule ID.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener protocol type.",
									},
									"listener_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Listener port.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Forwarding rule domain.",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Forwarding rule URL.",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Load balance name.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener name.",
									},
								},
							},
						},
						"target_group_instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of backend servers bound to the target group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of backend service.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of backend service.",
									},
									"server_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port of backend service.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Forwarding weight of back-end services.",
									},
									"public_ip_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        schema.TypeString,
										Description: "List of external network IP of back-end services.",
									},
									"private_ip_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        schema.TypeString,
										Description: "Intranet IP list of back-end services.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance name of the backend service.",
									},
									"registered_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time the backend service was bound.",
									},
									"eni_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of Elastic Network Interface.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_target_groups.read")()

	var (
		clbService = ClbService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		logId                = getLogId(contextNil)
		ctx                  = context.WithValue(context.TODO(), logIdKey, logId)
		instances            []*clb.TargetGroupBackend
		targetInfos          []*clb.TargetGroupInfo
		filters              = make(map[string]string, 2)
		targetGroupInstances []map[string]interface{}
		targetGroupId        string
		err                  error
	)

	if id, ok := d.GetOk("target_group_id"); ok {
		targetGroupId = id.(string)
	}
	if id, ok := d.GetOk("vpc_id"); ok {
		filters["TargetGroupVpcId"] = id.(string)
	}
	if name, ok := d.GetOk("target_group_name"); ok {
		filters["TargetGroupName"] = name.(string)
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		targetInfos, err = clbService.DescribeTargetGroups(ctx, targetGroupId, filters)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	var (
		list    = make([]map[string]interface{}, 0, len(targetInfos))
		ids     = make([]string, 0, len(targetInfos))
		isExist = make(map[string]bool)
	)

	for _, info := range targetInfos {
		targetId := *info.TargetGroupId
		ids = append(ids, targetId)
		if _, ok := isExist[targetId]; !ok {
			instances = []*clb.TargetGroupBackend{}
			targetGroupInstances = []map[string]interface{}{}

			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instances, err = clbService.DescribeTargetGroupInstances(ctx, map[string]string{
					"TargetGroupId": *info.TargetGroupId,
				})
				if err != nil {
					return retryError(err, InternalError)
				}
				return nil
			})
			if err != nil {
				return err
			}

			isExist[targetId] = true
			for _, instance := range instances {
				targetGroupInstances = append(targetGroupInstances, map[string]interface{}{
					"server_type":          instance.Type,
					"instance_id":          instance.InstanceId,
					"server_port":          instance.Port,
					"weight":               instance.Weight,
					"public_ip_addresses":  instance.PublicIpAddresses,
					"private_ip_addresses": instance.PrivateIpAddresses,
					"instance_name":        instance.InstanceName,
					"registered_time":      instance.RegisteredTime,
					"eni_id":               instance.EniId,
				})
			}
		}

		ruleInfo := make([]map[string]interface{}, 0, len(info.AssociatedRule))
		for _, rule := range info.AssociatedRule {
			ruleInfo = append(ruleInfo, map[string]interface{}{
				"load_balancer_id":   rule.LoadBalancerId,
				"listener_id":        rule.ListenerId,
				"location_id":        rule.LocationId,
				"protocol":           rule.Protocol,
				"listener_port":      rule.Port,
				"domain":             rule.Domain,
				"url":                rule.Url,
				"load_balancer_name": rule.LoadBalancerName,
				"listener_name":      rule.ListenerName,
			})
		}

		list = append(list, map[string]interface{}{
			"target_group_id":            targetId,
			"vpc_id":                     info.VpcId,
			"target_group_name":          info.TargetGroupName,
			"port":                       info.Port,
			"create_time":                info.CreatedTime,
			"update_time":                info.UpdatedTime,
			"associated_rule_list":       ruleInfo,
			"target_group_instance_list": targetGroupInstances,
		})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set target group list fail, reason:%s ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), list); err != nil {
			return err
		}
	}

	return nil
}
