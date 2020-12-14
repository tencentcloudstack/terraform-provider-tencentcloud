/*
Provides a resource to create a Ckafka Acl.

Example Usage

Ckafka Acl

```hcl
resource "tencentcloud_ckafka_acl" "foo" {
  instance_id     = "ckafka-f9ife4zz"
  resource_type   = "TOPIC"
  resource_name   = "topic-tf-test"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "*"
  principal       = tencentcloud_ckafka_user.foo.account_name
}
```

Import

Ckafka acl can be imported using the instance_id#permission_type#principal#host#operation_type#resource_type#resource_name, e.g.

```
$ terraform import tencentcloud_ckafka_acl.foo ckafka-f9ife4zz#ALLOW#test#*#WRITE#TOPIC#topic-tf-test
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudCkafkaAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaAclCreate,
		Read:   resourceTencentCloudCkafkaAclRead,
		Delete: resourceTencentCloudCkafkaAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ckafka instance.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TOPIC",
				ForceNew:    true,
				Description: "ACL resource type. Valid values are `UNKNOWN`, `ANY`, `TOPIC`, `GROUP`, `CLUSTER`, `TRANSACTIONAL_ID`. and `TOPIC` by default. Currently, only `TOPIC` is available, and other fields will be used for future ACLs compatible with open-source Kafka.",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL resource name, which is related to `resource_type`. For example, if `resource_type` is `TOPIC`, this field indicates the topic name; if `resource_type` is `GROUP`, this field indicates the group name.",
			},
			"operation_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL operation mode. Valid values: `UNKNOWN`, `ANY`, `ALL`, `READ`, `WRITE`, `CREATE`, `DELETE`, `ALTER`, `DESCRIBE`, `CLUSTER_ACTION`, `DESCRIBE_CONFIGS` and `ALTER_CONFIGS`.",
			},
			"permission_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALLOW",
				ForceNew:    true,
				Description: "ACL permission type. Valid values: `UNKNOWN`, `ANY`, `DENY`, `ALLOW`. and `ALLOW` by default. Currently, CKafka supports `ALLOW` (equivalent to allow list), and other fields will be used for future ACLs compatible with open-source Kafka.",
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*",
				ForceNew: true,
				ValidateFunc: func(i interface{}, s string) (rs []string, error []error) {
					if i.(string) == "*" {
						return
					}
					rs, error = validateIp(i, s)
					return
				},
				Description: "IP address allowed to access. The default value is `*`, which means that any host can access.",
			},
			"principal": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				ForceNew:    true,
				Description: "User list. The default value is `*`, which means that any user can access. The current user can only be one included in the user list.",
			},
		},
	}
}

func resourceTencentCloudCkafkaAclCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_acl.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	resourceType := d.Get("resource_type").(string)
	resourceName := d.Get("resource_name").(string)
	operation := d.Get("operation_type").(string)
	permissionType := d.Get("permission_type").(string)
	host := d.Get("host").(string)
	principal := d.Get("principal").(string)

	ckafkaService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	if err := ckafkaService.CreateAcl(ctx, instanceId, resourceType, resourceName, operation, permissionType, host, principal); err != nil {
		return fmt.Errorf("[CRITAL]%s create ckafka user failed, reason:%+v", logId, err)
	}
	d.SetId(instanceId + FILED_SP + permissionType + FILED_SP + principal + FILED_SP + host + FILED_SP + operation + FILED_SP + resourceType + FILED_SP + resourceName)

	return resourceTencentCloudCkafkaAclRead(d, meta)
}

func resourceTencentCloudCkafkaAclRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_acl.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	id := d.Id()
	info, has, err := ckafkaService.DescribeAclByAclId(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	items := strings.Split(id, FILED_SP)
	_ = d.Set("instance_id", items[0])
	_ = d.Set("resource_type", CKAFKA_ACL_RESOURCE_TYPE_TO_STRING[*info.ResourceType])
	_ = d.Set("resource_name", info.ResourceName)
	_ = d.Set("operation_type", CKAFKA_ACL_OPERATION_TO_STRING[*info.Operation])
	_ = d.Set("permission_type", CKAFKA_PERMISSION_TYPE_TO_STRING[*info.PermissionType])
	_ = d.Set("host", info.Host)
	_ = d.Set("principal", strings.TrimLeft(*info.Principal, CKAFKA_ACL_PRINCIPAL_STR))

	return nil
}

func resourceTencentCloudCkafkaAclDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_user.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if err := ckafkaService.DeleteAcl(ctx, d.Id()); err != nil {
		return err
	}
	return nil
}
