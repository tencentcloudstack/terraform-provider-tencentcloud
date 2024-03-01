package csip

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	csip "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/csip/v20221121"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCsipRiskCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCsipRiskCenterCreate,
		Read:   resourceTencentCloudCsipRiskCenterRead,
		Update: resourceTencentCloudCsipRiskCenterUpdate,
		Delete: resourceTencentCloudCsipRiskCenterDelete,

		Schema: map[string]*schema.Schema{
			"task_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task Name.",
			},
			"scan_plan_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(SCAN_PLAN_TYPE),
				Description:  "0- Periodic task, 1- immediate scan, 2- periodic scan, 3- Custom; 0, 2 and 3 are required for scan_plan_content.",
			},
			"scan_asset_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(SCAN_ASSET_TYPE),
				Description:  "0- Full scan, 1- Specify asset scan, 2- Exclude asset scan, 3- Manually fill in the scan. If 1 and 2 are required while task_mode not 1, the Assets field is required. If 3 is required, SelfDefiningAssets is required.",
			},
			"scan_item": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Scan Project. Example: port/poc/weakpass/webcontent/configrisk/exposedserver.",
			},
			"assets": {
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Scan the asset information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Asset nameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Asset typeNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"asset_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Asset classificationNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"asset": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ip/ domain name/asset id, database id, etc.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RegionNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"arn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Multi-cloud asset unique idNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"scan_plan_content": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Scan plan details.",
			},
			"self_defining_assets": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ip/domain/url array.",
			},
			"scan_from": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Request origin.",
			},
			"task_advance_cfg": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Advanced configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_risk": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Advanced Port Risk Configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port_sets": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Port collection, separated by commas.",
									},
									"check_type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Detection item type, 0-system defined, 1-user-defined.",
									},
									"detail": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Description of detection items.",
									},
									"enable": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Whether to enable, 0- No, 1- Enable.",
									},
								},
							},
						},
						"vul_risk": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Advanced vulnerability risk configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"risk_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Risk ID.",
									},
									"enable": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Whether to enable, 0- No, 1- Enable.",
									},
								},
							},
						},
						"weak_pwd_risk": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Weak password risk advanced configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"check_item_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Detection item ID.",
									},
									"enable": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Whether to enable, 0- No, 1- Enable.",
									},
								},
							},
						},
						"cfg_risk": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Configure advanced risk Settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"item_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Detection item ID.",
									},
									"enable": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Whether to enable, 0- No, 1- Enable.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource type.",
									},
								},
							},
						},
					},
				},
			},
			"task_mode": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(SCAN_TASK_MODE),
				Default:      SCAN_TASK_MODE_0,
				Description:  "Physical examination mode, 0-standard mode, 1-fast mode, 2-advanced mode, default standard mode.",
			},
		},
	}
}

func resourceTencentCloudCsipRiskCenterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_csip_risk_center.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = csip.NewCreateRiskCenterScanTaskRequest()
		response      = csip.NewCreateRiskCenterScanTaskResponse()
		taskId        string
		scanPlanType  int
		scanAssetType int
		taskMode      int
	)

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("scan_plan_type"); ok {
		request.ScanPlanType = helper.IntInt64(v.(int))
		scanPlanType = v.(int)
	}

	if v, ok := d.GetOkExists("scan_asset_type"); ok {
		request.ScanAssetType = helper.IntInt64(v.(int))
		scanAssetType = v.(int)
	}

	if v, ok := d.GetOkExists("task_mode"); ok {
		request.TaskMode = helper.IntInt64(v.(int))
		taskMode = v.(int)
		if taskMode == SCAN_TASK_MODE_1 {
			if scanPlanType != SCAN_PLAN_TYPE_1 || scanAssetType != SCAN_ASSET_TYPE_1 {
				return fmt.Errorf("If task_mode is `1`, scan_plan_type and scan_asset_type should be set to `1`.")
			}
		}
	}

	if v, ok := d.GetOk("scan_item"); ok {
		scanItemSet := v.(*schema.Set).List()
		for i := range scanItemSet {
			if scanItemSet[i] != nil {
				scanItem := scanItemSet[i].(string)
				request.ScanItem = append(request.ScanItem, &scanItem)
			}
		}
	}

	if v, ok := d.GetOk("assets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			taskAssetObject := csip.TaskAssetObject{}
			if v, ok := dMap["asset_name"]; ok {
				taskAssetObject.AssetName = helper.String(v.(string))
			}

			if v, ok := dMap["instance_type"]; ok {
				taskAssetObject.InstanceType = helper.String(v.(string))
			}

			if v, ok := dMap["asset_type"]; ok {
				taskAssetObject.AssetType = helper.String(v.(string))
			}

			if v, ok := dMap["asset"]; ok {
				taskAssetObject.Asset = helper.String(v.(string))
			}

			if v, ok := dMap["region"]; ok {
				taskAssetObject.Region = helper.String(v.(string))
			}

			if v, ok := dMap["arn"]; ok {
				taskAssetObject.Arn = helper.String(v.(string))
			}

			request.Assets = append(request.Assets, &taskAssetObject)
		}
	}

	if v, ok := d.GetOk("scan_plan_content"); ok {
		request.ScanPlanContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("self_defining_assets"); ok {
		selfDefiningAssetsSet := v.(*schema.Set).List()
		for i := range selfDefiningAssetsSet {
			if selfDefiningAssetsSet[i] != nil {
				selfDefiningAssets := selfDefiningAssetsSet[i].(string)
				request.SelfDefiningAssets = append(request.SelfDefiningAssets, &selfDefiningAssets)
			}
		}
	}

	request.ScanFrom = helper.String("csip")

	if dMap, ok := helper.InterfacesHeadMap(d, "task_advance_cfg"); ok {
		taskAdvanceCFG := csip.TaskAdvanceCFG{}
		if v, ok := dMap["port_risk"]; ok {
			for _, item := range v.([]interface{}) {
				portRiskMap := item.(map[string]interface{})
				portRiskAdvanceCFGParamItem := csip.PortRiskAdvanceCFGParamItem{}
				if v, ok := portRiskMap["port_sets"]; ok {
					portRiskAdvanceCFGParamItem.PortSets = helper.String(v.(string))
				}

				if v, ok := portRiskMap["check_type"]; ok {
					portRiskAdvanceCFGParamItem.CheckType = helper.IntInt64(v.(int))
				}

				if v, ok := portRiskMap["detail"]; ok {
					portRiskAdvanceCFGParamItem.Detail = helper.String(v.(string))
				}

				if v, ok := portRiskMap["enable"]; ok {
					portRiskAdvanceCFGParamItem.Enable = helper.IntInt64(v.(int))
				}
				taskAdvanceCFG.PortRisk = append(taskAdvanceCFG.PortRisk, &portRiskAdvanceCFGParamItem)
			}
		}

		if v, ok := dMap["vul_risk"]; ok {
			for _, item := range v.([]interface{}) {
				vulRiskMap := item.(map[string]interface{})
				taskCenterVulRiskInputParam := csip.TaskCenterVulRiskInputParam{}
				if v, ok := vulRiskMap["risk_id"]; ok {
					taskCenterVulRiskInputParam.RiskId = helper.String(v.(string))
				}

				if v, ok := vulRiskMap["enable"]; ok {
					taskCenterVulRiskInputParam.Enable = helper.IntInt64(v.(int))
				}
				taskAdvanceCFG.VulRisk = append(taskAdvanceCFG.VulRisk, &taskCenterVulRiskInputParam)
			}
		}

		if v, ok := dMap["weak_pwd_risk"]; ok {
			for _, item := range v.([]interface{}) {
				weakPwdRiskMap := item.(map[string]interface{})
				taskCenterWeakPwdRiskInputParam := csip.TaskCenterWeakPwdRiskInputParam{}
				if v, ok := weakPwdRiskMap["check_item_id"]; ok {
					taskCenterWeakPwdRiskInputParam.CheckItemId = helper.IntInt64(v.(int))
				}

				if v, ok := weakPwdRiskMap["enable"]; ok {
					taskCenterWeakPwdRiskInputParam.Enable = helper.IntInt64(v.(int))
				}

				taskAdvanceCFG.WeakPwdRisk = append(taskAdvanceCFG.WeakPwdRisk, &taskCenterWeakPwdRiskInputParam)
			}
		}
		if v, ok := dMap["cfg_risk"]; ok {
			for _, item := range v.([]interface{}) {
				cFGRiskMap := item.(map[string]interface{})
				taskCenterCFGRiskInputParam := csip.TaskCenterCFGRiskInputParam{}
				if v, ok := cFGRiskMap["item_id"]; ok {
					taskCenterCFGRiskInputParam.ItemId = helper.String(v.(string))
				}

				if v, ok := cFGRiskMap["enable"]; ok {
					taskCenterCFGRiskInputParam.Enable = helper.IntInt64(v.(int))
				}

				if v, ok := cFGRiskMap["resource_type"]; ok {
					taskCenterCFGRiskInputParam.ResourceType = helper.String(v.(string))
				}

				taskAdvanceCFG.CFGRisk = append(taskAdvanceCFG.CFGRisk, &taskCenterCFGRiskInputParam)
			}
		}

		request.TaskAdvanceCFG = &taskAdvanceCFG
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCsipClient().CreateRiskCenterScanTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.TaskId == nil {
			e = fmt.Errorf("csip riskCenter not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create csip riskCenter failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	// wait
	//waitRequest := csip.NewDescribeScanTaskListRequest()
	//waitRequest.Filter = &csip.Filter{
	//	Filters: []*csip.WhereFilter{
	//		{
	//			Name:         common.StringPtr("TaskId"),
	//			Values:       common.StringPtrs([]string{taskId}),
	//			OperatorType: common.Int64Ptr(1),
	//		},
	//	},
	//}
	//
	//err = resource.Retry(tccommon.ReadRetryTimeout*120, func() *resource.RetryError {
	//	result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCsipClient().DescribeScanTaskList(waitRequest)
	//	if e != nil {
	//		return tccommon.RetryError(e)
	//	} else {
	//		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	//	}
	//
	//	if result == nil || len(result.Response.Data) != 1 {
	//		e = fmt.Errorf("csip riskCenter not exists")
	//		return resource.NonRetryableError(e)
	//	}
	//
	//	if *result.Response.Data[0].Percent == 100 {
	//		return nil
	//	} else {
	//		return resource.RetryableError(fmt.Errorf("creating csip riskCenter, Percent %f ", *result.Response.Data[0].Percent))
	//	}
	//})
	//
	//if err != nil {
	//	log.Printf("[CRITAL]%s create csip riskCenter failed, reason:%+v", logId, err)
	//	return err
	//}

	return resourceTencentCloudCsipRiskCenterRead(d, meta)
}

func resourceTencentCloudCsipRiskCenterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_csip_risk_center.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CsipService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  = d.Id()
	)

	riskCenter, err := service.DescribeCsipRiskCenterById(ctx, taskId)
	if err != nil {
		return err
	}

	if riskCenter == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CsipRiskCenter` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if riskCenter.TaskName != nil {
		_ = d.Set("task_name", riskCenter.TaskName)
	}

	if riskCenter.TaskType != nil {
		_ = d.Set("scan_plan_type", riskCenter.TaskType)
	}

	if riskCenter.ScanAssetType != nil {
		_ = d.Set("scan_asset_type", riskCenter.ScanAssetType)
	}

	if riskCenter.ScanItem != nil {
		tmpList := strings.Split(*riskCenter.ScanItem, tccommon.COMMA_SP)
		_ = d.Set("scan_item", tmpList)
	}

	if riskCenter.Assets != nil {
		assetsList := []interface{}{}
		for _, assets := range riskCenter.Assets {
			assetsMap := map[string]interface{}{}

			if assets.AssetName != nil {
				assetsMap["asset_name"] = assets.AssetName
			}

			if assets.InstanceType != nil {
				assetsMap["instance_type"] = assets.InstanceType
			}

			if assets.AssetType != nil {
				assetsMap["asset_type"] = assets.AssetType
			}

			if assets.Asset != nil {
				assetsMap["asset"] = assets.Asset
			}

			if assets.Region != nil {
				assetsMap["region"] = assets.Region
			}

			if assets.Arn != nil {
				assetsMap["arn"] = assets.Arn
			}

			assetsList = append(assetsList, assetsMap)
		}

		_ = d.Set("assets", assetsList)

	}

	if riskCenter.ScanPlanContent != nil {
		_ = d.Set("scan_plan_content", riskCenter.ScanPlanContent)
	}

	if riskCenter.SelfDefiningAssets != nil {
		_ = d.Set("self_defining_assets", riskCenter.SelfDefiningAssets)
	}

	if riskCenter.ScanFrom != nil {
		_ = d.Set("scan_from", riskCenter.ScanFrom)
	}

	if riskCenter.TaskMode != nil {
		_ = d.Set("task_mode", riskCenter.TaskMode)
	}

	return nil
}

func resourceTencentCloudCsipRiskCenterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_csip_risk_center.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = csip.NewModifyRiskCenterScanTaskRequest()
		taskId  = d.Id()
	)

	request.TaskId = &taskId

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("scan_plan_type"); ok {
		request.ScanPlanType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("scan_asset_type"); ok {
		request.ScanAssetType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("scan_item"); ok {
		scanItemSet := v.(*schema.Set).List()
		for i := range scanItemSet {
			if scanItemSet[i] != nil {
				scanItem := scanItemSet[i].(string)
				request.ScanItem = append(request.ScanItem, &scanItem)
			}
		}
	}

	if v, ok := d.GetOk("assets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			taskAssetObject := csip.TaskAssetObject{}
			if v, ok := dMap["asset_name"]; ok {
				taskAssetObject.AssetName = helper.String(v.(string))
			}
			if v, ok := dMap["instance_type"]; ok {
				taskAssetObject.InstanceType = helper.String(v.(string))
			}
			if v, ok := dMap["asset_type"]; ok {
				taskAssetObject.AssetType = helper.String(v.(string))
			}
			if v, ok := dMap["asset"]; ok {
				taskAssetObject.Asset = helper.String(v.(string))
			}
			if v, ok := dMap["region"]; ok {
				taskAssetObject.Region = helper.String(v.(string))
			}
			if v, ok := dMap["arn"]; ok {
				taskAssetObject.Arn = helper.String(v.(string))
			}
			request.Assets = append(request.Assets, &taskAssetObject)
		}
	}

	if v, ok := d.GetOk("scan_plan_content"); ok {
		request.ScanPlanContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("self_defining_assets"); ok {
		selfDefiningAssetsSet := v.(*schema.Set).List()
		for i := range selfDefiningAssetsSet {
			if selfDefiningAssetsSet[i] != nil {
				selfDefiningAssets := selfDefiningAssetsSet[i].(string)
				request.SelfDefiningAssets = append(request.SelfDefiningAssets, &selfDefiningAssets)
			}
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_advance_cfg"); ok {
		taskAdvanceCFG := csip.TaskAdvanceCFG{}
		if v, ok := dMap["port_risk"]; ok {
			for _, item := range v.([]interface{}) {
				portRiskMap := item.(map[string]interface{})
				portRiskAdvanceCFGParamItem := csip.PortRiskAdvanceCFGParamItem{}
				if v, ok := portRiskMap["port_sets"]; ok {
					portRiskAdvanceCFGParamItem.PortSets = helper.String(v.(string))
				}

				if v, ok := portRiskMap["check_type"]; ok {
					portRiskAdvanceCFGParamItem.CheckType = helper.IntInt64(v.(int))
				}

				if v, ok := portRiskMap["detail"]; ok {
					portRiskAdvanceCFGParamItem.Detail = helper.String(v.(string))
				}

				if v, ok := portRiskMap["enable"]; ok {
					portRiskAdvanceCFGParamItem.Enable = helper.IntInt64(v.(int))
				}
				taskAdvanceCFG.PortRisk = append(taskAdvanceCFG.PortRisk, &portRiskAdvanceCFGParamItem)
			}
		}

		if v, ok := dMap["vul_risk"]; ok {
			for _, item := range v.([]interface{}) {
				vulRiskMap := item.(map[string]interface{})
				taskCenterVulRiskInputParam := csip.TaskCenterVulRiskInputParam{}
				if v, ok := vulRiskMap["risk_id"]; ok {
					taskCenterVulRiskInputParam.RiskId = helper.String(v.(string))
				}

				if v, ok := vulRiskMap["enable"]; ok {
					taskCenterVulRiskInputParam.Enable = helper.IntInt64(v.(int))
				}

				taskAdvanceCFG.VulRisk = append(taskAdvanceCFG.VulRisk, &taskCenterVulRiskInputParam)
			}
		}

		if v, ok := dMap["weak_pwd_risk"]; ok {
			for _, item := range v.([]interface{}) {
				weakPwdRiskMap := item.(map[string]interface{})
				taskCenterWeakPwdRiskInputParam := csip.TaskCenterWeakPwdRiskInputParam{}
				if v, ok := weakPwdRiskMap["check_item_id"]; ok {
					taskCenterWeakPwdRiskInputParam.CheckItemId = helper.IntInt64(v.(int))
				}

				if v, ok := weakPwdRiskMap["enable"]; ok {
					taskCenterWeakPwdRiskInputParam.Enable = helper.IntInt64(v.(int))
				}

				taskAdvanceCFG.WeakPwdRisk = append(taskAdvanceCFG.WeakPwdRisk, &taskCenterWeakPwdRiskInputParam)
			}
		}
		if v, ok := dMap["cfg_risk"]; ok {
			for _, item := range v.([]interface{}) {
				cFGRiskMap := item.(map[string]interface{})
				taskCenterCFGRiskInputParam := csip.TaskCenterCFGRiskInputParam{}
				if v, ok := cFGRiskMap["item_id"]; ok {
					taskCenterCFGRiskInputParam.ItemId = helper.String(v.(string))
				}

				if v, ok := cFGRiskMap["enable"]; ok {
					taskCenterCFGRiskInputParam.Enable = helper.IntInt64(v.(int))
				}

				if v, ok := cFGRiskMap["resource_type"]; ok {
					taskCenterCFGRiskInputParam.ResourceType = helper.String(v.(string))
				}

				taskAdvanceCFG.CFGRisk = append(taskAdvanceCFG.CFGRisk, &taskCenterCFGRiskInputParam)
			}
		}

		request.TaskAdvanceCFG = &taskAdvanceCFG
	}

	if v, ok := d.GetOkExists("task_mode"); ok {
		request.TaskMode = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCsipClient().ModifyRiskCenterScanTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || *result.Response.Status != 0 {
			e = fmt.Errorf("update csip riskCenter failed, status: %d", *result.Response.Status)
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update csip riskCenter failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCsipRiskCenterRead(d, meta)
}

func resourceTencentCloudCsipRiskCenterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_csip_risk_center.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CsipService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  = d.Id()
	)

	// stop
	if err := service.StopCsipRiskCenterById(ctx, taskId); err != nil {
		return err
	}

	// delete
	if err := service.DeleteCsipRiskCenterById(ctx, taskId); err != nil {
		return err
	}

	return nil
}
