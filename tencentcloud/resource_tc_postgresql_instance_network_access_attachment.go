/*
Provides a resource to create a postgresql instance_network_access_attachment

Example Usage

Vip assigned by system.
```hcl
resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  is_assign_vip = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

Vip specified by user.
```hcl
resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  vpc_id = local.my_vpc_id
  subnet_id = local.my_subnet_id
  is_assign_vip = true
  vip = "172.18.111.111"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgresql instance_network_access_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_network_access_attachment.instance_network_access_attachment instance_network_access_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlInstanceNetworkAccessAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentCreate,
		Read:   resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentRead,
		Delete: resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of postgres-6bwgamo3.",
			},

			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unified VPC ID.",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID.",
			},

			"is_assign_vip": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to manually assign the VIP. Valid values:true (manually assign), false (automatically assign).",
			},

			"vip": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target VIP.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_network_access_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgresql.NewCreateDBInstanceNetworkAccessRequest()
		dbInstanceId string
		vpcId        string
		vip          string
		port         string
		isUserAssign bool
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dbInstanceId = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_assign_vip"); ok {
		request.IsAssignVip = helper.Bool(v.(bool))
		isUserAssign = v.(bool)
	}

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
		vip = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().CreateDBInstanceNetworkAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgresql InstanceNetworkAccessAttachment failed, reason:%+v", logId, err)
		return err
	}

	id := strings.Join([]string{dbInstanceId, vpcId, vip, port}, FILED_SP)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"opened"}, 180*readRetryTimeout, time.Second, service.PostgresqlInstanceNetworkAccessAttachmentStateRefreshFunc(id, []string{}))

	var ret interface{}
	var e error
	if ret, e = conf.WaitForState(); e != nil {
		return e
	} else {
		object := ret.(*postgresql.DBInstanceNetInfo)
		// for _, info := range object {
		if object != nil {
			if isUserAssign {
				// find the port
				if *object.VpcId == vpcId && *object.Ip == vip {
					port = helper.UInt64ToStr(*object.Port)

				}
			} else {
				// find the port and vip when is_assign_vip is false
				if *object.VpcId == vpcId {
					port = helper.UInt64ToStr(*object.Port)
					vip = *object.Ip
				}
			}
		}
		// }
	}

	id = strings.Join([]string{dbInstanceId, vpcId, vip, port}, FILED_SP)
	d.SetId(id)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::postgres:%s:uin/:dbInstanceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_network_access_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s, location:%s", d.Id(), "resource.tencentcloud_postgresql_instance_network_access_attachment.read")
	}

	dbInstanceId := idSplit[0]
	vpcId := idSplit[1]
	vip := idSplit[2]
	port := idSplit[3]

	InstanceNetworkAccessAttachment, err := service.DescribePostgresqlInstanceNetworkAccessAttachmentById(ctx, dbInstanceId)
	if err != nil {
		return err
	}

	if InstanceNetworkAccessAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlInstanceNetworkAccessAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if InstanceNetworkAccessAttachment.DBInstanceId != nil {
		_ = d.Set("db_instance_id", InstanceNetworkAccessAttachment.DBInstanceId)
	}

	if InstanceNetworkAccessAttachment.VpcId != nil {
		_ = d.Set("vpc_id", InstanceNetworkAccessAttachment.VpcId)
	}

	if InstanceNetworkAccessAttachment.SubnetId != nil {
		_ = d.Set("subnet_id", InstanceNetworkAccessAttachment.SubnetId)
	}

	if vip == "" {
		// That's mean isUserAssign is false and need to set vip assigned by system
		if InstanceNetworkAccessAttachment.DBInstanceNetInfo != nil {
			for _, info := range InstanceNetworkAccessAttachment.DBInstanceNetInfo {
				if *info.VpcId == vpcId && helper.UInt64ToStr(*info.Port) == port {
					if info.Ip != nil {
						vip = *info.Ip
						log.Printf("[DEBUG]%s the id:[%s]'s filed vip[%s] updated successfully!\n", logId, d.Id(), vip)
						break
					}
				}
			}
		}
		// update the vip into unique id
		id := strings.Join([]string{dbInstanceId, vpcId, vip, port}, FILED_SP)
		d.SetId(id)
	}
	_ = d.Set("vip", vip)

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "dbInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance_network_access_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s, location:%s", d.Id(), "resource.tencentcloud_postgresql_instance_network_access_attachment.delete")
	}

	var subnetId string
	dbInstanceId := idSplit[0]
	vpcId := idSplit[1]
	vip := idSplit[2]
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}

	if err := service.DeletePostgresqlInstanceNetworkAccessAttachmentById(ctx, dbInstanceId, vpcId, subnetId, vip); err != nil {
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"closed"}, 180*readRetryTimeout, time.Second, service.PostgresqlInstanceNetworkAccessAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
