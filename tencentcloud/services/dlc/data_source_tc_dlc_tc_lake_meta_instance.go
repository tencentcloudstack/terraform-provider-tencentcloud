package dlc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudDlcTCLakeMetaInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcTCLakeMetaInstanceRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "TCLake activation status. Enumerated value: `Running` means activated successfully.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcTCLakeMetaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_tc_lake_meta_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var status *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcTCLakeMetaInstance(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil {
			log.Printf("[DATASOURCE] dlc tc_lake_meta_instance read empty, skip SetId")
			return resource.NonRetryableError(fmt.Errorf("dlc tc_lake_meta_instance read empty, skip SetId"))
		}
		status = result
		return nil
	})
	if err != nil {
		return err
	}

	if status != nil {
		_ = d.Set("status", status)
	}

	d.SetId("tc_lake_meta_instance")
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), status); e != nil {
			return e
		}
	}
	return nil
}
