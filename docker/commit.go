package main

import (
	"github.com/mark4z/ideas/docker/log"
	"os"
	"os/exec"
)

func commitContainer(imageName string) {
	pwd, _ := os.Getwd()
	mntURL := pwd + "/root/mnt"
	imageTar := pwd + "/root/" + imageName + ".tar"
	log.Infof("%s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
