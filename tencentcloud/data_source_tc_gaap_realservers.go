/*
Use this data source to query gaap realservers.

Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

data "tencentcloud_gaap_realservers" "foo" {
  ip = "${tencentcloud_gaap_realserver.foo.ip}"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapRealservers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapRealserversRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project within the GAAP realserver to be queried, default is '-1' means all projects.",
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ip"},
				Description:   "Domain of the GAAP realserver to be queried, conflict with `ip`.",
			},
			"ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain"},
				Description:   "IP of the GAAP realserver to be queried, conflict with `domain`.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the GAAP realserver to be queried, the maximum length is 30.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the GAAP proxy to be queried. Support up to 5, display the information as long as it matches one.",
			},

			// computed
			"realservers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of GAAP realserver. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the GAAP realserver.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the GAAP realserver.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the GAAP realserver.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain of the GAAP realserver.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project within the GAAP realserver.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the GAAP realserver.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapRealserversRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_realservers.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	projectId := -1
	if raw, ok := d.GetOk("project_id"); ok {
		projectId = raw.(int)
	}

	var (
		address *string
		name    *string
	)
	if raw, ok := d.GetOk("ip"); ok {
		address = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("domain"); ok {
		address = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}

	tags := getTags(d, "tags")

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	realservers, err := service.DescribeRealservers(ctx, address, name, tags, projectId)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(realservers))
	realserverList := make([]map[string]interface{}, 0, len(realservers))
	for _, rs := range realservers {
		if rs.RealServerId == nil {
			return errors.New("realserver id is nil")
		}
		if rs.RealServerName == nil {
			return errors.New("realserver name is nil")
		}
		if rs.RealServerIP == nil {
			return errors.New("realserver name is nil")
		}
		if rs.ProjectId == nil {
			return errors.New("realserver project id is nil")
		}

		ids = append(ids, *rs.RealServerId)

		m := map[string]interface{}{
			"id":         *rs.RealServerId,
			"name":       *rs.RealServerName,
			"project_id": *rs.ProjectId,
		}

		if net.ParseIP(*rs.RealServerIP) == nil {
			m["domain"] = *rs.RealServerIP
		} else {
			m["ip"] = *rs.RealServerIP
		}

		if len(rs.TagSet) > 0 {
			tags := make(map[string]string, len(rs.TagSet))
			for _, tag := range rs.TagSet {
				if tag.TagKey == nil {
					return errors.New("tag key is nil")
				}
				if tag.TagValue == nil {
					return errors.New("tag value is nil")
				}
				tags[*tag.TagKey] = *tag.TagValue
			}
			m["tags"] = tags
		}

		realserverList = append(realserverList, m)
	}

	d.Set("realservers", realserverList)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
