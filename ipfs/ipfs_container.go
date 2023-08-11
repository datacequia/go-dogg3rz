package ipfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"strings"
	"time"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
)

const (
	imageName      = "ipfs/kubo:v0.22.0"
	defaultTimeout = time.Second * 5
	cmdPull        = "pull"
	cmdRun         = "run"
	cmdStart       = "start"
	cmdStop        = "stop"
	cmdPs          = "ps"
)

type tempDirLocationFunc = func() string
type createTempDirFunc = func(dir string, pattern string) (*os.File, error)

//type hasLineFunc = func(r io.Reader) bool

type ipfsContainer struct {
	_tempDirLocation tempDirLocationFunc
	_createTempDir   createTempDirFunc

	_cmdRunner cmdRunner
	_cmdName   string
}

func newIpfsContainer() ipfsContainer {

	return ipfsContainer{
		_tempDirLocation: os.TempDir,
		_createTempDir:   os.CreateTemp,
	}

}

func (o ipfsContainer) cmdName(cmdName string) ipfsContainer {

	o._cmdName = cmdName

	return o
}

func (o ipfsContainer) cmdRunner(cmdRunner cmdRunner) ipfsContainer {

	o._cmdRunner = cmdRunner

	return o
}

func (o ipfsContainer) getStdoutStderrRedirectFiles(cmdName string) (*os.File, *os.File, error) {

	var stdoutFile, stderrFile *os.File

	var err error

	stdoutFile, err = o._createTempDir(o._tempDirLocation(), cmdName+"_*.stdout")
	if err != nil {
		return nil, nil, err
	}

	stderrFile, err = o._createTempDir(o._tempDirLocation(), cmdName+"_*.stderr")
	if err != nil {
		stdoutFile.Close() // CLOSE STDOUT FILE BEFORE RETURNING ERROR
		return nil, nil, err
	}

	return stdoutFile, stderrFile, nil

}

func (o ipfsContainer) pull() error {

	var stderrFile, stdoutFile *os.File
	var err error

	if len(o._cmdName) < 1 {
		// default if not overridden
		o._cmdName = cmdPull
	}

	stdoutFile, stderrFile, err = o.getStdoutStderrRedirectFiles(o._cmdName + "_cmd_*")
	if err != nil {
		return err
	}

	defer stderrFile.Close()
	defer stdoutFile.Close()

	if o._cmdRunner == nil {
		// set if not overridden
		o._cmdRunner = pullCmd(stdoutFile, stderrFile)
	}

	//cmd := o.getCmdRunner(pullCmd(stdoutFile, stderrFile))

	err = o._cmdRunner.Run()
	if err != nil {

		var wrappedErr error = dgrzerr.ExternalError.Wrap(err, "container pull request")

		var fi fs.FileInfo

		fi, err = stderrFile.Stat()

		if err != nil {
			wrappedErr = dgrzerr.AddErrorContext(wrappedErr,
				"stderr file redirect path",
				fmt.Sprintf("<not available - %v >", err))

		} else {
			wrappedErr = dgrzerr.AddErrorContext(wrappedErr,
				"stderr file redirect path", fi.Name())

		}

		fi, err = stdoutFile.Stat()

		if err != nil {
			wrappedErr = dgrzerr.AddErrorContext(wrappedErr, "stdout file redirect path", fmt.Sprintf("<not available - %v >", err))
		} else {
			wrappedErr = dgrzerr.AddErrorContext(wrappedErr, "stdout file redirect path", fi.Name())

		}

		wrappedErr = dgrzerr.AddErrorContext(wrappedErr, "command path", o._cmdRunner.String())

		return wrappedErr

	}

	return nil
}

func (o ipfsContainer) ps() ([]string, error) {

	var containerIdList []string

	var stdout bytes.Buffer
	var stderrFile *os.File
	var stdoutFile *os.File
	var err error

	stdoutFile, stderrFile, err = o.getStdoutStderrRedirectFiles(cmdPs + "_cmd_*")
	if err != nil {
		return containerIdList, err
	}

	defer stdoutFile.Close()
	defer stderrFile.Close()

	if o._cmdRunner == nil {
		// set if not overridden
		o._cmdRunner = pullCmd(stdoutFile, stderrFile)
	}

	err = o._cmdRunner.Run()
	if err != nil {
		err = dgrzerr.ExternalError.Wrap(err, "container ps command")
		err = dgrzerr.AddErrorContext(err, "command path", o._cmdRunner.String())

		var fi fs.FileInfo
		fi, err = stderrFile.Stat()

		if err != nil {
			err = dgrzerr.AddErrorContext(err, "stderr file redirect path",
				fmt.Sprintf("< not available - %v", err))
		} else {
			err = dgrzerr.AddErrorContext(err, "stderr file redirect path", fi.Name())
		}

	}

	if hasLine(&stdout) {
		be := dgrzerr.UnexpectedBehavior.Newf("expected container id in first line of stdout, nothing printed")
		dgrzerr.AddErrorContext(be, "command path", o._cmdRunner.String())
		dgrzerr.AddErrorContext(be, "stdout", stdout.String())

	}

	return containerIdList, nil
}

func hasLine(r io.Reader) bool {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	return scanner.Scan()

}

func (o ipfsContainer) run(runConfig RunConfig) (RunInfo, error) {

	runInfo := RunInfo{config: runConfig}

	cmd := runCmd(runConfig)

	// redirect stderr to file
	//tmpdir := os.TempDir()

	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout

	// DEFAULT TO BYTE STRING FOR STDERR
	//cmd.Stderr = &stderr

	stdoutFile, stderrFile, err := o.getStdoutStderrRedirectFiles(cmdRun + "_cmd_*")
	//stderrFile, err := o.createTempDir(tmpdir, "container_run_*.stderr")
	if err != nil {
		return runInfo, err
	}

	defer stderrFile.Close()
	defer stdoutFile.Close()

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	err = cmd.Run()
	if err != nil {
		fi, fi_err := stderrFile.Stat()
		var wrappedErr error

		if fi_err != nil {
			wrappedErr = dgrzerr.AddErrorContext(err, "stderr file redirect path",
				fmt.Sprintf("<not available - %v >", fi_err))
		} else {
			wrappedErr = dgrzerr.AddErrorContext(err, "stderr file redirect path", fi.Name())
		}

		return runInfo, wrappedErr
	}

	// VERIFY TCP LISTENER IS ALIVE

	err = waitForIpfsApiNetworkListener(runInfo)
	if err != nil {
		return runInfo, err
	}

	containerId := firstLine(stdoutFile)

	if len(containerId) < 1 {
		return runInfo, dgrzerr.UnexpectedBehavior.New("nerdctl run: expected to return container id in first line of stdout, found nothing")

	}

	return runInfo, nil
}

func firstLine(r io.Reader) string {

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	if scanner.Scan() {
		return scanner.Text()
	}

	return ""

}

func (o ipfsContainer) start(runInfo RunInfo) error {

	cmd := startCmd(runInfo.containerId)

	var stdout, stderr bytes.Buffer

	stdoutFile, stderrFile, err := o.getStdoutStderrRedirectFiles(cmdStart + "_cmd_*")
	if err != nil {
		return err
	}

	defer stderrFile.Close()
	defer stdoutFile.Close()

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	err = cmd.Run()
	if err != nil {
		fi, fi_err := stderrFile.Stat()
		var wrappedErr error

		if fi_err != nil {
			wrappedErr = dgrzerr.AddErrorContext(err, "stderr file redirect path",
				fmt.Sprintf("<not available - %v >", fi_err))
		} else {
			wrappedErr = dgrzerr.AddErrorContext(err, "stderr file redirect path", fi.Name())
		}

		return wrappedErr
	}

	err = waitForIpfsApiNetworkListener(runInfo)
	if err != nil {
		return err
	}

	outStr, _ := string(stdout.Bytes()), string(stderr.Bytes())

	containerId := strings.Trim(outStr, " \t\r\n")
	if containerId != runInfo.containerId {
		dgrzerr.UnexpectedValue.Newf("expected container id '%s' in stdout, found '%s'",
			runInfo.containerId, containerId)
	}

	return nil
}

func (o ipfsContainer) stop(runInfo RunInfo) error {

	cmd := stopCmd(runInfo.containerId)

	stdoutFile, stderrFile, err := o.getStdoutStderrRedirectFiles(cmdStop + "_cmd_*")
	//stderrFile, err := o.createTempDir(tmpdir, "container_run_*.stderr")
	if err != nil {
		return err
	}

	defer stderrFile.Close()
	defer stdoutFile.Close()

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	err = cmd.Run()
	if err != nil {
		return err
	}

	var containerId string

	containerId = firstLine(stdoutFile)

	if containerId != runInfo.containerId {
		dgrzerr.UnexpectedValue.Newf("expected container id '%s' in stdout, found '%s'",
			runInfo.containerId, containerId)
	}

	return nil
}

func waitForIpfsApiNetworkListener(runInfo RunInfo) error {

	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", runInfo.config.Host,
			runInfo.config.Port), defaultTimeout)
	if err != nil {
		return fmt.Errorf("error occurred while waiting for IPFS API Network Listener: "+
			"{ timeout = %f seconds, error = '%s' }",
			defaultTimeout.Seconds(), err.Error())

	}
	conn.Close()

	return nil

}
