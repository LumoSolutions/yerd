package utils

import "fmt"

func SystemdReload() error {
	if output, success := ExecuteCommand("systemctl", "daemon-reload"); !success {
		LogInfo("systemd", "Failed to reload systemd with daemon-reload")
		LogInfo("systemd", "Output: %s", output)
		return fmt.Errorf("unable to reload systemd")
	}

	return nil
}

func SystemdStopService(service string) {
	if output, success := ExecuteCommand("systemctl", "stop", service); !success {
		LogInfo("systemd", "Failed to reload systemd with daemon-reload")
		LogInfo("systemd", "Output: %s", output)
	}
}

func SystemdStartService(service string) error {
	if output, success := ExecuteCommand("systemctl", "start", service); !success {
		LogInfo("systemd", "Failed to reload systemd with daemon-reload")
		LogInfo("systemd", "Output: %s", output)
		return fmt.Errorf("unable to start systemd service %s", service)
	}

	return nil
}

func SystemdServiceActive(service string) bool {
	_, success := ExecuteCommand("systemctl", "is-active", service)
	return success
}
