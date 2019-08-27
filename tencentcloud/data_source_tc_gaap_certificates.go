package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func dataSourceTencentCloudGaapCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCertificatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 4),
			},

			// computed
			"certificates": {
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
						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"issuer_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_certificates.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var (
		id              *string
		name            *string
		certificateType *int
		ids             []string
		certificates    []map[string]interface{}
	)

	if raw, ok := d.GetOk("id"); ok {
		id = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("type"); ok {
		certificateType = common.IntPtr(raw.(int))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	respCertificates, err := service.DescribeCertificates(ctx, id, name, certificateType)
	if err != nil {
		return err
	}

	ids = make([]string, 0, len(respCertificates))
	certificates = make([]map[string]interface{}, 0, len(respCertificates))
	for _, certificate := range respCertificates {
		ids = append(ids, *certificate.CertificateId)

		m := map[string]interface{}{
			"id":          *certificate.CertificateId,
			"name":        *certificate.CertificateAlias,
			"type":        *certificate.CertificateType,
			"create_time": *certificate.CreateTime,
		}

		if certificate.BeginTime != nil {
			m["begin_time"] = *certificate.BeginTime
		}

		if certificate.EndTime != nil {
			m["end_time"] = *certificate.EndTime
		}

		if certificate.IssuerCN != nil {
			m["issuer_cn"] = *certificate.IssuerCN
		}

		if certificate.SubjectCN != nil {
			m["subject_cn"] = *certificate.SubjectCN
		}

		certificates = append(certificates, m)
	}

	d.Set("certificates", certificates)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
