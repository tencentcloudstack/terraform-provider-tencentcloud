package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlDatabaseCreate,
		Read:   resourceTencentCloudMysqlDatabaseRead,
		Update: resourceTencentCloudMysqlDatabaseUpdate,
		Delete: resourceTencentCloudMysqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of `cdb-c1nl9rpv`,  which is the same as the one displayed in the TencentDB console.",
			},

			"db_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of Database.",
			},

			"character_set_name": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"utf8", "gbk", "latin1", "utf8mb4"}),
				Description:  "Character set. Valid values:  `utf8`, `gbk`, `latin1`, `utf8mb4`.",
			},
		},
	}
}

func resourceTencentCloudMysqlDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_database.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mysql.NewCreateDatabaseRequest()
		instanceId string
		dBName     string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		dBName = v.(string)
		request.DBName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("character_set_name"); ok {
		request.CharacterSetName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDatabase(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql database failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{instanceId, dBName}, tccommon.FILED_SP))

	return resourceTencentCloudMysqlDatabaseRead(d, meta)
}

func resourceTencentCloudMysqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_database.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	dBName := idSplit[1]

	database, err := service.DescribeMysqlDatabaseById(ctx, instanceId, dBName)
	if err != nil {
		return err
	}

	if database == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlDatabase` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if database.DatabaseName != nil {
		_ = d.Set("db_name", database.DatabaseName)
	}

	if database.CharacterSet != nil {
		if *database.CharacterSet == "UTF8MB3" {
			_ = d.Set("character_set_name", "utf8")
		} else {
			_ = d.Set("character_set_name", strings.ToLower(*database.CharacterSet))
		}

	}

	return nil
}

func resourceTencentCloudMysqlDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_database.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"character_set_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudMysqlDatabaseRead(d, meta)
}

func resourceTencentCloudMysqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_database.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	dBName := idSplit[1]

	if err := service.DeleteMysqlDatabaseById(ctx, instanceId, dBName); err != nil {
		return err
	}

	return nil
}
