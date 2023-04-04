package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClickHoustInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudClickHouseInstanceRead,
		Create: resourceTencentCloudClickHouseInstanceCreate,
		Update: resourceTencentCloudClickHouseInstanceUpdate,
		Delete: resourceTencentCloudClickHouseInstanceDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone.",
			},
			"ha_flag": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Whether it is highly available.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Private network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet.",
			},
			"product_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product version.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance name.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Billing type: `PREPAID` prepaid, `POSTPAID_BY_HOUR` postpaid",
			},
			"renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "PREPAID needs to be passed. Whether to renew automatically. 1 means auto renewal is enabled.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Prepaid needs to be delivered, billing time length, how many months.",
			},
			"data_spec": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Data spec.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Spec name.",
						},
						"count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data spec count.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Disk size.",
						},
						"scale_out_cluster": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "`v_cluster` grouping. Must set when update NodeCount.The new expansion node will be added to the selected v_cluster packet, and the submission synchronization VIP will take effect.",
						},
					},
				},
			},
			"cls_log_set_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CLS log set id.",
			},

			"cos_bucket_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS bucket name.",
			},

			"mount_disk_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is mounted on a bare disk.",
			},

			"ha_zk": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether ZK is highly available.",
			},

			"common_spec": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "ZK node.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Spec name.",
						},
						"count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "count.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Disk size.",
						},
						"scale_out_cluster": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "`v_cluster` grouping. Must set when update NodeCount.The new expansion node will be added to the selected v_cluster packet, and the submission synchronization VIP will take effect.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expire time.",
			},
		},
	}
}

func resourceTencentCloudClickHouseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instanceInfo, err := service.DescribeInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource clickhouse instance [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("zone", instanceInfo.Zone)
	_ = d.Set("ha_flag", instanceInfo.HA)
	_ = d.Set("vpc_id", instanceInfo.VpcId)
	_ = d.Set("subnet_id", instanceInfo.SubnetId)
	_ = d.Set("product_version", instanceInfo.Version)
	_ = d.Set("instance_name", instanceInfo.InstanceName)
	_ = d.Set("charge_type", *instanceInfo.PayMode)
	_ = d.Set("renew_flag", instanceInfo.RenewFlag)
	_ = d.Set("expire_time", instanceInfo.ExpireTime)

	if instanceInfo.MasterSummary != nil {
		dataSpec := make(map[string]interface{})
		dataSpec["spec_name"] = instanceInfo.MasterSummary.Spec
		dataSpec["count"] = instanceInfo.MasterSummary.NodeSize
		dataSpec["disk_size"] = instanceInfo.MasterSummary.Disk
		_ = d.Set("data_spec", []map[string]interface{}{dataSpec})
	}

	_ = d.Set("cls_log_set_id", instanceInfo.ClsLogSetId)
	if instanceInfo.CommonSummary != nil {
		commonSpec := make(map[string]interface{})
		commonSpec["spec_name"] = instanceInfo.CommonSummary.Spec
		commonSpec["count"] = instanceInfo.CommonSummary.NodeSize
		commonSpec["disk_size"] = instanceInfo.CommonSummary.Disk
		_ = d.Set("common_spec", []map[string]interface{}{commonSpec})
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cdwch", "cdwchInstance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudClickHouseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdwch_tmp_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdwch.NewCreateInstanceNewRequest()
		response   = cdwch.NewCreateInstanceNewResponse()
		instanceId string
	)
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ha_flag"); v != nil {
		request.HaFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.UserVPCId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.UserSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_version"); ok {
		request.ProductVersion = helper.String(v.(string))
	}

	charge := cdwch.Charge{}
	if v, ok := d.GetOk("charge_type"); ok {
		charge.ChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("renew_flag"); ok {
		charge.RenewFlag = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("time_span"); ok {
		charge.TimeSpan = helper.IntInt64(v.(int))
	}
	request.ChargeProperties = &charge

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "data_spec"); ok {
		nodeSpec := cdwch.NodeSpec{}
		if v, ok := dMap["spec_name"]; ok {
			nodeSpec.SpecName = helper.String(v.(string))
		}
		if v, ok := dMap["count"]; ok {
			nodeSpec.Count = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["disk_size"]; ok {
			nodeSpec.DiskSize = helper.IntInt64(v.(int))
		}
		request.DataSpec = &nodeSpec
	}

	if v, ok := d.GetOk("cls_log_set_id"); ok {
		request.ClsLogSetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_bucket_name"); ok {
		request.CosBucketName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("mount_disk_type"); v != nil {
		request.MountDiskType = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("ha_zk"); v != nil {
		request.HAZk = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "common_spec"); ok {
		nodeSpec := cdwch.NodeSpec{}
		if v, ok := dMap["spec_name"]; ok {
			nodeSpec.SpecName = helper.String(v.(string))
		}
		if v, ok := dMap["count"]; ok {
			nodeSpec.Count = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["disk_size"]; ok {
			nodeSpec.DiskSize = helper.IntInt64(v.(int))
		}
		request.CommonSpec = &nodeSpec
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().CreateInstanceNew(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdwch tmpInstance failed, reason:%+v", logId, err)
		return err
	}
	instanceId = *response.Response.InstanceId
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	err = resource.Retry(5*writeRetryTimeout, func() *resource.RetryError {
		instanceInfo, innerErr := service.DescribeInstance(ctx, instanceId)
		if innerErr != nil {
			return retryError(innerErr)
		}
		if *instanceInfo.Status != "Serving" {
			return resource.RetryableError(fmt.Errorf("Still creating"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cdwch:%s:uin/:cdwchInstance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudClickHouseInstanceRead(d, meta)
}

func resourceTencentCloudClickHouseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceInfo, err := service.DescribeInstance(ctx, d.Id())
	if err != nil {
		return err
	}

	if instanceInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource clickhouse instance [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	subnetId := *instanceInfo.SubnetId
	subnetInfo, has, err := vpcService.DescribeSubnet(ctx, subnetId, nil, "", "")
	if err != nil {
		return err
	}

	if has == 0 {
		return fmt.Errorf("subnet subnet_id=%s not found", subnetId)
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("cdwch", "cdwchInstance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	specChangeFlag := false
	if d.HasChange("data_spec.0.spec_name") {
		dataSpecName := d.Get("data_spec.0.spec_name").(string)
		err := service.ScaleUpInstance(ctx, d.Id(), NODE_TYPE_CLICKHOUSE, dataSpecName)
		if err != nil {
			return err
		}
		specChangeFlag = true
	}
	if d.HasChange("data_spec.0.count") {
		dataSpecCount := d.Get("data_spec.0.count").(int)
		var scaleOutCluster string
		if v, ok := d.GetOk("data_spec.0.scale_out_cluster"); ok {
			scaleOutCluster = v.(string)
		} else {
			return fmt.Errorf("Must set scale_out_cluster when update node count.")
		}
		err := service.ScaleOutInstance(ctx, d.Id(), NODE_TYPE_CLICKHOUSE, scaleOutCluster, dataSpecCount, int(subnetInfo.availableIpCount))
		if err != nil {
			return err
		}
		specChangeFlag = true
	}
	if d.HasChange("data_spec.0.disk_size") {
		dataSpecDiskSize := d.Get("data_spec.0.disk_size").(int)
		err := service.ResizeDisk(ctx, d.Id(), NODE_TYPE_CLICKHOUSE, dataSpecDiskSize)
		if err != nil {
			return err
		}
		specChangeFlag = true
	}

	if d.HasChange("common_spec.0.spec_name") {
		commonSpecName := d.Get("common_spec.0.spec_name").(string)
		err := service.ScaleUpInstance(ctx, d.Id(), NODE_TYPE_ZOOKEEPER, commonSpecName)
		if err != nil {
			return err
		}
		specChangeFlag = true
	}
	if d.HasChange("common_spec.0.count") {
		commonSpecCount := d.Get("common_spec.0.count").(int)
		var scaleOutCluster string
		if v, ok := d.GetOk("common_spec.0.scale_out_cluster"); ok {
			scaleOutCluster = v.(string)
		} else {
			return fmt.Errorf("Must set scale_out_cluster when update node count.")
		}

		err = service.ScaleOutInstance(ctx, d.Id(), NODE_TYPE_ZOOKEEPER, scaleOutCluster, commonSpecCount, int(subnetInfo.availableIpCount))
		if err != nil {
			return err
		}
		specChangeFlag = true
	}
	if d.HasChange("common_spec.0.disk_size") {
		commonSpecDiskSize := d.Get("common_spec.0.disk_size").(int)
		err := service.ResizeDisk(ctx, d.Id(), NODE_TYPE_ZOOKEEPER, commonSpecDiskSize)
		if err != nil {
			return err
		}
		specChangeFlag = true
	}
	if specChangeFlag {
		changeState := false
		err := resource.Retry(3*writeRetryTimeout, func() *resource.RetryError {
			instanceInfo, innerErr := service.DescribeInstance(ctx, d.Id())
			if innerErr != nil {
				return retryError(innerErr)
			}
			if *instanceInfo.Status == "Changing" {
				changeState = true
			}
			if !changeState || *instanceInfo.Status != "Serving" {
				return resource.RetryableError(fmt.Errorf("Still updating"))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudClickHouseInstanceRead(d, meta)
}

func resourceTencentCloudClickHouseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DestroyInstance(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceInfo, innerErr := service.DescribeInstance(ctx, instanceId)
		if innerErr != nil {
			return retryError(innerErr)
		}
		if *instanceInfo.Status != "Deleted" {
			return resource.RetryableError(fmt.Errorf("Still destroying"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
