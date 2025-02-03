#!/bin/bash
print_black() {
    echo -e "\033[30m$1\033[0m"
}

print_red() {
    echo -e "\033[31m$1\033[0m"
}

print_green() {
    echo -e "\033[32m$1\033[0m"
}

print_yellow() {
    echo -e "\033[33m$1\033[0m"
}

print_blue() {
    echo -e "\033[34m$1\033[0m"
}

print_magenta() {
    echo -e "\033[35m$1\033[0m"
}

print_cyan() {
    echo -e "\033[36m$1\033[0m"
}

print_grey() {
    echo -e "\033[37m$1\033[0m"
}

print_white() {
    echo "$1"
}

PROGRAM="sleepy-dashboard"

config="/data/config.json"

# Read parameters
{
    # Parse parameters
    while [ $# -gt 0 ]; do
        case $1 in
        --config)
            config=$2
            shift
            ;;
        *)
            # Skip other flags
            ;;
        esac
        shift
    done
}

# Check install
[ -f "$config" ] && exit 0

print_cyan "$PROGRAM installation script"

print_yellow " ** Create dir..."

mkdir -p /data

print_yellow " ** Link custom..."

mv /opt/$PROGRAM/custom /data/custom
ln -s /data/custom /opt/$PROGRAM/custom

print_yellow " ** Copy config..."

if [ "$USE_TLS" == "1" ]; then
    cp -f /opt/$PROGRAM/examples/docker_https.json "$config"
else
    cp -f /opt/$PROGRAM/examples/docker_http.json "$config"
fi

print_yellow " ** Setting up..."

[ -n "$SECRET" ] && sed -i "s/PUT_YOUR_AGENT_SECRET_HERE/$SECRET/g" "$config"
[ -n "$SITE_CONFIG" ] && echo "$SITE_CONFIG" | base64 -d > /data/custom/site_config.json

if [ "$USE_TLS" == "1" ]; then
    print_yellow " ** Generate website ECC certificate and private key..."
    openssl ecparam -genkey -name prime256v1 -out /data/private.key
    openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes -keyout /data/private.key -x509 -days 36500 -subj "/" -out /data/cert.pem
fi

print_green "$PROGRAM installed successfully"
