#!/usr/bin/env bash

set -e

HUGE_PAGES=
HUGE_PAGES_OVERRIDE=
NVME_TCP=
INSTALL_ZFS=
INSTALL_LVM=
DRY_RUN=
SYSCTL="sudo sysctl"
MODPROBE="sudo modprobe"

help() {
  cat <<EOF
Usage: $(basename "$0") [COMMAND] [OPTIONS]

Options:
  -h, --help                            Display this text.
  --hugepages         <num>             Add <num> 2MiB hugepages.
  --nvme-tcp                            Load nvme_tcp kernel modules.
  --install-zfs                         Install ZFS utilities.
  --install-lvm                         Install LVM utilities.

Examples:
  $(basename "$0") --nvme-tcp --hugepages 2048 --install-zfs
EOF
}

echo_stderr() {
  echo -e "${1}" >&2
}

die() {
  local _return="${2:-1}"
  echo_stderr "$1"
  exit "${_return}"
}

setup_hugepages() {
  $SYSCTL -w vm.nr_hugepages="$1"
}

modprobe_nvme_tcp() {
  $MODPROBE nvme_tcp
}

nvme_ana_check() {
  cat /sys/module/nvme_core/parameters/multipath
}

distro() {
  cat /etc/os-release | awk -F= '/^NAME=/ {print $2}' | tr -d '"'
}

install_kernel_modules_nsup() {
  die "Installing kernel modules not supported for $1"
}

install_kernel_modules() {
  DISTRO="$(distro)"
  case "$DISTRO" in
    Ubuntu)
      sudo apt-get update
      sudo apt-get install -y linux-modules-extra-$(uname -r) xfsprogs quota
      ;;
    NixOS | *)
      install_kernel_modules_nsup "$DISTRO"
      ;;
  esac
}

install_zfs() {
  if ! command -v zfs &> /dev/null; then
    DISTRO="$(distro)"
    case "$DISTRO" in
      Ubuntu)
        sudo apt-get update
        if sudo apt-get install -y zfsutils-linux; then
          echo "Successfully installed zfsutils-linux"
        else
          die "Failed to install zfsutils-linux"
        fi
        ;;
      NixOS | *)
        die "Installation of zfsutils-linux not supported for $DISTRO"
        ;;
    esac
  else
    echo "ZFS utilities are already installed"
  fi
}

install_lvm() {
  if ! command -v lvm &> /dev/null; then
    DISTRO="$(distro)"
    case "$DISTRO" in
      Ubuntu)
        sudo apt-get update
        if sudo apt-get install -y lvm2; then
          echo "Successfully installed lvm2"
        else
          die "Failed to install lvm2"
        fi
        ;;
      NixOS | *)
        die "Installation of lvm2 not supported for $DISTRO"
        ;;
    esac
  else
    echo "LVM utilities are already installed"
  fi
}

load_lvm_modules() {
  # Load LVM snapshot and thin provisioning modules
  sudo modprobe dm-snapshot
  sudo modprobe dm-thin-pool
}

while [ "$#" -gt 0 ]; do
  case $1 in
    -h|--help)
      help
      exit 0
      ;;
    --hugepages)
      shift
      test $# -lt 1 && die "Missing hugepage number"
      HUGE_PAGES=$1
      shift
      ;;
    --hugepages-override)
      shift
      test $# -lt 1 && die "Missing hugepage number"
      HUGE_PAGES_OVERRIDE="y"
      HUGE_PAGES=$1
      shift
      ;;
    --nvme-tcp)
      NVME_TCP="y"
      shift
      ;;
    --install-zfs)
      INSTALL_ZFS="y"
      shift
      ;;
    --install-lvm)
      INSTALL_LVM="y"
      shift
      ;;
    --dry-run)
      if [ -z "$DRY_RUN" ]; then
        DRY_RUN="--dry-run"
        SYSCTL="echo $SYSCTL"
        MODPROBE="echo $MODPROBE"
      fi
      shift
      ;;
    *)
      die "Unknown argument $1!"
      ;;
  esac
done

if [ -n "$HUGE_PAGES" ]; then
  pages=$(sysctl -b vm.nr_hugepages)

  if [ "$HUGE_PAGES" -gt "$pages" ]; then
    setup_hugepages "$HUGE_PAGES"
  else
    if [ "$HUGE_PAGES" -lt "$pages" ] && [ -n "$HUGE_PAGES_OVERRIDE" ]; then
      echo "Overriding hugepages from $pages to $HUGE_PAGES, as requested"
      setup_hugepages "$HUGE_PAGES"
    else
      echo "Current hugepages ($pages) are sufficient"
    fi
  fi
fi

if [ -n "$NVME_TCP" ]; then
  if ! lsmod | grep "nvme_tcp" >/dev/null; then
    if ! modprobe_nvme_tcp >/dev/null; then
      install_kernel_modules
      if ! modprobe_nvme_tcp; then
        die "Failed to load nvme_tcp kernel module!"
      fi
    fi
    echo "Installed nvme_tcp kernel module"
  else
    echo "nvme-tcp kernel module already installed"
  fi

  if [ "$(nvme_ana_check)" != "Y" ]; then
    echo_stderr "NVMe multipath support is NOT enabled!"
  else
    echo "NVMe multipath support IS enabled"
  fi
fi

if [ -n "$INSTALL_ZFS" ]; then
  install_zfs
fi

if [ -n "$INSTALL_LVM" ]; then
  install_lvm
  load_lvm_modules
fi
