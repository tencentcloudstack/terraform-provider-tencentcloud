python ./scripts/internal-version-process.py
go mod vendor
yq eval 'keys' ./scripts/internal-version-code.yaml | awk '{print $2}' |grep -v 'go.mod'| xargs -I {} goimports -w {}
yq eval 'keys' ./scripts/internal-version-code.yaml | awk '{print $2}' |grep -v 'go.mod'| xargs -I {} go fmt {}
