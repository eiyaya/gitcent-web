package services

import (
	"fmt"

	"github.com/revel/revel"
)

//InitServices init background service
func InitServices(config *revel.MergedConfig) {
	fmt.Println(config.String("git.cmd"))
	initGitSshd()
}

func initGitSshd() {

}
