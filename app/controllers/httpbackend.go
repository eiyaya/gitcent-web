package controllers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/revel/revel"
)

// HTTPBackEnd (deal with git smart http protcol)
type HTTPBackEnd struct {
	App
}

// GitUploadPack (deal with git clone)
func (h HTTPBackEnd) GitUploadPack(repo string, group string) revel.Result {

	h.RequireUser()

	h.Response.Out.Header().Add("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	h.Response.Out.Header().Add("Pragma", "no-cache")
	h.Response.Out.Header().Add("Cache-Control", "no-cache, max-age=0, must-revalidate")
	h.Response.Out.Header().Del("Content-Type")
	h.Response.Out.Header().Add("Content-Type", "application/x-git-upload-pack-result")

	repoRoot, _ := revel.Config.String("repo.root")
	repo = repoRoot + repo
	git, _ := revel.Config.String("git.cmd")
	cmd := exec.Command(git, "upload-pack", "--stateless-rpc", repo)
	out, _ := cmd.StdoutPipe()
	input, _ := cmd.StdinPipe()
	cmd.Stderr = os.Stderr

	go func() {
		var err error
		var bytes = make([]byte, 1024)
		var read int
		for {
			read, err = h.Request.Body.Read(bytes)
			input.Write(bytes[:read])
			if err != nil {
				fmt.Println("read http request err:", err)
				input.Close()
				break
			}

		}
	}()
	cmd.Start()
	var err error
	var bytes = make([]byte, 1024)
	var read int
	for {
		read, err = out.Read(bytes)
		h.Response.Out.Write(bytes[:read])
		if err != nil {
			out.Close()
			fmt.Println("read cmd output err:", err)
			break
		}

	}
	cmd.Wait()

	return nil
}

// GitReceivePack (deal with git push)
func (h HTTPBackEnd) GitReceivePack(repo string, group string) revel.Result {
	h.Response.Out.Header().Add("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	h.Response.Out.Header().Add("Pragma", "no-cache")
	h.Response.Out.Header().Add("Cache-Control", "no-cache, max-age=0, must-revalidate")
	h.Response.Out.Header().Del("Content-Type")
	h.Response.Out.Header().Add("Content-Type", "application/x-git-receive-pack-result")

	repoRoot, _ := revel.Config.String("repo.root")
	repo = repoRoot + repo
	git, _ := revel.Config.String("git.cmd")
	cmd := exec.Command(git, "receive-pack", "--stateless-rpc", repo)
	out, _ := cmd.StdoutPipe()
	input, _ := cmd.StdinPipe()
	cmd.Stderr = os.Stderr
	go func() {
		var err error
		var bytes = make([]byte, 1024)
		var read int
		for {
			read, err = h.Request.Body.Read(bytes)
			input.Write(bytes[:read])
			if err != nil {
				fmt.Println("read http request err:", err)
				input.Close()
				break
			}

		}
	}()
	cmd.Start()
	var err error
	var bytes = make([]byte, 1024)
	var read int
	for {
		read, err = out.Read(bytes)
		h.Response.Out.Write(bytes[:read])
		if err != nil {
			out.Close()
			fmt.Println("read cmd output err:", err)
			break
		}

	}
	cmd.Wait()

	return nil
}

// GetInfoRefs (got git repo refs)
func (h HTTPBackEnd) GetInfoRefs(service string, repo string, group string) revel.Result {
	//service: git-receive-pack|git-upload-pack
	r, _ := regexp.Compile("git-")
	action := r.ReplaceAllString(service, "")
	r, _ = regexp.Compile("-pack")
	action = r.ReplaceAllString(action, "")
	h.Response.Out.Header().Add("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	h.Response.Out.Header().Add("Pragma", "no-cache")
	h.Response.Out.Header().Add("Cache-Control", "no-cache, max-age=0, must-revalidate")
	h.Response.Out.Header().Del("Content-Type")
	h.Response.Out.Header().Add("Content-Type", "application/x-git-"+action+"-pack-advertisement")

	repoRoot, _ := revel.Config.String("repo.root")
	repo = repoRoot + repo
	git, _ := revel.Config.String("git.cmd")
	refs, _ := exec.Command(git, action+"-pack", "--stateless-rpc", "--advertise-refs", repo).Output()
	act := "# service=git-" + action + "-pack\n"
	l := len(act) + 4
	s := "00" + fmt.Sprintf("%x", l) + act
	h.Response.Out.Write([]byte(s))
	h.Response.Out.Write([]byte{'0', '0', '0', '0'})
	h.Response.Out.Write(refs)

	return nil
}
