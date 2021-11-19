#!/bin/sh

log_info() {
  echo "info: ${1}"
}

log_err() {
  echo "error: ${1} (please install manually via https://keygen.sh/docs/cli/)"
  exit 1
}

get_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')

  case "${os}"
  in
  msys*)
    os='windows'
    ;;
  cygwin*)
    os='linux'
    ;;
  esac

  if [ -z "${os}" ]
  then
    log_fatal 'unable to detect operating system'
  fi

  echo "${os}"
}

get_arch() {
  arch=$(uname -m)

  case "${arch}"
  in
  x86_64|amd64p32)
    arch='amd64'
    ;;
  x86)
    arch='386'
    ;;
  aarch64)
    arch='arm64'
    ;;
  armv8)
    arch='arm64'
    ;;
  armv*)
    arch='arm'
    ;;
  i686|i386)
    arch='386'
    ;;
  esac

  if [ -z "${arch}" ]
  then
    log_fatal 'unable to detect architecture'
  fi

  echo "${arch}"
}

get_bin_version() {
  version=$(curl -sSL 'https://bin.keygen.sh/keygen/cli/version')

  echo "${version}"
}

get_bin_path() {
  echo '/usr/local/bin/keygen'
}

get_bin_url() {
  version=$(echo "${BIN_VERSION}" | sed 's/[-.+]/_/g')

  filename="keygen_${OS}_${ARCH}_${version}"
  if [ "${os}" = 'windows' ]
  then
    filename="${filename}.exe"
  fi

  echo "https://bin.keygen.sh/keygen/cli/${filename}"
}

assert_os_support() {
  case "${OS}"
  in
    darwin) return ;;
    dragonfly) return ;;
    freebsd) return ;;
    linux) return ;;
    netbsd) return ;;
    openbsd) return ;;
    windows) return ;;
  esac

  log_err "unsupported operating system: ${OS}"
}

assert_arch_support() {
  case "${ARCH}"
  in
    386) return ;;
    amd64) return ;;
    arm64) return ;;
    arm) return ;;
    ppc64) return ;;
    ppc64le) return ;;
    mips) return ;;
    mipsle) return ;;
    mips64) return ;;
    mips64le) return ;;
    s390x) return ;;
  esac

  log_err "unsupported architecture: ${ARCH}"
}

assert_platform_support() {
  assert_os_support
  assert_arch_support
}

main() {
  assert_platform_support

  status=$(curl -sSL "${BIN_URL}" --write-out "%{http_code}" -o keygen)

  if [ "${status}" -eq 200 ]
  then
    log_info "successfully downloaded v${BIN_VERSION} for ${PLATFORM}"
  else
    log_err "failed to download v${BIN_VERSION} for ${PLATFORM}"
  fi

  mv keygen "${BIN_PATH}" && \
    chmod +x "${BIN_PATH}"

  if [ $? -eq 0 ]
  then
    log_info "successfully installed v${BIN_VERSION} for ${PLATFORM}"
  else
    log_err "failed to installed v${BIN_VERSION} for ${PLATFORM}"
  fi

  ${BIN_PATH} --help
}

OS=$(get_os)
ARCH=$(get_arch)
PLATFORM="${OS}/${ARCH}"
BIN_VERSION=$(get_bin_version)
BIN_PATH=$(get_bin_path)
BIN_URL=$(get_bin_url)

main
