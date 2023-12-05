package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresqlReadonlyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlReadonlyGroupsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter condition. The primary ID must be specified in the format of db-master-instance-id to filter results, or else null will be returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "One or more filter values.",
						},
					},
				},
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting criterion. Valid values:ROGroupId, CreateTime, Name.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting order. Valid values:desc, asc.",
			},

			"read_only_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of read-only groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_only_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "read-only group idNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"read_only_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "read-only group nameNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project idNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"master_db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "master instance idNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"min_delay_eliminate_reserve": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum Number of Reserved InstancesNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"max_replay_latency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "delay space size threshold.",
						},
						"replay_latency_eliminate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "delay size switch.",
						},
						"max_replay_lag": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "delay time size threshold.",
						},
						"replay_lag_eliminate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "delay time switch.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "virtual network id.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet-idNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region id.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region id.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state.",
						},
						"read_only_db_instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "instance details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region to which the instance belongs, such as: ap-guangzhou, corresponding to the Region field of the RegionSet.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone to which the instance belongs, such as: ap-guangzhou-3, corresponding to the Zone field of ZoneSet.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "project ID.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private network ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet ID.",
									},
									"db_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance ID.",
									},
									"db_instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance name.",
									},
									"db_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance status, respectively: applying (applying), init (to be initialized), initing (initializing), running (running), limited run (limited run), isolated (isolated), recycling (recycling ), recycled (recycled), job running (task execution), offline (offline), migrating (migration), expanding (expanding), waitSwitch (waiting for switching), switching (switching), readonly (read-only ), restarting (restarting), network changing (network changing), upgrading (kernel version upgrade).",
									},
									"db_instance_memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the memory size allocated by the instance, unit: GB.",
									},
									"db_instance_storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the size of the storage space allocated by the instance, unit: GB.",
									},
									"db_instance_cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the number of CPUs allocated by the instance.",
									},
									"db_instance_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "sales specification ID.",
									},
									"db_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance type, the types are: 1. primary (primary instance); 2. readonly (read-only instance); 3. guard (disaster recovery instance); 4. temp (temporary instance).",
									},
									"db_instance_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance version, currently only supports standard (dual machine high availability version, one master and one slave).",
									},
									"db_charset": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance DB character set.",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "PostgreSQL version.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance creation time.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time when the instance performed the last update.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance expiration time.",
									},
									"isolated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance isolation time.",
									},
									"pay_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "billing mode, 1. prepaid (subscription, prepaid); 2. postpaid (billing by volume, postpaid).",
									},
									"auto_renew": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "auto-renew, 1: auto-renew, 0: no auto-renew.",
									},
									"db_instance_net_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "instance network connection information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DNS domain name.",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IP address.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "connection port address.",
												},
												"net_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "network type, 1. inner (intranet address of the basic network); 2. private (intranet address of the private network); 3. public (extranet address of the basic network or private network);.",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "network connection status, 1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "private network IDNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "subnet IDNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"protocol_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The protocol type for connecting to the database, currently supported: postgresql, mssql (MSSQL compatible syntax)Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "machine type.",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "user&#39;s AppId.",
									},
									"uid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Uid of the instance.",
									},
									"support_ipv6": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether the instance supports Ipv6, 1: support, 0: not support.",
									},
									"tag_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Label information bound to the instanceNote: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "label key.",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "tag value.",
												},
											},
										},
									},
									"master_db_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Master instance information, only returned when the instance is read-onlyNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"read_only_instance_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of read-only instancesNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"status_in_readonly_group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the read-only instance in the read-only groupNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"offline_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "offline timeNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"db_kernel_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database kernel versionNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"network_access_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Instance network information list (this field is obsolete)Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Network resource id, instance id or RO group idNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"resource_type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Resource type, 1-instance 2-RO groupNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "private network IDNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IPV4 addressNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"vip6": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IPV6 addressNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"vport": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "access portNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "subnet IDNote: This field may return null, indicating that no valid value can be obtained.",
												},
												"vpc_status": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Network status, 1-applying, 2-using, 3-deleting, 4-deletedNote: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"db_major_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "PostgreSQL major versionNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"db_node_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Instance node informationNote: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Node type, the value can be:Primary, representing the primary node;Standby, stands for standby node.",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Availability zone where the node is located, such as ap-guangzhou-1.",
												},
											},
										},
									},
									"is_support_t_d_e": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether the instance supports TDE data encryption 0: not supported, 1: supportedNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"db_engine": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database engine that supports:1. postgresql (cloud database PostgreSQL);2. mssql_compatible (MSSQL compatible - cloud database PostgreSQL);Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"db_engine_config": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration information for the database engineNote: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"rebalance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "automatic load balancing switch.",
						},
						"db_instance_net_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "network information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS domain name.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "connection port address.",
									},
									"net_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "network type, 1. inner (intranet address of the basic network); 2. private (intranet address of the private network); 3. public (extranet address of the basic network or private network);.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "network connection status, 1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private network IDNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet IDNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"protocol_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The protocol type for connecting to the database, currently supported: postgresql, mssql (MSSQL compatible syntax)Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"network_access_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Read-only list of group network information (this field is obsolete)Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network resource id, instance id or RO group idNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"resource_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource type, 1-instance 2-RO groupNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private network IDNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPV4 addressNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"vip6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPV6 addressNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "access portNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet IDNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"vpc_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network status, 1-applying, 2-using, 3-deleting, 4-deletedNote: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
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

func dataSourceTencentCloudPostgresqlReadonlyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgresql_readonly_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgresql.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgresql.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var readOnlyGroupList []*postgresql.ReadOnlyGroup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlReadonlyGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		readOnlyGroupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(readOnlyGroupList))
	tmpList := make([]map[string]interface{}, 0, len(readOnlyGroupList))

	if readOnlyGroupList != nil {
		for _, readOnlyGroup := range readOnlyGroupList {
			readOnlyGroupMap := map[string]interface{}{}

			if readOnlyGroup.ReadOnlyGroupId != nil {
				readOnlyGroupMap["read_only_group_id"] = readOnlyGroup.ReadOnlyGroupId
			}

			if readOnlyGroup.ReadOnlyGroupName != nil {
				readOnlyGroupMap["read_only_group_name"] = readOnlyGroup.ReadOnlyGroupName
			}

			if readOnlyGroup.ProjectId != nil {
				readOnlyGroupMap["project_id"] = readOnlyGroup.ProjectId
			}

			if readOnlyGroup.MasterDBInstanceId != nil {
				readOnlyGroupMap["master_db_instance_id"] = readOnlyGroup.MasterDBInstanceId
			}

			if readOnlyGroup.MinDelayEliminateReserve != nil {
				readOnlyGroupMap["min_delay_eliminate_reserve"] = readOnlyGroup.MinDelayEliminateReserve
			}

			if readOnlyGroup.MaxReplayLatency != nil {
				readOnlyGroupMap["max_replay_latency"] = readOnlyGroup.MaxReplayLatency
			}

			if readOnlyGroup.ReplayLatencyEliminate != nil {
				readOnlyGroupMap["replay_latency_eliminate"] = readOnlyGroup.ReplayLatencyEliminate
			}

			if readOnlyGroup.MaxReplayLag != nil {
				readOnlyGroupMap["max_replay_lag"] = readOnlyGroup.MaxReplayLag
			}

			if readOnlyGroup.ReplayLagEliminate != nil {
				readOnlyGroupMap["replay_lag_eliminate"] = readOnlyGroup.ReplayLagEliminate
			}

			if readOnlyGroup.VpcId != nil {
				readOnlyGroupMap["vpc_id"] = readOnlyGroup.VpcId
			}

			if readOnlyGroup.SubnetId != nil {
				readOnlyGroupMap["subnet_id"] = readOnlyGroup.SubnetId
			}

			if readOnlyGroup.Region != nil {
				readOnlyGroupMap["region"] = readOnlyGroup.Region
			}

			if readOnlyGroup.Zone != nil {
				readOnlyGroupMap["zone"] = readOnlyGroup.Zone
			}

			if readOnlyGroup.Status != nil {
				readOnlyGroupMap["status"] = readOnlyGroup.Status
			}

			if readOnlyGroup.ReadOnlyDBInstanceList != nil {
				readOnlyDbInstanceListList := []interface{}{}
				for _, readOnlyDbInstanceList := range readOnlyGroup.ReadOnlyDBInstanceList {
					readOnlyDbInstanceListMap := map[string]interface{}{}

					if readOnlyDbInstanceList.Region != nil {
						readOnlyDbInstanceListMap["region"] = readOnlyDbInstanceList.Region
					}

					if readOnlyDbInstanceList.Zone != nil {
						readOnlyDbInstanceListMap["zone"] = readOnlyDbInstanceList.Zone
					}

					if readOnlyDbInstanceList.ProjectId != nil {
						readOnlyDbInstanceListMap["project_id"] = readOnlyDbInstanceList.ProjectId
					}

					if readOnlyDbInstanceList.VpcId != nil {
						readOnlyDbInstanceListMap["vpc_id"] = readOnlyDbInstanceList.VpcId
					}

					if readOnlyDbInstanceList.SubnetId != nil {
						readOnlyDbInstanceListMap["subnet_id"] = readOnlyDbInstanceList.SubnetId
					}

					if readOnlyDbInstanceList.DBInstanceId != nil {
						readOnlyDbInstanceListMap["db_instance_id"] = readOnlyDbInstanceList.DBInstanceId
					}

					if readOnlyDbInstanceList.DBInstanceName != nil {
						readOnlyDbInstanceListMap["db_instance_name"] = readOnlyDbInstanceList.DBInstanceName
					}

					if readOnlyDbInstanceList.DBInstanceStatus != nil {
						readOnlyDbInstanceListMap["db_instance_status"] = readOnlyDbInstanceList.DBInstanceStatus
					}

					if readOnlyDbInstanceList.DBInstanceMemory != nil {
						readOnlyDbInstanceListMap["db_instance_memory"] = readOnlyDbInstanceList.DBInstanceMemory
					}

					if readOnlyDbInstanceList.DBInstanceStorage != nil {
						readOnlyDbInstanceListMap["db_instance_storage"] = readOnlyDbInstanceList.DBInstanceStorage
					}

					if readOnlyDbInstanceList.DBInstanceCpu != nil {
						readOnlyDbInstanceListMap["db_instance_cpu"] = readOnlyDbInstanceList.DBInstanceCpu
					}

					if readOnlyDbInstanceList.DBInstanceClass != nil {
						readOnlyDbInstanceListMap["db_instance_class"] = readOnlyDbInstanceList.DBInstanceClass
					}

					if readOnlyDbInstanceList.DBInstanceType != nil {
						readOnlyDbInstanceListMap["db_instance_type"] = readOnlyDbInstanceList.DBInstanceType
					}

					if readOnlyDbInstanceList.DBInstanceVersion != nil {
						readOnlyDbInstanceListMap["db_instance_version"] = readOnlyDbInstanceList.DBInstanceVersion
					}

					if readOnlyDbInstanceList.DBCharset != nil {
						readOnlyDbInstanceListMap["db_charset"] = readOnlyDbInstanceList.DBCharset
					}

					if readOnlyDbInstanceList.DBVersion != nil {
						readOnlyDbInstanceListMap["db_version"] = readOnlyDbInstanceList.DBVersion
					}

					if readOnlyDbInstanceList.CreateTime != nil {
						readOnlyDbInstanceListMap["create_time"] = readOnlyDbInstanceList.CreateTime
					}

					if readOnlyDbInstanceList.UpdateTime != nil {
						readOnlyDbInstanceListMap["update_time"] = readOnlyDbInstanceList.UpdateTime
					}

					if readOnlyDbInstanceList.ExpireTime != nil {
						readOnlyDbInstanceListMap["expire_time"] = readOnlyDbInstanceList.ExpireTime
					}

					if readOnlyDbInstanceList.IsolatedTime != nil {
						readOnlyDbInstanceListMap["isolated_time"] = readOnlyDbInstanceList.IsolatedTime
					}

					if readOnlyDbInstanceList.PayType != nil {
						readOnlyDbInstanceListMap["pay_type"] = readOnlyDbInstanceList.PayType
					}

					if readOnlyDbInstanceList.AutoRenew != nil {
						readOnlyDbInstanceListMap["auto_renew"] = readOnlyDbInstanceList.AutoRenew
					}

					if readOnlyDbInstanceList.DBInstanceNetInfo != nil {
						dbInstanceNetInfoList := []interface{}{}
						for _, dbInstanceNetInfo := range readOnlyDbInstanceList.DBInstanceNetInfo {
							dbInstanceNetInfoMap := map[string]interface{}{}

							if dbInstanceNetInfo.Address != nil {
								dbInstanceNetInfoMap["address"] = dbInstanceNetInfo.Address
							}

							if dbInstanceNetInfo.Ip != nil {
								dbInstanceNetInfoMap["ip"] = dbInstanceNetInfo.Ip
							}

							if dbInstanceNetInfo.Port != nil {
								dbInstanceNetInfoMap["port"] = dbInstanceNetInfo.Port
							}

							if dbInstanceNetInfo.NetType != nil {
								dbInstanceNetInfoMap["net_type"] = dbInstanceNetInfo.NetType
							}

							if dbInstanceNetInfo.Status != nil {
								dbInstanceNetInfoMap["status"] = dbInstanceNetInfo.Status
							}

							if dbInstanceNetInfo.VpcId != nil {
								dbInstanceNetInfoMap["vpc_id"] = dbInstanceNetInfo.VpcId
							}

							if dbInstanceNetInfo.SubnetId != nil {
								dbInstanceNetInfoMap["subnet_id"] = dbInstanceNetInfo.SubnetId
							}

							if dbInstanceNetInfo.ProtocolType != nil {
								dbInstanceNetInfoMap["protocol_type"] = dbInstanceNetInfo.ProtocolType
							}

							dbInstanceNetInfoList = append(dbInstanceNetInfoList, dbInstanceNetInfoMap)
						}

						readOnlyDbInstanceListMap["db_instance_net_info"] = dbInstanceNetInfoList
					}

					if readOnlyDbInstanceList.Type != nil {
						readOnlyDbInstanceListMap["type"] = readOnlyDbInstanceList.Type
					}

					if readOnlyDbInstanceList.AppId != nil {
						readOnlyDbInstanceListMap["app_id"] = readOnlyDbInstanceList.AppId
					}

					if readOnlyDbInstanceList.Uid != nil {
						readOnlyDbInstanceListMap["uid"] = readOnlyDbInstanceList.Uid
					}

					if readOnlyDbInstanceList.SupportIpv6 != nil {
						readOnlyDbInstanceListMap["support_ipv6"] = readOnlyDbInstanceList.SupportIpv6
					}

					if readOnlyDbInstanceList.TagList != nil {
						tagListList := []interface{}{}
						for _, tagList := range readOnlyDbInstanceList.TagList {
							tagListMap := map[string]interface{}{}

							if tagList.TagKey != nil {
								tagListMap["tag_key"] = tagList.TagKey
							}

							if tagList.TagValue != nil {
								tagListMap["tag_value"] = tagList.TagValue
							}

							tagListList = append(tagListList, tagListMap)
						}

						readOnlyDbInstanceListMap["tag_list"] = tagListList
					}

					if readOnlyDbInstanceList.MasterDBInstanceId != nil {
						readOnlyDbInstanceListMap["master_db_instance_id"] = readOnlyDbInstanceList.MasterDBInstanceId
					}

					if readOnlyDbInstanceList.ReadOnlyInstanceNum != nil {
						readOnlyDbInstanceListMap["read_only_instance_num"] = readOnlyDbInstanceList.ReadOnlyInstanceNum
					}

					if readOnlyDbInstanceList.StatusInReadonlyGroup != nil {
						readOnlyDbInstanceListMap["status_in_readonly_group"] = readOnlyDbInstanceList.StatusInReadonlyGroup
					}

					if readOnlyDbInstanceList.OfflineTime != nil {
						readOnlyDbInstanceListMap["offline_time"] = readOnlyDbInstanceList.OfflineTime
					}

					if readOnlyDbInstanceList.DBKernelVersion != nil {
						readOnlyDbInstanceListMap["db_kernel_version"] = readOnlyDbInstanceList.DBKernelVersion
					}

					if readOnlyDbInstanceList.NetworkAccessList != nil {
						networkAccessListList := []interface{}{}
						for _, networkAccessList := range readOnlyDbInstanceList.NetworkAccessList {
							networkAccessListMap := map[string]interface{}{}

							if networkAccessList.ResourceId != nil {
								networkAccessListMap["resource_id"] = networkAccessList.ResourceId
							}

							if networkAccessList.ResourceType != nil {
								networkAccessListMap["resource_type"] = networkAccessList.ResourceType
							}

							if networkAccessList.VpcId != nil {
								networkAccessListMap["vpc_id"] = networkAccessList.VpcId
							}

							if networkAccessList.Vip != nil {
								networkAccessListMap["vip"] = networkAccessList.Vip
							}

							if networkAccessList.Vip6 != nil {
								networkAccessListMap["vip6"] = networkAccessList.Vip6
							}

							if networkAccessList.Vport != nil {
								networkAccessListMap["vport"] = networkAccessList.Vport
							}

							if networkAccessList.SubnetId != nil {
								networkAccessListMap["subnet_id"] = networkAccessList.SubnetId
							}

							if networkAccessList.VpcStatus != nil {
								networkAccessListMap["vpc_status"] = networkAccessList.VpcStatus
							}

							networkAccessListList = append(networkAccessListList, networkAccessListMap)
						}

						readOnlyDbInstanceListMap["network_access_list"] = networkAccessListList
					}

					if readOnlyDbInstanceList.DBMajorVersion != nil {
						readOnlyDbInstanceListMap["db_major_version"] = readOnlyDbInstanceList.DBMajorVersion
					}

					if readOnlyDbInstanceList.DBNodeSet != nil {
						dbNodeSetList := []interface{}{}
						for _, dbNodeSet := range readOnlyDbInstanceList.DBNodeSet {
							dbNodeSetMap := map[string]interface{}{}

							if dbNodeSet.Role != nil {
								dbNodeSetMap["role"] = dbNodeSet.Role
							}

							if dbNodeSet.Zone != nil {
								dbNodeSetMap["zone"] = dbNodeSet.Zone
							}

							dbNodeSetList = append(dbNodeSetList, dbNodeSetMap)
						}

						readOnlyDbInstanceListMap["db_node_set"] = dbNodeSetList
					}

					if readOnlyDbInstanceList.IsSupportTDE != nil {
						readOnlyDbInstanceListMap["is_support_t_d_e"] = readOnlyDbInstanceList.IsSupportTDE
					}

					if readOnlyDbInstanceList.DBEngine != nil {
						readOnlyDbInstanceListMap["db_engine"] = readOnlyDbInstanceList.DBEngine
					}

					if readOnlyDbInstanceList.DBEngineConfig != nil {
						readOnlyDbInstanceListMap["db_engine_config"] = readOnlyDbInstanceList.DBEngineConfig
					}

					readOnlyDbInstanceListList = append(readOnlyDbInstanceListList, readOnlyDbInstanceListMap)
				}

				readOnlyGroupMap["read_only_db_instance_list"] = readOnlyDbInstanceListList
			}

			if readOnlyGroup.Rebalance != nil {
				readOnlyGroupMap["rebalance"] = readOnlyGroup.Rebalance
			}

			if readOnlyGroup.DBInstanceNetInfo != nil {
				dbInstanceNetInfoList := []interface{}{}
				for _, dbInstanceNetInfo := range readOnlyGroup.DBInstanceNetInfo {
					dbInstanceNetInfoMap := map[string]interface{}{}

					if dbInstanceNetInfo.Address != nil {
						dbInstanceNetInfoMap["address"] = dbInstanceNetInfo.Address
					}

					if dbInstanceNetInfo.Ip != nil {
						dbInstanceNetInfoMap["ip"] = dbInstanceNetInfo.Ip
					}

					if dbInstanceNetInfo.Port != nil {
						dbInstanceNetInfoMap["port"] = dbInstanceNetInfo.Port
					}

					if dbInstanceNetInfo.NetType != nil {
						dbInstanceNetInfoMap["net_type"] = dbInstanceNetInfo.NetType
					}

					if dbInstanceNetInfo.Status != nil {
						dbInstanceNetInfoMap["status"] = dbInstanceNetInfo.Status
					}

					if dbInstanceNetInfo.VpcId != nil {
						dbInstanceNetInfoMap["vpc_id"] = dbInstanceNetInfo.VpcId
					}

					if dbInstanceNetInfo.SubnetId != nil {
						dbInstanceNetInfoMap["subnet_id"] = dbInstanceNetInfo.SubnetId
					}

					if dbInstanceNetInfo.ProtocolType != nil {
						dbInstanceNetInfoMap["protocol_type"] = dbInstanceNetInfo.ProtocolType
					}

					dbInstanceNetInfoList = append(dbInstanceNetInfoList, dbInstanceNetInfoMap)
				}

				readOnlyGroupMap["db_instance_net_info"] = dbInstanceNetInfoList
			}

			if readOnlyGroup.NetworkAccessList != nil {
				networkAccessListList := []interface{}{}
				for _, networkAccessList := range readOnlyGroup.NetworkAccessList {
					networkAccessListMap := map[string]interface{}{}

					if networkAccessList.ResourceId != nil {
						networkAccessListMap["resource_id"] = networkAccessList.ResourceId
					}

					if networkAccessList.ResourceType != nil {
						networkAccessListMap["resource_type"] = networkAccessList.ResourceType
					}

					if networkAccessList.VpcId != nil {
						networkAccessListMap["vpc_id"] = networkAccessList.VpcId
					}

					if networkAccessList.Vip != nil {
						networkAccessListMap["vip"] = networkAccessList.Vip
					}

					if networkAccessList.Vip6 != nil {
						networkAccessListMap["vip6"] = networkAccessList.Vip6
					}

					if networkAccessList.Vport != nil {
						networkAccessListMap["vport"] = networkAccessList.Vport
					}

					if networkAccessList.SubnetId != nil {
						networkAccessListMap["subnet_id"] = networkAccessList.SubnetId
					}

					if networkAccessList.VpcStatus != nil {
						networkAccessListMap["vpc_status"] = networkAccessList.VpcStatus
					}

					networkAccessListList = append(networkAccessListList, networkAccessListMap)
				}

				readOnlyGroupMap["network_access_list"] = networkAccessListList
			}

			ids = append(ids, *readOnlyGroup.ReadOnlyGroupId)
			tmpList = append(tmpList, readOnlyGroupMap)
		}

		_ = d.Set("read_only_group_list", tmpList)
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
