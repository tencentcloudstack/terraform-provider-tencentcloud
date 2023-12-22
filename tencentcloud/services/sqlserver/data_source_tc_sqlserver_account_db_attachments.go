package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverAccountDBAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverAccountDBAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL Server instance ID that the account belongs to.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the SQL Server account to be queried.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the DB to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SQL Server account. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL Server instance ID that the account belongs to.",
						},
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL Server account name.",
						},
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL Server DB name.",
						},
						"privilege": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Privilege of the account on DB. Valid value are `ReadOnly`, `ReadWrite`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverAccountDBAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_account_db_attachments.read")()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		sqlserverService = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	dbName := d.Get("db_name").(string)

	attachments, err := sqlserverService.DescribeAccountDBAttachments(ctx, instanceId, accountName, dbName)

	if err != nil {
		return fmt.Errorf("api[DescribeAccountDBAttachments]fail, return %s", err.Error())
	}

	var list []map[string]interface{}
	var ids = make([]string, len(attachments))

	for _, item := range attachments {
		mapping := map[string]interface{}{
			"instance_id":  instanceId,
			"account_name": item["account_name"],
			"db_name":      item["db_name"],
			"privilege":    item["privilege"],
		}

		list = append(list, mapping)
		ids = append(ids, fmt.Sprintf("%s%s%s%s%s", instanceId, tccommon.FILED_SP, accountName, tccommon.FILED_SP, dbName))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
