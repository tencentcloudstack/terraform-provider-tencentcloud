package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeManagerDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeManagerDetailRead,
		Schema: map[string]*schema.Schema{
			"manager_id": {
				Required:    true,
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
				Description: "contact number.",
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
				Description: "Manager&amp;#39;s company information.",
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

			"company_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Manager Company ID.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_manager_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	managerId := d.Get("manager_id")
	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var response *ssl.DescribeManagerDetailResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeManagerDetailByFilter(ctx, helper.IntToStr(managerId.(int)))
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	tmpList := make([]map[string]interface{}, 1)

	if response != nil {
		if response.ManagerId != nil {
			_ = d.Set("manager_id", managerId)
		}

		if response.Status != nil {
			_ = d.Set("status", response.Status)
		}

		if response.ManagerFirstName != nil {
			_ = d.Set("manager_first_name", response.ManagerFirstName)
		}

		if response.ManagerMail != nil {
			_ = d.Set("manager_mail", response.ManagerMail)
		}

		if response.ContactFirstName != nil {
			_ = d.Set("contact_first_name", response.ContactFirstName)
		}

		if response.ManagerLastName != nil {
			_ = d.Set("manager_last_name", response.ManagerLastName)
		}

		if response.ContactPosition != nil {
			_ = d.Set("contact_position", response.ContactPosition)
		}

		if response.ManagerPosition != nil {
			_ = d.Set("manager_position", response.ManagerPosition)
		}

		if response.VerifyTime != nil {
			_ = d.Set("verify_time", response.VerifyTime)
		}

		if response.CreateTime != nil {
			_ = d.Set("create_time", response.CreateTime)
		}

		if response.ExpireTime != nil {
			_ = d.Set("expire_time", response.ExpireTime)
		}

		if response.ContactLastName != nil {
			_ = d.Set("contact_last_name", response.ContactLastName)
		}

		if response.ManagerPhone != nil {
			_ = d.Set("manager_phone", response.ManagerPhone)
		}

		if response.ContactPhone != nil {
			_ = d.Set("contact_phone", response.ContactPhone)
		}

		if response.ContactMail != nil {
			_ = d.Set("contact_mail", response.ContactMail)
		}

		if response.ManagerDepartment != nil {
			_ = d.Set("manager_department", response.ManagerDepartment)
		}

		if response.CompanyInfo != nil {
			companyInfoMap := map[string]interface{}{}

			if response.CompanyInfo.CompanyName != nil {
				companyInfoMap["company_name"] = response.CompanyInfo.CompanyName
			}

			if response.CompanyInfo.CompanyId != nil {
				companyInfoMap["company_id"] = response.CompanyInfo.CompanyId
			}

			if response.CompanyInfo.CompanyCountry != nil {
				companyInfoMap["company_country"] = response.CompanyInfo.CompanyCountry
			}

			if response.CompanyInfo.CompanyProvince != nil {
				companyInfoMap["company_province"] = response.CompanyInfo.CompanyProvince
			}

			if response.CompanyInfo.CompanyCity != nil {
				companyInfoMap["company_city"] = response.CompanyInfo.CompanyCity
			}

			if response.CompanyInfo.CompanyAddress != nil {
				companyInfoMap["company_address"] = response.CompanyInfo.CompanyAddress
			}

			if response.CompanyInfo.CompanyPhone != nil {
				companyInfoMap["company_phone"] = response.CompanyInfo.CompanyPhone
			}

			if response.CompanyInfo.IdType != nil {
				companyInfoMap["id_type"] = response.CompanyInfo.IdType
			}

			if response.CompanyInfo.IdNumber != nil {
				companyInfoMap["id_number"] = response.CompanyInfo.IdNumber
			}
			tmpList = append(tmpList, companyInfoMap)
			_ = d.Set("company_info", []interface{}{companyInfoMap})
		}

		if response.CompanyId != nil {
			_ = d.Set("company_id", response.CompanyId)
		}
	}

	d.SetId(helper.IntToStr(managerId.(int)))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
