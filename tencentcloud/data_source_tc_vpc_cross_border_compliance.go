/*
Use this data source to query detailed information of vpc cross_border_compliance

Example Usage

```hcl
data "tencentcloud_vpc_cross_border_compliance" "cross_border_compliance" {
  service_provider = "UNICOM"
  compliance_id = 10002
  company = "腾讯科技（广州）有限公司"
  uniform_social_credit_code = "91440101327598294H"
  legal_person = "张颖"
  issuing_authority = "广州市海珠区市场监督管理局"
  business_address = "广州市海珠区新港中路397号自编72号(商业街F5-1)"
  post_code = 510320
  manager = "李四"
  manager_id = "360732199007108888"
  manager_address = "广州市海珠区新港中路8888号"
  manager_telephone = "020-81167888"
  email = "test@tencent.com"
  service_start_date = "2020-07-29"
  service_end_date = "2021-07-29"
  state = "APPROVED"
  offset = 1
  limit = 2
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcCrossBorderCompliance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcCrossBorderComplianceRead,
		Schema: map[string]*schema.Schema{
			"service_provider": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) service provider, optional value: &amp;amp;#39;UNICOM&amp;amp;#39;.",
			},

			"compliance_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "(Exact match) compliance approval form: &amp;amp;#39;ID&amp;amp;#39;.",
			},

			"company": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy query) Company name.",
			},

			"uniform_social_credit_code": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) Uniform Social Credit Code.",
			},

			"legal_person": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy query) legal representative.",
			},

			"issuing_authority": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy query) Issuing authority.",
			},

			"business_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy query) business license address.",
			},

			"post_code": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "(Exact match) post code.",
			},

			"manager": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy query) Person in charge.",
			},

			"manager_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact query) ID number of the person in charge.",
			},

			"manager_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy query) ID card address of the person in charge.",
			},

			"manager_telephone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) contact number of the person in charge.",
			},

			"email": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) email.",
			},

			"service_start_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) service start date, such as: &amp;amp;#39;2020-07-28&amp;amp;#39;.",
			},

			"service_end_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) service end date, such as: &amp;amp;#39;2020-07-28&amp;amp;#39;.",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) status. Pending: PENDING, Passed: APPROVED, Denied: DENY.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return quantity.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcCrossBorderComplianceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_cross_border_compliance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_provider"); ok {
		paramMap["ServiceProvider"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("compliance_id"); v != nil {
		paramMap["ComplianceId"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("company"); ok {
		paramMap["Company"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniform_social_credit_code"); ok {
		paramMap["UniformSocialCreditCode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("legal_person"); ok {
		paramMap["LegalPerson"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("issuing_authority"); ok {
		paramMap["IssuingAuthority"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("business_address"); ok {
		paramMap["BusinessAddress"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("post_code"); v != nil {
		paramMap["PostCode"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("manager"); ok {
		paramMap["Manager"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_id"); ok {
		paramMap["ManagerId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_address"); ok {
		paramMap["ManagerAddress"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_telephone"); ok {
		paramMap["ManagerTelephone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		paramMap["Email"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_start_date"); ok {
		paramMap["ServiceStartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_end_date"); ok {
		paramMap["ServiceEndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["State"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var crossBorderComplianceSet []*vpc.CrossBorderCompliance

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCrossBorderComplianceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		crossBorderComplianceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(crossBorderComplianceSet))
	tmpList := make([]map[string]interface{}, 0, len(crossBorderComplianceSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
