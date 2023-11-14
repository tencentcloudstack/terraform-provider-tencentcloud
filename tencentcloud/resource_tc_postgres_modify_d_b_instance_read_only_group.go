/*
Provides a resource to create a postgres modify_d_b_instance_read_only_group

Example Usage

```hcl
resource "tencentcloud_postgres_modify_d_b_instance_read_only_group" "modify_d_b_instance_read_only_group" {
  d_b_instance_id = "postgres-6r233v55"
  read_only_group_id = "pgrogrp-test1"
  new_read_only_group_id = "pgrogrp-test2"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres modify_d_b_instance_read_only_group can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_modify_d_b_instance_read_only_group.modify_d_b_instance_read_only_group modify_d_b_instance_read_only_group_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupCreate,
		Read:   resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupRead,
		Delete: resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DbInstance ID.",
			},

			"read_only_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the RO group to which the read-only replica belongs.",
			},

			"new_read_only_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the new RO group into which the read-only replica will move.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_d_b_instance_read_only_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgres.NewModifyDBInstanceReadOnlyGroupRequest()
		response        = postgres.NewModifyDBInstanceReadOnlyGroupResponse()
		dBInstanceId    string
		readOnlyGroupId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_only_group_id"); ok {
		readOnlyGroupId = v.(string)
		request.ReadOnlyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("new_read_only_group_id"); ok {
		request.NewReadOnlyGroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyDBInstanceReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres ModifyDBInstanceReadOnlyGroup failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(strings.Join([]string{dBInstanceId, readOnlyGroupId}, FILED_SP))

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"ok"}, 120*readRetryTimeout, time.Second, service.PostgresModifyDBInstanceReadOnlyGroupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupRead(d, meta)
}

func resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_d_b_instance_read_only_group.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresModifyDBInstanceReadOnlyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_d_b_instance_read_only_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
