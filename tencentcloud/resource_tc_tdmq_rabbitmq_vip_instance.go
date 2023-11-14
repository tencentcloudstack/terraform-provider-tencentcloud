/*
Provides a resource to create a tdmq rabbitmq_vip_instance

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "rabbitmq_vip_instance" {
  zone_ids =
  vpc_id = ""
  subnet_id = ""
  cluster_name = ""
  node_spec = ""
  node_num =
  storage_size =
  enable_create_default_ha_mirror_queue =
  auto_renew_flag =
  time_span =
}
```

Import

tdmq rabbitmq_vip_instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_vip_instance.rabbitmq_vip_instance rabbitmq_vip_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTdmqRabbitmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRabbitmqVipInstanceCreate,
		Read:   resourceTencentCloudTdmqRabbitmqVipInstanceRead,
		Update: resourceTencentCloudTdmqRabbitmqVipInstanceUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqVipInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Availability zone.",
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
				Description: "Cluster name.",
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

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateRabbitMQVipInstanceRequest()
		response  = tdmq.NewCreateRabbitMQVipInstanceResponse()
		clusterId string
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

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	rabbitmqVipInstanceId := d.Id()

	rabbitmqVipInstance, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, clusterId)
	if err != nil {
		return err
	}

	if rabbitmqVipInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRabbitmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rabbitmqVipInstance.ZoneIds != nil {
		_ = d.Set("zone_ids", rabbitmqVipInstance.ZoneIds)
	}

	if rabbitmqVipInstance.VpcId != nil {
		_ = d.Set("vpc_id", rabbitmqVipInstance.VpcId)
	}

	if rabbitmqVipInstance.SubnetId != nil {
		_ = d.Set("subnet_id", rabbitmqVipInstance.SubnetId)
	}

	if rabbitmqVipInstance.ClusterName != nil {
		_ = d.Set("cluster_name", rabbitmqVipInstance.ClusterName)
	}

	if rabbitmqVipInstance.NodeSpec != nil {
		_ = d.Set("node_spec", rabbitmqVipInstance.NodeSpec)
	}

	if rabbitmqVipInstance.NodeNum != nil {
		_ = d.Set("node_num", rabbitmqVipInstance.NodeNum)
	}

	if rabbitmqVipInstance.StorageSize != nil {
		_ = d.Set("storage_size", rabbitmqVipInstance.StorageSize)
	}

	if rabbitmqVipInstance.EnableCreateDefaultHaMirrorQueue != nil {
		_ = d.Set("enable_create_default_ha_mirror_queue", rabbitmqVipInstance.EnableCreateDefaultHaMirrorQueue)
	}

	if rabbitmqVipInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", rabbitmqVipInstance.AutoRenewFlag)
	}

	if rabbitmqVipInstance.TimeSpan != nil {
		_ = d.Set("time_span", rabbitmqVipInstance.TimeSpan)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyRabbitMQVipInstanceRequest()

	rabbitmqVipInstanceId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"zone_ids", "vpc_id", "subnet_id", "cluster_name", "node_spec", "node_num", "storage_size", "enable_create_default_ha_mirror_queue", "auto_renew_flag", "time_span"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	rabbitmqVipInstanceId := d.Id()

	if err := service.DeleteTdmqRabbitmqVipInstanceById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
