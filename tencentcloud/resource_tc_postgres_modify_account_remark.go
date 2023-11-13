/*
Provides a resource to create a postgres modify_account_remark

Example Usage

```hcl
resource "tencentcloud_postgres_modify_account_remark" "modify_account_remark" {
  d_b_instance_id = ""
  user_name = ""
  remark = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres modify_account_remark can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_modify_account_remark.modify_account_remark modify_account_remark_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPostgresModifyAccountRemark() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresModifyAccountRemarkCreate,
		Read:   resourceTencentCloudPostgresModifyAccountRemarkRead,
		Delete: resourceTencentCloudPostgresModifyAccountRemarkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of postgres-4wdeb0zv.",
			},

			"user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance username.",
			},

			"remark": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "New remarks corresponding to user `UserName`.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresModifyAccountRemarkCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_account_remark.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewModifyAccountRemarkRequest()
		response     = postgres.NewModifyAccountRemarkResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyAccountRemark(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres ModifyAccountRemark failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresModifyAccountRemarkRead(d, meta)
}

func resourceTencentCloudPostgresModifyAccountRemarkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_account_remark.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresModifyAccountRemarkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_account_remark.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
