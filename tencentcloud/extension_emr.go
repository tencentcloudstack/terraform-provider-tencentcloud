package tencentcloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

const (
	EmrInternetStatusCreated int64 = 2
	EmrInternetStatusDeleted int64 = 201
)

const (
	DisplayStrategyIsclusterList = "clusterList"
)

func buildResourceSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"spec":         {Type: schema.TypeString, Optional: true},
				"storage_type": {Type: schema.TypeInt, Optional: true},
				"disk_type":    {Type: schema.TypeString, Optional: true},
				"mem_size":     {Type: schema.TypeInt, Optional: true},
				"cpu":          {Type: schema.TypeInt, Optional: true},
				"disk_size":    {Type: schema.TypeInt, Optional: true},
			},
		},
	}
}
