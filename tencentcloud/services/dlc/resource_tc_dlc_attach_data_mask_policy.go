package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcAttachDataMaskPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcAttachDataMaskPolicyCreate,
		Read:   resourceTencentCloudDlcAttachDataMaskPolicyRead,
		Delete: resourceTencentCloudDlcAttachDataMaskPolicyDelete,
		Schema: map[string]*schema.Schema{
			"data_mask_strategy_policy_set": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "A collection of data masking policy permission objects to be bound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_info": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Data masking permission object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The name of the database to be authorized. Use * to represent all databases under the current Catalog. For administrator-level authorization, only * is allowed. For data connection-level authorization, leave it empty. For other types, specify the database name.",
									},
									"catalog": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The name of the data source to be authorized. For administrator-level authorization, only * is allowed (representing all resources at this level). For data source-level and database-level authorization, only COSDataCatalog or * is allowed. For table-level authorization, custom data sources can be specified. Defaults to DataLakeCatalog if not specified. Note: For custom data sources, DLC can only manage a subset of permissions provided by the user during data source integration.",
									},
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The name of the table to be authorized. Use * to represent all tables under the current Database. For administrator-level authorization, only * is allowed. For data connection-level and database-level authorization, leave it empty. For other types, specify the table name.",
									},
									"column": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The name of the column to be authorized. Use * to represent all columns. For administrator-level authorization, only * is allowed.",
									},
								},
							},
						},
						"data_mask_strategy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The ID of the data masking strategy.",
						},
						"column_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The type of the bound field.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcAttachDataMaskPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_data_mask_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request            = dlcv20210125.NewAttachDataMaskPolicyRequest()
		catalog            string
		dataBase           string
		table              string
		column             string
		dataMaskStrategyId string
	)

	if v, ok := d.GetOk("data_mask_strategy_policy_set"); ok {
		for _, item := range v.([]interface{}) {
			dataMaskStrategyPolicySetMap := item.(map[string]interface{})
			dataMaskStrategyPolicy := dlcv20210125.DataMaskStrategyPolicy{}
			if policyInfoMap, ok := helper.ConvertInterfacesHeadToMap(dataMaskStrategyPolicySetMap["policy_info"]); ok {
				policy := dlcv20210125.Policy{}
				if v, ok := policyInfoMap["database"].(string); ok && v != "" {
					policy.Database = helper.String(v)
					dataBase = v
				}

				if v, ok := policyInfoMap["catalog"].(string); ok && v != "" {
					policy.Catalog = helper.String(v)
					catalog = v
				}

				if v, ok := policyInfoMap["table"].(string); ok && v != "" {
					policy.Table = helper.String(v)
					table = v
				}

				if v, ok := policyInfoMap["column"].(string); ok && v != "" {
					policy.Column = helper.String(v)
					column = v
				}

				policy.Operation = helper.String("SELECT")
				policy.PolicyType = helper.String("DATAMASK")
				dataMaskStrategyPolicy.PolicyInfo = &policy
			}

			if v, ok := dataMaskStrategyPolicySetMap["data_mask_strategy_id"].(string); ok && v != "" {
				dataMaskStrategyPolicy.DataMaskStrategyId = helper.String(v)
			}

			if v, ok := dataMaskStrategyPolicySetMap["column_type"].(string); ok && v != "" {
				dataMaskStrategyPolicy.ColumnType = helper.String(v)
			}

			request.DataMaskStrategyPolicySet = append(request.DataMaskStrategyPolicySet, &dataMaskStrategyPolicy)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AttachDataMaskPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc attach data mask policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	policyInfoStr := strings.Join([]string{catalog, dataBase, table, column}, tccommon.COMMA_SP)
	d.SetId(strings.Join([]string{policyInfoStr, dataMaskStrategyId}, tccommon.FILED_SP))
	return resourceTencentCloudDlcAttachDataMaskPolicyRead(d, meta)
}

func resourceTencentCloudDlcAttachDataMaskPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_data_mask_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	policyInfoStr := idSplit[0]
	secIdSplit := strings.Split(policyInfoStr, tccommon.COMMA_SP)
	if len(secIdSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	catalog := secIdSplit[0]
	dataBase := secIdSplit[1]
	table := secIdSplit[2]
	column := secIdSplit[3]

	respData, err := service.DescribeDlcAttachDataMaskPolicyById(ctx, catalog, dataBase, table)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_attach_data_mask_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Columns != nil {
		for _, item := range respData.Columns {
			if item.Name != nil && *item.Name == column {
				tmpList := make([]map[string]interface{}, 0, 1)
				dMap := make(map[string]interface{}, 0)
				policyInfoList := make([]map[string]interface{}, 0, 1)
				policyInfoMap := make(map[string]interface{}, 0)
				policyInfoMap["catalog"] = catalog
				policyInfoMap["database"] = dataBase
				policyInfoMap["table"] = table
				policyInfoMap["column"] = column

				policyInfoList = append(policyInfoList, policyInfoMap)
				dMap["policy_info"] = policyInfoList

				if item.DataMaskStrategyInfo != nil && item.DataMaskStrategyInfo.StrategyId != nil {
					dMap["data_mask_strategy_id"] = item.DataMaskStrategyInfo.StrategyId
				}

				if item.Type != nil {
					dMap["column_type"] = item.Type
				}

				tmpList = append(tmpList, dMap)
				_ = d.Set("data_mask_strategy_policy_set", tmpList)
			}
		}
	}

	return nil
}

func resourceTencentCloudDlcAttachDataMaskPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_data_mask_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlcv20210125.NewAttachDataMaskPolicyRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	policyInfoStr := idSplit[0]
	secIdSplit := strings.Split(policyInfoStr, tccommon.COMMA_SP)
	if len(secIdSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	catalog := secIdSplit[0]
	dataBase := secIdSplit[1]
	table := secIdSplit[2]
	column := secIdSplit[3]

	if v, ok := d.GetOk("data_mask_strategy_policy_set"); ok {
		for _, item := range v.([]interface{}) {
			dataMaskStrategyPolicySetMap := item.(map[string]interface{})
			dataMaskStrategyPolicy := dlcv20210125.DataMaskStrategyPolicy{}
			policy := dlcv20210125.Policy{}
			policy.Database = &dataBase
			policy.Catalog = &catalog
			policy.Table = &table
			policy.Column = &column
			policy.Operation = helper.String("SELECT")
			policy.PolicyType = helper.String("DATAMASK")
			dataMaskStrategyPolicy.PolicyInfo = &policy
			dataMaskStrategyPolicy.DataMaskStrategyId = helper.String("-1")

			if v, ok := dataMaskStrategyPolicySetMap["column_type"].(string); ok && v != "" {
				dataMaskStrategyPolicy.ColumnType = helper.String(v)
			}

			request.DataMaskStrategyPolicySet = append(request.DataMaskStrategyPolicySet, &dataMaskStrategyPolicy)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AttachDataMaskPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc attach data mask policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
