#!/bin/sh

log_info() {
  echo "[info] $1"
}

log_err() {
  echo "[error] $1"
  exit 1
}

main() {
  for platform in $PLATFORMS
  do
    version=$(echo "${VERSION}" | sed 's/[-.+]/_/g')
    read -r os arch \
      <<<$(echo "${platform}" | tr '/' ' ')

    filename="keygen_${os}_${arch}_${version}"
    if [ "${os}" = 'windows' ]
    then
      filename="${filename}.exe"
    fi

    log_info "publishing v${VERSION} for ${platform}: cli/${filename}"

    keygen dist "build/${filename}" \
      --filename "cli/${filename}" \
      --name "CLI v${VERSION}" \
      --platform "${platform}" \
      --channel "${CHANNEL}" \
      --version "${VERSION}" \
      --no-auto-upgrade

    if [ $? -eq 0 ]
    then
      log_info "successfully published v${VERSION} for ${platform}"
    else
      log_err "failed to publish v${VERSION} for ${platform}"
    fi
  done

  # We only want to update these releases for stable releases
  if [ "${CHANNEL}" = 'stable' ]
  then
    keygen dist "build/install.sh" \
      --filename "cli/install.sh" \
      --name "CLI Installer" \
      --platform '*' \
      --version "${VERSION}" \
      --no-auto-upgrade

    keygen dist "build/version" \
      --filename "cli/version" \
      --filetype 'txt' \
      --name "CLI Version" \
      --platform '*' \
      --version "${VERSION}" \
      --no-auto-upgrade
  fi
}

PLATFORMS=$(go tool dist list | grep -vE 'ios|android|js|aix|illumos|riscv64|plan9|solaris')
VERSION=$(cat VERSION)
CHANNEL='stable'

case "${VERSION}"
in
  *-rc.*)
    CHANNEL='rc'
    ;;
  *-beta.*)
    CHANNEL='beta'
    ;;
  *-alpha.*)
    CHANNEL='alpha'
    ;;
  *-dev.*)
    CHANNEL='dev'
    ;;
esac

main
