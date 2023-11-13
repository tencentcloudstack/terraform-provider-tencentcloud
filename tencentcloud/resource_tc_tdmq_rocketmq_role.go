/*
Provides a resource to create a tdmqRocketmq role

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = &lt;nil&gt;
  remark = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
      }
```

Import

tdmqRocketmq role can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_role.role role_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTdmqRocketmqRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqRoleCreate,
		Read:   resourceTencentCloudTdmqRocketmqRoleRead,
		Update: resourceTencentCloudTdmqRocketmqRoleUpdate,
		Delete: resourceTencentCloudTdmqRocketmqRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Role name, which can contain up to 32 letters, digits, hyphens, and underscores.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remarks (up to 128 characters).",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID (required).",
			},

			"token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Value of the role token.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tdmqRocketmq.NewCreateRoleRequest()
		response  = tdmqRocketmq.NewCreateRoleResponse()
		clusterId string
		roleName  string
	)
	if v, ok := d.GetOk("role_name"); ok {
		roleName = v.(string)
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().CreateRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq role failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{clusterId, roleName}, FILED_SP))

	return resourceTencentCloudTdmqRocketmqRoleRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]

	role, err := service.DescribeTdmqRocketmqRoleById(ctx, clusterId, roleName)
	if err != nil {
		return err
	}

	if role == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if role.RoleName != nil {
		_ = d.Set("role_name", role.RoleName)
	}

	if role.Remark != nil {
		_ = d.Set("remark", role.Remark)
	}

	if role.ClusterId != nil {
		_ = d.Set("cluster_id", role.ClusterId)
	}

	if role.Token != nil {
		_ = d.Set("token", role.Token)
	}

	if role.CreateTime != nil {
		_ = d.Set("create_time", role.CreateTime)
	}

	if role.UpdateTime != nil {
		_ = d.Set("update_time", role.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyRoleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]

	request.ClusterId = &clusterId
	request.RoleName = &roleName

	immutableArgs := []string{"role_name", "remark", "cluster_id", "token", "create_time", "update_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().ModifyRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmqRocketmq role failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqRoleRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]

	if err := service.DeleteTdmqRocketmqRoleById(ctx, clusterId, roleName); err != nil {
		return err
	}

	return nil
}
