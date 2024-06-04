package cvm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipCreate,
		Read:   resourceTencentCloudEipRead,
		Update: resourceTencentCloudEipUpdate,
		Delete: resourceTencentCloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"anti_ddos_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of anti DDos package, it must set when `type` is `AntiDDoSEIP`.",
			},

			"anycast_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The zone of anycast. Valid value: `ANYCAST_ZONE_GLOBAL` and `ANYCAST_ZONE_OVERSEAS`.",
			},

			"applicable_for_clb": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the anycast eip can be associated to a CLB.",
				Deprecated:  "It has been deprecated from version 1.27.0.",
			},

			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Auto renew flag.  0 - default state (manual renew); 1 - automatic renew; 2 - explicit no automatic renew. NOTES: Only supported prepaid EIP.",
			},

			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of bandwidth package, it will set when `internet_charge_type` is `BANDWIDTH_PACKAGE`.",
			},

			"egress": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Network egress. It defaults to `center_egress1`. If you want to try the egress feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).",
			},

			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The charge type of eip. Valid values: `BANDWIDTH_PACKAGE`, `BANDWIDTH_POSTPAID_BY_HOUR`, `BANDWIDTH_PREPAID_BY_MONTH` and `TRAFFIC_POSTPAID_BY_HOUR`.",
			},

			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The bandwidth limit of EIP, unit is Mbps.",
			},

			"internet_service_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Internet service provider of eip. Valid value: `BGP`, `CMCC`, `CTCC` and `CUCC`.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of eip.",
			},

			"prepaid_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Period of instance. Default value: `1`. Valid value: `1`, `2`, `3`, `4`, `6`, `7`, `8`, `9`, `12`, `24`, `36`. NOTES: must set when `internet_charge_type` is `BANDWIDTH_PREPAID_BY_MONTH`.",
			},

			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The elastic IP address.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The EIP current status.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of eip.",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "EIP",
				Description: "The type of eip. Valid value:  `EIP` and `AnycastEIP` and `HighQualityEIP` and `AntiDDoSEIP`. Default is `EIP`.",
			},
		},
	}
}

func resourceTencentCloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		eipId string
	)
	var (
		request  = vpc.NewAllocateAddressesRequest()
		response = vpc.NewAllocateAddressesResponse()
	)

	if v, ok := d.GetOk("internet_service_provider"); ok {
		request.InternetServiceProvider = helper.String(v.(string))
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.AddressType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("anycast_zone"); ok {
		request.AnycastZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			tag := vpc.Tag{}
			if v, ok := tagsMap["key"]; ok {
				tag.Key = helper.String(v.(string))
			}
			if v, ok := tagsMap["value"]; ok {
				tag.Value = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.BandwidthPackageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.AddressName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("egress"); ok {
		request.Egress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("anti_ddos_package_id"); ok {
		request.AntiDDoSPackageId = helper.String(v.(string))
	}

	if err := resourceTencentCloudEipCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AllocateAddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eip failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.AddressSet) < 1 {
		return fmt.Errorf("resource `tencentcloud_eip` create failed.")
	}

	eipId = *response.Response.AddressSet[0]

	if err := resourceTencentCloudEipCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(eipId)

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	eipId := d.Id()

	respData, err := service.DescribeEipById(ctx, eipId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `eip` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.AddressId != nil {
		_ = d.Set("address_id", respData.AddressId)
		eipId = *respData.AddressId
	}

	if respData.AddressName != nil {
		_ = d.Set("name", respData.AddressName)
	}

	if respData.AddressStatus != nil {
		_ = d.Set("status", respData.AddressStatus)
	}

	if respData.AddressIp != nil {
		_ = d.Set("public_ip", respData.AddressIp)
	}

	if respData.AddressType != nil {
		_ = d.Set("type", respData.AddressType)
	}

	if respData.InternetServiceProvider != nil {
		_ = d.Set("internet_service_provider", respData.InternetServiceProvider)
	}

	if respData.Bandwidth != nil {
		_ = d.Set("internet_max_bandwidth_out", respData.Bandwidth)
	}

	if respData.InternetChargeType != nil {
		_ = d.Set("internet_charge_type", respData.InternetChargeType)
	}

	if respData.Egress != nil {
		_ = d.Set("egress", respData.Egress)
	}

	if respData.AntiDDoSPackageId != nil {
		_ = d.Set("anti_ddos_package_id", respData.AntiDDoSPackageId)
	}

	if err := resourceTencentCloudEipReadPostHandleResponse0(ctx, respData); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"anti_ddos_package_id", "applicable_for_clb", "auto_renew_flag", "bandwidth_package_id", "egress", "prepaid_period", "tags"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	eipId := d.Id()

	if err := resourceTencentCloudEipUpdateOnStart(ctx); err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{"name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := vpc.NewModifyAddressAttributeRequest()

		request.AddressId = helper.String(eipId)

		if v, ok := d.GetOk("name"); ok {
			request.AddressName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyAddressAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update eip failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange1 := false
	mutableArgs1 := []string{"internet_charge_type", "internet_max_bandwidth_out"}
	for _, v := range mutableArgs1 {
		if d.HasChange(v) {
			needChange1 = true
			break
		}
	}

	if needChange1 {
		request1 := vpc.NewModifyAddressInternetChargeTypeRequest()

		request1.AddressId = helper.String(eipId)

		if v, ok := d.GetOk("internet_charge_type"); ok {
			request1.InternetChargeType = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
			request1.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
		}

		if err := resourceTencentCloudEipUpdatePostFillRequest1(ctx, request1); err != nil {
			return err
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyAddressInternetChargeTypeWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update eip failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange2 := false
	mutableArgs2 := []string{"internet_max_bandwidth_out"}
	for _, v := range mutableArgs2 {
		if d.HasChange(v) {
			needChange2 = true
			break
		}
	}

	if needChange2 {
		request2 := vpc.NewModifyAddressesBandwidthRequest()

		request2.AddressIds = []*string{helper.String(eipId)}

		if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
			request2.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyAddressesBandwidthWithContext(ctx, request2)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request2.GetAction(), request2.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update eip failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = eipId
	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	eipId := d.Id()

	var (
		request  = vpc.NewReleaseAddressesRequest()
		response = vpc.NewReleaseAddressesResponse()
	)

	request.AddressIds = []*string{helper.String(eipId)}

	if err := resourceTencentCloudEipDeletePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReleaseAddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete eip failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	if err := resourceTencentCloudEipDeletePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	_ = eipId
	return nil
}
