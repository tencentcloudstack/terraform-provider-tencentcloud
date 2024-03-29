package tcaplusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudTcaplusIdls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusIdlsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the TcaplusDB cluster to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File for saving results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of TcaplusDB table IDL. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the IDL.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusIdlsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcaplus_idls.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcaplusService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	clusterId := d.Get("cluster_id").(string)

	infos, err := service.DescribeIdlFileInfos(ctx, clusterId)
	if err != nil {
		infos, err = service.DescribeIdlFileInfos(ctx, clusterId)
	}
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(infos))

	for _, info := range infos {
		listItem := make(map[string]interface{})
		var tcaplusIdlId TcaplusIdlId
		tcaplusIdlId.ClusterId = clusterId
		tcaplusIdlId.FileName = *info.FileName
		tcaplusIdlId.FileType = *info.FileType
		tcaplusIdlId.FileExtType = *info.FileExtType
		tcaplusIdlId.FileSize = *info.FileSize
		tcaplusIdlId.FileId = *info.FileId
		id, err := json.Marshal(tcaplusIdlId)
		if err != nil {
			return fmt.Errorf("format idl id fail,%s", err.Error())
		}
		listItem["idl_id"] = string(id)
		list = append(list, listItem)
	}

	d.SetId("idl." + clusterId)
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
