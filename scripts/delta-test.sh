#!/bin/bash

range_sha=${BASE_SHA}
pr_id=${PR_ID}

update_source_count=`git diff --name-status ${range_sha}| awk '{print $2}' | egrep "^tencentcloud/resource_tc|^tencentcloud/data_source" | egrep -v "_test.go" | wc -l`
if [ $update_source_count -eq 0 ]; then
    printf "No source change, skip delta-test!"
    exit 0
fi

if [ ! -f ".changelog/${pr_id}.txt" ]; then
    printf "Not find changelog file!"
    exit 1
fi
source_names=`cat .changelog/${pr_id}.txt| grep -E "^(resource|datasource)\/(\w+)" | awk -F ":" '{print $1}' | sort | uniq`

test_files=""
for source_name in $source_names; do
    name=${source_name#*/}
    type=${source_name%/*}
    if [ $type == "datasource" ]; then
        type=dataSource
    fi 
    # echo $source_name $type $name
    function_name=$(cat tencentcloud/provider.go | grep "\"${name}\"" | grep "${type}")
    function_name=${function_name#*:}
    function_name=$(echo $(echo ${function_name%,*}))

    test_file=$(grep -r "func $function_name \*schema\.Resource" tencentcloud)
    test_file=${test_file#*/}
    test_file=${test_file%:*}
    test_files="$test_files $test_file"
done
echo "test files:" $test_files

for test_file in $test_files; do
    test_case_type=${test_file%_tc*}
    test_case_name=${test_file#*tc_}
    test_case_name=${test_case_name%.*}

    test_case_type=`echo $test_case_type | sed -r 's/(^|_)(\w)/\U\2/g'`
    test_case_name=`echo $test_case_name | sed -r 's/(^|_)(\w)/\U\2/g'`
   
    go_test_cmd="go test -v -run TestAccTencentCloud${test_case_name}${test_case_type} -timeout=0 ./tencentcloud/"
    echo $go_test_cmd
    $go_test_cmd
    if [ $? -ne 0 ]; then
        printf "[GO TEST FILED] ${go_test_cmd}"
        exit 1
    fi
done
