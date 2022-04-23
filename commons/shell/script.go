package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anthony-dong/go-sdk/commons/logs"
)

var (
	shell       string
	shellPrefix string
)

/**
获取当前用户执行的shell.
*/
func init() {
	shell = os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	shellPrefix = shell + " -c "
}

func GitClone(sshAddr, dir string) error {
	return GitCloneBranch("master", sshAddr, dir)
}

func GitCloneBranch(branch, sshAddr, dir string) error {
	gitCmd := fmt.Sprintf("git clone -b %s %s %s", branch, sshAddr, dir)
	return Cmd(gitCmd)
}

func Run(shell string) error {
	return Cmd(shell)
}

func Copy(src, dest string) (err error) {
	gitCmd := fmt.Sprintf("cp -R %s %s", src, dest)
	return Cmd(gitCmd)
}

func Mv(src, dest string) (err error) {
	gitCmd := fmt.Sprintf("mv '%s' '%s'", src, dest)
	return Cmd(gitCmd)
}

// delete file.
func Delete(file ...string) (err error) {
	if len(file) == 0 {
		return
	}
	for _, elem := range file {
		if elem == "/" || strings.Contains(elem, "*") {
			return errors.New("can not delete * file")
		}
	}
	gitCmd := fmt.Sprintf("rm -r '%s'", strings.Join(file, " "))
	return Cmd(gitCmd)
}

func Cmd(cmd string) error {
	command := exec.Command(shell, "-c", cmd)
	logs.Infof("exec: %s", strings.TrimPrefix(strings.TrimPrefix(command.String(), " "), shellPrefix))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
