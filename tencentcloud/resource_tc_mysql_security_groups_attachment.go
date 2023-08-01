/*
Provides a resource to create a mysql security_groups_attachment

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_security_groups_attachment" "example" {
  security_group_id = tencentcloud_security_group.security_group.id
  instance_id       = tencentcloud_mysql_instance.example.id
}
```

Import

mysql security_groups_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_security_groups_attachment.security_groups_attachment securityGroupId#instanceId
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
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlSecurityGroupsAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSecurityGroupsAttachmentCreate,
		Read:   resourceTencentCloudMysqlSecurityGroupsAttachmentRead,
		Delete: resourceTencentCloudMysqlSecurityGroupsAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of security group.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The id of instance.",
			},
		},
	}
}

func resourceTencentCloudMysqlSecurityGroupsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_security_groups_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = mysql.NewAssociateSecurityGroupsRequest()
		securityGroupId string
		instanceId      string
	)
	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql securityGroupsAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(securityGroupId + FILED_SP + instanceId)

	return resourceTencentCloudMysqlSecurityGroupsAttachmentRead(d, meta)
}

func resourceTencentCloudMysqlSecurityGroupsAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_security_groups_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	securityGroupsAttachment, err := service.DescribeMysqlSecurityGroupsAttachmentById(ctx, securityGroupId, instanceId)
	if err != nil {
		return err
	}

	if securityGroupsAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlSecurityGroupsAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil

	}
	_ = d.Set("instance_id", instanceId)
	if securityGroupsAttachment.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroupsAttachment.SecurityGroupId)
	}

	return nil
}

func resourceTencentCloudMysqlSecurityGroupsAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_security_groups_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteMysqlSecurityGroupsAttachmentById(ctx, securityGroupId, instanceId); err != nil {
		return err
	}

	return nil
}
