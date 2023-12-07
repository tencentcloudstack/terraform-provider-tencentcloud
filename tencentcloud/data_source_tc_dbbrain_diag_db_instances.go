package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainDiagDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagDbInstancesRead,
		Schema: map[string]*schema.Schema{
			"is_supported": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "whether it is an instance supported by DBbrain, always pass `true`.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database TDSQL-C for MySQL, the default is `mysql`.",
			},

			"instance_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "query based on the instance name condition.",
			},

			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "query based on the instance ID condition.",
			},

			"regions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "query based on geographical conditions.",
			},

			"db_scan_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "all-instance inspection status. `0`: All-instance inspection is enabled; `1`: All-instance inspection is not enabled.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "information about the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"health_score": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "health score.",
						},
						"product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "belongs to the product.",
						},
						"event_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the number of abnormal events.",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "instance type. `1`: MASTER; `2`: DR, `3`: RO, `4`: SDR.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "memory, in MB.",
						},
						"volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "hard disk storage, in GB.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "database version.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "intranet address.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "intranet port.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access source.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "group name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status: `0`: Shipping; `1`: Running normally; `4`: Destroying; `5`: Isolating.",
						},
						"uniq_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet uniform ID.",
						},
						"deploy_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cdb type.",
						},
						"init_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "cdb instance initialization flag: `0`: not initialized; `1`: initialized.",
						},
						"task_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "task status.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the unified ID of the private network.",
						},
						"instance_conf": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "status of instance inspection/overview.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"daily_inspection": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "database inspection switch, Yes/No.",
									},
									"overview_display": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance overview switch, Yes/No.",
									},
									"key_delimiters": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Custom separator for redis large key analysis, only used by `redis`. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"deadline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "resource expiration time.",
						},
						"is_supported": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether it is an instance supported by DBbrain.",
						},
						"sec_audit_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "enabled status of the instance security audit log. `ON`: security audit is enabled; `OFF`: security audit is not enabled.",
						},
						"audit_policy_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance audit log enable status. `ALL_AUDIT`: full audit is enabled; `RULE_AUDIT`: rule audit is enabled; `UNBOUND`: audit is not enabled.",
						},
						"audit_running_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance audit log running status. `normal`: running; `paused`: arrears suspended.",
						},
						"internal_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet VIPNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"internal_vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Intranet portNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDiagDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_diag_db_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("is_supported"); v != nil {
		paramMap["IsSupported"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_names"); ok {
		instanceNamesSet := v.(*schema.Set).List()
		paramMap["InstanceNames"] = helper.InterfacesStringsPoint(instanceNamesSet)
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("regions"); ok {
		regionsSet := v.(*schema.Set).List()
		paramMap["Regions"] = helper.InterfacesStringsPoint(regionsSet)
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		infos        []*dbbrain.InstanceInfo
		dbScanStatus *int64
		e            error
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, dbScanStatus, e = service.DescribeDbbrainDiagDbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(infos))
	tmpList := make([]map[string]interface{}, 0, len(infos))
	if dbScanStatus != nil {
		_ = d.Set("db_scan_status", dbScanStatus)
	}

	if infos != nil {

		for _, instanceInfo := range infos {
			instanceInfoMap := map[string]interface{}{}

			if instanceInfo.InstanceId != nil {
				instanceInfoMap["instance_id"] = instanceInfo.InstanceId
			}

			if instanceInfo.InstanceName != nil {
				instanceInfoMap["instance_name"] = instanceInfo.InstanceName
			}

			if instanceInfo.Region != nil {
				instanceInfoMap["region"] = instanceInfo.Region
			}

			if instanceInfo.HealthScore != nil {
				instanceInfoMap["health_score"] = instanceInfo.HealthScore
			}

			if instanceInfo.Product != nil {
				instanceInfoMap["product"] = instanceInfo.Product
			}

			if instanceInfo.EventCount != nil {
				instanceInfoMap["event_count"] = instanceInfo.EventCount
			}

			if instanceInfo.InstanceType != nil {
				instanceInfoMap["instance_type"] = instanceInfo.InstanceType
			}

			if instanceInfo.Cpu != nil {
				instanceInfoMap["cpu"] = instanceInfo.Cpu
			}

			if instanceInfo.Memory != nil {
				instanceInfoMap["memory"] = instanceInfo.Memory
			}

			if instanceInfo.Volume != nil {
				instanceInfoMap["volume"] = instanceInfo.Volume
			}

			if instanceInfo.EngineVersion != nil {
				instanceInfoMap["engine_version"] = instanceInfo.EngineVersion
			}

			if instanceInfo.Vip != nil {
				instanceInfoMap["vip"] = instanceInfo.Vip
			}

			if instanceInfo.Vport != nil {
				instanceInfoMap["vport"] = instanceInfo.Vport
			}

			if instanceInfo.Source != nil {
				instanceInfoMap["source"] = instanceInfo.Source
			}

			if instanceInfo.GroupId != nil {
				instanceInfoMap["group_id"] = instanceInfo.GroupId
			}

			if instanceInfo.GroupName != nil {
				instanceInfoMap["group_name"] = instanceInfo.GroupName
			}

			if instanceInfo.Status != nil {
				instanceInfoMap["status"] = instanceInfo.Status
			}

			if instanceInfo.UniqSubnetId != nil {
				instanceInfoMap["uniq_subnet_id"] = instanceInfo.UniqSubnetId
			}

			if instanceInfo.DeployMode != nil {
				instanceInfoMap["deploy_mode"] = instanceInfo.DeployMode
			}

			if instanceInfo.InitFlag != nil {
				instanceInfoMap["init_flag"] = instanceInfo.InitFlag
			}

			if instanceInfo.TaskStatus != nil {
				instanceInfoMap["task_status"] = instanceInfo.TaskStatus
			}

			if instanceInfo.UniqVpcId != nil {
				instanceInfoMap["uniq_vpc_id"] = instanceInfo.UniqVpcId
			}

			if instanceInfo.InstanceConf != nil {
				instanceConfMap := map[string]interface{}{}

				if instanceInfo.InstanceConf.DailyInspection != nil {
					instanceConfMap["daily_inspection"] = instanceInfo.InstanceConf.DailyInspection
				}

				if instanceInfo.InstanceConf.OverviewDisplay != nil {
					instanceConfMap["overview_display"] = instanceInfo.InstanceConf.OverviewDisplay
				}

				if instanceInfo.InstanceConf.KeyDelimiters != nil {
					instanceConfMap["key_delimiters"] = instanceInfo.InstanceConf.KeyDelimiters
				}

				instanceInfoMap["instance_conf"] = []interface{}{instanceConfMap}
			}

			if instanceInfo.DeadlineTime != nil {
				instanceInfoMap["deadline_time"] = instanceInfo.DeadlineTime
			}

			if instanceInfo.IsSupported != nil {
				instanceInfoMap["is_supported"] = instanceInfo.IsSupported
			}

			if instanceInfo.SecAuditStatus != nil {
				instanceInfoMap["sec_audit_status"] = instanceInfo.SecAuditStatus
			}

			if instanceInfo.AuditPolicyStatus != nil {
				instanceInfoMap["audit_policy_status"] = instanceInfo.AuditPolicyStatus
			}

			if instanceInfo.AuditRunningStatus != nil {
				instanceInfoMap["audit_running_status"] = instanceInfo.AuditRunningStatus
			}

			if instanceInfo.InternalVip != nil {
				instanceInfoMap["internal_vip"] = instanceInfo.InternalVip
			}

			if instanceInfo.InternalVport != nil {
				instanceInfoMap["internal_vport"] = instanceInfo.InternalVport
			}

			if instanceInfo.CreateTime != nil {
				instanceInfoMap["create_time"] = instanceInfo.CreateTime
			}

			ids = append(ids, *instanceInfo.InstanceId)
			tmpList = append(tmpList, instanceInfoMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
