/*
Provides a resource to create a mariadb copy_account_privileges

Example Usage

```hcl
resource "tencentcloud_mariadb_copy_account_privileges" "copy_account_privileges" {
  instance_id   = "tdsql-9vqvls95"
  src_user_name = "keep-modify-privileges"
  src_host      = "127.0.0.1"
  dst_user_name = "keep-copy-user"
  dst_host      = "127.0.0.1"
}
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMariadbCopyAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbCopyAccountPrivilegesCreate,
		Read:   resourceTencentCloudMariadbCopyAccountPrivilegesRead,
		Delete: resourceTencentCloudMariadbCopyAccountPrivilegesDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, which is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.",
			},
			"src_user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Source username.",
			},
			"src_host": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Access host allowed for source user.",
			},
			"dst_user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target username.",
			},
			"dst_host": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Access host allowed for target user.",
			},
			"src_read_only": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "`ReadOnly` attribute of source account.",
			},
			"dst_read_only": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "`ReadOnly` attribute of target account.",
			},
		},
	}
}

func resourceTencentCloudMariadbCopyAccountPrivilegesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_copy_account_privileges.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = mariadb.NewCopyAccountPrivilegesRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("src_user_name"); ok {
		request.SrcUserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_host"); ok {
		request.SrcHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_read_only"); ok {
		request.SrcReadOnly = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_user_name"); ok {
		request.DstUserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_host"); ok {
		request.DstHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_read_only"); ok {
		request.DstReadOnly = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CopyAccountPrivileges(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb copyAccountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbCopyAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudMariadbCopyAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_copy_account_privileges.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbCopyAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_copy_account_privileges.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
