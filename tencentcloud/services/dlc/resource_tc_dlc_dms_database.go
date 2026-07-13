package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDmsDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDmsDatabaseCreate,
		Read:   resourceTencentCloudDlcDmsDatabaseRead,
		Update: resourceTencentCloudDlcDmsDatabaseUpdate,
		Delete: resourceTencentCloudDlcDmsDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database name.",
			},

			"schema_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Schema name.",
			},

			"datasource_connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datasource connection name.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Db storage path.",
			},

			"asset": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Basic metadata object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary key.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name.",
						},
						"guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object GUID.",
						},
						"catalog": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data catalog.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"owner": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object owner.",
						},
						"owner_account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object owner account.",
						},
						"perm_values": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Permissions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configured key value.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configured value.",
									},
								},
							},
						},
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configured key value.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configured value.",
									},
								},
							},
						},
						"biz_params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional business attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configured key value.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configured value.",
									},
								},
							},
						},
						"data_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data version.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modified time.",
						},
						"datasource_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Datasource primary key.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDmsDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_database.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlc.NewCreateDMSDatabaseRequest()
	)

	name := d.Get("name").(string)
	schemaName := d.Get("schema_name").(string)
	datasourceConnectionName := d.Get("datasource_connection_name").(string)

	request.Name = helper.String(name)
	request.SchemaName = helper.String(schemaName)
	request.DatasourceConnectionName = helper.String(datasourceConnectionName)

	if v, ok := d.GetOk("location"); ok {
		request.Location = helper.String(v.(string))
	}

	if v, ok := d.GetOk("asset"); ok {
		request.Asset = buildDlcDmsDatabaseAsset(v.([]interface{}))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateDMSDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc_dms_database failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc_dms_database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[INFO]%s create dlc_dms_database success, name=%s", logId, name)
	d.SetId(strings.Join([]string{name, schemaName, datasourceConnectionName}, tccommon.FILED_SP))
	return resourceTencentCloudDlcDmsDatabaseRead(d, meta)
}

func resourceTencentCloudDlcDmsDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_database.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlc.NewDescribeDMSDatabaseRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	name := idSplit[0]
	schemaName := idSplit[1]
	datasourceConnectionName := idSplit[2]

	request.Name = helper.String(name)
	request.SchemaName = helper.String(schemaName)
	request.DatasourceConnectionName = helper.String(datasourceConnectionName)

	var response *dlc.DescribeDMSDatabaseResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDMSDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc_dms_database failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read dlc_dms_database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	resp := response.Response
	if resp == nil || resp.Name == nil {
		log.Printf("[CRUD] dlc_dms_database id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if resp.Name != nil {
		_ = d.Set("name", resp.Name)
	}

	if resp.SchemaName != nil {
		_ = d.Set("schema_name", resp.SchemaName)
	}

	_ = d.Set("datasource_connection_name", datasourceConnectionName)

	if resp.Location != nil {
		_ = d.Set("location", resp.Location)
	}

	if resp.Asset != nil {
		assetList := flattenDlcDmsDatabaseAsset(resp.Asset)
		if len(assetList) > 0 {
			_ = d.Set("asset", assetList)
		}
	}

	return nil
}

func resourceTencentCloudDlcDmsDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_database.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	name := idSplit[0]
	schemaName := idSplit[1]
	datasourceConnectionName := idSplit[2]

	immutableArgs := []string{"name", "schema_name", "datasource_connection_name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("dlc_dms_database argument `%s` is immutable, can not change.", v)
		}
	}

	needChange := false
	mutableArgs := []string{"location", "asset"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := dlc.NewAlterDMSDatabaseRequest()
		request.CurrentName = helper.String(name)
		request.SchemaName = helper.String(schemaName)
		request.DatasourceConnectionName = helper.String(datasourceConnectionName)

		if v, ok := d.GetOk("location"); ok {
			request.Location = helper.String(v.(string))
		}

		if v, ok := d.GetOk("asset"); ok {
			request.Asset = buildDlcDmsDatabaseAsset(v.([]interface{}))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AlterDMSDatabaseWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Alter dlc_dms_database failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc_dms_database failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDlcDmsDatabaseRead(d, meta)
}

func resourceTencentCloudDlcDmsDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_dms_database.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlc.NewDropDMSDatabaseRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	name := idSplit[0]
	datasourceConnectionName := idSplit[2]

	request.Name = helper.String(name)
	request.DatasourceConnectionName = helper.String(datasourceConnectionName)
	request.DeleteData = helper.Bool(d.Get("delete_data").(bool))
	request.Cascade = helper.Bool(d.Get("cascade").(bool))

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DropDMSDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Drop dlc_dms_database failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc_dms_database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildDlcDmsDatabaseAsset builds an *dlc.Asset from the HCL asset block list.
func buildDlcDmsDatabaseAsset(assetList []interface{}) *dlc.Asset {
	if len(assetList) == 0 {
		return nil
	}

	assetMap := assetList[0].(map[string]interface{})
	asset := &dlc.Asset{}

	if v, ok := assetMap["name"].(string); ok && v != "" {
		asset.Name = helper.String(v)
	}

	if v, ok := assetMap["catalog"].(string); ok && v != "" {
		asset.Catalog = helper.String(v)
	}

	if v, ok := assetMap["description"].(string); ok && v != "" {
		asset.Description = helper.String(v)
	}

	if v, ok := assetMap["owner"].(string); ok && v != "" {
		asset.Owner = helper.String(v)
	}

	if v, ok := assetMap["owner_account"].(string); ok && v != "" {
		asset.OwnerAccount = helper.String(v)
	}

	if v, ok := assetMap["perm_values"].([]interface{}); ok && len(v) > 0 {
		asset.PermValues = buildDlcDmsDatabaseKVPairs(v)
	}

	if v, ok := assetMap["params"].([]interface{}); ok && len(v) > 0 {
		asset.Params = buildDlcDmsDatabaseKVPairs(v)
	}

	if v, ok := assetMap["biz_params"].([]interface{}); ok && len(v) > 0 {
		asset.BizParams = buildDlcDmsDatabaseKVPairs(v)
	}

	return asset
}

// buildDlcDmsDatabaseKVPairs builds a []*dlc.KVPair from the HCL kv pair list.
func buildDlcDmsDatabaseKVPairs(kvList []interface{}) []*dlc.KVPair {
	kvPairs := make([]*dlc.KVPair, 0, len(kvList))
	for _, item := range kvList {
		kvMap := item.(map[string]interface{})
		kvPair := &dlc.KVPair{}
		if v, ok := kvMap["key"].(string); ok && v != "" {
			kvPair.Key = helper.String(v)
		}
		if v, ok := kvMap["value"].(string); ok && v != "" {
			kvPair.Value = helper.String(v)
		}
		kvPairs = append(kvPairs, kvPair)
	}
	return kvPairs
}

// flattenDlcDmsDatabaseAsset converts an *dlc.Asset to the HCL asset block list.
func flattenDlcDmsDatabaseAsset(asset *dlc.Asset) []map[string]interface{} {
	if asset == nil {
		return nil
	}

	assetMap := map[string]interface{}{}

	if asset.Id != nil {
		assetMap["id"] = *asset.Id
	}

	if asset.Name != nil {
		assetMap["name"] = *asset.Name
	}

	if asset.Guid != nil {
		assetMap["guid"] = *asset.Guid
	}

	if asset.Catalog != nil {
		assetMap["catalog"] = *asset.Catalog
	}

	if asset.Description != nil {
		assetMap["description"] = *asset.Description
	}

	if asset.Owner != nil {
		assetMap["owner"] = *asset.Owner
	}

	if asset.OwnerAccount != nil {
		assetMap["owner_account"] = *asset.OwnerAccount
	}

	if asset.PermValues != nil {
		assetMap["perm_values"] = flattenDlcDmsDatabaseKVPairs(asset.PermValues)
	}

	if asset.Params != nil {
		assetMap["params"] = flattenDlcDmsDatabaseKVPairs(asset.Params)
	}

	if asset.BizParams != nil {
		assetMap["biz_params"] = flattenDlcDmsDatabaseKVPairs(asset.BizParams)
	}

	if asset.DataVersion != nil {
		assetMap["data_version"] = *asset.DataVersion
	}

	if asset.CreateTime != nil {
		assetMap["create_time"] = *asset.CreateTime
	}

	if asset.ModifiedTime != nil {
		assetMap["modified_time"] = *asset.ModifiedTime
	}

	if asset.DatasourceId != nil {
		assetMap["datasource_id"] = *asset.DatasourceId
	}

	return []map[string]interface{}{assetMap}
}

// flattenDlcDmsDatabaseKVPairs converts a []*dlc.KVPair to the HCL kv pair list.
func flattenDlcDmsDatabaseKVPairs(kvPairs []*dlc.KVPair) []map[string]interface{} {
	kvList := make([]map[string]interface{}, 0, len(kvPairs))
	for _, kvPair := range kvPairs {
		kvMap := map[string]interface{}{}
		if kvPair.Key != nil {
			kvMap["key"] = *kvPair.Key
		}
		if kvPair.Value != nil {
			kvMap["value"] = *kvPair.Value
		}
		kvList = append(kvList, kvMap)
	}
	return kvList
}
