/*
Provides a resource to create a cvm reboot_instance

Example Usage

```hcl
resource "tencentcloud_cvm_reboot_instance" "reboot_instance" {
  instance_ids =
  force_reboot = false
  stop_type = "SOFT"
}
```

Import

cvm reboot_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_reboot_instance.reboot_instance reboot_instance_id
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

func resourceTencentCloudCvmRebootInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRebootInstanceCreate,
		Read:   resourceTencentCloudCvmRebootInstanceRead,
		Delete: resourceTencentCloudCvmRebootInstanceDelete,
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
				Description: "Instance IDs. To obtain the instance IDs, you can call DescribeInstances and look for InstanceId in the response. You can operate up to 100 instances in each request.",
			},

			"force_reboot": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "This parameter has been disused. We recommend using StopType instead. Note that ForceReboot and StopType parameters cannot be specified at the same time. Whether to forcibly restart an instance after a normal restart fails. Valid values:&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;TRUE: yes;&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;FALSE: no&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;Default value: FALSE. .",
			},

			"stop_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Shutdown type. Valid values:&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;SOFT：soft shutdown&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;HARD：hard shutdown&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;SOFT_FIRST：perform a soft shutdown first, and perform a hard shutdown if the soft shutdown fails&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;br&amp;amp;gt;Default value: SOFT.",
			},
		},
	}
}

func resourceTencentCloudCvmRebootInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_reboot_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cvm.NewRebootInstancesRequest()
		response   = cvm.NewRebootInstancesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, _ := d.GetOk("force_reboot"); v != nil {
		request.ForceReboot = helper.Bool(v.(bool))
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm rebootInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
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
