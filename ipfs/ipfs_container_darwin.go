//go:build darwin

package ipfs

import (
	"fmt"
	"io"
	"os/exec"
)

const (
	ctlCmdName    = "lima"
	ctlSubCmdName = "nerdctl"
)

type cmdRunner interface {
	fmt.Stringer

	Run() error
}

func pullCmd(stdout io.Writer, stderr io.Writer) cmdRunner {

	cmd := exec.Command(ctlCmdName,
		ctlSubCmdName,
		cmdPull,
		imageName)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd

}

func psCmd(stdout io.Writer, stderr io.Writer) *exec.Cmd {

	cmd := exec.Command(ctlCmdName,
		ctlSubCmdName,
		cmdPs)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd
}

func runCmd(runConfig RunConfig) *exec.Cmd {

	return exec.Command(ctlCmdName, ctlSubCmdName,
		cmdRun,
		"-d", // run container in background and print container id as first line of stdout
		"-p", // publish container's ports to host o/s
		fmt.Sprintf("%s:%d:%d", runConfig.Host, runConfig.Port, runConfig.Port),
		imageName)
}

func startCmd(containerId string) *exec.Cmd {

	return exec.Command(ctlCmdName,
		ctlSubCmdName,
		cmdStart,
		containerId)
}

func stopCmd(containerId string) *exec.Cmd {

	return exec.Command(ctlCmdName,
		ctlSubCmdName,
		cmdStart,
		containerId)
}
