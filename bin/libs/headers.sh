#!/usr/bin/env bash

# [Server Utils-@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# CentOS Packet Assets.
function rcs_packet_assert() {
  crs="$1" && ebc_debug "[系统安装检测($crs)]" && [ "$(yum list installed | grep -c "$crs")" != 1 ] && {
    yum list installed | grep "$crs"
    ebc_error "[软件包未安装($crs)]"
    exit 127
  }

  ebc_success "[软件包已安装($crs)]"
}

# Ubuntu Packet Assets.
function rus_packet_assert() {
  crs="$1" && ebc_debug "[系统安装检测($crs)]" && [ "$(dpkg -l | grep -c "$crs")" != 1 ] && {
    dpkg -l | grep "$crs"
    ebc_error "[软件包未安装($crs)]"
    exit 127
  }

  ebc_success "[软件包已安装($crs)]"
}

# Server Command Assets.
function rs_command_assert() {
  crs="$1" && ebc_debug "[系统安装检测($crs)]" && [ ! -x "$(command -v "$crs")" ] && {
    ebc_error "[软件包未安装($crs)]"
    exit 127
  }

  ebc_success "[软件包已安装($crs)]"
}

# [Colorizer-@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# 通用配色...
function ebc_success() {
  echo -e "\e[1;32m$1\e[0m"
}

function ebc_info() {
  echo -e "\e[1;36m$1\e[0m"
}

function ebc_warn() {
  echo -e "\e[1;33m$1\e[0m"
}

function ebc_error() {
  echo -e "\e[1;31m$1\e[0m"
}

function ebc_debug() {
  echo -e "\e[1;35m$1\e[0m"
}

# [Docker-@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# Docker System Cleaner
function docker_system_cleaner() {
  # Describe Docker System Info
  docker system df

  # Process Docker System Items
  docker system df --format '{{title .Size}}' | while read -r ds_row; do
    echo "Processing Docker System Row Size: $ds_row"

    # check row size
    [[ $ds_row != *"GB"* ]] && {
      echo "TooSmall: $ds_row,Skipping!"
      continue
    }

    # Get Docker System Size Number
    ds_size=$(echo $ds_row | grep -Eo '[+-]?[0-9]+([.][0-9]+)?')
    echo "[AssessDockerSystemSize][ds_row: $ds_row][ds_size: $ds_size]"
    if [ ${ds_size%.*} -ge 10 ]; then
      echo "[DockerSystemSizeTooLarge,Cleaning...][ds_size: $ds_size]"
      # Clearing Docker Data
      docker image prune -f -a
      break
    fi
  done
}

# [PushSvcMs-@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

# push_ms_token
push_ms_token="64e0cef81fdc027a60a34c2fc77aa774b2891f52b89f568d8fb93438f5c1061e"

# Push svc msg
function push_svc_ms() {
  # receive msg
  ns="$1"
  action="$2"
  rs_path="$3"
  branch="$4"
  env="$5"
  tag="$6"
  params="$7"

  echo "[推送参数][ns:$ns][action:$action][rs_path:$rs_path][branch:$branch][env:$env][tag:$tag][params:$params]"

  # 依赖 jq 软件
  rs_command_assert "jq" || sudo yum install -y jq

  # [Modify a key-value in a json using jq in-place](https://stackoverflow.com/questions/42716734/modify-a-key-value-in-a-json-using-jq-in-place)
  push="$(jq '.text.content = "[服务任务执行完成][ns:'$ns'][action:'$action'][rs_path:'$rs_path'][branch:'$branch'][env:'$env'][tag:'$tag'][host:'$(hostname)']"' $params)"

  # push msg
  curl 'https://oapi.dingtalk.com/robot/send?access_token='$push_ms_token -H 'Content-Type: application/json' -d "$push"
}

# @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
