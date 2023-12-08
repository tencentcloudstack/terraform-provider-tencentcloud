python ./scripts/yunti-process-code.py
go mod vendor
yq eval 'keys' ./scripts/yunti-code.yaml | awk '{print $2}' |grep -v 'go.mod'| xargs -I {} goimports -w {}
yq eval 'keys' ./scripts/yunti-code.yaml | awk '{print $2}' |grep -v 'go.mod'| xargs -I {} go fmt {}
