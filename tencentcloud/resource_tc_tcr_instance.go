/*
Provides a resource to create a tcr instance

Example Usage

```hcl
resource "tencentcloud_tcr_instance" "instance" {
  registry_id = "tcr-xxx"
  registry_charge_prepaid {
		period = 1
		renew_flag = 0

  }
  flag = 0
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr instance can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_instance.instance instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTcrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrInstanceCreate,
		Read:   resourceTencentCloudTcrInstanceRead,
		Delete: resourceTencentCloudTcrInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Registry id.",
			},

			"registry_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prepaid automatic continuous fee identification and long purchase time, 0: manual continuous fee, 1: automatic continuous fee, 2: discontinuous fee without notification; the unit is month.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The duration of purchasing an instance, unit: month.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Automatic renewal flag, 0: manual renewal, 1: automatic renewal, 2: no renewal and no notification.",
						},
					},
				},
			},

			"flag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "0 Renewal, 1 Subscription according to the amount of the year and month.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tcr.NewRenewInstanceRequest()
		response   = tcr.NewRenewInstanceResponse()
		registryId string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "registry_charge_prepaid"); ok {
		registryChargePrepaid := tcr.RegistryChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			registryChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			registryChargePrepaid.RenewFlag = helper.IntInt64(v.(int))
		}
		request.RegistryChargePrepaid = &registryChargePrepaid
	}

	if v, _ := d.GetOk("flag"); v != nil {
		request.Flag = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().RenewInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr Instance failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(registryId)

	return resourceTencentCloudTcrInstanceRead(d, meta)
}

func resourceTencentCloudTcrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
