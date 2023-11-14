/*
Provides a resource to create a postgres clone_d_b_instance

Example Usage

```hcl
resource "tencentcloud_postgres_clone_d_b_instance" "clone_d_b_instance" {
  d_b_instance_id = ""
  spec_code = ""
  storage =
  period =
  auto_renew_flag =
  vpc_id = ""
  subnet_id = ""
  name = ""
  instance_charge_type = ""
  security_group_ids =
  project_id =
  tag_list {
		tag_key = ""
		tag_value = ""

  }
  d_b_node_set {
		role = ""
		zone = ""

  }
  auto_voucher =
  voucher_ids = ""
  activity_id =
  backup_set_id = ""
  recovery_target_time = ""
}
```

Import

postgres clone_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_clone_d_b_instance.clone_d_b_instance clone_d_b_instance_id
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

func resourceTencentCloudPostgresCloneDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresCloneDBInstanceCreate,
		Read:   resourceTencentCloudPostgresCloneDBInstanceRead,
		Delete: resourceTencentCloudPostgresCloneDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the original instance to be cloned.",
			},

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

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Valid period in months of the purchased instance. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. This parameter is set to `1` when the pay-as-you-go billing mode is used.",
			},

			"auto_renew_flag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal flag. Valid values: `0` (manual renewal), `1` (auto-renewal). Default value: `0`.",
			},

			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPC ID.",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of a subnet in the VPC specified by `VpcId`.",
			},

			"name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of the purchased instance.",
			},

			"instance_charge_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance billing mode. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_HOUR` (pay-as-you-go).",
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

			"project_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"tag_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The information of tags to be bound with the purchased instance. This parameter is left empty by default.",
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

			"d_b_node_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "This parameter is required if you purchase a multi-AZ deployed instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type. Valid values:`Primary`;`Standby`.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AZ where the node resides, such as ap-guangzhou-1.",
						},
					},
				},
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to automatically use vouchers. Valid values: `1` (yes), `0` (no). Default value: `0`.",
			},

			"voucher_ids": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Voucher ID list.",
			},

			"activity_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Campaign ID.",
			},

			"backup_set_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Basic backup set ID.",
			},

			"recovery_target_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Restoration point in time.",
			},
		},
	}
}

func resourceTencentCloudPostgresCloneDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_clone_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewCloneDBInstanceRequest()
		response     = postgres.NewCloneDBInstanceResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
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

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tag_list"); ok {
		for _, item := range v.([]interface{}) {
			tag := postgres.Tag{}
			if v, ok := dMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			request.TagList = append(request.TagList, &tag)
		}
	}

	if v, ok := d.GetOk("d_b_node_set"); ok {
		for _, item := range v.([]interface{}) {
			dBNode := postgres.DBNode{}
			if v, ok := dMap["role"]; ok {
				dBNode.Role = helper.String(v.(string))
			}
			if v, ok := dMap["zone"]; ok {
				dBNode.Zone = helper.String(v.(string))
			}
			request.DBNodeSet = append(request.DBNodeSet, &dBNode)
		}
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		request.VoucherIds = helper.String(v.(string))
	}

	if v, _ := d.GetOk("activity_id"); v != nil {
		request.ActivityId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("backup_set_id"); ok {
		request.BackupSetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("recovery_target_time"); ok {
		request.RecoveryTargetTime = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().CloneDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres CloneDBInstance failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 600*readRetryTimeout, time.Second, service.PostgresCloneDBInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresCloneDBInstanceRead(d, meta)
}

func resourceTencentCloudPostgresCloneDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_clone_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresCloneDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_clone_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
