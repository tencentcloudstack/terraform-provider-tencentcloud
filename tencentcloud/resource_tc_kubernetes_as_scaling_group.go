/*
Provide a resource to create an auto scaling group for kubernetes cluster.

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
    instance_type      = "SN3ne.8XLARGE64"
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
}
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"strings"
)

func ResourceTencentCloudKubernetesAsScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesAsScalingGroupCreate,
		Read:   resourceKubernetesAsScalingGroupRead,
		Update: func(d *schema.ResourceData, meta interface{}) error {
			d.Partial(true)
			return fmt.Errorf("resource tencentcloud_kubernetes_as_scaling_group not support update")
		},
		Delete: resourceKubernetesAsScalingGroupDelete,
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
				Description: "Auto scaling group parameters",
			},
			"auto_scaling_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: kubernetesAsScalingConfigPara(),
				},
				Description: "Auto scaling config parameters",
			},
		},
	}
}

func kubernetesAsScalingGroupPara() map[string]*schema.Schema {
	removes := map[string]bool{
		"status":           true,
		"instance_count":   true,
		"create_time":      true,
		"configuration_id": true,
	}

	asGroupSchema := resourceTencentCloudAsScalingGroup().Schema
	needSchema := make(map[string]*schema.Schema, len(asGroupSchema)-len(removes))
	for name, schemaDef := range asGroupSchema {
		if removes[name] {
			continue
		}
		var newSchema = *schemaDef
		if len(newSchema.ConflictsWith) != 0 {
			newSchema.ConflictsWith = make([]string, len(newSchema.ConflictsWith))
			for index, cft := range schemaDef.ConflictsWith {
				newSchema.ConflictsWith[index] = "auto_scaling_group." + cft
			}
		}
		needSchema[name] = &newSchema
	}

	return needSchema
}

func kubernetesAsScalingGroupParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result   string
		errRet   error
		requires = []string{
			"scaling_group_name",
			"max_size",
			"min_size",
			"vpc_id",
		}
	)

	for _, require := range requires {
		if _, ok := dMap[require]; !ok {
			return "", fmt.Errorf("miss require param %s", require)
		}
	}

	request := as.NewCreateAutoScalingGroupRequest()
	request.AutoScalingGroupName = stringToPointer(dMap["scaling_group_name"].(string))
	request.MaxSize = intToPointer(dMap["max_size"].(int))
	request.MinSize = intToPointer(dMap["min_size"].(int))

	if *request.MinSize > *request.MaxSize {
		return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
	}

	request.VpcId = stringToPointer(dMap["vpc_id"].(string))

	if v, ok := dMap["default_cooldown"]; ok {
		request.DefaultCooldown = intToPointer(v.(int))
	}

	if v, ok := dMap["desired_capacity"]; ok {
		request.DesiredCapacity = intToPointer(v.(int))
		if *request.DesiredCapacity > *request.MaxSize ||
			*request.DesiredCapacity < *request.MinSize {
			return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
		}

	}

	if v, ok := dMap["retry_policy"]; ok {
		request.RetryPolicy = stringToPointer(v.(string))
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
				LoadBalancerId: stringToPointer(vv["load_balancer_id"].(string)),
				ListenerId:     stringToPointer(vv["listener_id"].(string)),
				LocationId:     stringToPointer(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   intToPointer(t["port"].(int)),
					Weight: intToPointer(t["weight"].(int)),
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
				Key:   stringToPointer(k),
				Value: stringToPointer(v.(string)),
			})
		}
	}
	result = request.ToJsonString()

	return result, errRet
}

func kubernetesAsScalingConfigPara() map[string]*schema.Schema {

	removes := map[string]bool{
		"image_id":         true,
		"keep_image_login": true,
		"user_data":        true,
		"status":           true,
		"create_time":      true,
		"instance_types":   true,
	}
	asConfigSchema := resourceTencentCloudAsScalingConfig().Schema
	needSchema := make(map[string]*schema.Schema, len(asConfigSchema)-len(removes))

	for name, schemaDef := range asConfigSchema {
		if removes[name] {
			continue
		}
		var newSchema = *schemaDef
		if len(newSchema.ConflictsWith) != 0 {
			newSchema.ConflictsWith = make([]string, len(newSchema.ConflictsWith))
			for index, cft := range schemaDef.ConflictsWith {
				newSchema.ConflictsWith[index] = "auto_scaling_config." + cft
			}
		}
		needSchema[name] = &newSchema
	}

	needSchema["key_ids"].ConflictsWith = []string{"auto_scaling_config.password"}
	needSchema["password"].ConflictsWith = []string{"auto_scaling_config.key_ids"}
	needSchema["instance_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specified types of CVM instance.",
	}

	return needSchema
}

func kubernetesAsScalingConfigParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result   string
		errRet   error
		requires = []string{
			"configuration_name",
		}
	)

	for _, require := range requires {
		if _, ok := dMap[require]; !ok {
			return "", fmt.Errorf("miss require param %s", require)
		}
	}

	request := as.NewCreateLaunchConfigurationRequest()

	request.LaunchConfigurationName = stringToPointer(dMap["configuration_name"].(string))

	if v, ok := dMap["project_id"]; ok {
		request.ProjectId = intToPointer(v.(int))
	}

	instanceType := dMap["instance_type"].(string)
	request.InstanceType = &instanceType

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = stringToPointer(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = intToPointer(v.(int))
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
		request.InternetAccessible.InternetChargeType = stringToPointer(v.(string))
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = intToPointer(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}

	if v, ok := dMap["password"]; ok {
		request.LoginSettings.Password = stringToPointer(v.(string))
	}
	if v, ok := dMap["key_ids"]; ok {
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
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
		request.InstanceTypesCheckPolicy = stringToPointer(v.(string))
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
	defer logElapsed("resource.resource_tc_kubernetes_as_scaling_group.read")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), "logId", logId)
		service = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		items   = strings.Split(d.Id(), ":")
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_as_scaling_group id  is broken")
	}
	clusterId := items[0]
	asGroupId := items[1]

	info, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = service.DescribeCluster(ctx, clusterId)
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
		scalingGroup *as.AutoScalingGroup
		number       int
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		scalingGroup, number, err = asService.DescribeAutoScalingGroupById(ctx, asGroupId)
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
	return nil
}

func resourceKubernetesAsScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.resource_tc_kubernetes_as_scaling_group.create")()
	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), "logId", logId)
		clusterId   = d.Get("cluster_id").(string)
		groupParas  = d.Get("auto_scaling_group").([]interface{})
		configParas = d.Get("auto_scaling_config").([]interface{})
		asService   = AsService{client: meta.(*TencentCloudClient).apiV3Conn}
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

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	asGroupId, err := service.CreateClusterAsGroup(ctx, clusterId, groupParaStr, configParaStr)
	if err != nil {
		return err
	}

	d.SetId(clusterId + ":" + asGroupId)

	// wait for status ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, asGroupId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
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

func resourceKubernetesAsScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_as_scaling_group.delete")()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), "logId", logId)
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
			return retryError(errRet, "InternalError")
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
			return retryError(err, "InternalError")
		}
		return nil
	})

	if hasDelete {
		return nil
	}

	// wait for delete ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		_, has, errRet := asService.DescribeAutoScalingGroupById(ctx, asGroupId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		if has > 0 {
			resource.RetryableError(fmt.Errorf("as group %s still alive", asGroupId))
		}
		return nil
	})

	return err
}
