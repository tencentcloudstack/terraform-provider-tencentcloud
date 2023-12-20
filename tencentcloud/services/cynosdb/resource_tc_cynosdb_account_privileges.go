package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbAccountPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbAccountPrivilegesCreate,
		Read:   resourceTencentCloudCynosdbAccountPrivilegesRead,
		Update: resourceTencentCloudCynosdbAccountPrivilegesUpdate,
		Delete: resourceTencentCloudCynosdbAccountPrivilegesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Host, default `%`.",
			},

			"global_privileges": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of global permissions.",
			},

			"database_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Array of database permissions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Database privileges.",
						},
					},
				},
			},

			"table_privileges": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "array of table permissions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Table name.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Table privileges.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCynosdbAccountPrivilegesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account_privileges.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		clusterId   string
		accountName string
		host        string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}
	if v, ok := d.GetOk("account_name"); ok {
		accountName = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}

	d.SetId(clusterId + tccommon.FILED_SP + accountName + tccommon.FILED_SP + host)

	return resourceTencentCloudCynosdbAccountPrivilegesUpdate(d, meta)
}

func resourceTencentCloudCynosdbAccountPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account_privileges.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	accountName := idSplit[1]
	host := idSplit[2]

	accountPrivileges, err := service.DescribeCynosdbAccountPrivilegesById(ctx, clusterId, accountName, host)
	if err != nil {
		return err
	}

	if accountPrivileges == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbAccountPrivileges` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("account_name", accountName)
	_ = d.Set("host", host)

	if accountPrivileges.GlobalPrivileges != nil {
		_ = d.Set("global_privileges", accountPrivileges.GlobalPrivileges)
	}

	if accountPrivileges.DatabasePrivileges != nil {
		databasePrivilegesList := []interface{}{}
		for _, databasePrivileges := range accountPrivileges.DatabasePrivileges {
			databasePrivilegesMap := map[string]interface{}{}

			if databasePrivileges.Db != nil {
				databasePrivilegesMap["db"] = databasePrivileges.Db
			}

			if databasePrivileges.Privileges != nil {
				databasePrivilegesMap["privileges"] = databasePrivileges.Privileges
			}

			databasePrivilegesList = append(databasePrivilegesList, databasePrivilegesMap)
		}

		_ = d.Set("database_privileges", databasePrivilegesList)

	}

	if accountPrivileges.TablePrivileges != nil {
		tablePrivilegesList := []interface{}{}
		for _, tablePrivileges := range accountPrivileges.TablePrivileges {
			tablePrivilegesMap := map[string]interface{}{}

			if tablePrivileges.Db != nil {
				tablePrivilegesMap["db"] = tablePrivileges.Db
			}

			if tablePrivileges.TableName != nil {
				tablePrivilegesMap["table_name"] = tablePrivileges.TableName
			}

			if tablePrivileges.Privileges != nil {
				tablePrivilegesMap["privileges"] = tablePrivileges.Privileges
			}

			tablePrivilegesList = append(tablePrivilegesList, tablePrivilegesMap)
		}

		_ = d.Set("table_privileges", tablePrivilegesList)

	}

	return nil
}

func resourceTencentCloudCynosdbAccountPrivilegesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account_privileges.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cynosdb.NewModifyAccountPrivilegesRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	accountName := idSplit[1]
	host := idSplit[2]

	request.ClusterId = &clusterId
	request.Account = &cynosdb.InputAccount{
		AccountName: &accountName,
		Host:        &host,
	}

	if d.HasChange("global_privileges") {
		if v, ok := d.GetOk("global_privileges"); ok {
			globalPrivilegesSet := v.(*schema.Set).List()
			for i := range globalPrivilegesSet {
				globalPrivileges := globalPrivilegesSet[i].(string)
				request.GlobalPrivileges = append(request.GlobalPrivileges, &globalPrivileges)
			}
		}
	}

	if d.HasChange("database_privileges") {
		if v, ok := d.GetOk("database_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				databasePrivileges := cynosdb.DatabasePrivileges{}
				if v, ok := dMap["db"]; ok {
					databasePrivileges.Db = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						databasePrivileges.Privileges = append(databasePrivileges.Privileges, &privileges)
					}
				}
				request.DatabasePrivileges = append(request.DatabasePrivileges, &databasePrivileges)
			}
		}
	}

	if d.HasChange("table_privileges") {
		if v, ok := d.GetOk("table_privileges"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				tablePrivileges := cynosdb.TablePrivileges{}
				if v, ok := dMap["db"]; ok {
					tablePrivileges.Db = helper.String(v.(string))
				}
				if v, ok := dMap["table_name"]; ok {
					tablePrivileges.TableName = helper.String(v.(string))
				}
				if v, ok := dMap["privileges"]; ok {
					privilegesSet := v.(*schema.Set).List()
					for i := range privilegesSet {
						privileges := privilegesSet[i].(string)
						tablePrivileges.Privileges = append(tablePrivileges.Privileges, &privileges)
					}
				}
				request.TablePrivileges = append(request.TablePrivileges, &tablePrivileges)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyAccountPrivileges(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb accountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbAccountPrivilegesRead(d, meta)
}

func resourceTencentCloudCynosdbAccountPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_account_privileges.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
