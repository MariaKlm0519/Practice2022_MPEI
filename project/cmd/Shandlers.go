package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/svc/mgr"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

var m Quote

func restroutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { return })
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) { return })
	mux.HandleFunc("/api/records", getrecords)
	mux.HandleFunc("/api/list", getlist)
	mux.HandleFunc("/api/ini", setini)
	mux.HandleFunc("/api/searchLog", searchLog)
	return mux
}

func getrecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.Method {
	case http.MethodGet:
		{
			newq, err := json.Marshal(Quote{2, m.Title, m.Text + time.Now().Format(" 2006-01-02 15:04:05")})
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Write(newq)
			w.WriteHeader(200)
		}
	case http.MethodPost:
		{
			body, _ := ioutil.ReadAll(r.Body)
			_ = json.Unmarshal(body, &m)
			w.WriteHeader(200)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			w.WriteHeader(423)
			return
		}
	}
}

func getlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.Method {
	case http.MethodPost:
		{
			body, _ := ioutil.ReadAll(r.Body)
			var s struct {
				Stat string `json:"status"`
			}
			_ = json.Unmarshal(body, &s)
			get_status, _ := strconv.Atoi(s.Stat)
			data := ListServices(uint32(get_status))
			newq, _ := json.Marshal(data)
			w.Write(newq)
			w.WriteHeader(200)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			w.WriteHeader(423)
			return
		}
	}
}

func setini(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.Method {
	case http.MethodPost:
		{
			body, _ := ioutil.ReadAll(r.Body)
			var n struct {
				Name string `json:"name"`
			}
			_ = json.Unmarshal(body, &n)
			cfg, err := ini.Load(n.Name)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			status := cfg.Section("Options").Key("Enabled").Value()
			var new_status string
			if status == "1" {
				new_status = "0"
			} else {
				new_status = "1"
			}
			cfg.Section("Options").Key("Enabled").SetValue(new_status)
			err = cfg.SaveTo(n.Name)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			w.WriteHeader(423)
			return
		}
	}
}

func searchLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "application/zip")
	switch r.Method {
	case http.MethodPost:
		{
			body, _ := ioutil.ReadAll(r.Body)
			var dat struct {
				Date_start string `json:"date_start"`
				Time_start string `json:"time_start"`
				Date_end   string `json:"date_end"`
				Time_end   string `json:"time_end"`
			}
			_ = json.Unmarshal(body, &dat)
			time_before, _ := time.ParseInLocation("2006-01-02 15:04", dat.Date_start+" "+dat.Time_start, time.Local)
			time_after, _ := time.ParseInLocation("2006-01-02 15:04", dat.Date_end+" "+dat.Time_end, time.Local)
			outfile, err := listDirByWalk("\\Документы\\Project_goland\\logs", "\\Документы\\Project_goland", time_before, time_after)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			fileBytes, _ := ioutil.ReadFile(outfile.Name())
			w.Write(fileBytes)
			w.WriteHeader(200)
		}
	case http.MethodGet:
		{
			hostStat, _ := host.Info()
			cpuStat, _ := cpu.Info()
			vmStat, _ := mem.VirtualMemory()
			diskStat, _ := disk.Usage("\\")

			var info SysInfo
			info.Hostname = hostStat.Hostname
			info.Platform = hostStat.Platform
			info.CPU = cpuStat[0].ModelName
			info.RAM = vmStat.Total / 1024 / 1024
			info.Disk = diskStat.Free / 1024 / 1024

			data, err := json.Marshal(info)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Write(data)
			w.WriteHeader(200)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			w.WriteHeader(423)
			return
		}
	}
}

func ListServices(get_status uint32) []Service {
	m, _ := mgr.Connect()
	names, _ := m.ListServices()
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
	return result
}

func listDirByWalk(file_path string, zip_path string, t1 time.Time, t2 time.Time) (*os.File, error) {

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
