/*
Provides a resource to create a tdmq rocketmq_vip_instance

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_vip_instance" "rocketmq_vip_instance" {
  name = ""
  spec = ""
  node_count =
  storage_size =
  zone_ids =
  vpc_info {
		vpc_id = ""
		subnet_id = ""

  }
  time_span =
}
```

Import

tdmq rocketmq_vip_instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_vip_instance.rocketmq_vip_instance rocketmq_vip_instance_id
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

func resourceTencentCloudTdmqRocketmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqVipInstanceCreate,
		Read:   resourceTencentCloudTdmqRocketmqVipInstanceRead,
		Update: resourceTencentCloudTdmqRocketmqVipInstanceUpdate,
		Delete: resourceTencentCloudTdmqRocketmqVipInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"spec": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance specification:Basic type, rocket-vip-basic-1Standard type, rocket-vip-basic-2Advanced Type I, rocket-vip-basic-3Advanced Type II, rocket-vip-basic-4.",
			},

			"node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of nodes, minimum 2, maximum 20.",
			},

			"storage_size": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Single node storage space, in GB, minimum 200GB.",
			},

			"zone_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The Zone ID list for node deployment, such as Guangzhou Zone 1, is 100001. For details, please refer to the official website of Tencent Cloud.",
			},

			"vpc_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "VPC information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},

			"time_span": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Purchase period, in months.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqVipInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateRocketMQVipInstanceRequest()
		response  = tdmq.NewCreateRocketMQVipInstanceResponse()
		clusterId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec"); ok {
		request.Spec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("node_count"); ok {
		request.NodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.(*schema.Set).List()
		for i := range zoneIdsSet {
			zoneIds := zoneIdsSet[i].(string)
			request.ZoneIds = append(request.ZoneIds, &zoneIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "vpc_info"); ok {
		vpcInfo := tdmq.VpcInfo{}
		if v, ok := dMap["vpc_id"]; ok {
			vpcInfo.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			vpcInfo.SubnetId = helper.String(v.(string))
		}
		request.VpcInfo = &vpcInfo
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRocketMQVipInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rocketmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudTdmqRocketmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	rocketmqVipInstanceId := d.Id()

	rocketmqVipInstance, err := service.DescribeTdmqRocketmqVipInstanceById(ctx, clusterId)
	if err != nil {
		return err
	}

	if rocketmqVipInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rocketmqVipInstance.Name != nil {
		_ = d.Set("name", rocketmqVipInstance.Name)
	}

	if rocketmqVipInstance.Spec != nil {
		_ = d.Set("spec", rocketmqVipInstance.Spec)
	}

	if rocketmqVipInstance.NodeCount != nil {
		_ = d.Set("node_count", rocketmqVipInstance.NodeCount)
	}

	if rocketmqVipInstance.StorageSize != nil {
		_ = d.Set("storage_size", rocketmqVipInstance.StorageSize)
	}

	if rocketmqVipInstance.ZoneIds != nil {
		_ = d.Set("zone_ids", rocketmqVipInstance.ZoneIds)
	}

	if rocketmqVipInstance.VpcInfo != nil {
		vpcInfoMap := map[string]interface{}{}

		if rocketmqVipInstance.VpcInfo.VpcId != nil {
			vpcInfoMap["vpc_id"] = rocketmqVipInstance.VpcInfo.VpcId
		}

		if rocketmqVipInstance.VpcInfo.SubnetId != nil {
			vpcInfoMap["subnet_id"] = rocketmqVipInstance.VpcInfo.SubnetId
		}

		_ = d.Set("vpc_info", []interface{}{vpcInfoMap})
	}

	if rocketmqVipInstance.TimeSpan != nil {
		_ = d.Set("time_span", rocketmqVipInstance.TimeSpan)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyRocketMQInstanceSpecRequest()

	rocketmqVipInstanceId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"name", "spec", "node_count", "storage_size", "zone_ids", "vpc_info", "time_span"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("node_count") {
		if v, ok := d.GetOkExists("node_count"); ok {
			request.NodeCount = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("storage_size") {
		if v, ok := d.GetOkExists("storage_size"); ok {
			request.StorageSize = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRocketMQInstanceSpec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmq rocketmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqVipInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	rocketmqVipInstanceId := d.Id()

	if err := service.DeleteTdmqRocketmqVipInstanceById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
