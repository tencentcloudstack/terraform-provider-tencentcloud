/*
Provides a resource to create a cynosdb cluster_databases

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  db_name       = "terraform-test"
  character_set = "utf8"
  collate_rule  = "utf8_general_ci"
  user_host_privileges {
    db_user_name = "root"
    db_host      = "%"
    db_privilege = "READ_ONLY"
  }
  description = "terraform test"
}
```

Import

cynosdb cluster_databases can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_databases.cluster_databases cluster_databases_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbClusterDatabases() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterDatabasesCreate,
		Read:   resourceTencentCloudCynosdbClusterDatabasesRead,
		Update: resourceTencentCloudCynosdbClusterDatabasesUpdate,
		Delete: resourceTencentCloudCynosdbClusterDatabasesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},

			"character_set": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Character Set Type.",
			},

			"collate_rule": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sort Rules.",
			},

			"user_host_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Authorize user host permissions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authorized Users.",
						},
						"db_host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: ".",
						},
						"db_privilege": {
							Type:        schema.TypeString,
							Required:    true,
							Description: ".",
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterDatabasesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_databases.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewCreateClusterDatabaseRequest()
		clusterId string
		dbName    string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		dbName = v.(string)
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("character_set"); ok {
		request.CharacterSet = helper.String(v.(string))
	}

	if v, ok := d.GetOk("collate_rule"); ok {
		request.CollateRule = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_host_privileges"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			userHostPrivilege := cynosdb.UserHostPrivilege{}
			if v, ok := dMap["db_user_name"]; ok {
				userHostPrivilege.DbUserName = helper.String(v.(string))
			}
			if v, ok := dMap["db_host"]; ok {
				userHostPrivilege.DbHost = helper.String(v.(string))
			}
			if v, ok := dMap["db_privilege"]; ok {
				userHostPrivilege.DbPrivilege = helper.String(v.(string))
			}
			request.UserHostPrivileges = append(request.UserHostPrivileges, &userHostPrivilege)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateClusterDatabase(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterDatabases failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + dbName)

	return resourceTencentCloudCynosdbClusterDatabasesRead(d, meta)
}

func resourceTencentCloudCynosdbClusterDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_databases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	dbName := idSplit[1]

	clusterDatabases, err := service.DescribeCynosdbClusterDatabasesById(ctx, clusterId, dbName)
	if err != nil {
		return err
	}

	if clusterDatabases == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbClusterDatabases` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterDatabases.ClusterId != nil {
		_ = d.Set("cluster_id", clusterDatabases.ClusterId)
	}

	if clusterDatabases.DbName != nil {
		_ = d.Set("db_name", clusterDatabases.DbName)
	}

	if clusterDatabases.CharacterSet != nil {
		_ = d.Set("character_set", clusterDatabases.CharacterSet)
	}

	if clusterDatabases.CollateRule != nil {
		_ = d.Set("collate_rule", clusterDatabases.CollateRule)
	}

	if clusterDatabases.UserHostPrivileges != nil {
		userHostPrivilegesList := []interface{}{}
		for _, userHostPrivileges := range clusterDatabases.UserHostPrivileges {
			userHostPrivilegesMap := map[string]interface{}{}

			if userHostPrivileges.DbUserName != nil {
				userHostPrivilegesMap["db_user_name"] = userHostPrivileges.DbUserName
			}

			if userHostPrivileges.DbHost != nil {
				userHostPrivilegesMap["db_host"] = userHostPrivileges.DbHost
			}

			if userHostPrivileges.DbPrivilege != nil {
				userHostPrivilegesMap["db_privilege"] = userHostPrivileges.DbPrivilege
			}

			userHostPrivilegesList = append(userHostPrivilegesList, userHostPrivilegesMap)
		}

		_ = d.Set("user_host_privileges", userHostPrivilegesList)

	}

	if clusterDatabases.Description != nil {
		_ = d.Set("description", clusterDatabases.Description)
	}

	return nil
}

func resourceTencentCloudCynosdbClusterDatabasesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_databases.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyClusterDatabaseRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	dbName := idSplit[1]

	request.ClusterId = &clusterId
	request.DbName = &dbName

	immutableArgs := []string{"cluster_id", "character_set", "collate_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("user_host_privileges") {
		oldPrivileges, privileges := d.GetChange("user_host_privileges")
		for _, item := range privileges.([]interface{}) {
			dMap := item.(map[string]interface{})
			userHostPrivilege := cynosdb.UserHostPrivilege{}
			if v, ok := dMap["db_user_name"]; ok {
				userHostPrivilege.DbUserName = helper.String(v.(string))
			}
			if v, ok := dMap["db_host"]; ok {
				userHostPrivilege.DbHost = helper.String(v.(string))
			}
			if v, ok := dMap["db_privilege"]; ok {
				userHostPrivilege.DbPrivilege = helper.String(v.(string))
			}
			request.NewUserHostPrivileges = append(request.NewUserHostPrivileges, &userHostPrivilege)
		}

		for _, item := range oldPrivileges.([]interface{}) {
			dMap := item.(map[string]interface{})
			userHostPrivilege := cynosdb.UserHostPrivilege{}
			if v, ok := dMap["db_user_name"]; ok {
				userHostPrivilege.DbUserName = helper.String(v.(string))
			}
			if v, ok := dMap["db_host"]; ok {
				userHostPrivilege.DbHost = helper.String(v.(string))
			}
			if v, ok := dMap["db_privilege"]; ok {
				userHostPrivilege.DbPrivilege = helper.String(v.(string))
			}
			request.OldUserHostPrivileges = append(request.OldUserHostPrivileges, &userHostPrivilege)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyClusterDatabase(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb clusterDatabases failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbClusterDatabasesRead(d, meta)
}

func resourceTencentCloudCynosdbClusterDatabasesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_databases.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	dbName := idSplit[1]

	if err := service.DeleteCynosdbClusterDatabasesById(ctx, clusterId, dbName); err != nil {
		return err
	}

	return nil
}
