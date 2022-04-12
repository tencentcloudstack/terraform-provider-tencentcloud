/*
Use this resource to create postgresql readonly attachment.

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_attachment" "attach" {
  db_instance_id = tencentcloud_postgresql_readonly_instance.foo.id
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlReadonlyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadOnlyAttachmentCreate,
		Read:   resourceTencentCloudPostgresqlReadOnlyAttachmentRead,
		//Update: resourceTencentCloudPostgresqlReadOnlyAttachmentUpdate,
		Delete: resourceTencentCLoudPostgresqlReadOnlyAttachmentDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Read only instance ID.",
			},
			"read_only_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Read only group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadOnlyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_attachment.create")()

	logId := getLogId(contextNil)

	var (
		request      = postgresql.NewAddDBInstanceToReadOnlyGroupRequest()
		dbInstanceId string
		groupId      string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		dbInstanceId = v.(string)
		request.DBInstanceId = helper.String(dbInstanceId)

	}
	if v, ok := d.GetOk("read_only_group_id"); ok {
		groupId = v.(string)
		request.ReadOnlyGroupId = helper.String(groupId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().AddDBInstanceToReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}
	instanceId := helper.IdFormat(dbInstanceId, groupId)
	d.SetId(instanceId)

	return nil
}

func resourceTencentCloudPostgresqlReadOnlyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, err := postgresqlService.DescribePostgresqlReadOnlyGroupById(ctx, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCLoudPostgresqlReadOnlyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_attachment.delete")()

	logId := getLogId(contextNil)
	request := postgresql.NewRemoveDBInstanceFromReadOnlyGroupRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dbInstanceId := idSplit[0]
	groupId := idSplit[1]
	request.ReadOnlyGroupId = helper.String(groupId)
	request.DBInstanceId = helper.String(dbInstanceId)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().RemoveDBInstanceFromReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
