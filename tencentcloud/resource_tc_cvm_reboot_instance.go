package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmRebootInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRebootInstanceCreate,
		Read:   resourceTencentCloudCvmRebootInstanceRead,
		Delete: resourceTencentCloudCvmRebootInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"force_reboot": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeBool,
				ConflictsWith: []string{"stop_type"},
				Deprecated:    "It has been deprecated from version 1.81.21. Please use `stop_type` instead.",
				Description:   "This parameter has been disused. We recommend using StopType instead. Note that ForceReboot and StopType parameters cannot be specified at the same time. Whether to forcibly restart an instance after a normal restart fails. Valid values are `TRUE` and `FALSE`. Default value: FALSE.",
			},

			"stop_type": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"force_reboot"},
				Description:   "Shutdown type. Valid values: `SOFT`: soft shutdown; `HARD`: hard shutdown; `SOFT_FIRST`: perform a soft shutdown first, and perform a hard shutdown if the soft shutdown fails. Default value: SOFT.",
			},
		},
	}
}

func resourceTencentCloudCvmRebootInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_reboot_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cvm.NewRebootInstancesRequest()
	instanceId := d.Get("instance_id").(string)
	request.InstanceIds = []*string{&instanceId}

	if v, _ := d.GetOk("force_reboot"); v != nil {
		if _, ok := d.GetOk("stop_type"); !ok {
			request.ForceReboot = helper.Bool(v.(bool))
		}
	}

	if v, ok := d.GetOk("stop_type"); ok {
		request.StopType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().RebootInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm rebootInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudCvmRebootInstanceRead(d, meta)
}

func resourceTencentCloudCvmRebootInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_reboot_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRebootInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_reboot_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
