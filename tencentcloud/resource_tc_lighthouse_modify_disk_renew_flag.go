/*
Provides a resource to create a lighthouse modify_disk_renew_flag

Example Usage

```hcl
resource "tencentcloud_lighthouse_modify_disk_renew_flag" "modify_disk_renew_flag" {
  disk_ids =
  renew_flag = ""
}
```

Import

lighthouse modify_disk_renew_flag can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_modify_disk_renew_flag.modify_disk_renew_flag modify_disk_renew_flag_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLighthouseModifyDiskRenewFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseModifyDiskRenewFlagCreate,
		Read:   resourceTencentCloudLighthouseModifyDiskRenewFlagRead,
		Delete: resourceTencentCloudLighthouseModifyDiskRenewFlagDelete,
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
				Description: "List of cloud disk IDs.",
			},

			"renew_flag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Whether Auto-Renewal is enabled.",
			},
		},
	}
}

func resourceTencentCloudLighthouseModifyDiskRenewFlagCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_modify_disk_renew_flag.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewModifyDisksRenewFlagRequest()
		response   = lighthouse.NewModifyDisksRenewFlagResponse()
		instanceId string
	)
	if v, ok := d.GetOk("disk_ids"); ok {
		diskIdsSet := v.(*schema.Set).List()
		for i := range diskIdsSet {
			diskIds := diskIdsSet[i].(string)
			request.DiskIds = append(request.DiskIds, &diskIds)
		}
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyDisksRenewFlag(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse modifyDiskRenewFlag failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudLighthouseModifyDiskRenewFlagRead(d, meta)
}

func resourceTencentCloudLighthouseModifyDiskRenewFlagRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_modify_disk_renew_flag.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseModifyDiskRenewFlagDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_modify_disk_renew_flag.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
