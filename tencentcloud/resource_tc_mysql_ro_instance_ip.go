/*
Provides a resource to create a mysql ro_instance_ip

Example Usage

```hcl
resource "tencentcloud_mysql_ro_instance_ip" "ro_instance_ip" {
  instance_id = "cdbro-bdlvcfpj"
  uniq_subnet_id = "subnet-dwj7ipnc"
  uniq_vpc_id = "vpc-4owdpnwr"
}

```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlRoInstanceIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRoInstanceIpCreate,
		Read:   resourceTencentCloudMysqlRoInstanceIpRead,
		Delete: resourceTencentCloudMysqlRoInstanceIpDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Read-only instance ID, in the format: cdbro-3i70uj0k, which is the same as the read-only instance ID displayed on the cloud database console page.",
			},

			"uniq_subnet_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subnet descriptor, for example: subnet-1typ0s7d.",
			},

			"uniq_vpc_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "vpc descriptor, for example: vpc-a23yt67j, if this field is passed, UniqSubnetId must be passed.",
			},

			"ro_vip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet IP address of the read-only instance.",
			},

			"ro_vport": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Intranet port number of the read-only instance.",
			},
		},
	}
}

func resourceTencentCloudMysqlRoInstanceIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_instance_ip.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mysql.NewCreateRoInstanceIpRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_subnet_id"); ok {
		request.UniqSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateRoInstanceIp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql roInstanceIp failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlRoInstanceIpRead(d, meta)
}

func resourceTencentCloudMysqlRoInstanceIpRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_instance_ip.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	switchForUpgrade, err := service.DescribeDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if switchForUpgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlSwitchForUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if switchForUpgrade.InstanceId != nil {
		_ = d.Set("instance_id", switchForUpgrade.InstanceId)
	}

	if switchForUpgrade.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", switchForUpgrade.UniqVpcId)
	}

	if switchForUpgrade.UniqSubnetId != nil {
		_ = d.Set("uniq_subnet_id", switchForUpgrade.UniqSubnetId)
	}

	if switchForUpgrade.RoVipInfo != nil {
		if switchForUpgrade.RoVipInfo.RoVip != nil {
			_ = d.Set("ro_vip", switchForUpgrade.RoVipInfo.RoVip)
		}

		if switchForUpgrade.RoVipInfo.RoVport != nil {
			_ = d.Set("ro_vport", switchForUpgrade.RoVipInfo.RoVport)
		}
	}

	return nil
}

func resourceTencentCloudMysqlRoInstanceIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_ro_instance_ip.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
