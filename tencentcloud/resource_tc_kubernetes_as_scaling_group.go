/*
Provide a resource to create an auto scaling group for kubernetes cluster.

~> **NOTE:**  It has been deprecated and replaced by `tencentcloud_cluster_node_pool`.
~> **NOTE:** To use the custom Kubernetes component startup parameter function (parameter `extra_args`), you need to submit a ticket for application.

Example Usage

```hcl

resource "tencentcloud_kubernetes_as_scaling_group" "test" {

  cluster_id = "cls-kb32pbv4"

  auto_scaling_group {
    scaling_group_name   = "tf-guagua-as-group"
    max_size             = "5"
    min_size             = "0"
    vpc_id               = "vpc-dk8zmwuf"
    subnet_ids           = ["subnet-pqfek0t8"]
    project_id           = 0
    default_cooldown     = 400
    desired_capacity     = "0"
    termination_policies = ["NEWEST_INSTANCE"]
    retry_policy         = "INCREMENTAL_INTERVALS"

    tags = {
      "test" = "test"
    }

  }


  auto_scaling_config {
    configuration_name = "tf-guagua-as-config"
    instance_type      = "S1.SMALL1"
    project_id         = 0
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

    instance_tags = {
      tag = "as"
    }

  }

  labels = {
    "test1" = "test1",
    "test1" = "test2",
  }
}
```

Use Kubelet

```hcl

resource "tencentcloud_kubernetes_as_scaling_group" "test" {

  cluster_id = "cls-kb32pbv4"

  auto_scaling_group {
    scaling_group_name   = "tf-guagua-as-group"
    max_size             = "5"
    min_size             = "0"
    vpc_id               = "vpc-dk8zmwuf"
    subnet_ids           = ["subnet-pqfek0t8"]
    project_id           = 0
    default_cooldown     = 400
    desired_capacity     = "0"
    termination_policies = ["NEWEST_INSTANCE"]
    retry_policy         = "INCREMENTAL_INTERVALS"

    tags = {
      "test" = "test"
    }

  }


  auto_scaling_config {
    configuration_name = "tf-guagua-as-config"
    instance_type      = "S1.SMALL1"
    project_id         = 0
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

    instance_tags = {
      tag = "as"
    }

  }

  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]

  labels = {
    "test1" = "test1",
    "test1" = "test2",
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesAsScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesAsScalingGroupCreate,
		Read:   resourceKubernetesAsScalingGroupRead,
		Delete: resourceKubernetesAsScalingGroupDelete,
		Update: resourceKubernetesAsScalingGroupUpdate,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the cluster.",
			},
			"auto_scaling_group": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: kubernetesAsScalingGroupPara(),
				},
				Description: "Auto scaling group parameters.",
			},
			"auto_scaling_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: kubernetesAsScalingConfigPara(),
				},
				Description: "Auto scaling config parameters.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels of kubernetes AS Group created nodes.",
			},
			"extra_args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom parameter information related to the node.",
			},
			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
			},
		},
	}
}

func kubernetesAsScalingConfigPara() map[string]*schema.Schema {
	needSchema := map[string]*schema.Schema{
		"configuration_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateStringLengthInRange(1, 60),
			Description:  "Name of a launch configuration.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Specifys to which project the configuration belongs.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specified types of CVM instance.",
		},
		"system_disk_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
			Description:  "Type of a CVM disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(50, 500),
			Description:  "Volume of system disk in GB. Default is `50`.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    11,
			Description: "Configurations of data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
						ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
						Description:  "Types of disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						Optional:    true,
						ForceNew:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Data disk snapshot ID.",
					},
				},
			},
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
			ValidateFunc: validateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
			Description:  "Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
		},
		"internet_max_bandwidth_out": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Max bandwidth of Internet access in Mbps. Default is `0`.",
		},
		"public_ip_assigned": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Description: "Specify whether to assign an Internet IP address.",
		},
		"password": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ForceNew:      true,
			ValidateFunc:  validateAsConfigPassword,
			ConflictsWith: []string{"auto_scaling_config.0.key_ids"},
			Description:   "Password to access.",
		},
		"key_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"auto_scaling_config.0.password"},
			Description:   "ID list of keys.",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Security groups to which a CVM instance belongs.",
		},
		"enhanced_security_service": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "To specify whether to enable cloud security service. Default is TRUE.",
		},
		"enhanced_monitor_service": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
		},
		"instance_tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "A list of tags used to associate different resources.",
		},
	}

	return needSchema
}

func kubernetesAsScalingGroupPara() map[string]*schema.Schema {

	asGroupSchema := map[string]*schema.Schema{
		"scaling_group_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateStringLengthInRange(1, 55),
			Description:  "Name of a scaling group.",
		},

		"max_size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validateIntegerInRange(0, 2000),
			Description:  "Maximum number of CVM instances (0~2000).",
		},
		"min_size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validateIntegerInRange(0, 2000),
			Description:  "Minimum number of CVM instances (0~2000).",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of VPC network.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Specifys to which project the scaling group belongs.",
		},
		"subnet_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "ID list of subnet, and for VPC it is required.",
		},
		"zones": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of available zones, for Basic network it is required.",
		},
		"default_cooldown": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     300,
			Description: "Default cooldown time in second, and default value is 300.",
		},
		"desired_capacity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Desired volume of CVM instances, which is between max_size and min_size.",
		},
		"load_balancer_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"auto_scaling_group.0.forward_balancer_ids"},
			Description:   "ID list of traditional load balancers.",
		},
		"forward_balancer_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"auto_scaling_group.0.load_balancer_ids"},
			Description:   "List of application load balancers, which can't be specified with load_balancer_ids together.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"load_balancer_id": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "ID of available load balancers.",
					},
					"listener_id": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "Listener ID for application load balancers.",
					},
					"rule_id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "ID of forwarding rules.",
					},
					"target_attribute": {
						Type:        schema.TypeList,
						Required:    true,
						ForceNew:    true,
						Description: "Attribute list of target rules.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"port": {
									Type:        schema.TypeInt,
									Required:    true,
									ForceNew:    true,
									Description: "Port number.",
								},
								"weight": {
									Type:        schema.TypeInt,
									Required:    true,
									ForceNew:    true,
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
			ForceNew:    true,
			MaxItems:    1,
			Description: "Available values for termination policies include `OLDEST_INSTANCE` and `NEWEST_INSTANCE`.",
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
			ForceNew:    true,
			Description: "Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.",
			Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
			ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS}),
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Tags of a scaling group.",
		},
	}

	return asGroupSchema
}

func kubernetesAsScalingGroupParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateAutoScalingGroupRequest()
	request.AutoScalingGroupName = helper.String(dMap["scaling_group_name"].(string))
	request.MaxSize = helper.IntUint64(dMap["max_size"].(int))
	request.MinSize = helper.IntUint64(dMap["min_size"].(int))

	if *request.MinSize > *request.MaxSize {
		return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
	}

	request.VpcId = helper.String(dMap["vpc_id"].(string))

	if v, ok := dMap["default_cooldown"]; ok {
		request.DefaultCooldown = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["desired_capacity"]; ok {
		request.DesiredCapacity = helper.IntUint64(v.(int))
		if *request.DesiredCapacity > *request.MaxSize ||
			*request.DesiredCapacity < *request.MinSize {
			return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
		}

	}

	if v, ok := dMap["retry_policy"]; ok {
		request.RetryPolicy = helper.String(v.(string))
	}

	if v, ok := dMap["subnet_ids"]; ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}

	if v, ok := dMap["zones"]; ok {
		zones := v.([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}

	if v, ok := dMap["load_balancer_ids"]; ok {
		loadBalancerIds := v.([]interface{})
		request.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerId)
		}
	}

	if v, ok := dMap["forward_balancer_ids"]; ok {
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

	if v, ok := dMap["termination_policies"]; ok {
		terminationPolicies := v.([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if v, ok := dMap["tags"]; ok {
		for k, v := range v.(map[string]interface{}) {
			request.Tags = append(request.Tags, &as.Tag{
				Key:   helper.String(k),
				Value: helper.String(v.(string)),
			})
		}
	}
	result = request.ToJsonString()

	return result, errRet
}

func kubernetesAsScalingConfigParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateLaunchConfigurationRequest()

	request.LaunchConfigurationName = helper.String(dMap["configuration_name"].(string))

	if v, ok := dMap["project_id"]; ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	instanceType := dMap["instance_type"].(string)
	request.InstanceType = &instanceType

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
				DiskSize: &diskSize,
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := dMap["internet_charge_type"]; ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}

	if v, ok := dMap["password"]; ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := dMap["key_ids"]; ok {
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}

	if request.LoginSettings.Password != nil && *request.LoginSettings.Password == "" {
		request.LoginSettings.Password = nil
	}

	if request.LoginSettings.Password == nil && len(request.LoginSettings.KeyIds) == 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` should be set one")
		return result, errRet
	}

	if request.LoginSettings.Password != nil && len(request.LoginSettings.KeyIds) != 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` can only be supported one")
		return result, errRet
	}

	if v, ok := dMap["security_group_ids"]; ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := dMap["enhanced_security_service"]; ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	chargeType := INSTANCE_CHARGE_TYPE_POSTPAID
	request.InstanceChargeType = &chargeType

	if v, ok := dMap["instance_types_check_policy"]; ok {
		request.InstanceTypesCheckPolicy = helper.String(v.(string))
	}

	if v, ok := dMap["instance_tags"]; ok {
		tags := v.(map[string]interface{})
		request.InstanceTags = make([]*as.InstanceTag, 0, len(tags))
		for k, t := range tags {
			key := k
			value := t.(string)
			tag := as.InstanceTag{
				Key:   &key,
				Value: &value,
			}
			request.InstanceTags = append(request.InstanceTags, &tag)
		}
	}
	result = request.ToJsonString()
	return result, errRet
}

func resourceKubernetesAsScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_as_scaling_group.read")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		items   = strings.Split(d.Id(), ":")
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_as_scaling_group id  is broken")
	}
	clusterId := items[0]
	asGroupId := items[1]

	_, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		d.SetId("")
		return nil
	}

	var (
		asService = AsService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		number int
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, number, err = asService.DescribeAutoScalingGroupById(ctx, asGroupId)
		if err != nil {
			return retryError(err)
		}
		if number == 0 {
			d.SetId("")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if number == 0 {
		return nil
	}

	var clusterAsGroupSet *tke.ClusterAsGroup
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		clusterAsGroupSet, err = service.DescribeClusterAsGroupsByGroupId(ctx, clusterId, asGroupId)
		if err != nil {
			return retryError(err)
		}

		if clusterAsGroupSet == nil {
			return nil
		}

		labels := clusterAsGroupSet.Labels
		var labelsMap = make(map[string]string, len(labels))

		for _, v := range labels {
			labelsMap[*v.Name] = *v.Value
		}
		d.Set("labels", labelsMap)
		return nil
	})

	if err != nil {
		return err
	}

	if clusterAsGroupSet == nil {
		d.SetId("")
	}
	return nil
}

func resourceKubernetesAsScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_kubernetes_as_scaling_group.create")()
	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		clusterId   = d.Get("cluster_id").(string)
		groupParas  = d.Get("auto_scaling_group").([]interface{})
		configParas = d.Get("auto_scaling_config").([]interface{})
		asService   = AsService{client: meta.(*TencentCloudClient).apiV3Conn}
		iAdvanced   InstanceAdvancedSettings
	)
	if len(groupParas) != 1 || len(configParas) != 1 {
		return fmt.Errorf("need only one auto_scaling_group and one auto_scaling_config")
	}

	groupParaStr, err := kubernetesAsScalingGroupParaSerial(groupParas[0].(map[string]interface{}), meta)
	if err != nil {
		return err
	}

	configParaStr, err := kubernetesAsScalingConfigParaSerial(configParas[0].(map[string]interface{}), meta)
	if err != nil {
		return err
	}

	labels := GetTkeLabels(d, "labels")
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = int64(temp.(int))
	}
	if temp, ok := d.GetOk("extra_args"); ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		for _, extraArg := range extraArgs {
			iAdvanced.ExtraArgs.Kubelet = append(iAdvanced.ExtraArgs.Kubelet, &extraArg)
		}
	}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	asGroupId, err := service.CreateClusterAsGroup(ctx, clusterId, groupParaStr, configParaStr, labels, iAdvanced)
	if err != nil {
		return err
	}

	d.SetId(clusterId + ":" + asGroupId)

	// wait for status ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, asGroupId)
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

	return nil
}

func resourceKubernetesAsScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_as_scaling_group.update")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		items   = strings.Split(d.Id(), ":")
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_as_scaling_group id  is broken")
	}
	clusterId := items[0]
	asGroupId := items[1]
	if d.HasChange("auto_scaling_group") {
		asGroups := d.Get("auto_scaling_group").([]interface{})
		for _, d := range asGroups {
			value := d.(map[string]interface{})
			maxSize := int64(value["max_size"].(int))
			minSize := int64(value["min_size"].(int))
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				errRet := service.ModifyClusterAsGroupAttribute(ctx, clusterId, asGroupId, maxSize, minSize)
				if errRet != nil {
					return retryError(errRet)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return resourceKubernetesAsScalingGroupRead(d, meta)
}

func resourceKubernetesAsScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_as_scaling_group.delete")()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		asService = AsService{client: meta.(*TencentCloudClient).apiV3Conn}
		items     = strings.Split(d.Id(), ":")
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_as_scaling_group id  is broken")
	}
	clusterId := items[0]
	asGroupId := items[1]

	//set len(cvm)==0
	scalingGroup, has, err := asService.DescribeAutoScalingGroupById(ctx, asGroupId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	if *scalingGroup.InstanceCount > 0 || *scalingGroup.DesiredCapacity > 0 {
		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			return retryError(asService.ClearScalingGroupInstance(ctx, asGroupId))
		}); err != nil {
			return err
		}
	}

	// wait  set finish
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, asGroupId)
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

	//delete as group
	hasDelete := false
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := service.DeleteClusterAsGroups(ctx, clusterId, asGroupId)

		if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkErr.Code == "InternalError.Param" && strings.Contains(sdkErr.Message, "Not Found") {
				hasDelete = true
				return nil
			}
		}
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if hasDelete {
		return nil
	}

	// wait for delete ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		_, has, errRet := asService.DescribeAutoScalingGroupById(ctx, asGroupId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if has > 0 {
			resource.RetryableError(fmt.Errorf("as group %s still alive", asGroupId))
		}
		return nil
	})

	return err
}
