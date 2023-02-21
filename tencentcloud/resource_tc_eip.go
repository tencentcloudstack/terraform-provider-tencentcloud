/*
Provides an EIP resource.

Example Usage

```hcl
resource "tencentcloud_eip" "foo" {
  name                 = "awesome_gateway_ip"
  bandwidth_package_id = "bwp-jtvzuky6"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  type                 = "EIP"
}
```

Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipCreate,
		Read:   resourceTencentCloudEipRead,
		Update: resourceTencentCloudEipUpdate,
		Delete: resourceTencentCloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
				Description:  "The name of eip.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     EIP_TYPE_EIP,
				ForceNew:    true,
				Description: "The type of eip. Valid value:  `EIP` and `AnycastEIP` and `HighQualityEIP`. Default is `EIP`.",
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
			"internet_service_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Internet service provider of eip. Valid value: `BGP`, `CMCC`, `CTCC` and `CUCC`.",
			},
			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The charge type of eip. Valid values: `BANDWIDTH_PACKAGE`, `BANDWIDTH_POSTPAID_BY_HOUR`, `BANDWIDTH_PREPAID_BY_MONTH` and `TRAFFIC_POSTPAID_BY_HOUR`.",
			},

			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue(EIP_AVAILABLE_PERIOD),
				Description:  "Period of instance. Default value: `1`. Valid value: `1`, `2`, `3`, `4`, `6`, `7`, `8`, `9`, `12`, `24`, `36`. NOTES: must set when `internet_charge_type` is `BANDWIDTH_PREPAID_BY_MONTH`.",
			},

			"auto_renew_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
				Description:  "Auto renew flag.  0 - default state (manual renew); 1 - automatic renew; 2 - explicit no automatic renew. NOTES: Only supported prepaid EIP.",
			},

			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The bandwidth limit of EIP, unit is Mbps.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of eip.",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of bandwidth package, it will set when `internet_charge_type` is `BANDWIDTH_PACKAGE`.",
			},
			// computed
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
		},
	}
}

func resourceTencentCloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	var internetChargeType string

	request := vpc.NewAllocateAddressesRequest()
	if v, ok := d.GetOk("type"); ok {
		request.AddressType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("anycast_zone"); ok {
		request.AnycastZone = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_service_provider"); ok {
		request.InternetServiceProvider = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetChargeType = v.(string)
		request.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
	}

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		addressChargePrepaid := vpc.AddressChargePrepaid{}
		period := d.Get("prepaid_period")
		renewFlag := d.Get("auto_renew_flag")
		addressChargePrepaid.Period = helper.IntInt64(period.(int))
		addressChargePrepaid.AutoRenewFlag = helper.IntInt64(renewFlag.(int))
		request.AddressChargePrepaid = &addressChargePrepaid
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
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.BandwidthPackageId = helper.String(v.(string))
	}

	eipId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := client.UseVpcClient().AllocateAddresses(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if len(response.Response.AddressSet) < 1 {
			return resource.RetryableError(fmt.Errorf("eip id is nil"))
		}
		eipId = *response.Response.AddressSet[0]
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(eipId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName(VPC_SERVICE_TYPE, EIP_RESOURCE_TYPE, region, eipId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			log.Printf("[CRITAL]%s set eip tags failed: %+v", logId, err)
			return err
		}
	}

	// wait for status
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet)
		}
		if eip != nil && *eip.AddressStatus == EIP_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("eip is still creating"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := vpcService.ModifyEipName(ctx, eipId, name)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	eipId := d.Id()
	var eip *vpc.Address
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet)
		}
		eip = instance
		return nil
	})
	if err != nil {
		return err
	}
	if eip == nil {
		d.SetId("")
		return nil
	}

	tags, err := tagService.DescribeResourceTags(ctx, VPC_SERVICE_TYPE, EIP_RESOURCE_TYPE, region, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}

	bgp, err := vpcService.DescribeVpcBandwidthPackageByEip(ctx, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}
	_ = d.Set("name", eip.AddressName)
	_ = d.Set("type", eip.AddressType)
	_ = d.Set("public_ip", eip.AddressIp)
	_ = d.Set("status", eip.AddressStatus)
	_ = d.Set("internet_charge_type", eip.InternetChargeType)
	_ = d.Set("tags", tags)

	if eip.Bandwidth != nil {
		_ = d.Set("internet_max_bandwidth_out", eip.Bandwidth)
	}

	if bgp != nil {
		_ = d.Set("bandwidth_package_id", bgp.BandwidthPackageId)
	}
	return nil
}

func resourceTencentCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	eipId := d.Id()

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"bandwidth_package_id",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_eip update on %s is not support yet", field)
		}
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		err := vpcService.ModifyEipName(ctx, eipId, name)
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
			bandwidthOut := v.(int)
			err := vpcService.ModifyEipBandwidthOut(ctx, eipId, bandwidthOut)
			if err != nil {
				return err
			}
			d.SetPartial("internet_max_bandwidth_out")
		}
	}

	if d.HasChange("prepaid_period") || d.HasChange("auto_renew_flag") {
		period := d.Get("prepaid_period").(int)
		renewFlag := d.Get("auto_renew_flag").(int)
		err := vpcService.RenewAddress(ctx, eipId, period, renewFlag)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName(VPC_SERVICE_TYPE, EIP_RESOURCE_TYPE, region, eipId)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			log.Printf("[CRITAL]%s update eip tags failed: %+v", logId, err)
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	eipId := d.Id()
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.UnattachEip(ctx, eipId)
		if errRet != nil {
			return retryError(errRet, "DesOperation.MutexTaskRunning")
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.DeleteEip(ctx, eipId)
		if errRet != nil {
			return retryError(errRet, "DesOperation.MutexTaskRunning")
		}
		return nil
	})
	if err != nil {
		return err
	}

	var internetChargeType string
	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetChargeType = v.(string)
	}

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		// isolated
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			eip, errRet := vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return retryError(errRet)
			}
			if !*eip.IsArrears {
				return resource.RetryableError(fmt.Errorf("eip is still isolate"))
			}
			return nil
		})
		if err != nil {
			return err
		}

		// release
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := vpcService.DeleteEip(ctx, eipId)
			if errRet != nil {
				return retryError(errRet, "DesOperation.MutexTaskRunning")
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet)
		}
		if eip != nil {
			return resource.RetryableError(fmt.Errorf("eip is still deleting"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
