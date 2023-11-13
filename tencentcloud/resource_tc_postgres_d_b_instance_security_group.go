/*
Provides a resource to create a postgres d_b_instance_security_group

Example Usage

```hcl
resource "tencentcloud_postgres_d_b_instance_security_group" "d_b_instance_security_group" {
  security_group_id_set =
  d_b_instance_id = ""
  read_only_group_id = ""
}
```

Import

postgres d_b_instance_security_group can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_d_b_instance_security_group.d_b_instance_security_group d_b_instance_security_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"log"
	"strings"
)

func resourceTencentCloudPostgresDBInstanceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresDBInstanceSecurityGroupCreate,
		Read:   resourceTencentCloudPostgresDBInstanceSecurityGroupRead,
		Update: resourceTencentCloudPostgresDBInstanceSecurityGroupUpdate,
		Delete: resourceTencentCloudPostgresDBInstanceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id_set": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Information of security groups in array.",
			},

			"d_b_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance ID. Either this parameter or ReadOnlyGroupId must be passed in. If both parameters are passed in, ReadOnlyGroupId will be ignored.",
			},

			"read_only_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "RO group ID. Either this parameter or DBInstanceId must be passed in. To query the security groups associated with the RO groups, only pass in ReadOnlyGroupId.",
			},
		},
	}
}

func resourceTencentCloudPostgresDBInstanceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance_security_group.create")()
	defer inconsistentCheck(d, meta)()

	var dBInstanceId string
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
	}

	var readOnlyGroupId string
	if v, ok := d.GetOk("read_only_group_id"); ok {
		readOnlyGroupId = v.(string)
	}

	d.SetId(strings.Join([]string{dBInstanceId, readOnlyGroupId}, FILED_SP))

	return resourceTencentCloudPostgresDBInstanceSecurityGroupUpdate(d, meta)
}

func resourceTencentCloudPostgresDBInstanceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance_security_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	readOnlyGroupId := idSplit[1]

	DBInstanceSecurityGroup, err := service.DescribePostgresDBInstanceSecurityGroupById(ctx, dBInstanceId, readOnlyGroupId)
	if err != nil {
		return err
	}

	if DBInstanceSecurityGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresDBInstanceSecurityGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DBInstanceSecurityGroup.SecurityGroupIdSet != nil {
		_ = d.Set("security_group_id_set", DBInstanceSecurityGroup.SecurityGroupIdSet)
	}

	if DBInstanceSecurityGroup.DBInstanceId != nil {
		_ = d.Set("d_b_instance_id", DBInstanceSecurityGroup.DBInstanceId)
	}

	if DBInstanceSecurityGroup.ReadOnlyGroupId != nil {
		_ = d.Set("read_only_group_id", DBInstanceSecurityGroup.ReadOnlyGroupId)
	}

	return nil
}

func resourceTencentCloudPostgresDBInstanceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance_security_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyDBInstanceSecurityGroupsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	readOnlyGroupId := idSplit[1]

	request.DBInstanceId = &dBInstanceId
	request.ReadOnlyGroupId = &readOnlyGroupId

	immutableArgs := []string{"security_group_id_set", "d_b_instance_id", "read_only_group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyDBInstanceSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres DBInstanceSecurityGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresDBInstanceSecurityGroupRead(d, meta)
}

func resourceTencentCloudPostgresDBInstanceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance_security_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
