package vpc

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcBandwidthPackageCreate,
		Read:   resourceTencentCloudVpcBandwidthPackageRead,
		Update: resourceTencentCloudVpcBandwidthPackageUpdate,
		Delete: resourceTencentCloudVpcBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Bandwidth packet type, default: `BGP`. " +
					"Optional value: `BGP`: common BGP shared bandwidth package; `HIGH_QUALITY_BGP`: High Quality BGP Shared Bandwidth Package; " +
					"`SINGLEISP_CMCC`: CMCC shared bandwidth package; `SINGLEISP_CTCC:`: CTCC shared bandwidth package; `SINGLEISP_CUCC`: CUCC shared bandwidth package.",
			},

			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Bandwidth package billing type, default: `TOP5_POSTPAID_BY_MONTH`." +
					" Optional value: `TOP5_POSTPAID_BY_MONTH`: TOP5 billed by monthly postpaid; `PERCENT95_POSTPAID_BY_MONTH`: 95 billed monthly postpaid;" +
					" `FIXED_PREPAID_BY_MONTH`: Monthly prepaid billing (Type FIXED_PREPAID_BY_MONTH product API capability is under construction);" +
					" `BANDWIDTH_POSTPAID_BY_DAY`: bandwidth billed by daily postpaid; `ENHANCED95_POSTPAID_BY_MONTH`: enhanced 95 billed monthly postpaid.",
			},

			"bandwidth_package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bandwidth package name.",
			},

			"internet_max_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Bandwidth packet speed limit size. Unit: Mbps, -1 means no speed limit.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},

			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The purchase duration of the prepaid monthly bandwidth package, unit: month, value range: 1~60.",
			},

			"egress": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Network egress. It defaults to `center_egress1`. If you want to try the egress feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).",
			},
		},
	}
}

func resourceTencentCloudVpcBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = vpc.NewCreateBandwidthPackageRequest()
		response *vpc.CreateBandwidthPackageResponse
	)

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bandwidth_package_name"); ok {
		request.BandwidthPackageName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth"); ok {
		request.InternetMaxBandwidth = helper.IntInt64(v.(int))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("egress"); ok {
		request.Egress = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateBandwidthPackage(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vpc bandwidthPackage failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc bandwidthPackage failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.BandwidthPackageId == nil {
		return fmt.Errorf("BandwidthPackageId is nil.")
	}

	bandwidthPackageId := *response.Response.BandwidthPackageId
	d.SetId(bandwidthPackageId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:bandwidthPackage/%s", region, bandwidthPackageId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	// wait
	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeVpcBandwidthPackage(ctx, bandwidthPackageId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance == nil {
			return resource.RetryableError(fmt.Errorf("vpc bandwidthPackage instance is being created, retry..."))
		}

		if *instance.Status == "CREATED" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("vpc bandwidthPackage instance status is %v, retry...", *instance.Status))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudVpcBandwidthPackageRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package.read")()
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		bandwidthPackageId = d.Id()
	)

	bandwidthPackage, err := service.DescribeVpcBandwidthPackage(ctx, bandwidthPackageId)
	if err != nil {
		return err
	}

	if bandwidthPackage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_vpc_bandwidth_package` [%s] not found, please check if it has been deleted.", logId, bandwidthPackageId)
		return nil
	}

	if bandwidthPackage.NetworkType != nil {
		_ = d.Set("network_type", bandwidthPackage.NetworkType)
	}

	if bandwidthPackage.ChargeType != nil {
		_ = d.Set("charge_type", bandwidthPackage.ChargeType)
	}

	if bandwidthPackage.BandwidthPackageName != nil {
		_ = d.Set("bandwidth_package_name", bandwidthPackage.BandwidthPackageName)
	}

	if bandwidthPackage.Bandwidth != nil {
		_ = d.Set("internet_max_bandwidth", bandwidthPackage.Bandwidth)
	}

	if bandwidthPackage.Egress != nil {
		_ = d.Set("egress", bandwidthPackage.Egress)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "bandwidthPackage", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpcBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		bandwidthPackageId = d.Id()
	)

	immutableArgs := []string{
		"network_type",
		"egress",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("bandwidth_package_name") || d.HasChange("charge_type") {
		request := vpc.NewModifyBandwidthPackageAttributeRequest()
		request.BandwidthPackageId = &bandwidthPackageId
		if v, ok := d.GetOk("bandwidth_package_name"); ok {
			request.BandwidthPackageName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("charge_type"); ok {
			request.ChargeType = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyBandwidthPackageAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s Modify vpc bandwidthPackage attribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("internet_max_bandwidth") {
		request := vpc.NewModifyBandwidthPackageBandwidthRequest()
		request.BandwidthPackageId = &bandwidthPackageId
		if v, ok := d.GetOkExists("internet_max_bandwidth"); ok {
			request.InternetMaxBandwidth = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyBandwidthPackageBandwidth(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s Modify vpc bandwidthPackage bandWidth failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("vpc", "bandwidthPackage", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcBandwidthPackageRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		bandwidthPackageId = d.Id()
	)

	if err := service.DeleteVpcBandwidthPackageById(ctx, bandwidthPackageId); err != nil {
		return err
	}

	return nil
}
