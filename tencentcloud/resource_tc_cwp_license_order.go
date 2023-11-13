/*
Provides a resource to create a cwp license_order

Example Usage

```hcl
resource "tencentcloud_cwp_license_order" "license_order" {
  license_type =
  license_num =
  region_id =
  project_id =
}
```

Import

cwp license_order can be imported using the id, e.g.

```
terraform import tencentcloud_cwp_license_order.license_order license_order_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCwpLicenseOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCwpLicenseOrderCreate,
		Read:   resourceTencentCloudCwpLicenseOrderRead,
		Update: resourceTencentCloudCwpLicenseOrderUpdate,
		Delete: resourceTencentCloudCwpLicenseOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"license_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "LicenseType, 0 CWP Pro - Pay as you go ,1 CWP Pro - Monthly subscription , 2 CWP Ultimate - Monthly subscriptionDefault is 0.",
			},

			"license_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "License quantity , Quantity to be purchased.Default is 1.",
			},

			"region_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase order region, only 1 Guangzhou,9 Singapore is supported here.Guangzhou is recommended. Singapore is whitelisted.Default is 1.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID .Default is 0.",
			},
		},
	}
}

func resourceTencentCloudCwpLicenseOrderCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_order.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cwp.NewCreateLicenseOrderRequest()
		response   = cwp.NewCreateLicenseOrderResponse()
		resourceId string
	)
	if v, ok := d.GetOkExists("license_type"); ok {
		request.LicenseType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("license_num"); ok {
		request.LicenseNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("region_id"); ok {
		request.RegionId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCwpClient().CreateLicenseOrder(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cwp licenseOrder failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.ResourceId
	d.SetId(resourceId)

	return resourceTencentCloudCwpLicenseOrderRead(d, meta)
}

func resourceTencentCloudCwpLicenseOrderRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_order.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CwpService{client: meta.(*TencentCloudClient).apiV3Conn}

	licenseOrderId := d.Id()

	licenseOrder, err := service.DescribeCwpLicenseOrderById(ctx, resourceId)
	if err != nil {
		return err
	}

	if licenseOrder == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CwpLicenseOrder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if licenseOrder.LicenseType != nil {
		_ = d.Set("license_type", licenseOrder.LicenseType)
	}

	if licenseOrder.LicenseNum != nil {
		_ = d.Set("license_num", licenseOrder.LicenseNum)
	}

	if licenseOrder.RegionId != nil {
		_ = d.Set("region_id", licenseOrder.RegionId)
	}

	if licenseOrder.ProjectId != nil {
		_ = d.Set("project_id", licenseOrder.ProjectId)
	}

	return nil
}

func resourceTencentCloudCwpLicenseOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_order.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cwp.NewModifyLicenseOrderRequest()

	licenseOrderId := d.Id()

	request.ResourceId = &resourceId

	immutableArgs := []string{"license_type", "license_num", "region_id", "project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOkExists("project_id"); ok {
			request.ProjectId = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCwpClient().ModifyLicenseOrder(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cwp licenseOrder failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCwpLicenseOrderRead(d, meta)
}

func resourceTencentCloudCwpLicenseOrderDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_order.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CwpService{client: meta.(*TencentCloudClient).apiV3Conn}
	licenseOrderId := d.Id()

	if err := service.DeleteCwpLicenseOrderById(ctx, resourceId); err != nil {
		return err
	}

	return nil
}
