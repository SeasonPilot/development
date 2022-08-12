package main

import (
	"fmt"

	"github.com/marmotedu/iam/pkg/log"
	"github.com/spf13/pflag"
)

func main() {
	fs := pflag.NewFlagSet("test", pflag.ExitOnError)
	opt := log.NewOptions()
	opt.AddFlags(fs)

	fmt.Printf("%#v\n", opt.Level) // "info"

	args := []string{"--log.level=debug"}
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v\n", "debug" == opt.Level)
}
