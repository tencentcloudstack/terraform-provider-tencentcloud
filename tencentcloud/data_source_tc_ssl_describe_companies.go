package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeCompanies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeCompaniesRead,
		Schema: map[string]*schema.Schema{
			"company_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Company ID.",
			},

			"companies": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Company list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"company_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Company Name.",
						},
						"company_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Company ID.",
						},
						"company_country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Company country.",
						},
						"company_province": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Province where the company is located.",
						},
						"company_city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The city where the company is.",
						},
						"company_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detailed address where the company is located.",
						},
						"company_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "company phone.",
						},
						"id_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "typeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"id_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID numberNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeCompaniesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_companies.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("company_id"); v != nil {
		paramMap["CompanyId"] = helper.IntInt64(v.(int))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var companies []*ssl.CompanyInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCompaniesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		companies = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(companies))
	tmpList := make([]map[string]interface{}, 0, len(companies))

	if companies != nil {
		for _, companyInfo := range companies {
			companyInfoMap := map[string]interface{}{}

			if companyInfo.CompanyName != nil {
				companyInfoMap["company_name"] = companyInfo.CompanyName
			}

			if companyInfo.CompanyId != nil {
				companyInfoMap["company_id"] = companyInfo.CompanyId
			}

			if companyInfo.CompanyCountry != nil {
				companyInfoMap["company_country"] = companyInfo.CompanyCountry
			}

			if companyInfo.CompanyProvince != nil {
				companyInfoMap["company_province"] = companyInfo.CompanyProvince
			}

			if companyInfo.CompanyCity != nil {
				companyInfoMap["company_city"] = companyInfo.CompanyCity
			}

			if companyInfo.CompanyAddress != nil {
				companyInfoMap["company_address"] = companyInfo.CompanyAddress
			}

			if companyInfo.CompanyPhone != nil {
				companyInfoMap["company_phone"] = companyInfo.CompanyPhone
			}

			if companyInfo.IdType != nil {
				companyInfoMap["id_type"] = companyInfo.IdType
			}

			if companyInfo.IdNumber != nil {
				companyInfoMap["id_number"] = companyInfo.IdNumber
			}

			ids = append(ids, helper.Int64ToStr(*companyInfo.CompanyId))
			tmpList = append(tmpList, companyInfoMap)
		}

		_ = d.Set("companies", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
