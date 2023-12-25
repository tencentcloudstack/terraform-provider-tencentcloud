package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmRenewHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRenewHostCreate,
		Read:   resourceTencentCloudCvmRenewHostRead,
		Delete: resourceTencentCloudCvmRenewHostDelete,
		Schema: map[string]*schema.Schema{
			"host_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CDH instance ID.",
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
							Description: "Auto renewal flag. Valid values:&lt;br&gt;&lt;li&gt;NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically&lt;br&gt;&lt;li&gt;NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically&lt;br&gt;&lt;li&gt;DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically&lt;br&gt;&lt;br&gt;Default value: NOTIFY_AND_AUTO_RENEW。If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCvmRenewHostCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_host.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cvm.NewRenewHostsRequest()
	)
	hostId := d.Get("host_id").(string)
	request.HostIds = []*string{&hostId}

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RenewHosts(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm renewHost failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(hostId)

	return resourceTencentCloudCvmRenewHostRead(d, meta)
}

func resourceTencentCloudCvmRenewHostRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_host.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRenewHostDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_host.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
