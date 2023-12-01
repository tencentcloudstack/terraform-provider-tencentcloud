/*
Provide a resource to create a TDMQ role.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}
```

Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.test tdmq_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudTdmqRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRoleCreate,
		Read:   resourceTencentCloudTdmqRoleRead,
		Update: resourceTencentCloudTdmqRoleUpdate,
		Delete: resourceTencentCloudTdmqRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq role.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of tdmq cluster.",
			},
			"remark": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of tdmq role.",
			},
		},
	}
}

func resourceTencentCloudTdmqRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_role.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		roleName  string
		clusterId string
		remark    string
	)
	if temp, ok := d.GetOk("role_name"); ok {
		roleName = temp.(string)
		if len(roleName) < 1 {
			return fmt.Errorf("role_name should be not empty string")
		}
	}

	if temp, ok := d.GetOk("cluster_id"); ok {
		clusterId = temp.(string)
		if len(clusterId) < 1 {
			return fmt.Errorf("cluster_id should be not empty string")
		}
	}

	if temp, ok := d.GetOk("remark"); ok {
		remark = temp.(string)
	}

	clusterId, err := tdmqService.CreateTdmqRole(ctx, roleName, clusterId, remark)
	if err != nil {
		return err
	}
	d.SetId(clusterId)

	return resourceTencentCloudTdmqRoleRead(d, meta)
}

func resourceTencentCloudTdmqRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_role.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	roleName := d.Id()
	clusterId := d.Get("cluster_id").(string)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqRoleById(ctx, roleName, clusterId)
		if e != nil {
			return retryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("role_name", info.RoleName)
		_ = d.Set("remark", info.Remark)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTdmqRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_role.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	roleName := d.Id()
	clusterId := d.Get("cluster_id").(string)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		remark string
	)
	old, now := d.GetChange("remark")
	if d.HasChange("remark") {
		remark = now.(string)
	} else {
		remark = old.(string)
	}

	d.Partial(true)

	if err := service.ModifyTdmqRoleAttribute(ctx, roleName, clusterId, remark); err != nil {
		return err
	}

	d.Partial(false)
	return resourceTencentCloudTdmqRoleRead(d, meta)
}

func resourceTencentCloudTdmqRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_role.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	roleName := d.Id()
	clusterId := d.Get("cluster_id").(string)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqRole(ctx, roleName, clusterId); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
