/*
Use this data source to query detailed information of ssl describe_manager_detail

Example Usage

```hcl
data "tencentcloud_ssl_describe_manager_detail" "describe_manager_detail" {
                                      }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeManagerDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeManagerDetailRead,
		Schema: map[string]*schema.Schema{
			"manager_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Manager ID.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status: Audit: OK during the review: review passed inValid: expired expiRing: is about to expire Expired: expired.",
			},

			"manager_first_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Manager name.",
			},

			"manager_mail": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Manager mailbox.",
			},

			"contact_first_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Contact name.",
			},

			"manager_last_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Manager name.",
			},

			"contact_position": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Contact position.",
			},

			"manager_position": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Manager position.",
			},

			"verify_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Verify timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"expire_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Verify expiration timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"contact_last_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Contact name.",
			},

			"manager_phone": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Manager phone call.",
			},

			"contact_phone": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Contact number.",
			},

			"contact_mail": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Contact mailbox.",
			},

			"manager_department": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Administrator department.",
			},

			"company_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Manager&amp;amp;#39;s company information.",
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
							Description: "Company phone.",
						},
						"id_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TypeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"id_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID numberNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
					},
				},
			},

			"company_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Manager Company ID.",
			},

			"status_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Review status details.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeManagerDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_manager_detail.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeManagerDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		status = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(status))
	if managerId != nil {
		_ = d.Set("manager_id", managerId)
	}

	if status != nil {
		_ = d.Set("status", status)
	}

	if managerFirstName != nil {
		_ = d.Set("manager_first_name", managerFirstName)
	}

	if managerMail != nil {
		_ = d.Set("manager_mail", managerMail)
	}

	if contactFirstName != nil {
		_ = d.Set("contact_first_name", contactFirstName)
	}

	if managerLastName != nil {
		_ = d.Set("manager_last_name", managerLastName)
	}

	if contactPosition != nil {
		_ = d.Set("contact_position", contactPosition)
	}

	if managerPosition != nil {
		_ = d.Set("manager_position", managerPosition)
	}

	if verifyTime != nil {
		_ = d.Set("verify_time", verifyTime)
	}

	if createTime != nil {
		_ = d.Set("create_time", createTime)
	}

	if expireTime != nil {
		_ = d.Set("expire_time", expireTime)
	}

	if contactLastName != nil {
		_ = d.Set("contact_last_name", contactLastName)
	}

	if managerPhone != nil {
		_ = d.Set("manager_phone", managerPhone)
	}

	if contactPhone != nil {
		_ = d.Set("contact_phone", contactPhone)
	}

	if contactMail != nil {
		_ = d.Set("contact_mail", contactMail)
	}

	if managerDepartment != nil {
		_ = d.Set("manager_department", managerDepartment)
	}

	if companyInfo != nil {
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

		ids = append(ids, *companyInfo.ManagerId)
		_ = d.Set("company_info", companyInfoMap)
	}

	if companyId != nil {
		_ = d.Set("company_id", companyId)
	}

	if statusInfo != nil {
		_ = d.Set("status_info", statusInfo)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
