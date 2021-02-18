#!/usr/bin/env bash
#set -xe

# [How can I repeat a character in Bash?](https://stackoverflow.com/questions/5349718/how-can-i-repeat-a-character-in-bash)
function println() {
  # print titles
  PTS=${1:-"@"}
  # print chars
  PCS=${2:-"#"}

  #  printf "%s %0.s$PCS" "[$PTS]" {1..100}

  PTS="[$PTS]"
  left=$(($(tput cols) - ${#PTS}))
  printf "\e[1;36m$PTS\e[0m\e[1;35m"
  printf "$PCS%.0s" $(seq 1 $left)
  echo -e "\e[0m"
}

# protocol
grpcs="grpcurl -plaintext"
#grpcs="grpcurl -insecure"
#grpcs="grpcurl -cacert=/Users/coam/Run/Test/grpcs/tlscert/cert.crt"
#grpcs="grpcurl -cert=/Users/coam/Run/Test/grpcs/tlscert/cert.crt -key=/Users/coam/Run/Test/grpcs/tlscert/private.key"

#random=$(((RANDOM % 10000) + 1))
#println "random: ${random}"

println "GRpc Service List"

$grpcs -H 'trace: 111111111111111111' localhost:8080 list

println "GRpc Service Call"

$grpcs -H 'user: zhangyafei' -H 'pass: 123456789' -H 'trace: 222222222222222222222222' -vv -d @ localhost:8080 helloworld.Greeter/SayHello <<EOM
{
  "name": "zhangsan"
}
EOM

println "GRpc Service Describe"

$grpcs -H 'trace: 33333333333333333333333' localhost:8080 describe helloworld.Greeter.SayHello
$grpcs -H 'trace: 44444444444444444444444' localhost:8080 describe helloworld.Greeter.SayList

println "Server-Client"

# server
go run src/rpc/echo/server/main.go

# client
go run src/rpc/echo/client/main.go

println "Ended"
