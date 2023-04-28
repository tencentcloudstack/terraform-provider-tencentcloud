module github.com/tencentcloudstack/terraform-provider-tencentcloud

go 1.13

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/aws/aws-sdk-go v1.31.8
	github.com/bflad/tfproviderlint v0.14.0
	github.com/client9/misspell v0.3.4
	github.com/fatih/color v1.9.0
	github.com/golangci/golangci-lint v1.27.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/hcl/v2 v2.6.0
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/katbyte/terrafmt v0.2.0
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mozillazg/go-httpheader v0.3.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	github.com/tencentcloud/tencentcloud-sdk-go-intl-en v3.0.646+incompatible
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos v1.0.358
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api v1.0.285
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway v1.0.571
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm v1.0.624
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as v1.0.466
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam v1.0.409
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat v1.0.520
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs v1.0.591
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb v1.0.576
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn v1.0.539
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs v1.0.627
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs v1.0.600
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka v1.0.634
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb v1.0.599
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit v1.0.544
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls v1.0.412
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.648
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.624
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp v1.0.589
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb v1.0.572
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu v1.0.335
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain v1.0.634
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb v1.0.572
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.539
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain v1.0.414
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts v1.0.628
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr v1.0.287
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es v1.0.383
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap v1.0.514
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse v1.0.644
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live v1.0.535
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb v1.0.532
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb v1.0.638
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor v1.0.616
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps v1.0.584
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization v1.0.540
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres v1.0.625
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns v1.0.290
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts v1.0.533
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis v1.0.633
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum v1.0.542
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf v1.0.275
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses v1.0.529
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.486
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver v1.0.581
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts v1.0.524
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat v1.0.634
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm v1.0.547
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr v1.0.593
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg v1.0.533
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq v1.0.564
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem v1.0.578
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo v1.0.529
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke v1.0.644
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf v1.0.645
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc v1.0.648
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss v1.0.199
	github.com/tencentyun/cos-go-sdk-v5 v0.7.40
	github.com/yangwenmai/ratelimit v0.0.0-20180104140304-44221c2292e1
	github.com/zclconf/go-cty v1.4.2 // indirect
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	gopkg.in/yaml.v2 v2.2.8
)
