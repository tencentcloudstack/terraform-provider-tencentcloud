/*
Provide a resource to create a TDMQ role.

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark = "this is description."
}

resource "tencentcloud_tdmq_namespace" "bar" {
  environ_name = "example"
  msg_ttl = 300
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark = "this is description."
}

resource "tencentcloud_tdmq_topic" "bar" {
  environ_id = "${tencentcloud_tdmq_namespace.bar.id}"
  topic_name = "example"
  partitions = 6
  topic_type = 0
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark = "this is description."
}

resource "tencentcloud_tdmq_role" "bar" {
  role_name = "example"
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark = "this is description world"
}

resource "tencentcloud_tdmq_namespace_role_attachment" "bar" {
  environ_id = "${tencentcloud_tdmq_namespace.bar.id}"
  role_name = "${tencentcloud_tdmq_role.bar.role_name}"
  permissions = ["produce", "consume"]
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqNamespaceRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqNamespaceRoleAttachmentCreate,
		Read:   resourceTencentCloudTdmqNamespaceRoleAttachmentRead,
		Update: resourceTencentCloudTdmqNamespaceRoleAttachmentUpdate,
		Delete: resourceTencentCloudTdmqNamespaceRoleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environ_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq namespace.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq role.",
			},
			"permissions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The permissions of tdmq role.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of tdmq cluster.",
			},
			//compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
		},
	}
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		environId   string
		roleName    string
		permissions []*string
		clusterId   string
	)
	if temp, ok := d.GetOk("environ_id"); ok {
		environId = temp.(string)
		if len(environId) < 1 {
			return fmt.Errorf("environ_id should be not empty string")
		}
	}

	if temp, ok := d.GetOk("role_name"); ok {
		roleName = temp.(string)
		if len(roleName) < 1 {
			return fmt.Errorf("role_name should be not empty string")
		}
	}

	if v, ok := d.GetOk("permissions"); ok {
		for _, id := range v.([]interface{}) {
			permissions = append(permissions, helper.String(id.(string)))
		}
	}

	if temp, ok := d.GetOk("cluster_id"); ok {
		clusterId = temp.(string)
		if len(clusterId) < 1 {
			return fmt.Errorf("cluster_id should be not empty string")
		}
	}

	err := tdmqService.CreateTdmqNamespaceRoleAttachment(ctx, environId, roleName, permissions, clusterId)
	if err != nil {
		return err
	}

	d.SetId(environId + FILED_SP + roleName)

	return resourceTencentCloudTdmqNamespaceRoleAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("environment role id is borken, id is %s", d.Id())
	}
	environId := idSplit[0]
	roleName := idSplit[1]
	clusterId := d.Get("cluster_id").(string)

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqNamespaceRoleAttachment(ctx, environId, roleName, clusterId)
		if e != nil {
			return retryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}
		_ = d.Set("environ_id", info.EnvironmentId)
		_ = d.Set("role_name", info.RoleName)
		_ = d.Set("permissions", info.Permissions)
		_ = d.Set("create_time", info.CreateTime)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("environment role id is borken, id is %s", d.Id())
	}
	environId := idSplit[0]
	roleName := idSplit[1]
	clusterId := d.Get("cluster_id").(string)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		permissions []*string
	)
	old, now := d.GetChange("permissions")
	if d.HasChange("permissions") {
		for _, id := range now.([]interface{}) {
			permissions = append(permissions, helper.String(id.(string)))
		}
	} else {
		for _, id := range old.([]interface{}) {
			permissions = append(permissions, helper.String(id.(string)))
		}
	}

	d.Partial(true)

	if err := service.ModifyTdmqNamespaceRoleAttachment(ctx, environId, roleName, permissions, clusterId); err != nil {
		return err
	}
	d.SetPartial("permissions")

	d.Partial(false)
	return resourceTencentCloudTdmqNamespaceRoleAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("environment role id is borken, id is %s", d.Id())
	}
	environId := idSplit[0]
	roleName := idSplit[1]
	clusterId := d.Get("cluster_id").(string)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqNamespaceRoleAttachment(ctx, environId, roleName, clusterId); err != nil {
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
