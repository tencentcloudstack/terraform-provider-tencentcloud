/*
Provides a resource to create a postgres rebalance_read_only_group

Example Usage

```hcl
resource "tencentcloud_postgres_rebalance_read_only_group" "rebalance_read_only_group" {
  read_only_group_id = "pgrogrp-test"
}
```

Import

postgres rebalance_read_only_group can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_rebalance_read_only_group.rebalance_read_only_group rebalance_read_only_group_id
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

func resourceTencentCloudPostgresRebalanceReadOnlyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresRebalanceReadOnlyGroupCreate,
		Read:   resourceTencentCloudPostgresRebalanceReadOnlyGroupRead,
		Delete: resourceTencentCloudPostgresRebalanceReadOnlyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"read_only_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Readonly Group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresRebalanceReadOnlyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_rebalance_read_only_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgres.NewRebalanceReadOnlyGroupRequest()
		response        = postgres.NewRebalanceReadOnlyGroupResponse()
		readOnlyGroupId string
	)
	if v, ok := d.GetOk("read_only_group_id"); ok {
		readOnlyGroupId = v.(string)
		request.ReadOnlyGroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().RebalanceReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres RebalanceReadOnlyGroup failed, reason:%+v", logId, err)
		return err
	}

	readOnlyGroupId = *response.Response.ReadOnlyGroupId
	d.SetId(readOnlyGroupId)

	return resourceTencentCloudPostgresRebalanceReadOnlyGroupRead(d, meta)
}

func resourceTencentCloudPostgresRebalanceReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_rebalance_read_only_group.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresRebalanceReadOnlyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_rebalance_read_only_group.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
