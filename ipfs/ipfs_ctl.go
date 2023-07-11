//go:build darwin

package ipfs

type RunConfig struct {
	Host string // BIND HOST ADDR/NAME
	Port int    // BIND HOST PORT
	Tag  string // OPTIONAL
}

type RunInfo struct {
	config      RunConfig
	containerId string
}

func PullDefault() error {
	return newIpfsContainer().pull()

}

func RunDefault() (RunInfo, error) {

	runConfig := RunConfig{Host: "127.0.0.1",
		Port: 5001,
		Tag:  "ipfs"}
	return newIpfsContainer().run(runConfig)

}
func Run(runConfig RunConfig) (RunInfo, error) {

	return newIpfsContainer().run(runConfig)
}

func Start(runInfo RunInfo) error {

	return newIpfsContainer().start(runInfo)
}

func Stop(runInfo RunInfo) error {

	return newIpfsContainer().stop(runInfo)
}
