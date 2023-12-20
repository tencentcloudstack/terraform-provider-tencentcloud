package cynosdb

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbExportInstanceSlowQueries() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbExportInstanceSlowQueriesCreate,
		Read:   resourceTencentCloudCynosdbExportInstanceSlowQueriesRead,
		Delete: resourceTencentCloudCynosdbExportInstanceSlowQueriesDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"start_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Earliest transaction start time.",
			},

			"end_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Latest transaction start time.",
			},

			"username": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "user name.",
			},

			"host": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Client host.",
			},

			"database": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},

			"file_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File type, optional values: csv, original.",
			},

			"file_content": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Slow query export content.",
			},
		},
	}
}

func resourceTencentCloudCynosdbExportInstanceSlowQueriesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_export_instance_slow_queries.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cynosdb.NewExportInstanceSlowQueriesRequest()
		response   = cynosdb.NewExportInstanceSlowQueriesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database"); ok {
		request.Database = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_type"); ok {
		request.FileType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ExportInstanceSlowQueries(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb exportInstanceSlowQueries failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	if response.Response.FileContent != nil {
		_ = d.Set("file_content", response.Response.FileContent)
	}

	return resourceTencentCloudCynosdbExportInstanceSlowQueriesRead(d, meta)
}

func resourceTencentCloudCynosdbExportInstanceSlowQueriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_export_instance_slow_queries.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbExportInstanceSlowQueriesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_export_instance_slow_queries.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
