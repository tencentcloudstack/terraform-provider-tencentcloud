/*
Provides a resource to create a sqlserver config_instance_security_groups

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_security_groups" "config_instance_security_groups" {
  instance_id = "mssql-i1z41iwd"
  security_group_id_set =
}
```

Import

sqlserver config_instance_security_groups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_security_groups.config_instance_security_groups config_instance_security_groups_id
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

func resourceTencentCloudSqlserverConfigInstanceSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigInstanceSecurityGroupsCreate,
		Read:   resourceTencentCloudSqlserverConfigInstanceSecurityGroupsRead,
		Update: resourceTencentCloudSqlserverConfigInstanceSecurityGroupsUpdate,
		Delete: resourceTencentCloudSqlserverConfigInstanceSecurityGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"security_group_id_set": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of security group IDs to modify, an array of one or more security group IDs.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigInstanceSecurityGroupsUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configInstanceSecurityGroupsId := d.Id()

	configInstanceSecurityGroups, err := service.DescribeSqlserverConfigInstanceSecurityGroupsById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configInstanceSecurityGroups == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceSecurityGroups` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configInstanceSecurityGroups.InstanceId != nil {
		_ = d.Set("instance_id", configInstanceSecurityGroups.InstanceId)
	}

	if configInstanceSecurityGroups.SecurityGroupIdSet != nil {
		_ = d.Set("security_group_id_set", configInstanceSecurityGroups.SecurityGroupIdSet)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDBInstanceSecurityGroupsRequest()

	configInstanceSecurityGroupsId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "security_group_id_set"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBInstanceSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configInstanceSecurityGroups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigInstanceSecurityGroupsRead(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceSecurityGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_security_groups.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
