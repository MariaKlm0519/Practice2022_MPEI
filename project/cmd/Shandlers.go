package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

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
	switch r.Method {
	case http.MethodPost:
		{
			body, _ := ioutil.ReadAll(r.Body)
			var s struct {
				Stat string `json:"status"`
			}
			err := json.Unmarshal(body, &s)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			get_status, err := strconv.Atoi(s.Stat)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			data, err := ListServices(uint32(get_status))
			if err != nil {
				w.WriteHeader(500)
				return
			}
			newq, err := json.Marshal(data)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "1000")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
	switch r.Method {
	case http.MethodPost:
		{
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			var n struct {
				Name string `json:"name"`
			}
			err = json.Unmarshal(body, &n)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			err = ChangeIni(n.Name)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "1000")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
			outfile, err := ListDirByWalk("\\Документы\\Project_goland\\logs", "\\Документы\\Project_goland", time_before, time_after)
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
			info, err := GetSystemInfo()
			if err != nil {
				w.WriteHeader(500)
				return
			}
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
