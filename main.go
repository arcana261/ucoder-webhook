package main

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/go-playground/webhooks.v4"
	"gopkg.in/go-playground/webhooks.v4/github"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	path = "/webhooks"
)

func main() {
	flag.String("port", "3000", "port to listen to")
	flag.String("secret", "", "github secret to use")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	port := viper.GetString("port")
	secret := viper.GetString("secret")

	hook := github.New(&github.Config{Secret: secret})
	hook.RegisterEvents(handlePush, github.PushEvent)

	err := webhooks.Run(hook, ":"+port, path)
	if err != nil {
		fmt.Println(err)
	}
}

func handlePush(payload interface{}, header webhooks.Header) {
	pl := payload.(github.PushPayload)

	switch pl.Repository.FullName {
	case "arcana261/ucoder":
		log.Printf("Handling push event for arcana261/ucoder\n")
		handleUcoder()
	default:
		log.Printf("Repository not recognized: %s\n", pl.Repository.FullName)
	}
}

func handleUcoder() {
	cwd := "/usr/share/nginx/ucoder.ir"
	run(cwd, "/bin/git", "pull", "origin", "master")
	run(cwd, "/usr/local/bin/hexo", "generate")
	run(cwd, "/bin/chcon", "-Rt", "httpd_sys_content_t", cwd)
	run(cwd, "/bin/systemctl", "restart", "nginx")
}
