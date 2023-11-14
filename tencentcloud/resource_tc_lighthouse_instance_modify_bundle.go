/*
Provides a resource to create a lighthouse instance_modify_bundle

Example Usage

```hcl
resource "tencentcloud_lighthouse_instance_modify_bundle" "instance_modify_bundle" {
  instance_ids =
  bundle_id = "bundle_gen_03"
  auto_voucher = true
}
```

Import

lighthouse instance_modify_bundle can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_instance_modify_bundle.instance_modify_bundle instance_modify_bundle_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"log"
	"time"
)

func resourceTencentCloudLighthouseInstanceModifyBundle() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseInstanceModifyBundleCreate,
		Read:   resourceTencentCloudLighthouseInstanceModifyBundleRead,
		Update: resourceTencentCloudLighthouseInstanceModifyBundleUpdate,
		Delete: resourceTencentCloudLighthouseInstanceModifyBundleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of instance IDs.",
			},

			"bundle_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bundle ID.",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically deduct vouchers. Valid values:true：Automatically deduct vouchers.false：Do not automatically deduct vouchers.Default value: false.",
			},
		},
	}
}

func resourceTencentCloudLighthouseInstanceModifyBundleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_modify_bundle.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudLighthouseInstanceModifyBundleUpdate(d, meta)
}

func resourceTencentCloudLighthouseInstanceModifyBundleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_modify_bundle.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceModifyBundleId := d.Id()

	instanceModifyBundle, err := service.DescribeLighthouseInstanceModifyBundleById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceModifyBundle == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseInstanceModifyBundle` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceModifyBundle.InstanceIds != nil {
		_ = d.Set("instance_ids", instanceModifyBundle.InstanceIds)
	}

	if instanceModifyBundle.BundleId != nil {
		_ = d.Set("bundle_id", instanceModifyBundle.BundleId)
	}

	if instanceModifyBundle.AutoVoucher != nil {
		_ = d.Set("auto_voucher", instanceModifyBundle.AutoVoucher)
	}

	return nil
}

func resourceTencentCloudLighthouseInstanceModifyBundleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_modify_bundle.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifyInstancesBundleRequest()

	instanceModifyBundleId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_ids", "bundle_id", "auto_voucher"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyInstancesBundle(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse instanceModifyBundle failed, reason:%+v", logId, err)
		return err
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseInstanceModifyBundleStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseInstanceModifyBundleRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceModifyBundleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance_modify_bundle.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
