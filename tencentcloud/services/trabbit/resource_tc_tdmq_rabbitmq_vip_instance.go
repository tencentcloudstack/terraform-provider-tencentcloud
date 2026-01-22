package trabbit

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqRabbitmqVipInstance() *schema.Resource {
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
				Description: "Node specifications. Valid values: rabbit-vip-basic-5 (for 2C4G), rabbit-vip-profession-2c8g (for 2C8G), rabbit-vip-basic-1 (for 4C8G), rabbit-vip-profession-4c16g (for 4C16G), rabbit-vip-basic-2 (for 8C16G), rabbit-vip-profession-8c32g (for 8C32G), rabbit-vip-basic-4 (for 16C32G), rabbit-vip-profession-16c64g (for 16C64G). The default is rabbit-vip-basic-1. NOTE: The above specifications may be sold out or removed from the shelves.",
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
			"pay_mode": {
				Optional:    true,
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Payment method: 0 indicates postpaid; 1 indicates prepaid. Default: prepaid.",
			},
			"cluster_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster version, the default is `3.8.30`, valid values: `3.8.30`, `3.11.8` and `3.13.7`.",
			},
			"public_access_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public Network Access Point.",
			},
			"vpcs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of VPC Access Points.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"vpc_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC Endpoint.",
						},
						"vpc_data_stream_endpoint_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status Of Vpc Endpoint.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqVipInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
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

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("cluster_version"); ok {
		request.ClusterVersion = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateRabbitMQVipInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tdmq rabbitmqVipInstance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	// wait
	paramMap := make(map[string]interface{})
	tmpSet := make([]*tdmq.Filter, 0)
	filter := tdmq.Filter{}
	filter.Name = helper.String("instanceIds")
	filter.Values = helper.Strings([]string{instanceId})
	tmpSet = append(tmpSet, &filter)
	paramMap["filters"] = tmpSet
	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("resource `tencentcloud_tdmq_rabbitmq_vip_instance` %s does not exist", instanceId))
		}

		if len(result) != 1 {
			return resource.NonRetryableError(fmt.Errorf("resource `tencentcloud_tdmq_rabbitmq_vip_instance` %s id error", instanceId))
		}

		if *result[0].Status == svctdmq.RabbitMQVipInstanceRunning {
			return resource.RetryableError(fmt.Errorf("rabbitmq_vip_instance status is creating"))
		} else if *result[0].Status == svctdmq.RabbitMQVipInstanceSuccess {
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("rabbitmq_vip_instance status illegal"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		instanceId = d.Id()
	)

	rabbitmqVipInstance, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if rabbitmqVipInstance == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tdmq_rabbitmq_vip_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
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

	if rabbitmqVipInstance.ClusterInfo.PayMode != nil {
		_ = d.Set("pay_mode", rabbitmqVipInstance.ClusterInfo.PayMode)
	}

	if rabbitmqVipInstance.ClusterInfo.ClusterVersion != nil {
		_ = d.Set("cluster_version", rabbitmqVipInstance.ClusterInfo.ClusterVersion)
	}

	paramMap := make(map[string]interface{})
	tmpSet := make([]*tdmq.Filter, 0)
	filter := tdmq.Filter{}
	filter.Name = helper.String("instanceIds")
	filter.Values = helper.Strings([]string{instanceId})
	tmpSet = append(tmpSet, &filter)
	paramMap["filters"] = tmpSet
	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result[0].SpecName != nil {
			_ = d.Set("node_spec", result[0].SpecName)
		}

		if result[0].InstanceName != nil {
			_ = d.Set("cluster_name", result[0].InstanceName)
		}

		if result[0].AutoRenewFlag != nil {
			if *result[0].AutoRenewFlag == svctdmq.AutoRenewFlagTrue {
				_ = d.Set("auto_renew_flag", true)
			} else {
				_ = d.Set("auto_renew_flag", false)
			}
		}

		if result[0].PublicAccessEndpoint != nil {
			_ = d.Set("public_access_endpoint", result[0].PublicAccessEndpoint)
		}

		if result[0].Vpcs != nil {
			tmpList := make([]map[string]interface{}, 0, len(result[0].Vpcs))
			for _, vpc := range result[0].Vpcs {
				vpcMap := map[string]interface{}{}
				if vpc.VpcId != nil {
					vpcMap["vpc_id"] = vpc.VpcId
				}
				if vpc.SubnetId != nil {
					vpcMap["subnet_id"] = vpc.SubnetId
				}
				if vpc.VpcEndpoint != nil {
					vpcMap["vpc_endpoint"] = vpc.VpcEndpoint
				}
				if vpc.VpcDataStreamEndpointStatus != nil {
					vpcMap["vpc_data_stream_endpoint_status"] = vpc.VpcDataStreamEndpointStatus
				}
				tmpList = append(tmpList, vpcMap)
			}
			_ = d.Set("vpcs", tmpList)
		}

		return nil
	})

	if err != nil {
		log.Printf("[WARN]%s resource `TdmqRabbitmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = tdmq.NewModifyRabbitMQVipInstanceRequest()
		instanceId = d.Id()
	)

	immutableArgs := []string{
		"zone_ids", "vpc_id", "subnet_id", "node_spec", "node_num",
		"storage_size", "enable_create_default_ha_mirror_queue",
		"auto_renew_flag", "time_span", "pay_mode", "cluster_version",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_name") {
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}

		request.InstanceId = &instanceId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVipInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tdmq rabbitmqVipInstance failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVipInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_vip_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		instanceId = d.Id()
	)

	if err := service.DeleteTdmqRabbitmqVipInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
