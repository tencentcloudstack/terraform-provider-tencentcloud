/*
Provides a resource to create a lighthouse modify_instance_renew_flag

Example Usage

```hcl
resource "tencentcloud_lighthouse_modify_instance_renew_flag" "modify_instance_renew_flag" {
  instance_ids =
  renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
```

Import

lighthouse modify_instance_renew_flag can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_modify_instance_renew_flag.modify_instance_renew_flag modify_instance_renew_flag_id
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

func resourceTencentCloudLighthouseModifyInstanceRenewFlag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseModifyInstanceRenewFlagCreate,
		Read:   resourceTencentCloudLighthouseModifyInstanceRenewFlagRead,
		Delete: resourceTencentCloudLighthouseModifyInstanceRenewFlagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
			},

			"renew_flag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Automatic renewal flag. Valid values are (NOTIFY_AND_AUTO_RENEW, NOTIFY_AND_MANUAL_RENEW, DISABLE_NOTIFY_AND_AUTO_RENEW).Default value: NOTIFY_AND_MANUAL_RENEW。If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis when the account balance is sufficient.",
			},
		},
	}
}

func resourceTencentCloudLighthouseModifyInstanceRenewFlagCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_modify_instance_renew_flag.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewModifyInstancesRenewFlagRequest()
		response   = lighthouse.NewModifyInstancesRenewFlagResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyInstancesRenewFlag(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse modifyInstanceRenewFlag failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudLighthouseModifyInstanceRenewFlagRead(d, meta)
}

func resourceTencentCloudLighthouseModifyInstanceRenewFlagRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_modify_instance_renew_flag.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseModifyInstanceRenewFlagDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_modify_instance_renew_flag.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
