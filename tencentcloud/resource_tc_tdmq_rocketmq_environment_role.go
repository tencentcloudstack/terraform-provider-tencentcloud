/*
Provides a resource to create a tdmqRocketmq environment_role

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_environment_role" "environment_role" {
  environment_id = &lt;nil&gt;
  role_name = &lt;nil&gt;
  permissions = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
}
```

Import

tdmqRocketmq environment_role can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rocketmq_environment_role.environment_role environment_role_id
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

func resourceTencentCloudTdmqRocketmqEnvironmentRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqEnvironmentRoleCreate,
		Read:   resourceTencentCloudTdmqRocketmqEnvironmentRoleRead,
		Update: resourceTencentCloudTdmqRocketmqEnvironmentRoleUpdate,
		Delete: resourceTencentCloudTdmqRocketmqEnvironmentRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment (namespace) name.",
			},

			"role_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Role Name.",
			},

			"permissions": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Permissions, which is a non-empty string array of `produce` and `consume` at the most.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID (required).",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_environment_role.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tdmqRocketmq.NewCreateEnvironmentRoleRequest()
		response      = tdmqRocketmq.NewCreateEnvironmentRoleResponse()
		clusterId     string
		roleName      string
		environmentId string
	)
	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_name"); ok {
		roleName = v.(string)
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("permissions"); ok {
		permissionsSet := v.(*schema.Set).List()
		for i := range permissionsSet {
			permissions := permissionsSet[i].(string)
			request.Permissions = append(request.Permissions, &permissions)
		}
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().CreateEnvironmentRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq environmentRole failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{clusterId, roleName, environmentId}, FILED_SP))

	return resourceTencentCloudTdmqRocketmqEnvironmentRoleRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_environment_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]
	environmentId := idSplit[2]

	environmentRole, err := service.DescribeTdmqRocketmqEnvironmentRoleById(ctx, clusterId, roleName, environmentId)
	if err != nil {
		return err
	}

	if environmentRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqEnvironmentRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if environmentRole.EnvironmentId != nil {
		_ = d.Set("environment_id", environmentRole.EnvironmentId)
	}

	if environmentRole.RoleName != nil {
		_ = d.Set("role_name", environmentRole.RoleName)
	}

	if environmentRole.Permissions != nil {
		_ = d.Set("permissions", environmentRole.Permissions)
	}

	if environmentRole.ClusterId != nil {
		_ = d.Set("cluster_id", environmentRole.ClusterId)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_environment_role.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmqRocketmq.NewModifyEnvironmentRoleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]
	environmentId := idSplit[2]

	request.ClusterId = &clusterId
	request.RoleName = &roleName
	request.EnvironmentId = &environmentId

	immutableArgs := []string{"environment_id", "role_name", "permissions", "cluster_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("permissions") {
		if v, ok := d.GetOk("permissions"); ok {
			permissionsSet := v.(*schema.Set).List()
			for i := range permissionsSet {
				permissions := permissionsSet[i].(string)
				request.Permissions = append(request.Permissions, &permissions)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqRocketmqClient().ModifyEnvironmentRole(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmqRocketmq environmentRole failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqEnvironmentRoleRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rocketmq_environment_role.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]
	environmentId := idSplit[2]

	if err := service.DeleteTdmqRocketmqEnvironmentRoleById(ctx, clusterId, roleName, environmentId); err != nil {
		return err
	}

	return nil
}
