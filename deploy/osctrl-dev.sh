#!/bin/sh
#
# To run this script execute:
#   `curl -ks https:/__TLSHOST/dev/__SECRETMD5/osctrl.sh | sh`
#
# IMPORTANT! If osquery is not installed, it will be installed.

_PROJECT="osctrl"
_SECRET="__OSQUERYSECRET"
_SECRET_LINUX=/etc/osquery/osquery.secret
_FLAGS_LINUX=/etc/osquery/osquery.flags
_CERT_LINUX=/etc/osquery/certs/osctrl.crt
_SECRET_OSX=/private/var/osquery/osquery.secret
_FLAGS_OSX=/private/var/osquery/osquery.flags
_CERT_OSX=/private/var/osquery/certs/osctrl.crt
_OSQUERY_PKG="https://osquery-packages.s3.amazonaws.com/darwin/osquery-3.3.0.pkg"
_OSQUERY_DEB="https://osquery-packages.s3.amazonaws.com/deb/osquery_3.3.0_1.linux.amd64.deb"
_OSQUERY_RPM="https://osquery-packages.s3.amazonaws.com/rpm/osquery-3.3.0-1.linux.x86_64.rpm"

fail() {
  echo "[!] $1"
  exit 1
}

log() {
  echo "[+] $1"
}

installOsquery() {
  log "Installing osquery for $OS"
  if [ "$OS" = "linux" ]; then
    distro=$(/usr/bin/rpm -q -f /usr/bin/rpm >/dev/null 2>&1)
    if [ "$?" = "0" ]; then
      log "RPM based system detected"
      _RPM="$(echo $_OSQUERY_RPM | cut -d"/" -f5)"
      sudo curl -# "$_OSQUERY_RPM" -o "/tmp/$_RPM"
      sudo rpm -ivh "/tmp/$_RPM"
    else
      log "DEB based system detected"
      _DEB="$(echo $_OSQUERY_DEB | cut -d"/" -f5)"
      sudo curl -# "$_OSQUERY_DEB" -o "/tmp/$_DEB"
      sudo dpkg -i "/tmp/$_DEB"
    fi
  fi
  if [ "$OS" = "darwin" ]; then
    _PKG="$(echo $_OSQUERY_PKG | cut -d"/" -f5)"
    sudo curl -# "$_OSQUERY_PKG" -o "/tmp/$_PKG"
    sudo installer -pkg "/tmp/$_PKG" -target /
  fi
}

verifyOsquery() {
  osqueryi=$(which osqueryi)
  if [ "$?" = "1" ]; then
    #read -p "[+] $_PROJECT needs osquery. Do you want to install it? [y/n]" yn
    #case $yn in
    #  [Yy]* ) installOsquery;;
    #  [Nn]* ) exit 1;;
    #  * ) exit 1;;
    #esac
    log "[+] $_PROJECT needs osquery"
    installOsquery
  else
    osqueryi -version
  fi
}

whatOS() {
	OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')
  log "OS=$OS"
  if [ "$OS" = "linux" ]; then
    _SECRET_FILE="$_SECRET_LINUX"
    _FLAGS="$_FLAGS_LINUX"
    _CERT="$_CERT_LINUX"
  fi
  if [ "$OS" = "darwin" ]; then
    _SECRET_FILE="$_SECRET_OSX"
    _FLAGS="$_FLAGS_OSX"
    _CERT="$_CERT_OSX"
  fi
  log "_SECRET_FILE=$_SECRET_FILE"
  log "_FLAGS=$_FLAGS"
  log "_CERT=$_CERT"
}

stopOsquery() {
  log "Stopping osqueryd"
  if [ "$OS" = "linux" ]; then
    sudo systemctl stop osqueryd
  fi
  if [ "$OS" = "darwin" ]; then
    if launchctl list | grep -qcm1 com.facebook.osqueryd; then
      sudo launchctl unload /Library/LaunchDaemons/com.facebook.osqueryd.plist
    fi
  fi
}

prepareSecret() {
  log "Preparing osquery secret"
  echo "$_SECRET" | sudo tee "$_SECRET_FILE"
  sudo chmod 700 "$_SECRET_FILE"
}

prepareFlags() {
  log "Preparing osquery flags"
  sudo sh -c "cat <<EOF > $_FLAGS
--host_identifier=uuid
--force=true
--verbose=true
--debug
--utc
--pidfile=/tmp/osquery.pid
--database_path=/tmp/osquery.db
--enroll_secret_path=$_SECRET_FILE
--enroll_tls_endpoint=/dev/osquery_enroll
--config_plugin=tls
--config_tls_endpoint=/dev/osquery_config
--config_tls_refresh=10
--logger_plugin=tls
--logger_tls_compress=false
--logger_tls_endpoint=/dev/osquery_log
--logger_tls_period=10
--disable_distributed=false
--distributed_interval=10
--distributed_plugin=tls
--distributed_tls_read_endpoint=/dev/osquery_read
--distributed_tls_write_endpoint=/dev/osquery_write
--tls_dump=true
--tls_hostname=__TLSHOST
--tls_server_certs=$_CERT
EOF"
}

prepareCert() {
  log "Preparing osquery certificate"
  sudo mkdir -p $(dirname "$_CERT")
  sudo sh -c "cat <<EOF > $_CERT
__CERT_CONTENT
EOF"
}

startOsquery() {
  log "Starting osqueryd"
  if [ "$OS" = "linux" ]; then
    sudo systemctl start osqueryd
    sudo systemctl enable osqueryd
  fi
  if [ "$OS" = "darwin" ]; then
    sudo cp /private/var/osquery/com.facebook.osqueryd.plist /Library/LaunchDaemons/com.facebook.osqueryd.plist
    sudo launchctl load /Library/LaunchDaemons/com.facebook.osqueryd.plist
  fi
}

bye() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "[!] Fail to enroll $_PROJECT node"
  fi
  exit $result
}

trap "bye" EXIT
whatOS
verifyOsquery
set -e
stopOsquery
prepareSecret
prepareFlags
prepareCert
startOsquery

log "Congratulations! The node has been enrolled in $_PROJECT"

# EOF