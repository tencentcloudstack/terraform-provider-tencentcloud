/*
Use this data source to query detailed information of security group.

Example Usage

```hcl
resource "tencentcloud_security_group" "sglab" {
  name        = "mysg"
  description = "favourite sg"
  project_id  = 0
}
data "tencentcloud_security_group" "sglab" {
  security_group_id = "${tencentcloud_security_group.sglab.id}"
}
```
*/
package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the security group to be queried.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Name of the security group to be queried.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(2, 100),
				Description:  "Description of the security group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of security group.",
			},
			"be_associate_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of security group binding resources.",
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn
	params := map[string]string{
		"Action": "DescribeSecurityGroupEx",
		"sgId":   d.Get("security_group_id").(string),
	}

	log.Printf("[DEBUG] resource_tc_security_group read params:%v", params)

	response, err := client.SendRequest("dfw", params)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group read client.SendRequest error:%v", err)
		return err
	}

	var jsonresp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TotalNum int `json:"totalNum"`
			Detail   []struct {
				SgId             string      `json:"sgId"`
				SgName           string      `json:"sgName"`
				SgRemark         string      `json:"sgRemark"`
				BeAssociateCount int         `json:"beAssociateCount"`
				CreateTime       string      `json:"createTime"`
				ProjectId        interface{} `json:"projectId"`
			}
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		log.Printf("[ERROR] resource_tc_security_group read json.Unmarshal error:%v", err)
		return err
	}
	if jsonresp.Code != 0 {
		log.Printf("[ERROR] resource_tc_security_group read error, code:%v, message:%v", jsonresp.Code, jsonresp.Message)
		return errors.New(jsonresp.Message)
	} else if jsonresp.Data.TotalNum <= 0 || len(jsonresp.Data.Detail) <= 0 {
		return errors.New("Security group not found")
	}

	sg := jsonresp.Data.Detail[0]

	d.SetId(sg.SgId)
	d.Set("security_group_id", sg.SgId)
	d.Set("name", sg.SgName)
	d.Set("description", sg.SgRemark)
	d.Set("create_time", sg.CreateTime)
	d.Set("be_associate_count", sg.BeAssociateCount)

	if "string" == reflect.TypeOf(sg.ProjectId).String() {
		if intVal, err := strconv.ParseInt(sg.ProjectId.(string), 10, 64); err != nil {
			return fmt.Errorf("create security_group project ParseInt  error ,%s", err.Error())
		} else {
			d.Set("project_id", int(intVal))
		}

	} else {
		d.Set("project_id", sg.ProjectId.(int))
	}
	return nil
}
