/*
Provides a resource to create a mariadb reset_password

Example Usage

```hcl
resource "tencentcloud_mariadb_reset_password" "reset_password" {
  instance_id = ""
  user_name = ""
  host = ""
  password = ""
}
```

Import

mariadb reset_password can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_reset_password.reset_password reset_password_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudMariadbResetPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbResetPasswordCreate,
		Read:   resourceTencentCloudMariadbResetPasswordRead,
		Delete: resourceTencentCloudMariadbResetPasswordDelete,
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

func resourceTencentCloudMariadbResetPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_reset_password.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewResetAccountPasswordRequest()
		response   = mariadb.NewResetAccountPasswordResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ResetAccountPassword(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb resetPassword failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, userName, host}, FILED_SP))

	return resourceTencentCloudMariadbResetPasswordRead(d, meta)
}

func resourceTencentCloudMariadbResetPasswordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_reset_password.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbResetPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_reset_password.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
