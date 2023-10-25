/*
Use this data source to query detailed information of dlc describe_data_engine

Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine" "describe_data_engine" {
  data_engine_name = "testSpark"
  }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
)

func dataSourceTencentCloudDlcDescribeDataEngine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeDataEngineRead,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine name.",
			},

			"data_engine": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data engine details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_engine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine name.",
						},
						"engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine type, only support: spark/presto.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine cluster type, only support: spark_cu/presto_cu.",
						},
						"quota_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference ID.",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine state, only support: 0:Init/-1:Failed/-2:Deleted/1:Pause/2:Running/3:ToBeDelete/4:Deleting.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Create time.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update time.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine size.",
						},
						"mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine mode, only support 1: ByAmount, 2: YearlyAndMonthly.",
						},
						"min_clusters": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine min size, greater than or equal to 1 and MaxClusters bigger than MinClusters.",
						},
						"max_clusters": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine max cluster size,  MaxClusters less than or equal to 10 and MaxClusters bigger than MinClusters.",
						},
						"auto_resume": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to automatically start the cluster, prepay not support.",
						},
						"spend_after": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Automatic recovery time, prepay not support.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine VPC network segment, just like 192.0.2.1/24.",
						},
						"default_data_engine": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is the default virtual cluster.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine description information.",
						},
						"data_engine_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine unique id.",
						},
						"sub_account_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operator.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expire time.",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Isolated time.",
						},
						"reversal_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reversal time.",
						},
						"user_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User&amp;#39;s nickname.",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"permissions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Engine permissions.",
						},
						"auto_suspend": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to automatically suspend the cluster, prepay not support.",
						},
						"crontab_resume_suspend": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine crontab resume or suspend strategy, only support: 0: Wait(default), 1: Kill.",
						},
						"crontab_resume_suspend_strategy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Engine auto suspend strategy, when AutoSuspend is true, CrontabResumeSuspend must stop.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resume_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scheduled pull-up time: For example: 8 o&amp;#39;clock on Monday is expressed as 1000000-08:00:00.",
									},
									"suspend_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scheduled suspension time: For example: 20 o&amp;#39;clock on Monday is expressed as 1000000-20:00:00.",
									},
									"suspend_strategy": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Suspend configuration: 0 (default): wait for the task to end before suspending, 1: force suspend.",
									},
								},
							},
						},
						"engine_exec_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine exec type, only support SQL(default) or BATCH.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Automatic renewal flag, 0, initial state, automatic renewal is not performed by default. If the user has prepaid non-stop service privileges, automatic renewal will occur. 1: Automatic renewal. 2: Make it clear that there will be no automatic renewal.",
						},
						"auto_suspend_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster automatic suspension time, default 10 minutes.",
						},
						"network_connection_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network connection configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network configuration id.",
									},
									"associate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network configuration unique identifier.",
									},
									"house_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data engine id.",
									},
									"datasource_connection_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data source id (obsolete).",
									},
									"state": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network configuration status (0-initialization, 1-normal).",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Create time.",
									},
									"update_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Update time.",
									},
									"appid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "User appid.",
									},
									"house_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data engine name.",
									},
									"datasource_connection_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network configuration name.",
									},
									"network_connection_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network configuration type.",
									},
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User uin.",
									},
									"sub_account_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User sub uin.",
									},
									"network_connection_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network configuration description.",
									},
									"datasource_connection_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource vpcid.",
									},
									"datasource_connection_subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource subnetId.",
									},
									"datasource_connection_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource connection cidr block.",
									},
									"datasource_connection_subnet_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource connection subnet cidr block.",
									},
								},
							},
						},
						"ui_u_r_l": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Jump address of ui.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine resource type not match, only support: Standard_CU/Memory_CU(only BATCH ExecType).",
						},
						"image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine major version id.",
						},
						"child_image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine Image version id.",
						},
						"image_version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine image version name.",
						},
						"start_standby_cluster": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable the backup cluster.",
						},
						"elastic_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "For spark Batch ExecType, yearly and monthly cluster whether to enable elasticity.",
						},
						"elastic_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For spark Batch ExecType, yearly and monthly cluster elastic limit.",
						},
						"default_house": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is it the default engine?.",
						},
						"max_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of concurrent tasks in a single cluster, default 5.",
						},
						"tolerable_queue_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Tolerable queuing time, default 0. scaling may be triggered when tasks are queued for longer than the tolerable time. if this parameter is 0, it means that capacity expansion may be triggered immediately once a task is queued.",
						},
						"user_app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User appid.",
						},
						"user_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User uin.",
						},
						"session_resource_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "For spark Batch ExecType, cluster session resource configuration template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"driver_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine driver size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.",
									},
									"executor_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine executor size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.",
									},
									"executor_nums": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specify the number of executors. The minimum value is 1 and the maximum value is less than the cluster specification.",
									},
									"executor_max_numbers": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specify the executor max number (in a dynamic configuration scenario), the minimum value is 1, and the maximum value is less than the cluster specification (when ExecutorMaxNumbers is less than ExecutorNums, the value is set to ExecutorNums).",
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

func dataSourceTencentCloudDlcDescribeDataEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_describe_data_engine.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var dataEngineName string
	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var dataEngine *dlc.DataEngineInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineByName(ctx, dataEngineName)
		if e != nil {
			return retryError(e)
		}
		dataEngine = result
		return nil
	})
	if err != nil {
		return err
	}
	dataEngineInfoMap := map[string]interface{}{}

	if dataEngine != nil {

		if dataEngine.DataEngineName != nil {
			dataEngineInfoMap["data_engine_name"] = dataEngine.DataEngineName
		}

		if dataEngine.EngineType != nil {
			dataEngineInfoMap["engine_type"] = dataEngine.EngineType
		}

		if dataEngine.ClusterType != nil {
			dataEngineInfoMap["cluster_type"] = dataEngine.ClusterType
		}

		if dataEngine.QuotaId != nil {
			dataEngineInfoMap["quota_id"] = dataEngine.QuotaId
		}

		if dataEngine.State != nil {
			dataEngineInfoMap["state"] = dataEngine.State
		}

		if dataEngine.CreateTime != nil {
			dataEngineInfoMap["create_time"] = dataEngine.CreateTime
		}

		if dataEngine.UpdateTime != nil {
			dataEngineInfoMap["update_time"] = dataEngine.UpdateTime
		}

		if dataEngine.Size != nil {
			dataEngineInfoMap["size"] = dataEngine.Size
		}

		if dataEngine.Mode != nil {
			dataEngineInfoMap["mode"] = dataEngine.Mode
		}

		if dataEngine.MinClusters != nil {
			dataEngineInfoMap["min_clusters"] = dataEngine.MinClusters
		}

		if dataEngine.MaxClusters != nil {
			dataEngineInfoMap["max_clusters"] = dataEngine.MaxClusters
		}

		if dataEngine.AutoResume != nil {
			dataEngineInfoMap["auto_resume"] = dataEngine.AutoResume
		}

		if dataEngine.SpendAfter != nil {
			dataEngineInfoMap["spend_after"] = dataEngine.SpendAfter
		}

		if dataEngine.CidrBlock != nil {
			dataEngineInfoMap["cidr_block"] = dataEngine.CidrBlock
		}

		if dataEngine.DefaultDataEngine != nil {
			dataEngineInfoMap["default_data_engine"] = dataEngine.DefaultDataEngine
		}

		if dataEngine.Message != nil {
			dataEngineInfoMap["message"] = dataEngine.Message
		}

		if dataEngine.DataEngineId != nil {
			dataEngineInfoMap["data_engine_id"] = dataEngine.DataEngineId
		}

		if dataEngine.SubAccountUin != nil {
			dataEngineInfoMap["sub_account_uin"] = dataEngine.SubAccountUin
		}

		if dataEngine.ExpireTime != nil {
			dataEngineInfoMap["expire_time"] = dataEngine.ExpireTime
		}

		if dataEngine.IsolatedTime != nil {
			dataEngineInfoMap["isolated_time"] = dataEngine.IsolatedTime
		}

		if dataEngine.ReversalTime != nil {
			dataEngineInfoMap["reversal_time"] = dataEngine.ReversalTime
		}

		if dataEngine.UserAlias != nil {
			dataEngineInfoMap["user_alias"] = dataEngine.UserAlias
		}

		if dataEngine.TagList != nil {
			var tagListList []interface{}
			for _, tagList := range dataEngine.TagList {
				tagListMap := map[string]interface{}{}

				if tagList.TagKey != nil {
					tagListMap["tag_key"] = tagList.TagKey
				}

				if tagList.TagValue != nil {
					tagListMap["tag_value"] = tagList.TagValue
				}

				tagListList = append(tagListList, tagListMap)
			}

			dataEngineInfoMap["tag_list"] = tagListList
		}

		if dataEngine.Permissions != nil {
			dataEngineInfoMap["permissions"] = dataEngine.Permissions
		}

		if dataEngine.AutoSuspend != nil {
			dataEngineInfoMap["auto_suspend"] = dataEngine.AutoSuspend
		}

		if dataEngine.CrontabResumeSuspend != nil {
			dataEngineInfoMap["crontab_resume_suspend"] = dataEngine.CrontabResumeSuspend
		}

		if dataEngine.CrontabResumeSuspendStrategy != nil {
			crontabResumeSuspendStrategyMap := map[string]interface{}{}

			if dataEngine.CrontabResumeSuspendStrategy.ResumeTime != nil {
				crontabResumeSuspendStrategyMap["resume_time"] = dataEngine.CrontabResumeSuspendStrategy.ResumeTime
			}

			if dataEngine.CrontabResumeSuspendStrategy.SuspendTime != nil {
				crontabResumeSuspendStrategyMap["suspend_time"] = dataEngine.CrontabResumeSuspendStrategy.SuspendTime
			}

			if dataEngine.CrontabResumeSuspendStrategy.SuspendStrategy != nil {
				crontabResumeSuspendStrategyMap["suspend_strategy"] = dataEngine.CrontabResumeSuspendStrategy.SuspendStrategy
			}

			dataEngineInfoMap["crontab_resume_suspend_strategy"] = []interface{}{crontabResumeSuspendStrategyMap}
		}

		if dataEngine.EngineExecType != nil {
			dataEngineInfoMap["engine_exec_type"] = dataEngine.EngineExecType
		}

		if dataEngine.RenewFlag != nil {
			dataEngineInfoMap["renew_flag"] = dataEngine.RenewFlag
		}

		if dataEngine.AutoSuspendTime != nil {
			dataEngineInfoMap["auto_suspend_time"] = dataEngine.AutoSuspendTime
		}

		if dataEngine.NetworkConnectionSet != nil {
			var networkConnectionSetList []interface{}
			for _, networkConnectionSet := range dataEngine.NetworkConnectionSet {
				networkConnectionSetMap := map[string]interface{}{}

				if networkConnectionSet.Id != nil {
					networkConnectionSetMap["id"] = networkConnectionSet.Id
				}

				if networkConnectionSet.AssociateId != nil {
					networkConnectionSetMap["associate_id"] = networkConnectionSet.AssociateId
				}

				if networkConnectionSet.HouseId != nil {
					networkConnectionSetMap["house_id"] = networkConnectionSet.HouseId
				}

				if networkConnectionSet.DatasourceConnectionId != nil {
					networkConnectionSetMap["datasource_connection_id"] = networkConnectionSet.DatasourceConnectionId
				}

				if networkConnectionSet.State != nil {
					networkConnectionSetMap["state"] = networkConnectionSet.State
				}

				if networkConnectionSet.CreateTime != nil {
					networkConnectionSetMap["create_time"] = networkConnectionSet.CreateTime
				}

				if networkConnectionSet.UpdateTime != nil {
					networkConnectionSetMap["update_time"] = networkConnectionSet.UpdateTime
				}

				if networkConnectionSet.Appid != nil {
					networkConnectionSetMap["appid"] = networkConnectionSet.Appid
				}

				if networkConnectionSet.HouseName != nil {
					networkConnectionSetMap["house_name"] = networkConnectionSet.HouseName
				}

				if networkConnectionSet.DatasourceConnectionName != nil {
					networkConnectionSetMap["datasource_connection_name"] = networkConnectionSet.DatasourceConnectionName
				}

				if networkConnectionSet.NetworkConnectionType != nil {
					networkConnectionSetMap["network_connection_type"] = networkConnectionSet.NetworkConnectionType
				}

				if networkConnectionSet.Uin != nil {
					networkConnectionSetMap["uin"] = networkConnectionSet.Uin
				}

				if networkConnectionSet.SubAccountUin != nil {
					networkConnectionSetMap["sub_account_uin"] = networkConnectionSet.SubAccountUin
				}

				if networkConnectionSet.NetworkConnectionDesc != nil {
					networkConnectionSetMap["network_connection_desc"] = networkConnectionSet.NetworkConnectionDesc
				}

				if networkConnectionSet.DatasourceConnectionVpcId != nil {
					networkConnectionSetMap["datasource_connection_vpc_id"] = networkConnectionSet.DatasourceConnectionVpcId
				}

				if networkConnectionSet.DatasourceConnectionSubnetId != nil {
					networkConnectionSetMap["datasource_connection_subnet_id"] = networkConnectionSet.DatasourceConnectionSubnetId
				}

				if networkConnectionSet.DatasourceConnectionCidrBlock != nil {
					networkConnectionSetMap["datasource_connection_cidr_block"] = networkConnectionSet.DatasourceConnectionCidrBlock
				}

				if networkConnectionSet.DatasourceConnectionSubnetCidrBlock != nil {
					networkConnectionSetMap["datasource_connection_subnet_cidr_block"] = networkConnectionSet.DatasourceConnectionSubnetCidrBlock
				}

				networkConnectionSetList = append(networkConnectionSetList, networkConnectionSetMap)
			}

			dataEngineInfoMap["network_connection_set"] = networkConnectionSetList
		}

		if dataEngine.UiURL != nil {
			dataEngineInfoMap["ui_u_r_l"] = dataEngine.UiURL
		}

		if dataEngine.ResourceType != nil {
			dataEngineInfoMap["resource_type"] = dataEngine.ResourceType
		}

		if dataEngine.ImageVersionId != nil {
			dataEngineInfoMap["image_version_id"] = dataEngine.ImageVersionId
		}

		if dataEngine.ChildImageVersionId != nil {
			dataEngineInfoMap["child_image_version_id"] = dataEngine.ChildImageVersionId
		}

		if dataEngine.ImageVersionName != nil {
			dataEngineInfoMap["image_version_name"] = dataEngine.ImageVersionName
		}

		if dataEngine.StartStandbyCluster != nil {
			dataEngineInfoMap["start_standby_cluster"] = dataEngine.StartStandbyCluster
		}

		if dataEngine.ElasticSwitch != nil {
			dataEngineInfoMap["elastic_switch"] = dataEngine.ElasticSwitch
		}

		if dataEngine.ElasticLimit != nil {
			dataEngineInfoMap["elastic_limit"] = dataEngine.ElasticLimit
		}

		if dataEngine.DefaultHouse != nil {
			dataEngineInfoMap["default_house"] = dataEngine.DefaultHouse
		}

		if dataEngine.MaxConcurrency != nil {
			dataEngineInfoMap["max_concurrency"] = dataEngine.MaxConcurrency
		}

		if dataEngine.TolerableQueueTime != nil {
			dataEngineInfoMap["tolerable_queue_time"] = dataEngine.TolerableQueueTime
		}

		if dataEngine.UserAppId != nil {
			dataEngineInfoMap["user_app_id"] = dataEngine.UserAppId
		}

		if dataEngine.UserUin != nil {
			dataEngineInfoMap["user_uin"] = dataEngine.UserUin
		}

		if dataEngine.SessionResourceTemplate != nil {
			sessionResourceTemplateMap := map[string]interface{}{}

			if dataEngine.SessionResourceTemplate.DriverSize != nil {
				sessionResourceTemplateMap["driver_size"] = dataEngine.SessionResourceTemplate.DriverSize
			}

			if dataEngine.SessionResourceTemplate.ExecutorSize != nil {
				sessionResourceTemplateMap["executor_size"] = dataEngine.SessionResourceTemplate.ExecutorSize
			}

			if dataEngine.SessionResourceTemplate.ExecutorNums != nil {
				sessionResourceTemplateMap["executor_nums"] = dataEngine.SessionResourceTemplate.ExecutorNums
			}

			if dataEngine.SessionResourceTemplate.ExecutorMaxNumbers != nil {
				sessionResourceTemplateMap["executor_max_numbers"] = dataEngine.SessionResourceTemplate.ExecutorMaxNumbers
			}

			dataEngineInfoMap["session_resource_template"] = []interface{}{sessionResourceTemplateMap}
		}

		_ = d.Set("data_engine", []interface{}{dataEngineInfoMap})
	}

	d.SetId(dataEngineName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), dataEngineInfoMap); e != nil {
			return e
		}
	}
	return nil
}
