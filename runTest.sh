#!/bin/bash

echo "authentication tests are running................."
test1=$(go test authentication/*.go | awk '{print$1}')
echo $test1 |  sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
result1=$(echo ${test1: -2})
echo $result1


echo "handlers tests are running................."
test2=$(go test -v handler/*.go | awk '{print$1}')
echo $test2 |  sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
result2=$(echo ${test2: -2})
echo $result2

echo "util tests are running................."
test3=$(go test -v util/*.go | awk '{print$1}')
echo $test3 |  sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
result3=$(echo ${test3: -2})
echo $result3


if [ \( "$result1" = "ok" -a "$result2" = "ok" -a "$result3" = "ok" \) ]
then
    echo "pass"
else
    echo "fail"
fi


