/*
Provides a resource to create a lighthouse renew_disks

Example Usage

```hcl
resource "tencentcloud_lighthouse_renew_disks" "renew_disks" {
  disk_ids =
  renew_disk_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"
		time_unit = "m"
		cur_instance_deadline = "2018-01-01 00:00:00"

  }
  auto_voucher = true
}
```

Import

lighthouse renew_disks can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_renew_disks.renew_disks renew_disks_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudLighthouseRenewDisks() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseRenewDisksCreate,
		Read:   resourceTencentCloudLighthouseRenewDisksRead,
		Delete: resourceTencentCloudLighthouseRenewDisksDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disk_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of disk IDs. One or more cloud disk IDs to be operated. It can be obtained through the DiskId in the return value of the DescribeDisks interface. The total number of renewal data disks for each request is limited to 50.",
			},

			"renew_disk_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Renew cloud hard disk subscription related parameter settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Renewal period.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automatic renewal falg. Value:NOTIFY_AND_AUTO_RENEW: Notice expires and auto-renews.NOTIFY_AND_MANUAL_RENEW: Notification expires without automatic renewal, users need to manually renew.DISABLE_NOTIFY_AND_AUTO_RENEW: No automatic renewal and no notification.Default: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the disk will be automatically renewed monthly when the account balance is sufficient.",
						},
						"time_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Newly purchased unit. Default: m.",
						},
						"cur_instance_deadline": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Current instance expiration time. Such as 2018-01-01 00:00:00. Specifying this parameter can align the expiration time of the instance attached to the disk. One of this parameter and Period must be specified, and cannot be specified at the same time.",
						},
					},
				},
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically use the voucher. Not used by default.",
			},
		},
	}
}

func resourceTencentCloudLighthouseRenewDisksCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_renew_disks.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = lighthouse.NewRenewDisksRequest()
		response = lighthouse.NewRenewDisksResponse()
		diskId   string
	)
	if v, ok := d.GetOk("disk_ids"); ok {
		diskIdsSet := v.(*schema.Set).List()
		for i := range diskIdsSet {
			diskIds := diskIdsSet[i].(string)
			request.DiskIds = append(request.DiskIds, &diskIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "renew_disk_charge_prepaid"); ok {
		renewDiskChargePrepaid := lighthouse.RenewDiskChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			renewDiskChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			renewDiskChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		if v, ok := dMap["time_unit"]; ok {
			renewDiskChargePrepaid.TimeUnit = helper.String(v.(string))
		}
		if v, ok := dMap["cur_instance_deadline"]; ok {
			renewDiskChargePrepaid.CurInstanceDeadline = helper.String(v.(string))
		}
		request.RenewDiskChargePrepaid = &renewDiskChargePrepaid
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().RenewDisks(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse renewDisks failed, reason:%+v", logId, err)
		return err
	}

	diskId = *response.Response.DiskId
	d.SetId(diskId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseRenewDisksStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseRenewDisksRead(d, meta)
}

func resourceTencentCloudLighthouseRenewDisksRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_renew_disks.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseRenewDisksDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_renew_disks.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
