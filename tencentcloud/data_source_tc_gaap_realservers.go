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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ip"},
			},
			"ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// computed
			"realservers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
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

	projectId := d.Get("project_id").(int)

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
