/*
Provides a resource to create a tdmq rabbitmq_role

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_role" "rabbitmq_role" {
  role_name = ""
  cluster_id = ""
  remark = ""
}

```
Import

tdmq rabbitmq_role can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rabbitmq_role.rabbitmq_role rabbitmqRole_id
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
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqRabbitmqRole() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRabbitmqRoleRead,
		Create: resourceTencentCloudTdmqRabbitmqRoleCreate,
		Update: resourceTencentCloudTdmqRabbitmqRoleUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "role name.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster id.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "role description, 128 characters or less.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmq.NewCreateRoleRequest()
		clusterId string
		roleName  string
	)

	if v, ok := d.GetOk("role_name"); ok {
		roleName = v.(string)
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {

		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqRole failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + roleName)
	return resourceTencentCloudTdmqRabbitmqRoleRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]

	rabbitmqRole, err := service.DescribeTdmqRabbitmqRole(ctx, clusterId, roleName)

	if err != nil {
		return err
	}

	if rabbitmqRole == nil {
		d.SetId("")
		return fmt.Errorf("resource `rabbitmqRole` %s does not exist", clusterId)
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("role_name", roleName)

	if rabbitmqRole.Remark != nil {
		_ = d.Set("remark", rabbitmqRole.Remark)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyRoleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]

	request.ClusterId = &clusterId
	request.RoleName = &roleName

	if d.HasChange("role_name") {
		return fmt.Errorf("`role_name` do not support change now.")
	}

	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqRole failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqRoleRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]

	if err := service.DeleteTdmqRabbitmqRoleById(ctx, clusterId, roleName); err != nil {
		return err
	}

	return nil
}
