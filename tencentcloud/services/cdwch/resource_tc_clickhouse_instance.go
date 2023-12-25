package cdwch

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClickhouseInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudClickhouseInstanceRead,
		Create: resourceTencentCloudClickhouseInstanceCreate,
		Update: resourceTencentCloudClickhouseInstanceUpdate,
		Delete: resourceTencentCloudClickhouseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Availability zone.",
			},
			"ha_flag": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether it is highly available.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet.",
			},
			"product_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product version.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Billing type: `PREPAID` prepaid, `POSTPAID_BY_HOUR` postpaid.",
			},
			"renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "PREPAID needs to be passed. Whether to renew automatically. 1 means auto renewal is enabled.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
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
					},
				},
			},
			"cls_log_set_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "CLS log set id.",
			},

			"cos_bucket_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "COS bucket name.",
			},

			"mount_disk_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is mounted on a bare disk.",
			},

			"ha_zk": {
				Optional:    true,
				Computed:    true,
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
							Description: "Node count. NOTE: Only support value 3.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Disk size.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
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

func resourceTencentCloudClickhouseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CdwchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	instanceInfos, err := service.DescribeInstancesNew(ctx, instanceId)
	if err != nil {
		return err
	}

	if len(instanceInfos) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource clickhouse instance [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	instanceInfo := instanceInfos[0]
	_ = d.Set("zone", instanceInfo.Zone)
	haFlag, err := strconv.ParseBool(*instanceInfo.HA)
	if err != nil {
		return err
	}
	_ = d.Set("ha_flag", haFlag)
	_ = d.Set("vpc_id", instanceInfo.VpcId)
	_ = d.Set("subnet_id", instanceInfo.SubnetId)
	_ = d.Set("product_version", instanceInfo.Version)
	_ = d.Set("instance_name", instanceInfo.InstanceName)
	_ = d.Set("charge_type", *instanceInfo.PayMode)
	if *instanceInfo.RenewFlag {
		_ = d.Set("renew_flag", 1)
	}
	_ = d.Set("expire_time", instanceInfo.ExpireTime)
	_ = d.Set("cos_bucket_name", instanceInfo.CosBucketName)
	_ = d.Set("mount_disk_type", instanceInfo.MountDiskType)
	_ = d.Set("ha_zk", instanceInfo.HAZk)

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

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "cdwch", "cdwchInstance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudClickhouseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwch_tmp_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwchClient().CreateInstanceNew(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := CdwchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instanceInfo, innerErr := service.DescribeInstance(ctx, instanceId)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
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
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::cdwch:%s:uin/:cdwchInstance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudClickhouseInstanceRead(d, meta)
}

func resourceTencentCloudClickhouseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("cdwch", "cdwchInstance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	immutableArgs := []string{"zone", "ha_flag", "vpc_id", "subnet_id", "product_version", "instance_name", "charge_type", "renew_flag", "time_span", "data_spec", "cls_log_set_id", "cos_bucket_name", "mount_disk_type", "ha_zk", "common_spec"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudClickhouseInstanceRead(d, meta)
}

func resourceTencentCloudClickhouseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CdwchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()

	if d.Get("charge_type").(string) == "PREPAID" {
		if err := service.DestroyInstance(ctx, instanceId); err != nil {
			return err
		}

		err := resource.Retry(5*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			instanceInfo, innerErr := service.DescribeInstance(ctx, instanceId)
			if innerErr != nil {
				return tccommon.RetryError(innerErr)
			}
			if *instanceInfo.Status != "Isolated" {
				return resource.RetryableError(fmt.Errorf("Still isolating"))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if err := service.DestroyInstance(ctx, instanceId); err != nil {
		return err
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instancesList, innerErr := service.DescribeInstancesNew(ctx, instanceId)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
		}
		if len(instancesList) != 0 {
			return resource.RetryableError(fmt.Errorf("Still destroying"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
