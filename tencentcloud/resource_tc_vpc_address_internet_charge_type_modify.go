/*
Provides a resource to create a vpc address_internet_charge_type_modify

Example Usage

```hcl
resource "tencentcloud_vpc_address_internet_charge_type_modify" "address_internet_charge_type_modify" {
  address_id = "eip-3456tghy"
  internet_charge_type = "BANDWIDTH_PREPAID_BY_MONTH"
  internet_max_bandwidth_out = 10
  address_charge_prepaid {
		period = 1
		auto_renew_flag = 1

  }
}
```

Import

vpc address_internet_charge_type_modify can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_address_internet_charge_type_modify.address_internet_charge_type_modify address_internet_charge_type_modify_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudVpcAddressInternetChargeTypeModify() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcAddressInternetChargeTypeModifyCreate,
		Read:   resourceTencentCloudVpcAddressInternetChargeTypeModifyRead,
		Delete: resourceTencentCloudVpcAddressInternetChargeTypeModifyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"address_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the elastic public network IP, in the form of eip-xxxxxxxx.",
			},

			"internet_charge_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The target charge type of elastic public network IP.",
			},

			"internet_max_bandwidth_out": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The target bandwidth value of elastic public network IP.",
			},

			"address_charge_prepaid": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Monthly bandwidth network charge type parameters. When the target charge type of the EIP is BANDWIDTH_PREPAID_BY_MONTH, this parameter must be passed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The duration of purchasing an instance, in months. Available duration: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Auto-renew flag. 0 indicates manual renewal, 1 indicates automatic renewal, and 2 indicates no renewal upon expiration. The default is 0.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcAddressInternetChargeTypeModifyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_address_internet_charge_type_modify.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = vpc.NewModifyAddressInternetChargeTypeRequest()
		response  = vpc.NewModifyAddressInternetChargeTypeResponse()
		addressId string
	)
	if v, ok := d.GetOk("address_id"); ok {
		addressId = v.(string)
		request.AddressId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetChargeType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("internet_max_bandwidth_out"); v != nil {
		request.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "address_charge_prepaid"); ok {
		addressChargePrepaid := vpc.AddressChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			addressChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["auto_renew_flag"]; ok {
			addressChargePrepaid.AutoRenewFlag = helper.IntInt64(v.(int))
		}
		request.AddressChargePrepaid = &addressChargePrepaid
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyAddressInternetChargeType(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc addressInternetChargeTypeModify failed, reason:%+v", logId, err)
		return err
	}

	addressId = *response.Response.AddressId
	d.SetId(addressId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"RUNNING"}, 1*readRetryTimeout, time.Second, service.VpcAddressInternetChargeTypeModifyStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudVpcAddressInternetChargeTypeModifyRead(d, meta)
}

func resourceTencentCloudVpcAddressInternetChargeTypeModifyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_address_internet_charge_type_modify.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcAddressInternetChargeTypeModifyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_address_internet_charge_type_modify.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
