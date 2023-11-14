/*
Provides a resource to create a vpc ipv6_address

Example Usage

```hcl
resource "tencentcloud_vpc_ipv6_address" "ipv6_address" {
  ip6_addresses =
  internet_max_bandwidth_out = 200
  internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  bandwidth_package_id = "bwp-34rfgt56"
}
```

Import

vpc ipv6_address can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_address.ipv6_address ipv6_address_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudVpcIpv6Address() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcIpv6AddressCreate,
		Read:   resourceTencentCloudVpcIpv6AddressRead,
		Update: resourceTencentCloudVpcIpv6AddressUpdate,
		Delete: resourceTencentCloudVpcIpv6AddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ip6_addresses": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IPV6 address that needs to be enabled for public network access.",
			},

			"internet_max_bandwidth_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Bandwidth, in Mbps. The default is 1Mbps.",
			},

			"internet_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Network billing mode. IPV6 currently supports &amp;amp;#39;TRAFFIC_POSTPAID_BY_HOUR&amp;amp;#39; for standard account types and &amp;amp;#39;BANDWIDTH_PACKAGE&amp;amp;#39; for traditional account types. The default network billing mode is &amp;amp;#39;TRAFFIC_POSTPAID_BY_HOUR&amp;amp;#39;.",
			},

			"bandwidth_package_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The bandwidth package id, the Legacy account and the ipv6 address to apply for the bandwidth package charge type need to be passed in.",
			},
		},
	}
}

func resourceTencentCloudVpcIpv6AddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_address.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = vpc.NewAllocateIp6AddressesBandwidthRequest()
		response = vpc.NewAllocateIp6AddressesBandwidthResponse()
	)
	if v, ok := d.GetOk("ip6_addresses"); ok {
		ip6AddressesSet := v.(*schema.Set).List()
		for i := range ip6AddressesSet {
			ip6Addresses := ip6AddressesSet[i].(string)
			request.Ip6Addresses = append(request.Ip6Addresses, &ip6Addresses)
		}
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AllocateIp6AddressesBandwidth(request)
		if e != nil {
			return retryError(e)
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

	ip6AddressIds = *response.Response.Ip6AddressIds
	d.SetId()

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"RUNNING"}, 1*readRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudVpcIpv6AddressRead(d, meta)
}

func resourceTencentCloudVpcIpv6AddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	ipv6AddressId := d.Id()

	ipv6Address, err := service.DescribeVpcIpv6AddressById(ctx, ip6AddressIds)
	if err != nil {
		return err
	}

	if ipv6Address == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6Address` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ipv6Address.Ip6Addresses != nil {
		_ = d.Set("ip6_addresses", ipv6Address.Ip6Addresses)
	}

	if ipv6Address.InternetMaxBandwidthOut != nil {
		_ = d.Set("internet_max_bandwidth_out", ipv6Address.InternetMaxBandwidthOut)
	}

	if ipv6Address.InternetChargeType != nil {
		_ = d.Set("internet_charge_type", ipv6Address.InternetChargeType)
	}

	if ipv6Address.BandwidthPackageId != nil {
		_ = d.Set("bandwidth_package_id", ipv6Address.BandwidthPackageId)
	}

	return nil
}

func resourceTencentCloudVpcIpv6AddressUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_address.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyIp6AddressesBandwidthRequest()

	ipv6AddressId := d.Id()

	request.Ip6AddressIds = &ip6AddressIds

	immutableArgs := []string{"ip6_addresses", "internet_max_bandwidth_out", "internet_charge_type", "bandwidth_package_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("ip6_addresses") {
		if v, ok := d.GetOk("ip6_addresses"); ok {
			ip6AddressesSet := v.(*schema.Set).List()
			for i := range ip6AddressesSet {
				ip6Addresses := ip6AddressesSet[i].(string)
				request.Ip6Addresses = append(request.Ip6Addresses, &ip6Addresses)
			}
		}
	}

	if d.HasChange("internet_max_bandwidth_out") {
		if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
			request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyIp6AddressesBandwidth(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc ipv6Address failed, reason:%+v", logId, err)
		return err
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"RUNNING"}, 1*readRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudVpcIpv6AddressRead(d, meta)
}

func resourceTencentCloudVpcIpv6AddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	ipv6AddressId := d.Id()

	if err := service.DeleteVpcIpv6AddressById(ctx, ip6AddressIds); err != nil {
		return err
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"RUNNING"}, 1*readRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
