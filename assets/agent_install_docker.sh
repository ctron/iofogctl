#!/bin/sh
set -x
set -e

start_docker() {
	set +e
	# check if docker is running
	if [ ! -f /var/run/docker.pid ]; then
		$sh_c "/etc/init.d/docker start"
		local err_code=$?
		if [ $err_code -ne 0 ]; then
			$sh_c "service docker start"
			err_code=$?
		fi
		if [ $err_code -ne 0 ]; then
			echo "Could not start Docker daemon"
			exit 1
		fi
	fi
	set -e
}

do_configure_overlay() {
	echo "# Configuring /etc/systemd/system/docker.service.d/overlay.conf..."
	if [ "$lsb_dist" = "raspbian" ] || [ "$(uname -m)" = "armv7l" ] || [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "armv8" ]; then
		if [ ! -d "/etc/systemd/system/docker.service.d" ]; then
			$sh_c "mkdir -p /etc/systemd/system/docker.service.d"
		fi
		if [ -f "/etc/systemd/system/docker.service.d/overlay.conf" ] && ! grep -Fxq "ExecStart=/usr/bin/dockerd --storage-driver overlay -H unix:// -H tcp://127.0.0.1:2375" "/etc/systemd/system/docker.service.d/overlay.conf"; then
			$sh_c 'echo "ExecStart=/usr/bin/dockerd --storage-driver overlay -H unix:// -H tcp://127.0.0.1:2375" >> /etc/systemd/system/docker.service.d/overlay.conf'
		elif [ ! -f "/etc/systemd/system/docker.service.d/overlay.conf" ]; then
			$sh_c 'echo "[Service]" > /etc/systemd/system/docker.service.d/overlay.conf'
			$sh_c 'echo "ExecStart=" >> /etc/systemd/system/docker.service.d/overlay.conf'
			$sh_c 'echo "ExecStart=/usr/bin/dockerd --storage-driver overlay -H unix:// -H tcp://127.0.0.1:2375" >> /etc/systemd/system/docker.service.d/overlay.conf'
		fi
		$sh_c "systemctl daemon-reload"
		$sh_c "service docker restart"
	fi
}

do_install_docker() {
	# Check that Docker 18.09.2 or greater is installed
	if command_exists docker; then
		docker_version=$(docker -v | sed 's/.*version \(.*\),.*/\1/' | tr -d '.')
		if [ "$docker_version" -ge 18090 ]; then
			echo "# Docker $docker_version already installed"
			start_docker
			do_configure_overlay
			return
		fi
	fi
	echo "# Installing Docker..."
	curl -fsSL https://get.docker.com/ | sh
	
	if ! command_exists docker; then
		echo "Failed to install Docker"
		exit 1
	fi
	start_docker
	do_configure_overlay
}

. /tmp/agent_init.sh
init
do_install_docker