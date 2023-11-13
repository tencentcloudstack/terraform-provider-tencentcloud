/*
Provides a resource to create a dcdb reset_account_password

Example Usage

```hcl
resource "tencentcloud_dcdb_reset_account_password" "reset_account_password" {
  instance_id = ""
  user_name = ""
  host = ""
  password = ""
}
```

Import

dcdb reset_account_password can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_reset_account_password.reset_account_password reset_account_password_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDcdbResetAccountPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbResetAccountPasswordCreate,
		Read:   resourceTencentCloudDcdbResetAccountPasswordRead,
		Delete: resourceTencentCloudDcdbResetAccountPasswordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, which is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.",
			},

			"user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Login username.",
			},

			"host": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Access host allowed for user. An account is uniquely identified by username and host.",
			},

			"password": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "New password, which can contain 6-32 letters, digits, and common symbols but not semicolons, single quotation marks, and double quotation marks.",
			},
		},
	}
}

func resourceTencentCloudDcdbResetAccountPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_reset_account_password.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewResetAccountPasswordRequest()
		response   = dcdb.NewResetAccountPasswordResponse()
		instanceId string
		userName   string
		host       string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		userName = v.(string)
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ResetAccountPassword(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb resetAccountPassword failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, userName, host}, FILED_SP))

	return resourceTencentCloudDcdbResetAccountPasswordRead(d, meta)
}

func resourceTencentCloudDcdbResetAccountPasswordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_reset_account_password.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbResetAccountPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_reset_account_password.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
