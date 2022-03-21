module github.com/tencentcloudstack/terraform-provider-tencentcloud

go 1.13

require (
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/storage v1.21.0 // indirect
	github.com/Azure/azure-sdk-for-go v50.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v12 v12.0.0 // indirect
	github.com/aws/aws-sdk-go v1.43.21
	github.com/bflad/tfproviderlint v0.14.0
	github.com/bketelsen/crypt v0.0.4 // indirect
	github.com/client9/misspell v0.3.4
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.7.0 // indirect
	github.com/docker/cli v20.10.7+incompatible // indirect
	github.com/docker/docker v20.10.7+incompatible // indirect
	github.com/fatih/color v1.9.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/golangci/golangci-lint v1.27.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.2.0 // indirect
	github.com/googleapis/google-cloud-go-testing v0.0.0-20200911160855-bcd43fbb19e8 // indirect
	github.com/gruntwork-io/terratest v0.37.13
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-getter v1.5.11 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/hashicorp/terraform-json v0.13.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/hashicorp/terraform-plugin-test v1.3.0 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/jstemmer/go-junit-report v1.0.0 // indirect
	github.com/katbyte/terrafmt v0.2.0
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mozillazg/go-httpheader v0.3.0 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.6.2 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos v1.0.358
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api v1.0.285
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as v1.0.363
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam v1.0.357
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka v1.0.310
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb v1.0.283
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls v1.0.291
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.363
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.351
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb v1.0.359
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu v1.0.335
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.294
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr v1.0.287
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor v1.0.329
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres v1.0.332
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns v1.0.290
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf v1.0.275
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr v1.0.267
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq v1.0.268
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke v1.0.302
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc v1.0.357
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss v1.0.199
	github.com/tencentyun/cos-go-sdk-v5 v0.7.33
	github.com/tmccombs/hcl2json v0.3.4 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/yangwenmai/ratelimit v0.0.0-20180104140304-44221c2292e1
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/tools v0.1.10 // indirect
	google.golang.org/api v0.73.0 // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)
