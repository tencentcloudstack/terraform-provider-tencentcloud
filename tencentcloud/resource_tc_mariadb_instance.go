/*
Provides a resource to create a mariadb instance

Example Usage

```hcl
resource "tencentcloud_mariadb_instance" "instance" {
  zones = ["ap-guangzhou-3",]
  node_count = 2
  memory = 8
  storage = 10
  period = 1
  # auto_voucher =
  # voucher_ids =
  vpc_id = "vpc-ii1jfbhl"
  subnet_id = "subnet-3ku415by"
  # project_id = ""
  db_version_id = "8.0"
  instance_name = "terraform-test"
  # security_group_ids = ""
  auto_renew_flag = 1
  ipv6_flag = 0
  tags = {
    "createby" = "terrafrom-2"
  }
  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params {
    param = "lower_case_table_names"
    value = "0"
  }
  init_params {
    param = "innodb_page_size"
    value = "16384"
  }
  init_params {
    param = "sync_mode"
    value = "1"
  }
  dcn_region = ""
  dcn_instance_id = ""
}
```
Import

mariadb tencentcloud_mariadb_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_instance.instance tdsql-4pzs5b67
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbInstanceCreate,
		Read:   resourceTencentCloudMariadbInstanceRead,
		Update: resourceTencentCloudMariadbInstanceUpdate,
		Delete: resourceTencentCloudMariadbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID, uniquely identifies a TDSQL instance.",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance name, you can set the name of the instance independently through this field.",
			},

			"zones": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance node availability zone distribution, up to two availability zones can be filled. When the shard specification is one master and two slaves, two of the nodes are in the first availability zone.",
			},

			"node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of nodes, 2 is one master and one slave, 3 is one master and two slaves.",
			},

			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory size, unit: GB, can be obtained by querying instance specifications through DescribeDBInstanceSpecs.",
			},

			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Storage size, unit: GB. You can query instance specifications through DescribeDBInstanceSpecs to obtain the lower and upper limits of disk specifications corresponding to different memory sizes.",
			},

			"period": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The duration of the purchase, unit: month.",
			},

			// "count": {
			// 	Optional:    true,
			// 	Type:        schema.TypeInt,
			// 	Description: "The quantity to be purchased, the price of purchasing 1 instance is queried by default.",
			// },

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically use the voucher for payment, the default is not used.",
			},

			"voucher_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of voucher IDs. Currently, only one voucher can be specified.",
			},

			"vpc_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Virtual private network ID, if not passed, it means that it is created as a basic network.",
			},

			"subnet_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Virtual private network subnet ID, required when VpcId is not empty.",
			},

			"project_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Project ID, which can be obtained by viewing the project list, if not passed, it will be associated with the default project.",
			},

			"db_version_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Database engine version, currently available: 8.0.18, 10.1.9, 5.7.17. If not passed, the default is Percona 5.7.17.",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group ID list.",
			},

			"auto_renew_flag": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag, 1: automatic renewal, 2: no automatic renewal.",
			},

			"ipv6_flag": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Whether IPv6 is supported.",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the application to which the instance belongs.",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the region where the instance is located, such as ap-shanghai.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance status: 0 creating, 1 process processing, 2 running, 3 instance not initialized, -1 instance isolated, 4 instance initializing, 5 instance deleting, 6 instance restarting, 7 data migration.",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet IP address.",
			},
			"vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Intranet port.",
			},
			"wan_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name accessed from the external network, which can be resolved by the public network.",
			},
			"wan_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Extranet IP address, accessible from the public network.",
			},
			"wan_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internet port.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance creation time, the format is 2006-01-02 15:04:05.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the instance in the format of 2006-01-02 15:04:05.",
			},
			"period_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance expiration time, the format is 2006-01-02 15:04:05.",
			},
			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account to which the instance belongs.",
			},
			"tdsql_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TDSQL version information.",
			},
			// "unique_vpc_id": {
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "String private network ID.",
			// },
			// "unique_subnet_id": {
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "String private network subnet ID.",
			// },
			"is_tmp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether it is a temporary instance, 0 means no, non-zero means yes.",
			},
			"excluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exclusive cluster ID, if it is empty, it means a normal instance.",
			},
			"pid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Product Type ID.",
			},
			"qps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum Qps value.",
			},
			"paymode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Payment Mode.",
			},
			"locker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Asynchronous task process ID when the instance is in an asynchronous task.",
			},
			"status_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the current running state of the instance.",
			},
			"wan_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "External network status, 0-unopened; 1-opened; 2-closed; 3-opening.",
			},
			"is_audit_supported": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the instance supports auditing. 1-supported; 0-not supported.",
			},
			"machine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Machine Model.",
			},
			"is_encrypt_supported": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether data encryption is supported. 1-supported; 0-not supported.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of CPU cores of the instance.",
			},
			"vipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet IPv6.",
			},
			"wan_vipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internet IPv6.",
			},
			"wan_port_ipv6": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internet IPv6 port.",
			},
			"wan_status_ipv6": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internet IPv6 status.",
			},
			"db_engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database Engine.",
			},
			"dcn_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DCN flag, 0-none, 1-primary instance, 2-disaster backup instance.",
			},
			"dcn_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DCN status, 0-none, 1-creating, 2-synchronizing, 3-disconnected.",
			},

			"dcn_dst_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of DCN disaster recovery instances.",
			},
			"instance_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "1: primary instance (exclusive), 2: primary instance, 3: disaster recovery instance, 4: disaster recovery instance (exclusive type).",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "tag list.",
			},

			"init_params": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Parameter list. The optional values of this interface are: character_set_server (character set, required) enum: utf8,latin1,gbk,utf8mb4,gb18030, lower_case_table_names (table name case sensitive, required, 0 - sensitive; 1 - insensitive), innodb_page_size (innodb data page, Default 16K), sync_mode (sync mode: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous can degenerate. The default is strong synchronous can degenerate).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "parameter name.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "parameter value.",
						},
					},
				},
			},

			"dcn_region": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DCN source region.",
			},

			"dcn_instance_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DCN source instance ID.",
			},
		},
	}
}

func resourceTencentCloudMariadbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mariadb.NewCreateDBInstanceRequest()
		response   = mariadb.NewCreateDBInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, _ := d.GetOk("node_count"); v != nil {
		request.NodeCount = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("memory"); v != nil {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, _ := d.GetOk("auto_renew_flag"); v != nil {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("ipv6_flag"); v != nil {
		request.Ipv6Flag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			resourceTag := mariadb.ResourceTag{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("init_params"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dBParamValue := mariadb.DBParamValue{}
			if v, ok := dMap["param"]; ok {
				dBParamValue.Param = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				dBParamValue.Value = helper.String(v.(string))
			}
			request.InitParams = append(request.InitParams, &dBParamValue)
		}
	}

	if v, ok := d.GetOk("dcn_region"); ok {
		request.DcnRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dcn_instance_id"); ok {
		request.DcnInstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CreateDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mariadb instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeMariadbInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if instance == nil {
			err = fmt.Errorf("mariadb %s instance not exists", instanceId)
			return resource.NonRetryableError(err)
		}
		if *instance.Status == 0 || *instance.Status == 1 || *instance.Status == 4 {
			return resource.RetryableError(fmt.Errorf("create mariadb status is %v,start retrying ...", *instance.Status))
		}
		if *instance.Status == 2 {
			return nil
		}
		err = fmt.Errorf("create mariadb status is %v,we won't wait for it finish", *instance.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb fail, reason:%s\n ", logId, err.Error())
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instance, err := service.DescribeMariadbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance != nil {
		if instance.InstanceId != nil {
			_ = d.Set("instance_id", instance.InstanceId)
		}

		if instance.InstanceName != nil {
			_ = d.Set("instance_name", instance.InstanceName)
		}

		if instance.AppId != nil {
			_ = d.Set("app_id", instance.AppId)
		}

		if instance.ProjectId != nil {
			_ = d.Set("project_id", instance.ProjectId)
		}

		if instance.Region != nil {
			_ = d.Set("region", instance.Region)
		}

		if instance.Zone != nil {
			_ = d.Set("zones", []*string{instance.Zone})
		}

		if instance.VpcId != nil {
			_ = d.Set("vpc_id", instance.UniqueVpcId)
		}

		if instance.SubnetId != nil {
			_ = d.Set("subnet_id", instance.UniqueSubnetId)
		}

		// if instance.Period != nil {
		// 	_ = d.Set("period", instance.Period)
		// }

		// if instance.AutoVoucher != nil {
		// 	_ = d.Set("auto_voucher", instance.AutoVoucher)
		// }

		// if instance.VoucherIds != nil {
		// 	_ = d.Set("voucher_ids", instance.VoucherIds)
		// }

		if instance.Status != nil {
			_ = d.Set("status", instance.Status)
		}

		if instance.Vip != nil {
			_ = d.Set("vip", instance.Vip)
		}

		if instance.Vport != nil {
			_ = d.Set("vport", instance.Vport)
		}

		if instance.WanDomain != nil {
			_ = d.Set("wan_domain", instance.WanDomain)
		}

		if instance.WanVip != nil {
			_ = d.Set("wan_vip", instance.WanVip)
		}

		if instance.WanPort != nil {
			_ = d.Set("wan_port", instance.WanPort)
		}

		if instance.CreateTime != nil {
			_ = d.Set("create_time", instance.CreateTime)
		}

		if instance.UpdateTime != nil {
			_ = d.Set("update_time", instance.UpdateTime)
		}

		if instance.AutoRenewFlag != nil {
			_ = d.Set("auto_renew_flag", instance.AutoRenewFlag)
		}

		if instance.PeriodEndTime != nil {
			_ = d.Set("period_end_time", instance.PeriodEndTime)
		}

		if instance.Uin != nil {
			_ = d.Set("uin", instance.Uin)
		}

		if instance.TdsqlVersion != nil {
			_ = d.Set("tdsql_version", instance.TdsqlVersion)
		}

		if instance.Memory != nil {
			_ = d.Set("memory", instance.Memory)
		}

		if instance.Storage != nil {
			_ = d.Set("storage", instance.Storage)
		}

		// if instance.UniqueVpcId != nil {
		// 	_ = d.Set("unique_vpc_id", instance.UniqueVpcId)
		// }

		// if instance.UniqueSubnetId != nil {
		// 	_ = d.Set("unique_subnet_id", instance.UniqueSubnetId)
		// }

		if instance.NodeCount != nil {
			_ = d.Set("node_count", instance.NodeCount)
		}

		if instance.IsTmp != nil {
			_ = d.Set("is_tmp", instance.IsTmp)
		}

		if instance.ExclusterId != nil {
			_ = d.Set("excluster_id", instance.ExclusterId)
		}

		if instance.Pid != nil {
			_ = d.Set("pid", instance.Pid)
		}

		if instance.Qps != nil {
			_ = d.Set("qps", instance.Qps)
		}

		if instance.Paymode != nil {
			_ = d.Set("paymode", instance.Paymode)
		}

		if instance.Locker != nil {
			_ = d.Set("locker", instance.Locker)
		}

		if instance.StatusDesc != nil {
			_ = d.Set("status_desc", instance.StatusDesc)
		}

		if instance.WanStatus != nil {
			_ = d.Set("wan_status", instance.WanStatus)
		}

		if instance.IsAuditSupported != nil {
			_ = d.Set("is_audit_supported", instance.IsAuditSupported)
		}

		if instance.Machine != nil {
			_ = d.Set("machine", instance.Machine)
		}

		if instance.IsEncryptSupported != nil {
			_ = d.Set("is_encrypt_supported", instance.IsEncryptSupported)
		}

		if instance.Cpu != nil {
			_ = d.Set("cpu", instance.Cpu)
		}

		if instance.Ipv6Flag != nil {
			_ = d.Set("ipv6_flag", instance.Ipv6Flag)
		}

		if instance.Vipv6 != nil {
			_ = d.Set("vipv6", instance.Vipv6)
		}

		if instance.WanVipv6 != nil {
			_ = d.Set("wan_vipv6", instance.WanVipv6)
		}

		if instance.WanPortIpv6 != nil {
			_ = d.Set("wan_port_ipv6", instance.WanPortIpv6)
		}

		if instance.WanStatusIpv6 != nil {
			_ = d.Set("wan_status_ipv6", instance.WanStatusIpv6)
		}

		if instance.DbEngine != nil {
			_ = d.Set("db_engine", instance.DbEngine)
		}

		if instance.DbVersion != nil {
			_ = d.Set("db_version", instance.DbVersion)
		}

		if instance.DcnFlag != nil {
			_ = d.Set("dcn_flag", instance.DcnFlag)
		}

		if instance.DcnStatus != nil {
			_ = d.Set("dcn_status", instance.DcnStatus)
		}

		if instance.DcnDstNum != nil {
			_ = d.Set("dcn_dst_num", instance.DcnDstNum)
		}

		if instance.InstanceType != nil {
			_ = d.Set("instance_type", instance.InstanceType)
		}

		if instance.DbVersionId != nil {
			_ = d.Set("db_version_id", instance.DbVersionId)
		}

		// if instance.SecurityGroupIds != nil {
		// 	_ = d.Set("security_group_ids", instance.SecurityGroupIds)
		// }

		// if instance.InitParams != nil {
		// 	initParamsList := []interface{}{}
		// 	for _, initParams := range instance.InitParams {
		// 		initParamsMap := map[string]interface{}{}

		// 		if initParams.Param != nil {
		// 			initParamsMap["param", initParams.Param
		// 		}

		// 		if initParams.Value != nil {
		// 			initParamsMap["value", initParams.Value
		// 		}

		// 		initParamsList = append(initParamsList, initParamsMap)
		// 	}

		// 	_ = d.Set("init_params", initParamsList)

		// }

		// if instance.DcnRegion != nil {
		// 	_ = d.Set("dcn_region", instance.DcnRegion)
		// }

		// if instance.DcnInstanceId != nil {
		// 	_ = d.Set("dcn_instance_id", instance.DcnInstanceId)
		// }
	}

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region
	tags, err := tagService.DescribeResourceTags(ctx, "mariadb", "instance", region, instanceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMariadbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := mariadb.NewModifyDBInstanceNameRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"zones", "node_count", "memory", "storage", "period", "count", "auto_voucher", "voucher_ids", "vpc_id", "subnet_id", "project_id", "db_version_id", "security_group_ids", "auto_renew_flag", "ipv6_flag", "init_params", "dcn_region", "dcn_instance_id", "total_count", "instances"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBInstanceName(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mariadb instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("mariadb", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.IsolateDBInstanceById(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeMariadbInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if instance == nil {
			return nil
		}
		if *instance.Status == 2 {
			return resource.RetryableError(fmt.Errorf("isolate mariadb status is %v,start retrying ...", *instance.Status))
		}
		if *instance.Status == -1 {
			return nil
		}
		err := fmt.Errorf("isolate mariadb status is %v,we won't wait for it finish", *instance.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s isolate mariadb fail, reason:%s\n ", logId, err.Error())
		return err
	}

	if err := service.DeleteMariadbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeMariadbInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if instance == nil {
			return nil
		}

		if *instance.Status == -1 {
			return resource.RetryableError(fmt.Errorf("delete mariadb status is %v,start retrying ...", *instance.Status))
		}

		err := fmt.Errorf("delete mariadb status is %v,we won't wait for it finish", *instance.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mariadb fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
