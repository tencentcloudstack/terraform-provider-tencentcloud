package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentCloudDcGatewayCCNRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcGatewayCCNRoutesRead,
		Schema: map[string]*schema.Schema{
			"dcg_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the DCG to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the DCG route entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DCG.",
						},
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DCG route.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network address segment of IDC.",
						},
						"as_path": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "As path list of the BGP.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcGatewayCCNRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dc_gateway_ccn_routes.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id = d.Get("dcg_id").(string)
	)

	var infos, err = service.DescribeDirectConnectGatewayCcnRoutes(ctx, id)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["dcg_id"] = item.dcgId
		infoMap["route_id"] = item.routeId
		infoMap["cidr_block"] = item.cidrBlock
		infoMap["as_path"] = item.asPaths
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  dcg  ccn routes fail, reason:%s\n ",
			logId,
			err.Error())
		return err
	}

	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId,
				output.(string),
				err.Error())
			return err
		}
	}
	return nil

}
