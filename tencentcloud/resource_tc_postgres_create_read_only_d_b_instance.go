/*
Provides a resource to create a postgres create_read_only_d_b_instance

Example Usage

```hcl
resource "tencentcloud_postgres_create_read_only_d_b_instance" "create_read_only_d_b_instance" {
  spec_code = ""
  storage =
  instance_count =
  period =
  master_d_b_instance_id = ""
  zone = ""
  project_id =
  d_b_version = ""
  instance_charge_type = ""
  auto_voucher =
  voucher_ids =
  auto_renew_flag =
  vpc_id = ""
  subnet_id = ""
  activity_id =
  name = ""
  need_support_ipv6 =
  read_only_group_id = ""
  tag_list {
		tag_key = ""
		tag_value = ""

  }
  security_group_ids =
}
```

Import

postgres create_read_only_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_create_read_only_d_b_instance.create_read_only_d_b_instance create_read_only_d_b_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudPostgresCreateReadOnlyDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresCreateReadOnlyDBInstanceCreate,
		Read:   resourceTencentCloudPostgresCreateReadOnlyDBInstanceRead,
		Delete: resourceTencentCloudPostgresCreateReadOnlyDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"spec_code": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Purchasable specification ID, which can be obtained through the `SpecCode` field in the returned value of the `DescribeProductConfig` API.",
			},

			"storage": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Instance storage capacity in GB.",
			},

			"instance_count": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of instances purchased at a time. Value range: 1–100.",
			},

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Valid period in months of purchased instances. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. This parameter is set to `1` when the pay-as-you-go billing mode is used.",
			},

			"master_d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the primary instance to which the read-only replica belongs.",
			},

			"zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Availability zone ID, which can be obtained through the `Zone` field in the returned value of the `DescribeZones` API.",
			},

			"project_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"d_b_version": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "(Disused) You don’t need to specify a version, as the kernel version is as the same as that of the instance.",
			},

			"instance_charge_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance billing mode. Valid value: `POSTPAID_BY_HOUR` (pay-as-you-go). If the source instance is pay-as-you-go, so is the read-only instance.",
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to automatically use vouchers. Valid values: `1` (yes), `0` (no). Default value: `0`.",
			},

			"voucher_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID list. Currently, you can specify only one voucher.",
			},

			"auto_renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal flag. Valid values: `0` (manual renewal), `1` (auto-renewal). Default value: `0`.",
			},

			"vpc_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPC ID.",
			},

			"subnet_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPC subnet ID.",
			},

			"activity_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Special offer ID.",
			},

			"name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance name (which will be supported in the future).",
			},

			"need_support_ipv6": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to support IPv6 address access. Valid values: `1` (yes), `0` (no).",
			},

			"read_only_group_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "RO group ID.",
			},

			"tag_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The information of tags to be bound with the purchased instance, which is left empty by default (type: tag array).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresCreateReadOnlyDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_create_read_only_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewCreateReadOnlyDBInstanceRequest()
		response     = postgres.NewCreateReadOnlyDBInstanceResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		request.Storage = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("instance_count"); v != nil {
		request.InstanceCount = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("master_d_b_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("d_b_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, _ := d.GetOk("auto_renew_flag"); v != nil {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("activity_id"); v != nil {
		request.ActivityId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, _ := d.GetOk("need_support_ipv6"); v != nil {
		request.NeedSupportIpv6 = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "tag_list"); ok {
		tag := postgres.Tag{}
		if v, ok := dMap["tag_key"]; ok {
			tag.TagKey = helper.String(v.(string))
		}
		if v, ok := dMap["tag_value"]; ok {
			tag.TagValue = helper.String(v.(string))
		}
		request.TagList = &tag
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().CreateReadOnlyDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres CreateReadOnlyDBInstance failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 600*readRetryTimeout, time.Second, service.PostgresCreateReadOnlyDBInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresCreateReadOnlyDBInstanceRead(d, meta)
}

func resourceTencentCloudPostgresCreateReadOnlyDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_create_read_only_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresCreateReadOnlyDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_create_read_only_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
