package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

	"log"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if apikey != "" && r.Header.Get("Authorization") != "Bearer "+apikey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.RequestURI == "/stats" {
		statsHandler(w, r)
	} else if r.RequestURI[:6] == "/items" {
		itemsHandler(w, r)
	} else if r.RequestURI[:28] == "/.well-known/acme-challenge/" {
		certManager.HTTPHandler(nil).ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/items" {
		switch r.Method {
		case http.MethodDelete:
			cache.Clear()
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	} else if r.RequestURI[:7] != "/items/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	key := []byte(r.RequestURI[7:])

	switch r.Method {
	case http.MethodGet:
		if val, err := cache.Get(key); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(val)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		exp, _ := strconv.Atoi(r.Header.Get("X-Cache-Expire"))
		buf := &bytes.Buffer{}
		buf.ReadFrom(r.Body)
		if err := cache.Set(key, buf.Bytes(), exp); err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	case http.MethodDelete:
		if cache.Del(key) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s := &struct {
			Cache struct {
				AverageAccessTime int64
				EntryCount        int64
				ExpiredCount      int64
				HitCount          int64
				HitRate           float64
				LookupCount       int64
				MissCount         int64
				OverwriteCount    int64
			}
			Mem runtime.MemStats
		}{}

		s.Cache.AverageAccessTime = cache.AverageAccessTime()
		s.Cache.EntryCount = cache.EntryCount()
		s.Cache.ExpiredCount = cache.ExpiredCount()
		s.Cache.HitCount = cache.HitCount()
		s.Cache.HitRate = cache.HitRate()
		s.Cache.LookupCount = cache.LookupCount()
		s.Cache.MissCount = cache.MissCount()
		s.Cache.OverwriteCount = cache.OverwriteCount()

		runtime.ReadMemStats(&s.Mem)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(s); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	case http.MethodDelete:
		cache.ResetStatistics()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
