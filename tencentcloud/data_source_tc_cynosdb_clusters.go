/*
Use this data source to query detailed information of cynosdb clusters

Example Usage

```hcl
data "tencentcloud_cynosdb_clusters" "clusters" {
  db_type = "MYSQL"
  limit = 20
  offset = 0
  order_by = &lt;nil&gt;
  order_by_type = &lt;nil&gt;
  filters {
		names = &lt;nil&gt;
		values = &lt;nil&gt;
		exact_match = &lt;nil&gt;
		name = &lt;nil&gt;
		operator = &lt;nil&gt;

  }
  total_count = &lt;nil&gt;
  cluster_set {
		status = ""
		update_time = &lt;nil&gt;
		zone = &lt;nil&gt;
		cluster_name = &lt;nil&gt;
		region = &lt;nil&gt;
		db_version = &lt;nil&gt;
		cluster_id = &lt;nil&gt;
		instance_num = &lt;nil&gt;
		uin = &lt;nil&gt;
		db_type = &lt;nil&gt;
		app_id = &lt;nil&gt;
		status_desc = &lt;nil&gt;
		create_time = ""
		pay_mode = &lt;nil&gt;
		period_end_time = &lt;nil&gt;
		vip = &lt;nil&gt;
		vport = &lt;nil&gt;
		project_i_d = &lt;nil&gt;
		vpc_id = &lt;nil&gt;
		subnet_id = &lt;nil&gt;
		cynos_version = &lt;nil&gt;
		storage_limit = &lt;nil&gt;
		renew_flag = &lt;nil&gt;
		processing_task = &lt;nil&gt;
		tasks {
			task_id = &lt;nil&gt;
			task_type = &lt;nil&gt;
			task_status = &lt;nil&gt;
			object_id = &lt;nil&gt;
			object_type = &lt;nil&gt;
		}
		resource_tags {
			tag_key = &lt;nil&gt;
			tag_value = &lt;nil&gt;
		}
		db_mode = &lt;nil&gt;
		serverless_status = &lt;nil&gt;
		storage = &lt;nil&gt;
		storage_id = &lt;nil&gt;
		storage_pay_mode = &lt;nil&gt;
		min_storage_size = &lt;nil&gt;
		max_storage_size = &lt;nil&gt;
		net_addrs {
			vip = &lt;nil&gt;
			vport = &lt;nil&gt;
			wan_domain = &lt;nil&gt;
			wan_port = &lt;nil&gt;
			net_type = &lt;nil&gt;
			uniq_subnet_id = &lt;nil&gt;
			uniq_vpc_id = &lt;nil&gt;
			description = &lt;nil&gt;
			wan_i_p = &lt;nil&gt;
			wan_status = &lt;nil&gt;
		}
		physical_zone = &lt;nil&gt;
		master_zone = &lt;nil&gt;
		has_slave_zone = &lt;nil&gt;
		slave_zones = &lt;nil&gt;
		business_type = &lt;nil&gt;
		is_freeze = &lt;nil&gt;
		order_source = &lt;nil&gt;
		ability {
			is_support_slave_zone = &lt;nil&gt;
			nonsupport_slave_zone_reason = &lt;nil&gt;
			is_support_ro = &lt;nil&gt;
			nonsupport_ro_reason = &lt;nil&gt;
		}

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClustersRead,
		Schema: map[string]*schema.Schema{
			"db_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine type: currently supports MYSQL, POSTGRESQL.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number to return, the default is 20, the maximum is 100.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Record offset, the default value is 0.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field, range of values:CREATETIME, PERIODENDTIME.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort type, range of values:ASC, DESC.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Search conditions, if there are multiple Filters, the relationship between Filters is a logical AND (AND) relationship.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"names": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Search field, currently supports: InstanceId, ProjectId, InstanceName, Vip.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Search string.",
						},
						"exact_match": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is it an exact match.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Search field.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator.",
						},
					},
				},
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The total count of clusters.",
			},

			"cluster_set": {
				Type:        schema.TypeList,
				Description: "The list of clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Description: "Cluster status, optional values are as follows:creating, running, isolating, isolated, activating, offlining, offlined, deleting, deleted.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Update time.",
						},
						"zone": {
							Type:        schema.TypeString,
							Description: "Availability zone.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Description: "The name of cluster.",
						},
						"region": {
							Type:        schema.TypeString,
							Description: "Region.",
						},
						"db_version": {
							Type:        schema.TypeString,
							Description: "Database version.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Description: "The ID of cluster.",
						},
						"instance_num": {
							Type:        schema.TypeInt,
							Description: "The number of instances.",
						},
						"uin": {
							Type:        schema.TypeString,
							Description: "User uin.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"db_type": {
							Type:        schema.TypeString,
							Description: "Engine type.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Description: "User app id.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Description: "Cluster Status Description.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Cluster creation time.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Description: "Payment mode. 0 - pay as your go, 1 - yearly and monthly.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"period_end_time": {
							Type:        schema.TypeString,
							Description: "Deadline.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"vip": {
							Type:        schema.TypeString,
							Description: "Cluster read and write vip.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Description: "Cluster read and write vport.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"project_i_d": {
							Type:        schema.TypeInt,
							Description: "Project id.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "Vpc id.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet id.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"cynos_version": {
							Type:        schema.TypeString,
							Description: "Cynos kernel version.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_limit": {
							Type:        schema.TypeInt,
							Description: "Storage limit.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Description: "Renew flag.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"processing_task": {
							Type:        schema.TypeString,
							Description: "Task being processed.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tasks": {
							Type:        schema.TypeList,
							Description: "Array of tasks for the cluster.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_id": {
										Type:        schema.TypeInt,
										Description: "Task auto-increment ID.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"task_type": {
										Type:        schema.TypeString,
										Description: "Task type.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"task_status": {
										Type:        schema.TypeString,
										Description: "Task status.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"object_id": {
										Type:        schema.TypeString,
										Description: "Task ID (cluster ID|instance group ID|instance ID).Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"object_type": {
										Type:        schema.TypeString,
										Description: "Object type.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"resource_tags": {
							Type:        schema.TypeList,
							Description: "Cluster bound tag array.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Description: "The key of tag.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Description: "The value of tag.",
									},
								},
							},
						},
						"db_mode": {
							Type:        schema.TypeString,
							Description: "Db mode (NORMAL, SERVERLESS).Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"serverless_status": {
							Type:        schema.TypeString,
							Description: "When the Db mode is SERVERLESS, serverless cluster status, optional value:resume, pause.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Description: "Cluster prepaid storage value size.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_id": {
							Type:        schema.TypeString,
							Description: "The storage ID when the cluster storage is prepaid, which is used for prepaid storage allocation.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"storage_pay_mode": {
							Type:        schema.TypeInt,
							Description: "Cluster storage payment model. 0 - pay as your go, 1 - yearly and monthly.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"min_storage_size": {
							Type:        schema.TypeInt,
							Description: "The minimum storage size corresponding to the cluster computing specification.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"max_storage_size": {
							Type:        schema.TypeInt,
							Description: "The maximum storage size corresponding to the cluster computing specification.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"net_addrs": {
							Type:        schema.TypeList,
							Description: "Cluster network information.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vip": {
										Type:        schema.TypeString,
										Description: "Intranet ip.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Description: "Intranet vport.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"wan_domain": {
										Type:        schema.TypeString,
										Description: "Public domain.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"wan_port": {
										Type:        schema.TypeInt,
										Description: "Public port.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"net_type": {
										Type:        schema.TypeString,
										Description: "Network type (ro-read-only, rw/ha-read-write).Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"uniq_subnet_id": {
										Type:        schema.TypeString,
										Description: "Subnet id.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Description: "Vpc id.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"description": {
										Type:        schema.TypeString,
										Description: "Description information.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"wan_i_p": {
										Type:        schema.TypeString,
										Description: "Public ip.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"wan_status": {
										Type:        schema.TypeString,
										Description: "Public status.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"physical_zone": {
							Type:        schema.TypeString,
							Description: "Physical availability zone.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"master_zone": {
							Type:        schema.TypeString,
							Description: "Master availability zone.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"has_slave_zone": {
							Type:        schema.TypeString,
							Description: "Is there a slave availability zone.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"slave_zones": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Slave availability zone.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"business_type": {
							Type:        schema.TypeString,
							Description: "Business type.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"is_freeze": {
							Type:        schema.TypeString,
							Description: "Whether to freeze.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"order_source": {
							Type:        schema.TypeString,
							Description: "Order source.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ability": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ability.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_support_slave_zone": {
										Type:        schema.TypeString,
										Description: "Whether to support slave availability zone.",
									},
									"nonsupport_slave_zone_reason": {
										Type:        schema.TypeString,
										Description: "The reason for not supporting slave availability zone.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"is_support_ro": {
										Type:        schema.TypeString,
										Description: "Whether to support Read-Only instance.",
									},
									"nonsupport_ro_reason": {
										Type:        schema.TypeString,
										Description: "Reasons for not supporting Read-Only instances.Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudCynosdbClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_clusters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("db_type"); ok {
		paramMap["DbType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cynosdb.QueryFilter, 0, len(filtersSet))

		for _, item := range filtersSet {
			queryFilter := cynosdb.QueryFilter{}
			queryFilterMap := item.(map[string]interface{})

			if v, ok := queryFilterMap["names"]; ok {
				namesSet := v.(*schema.Set).List()
				queryFilter.Names = helper.InterfacesStringsPoint(namesSet)
			}
			if v, ok := queryFilterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				queryFilter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			if v, ok := queryFilterMap["exact_match"]; ok {
				queryFilter.ExactMatch = helper.Bool(v.(bool))
			}
			if v, ok := queryFilterMap["name"]; ok {
				queryFilter.Name = helper.String(v.(string))
			}
			if v, ok := queryFilterMap["operator"]; ok {
				queryFilter.Operator = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &queryFilter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_set"); ok {
		clusterSetSet := v.([]interface{})
		tmpSet := make([]*cynosdb.CynosdbCluster, 0, len(clusterSetSet))

		for _, item := range clusterSetSet {
			cynosdbCluster := cynosdb.CynosdbCluster{}
			cynosdbClusterMap := item.(map[string]interface{})

			if v, ok := cynosdbClusterMap["status"]; ok {
				cynosdbCluster.Status = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["update_time"]; ok {
				cynosdbCluster.UpdateTime = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["zone"]; ok {
				cynosdbCluster.Zone = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["cluster_name"]; ok {
				cynosdbCluster.ClusterName = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["region"]; ok {
				cynosdbCluster.Region = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["db_version"]; ok {
				cynosdbCluster.DbVersion = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["cluster_id"]; ok {
				cynosdbCluster.ClusterId = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["instance_num"]; ok {
				cynosdbCluster.InstanceNum = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["uin"]; ok {
				cynosdbCluster.Uin = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["db_type"]; ok {
				cynosdbCluster.DbType = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["app_id"]; ok {
				cynosdbCluster.AppId = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["status_desc"]; ok {
				cynosdbCluster.StatusDesc = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["create_time"]; ok {
				cynosdbCluster.CreateTime = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["pay_mode"]; ok {
				cynosdbCluster.PayMode = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["period_end_time"]; ok {
				cynosdbCluster.PeriodEndTime = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["vip"]; ok {
				cynosdbCluster.Vip = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["vport"]; ok {
				cynosdbCluster.Vport = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["project_i_d"]; ok {
				cynosdbCluster.ProjectID = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["vpc_id"]; ok {
				cynosdbCluster.VpcId = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["subnet_id"]; ok {
				cynosdbCluster.SubnetId = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["cynos_version"]; ok {
				cynosdbCluster.CynosVersion = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["storage_limit"]; ok {
				cynosdbCluster.StorageLimit = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["renew_flag"]; ok {
				cynosdbCluster.RenewFlag = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["processing_task"]; ok {
				cynosdbCluster.ProcessingTask = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["tasks"]; ok {
				for _, item := range v.([]interface{}) {
					tasksMap := item.(map[string]interface{})
					objectTask := cynosdb.ObjectTask{}
					if v, ok := tasksMap["task_id"]; ok {
						objectTask.TaskId = helper.IntInt64(v.(int))
					}
					if v, ok := tasksMap["task_type"]; ok {
						objectTask.TaskType = helper.String(v.(string))
					}
					if v, ok := tasksMap["task_status"]; ok {
						objectTask.TaskStatus = helper.String(v.(string))
					}
					if v, ok := tasksMap["object_id"]; ok {
						objectTask.ObjectId = helper.String(v.(string))
					}
					if v, ok := tasksMap["object_type"]; ok {
						objectTask.ObjectType = helper.String(v.(string))
					}
					cynosdbCluster.Tasks = append(cynosdbCluster.Tasks, &objectTask)
				}
			}
			if v, ok := cynosdbClusterMap["resource_tags"]; ok {
				for _, item := range v.([]interface{}) {
					resourceTagsMap := item.(map[string]interface{})
					tag := cynosdb.Tag{}
					if v, ok := resourceTagsMap["tag_key"]; ok {
						tag.TagKey = helper.String(v.(string))
					}
					if v, ok := resourceTagsMap["tag_value"]; ok {
						tag.TagValue = helper.String(v.(string))
					}
					cynosdbCluster.ResourceTags = append(cynosdbCluster.ResourceTags, &tag)
				}
			}
			if v, ok := cynosdbClusterMap["db_mode"]; ok {
				cynosdbCluster.DbMode = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["serverless_status"]; ok {
				cynosdbCluster.ServerlessStatus = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["storage"]; ok {
				cynosdbCluster.Storage = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["storage_id"]; ok {
				cynosdbCluster.StorageId = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["storage_pay_mode"]; ok {
				cynosdbCluster.StoragePayMode = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["min_storage_size"]; ok {
				cynosdbCluster.MinStorageSize = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["max_storage_size"]; ok {
				cynosdbCluster.MaxStorageSize = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbClusterMap["net_addrs"]; ok {
				for _, item := range v.([]interface{}) {
					netAddrsMap := item.(map[string]interface{})
					netAddr := cynosdb.NetAddr{}
					if v, ok := netAddrsMap["vip"]; ok {
						netAddr.Vip = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["vport"]; ok {
						netAddr.Vport = helper.IntInt64(v.(int))
					}
					if v, ok := netAddrsMap["wan_domain"]; ok {
						netAddr.WanDomain = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["wan_port"]; ok {
						netAddr.WanPort = helper.IntInt64(v.(int))
					}
					if v, ok := netAddrsMap["net_type"]; ok {
						netAddr.NetType = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["uniq_subnet_id"]; ok {
						netAddr.UniqSubnetId = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["uniq_vpc_id"]; ok {
						netAddr.UniqVpcId = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["description"]; ok {
						netAddr.Description = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["wan_i_p"]; ok {
						netAddr.WanIP = helper.String(v.(string))
					}
					if v, ok := netAddrsMap["wan_status"]; ok {
						netAddr.WanStatus = helper.String(v.(string))
					}
					cynosdbCluster.NetAddrs = append(cynosdbCluster.NetAddrs, &netAddr)
				}
			}
			if v, ok := cynosdbClusterMap["physical_zone"]; ok {
				cynosdbCluster.PhysicalZone = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["master_zone"]; ok {
				cynosdbCluster.MasterZone = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["has_slave_zone"]; ok {
				cynosdbCluster.HasSlaveZone = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["slave_zones"]; ok {
				slaveZonesSet := v.(*schema.Set).List()
				cynosdbCluster.SlaveZones = helper.InterfacesStringsPoint(slaveZonesSet)
			}
			if v, ok := cynosdbClusterMap["business_type"]; ok {
				cynosdbCluster.BusinessType = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["is_freeze"]; ok {
				cynosdbCluster.IsFreeze = helper.String(v.(string))
			}
			if v, ok := cynosdbClusterMap["order_source"]; ok {
				cynosdbCluster.OrderSource = helper.String(v.(string))
			}
			if abilityMap, ok := helper.InterfaceToMap(cynosdbClusterMap, "ability"); ok {
				ability := cynosdb.Ability{}
				if v, ok := abilityMap["is_support_slave_zone"]; ok {
					ability.IsSupportSlaveZone = helper.String(v.(string))
				}
				if v, ok := abilityMap["nonsupport_slave_zone_reason"]; ok {
					ability.NonsupportSlaveZoneReason = helper.String(v.(string))
				}
				if v, ok := abilityMap["is_support_ro"]; ok {
					ability.IsSupportRo = helper.String(v.(string))
				}
				if v, ok := abilityMap["nonsupport_ro_reason"]; ok {
					ability.NonsupportRoReason = helper.String(v.(string))
				}
				cynosdbCluster.Ability = &ability
			}
			tmpSet = append(tmpSet, &cynosdbCluster)
		}
		paramMap["cluster_set"] = tmpSet
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterSet []*cynosdb.CynosdbCluster

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClustersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		clusterSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
