package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqProfessionalCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqProfessionalClusterCreate,
		Read:   resourceTencentCloudTdmqProfessionalClusterRead,
		Update: resourceTencentCloudTdmqProfessionalClusterUpdate,
		Delete: resourceTencentCloudTdmqProfessionalClusterDelete,
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
				Description: "Multi-AZ deployment select three Availability Zones, like: [200002,200003,200004]. Single availability zone deployment selects an availability zone, like [200002].",
			},

			"product_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster specification code. Reference[Professional Cluster Specifications](https://cloud.tencent.com/document/product/1179/83705).",
			},

			"storage_size": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Storage specifications. Reference[Professional Cluster Specifications](https://cloud.tencent.com/document/product/1179/83705).",
			},

			"auto_renew_flag": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to turn on automatic monthly renewal. `1`: turn on, `0`: turn off.",
			},

			"time_span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase duration, value range: 1~50. Default: 1.",
				ForceNew:    true,
				Computed:    true,
			},

			"cluster_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of cluster. It does not support Chinese characters and special characters except dashes and underscores and cannot exceed 64 characters.",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to automatically select vouchers. `1`: Yes, `0`: No. Default is `0`.",
				ForceNew:    true,
				Computed:    true,
			},

			"vpc": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Label of VPC network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Id of VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Id of Subnet.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTdmqProfessionalClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_professional_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = tdmq.NewCreateProClusterRequest()
		response  = tdmq.NewCreateProClusterResponse()
		clusterId string
	)
	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.(*schema.Set).List()
		for i := range zoneIdsSet {
			zoneIds := zoneIdsSet[i].(int)
			request.ZoneIds = append(request.ZoneIds, helper.IntInt64(zoneIds))
		}
	}

	if v, ok := d.GetOk("product_name"); ok {
		request.ProductName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	} else {
		request.TimeSpan = helper.Int64(1)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntInt64(v.(int))
	} else {
		request.AutoVoucher = helper.Int64(0)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "vpc"); ok {
		vpcInfo := tdmq.VpcInfo{}
		if v, ok := dMap["vpc_id"]; ok {
			vpcInfo.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			vpcInfo.SubnetId = helper.String(v.(string))
		}
		request.Vpc = &vpcInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateProCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq professionalCluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	conf := tccommon.BuildStateChangeConf([]string{"0"}, []string{"1"}, 8*tccommon.ReadRetryTimeout, time.Second, service.TdmqProfessionalClusterStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::tdmq:%s:uin/:cluster/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTdmqProfessionalClusterRead(d, meta)
}

func resourceTencentCloudTdmqProfessionalClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_professional_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	clusterId := d.Id()

	clusterInfo, err := service.DescribeTdmqProfessionalClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	professionalCluster, err := service.DescribePulsarProInstances(ctx, clusterId)
	if err != nil {
		return err
	}

	if professionalCluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqProfessionalCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterInfo.NodeDistribution != nil {
		var zoneIds []int64
		for _, node := range clusterInfo.NodeDistribution {
			zoneIds = append(zoneIds, helper.StrToInt64(*node.ZoneId))
		}
		_ = d.Set("zone_ids", zoneIds)
	}

	if professionalCluster.SpecName != nil {
		_ = d.Set("product_name", professionalCluster.SpecName)
	}

	if professionalCluster.MaxStorage != nil {
		_ = d.Set("storage_size", *professionalCluster.MaxStorage/3)
	}

	if professionalCluster.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", professionalCluster.AutoRenewFlag)
	}

	//if professionalCluster.TimeSpan != nil {
	//	_ = d.Set("time_span", professionalCluster.TimeSpan)
	//}

	if professionalCluster.InstanceName != nil {
		_ = d.Set("cluster_name", professionalCluster.InstanceName)
	}

	//if professionalCluster.AutoVoucher != nil {
	//	_ = d.Set("auto_voucher", professionalClusterInfo.AutoVoucher)
	//}

	if professionalCluster.VpcId != nil && professionalCluster.SubnetId != nil {
		vpcMap := map[string]interface{}{}

		vpcMap["vpc_id"] = professionalCluster.VpcId
		vpcMap["subnet_id"] = professionalCluster.SubnetId
		_ = d.Set("vpc", []interface{}{vpcMap})
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "tdmq", "cluster", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTdmqProfessionalClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_professional_cluster.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tdmq.NewModifyClusterRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"zone_ids", "product_name", "storage_size", "auto_renew_flag", "time_span", "auto_voucher", "vpc"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_name") {
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyCluster(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update tdmq professionalCluster failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tdmq", "cluster", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTdmqProfessionalClusterRead(d, meta)
}

func resourceTencentCloudTdmqProfessionalClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_professional_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	clusterId := d.Id()

	if err := service.DeleteTdmqProfessionalClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
