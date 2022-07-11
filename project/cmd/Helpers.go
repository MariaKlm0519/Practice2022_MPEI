package main

import (
	"archive/zip"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/svc/mgr"
	"gopkg.in/ini.v1"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Quote struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
type SysInfo struct {
	Hostname string `json:"hostname"`
	Platform string `json:"platform"`
	CPU      string `json:"cpu"`
	RAM      uint64 `json:"ram"`
	Disk     uint64 `json:"disk"`
}
type Service struct {
	Name   string       `json:"name"`
	Config mgr.Config   `json:"config"`
	Status uint32       `json:"status"`
	srv    *mgr.Service `json:"service"`
}

type templateData interface{}

func render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
	files := []string{
		name,
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	rs, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = rs.Execute(w, td)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func ListServices(get_status uint32) ([]Service, error) {
	m, err := mgr.Connect()
	defer m.Disconnect()
	if err != nil {
		return nil, errors.New("can't connect to service control manager")
	}
	names, err := m.ListServices()
	if err != nil {
		return nil, errors.New("can't get service list")
	}
	var result []Service
	for i := 0; i < len(names); i++ {
		serv, err := m.OpenService(names[i])
		if err != nil {
			continue
		}
		status, err := serv.Query()
		if err != nil {
			continue
		}
		if uint32(status.State) == get_status {
			config, err := serv.Config()
			if err != nil {
				continue
			}
			newserv := Service{names[i], config, uint32(status.State), serv}
			result = append(result, newserv)
		}
	}
	return result, nil
}

func ChangeIni(name string) error {
	cfg, err := ini.Load(name)
	if err != nil {
		return errors.New("can't load ini-file")
	}
	status := cfg.Section("Options").Key("Enabled").Value()
	var new_status string
	if status == "1" {
		new_status = "0"
	} else {
		new_status = "1"
	}
	cfg.Section("Options").Key("Enabled").SetValue(new_status)
	err = cfg.SaveTo(name)
	if err != nil {
		return errors.New("can't save ini-file")
	}
	return nil
}

func ListDirByWalk(file_path string, zip_path string, t1 time.Time, t2 time.Time) (*os.File, error) {
	name := time.Now().Format("02012006150405") + ".zip"
	outFile, err := os.Create(zip_path + "\\" + name)
	if err != nil {
		return nil, errors.New("can't create output file")
	}
	zipW := zip.NewWriter(outFile)

	filepath.Walk(file_path, func(wPath string, info os.FileInfo, err error) error {
		if wPath == file_path {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().After(t1) && info.ModTime().Before(t2) {
			dat, _ := ioutil.ReadFile(wPath)
			f, _ := zipW.Create(info.Name())
			f.Write(dat)
		}
		return nil
	})
	err = zipW.Close()
	if err != nil {
		return nil, errors.New("can't close zip writer")
	}
	err = outFile.Close()
	if err != nil {
		return nil, errors.New("can't close output file")
	}
	return outFile, nil
}

func GetSystemInfo() (SysInfo, error) {
	var info SysInfo
	hostStat, err := host.Info()
	if err != nil {
		return info, errors.New("can't get system host info")
	}
	cpuStat, err := cpu.Info()
	if err != nil {
		return info, errors.New("can't get system cpu info")
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return info, errors.New("can't get system memory info")
	}
	diskStat, err := disk.Usage("\\")
	if err != nil {
		return info, errors.New("can't get system disk info")
	}

	info.Hostname = hostStat.Hostname
	info.Platform = hostStat.Platform
	info.CPU = cpuStat[0].ModelName
	info.RAM = vmStat.Total / 1024 / 1024
	info.Disk = diskStat.Free / 1024 / 1024
	return info, nil
}
