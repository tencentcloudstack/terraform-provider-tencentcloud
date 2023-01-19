/*
Provides a resource to create a cfs sign_up_cfs_service

Example Usage

```hcl
resource "tencentcloud_cfs_sign_up_cfs_service" "sign_up_cfs_service" {}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
)

func resourceTencentCloudCfsSignUpCfsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsSignUpCfsServiceCreate,
		Read:   resourceTencentCloudCfsSignUpCfsServiceRead,
		Delete: resourceTencentCloudCfsSignUpCfsServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cfs_service_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current status of the CFS service for this user. Valid values: creating (activating); created (activated).",
			},
		},
	}
}

func resourceTencentCloudCfsSignUpCfsServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfs_sign_up_cfs_service.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = cfs.NewSignUpCfsServiceRequest()
		response         = cfs.NewSignUpCfsServiceResponse()
		cfsServiceStatus string
	)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().SignUpCfsService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate cfs signUpCfsService failed, reason:%+v", logId, err)
		return nil
	}

	cfsServiceStatus = *response.Response.CfsServiceStatus
	d.SetId(cfsServiceStatus)

	return resourceTencentCloudCfsSignUpCfsServiceRead(d, meta)
}

func resourceTencentCloudCfsSignUpCfsServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_sign_up_cfs_service.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cfs.NewDescribeCfsServiceStatusRequest()
		response = cfs.NewDescribeCfsServiceStatusResponse()
	)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().DescribeCfsServiceStatus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Println("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Println("[CRITAL]%s operate cfs signUpCfsService failed, reason:%+v", logId, err)
		return nil
	}

	_ = d.Set("cfs_service_status", response.Response.CfsServiceStatus)

	return nil
}

func resourceTencentCloudCfsSignUpCfsServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_sign_up_cfs_service.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
