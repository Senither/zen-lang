package objects

import "github.com/senither/zen-lang/objects/os"

func globalOSHostname(args ...Object) (Object, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &String{Value: name}, nil
}

func globalOSPlatform(args ...Object) (Object, error) {
	name, err := os.Platform()
	if err != nil {
		return nil, err
	}

	return &String{Value: name}, nil
}

func globalOSArch(args ...Object) (Object, error) {
	name, err := os.Arch()
	if err != nil {
		return nil, err
	}

	return &String{Value: name}, nil
}

func globalOSCPUs(args ...Object) (Object, error) {
	return &Integer{Value: os.CPUs()}, nil
}
