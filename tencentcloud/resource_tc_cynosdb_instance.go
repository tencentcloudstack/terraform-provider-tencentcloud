/*
Provides a resource to create a cynosdb instance

Example Usage

```hcl
resource "tencentcloud_cynosdb_instance" "instance" {
  cluster_id = "cynosdbmysql-6gtlgm5l"
  cpu = 2
  memory = 4
  read_only_count = 1
  instance_grp_id = "cynosmysql-grp-xxxxxxxx"
  vpc_id = "vpc-1ptuei0b"
  subnet_id = "subnet-1tmw9t4o"
  port = 2000
  instance_name = "cynosmysql-xxxxxxxx"
  auto_voucher = 0
  db_type = "MYSQL"
  order_source = "api"
  deal_mode = 0
  param_template_id = 0
  instance_params {
		param_name = ""
		current_value = ""
		old_value = ""

  }
  security_group_ids =
}
```

Import

cynosdb instance can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_instance.instance instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbInstanceCreate,
		Read:   resourceTencentCloudCynosdbInstanceRead,
		Delete: resourceTencentCloudCynosdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cpu": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Cpu Kernel Number.",
			},

			"memory": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Memory, in GB.",
			},

			"read_only_count": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of new read-only instances, with a value range of [0,4].",
			},

			"instance_grp_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance group ID is used when adding an instance to an existing RO group. If it is not passed, the RO group will be added. The current version does not recommend transferring this value. The current version is obsolete.",
			},

			"vpc_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The VPC network ID to which it belongs.",
			},

			"subnet_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The subnet ID to which it belongs. If VpcId is set, SubnetId is required.",
			},

			"port": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The Port used when adding RO groups, with a value range of [065535].",
			},

			"instance_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance name, string length range is [0,64], value range is uppercase and lowercase letters, 0-9 digits, &amp;amp;#39;_&amp;amp;#39;, &amp;amp;#39;-&amp;amp;#39;, &amp;amp;#39;.&amp;amp;#39;.",
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Automatically select voucher 1 Yes 0 No Default is 0.",
			},

			"db_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database type, value range:&amp;amp;lt;li&amp;amp;gt;MYSQL&amp;amp;lt;/li&amp;amp;gt;.",
			},

			"order_source": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Order source, string length range is [0,64].",
			},

			"deal_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Transaction Mode 0- Place Order and Pay 1- Place Order.",
			},

			"param_template_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Parameter template ID.",
			},

			"instance_params": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Parameter list, InstanceParams is only valid when ParamTemplateId is passed in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Current value of parameter.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter old value (only useful when generating parameters) Note: This field may return null, indicating that a valid value cannot be obtained.",
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
				Description: "Security group ID, which can be specified when creating a read-only instance.",
			},
		},
	}
}

func resourceTencentCloudCynosdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewAddInstancesRequest()
		response   = cynosdb.NewAddInstancesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("cpu"); v != nil {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("memory"); v != nil {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("read_only_count"); v != nil {
		request.ReadOnlyCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_grp_id"); ok {
		request.InstanceGrpId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("port"); v != nil {
		request.Port = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_type"); ok {
		request.DbType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_source"); ok {
		request.OrderSource = helper.String(v.(string))
	}

	if v, _ := d.GetOk("deal_mode"); v != nil {
		request.DealMode = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("param_template_id"); v != nil {
		request.ParamTemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_params"); ok {
		for _, item := range v.([]interface{}) {
			modifyParamItem := cynosdb.ModifyParamItem{}
			if v, ok := dMap["param_name"]; ok {
				modifyParamItem.ParamName = helper.String(v.(string))
			}
			if v, ok := dMap["current_value"]; ok {
				modifyParamItem.CurrentValue = helper.String(v.(string))
			}
			if v, ok := dMap["old_value"]; ok {
				modifyParamItem.OldValue = helper.String(v.(string))
			}
			request.InstanceParams = append(request.InstanceParams, &modifyParamItem)
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().AddInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCynosdbInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
