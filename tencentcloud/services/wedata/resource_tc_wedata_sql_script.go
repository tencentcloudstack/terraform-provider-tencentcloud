package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataSqlScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataSqlScriptCreate,
		Read:   resourceTencentCloudWedataSqlScriptRead,
		Update: resourceTencentCloudWedataSqlScriptUpdate,
		Delete: resourceTencentCloudWedataSqlScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"script_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Script name.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Parent folder path, /aaa/bbb/ccc, root directory is empty string or /.",
			},

			"script_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Data exploration script configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datasource_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data source ID.",
						},
						"datasource_env": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data source environment.",
						},
						"compute_resource": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Computing resource.",
						},
						"executor_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution resource group.",
						},
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Advanced runtime parameters, variable substitution, map-json String,String.",
						},
						"advance_config": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Advanced settings, execution configuration parameters, map-json String,String. Encoded in Base64.",
						},
					},
				},
			},

			"script_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script content, if there is a value.",
			},

			"access_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Permission scope: SHARED, PRIVATE.",
			},

			// computed
			"script_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script ID.",
			},

			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path of the node, /aaa/bbb/ccc.ipynb, consists of the names of each node.",
			},
		},
	}
}

func resourceTencentCloudWedataSqlScriptCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_script.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewCreateSQLScriptRequest()
		response  = wedatav20250806.NewCreateSQLScriptResponse()
		projectId string
		scriptId  string
	)

	if v, ok := d.GetOk("script_name"); ok {
		request.ScriptName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	if scriptConfigMap, ok := helper.InterfacesHeadMap(d, "script_config"); ok {
		sQLScriptConfig := wedatav20250806.SQLScriptConfig{}
		if v, ok := scriptConfigMap["datasource_id"].(string); ok && v != "" {
			sQLScriptConfig.DatasourceId = helper.String(v)
		}

		if v, ok := scriptConfigMap["datasource_env"].(string); ok && v != "" {
			sQLScriptConfig.DatasourceEnv = helper.String(v)
		}

		if v, ok := scriptConfigMap["compute_resource"].(string); ok && v != "" {
			sQLScriptConfig.ComputeResource = helper.String(v)
		}

		if v, ok := scriptConfigMap["executor_group_id"].(string); ok && v != "" {
			sQLScriptConfig.ExecutorGroupId = helper.String(v)
		}

		if v, ok := scriptConfigMap["params"].(string); ok && v != "" {
			sQLScriptConfig.Params = helper.String(v)
		}

		if v, ok := scriptConfigMap["advance_config"].(string); ok && v != "" {
			sQLScriptConfig.AdvanceConfig = helper.String(v)
		}

		request.ScriptConfig = &sQLScriptConfig
	}

	if v, ok := d.GetOk("script_content"); ok {
		request.ScriptContent = helper.String(tccommon.StringToBase64(v.(string)))
	}

	if v, ok := d.GetOk("access_scope"); ok {
		request.AccessScope = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateSQLScriptWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata sql script failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata sql script failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.ScriptId == nil {
		return fmt.Errorf("ScriptId is nil.")
	}

	scriptId = *response.Response.Data.ScriptId
	d.SetId(strings.Join([]string{projectId, scriptId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataSqlScriptRead(d, meta)
}

func resourceTencentCloudWedataSqlScriptRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_script.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	scriptId := idSplit[1]

	respData, err := service.DescribeWedataSqlScriptById(ctx, projectId, scriptId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_sql_script` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ScriptName != nil {
		_ = d.Set("script_name", respData.ScriptName)
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.ParentFolderPath != nil {
		_ = d.Set("parent_folder_path", respData.ParentFolderPath)
	}

	if respData.ScriptConfig != nil {
		scriptConfigMap := map[string]interface{}{}
		if respData.ScriptConfig.DatasourceId != nil {
			scriptConfigMap["datasource_id"] = respData.ScriptConfig.DatasourceId
		}

		if respData.ScriptConfig.DatasourceEnv != nil {
			scriptConfigMap["datasource_env"] = respData.ScriptConfig.DatasourceEnv
		}

		if respData.ScriptConfig.ComputeResource != nil {
			scriptConfigMap["compute_resource"] = respData.ScriptConfig.ComputeResource
		}

		if respData.ScriptConfig.ExecutorGroupId != nil {
			scriptConfigMap["executor_group_id"] = respData.ScriptConfig.ExecutorGroupId
		}

		if respData.ScriptConfig.Params != nil {
			scriptConfigMap["params"] = respData.ScriptConfig.Params
		}

		if respData.ScriptConfig.AdvanceConfig != nil {
			scriptConfigMap["advance_config"] = respData.ScriptConfig.AdvanceConfig
		}

		_ = d.Set("script_config", []interface{}{scriptConfigMap})
	}

	if respData.ScriptContent != nil {
		sqlStr, err := tccommon.Base64ToString(*respData.ScriptContent)
		if err != nil {
			log.Printf("[ERROR]%s base64 decode failed, reason:%+v", logId, err)
			return err
		}

		_ = d.Set("script_content", sqlStr)
	}

	if respData.AccessScope != nil {
		_ = d.Set("access_scope", respData.AccessScope)
	}

	if respData.ScriptId != nil {
		_ = d.Set("script_id", respData.ScriptId)
	}

	if respData.Path != nil {
		_ = d.Set("path", *respData.Path)
	}

	return nil
}

func resourceTencentCloudWedataSqlScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_script.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	scriptId := idSplit[1]

	needChange := false
	mutableArgs := []string{"script_config", "script_content"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateSQLScriptRequest()
		if scriptConfigMap, ok := helper.InterfacesHeadMap(d, "script_config"); ok {
			sQLScriptConfig := wedatav20250806.SQLScriptConfig{}
			if v, ok := scriptConfigMap["datasource_id"].(string); ok && v != "" {
				sQLScriptConfig.DatasourceId = helper.String(v)
			}

			if v, ok := scriptConfigMap["datasource_env"].(string); ok && v != "" {
				sQLScriptConfig.DatasourceEnv = helper.String(v)
			}

			if v, ok := scriptConfigMap["compute_resource"].(string); ok && v != "" {
				sQLScriptConfig.ComputeResource = helper.String(v)
			}

			if v, ok := scriptConfigMap["executor_group_id"].(string); ok && v != "" {
				sQLScriptConfig.ExecutorGroupId = helper.String(v)
			}

			if v, ok := scriptConfigMap["params"].(string); ok && v != "" {
				sQLScriptConfig.Params = helper.String(v)
			}

			if v, ok := scriptConfigMap["advance_config"].(string); ok && v != "" {
				sQLScriptConfig.AdvanceConfig = helper.String(v)
			}

			request.ScriptConfig = &sQLScriptConfig
		}

		if v, ok := d.GetOk("script_content"); ok {
			request.ScriptContent = helper.String(v.(string))
		}

		request.ProjectId = &projectId
		request.ScriptId = &scriptId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateSQLScriptWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata sql script failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWedataSqlScriptRead(d, meta)
}

func resourceTencentCloudWedataSqlScriptDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_sql_script.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDeleteSQLScriptRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	scriptId := idSplit[1]

	request.ProjectId = &projectId
	request.ScriptId = &scriptId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteSQLScriptWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata sql script failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata sql script failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata sql script failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
