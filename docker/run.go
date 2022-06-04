package main

import (
	"github.com/mark4z/ideas/docker/cgroups"
	"github.com/mark4z/ideas/docker/cgroups/subsystems"
	"github.com/mark4z/ideas/docker/container"
	"github.com/mark4z/ideas/docker/log"
	"os"
	"strings"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("docker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()

	pwd, _ := os.Getwd()
	mntURL := pwd + "/root/mnt/"
	rootURL := pwd + "/root/"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
