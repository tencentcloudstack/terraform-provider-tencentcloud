/*
Provides a resource to create a cvm renew_hosts

Example Usage

```hcl
resource "tencentcloud_cvm_renew_hosts" "renew_hosts" {
  host_ids =
  host_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  }
}
```

Import

cvm renew_hosts can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_renew_hosts.renew_hosts renew_hosts_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCvmRenewHosts() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRenewHostsCreate,
		Read:   resourceTencentCloudCvmRenewHostsRead,
		Delete: resourceTencentCloudCvmRenewHostsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"host_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "One or more CDH instance IDs to be operated on. The upper limit of CDH instances per request is 100.",
			},

			"host_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prepaid mode, that is, yearly and monthly subscription related parameter settings. Through this parameter, you can specify attributes such as the purchase duration of the Subscription instance and whether to set automatic renewal. If the payment mode of the specified instance is prepaid, this parameter must be passed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The duration of purchasing an instance, unit: month. Value range: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auto renewal flag. Valid values:&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically&amp;lt;br&amp;gt;&amp;lt;br&amp;gt;Default value: NOTIFY_AND_AUTO_RENEW。If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCvmRenewHostsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_renew_hosts.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cvm.NewRenewHostsRequest()
		response = cvm.NewRenewHostsResponse()
		hostId   string
	)
	if v, ok := d.GetOk("host_ids"); ok {
		hostIdsSet := v.(*schema.Set).List()
		for i := range hostIdsSet {
			hostIds := hostIdsSet[i].(string)
			request.HostIds = append(request.HostIds, &hostIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "host_charge_prepaid"); ok {
		chargePrepaid := cvm.ChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			chargePrepaid.Period = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			chargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.HostChargePrepaid = &chargePrepaid
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().RenewHosts(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm renewHosts failed, reason:%+v", logId, err)
		return err
	}

	hostId = *response.Response.HostId
	d.SetId(hostId)

	return resourceTencentCloudCvmRenewHostsRead(d, meta)
}

func resourceTencentCloudCvmRenewHostsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_renew_hosts.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRenewHostsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_renew_hosts.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
