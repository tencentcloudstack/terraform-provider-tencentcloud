/*
Provides a resource to create a tdmq rabbitmq_vip_instance

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "rabbitmq-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "rabbitmq-subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRabbitmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRabbitmqVipInstanceCreate,
		Read:   resourceTencentCloudTdmqRabbitmqVipInstanceRead,
		Update: resourceTencentCloudTdmqRabbitmqVipInstanceUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqVipInstanceDelete,

		Schema: map[string]*schema.Schema{
			"zone_ids": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "availability zone.",
			},
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Private network VpcId.",
			},
			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Private network SubnetId.",
			},
			"cluster_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "cluster name.",
			},
			"node_spec": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Node specifications, basic type rabbit-vip-basic-1, standard type rabbit-vip-basic-2, high-level type 1 rabbit-vip-basic-3, high-level type 2 rabbit-vip-basic-4. If not passed, the default is the basic type.",
			},
			"node_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of nodes, a minimum of 3 nodes for a multi-availability zone. If not passed, the default single availability zone is 1, and the multi-availability zone is 3.",
			},
			"storage_size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Single node storage specification, the default is 200G.",
			},
			"enable_create_default_ha_mirror_queue": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Mirrored queue, the default is false.",
			},
			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Automatic renewal, the default is true.",
			},
			"time_span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase duration, the default is 1 (month).",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqVipInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = tdmq.NewCreateRabbitMQVipInstanceRequest()
		response   = tdmq.NewCreateRabbitMQVipInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.(*schema.Set).List()
		for i := range zoneIdsSet {
			zoneIds := zoneIdsSet[i].(int)
			request.ZoneIds = append(request.ZoneIds, helper.IntInt64(zoneIds))
		}
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_spec"); ok {
		request.NodeSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("node_num"); ok {
		request.NodeNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_create_default_ha_mirror_queue"); ok {
		request.EnableCreateDefaultHaMirrorQueue = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRabbitMQVipInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId

	// wait
	paramMap := make(map[string]interface{})
	tmpSet := make([]*tdmq.Filter, 0)
	filter := tdmq.Filter{}
	filter.Name = helper.String("instanceIds")
	filter.Values = helper.Strings([]string{instanceId})
	tmpSet = append(tmpSet, &filter)
	paramMap["filters"] = tmpSet
	err = resource.Retry(readRetryTimeout*10, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("resource `tencentcloud_tdmq_rabbitmq_vip_instance` %s does not exist", instanceId))
		}

		if len(result) != 1 {
			return resource.NonRetryableError(fmt.Errorf("resource `tencentcloud_tdmq_rabbitmq_vip_instance` %s id error", instanceId))
		}

		if *result[0].Status == RabbitMQVipInstanceRunning {
			return resource.RetryableError(fmt.Errorf("rabbitmq_vip_instance status is running"))
		} else if *result[0].Status == RabbitMQVipInstanceSuccess {
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("rabbitmq_vip_instance status illegal"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	rabbitmqVipInstance, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if rabbitmqVipInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRabbitmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rabbitmqVipInstance.ClusterInfo.ZoneIds != nil {
		_ = d.Set("zone_ids", rabbitmqVipInstance.ClusterInfo.ZoneIds)
	}

	if rabbitmqVipInstance.ClusterInfo.Vpcs != nil {
		_ = d.Set("vpc_id", rabbitmqVipInstance.ClusterInfo.Vpcs[0].VpcId)
		_ = d.Set("subnet_id", rabbitmqVipInstance.ClusterInfo.Vpcs[0].SubnetId)
	}

	if rabbitmqVipInstance.ClusterSpecInfo.NodeCount != nil {
		_ = d.Set("node_num", rabbitmqVipInstance.ClusterSpecInfo.NodeCount)
	}

	if rabbitmqVipInstance.ClusterSpecInfo.MaxStorage != nil {
		_ = d.Set("storage_size", rabbitmqVipInstance.ClusterSpecInfo.MaxStorage)
	}

	paramMap := make(map[string]interface{})
	tmpSet := make([]*tdmq.Filter, 0)
	filter := tdmq.Filter{}
	filter.Name = helper.String("instanceIds")
	filter.Values = helper.Strings([]string{instanceId})
	tmpSet = append(tmpSet, &filter)
	paramMap["filters"] = tmpSet
	err = resource.Retry(readRetryTimeout*10, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		if result[0].SpecName != nil {
			_ = d.Set("node_spec", result[0].SpecName)
		}

		if result[0].InstanceName != nil {
			_ = d.Set("cluster_name", result[0].InstanceName)
		}

		if result[0].AutoRenewFlag != nil {
			if *result[0].AutoRenewFlag == AutoRenewFlagTrue {
				_ = d.Set("auto_renew_flag", true)
			} else {
				_ = d.Set("auto_renew_flag", false)
			}

		}

		return nil
	})

	if err != nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRabbitmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = tdmq.NewModifyRabbitMQVipInstanceRequest()
		instanceId = d.Id()
	)

	immutableArgs := []string{"zone_ids", "vpc_id", "subnet_id", "node_spec", "node_num", "storage_size", "enable_create_default_ha_mirror_queue", "auto_renew_flag", "time_span"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.InstanceId = &instanceId

	if d.HasChange("cluster_name") {
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRabbitMQVipInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVipInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	if err := service.DeleteTdmqRabbitmqVipInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
