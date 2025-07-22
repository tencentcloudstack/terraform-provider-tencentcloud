package emr

import (
	"context"
	innerErr "errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
)

func ResourceTencentCloudEmrCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrClusterCreate,
		Read:   resourceTencentCloudEmrClusterRead,
		Delete: resourceTencentCloudEmrClusterDelete,
		Update: resourceTencentCloudEmrClusterUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "It will be deprecated in later versions.",
				Description: "Display strategy of EMR instance.",
			},
			"product_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				Description: "Product ID. Different products ID represents different EMR product versions. Value range:\n" +
					"	- 16: represents EMR-V2.3.0\n" +
					"	- 20: represents EMR-V2.5.0\n" +
					"	- 25: represents EMR-V3.1.0\n" +
					"	- 27: represents KAFKA-V1.0.0\n" +
					"	- 30: represents EMR-V2.6.0\n" +
					"	- 33: represents EMR-V3.2.1\n" +
					"	- 34: represents EMR-V3.3.0\n" +
					"	- 37: represents EMR-V3.4.0\n" +
					"	- 38: represents EMR-V2.7.0\n" +
					"	- 44: represents EMR-V3.5.0\n" +
					"	- 50: represents KAFKA-V2.0.0\n" +
					"	- 51: represents STARROCKS-V1.4.0\n" +
					"	- 53: represents EMR-V3.6.0\n" +
					"	- 54: represents STARROCKS-V2.0.0.",
			},
			"vpc_settings": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "The private net config of EMR instance.",
			},
			"softwares": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The softwares of a EMR instance.",
			},
			"resource_spec": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"multi_zone_setting", "multi_zone"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_resource_spec": buildResourceSpecSchema(),
						"core_resource_spec":   buildResourceSpecSchema(),
						"task_resource_spec":   buildResourceSpecSchema(),
						"master_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of master node.",
						},
						"core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of core node.",
						},
						"task_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of core node.",
						},
						"common_resource_spec": buildResourceSpecSchema(),
						"common_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: "The number of common node.",
						},
					},
				},
				Description: "Resource specification of EMR instance.",
			},
			"terminate_node_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Terminate nodes. Note: it only works when the number of nodes decreases.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_instance_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Destroy resource list.",
						},
						"node_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value range of destruction node type: `MASTER`, `TASK`, `CORE`, `ROUTER`.",
						},
					},
				},
			},
			"support_ha": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
				Description:  "The flag whether the instance support high availability.(0=>not support, 1=>support).",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(6, 36),
				Description:  "Name of the instance, which can contain 6 to 36 English letters, Chinese characters, digits, dashes(-), or underscores(_).",
			},
			"pay_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
				Description:  "The pay mode of instance. 0 represent POSTPAID_BY_HOUR, 1 represent PREPAID.",
			},
			"placement": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"placement", "placement_info"},
				Deprecated:   "It will be deprecated in later versions. Use `placement_info` instead.",
				Description:  "The location of the instance.",
			},
			"placement_info": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"placement", "placement_info"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "Project id.",
						},
					},
				},
				Description: "The location of the instance.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The length of time the instance was purchased. Use with TimeUnit.When TimeUnit is s, the parameter can only be filled in at 3600, representing a metered instance.\nWhen TimeUnit is m, the number filled in by this parameter indicates the length of purchase of the monthly instance of the package year, such as 1 for one month of purchase.",
			},
			"time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unit of time in which the instance was purchased. When PayMode is 0, TimeUnit can only take values of s(second). When PayMode is 1, TimeUnit can only take the value m(month).",
			},
			"login_settings": {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
				Description: "Instance login settings. There are two optional fields:" +
					"- password: Instance login password: 8-16 characters, including uppercase letters, lowercase letters, numbers and special characters. Special symbols only support! @% ^ *. The first bit of the password cannot be a special character;" +
					"- public_key_id: Public key id. After the key is associated, the instance can be accessed through the corresponding private key.",
			},
			"extend_fs_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access the external file system.",
			},
			"scene_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Scene-based value:\n" +
					"	- Hadoop-Kudu\n" +
					"	- Hadoop-Zookeeper\n" +
					"	- Hadoop-Presto\n" +
					"	- Hadoop-Hbase.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created EMR instance id.",
			},
			"need_master_wan": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      EMR_MASTER_WAN_TYPE_NEED_MASTER_WAN,
				ValidateFunc: tccommon.ValidateAllowedStringValue(EMR_MASTER_WAN_TYPES),
				Description: `Whether to enable the cluster Master node public network. Value range:
				- NEED_MASTER_WAN: Indicates that the cluster Master node public network is enabled.
				- NOT_NEED_MASTER_WAN: Indicates that it is not turned on.
				By default, the cluster Master node internet is enabled.`,
			},
			"sg_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the security group to which the instance belongs, in the form of sg-xxxxxxxx.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tag description list.",
			},
			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "0 means turn off automatic renewal, 1 means turn on automatic renewal. Default is 0.",
			},
			"pre_executed_file_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Pre executed file settings. It can only be set at the time of creation, and cannot be modified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Execution script parameters.",
						},
						"run_order": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Run order.",
						},
						"when_run": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "`resourceAfter` or `clusterAfter`.",
						},
						"cos_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Script file name.",
						},
						"cos_file_uri": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The cos address of the script.",
						},
						"cos_secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Cos secretId.",
						},
						"cos_secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Cos secretKey.",
						},
						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Remark.",
						},
					},
				},
			},
			"multi_zone": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"multi_zone_setting"},
				Description:  "true means that cross-AZ deployment is enabled; it is only a user parameter when creating a new cluster, and no subsequent adjustment is supported.",
			},
			"multi_zone_setting": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"multi_zone"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_settings": {
							Type:        schema.TypeMap,
							Required:    true,
							ForceNew:    true,
							Description: "The private net config of EMR instance.",
						},
						"placement": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Zone.",
									},
								},
							},
							Description: "The location of the instance.",
						},
						"resource_spec": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_resource_spec": buildResourceSpecSchema(),
									"core_resource_spec":   buildResourceSpecSchema(),
									"task_resource_spec":   buildResourceSpecSchema(),
									"master_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										Description: "The number of master node.",
									},
									"core_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										Description: "The number of core node.",
									},
									"task_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										Description: "The number of core node.",
									},
									"common_resource_spec": buildResourceSpecSchema(),
									"common_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										ForceNew:    true,
										Description: "The number of common node.",
									},
								},
							},
							Description: "Resource specification of EMR instance.",
						},
					},
				},
				Description: "The specification of node resources is as follows: fill in a few available areas. In order, the first one is the main available area, the second one is the backup available area, and the third one is the arbitration available area.",
			},
		},
	}
}

func resourceTencentCloudEmrClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	immutableFields := []string{"auto_renew", "placement", "placement_info", "display_strategy", "login_settings", "extend_fs_field", "scene_name", "pay_mode"}
	for _, f := range immutableFields {
		if d.HasChange(f) {
			return fmt.Errorf("cannot update argument `%s`", f)
		}
	}

	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()
	timeUnit, hasTimeUnit := d.GetOkExists("time_unit")
	timeSpan, hasTimeSpan := d.GetOkExists("time_span")
	payMode, hasPayMode := d.GetOkExists("pay_mode")
	if !hasTimeUnit || !hasTimeSpan || !hasPayMode {
		return innerErr.New("Time_unit, time_span or pay_mode must be set.")
	}
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		err := emrService.ModifyResourcesTags(ctx, meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region, instanceId, oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		if err != nil {
			return err
		}
	}

	if d.HasChange("resource_spec.0.master_count") {
		request := emr.NewScaleOutInstanceRequest()
		request.TimeUnit = common.StringPtr(timeUnit.(string))
		request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
		request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
		request.InstanceId = common.StringPtr(instanceId)

		o, n := d.GetChange("resource_spec.0.master_count")
		if o.(int) < n.(int) {
			request.MasterCount = common.Uint64Ptr((uint64)(n.(int) - o.(int)))
			traceId, err := emrService.ScaleOutInstance(ctx, request)
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}
		}
	}
	if d.HasChange("resource_spec.0.task_count") {
		request := emr.NewScaleOutInstanceRequest()
		request.TimeUnit = common.StringPtr(timeUnit.(string))
		request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
		request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
		request.InstanceId = common.StringPtr(instanceId)

		o, n := d.GetChange("resource_spec.0.task_count")
		if o.(int) < n.(int) {
			request.TaskCount = common.Uint64Ptr((uint64)(n.(int) - o.(int)))
			traceId, err := emrService.ScaleOutInstance(ctx, request)
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}
		}
	}
	if d.HasChange("resource_spec.0.core_count") {
		request := emr.NewScaleOutInstanceRequest()
		request.TimeUnit = common.StringPtr(timeUnit.(string))
		request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
		request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
		request.InstanceId = common.StringPtr(instanceId)

		o, n := d.GetChange("resource_spec.0.core_count")
		if o.(int) < n.(int) {
			request.CoreCount = common.Uint64Ptr((uint64)(n.(int) - o.(int)))
			traceId, err := emrService.ScaleOutInstance(ctx, request)
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}
		}
	}

	// multi_zone_setting
	var isChangeMultiZoneNodeCount bool
	if d.HasChange("multi_zone_setting") {
		if v, ok := d.GetOk("multi_zone_setting"); ok {
			multiZoneSettings := v.([]interface{})
			for idx, multiZoneSetting := range multiZoneSettings {
				multiZoneSettingMap := multiZoneSetting.(map[string]interface{})
				placement, placementOk := multiZoneSettingMap["placement"]
				if !placementOk {
					return fmt.Errorf("Argument `multi_zone_setting.%d.placement` must be set", idx)
				}
				placementList := placement.([]interface{})
				if len(placementList) == 0 {
					return fmt.Errorf("Argument `multi_zone_setting.%d.placement` must be set", idx)
				}
				if _, ok := placementList[0].(map[string]interface{})["zone"]; !ok {
					return fmt.Errorf("Argument `multi_zone_setting.%d.placement.zone` must be set", idx)
				}
				zone := placementList[0].(map[string]interface{})["zone"].(string)
				var zoneId int64
				cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
				err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					zones, errRet := cvmService.DescribeZones(ctx)
					if errRet != nil {
						return tccommon.RetryError(errRet, tccommon.InternalError)
					}
					for _, z := range zones {
						if *z.Zone == zone {
							zoneId = helper.StrToInt64(*z.ZoneId)
						}
					}
					return nil
				})
				if err != nil {
					return err
				}
				if zoneId == 0 {
					return fmt.Errorf("Argument `multi_zone_setting.%d.placement.zone` %s not found", idx, zone)
				}
				vpcSetting, vpcSettingOk := multiZoneSettingMap["vpc_settings"]
				if !vpcSettingOk {
					return fmt.Errorf("Argument `multi_zone_setting.%d.vpc_settings` must be set", idx)
				}
				vpcSettingMap, vpcSettingMapOk := vpcSetting.(map[string]interface{})
				if !vpcSettingMapOk {
					return fmt.Errorf("Argument `multi_zone_setting.%d.vpc_settings` must be a map", idx)
				}
				subnetId := vpcSettingMap["subnet_id"].(string)

				resourceSpec, resourceSpecOk := multiZoneSettingMap["resource_spec"]
				if !resourceSpecOk {
					return fmt.Errorf("Argument `multi_zone_setting.%d.resource_spec` must be set", idx)
				}
				resourceSpecList := resourceSpec.([]interface{})
				if len(resourceSpecList) == 0 {
					return fmt.Errorf("Argument `multi_zone_setting.%d.resource_spec` must be set", idx)
				}
				if d.HasChange(fmt.Sprintf("multi_zone_setting.%d.resource_spec.0.master_count", idx)) {
					request := emr.NewScaleOutInstanceRequest()
					request.TimeUnit = common.StringPtr(timeUnit.(string))
					request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
					request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
					request.InstanceId = common.StringPtr(instanceId)
					request.ZoneId = helper.Int64(zoneId)
					request.SubnetId = common.StringPtr(subnetId)

					o, n := d.GetChange(fmt.Sprintf("multi_zone_setting.%d.resource_spec.0.master_count", idx))
					if o.(int) < n.(int) {
						request.MasterCount = common.Uint64Ptr((uint64)(n.(int) - o.(int)))
						traceId, err := emrService.ScaleOutInstance(ctx, request)
						if err != nil {
							return err
						}
						time.Sleep(5 * time.Second)
						conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if _, e := conf.WaitForState(); e != nil {
							return e
						}
					}
					isChangeMultiZoneNodeCount = true
				}
				if d.HasChange(fmt.Sprintf("multi_zone_setting.%d.resource_spec.0.task_count", idx)) {
					request := emr.NewScaleOutInstanceRequest()
					request.TimeUnit = common.StringPtr(timeUnit.(string))
					request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
					request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
					request.InstanceId = common.StringPtr(instanceId)
					request.ZoneId = helper.Int64(zoneId)
					request.SubnetId = common.StringPtr(subnetId)

					o, n := d.GetChange(fmt.Sprintf("multi_zone_setting.%d.resource_spec.0.task_count", idx))
					if o.(int) < n.(int) {
						request.TaskCount = common.Uint64Ptr((uint64)(n.(int) - o.(int)))
						traceId, err := emrService.ScaleOutInstance(ctx, request)
						if err != nil {
							return err
						}
						time.Sleep(5 * time.Second)
						conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if _, e := conf.WaitForState(); e != nil {
							return e
						}
					}
					isChangeMultiZoneNodeCount = true
				}
				if d.HasChange(fmt.Sprintf("multi_zone_setting.%d.resource_spec.0.core_count", idx)) {
					request := emr.NewScaleOutInstanceRequest()
					request.TimeUnit = common.StringPtr(timeUnit.(string))
					request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
					request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
					request.InstanceId = common.StringPtr(instanceId)
					request.ZoneId = helper.Int64(zoneId)
					request.SubnetId = common.StringPtr(subnetId)

					o, n := d.GetChange(fmt.Sprintf("multi_zone_setting.%d.resource_spec.0.core_count", idx))
					if o.(int) < n.(int) {
						request.CoreCount = common.Uint64Ptr((uint64)(n.(int) - o.(int)))
						traceId, err := emrService.ScaleOutInstance(ctx, request)
						if err != nil {
							return err
						}
						time.Sleep(5 * time.Second)
						conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if _, e := conf.WaitForState(); e != nil {
							return e
						}
					}
					isChangeMultiZoneNodeCount = true
				}

			}
		}
	}
	if d.HasChange("resource_spec.0.master_count") || d.HasChange("resource_spec.0.task_count") || d.HasChange("resource_spec.0.core_count") || isChangeMultiZoneNodeCount {
		if v, ok := d.GetOk("terminate_node_info"); ok {
			terminateNodeInfos := v.([]interface{})
			for _, terminateNodeInfo := range terminateNodeInfos {
				terminateNodeInfoMap := terminateNodeInfo.(map[string]interface{})
				instanceIds := make([]string, 0)
				for _, instanceId := range terminateNodeInfoMap["cvm_instance_ids"].([]interface{}) {
					instanceIds = append(instanceIds, instanceId.(string))
				}
				flowId, err := emrService.TerminateClusterNodes(ctx, instanceIds, instanceId, terminateNodeInfoMap["node_flag"].(string))
				if err != nil {
					return err
				}
				time.Sleep(5 * time.Second)
				conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
				if _, e := conf.WaitForState(); e != nil {
					return e
				}
			}
		}
	}

	return resourceTencentCloudEmrClusterRead(d, meta)
}

func resourceTencentCloudEmrClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	immutableFields := []string{"time_unit", "time_span", "login_settings"}
	for _, f := range immutableFields {
		if _, ok := d.GetOkExists(f); !ok {
			return fmt.Errorf("Argument `%s` must be set", f)
		}
	}

	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId, err := emrService.CreateInstance(ctx, d)
	if err != nil {
		return err
	}
	d.SetId(instanceId)
	_ = d.Set("instance_id", instanceId)
	var displayStrategy string
	if v, ok := d.GetOk("display_strategy"); ok {
		displayStrategy = v.(string)
	} else {
		displayStrategy = "clusterList"
	}
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, displayStrategy)

		if err != nil {
			return resource.RetryableError(err)
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusCreated {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudEmrClusterRead(d, meta)
}

func resourceTencentCloudEmrClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()
	clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)
	if len(clusters) == 0 {
		return innerErr.New("Not find clusters.")
	}
	metaDB := clusters[0].MetaDb
	if err != nil {
		return err
	}
	if err = emrService.DeleteInstance(ctx, d); err != nil {
		return err
	}
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "ResourceNotFound.InstanceNotFound" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusDeleted {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 预付费删除
	payMode := d.Get("pay_mode").(int)
	if payMode == 1 {
		if err = emrService.DeleteInstance(ctx, d); err != nil {
			return err
		}
		err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}

			if len(clusters) > 0 {
				return resource.RetryableError(fmt.Errorf("%v being destroyed", instanceId))
			}

			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if metaDB != nil && *metaDB != "" {
		// remove metadb
		mysqlService := svccdb.NewMysqlService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := mysqlService.OfflineIsolatedInstances(ctx, *metaDB)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudEmrClusterRead(d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()
	var instance *emr.ClusterInstancesInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}

		if len(result) > 0 {
			instance = result[0]
		}

		return nil
	})

	if err != nil {
		return err
	}

	_ = d.Set("instance_id", instanceId)
	clusterNodeMap := make(map[string]*emr.NodeHardwareInfo)
	clusterNodeNum := make(map[string]int)
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := emrService.DescribeClusterNodes(ctx, instanceId, "all", "all", 0, 10)

		if err != nil {
			return resource.RetryableError(err)
		}

		if len(result) > 0 {
			_ = d.Set("auto_renew", result[0].IsAutoRenew)
			for _, item := range result {
				node := item
				var nodeFlag string
				if node.ZoneId != nil {
					nodeFlag = fmt.Sprintf("%d-", *node.ZoneId)
				}
				// 节点类型 0:common节点；1:master节点；2:core节点；3:task节点
				if node.Flag != nil {
					nodeFlag = nodeFlag + strconv.FormatInt(*node.Flag, 10)
					clusterNodeMap[nodeFlag] = node
					if v, ok := clusterNodeNum[nodeFlag]; ok {
						clusterNodeNum[nodeFlag] = v + 1
					} else {
						clusterNodeNum[nodeFlag] = 1
					}

				}

			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	if instance != nil {
		_ = d.Set("scene_name", instance.SceneName)
		_ = d.Set("product_id", instance.ProductId)
		_ = d.Set("vpc_settings", map[string]interface{}{
			"vpc_id":    *instance.UniqVpcId,
			"subnet_id": *instance.UniqSubnetId,
		})
		if instance.Config != nil {
			if instance.Config.SoftInfo != nil {
				_ = d.Set("softwares", helper.PStrings(instance.Config.SoftInfo))
			}

			if instance.Config.SupportHA != nil {
				if *instance.Config.SupportHA {
					_ = d.Set("support_ha", 1)
				} else {
					_ = d.Set("support_ha", 0)
				}
			}

			if instance.Config.SecurityGroup != nil {
				_ = d.Set("sg_id", instance.Config.SecurityGroup)
			}

			multiZoneSetting := buildMultiZoneSettingList(instance, clusterNodeMap, clusterNodeNum)
			_ = d.Set("multi_zone", len(multiZoneSetting) > 1)

			if len(multiZoneSetting) > 1 {
				multiZoneSetting := buildMultiZoneSettingList(instance, clusterNodeMap, clusterNodeNum)
				_ = d.Set("multi_zone_setting", multiZoneSetting)
			} else {
				resourceSpec := buildResourceSpec(instance, clusterNodeMap, clusterNodeNum)
				_ = d.Set("resource_spec", []interface{}{resourceSpec})
			}
		}

		_ = d.Set("instance_name", instance.ClusterName)
		_ = d.Set("pay_mode", instance.ChargeType)
		placement := map[string]interface{}{
			"zone":       *instance.Zone,
			"project_id": *instance.ProjectId,
		}
		_ = d.Set("placement", map[string]interface{}{
			"zone": *instance.Zone,
		})
		_ = d.Set("placement_info", []interface{}{placement})
		if instance.MasterIp != nil && len(*instance.MasterIp) > 2 {
			_ = d.Set("need_master_wan", "NEED_MASTER_WAN")
		} else {
			_ = d.Set("need_master_wan", "NOT_NEED_MASTER_WAN")
		}
	}

	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "emr", "emr-instance", region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func buildResourceSpec(instance *emr.ClusterInstancesInfo, clusterNodeMap map[string]*emr.NodeHardwareInfo, clusterNodeNum map[string]int) map[string]interface{} {
	resourceSpec := make(map[string]interface{})

	if v, ok := clusterNodeNum[fmt.Sprintf("%d-1", *instance.ZoneId)]; ok {
		resourceSpec["master_count"] = v
	}

	if v, ok := clusterNodeMap[fmt.Sprintf("%d-1", *instance.ZoneId)]; ok && v != nil {
		masterResourceSpec := make(map[string]interface{})
		masterResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
		masterResourceSpec["cpu"] = v.CpuNum
		if instance.Config.MasterResource != nil {
			masterResource := instance.Config.MasterResource
			masterResourceSpec["disk_size"] = masterResource.DiskSize
			masterResourceSpec["multi_disks"] = fetchMultiDisks(v, masterResource)

		}
		if v.StorageType != nil {
			masterResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
		}
		masterResourceSpec["spec"] = v.Spec
		masterResourceSpec["storage_type"] = v.RootStorageType
		masterResourceSpec["root_size"] = v.RootSize
		resourceSpec["master_resource_spec"] = []interface{}{masterResourceSpec}
	}

	if v, ok := clusterNodeNum[fmt.Sprintf("%d-2", *instance.ZoneId)]; ok {
		resourceSpec["core_count"] = v
	}
	if v, ok := clusterNodeMap[fmt.Sprintf("%d-2", *instance.ZoneId)]; ok && v != nil {
		coreResourceSpec := make(map[string]interface{})
		coreResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
		coreResourceSpec["cpu"] = v.CpuNum
		if instance.Config.CoreResource != nil {
			coreResource := instance.Config.CoreResource
			coreResourceSpec["disk_size"] = coreResource.DiskSize
			coreResourceSpec["multi_disks"] = fetchMultiDisks(v, coreResource)
		}
		if v.StorageType != nil {
			coreResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
		}
		coreResourceSpec["spec"] = v.Spec
		coreResourceSpec["storage_type"] = v.RootStorageType
		coreResourceSpec["root_size"] = v.RootSize
		resourceSpec["core_resource_spec"] = []interface{}{coreResourceSpec}
	}

	if v, ok := clusterNodeNum[fmt.Sprintf("%d-3", *instance.ZoneId)]; ok {
		resourceSpec["task_count"] = v
	}
	if v, ok := clusterNodeMap[fmt.Sprintf("%d-3", *instance.ZoneId)]; ok && v != nil {
		taskResourceSpec := make(map[string]interface{})
		taskResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
		taskResourceSpec["cpu"] = v.CpuNum
		if instance.Config.TaskResource != nil {
			taskResource := instance.Config.TaskResource
			taskResourceSpec["disk_size"] = taskResource.DiskSize
			taskResourceSpec["multi_disks"] = fetchMultiDisks(v, taskResource)
		}
		if v.StorageType != nil {
			taskResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
		}
		taskResourceSpec["spec"] = v.Spec
		taskResourceSpec["storage_type"] = v.RootStorageType
		taskResourceSpec["root_size"] = v.RootSize
		resourceSpec["task_resource_spec"] = []interface{}{taskResourceSpec}
	}

	if v, ok := clusterNodeNum[fmt.Sprintf("%d-0", *instance.ZoneId)]; ok {
		resourceSpec["common_count"] = v
	}
	if v, ok := clusterNodeMap[fmt.Sprintf("%d-0", *instance.ZoneId)]; ok && v != nil {
		comResourceSpec := make(map[string]interface{})
		comResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
		comResourceSpec["cpu"] = v.CpuNum
		if instance.Config.ComResource != nil {
			comResource := instance.Config.ComResource
			comResourceSpec["disk_size"] = comResource.DiskSize
			comResourceSpec["multi_disks"] = fetchMultiDisks(v, comResource)
		}
		if v.StorageType != nil {
			comResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
		}
		comResourceSpec["spec"] = v.Spec
		comResourceSpec["storage_type"] = v.RootStorageType
		comResourceSpec["root_size"] = v.RootSize
		resourceSpec["common_resource_spec"] = []interface{}{comResourceSpec}
	}
	return resourceSpec
}

func buildMultiZoneSettingList(instance *emr.ClusterInstancesInfo, clusterNodeMap map[string]*emr.NodeHardwareInfo, clusterNodeNum map[string]int) []interface{} {
	firstZoneID := helper.Int64ToStr(*instance.ZoneId)
	var secondZoneID, thirdZoneID string
	zoneNodeNums := make(map[string]int)
	for k, v := range clusterNodeNum {
		zoneID := strings.Split(k, "-")[0]
		zoneNodeNums[zoneID] = zoneNodeNums[zoneID] + v
	}
	for k, v := range zoneNodeNums {
		if v == 1 {
			zoneID := strings.Split(k, "-")[0]
			thirdZoneID = zoneID
		}
	}
	for k := range zoneNodeNums {
		if thirdZoneID == "" {
			if !strings.HasPrefix(k, firstZoneID) {
				zoneID := strings.Split(k, "-")[0]
				secondZoneID = zoneID
				break
			}
		} else {
			if !strings.HasPrefix(k, firstZoneID) && !strings.HasPrefix(k, thirdZoneID) {
				zoneID := strings.Split(k, "-")[0]
				secondZoneID = zoneID
				break
			}
		}
	}

	multiZoneSettingList := []interface{}{}
	for _, z := range []string{firstZoneID, secondZoneID, thirdZoneID} {
		resourceSpec := make(map[string]interface{})
		vpcSettings := make(map[string]interface{})
		placement := make(map[string]interface{})
		vpcSettings["vpc_id"] = instance.UniqVpcId

		if v, ok := clusterNodeNum[fmt.Sprintf("%s-1", z)]; ok {
			resourceSpec["master_count"] = v
		}

		if v, ok := clusterNodeMap[fmt.Sprintf("%s-1", z)]; ok && v != nil {
			vpcSettings["subnet_id"] = v.SubnetInfo.SubnetId
			placement["zone"] = v.Zone
			masterResourceSpec := make(map[string]interface{})
			masterResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
			masterResourceSpec["cpu"] = v.CpuNum
			if instance.Config.MasterResource != nil {
				masterResource := instance.Config.MasterResource
				masterResourceSpec["disk_size"] = masterResource.DiskSize
				masterResourceSpec["multi_disks"] = fetchMultiDisks(v, masterResource)

			}
			if v.StorageType != nil {
				masterResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
			}
			masterResourceSpec["spec"] = v.Spec
			masterResourceSpec["storage_type"] = v.RootStorageType
			masterResourceSpec["root_size"] = v.RootSize
			resourceSpec["master_resource_spec"] = []interface{}{masterResourceSpec}
		}

		if v, ok := clusterNodeNum[fmt.Sprintf("%s-2", z)]; ok {
			resourceSpec["core_count"] = v
		}
		if v, ok := clusterNodeMap[fmt.Sprintf("%s-2", z)]; ok && v != nil {
			vpcSettings["subnet_id"] = v.SubnetInfo.SubnetId
			placement["zone"] = v.Zone
			coreResourceSpec := make(map[string]interface{})
			coreResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
			coreResourceSpec["cpu"] = v.CpuNum
			if instance.Config.CoreResource != nil {
				coreResource := instance.Config.CoreResource
				coreResourceSpec["disk_size"] = coreResource.DiskSize
				coreResourceSpec["multi_disks"] = fetchMultiDisks(v, coreResource)
			}
			if v.StorageType != nil {
				coreResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
			}
			coreResourceSpec["spec"] = v.Spec
			coreResourceSpec["storage_type"] = v.RootStorageType
			coreResourceSpec["root_size"] = v.RootSize
			resourceSpec["core_resource_spec"] = []interface{}{coreResourceSpec}
		}

		if v, ok := clusterNodeNum[fmt.Sprintf("%s-3", z)]; ok {
			resourceSpec["task_count"] = v
		}
		if v, ok := clusterNodeMap[fmt.Sprintf("%s-3", z)]; ok && v != nil {
			vpcSettings["subnet_id"] = v.SubnetInfo.SubnetId
			placement["zone"] = v.Zone
			taskResourceSpec := make(map[string]interface{})
			taskResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
			taskResourceSpec["cpu"] = v.CpuNum
			if instance.Config.TaskResource != nil {
				taskResource := instance.Config.TaskResource
				taskResourceSpec["disk_size"] = taskResource.DiskSize
				taskResourceSpec["multi_disks"] = fetchMultiDisks(v, taskResource)
			}
			if v.StorageType != nil {
				taskResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
			}
			taskResourceSpec["spec"] = v.Spec
			taskResourceSpec["storage_type"] = v.RootStorageType
			taskResourceSpec["root_size"] = v.RootSize
			resourceSpec["task_resource_spec"] = []interface{}{taskResourceSpec}
		}

		if v, ok := clusterNodeNum[fmt.Sprintf("%s-0", z)]; ok {
			resourceSpec["common_count"] = v
		}
		if v, ok := clusterNodeMap[fmt.Sprintf("%s-0", z)]; ok && v != nil {
			vpcSettings["subnet_id"] = v.SubnetInfo.SubnetId
			placement["zone"] = v.Zone
			comResourceSpec := make(map[string]interface{})
			comResourceSpec["mem_size"] = int(*v.MemSize / 1024 / 1024)
			comResourceSpec["cpu"] = v.CpuNum
			if instance.Config.ComResource != nil {
				comResource := instance.Config.ComResource
				comResourceSpec["disk_size"] = comResource.DiskSize
				comResourceSpec["multi_disks"] = fetchMultiDisks(v, comResource)
			}
			if v.StorageType != nil {
				comResourceSpec["disk_type"] = translateDiskType(*v.StorageType)
			}
			comResourceSpec["spec"] = v.Spec
			comResourceSpec["storage_type"] = v.RootStorageType
			comResourceSpec["root_size"] = v.RootSize
			resourceSpec["common_resource_spec"] = []interface{}{comResourceSpec}
		}
		if len(resourceSpec) == 0 {
			continue
		}
		multiZoneSettingList = append(multiZoneSettingList, map[string]interface{}{
			"placement":     []interface{}{placement},
			"vpc_settings":  vpcSettings,
			"resource_spec": []interface{}{resourceSpec},
		})
	}

	return multiZoneSettingList
}
