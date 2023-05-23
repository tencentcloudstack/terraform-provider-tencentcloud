/*
Provides a resource to create a postgresql security_group_config

Example Usage

Set security group for the sepcified postgres instance
```hcl
resource "tencentcloud_postgresql_security_group_config" "security_group_config" {
  security_group_id_set = [local.sg_id, local.sg_id2]
  db_instance_id = local.pgsql_id
}
```

Set security group for the specified readonly group
```hcl
resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_test_ro_sg"
	project_id = 0
	subnet_id             = local.subnet_id
	vpc_id                = local.vpc_id
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
  }

resource "tencentcloud_postgresql_security_group_config" "security_group_config" {
  security_group_id_set = [local.sg_id, local.sg_id2]
  read_only_group_id = tencentcloud_postgresql_readonly_group.group.id
}
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
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlSecurityGroupConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlSecurityGroupConfigCreate,
		Read:   resourceTencentCloudPostgresqlSecurityGroupConfigRead,
		Update: resourceTencentCloudPostgresqlSecurityGroupConfigUpdate,
		Delete: resourceTencentCloudPostgresqlSecurityGroupConfigDelete,
		Schema: map[string]*schema.Schema{
			"security_group_id_set": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Information of security groups in array.",
			},

			"db_instance_id": {
				Optional:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"read_only_group_id"},
				Description:   "Instance ID. Either this parameter or ReadOnlyGroupId must be passed in. If both parameters are passed in, ReadOnlyGroupId will be ignored.",
			},

			"read_only_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "RO group ID. Either this parameter or DBInstanceId must be passed in. To query the security groups associated with the RO groups, only pass in ReadOnlyGroupId.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlSecurityGroupConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_security_group_config.create")()
	defer inconsistentCheck(d, meta)()

	pgInstanceId := ""
	roGroupId := ""

	if v, ok := d.GetOk("db_instance_id"); ok {
		pgInstanceId = v.(string)
	}

	if v, ok := d.GetOk("read_only_group_id"); ok {
		roGroupId = v.(string)
	}

	d.SetId(strings.Join([]string{pgInstanceId, roGroupId}, FILED_SP))

	return resourceTencentCloudPostgresqlSecurityGroupConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlSecurityGroupConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_security_group_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	readOnlyGroupId := idSplit[1]

	SecurityGroupConfigs, err := service.DescribePostgresqlSecurityGroupConfigById(ctx, dBInstanceId, readOnlyGroupId)
	if err != nil {
		return err
	}

	if len(SecurityGroupConfigs) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlSecurityGroupConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	sgIDList := make([]*string, 0)
	for _, sg := range SecurityGroupConfigs {
		if sg != nil {
			sgIDList = append(sgIDList, sg.SecurityGroupId)
		}
	}
	_ = d.Set("security_group_id_set", sgIDList)

	return nil
}

func resourceTencentCloudPostgresqlSecurityGroupConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_security_group_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgresql.NewModifyDBInstanceSecurityGroupsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	readOnlyGroupId := idSplit[1]

	if dBInstanceId != "" {
		request.DBInstanceId = &dBInstanceId
	}
	if readOnlyGroupId != "" {
		request.ReadOnlyGroupId = &readOnlyGroupId
	}

	if d.HasChange("security_group_id_set") {
		if v, ok := d.GetOk("security_group_id_set"); ok {
			request.SecurityGroupIdSet = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyDBInstanceSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgresql SecurityGroupConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlSecurityGroupConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlSecurityGroupConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_security_group_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
