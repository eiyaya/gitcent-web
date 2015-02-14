package controllers

import (
	"fmt"
	"os/exec"

	"github.com/revel/revel"
)

// HTTPBackEnd (deal with git smart http protcol)
type HTTPBackEnd struct {
	*revel.Controller
}

// GetInfoRefs (got git repo refs)
func (h HTTPBackEnd) GetInfoRefs(service string) revel.Result {
	h.Response.Out.Header().Add("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	h.Response.Out.Header().Add("Pragma", "no-cache")
	h.Response.Out.Header().Add("Cache-Control", "no-cache, max-age=0, must-revalidate")
	h.Response.Out.Header().Add("Content-Type", "application/x-git-receive-pack-advertisement")

	repo := "/Users/stephenzhen/gitcent-repos/test"
	git := "/usr/bin/git"

	refs, _ := exec.Command(git, "receive-pack", repo, "--stateless-rpc", "--advertise-refs").Output()

	act := "# service=git-receive-pack\n"
	l := len(act) + 4
	s := "00" + fmt.Sprintf("%x", l) + act
	h.Response.Out.Write([]byte(s))
	h.Response.Out.Write([]byte{'0', '0', '0', '0'})
	h.Response.Out.Write(refs)
	return h.RenderText("")
}
