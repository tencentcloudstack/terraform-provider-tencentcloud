/*
Provides a resource to create a cynosdb cluster_databases

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  db_name = "test"
  character_set = "utf8"
  collate_rule = " utf8_general_ci "
  user_host_privileges {
		db_user_name = ""
		db_host = ""
		db_privilege = ""

  }
  description = "test"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
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
							Description: "Client IP Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"db_privilege": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User permission note: This field may return null, indicating that a valid value cannot be obtained.",
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
		response  = cynosdb.NewCreateClusterDatabaseResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterDatabases failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudCynosdbClusterDatabasesRead(d, meta)
}

func resourceTencentCloudCynosdbClusterDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_databases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterDatabasesId := d.Id()

	clusterDatabases, err := service.DescribeCynosdbClusterDatabasesById(ctx, clusterId)
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

			if clusterDatabases.UserHostPrivileges.DbUserName != nil {
				userHostPrivilegesMap["db_user_name"] = clusterDatabases.UserHostPrivileges.DbUserName
			}

			if clusterDatabases.UserHostPrivileges.DbHost != nil {
				userHostPrivilegesMap["db_host"] = clusterDatabases.UserHostPrivileges.DbHost
			}

			if clusterDatabases.UserHostPrivileges.DbPrivilege != nil {
				userHostPrivilegesMap["db_privilege"] = clusterDatabases.UserHostPrivileges.DbPrivilege
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

	clusterDatabasesId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "db_name", "character_set", "collate_rule", "user_host_privileges", "description"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("db_name") {
		if v, ok := d.GetOk("db_name"); ok {
			request.DbName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
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
	clusterDatabasesId := d.Id()

	if err := service.DeleteCynosdbClusterDatabasesById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
