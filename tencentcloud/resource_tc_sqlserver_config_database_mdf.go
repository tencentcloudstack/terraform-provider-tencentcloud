/*
Provides a resource to create a sqlserver config_database_mdf

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_database_mdf" "config_database_mdf" {
  d_b_names =
  instance_id = "mssql-i1z41iwd"
}
```

Import

sqlserver config_database_mdf can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_mdf.config_database_mdf config_database_mdf_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
)

func resourceTencentCloudSqlserverConfigDatabaseMdf() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigDatabaseMdfCreate,
		Read:   resourceTencentCloudSqlserverConfigDatabaseMdfRead,
		Update: resourceTencentCloudSqlserverConfigDatabaseMdfUpdate,
		Delete: resourceTencentCloudSqlserverConfigDatabaseMdfDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of database names.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigDatabaseMdfCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_mdf.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigDatabaseMdfUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseMdfRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_mdf.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configDatabaseMdfId := d.Id()

	configDatabaseMdf, err := service.DescribeSqlserverConfigDatabaseMdfById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configDatabaseMdf == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigDatabaseMdf` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configDatabaseMdf.DBNames != nil {
		_ = d.Set("d_b_names", configDatabaseMdf.DBNames)
	}

	if configDatabaseMdf.InstanceId != nil {
		_ = d.Set("instance_id", configDatabaseMdf.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigDatabaseMdfUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_mdf.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDatabaseMdfRequest()

	configDatabaseMdfId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"d_b_names", "instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDatabaseMdf(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configDatabaseMdf failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigDatabaseMdfRead(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseMdfDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_database_mdf.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
