/*
Provides a resource to create a cdb root_account

Example Usage

```hcl
resource "tencentcloud_cdb_root_account" "root_account" {
  instance_id = ""
  password = ""
}
```

Import

cdb root_account can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_root_account.root_account root_account_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbRootAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRootAccountCreate,
		Read:   resourceTencentCloudCdbRootAccountRead,
		Delete: resourceTencentCloudCdbRootAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"password": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The password of the ROOT account of the instance.",
			},
		},
	}
}

func resourceTencentCloudCdbRootAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_root_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewVerifyRootAccountRequest()
		response   = cdb.NewVerifyRootAccountResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().VerifyRootAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb rootAccount failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbRootAccountStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbRootAccountRead(d, meta)
}

func resourceTencentCloudCdbRootAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_root_account.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbRootAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_root_account.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
