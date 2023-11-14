/*
Provides a resource to create a cvm chc_attribute

Example Usage

```hcl
resource "tencentcloud_cvm_chc_attribute" "chc_attribute" {
  chc_ids =
  instance_name = ""
  device_type = ""
  bmc_user = ""
  password = ""
  bmc_security_group_ids =
}
```

Import

cvm chc_attribute can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_attribute.chc_attribute chc_attribute_id
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

func resourceTencentCloudCvmChcAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmChcAttributeCreate,
		Read:   resourceTencentCloudCvmChcAttributeRead,
		Delete: resourceTencentCloudCvmChcAttributeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"chc_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CHC host IDs.",
			},

			"instance_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CHC host name.",
			},

			"device_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Server type.",
			},

			"bmc_user": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Valid characters: Letters, numbers, hyphens and underscores.",
			},

			"password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&amp;amp;amp;*-+=_|{}).",
			},

			"bmc_security_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "BMC network security group list.",
			},
		},
	}
}

func resourceTencentCloudCvmChcAttributeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_attribute.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cvm.NewModifyChcAttributeRequest()
		response = cvm.NewModifyChcAttributeResponse()
		chcId    string
	)
	if v, ok := d.GetOk("chc_ids"); ok {
		chcIdsSet := v.(*schema.Set).List()
		for i := range chcIdsSet {
			chcIds := chcIdsSet[i].(string)
			request.ChcIds = append(request.ChcIds, &chcIds)
		}
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("device_type"); ok {
		request.DeviceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bmc_user"); ok {
		request.BmcUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bmc_security_group_ids"); ok {
		bmcSecurityGroupIdsSet := v.(*schema.Set).List()
		for i := range bmcSecurityGroupIdsSet {
			bmcSecurityGroupIds := bmcSecurityGroupIdsSet[i].(string)
			request.BmcSecurityGroupIds = append(request.BmcSecurityGroupIds, &bmcSecurityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
		return err
	}

	chcId = *response.Response.ChcId
	d.SetId(chcId)

	return resourceTencentCloudCvmChcAttributeRead(d, meta)
}

func resourceTencentCloudCvmChcAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_attribute.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmChcAttributeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_attribute.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
