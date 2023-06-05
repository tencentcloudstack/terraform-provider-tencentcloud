/*
Provides a resource to create a mariadb grant_account_privileges

Example Usage

```hcl
resource "tencentcloud_mariadb_grant_account_privileges" "grant_account_privileges" {
  instance_id = "tdsql-9vqvls95"
  user_name   = "keep-modify-privileges"
  host        = "127.0.0.1"
  db_name     = "*"
  privileges  = ["SELECT", "INSERT"]
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

func resourceTencentCloudMariadbGrantAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbGrantAccountPrivilegesCreate,
		Read:   resourceTencentCloudMariadbGrantAccountPrivilegesRead,
		Delete: resourceTencentCloudMariadbGrantAccountPrivilegesDelete,

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
			"db_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database name. `*` indicates that global permissions will be set (i.e., `*.*`), in which case the `Type` and `Object ` parameters will be ignored. If `DbName` is not `*`, the input parameter `Type` is required.",
			},
			"privileges": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Global permission: SELECT，INSERT，UPDATE，DELETE，CREATE，DROP，REFERENCES，INDEX，ALTER，CREATE TEMPORARY TABLES，LOCK TABLES，EXECUTE，CREATE VIEW，SHOW VIEW，CREATE ROUTINE，ALTER ROUTINE，EVENT，TRIGGER，SHOW DATABASES，REPLICATION CLIENT，REPLICATION SLAVE 库权限： SELECT，INSERT，UPDATE，DELETE，CREATE，DROP，REFERENCES，INDEX，ALTER，CREATE TEMPORARY TABLES，LOCK TABLES，EXECUTE，CREATE VIEW，SHOW VIEW，CREATE ROUTINE，ALTER ROUTINE，EVENT，TRIGGER 表/视图权限： SELECT，INSERT，UPDATE，DELETE，CREATE，DROP，REFERENCES，INDEX，ALTER，CREATE VIEW，SHOW VIEW，TRIGGER 存储过程/函数权限： ALTER ROUTINE，EXECUTE 字段权限： INSERT，REFERENCES，SELECT，UPDATE.",
			},
			"type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Type. Valid values: table, view, proc, func, *. If `DbName` is a specific database name and `Type` is `*`, the permissions of the database will be set (i.e., `db.*`), in which case the `Object` parameter will be ignored.",
			},
			"object": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Type name. For example, if `Type` is `table`, `Object` indicates a specific table name; if both `DbName` and `Type` are specific names, it indicates a specific object name and cannot be `*` or empty.",
			},
			"col_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "If `Type` is `table` and `ColName` is `*`, the permissions will be granted to the table; if `ColName` is a specific field name, the permissions will be granted to the field.",
			},
		},
	}
}

func resourceTencentCloudMariadbGrantAccountPrivilegesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_grant_account_privileges.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = mariadb.NewGrantAccountPrivilegesRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("privileges"); ok {
		privilegesSet := v.(*schema.Set).List()
		for i := range privilegesSet {
			privileges := privilegesSet[i].(string)
			request.Privileges = append(request.Privileges, &privileges)
		}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("object"); ok {
		request.Object = helper.String(v.(string))
	}

	if v, ok := d.GetOk("col_name"); ok {
		request.ColName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().GrantAccountPrivileges(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb grantAccountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbGrantAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudMariadbGrantAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_grant_account_privileges.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbGrantAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_grant_account_privileges.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
