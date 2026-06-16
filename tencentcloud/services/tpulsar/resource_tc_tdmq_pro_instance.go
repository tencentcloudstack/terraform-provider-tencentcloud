package tpulsar

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"
)

func ResourceTencentCloudTdmqProInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqProInstanceCreate,
		Read:   resourceTencentCloudTdmqProInstanceRead,
		Update: resourceTencentCloudTdmqProInstanceUpdate,
		Delete: resourceTencentCloudTdmqProInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Multi-AZ deployment: select three availability zones, e.g. [200002,200003,200004]. Single-AZ deployment: select one availability zone, e.g. [200002].",
			},
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster specification code. Reference [Professional Cluster Specifications](https://cloud.tencent.com/document/product/1179/83705).",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Storage specification. Reference [Professional Cluster Specifications](https://cloud.tencent.com/document/product/1179/83705).",
			},
			"vpc": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
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
				Description: "VPC network configuration.",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster name. It does not support Chinese characters and special characters except dashes and underscores and cannot exceed 64 characters.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Purchase duration, value range: 1~50.",
			},
			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable auto-renewal. `1`: enable, `0`: disable.",
			},
			"auto_voucher": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to automatically select vouchers. `1`: yes, `0`: no. Default is `0`.",
			},
			"instance_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cluster version information.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster description, up to 128 characters.",
			},
			"public_access_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable public network access. Can only be set to true.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Cluster status. 0: creating, 1: normal, 2: deleting, 3: deleted, 4: isolated, 5: creation failed, 6: deletion failed.",
			},
			"deal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub-order number.",
			},
			"big_deal_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Order number.",
			},
		},
	}
}

func resourceTencentCloudTdmqProInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_pro_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = tdmq.NewCreateProClusterRequest()
		response *tdmq.CreateProClusterResponse
	)

	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.(*schema.Set).List()
		for _, item := range zoneIdsSet {
			request.ZoneIds = append(request.ZoneIds, helper.IntInt64(item.(int)))
		}
	}

	if v, ok := d.GetOk("product_name"); ok {
		request.ProductName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_version"); ok {
		request.InstanceVersion = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateProCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_tdmq_pro_instance failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_tdmq_pro_instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s logId=%s, d.Id()=%s", logId, logId, d.Id())

	if response.Response.ClusterId == nil || *response.Response.ClusterId == "" {
		return fmt.Errorf("Create tencentcloud_tdmq_pro_instance failed, ClusterId is empty")
	}

	clusterId := *response.Response.ClusterId
	d.SetId(clusterId)

	if response.Response.DealName != nil {
		_ = d.Set("deal_name", response.Response.DealName)
	}

	if response.Response.BigDealId != nil {
		_ = d.Set("big_deal_id", response.Response.BigDealId)
	}

	// Wait for cluster to be ready
	service := svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	conf := tccommon.BuildStateChangeConf([]string{"0"}, []string{"1"}, 8*tccommon.ReadRetryTimeout, time.Second, service.TdmqProfessionalClusterStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTdmqProInstanceRead(d, meta)
}

func resourceTencentCloudTdmqProInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_pro_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clusterId = d.Id()
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	var clusterInfo *tdmq.Cluster
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := tdmq.NewDescribeClustersRequest()
		request.ClusterIdList = []*string{&clusterId}

		var iacExtInfo connectivity.IacExtInfo
		iacExtInfo.InstanceId = clusterId
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient(iacExtInfo).DescribeClusters(request)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || len(result.Response.ClusterSet) == 0 {
			return nil
		}

		clusterInfo = result.Response.ClusterSet[0]
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read tencentcloud_tdmq_pro_instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if clusterInfo == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tdmq_pro_instance` id=%s not found, removing from state", logId, d.Id())
		d.SetId("")
		return nil
	}

	if clusterInfo.ClusterId != nil {
		_ = d.Set("cluster_id", clusterInfo.ClusterId)
	}

	if clusterInfo.ClusterName != nil {
		_ = d.Set("cluster_name", clusterInfo.ClusterName)
	}

	if clusterInfo.Remark != nil {
		_ = d.Set("remark", clusterInfo.Remark)
	}

	if clusterInfo.Status != nil {
		_ = d.Set("status", clusterInfo.Status)
	}

	if clusterInfo.PublicEndPoint != nil && *clusterInfo.PublicEndPoint != "" {
		_ = d.Set("public_access_enabled", true)
	}

	// Read VPC info and other details from professional cluster instance detail
	clusterDetail, err := service.DescribeTdmqProfessionalClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if clusterDetail != nil {
		if clusterDetail.NodeDistribution != nil {
			var zoneIds []int64
			for _, node := range clusterDetail.NodeDistribution {
				if node.ZoneId != nil {
					zoneIds = append(zoneIds, helper.StrToInt64(*node.ZoneId))
				}
			}
			if len(zoneIds) > 0 {
				_ = d.Set("zone_ids", zoneIds)
			}
		}
	}

	// Read instance-level details
	professionalCluster, err := service.DescribePulsarProInstances(ctx, clusterId)
	if err != nil {
		return err
	}

	if professionalCluster != nil {
		if professionalCluster.VpcId != nil && professionalCluster.SubnetId != nil {
			vpcMap := map[string]interface{}{
				"vpc_id":    professionalCluster.VpcId,
				"subnet_id": professionalCluster.SubnetId,
			}
			_ = d.Set("vpc", []interface{}{vpcMap})
		}

		if professionalCluster.SpecName != nil {
			_ = d.Set("product_name", professionalCluster.SpecName)
		}

		if professionalCluster.MaxStorage != nil {
			_ = d.Set("storage_size", int(*professionalCluster.MaxStorage/3))
		}

		if professionalCluster.AutoRenewFlag != nil {
			_ = d.Set("auto_renew_flag", int(*professionalCluster.AutoRenewFlag))
		}
	}

	return nil
}

func resourceTencentCloudTdmqProInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_pro_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		clusterId = d.Id()
	)

	needChange := false
	request := tdmq.NewModifyClusterRequest()
	request.ClusterId = &clusterId

	if d.HasChange("cluster_name") {
		needChange = true
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		needChange = true
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		} else {
			request.Remark = helper.String("")
		}
	}

	if d.HasChange("public_access_enabled") {
		needChange = true
		if v, ok := d.GetOk("public_access_enabled"); ok {
			request.PublicAccessEnabled = helper.Bool(v.(bool))
		}
	}

	if needChange {
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyCluster(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tencentcloud_tdmq_pro_instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTdmqProInstanceRead(d, meta)
}

func resourceTencentCloudTdmqProInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_pro_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		clusterId = d.Id()
		request   = tdmq.NewDeleteClusterRequest()
	)

	request.ClusterId = &clusterId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().DeleteCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_tdmq_pro_instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
