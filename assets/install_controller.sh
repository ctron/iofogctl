#!/bin/sh
set -x
set -e

INSTALL_DIR="/opt/iofog"
TMP_DIR="/tmp/iofog"

load_existing_nvm() {
	set +e
	if [ -z "$(command -v nvm)" ]; then
		export NVM_DIR="${HOME}/.nvm"
		mkdir -p $NVM_DIR
		if [ -f "$NVM_DIR/nvm.sh" ]; then
			[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
		fi
	fi
	set -e
}

controller_service() {
    USE_SYSTEMD=`grep -m1 -c systemd /proc/1/comm`
    USE_INITCTL=`which initctl | wc -l`
    USE_SERVICE=`which service | wc -l`

    if [ $USE_SYSTEMD -eq 1 ]; then
        cp /tmp/iofog-controller-service/iofog-controller.systemd /etc/systemd/system/iofog-controller.service
        chmod 644 /etc/systemd/system/iofog-controller.service
        systemctl daemon-reload
        systemctl enable iofog-controller.service
    elif [ $USE_INITCTL -eq 1 ]; then
        cp /tmp/iofog-controller-service/iofog-controller.initctl /etc/init/iofog-controller.conf
        initctl reload-configuration
    elif [ $USE_SERVICE -eq 1 ]; then
        cp /tmp/iofog-controller-service/iofog-controller.update-rc /etc/init.d/iofog-controller
        chmod +x /etc/init.d/iofog-controller
        update-rc.d iofog-controller defaults
    else
        echo "Unable to setup Controller startup script."
    fi
}

install_deps() {
	if [ -z "$(command -v lsof)" ]; then
		if [ -z "$(command -v apt)" ]; then
			echo "Unsupported distro"
			exit 1
		fi
		apt update -qq
		apt install -y lsof
	fi
}

deploy_controller() {
	# Nuke any existing instances
	if [ ! -z $(lsof -ti tcp:51121) ]; then
		lsof -ti tcp:51121 | xargs kill
	fi

	# nvm
	load_existing_nvm
	if [ -z "$(command -v nvm)" ]; then
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.34.0/install.sh | bash
		export NVM_DIR="${HOME}/.nvm"
		[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
	fi
	nvm install lts/dubnium
	nvm use lts/dubnium
	ln -Ffs $(which node) /usr/local/bin/node

	# npmrc
	if [ -z "$(command -v npmrc)" ]; then
		npm i npmrc -g
	fi

	# If token is provided, set up private repo
	if [ ! -z $token ]; then
		if [ ! -z $(npmrc | grep iofog)]; then
			npmrc -c iofog
			npmrc iofog
		fi
		curl -s https://"$token":@packagecloud.io/install/repositories/"$repo"/script.node.sh | force_npm=1 bash
		mv ~/.npmrc ~/.npmrcs/npmrc
		ln -s ~/.npmrcs/npmrc ~/.npmrc
	else
		npmrc default
	fi

	# Install in temporary location
	mkdir -p "$TMP_DIR/controller"
	chmod 0777 "$TMP_DIR/controller"
	if [ -z $version ]; then
		npm install -g -f minipass@2.7.0 iofogcontroller --unsafe-perm --prefix "$TMP_DIR/controller"
	else
		npm install -g -f minipass@2.7.0 "iofogcontroller@$version" --unsafe-perm --prefix "$TMP_DIR/controller"
	fi

	# Move files into $INSTALL_DIR/controller
	mkdir -p "$INSTALL_DIR/"
	rm -rf "$INSTALL_DIR/controller" # Clean possible previous install
	mv "$TMP_DIR/controller/" "$INSTALL_DIR/"

	# Symbolic links
	if [ ! -f "/usr/local/bin/iofog-controller" ]; then
		ln -fFs "$INSTALL_DIR/controller/bin/iofog-controller" /usr/local/bin/iofog-controller
	fi

    # Set controller permissions
    chmod 744 -R "$INSTALL_DIR/controller"

    # Startup script
    controller_service

    # Run controller
    iofog-controller start
}

# main
version="$1"
repo=$([ -z "$2" ] && echo "iofog/iofog-controller-snapshots" || echo "$2")
token="$3"
# Optional args
export DB_PROVIDER="$4"
export DB_HOST="$5"
export DB_USER="$6"
export DB_PASSWORD="$7"
export DB_PORT="$8"

install_deps
deploy_controller