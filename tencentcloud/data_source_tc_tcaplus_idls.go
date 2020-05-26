/*
Use this data source to query tcaplus idl files

Example Usage

```hcl
data "tencentcloud_tcaplus_idls" "id_test" {
  cluster_id = "19162256624"
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudTcaplusIdls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusIdlsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the tcaplus cluster to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tcaplus idls. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of this idl.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusIdlsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_idls.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
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
		return writeToFile(output.(string), list)
	}
	return nil

}
