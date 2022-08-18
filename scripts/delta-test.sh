#!/bin/bash

range_sha=${BASE_SHA}
echo $(git diff --name-status ${range_sha} | awk '{print $2}')
#service files
update_service_functions=`git diff --name-status ${range_sha} | awk '{print $2}' | grep "^tencentcloud/service*" | xargs git diff ${range_sha} | grep "@@" | grep "func" | awk -F ")" '{print $2}' | awk -F "(" '{print $1}' | tr -d ' '`
echo "update_service_functions: $update_service_functions"
need_test_files=""
for update_service_function in $update_service_functions; do
    tmp_files=`grep -r --with-filename $update_service_function ./tencentcloud | awk -F ":" '{print $1}' | grep -v "service_tencent*" | awk -F "/" '{print $3}' | sort | uniq | egrep "^resource_tc_|^data_source_tc" | awk -F "." '{print $1}' | awk '/_test$/{print "tencentcloud/"$0".go"} !/_test$/{print "tencentcloud/"$0"_test.go"}'`
    need_test_files="$need_test_files $tmp_files"
done
echo "need_test_files: $need_test_files"

# resource&&data_source files
update_sources=`git diff --name-status ${range_sha}| awk '{print $2}' | egrep "^tencentcloud/resource_tc|^tencentcloud/data_source" | egrep -v "_test.go" | awk -F "." '{print $1"_test.go"}'`
echo "update_sources: $update_sources"
# test files
delta_test_files=`git diff --name-status ${range_sha} | egrep "_test\.go$" | awk '{print $2}'`
echo "delta_test_files: $delta_test_files"
# all test files
delta_test_files="$delta_test_files $need_test_files $update_sources"
delta_test_files=`echo $delta_test_files | xargs -n1 | sort | uniq`
echo "all delta_test_files: $delta_test_files"
for delta_test_file in ${delta_test_files}; do
    test_casts=`egrep "func TestAcc.+\(" ${delta_test_file} | awk -F "(" '{print $1}' | awk '{print $2}' | grep -v "NeedFix"`
    echo "[$delta_test_file] \n$test_casts"
    for test_cast in ${test_casts}; do
        go_test_cmd="go test -v -run ${test_cast} -timeout=0 ./tencentcloud/"
        $go_test_cmd
        if [ $? -ne 0 ]; then
            printf "[GO TEST FILED] ${go_test_cmd}"
            exit 1
        fi
    done
done
