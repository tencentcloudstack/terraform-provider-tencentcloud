package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIpv6AddressBandwidth() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIpv6AddressBandwidthCreate,
		Read:   resourceTencentCloudIpv6AddressBandwidthRead,
		Update: resourceTencentCloudIpv6AddressBandwidthUpdate,
		Delete: resourceTencentCloudIpv6AddressBandwidthDelete,
		// it can support import because
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		Schema: map[string]*schema.Schema{
			"ipv6_address": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "IPV6 address that needs to be enabled for public network access.",
			},

			"internet_max_bandwidth_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     1,
				Description: "Bandwidth, in Mbps. The default is 1Mbps.",
			},

			"internet_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Network billing mode. IPV6 currently supports: `TRAFFIC_POSTPAID_BY_HOUR`, for standard account types; `BANDWIDTH_PACKAGE`, for traditional account types. The default network billing mode is: `TRAFFIC_POSTPAID_BY_HOUR`.",
			},

			"bandwidth_package_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The bandwidth package id, the Legacy account and the ipv6 address to apply for the bandwidth package charge type need to be passed in.",
			},
		},
	}
}

func resourceTencentCloudIpv6AddressBandwidthCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ipv6_address_bandwidth.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = vpc.NewAllocateIp6AddressesBandwidthRequest()
		response     = vpc.NewAllocateIp6AddressesBandwidthResponse()
		ip6AddressId string
	)
	if v, ok := d.GetOk("ipv6_address"); ok {
		ip6AddressId := helper.String(v.(string))
		request.Ip6Addresses = []*string{ip6AddressId}
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.BandwidthPackageId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AllocateIp6AddressesBandwidth(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ipv6Address failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.TaskId
	d.SetId(ip6AddressId)

	ip6AddressId = *response.Response.AddressSet[0]
	d.SetId(ip6AddressId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*tccommon.ReadRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(taskId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudIpv6AddressBandwidthRead(d, meta)
}

func resourceTencentCloudIpv6AddressBandwidthRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ipv6_address_bandwidth.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ipv6AddressId := d.Id()

	ipv6Address, err := service.DescribeVpcIpv6AddressById(ctx, ipv6AddressId)
	if err != nil {
		return err
	}

	if ipv6Address == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6Address` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ipv6_address", ipv6Address.AddressIp)

	if ipv6Address.Bandwidth != nil {
		_ = d.Set("internet_max_bandwidth_out", ipv6Address.Bandwidth)
	}

	if ipv6Address.InternetChargeType != nil {
		_ = d.Set("internet_charge_type", ipv6Address.InternetChargeType)
	}

	//if ipv6Address.BandwidthPackageId != nil {
	//	_ = d.Set("bandwidth_package_id", ipv6Address.BandwidthPackageId)
	//}

	return nil
}

func resourceTencentCloudIpv6AddressBandwidthUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ipv6_address_bandwidth.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyIp6AddressesBandwidthRequest()

	ipv6AddressId := d.Id()

	request.Ip6AddressIds = []*string{&ipv6AddressId}

	immutableArgs := []string{"internet_charge_type", "bandwidth_package_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false

	if d.HasChange("internet_max_bandwidth_out") {
		needChange = true
		if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
			request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}
	}

	if needChange {
		var taskId string

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyIp6AddressesBandwidth(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			taskId = *result.Response.TaskId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc ipv6Address failed, reason:%+v", logId, err)
			return err
		}

		service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*tccommon.ReadRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(taskId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}

	}

	return resourceTencentCloudIpv6AddressBandwidthRead(d, meta)
}

func resourceTencentCloudIpv6AddressBandwidthDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ipv6_address_bandwidth.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	ipv6AddressId := d.Id()

	if err := service.DeleteVpcIpv6AddressById(ctx, ipv6AddressId); err != nil {
		return err
	}

	return nil
}
