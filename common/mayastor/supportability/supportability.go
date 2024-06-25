package supportability

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	v1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"

	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

var TmpDir = fmt.Sprintf("/tmp/%s/", RandomString(5))

func GetDumpDirectory() (string, error) {
	var e *fs.PathError
	files, err := os.ReadDir(TmpDir)
	if err != nil {
		if errors.As(err, &e) {
			return "", nil
		}
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("unable to find any files")
	}
	sort.Slice(files, func(i, j int) bool {
		infoJ, _ := files[j].Info()
		infoI, _ := files[i].Info()
		return infoJ.ModTime().Before(infoI.ModTime())
	})
	r := regexp.MustCompile(`([a-zA-Z].*)-2(\d{3}-\d{2}-\d{2}--\d{2}-\d{2}-\d{2})-UTC.*`)
	for _, f := range files {
		match := r.MatchString(f.Name())
		if match {
			if !f.IsDir() {
				err = unTar(f.Name())
				if err != nil {
					logf.Log.Error(err, "can't untar dump from dir", "dir_name", TmpDir, "file_name", f.Name())
					return "", err
				}
				return strings.Split(f.Name(), ".")[0], nil
			}
			return f.Name(), nil
		}
	}
	return "", errors.New("dump file not found")
}

var FilesMap map[string][]fs.FileInfo

func GetFiles(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		logf.Log.Error(err, "problem to get a files from directory", "path", dirPath)
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			err = GetFiles(fmt.Sprintf("%s/%s", dirPath, file.Name()))
			if err != nil {
				return err
			}
		} else {
			/*
				cutting off first part of path e.g.:
				/tmp/xyari/mayastor-2022-08-02--11-48-46-UTC/logs/operator-diskpool/node-0-11266_mayastor-operator-diskpool-c5f4b9447-tf9pc -> logs/operator-diskpool/node-0-11266_mayastor-operator-diskpool-c5f4b9447-tf9pc
			*/
			fileInfo, _ := file.Info()
			dir := strings.SplitN(dirPath, "/", 5)
			FilesMap[dir[len(dir)-1]] = append(FilesMap[dir[len(dir)-1]], fileInfo)
		}
	}
	return nil
}

func unTar(fileName string) error {
	fileName = TmpDir + fileName
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(fileName[:len(fileName)-7], header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
}

func GetPodsWithLoggingLabel() ([]coreV1.Pod, error) {
	pl, err := k8stest.ListPod(common.NSMayastor())
	if err != nil {
		logf.Log.Error(err, "can't get pods", "namespace", common.NSMayastor())
		return nil, err
	}
	dpl, err := k8stest.ListPod(common.NSDefault)
	var pods []coreV1.Pod
	if err != nil {
		logf.Log.Error(err, "can't get pods", "namespace", common.NSDefault)
		return nil, err
	}
	for _, p := range pl.Items {
		if p.Labels[e2e_config.GetConfig().Product.LoggingLabel] == "true" {
			pods = append(pods, p)
		}
	}
	for _, dp := range dpl.Items {
		if dp.Labels[e2e_config.GetConfig().Product.LoggingLabel] == "true" {
			pods = append(pods, dp)
		}
	}
	if len(pods) == 0 {
		return nil, errors.New("unable to find any pods")
	}
	return pods, nil
}

func SystemDump() (string, error) {
	bp := v1.GetPluginPath()
	cmd := exec.Command(bp, "dump", "system", "-n", common.NSMayastor(), "-d", TmpDir)
	start := time.Now()
	logf.Log.Info("Collecting system log dump starts", "namespace", common.NSMayastor(), "start_time", start)
	err := cmd.Run()
	if err != nil {
		logf.Log.Error(err, "the plugin raises an exception")
		return "", err
	}
	logf.Log.Info("Collecting system log dump ends", "namespace", common.NSMayastor(), "end_time", time.Now(), "duration", time.Since(start))
	return GetDumpDirectory()
}

func RemoveDumpFile(fn string) error {
	err := os.RemoveAll(fmt.Sprintf("%s%s", TmpDir, fn))
	if err != nil {
		return err
	}
	err = os.RemoveAll(fmt.Sprintf("%s%s.tar.gz", TmpDir, fn))
	if err != nil {
		return err
	}
	return nil
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
