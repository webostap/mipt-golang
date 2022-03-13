package urlshortener

import (
	"fmt"
	"github.com/go-chi/chi"
	"hash/fnv"
	"net/http"
	"strconv"
)

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

type URLShortener struct {
	addr string
	m    map[uint64]string
}

func NewShortener(addr string) *URLShortener {
	return &URLShortener{
		addr: addr,
		m:    make(map[uint64]string),
	}
}

func (s *URLShortener) HandleSave(rw http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("u")
	h := hash(url)
	_, found := s.m[h]
	if !found {
		s.m[h] = url
		fmt.Fprintf(rw, "%s/%x", s.addr, h)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *URLShortener) HandleExpand(rw http.ResponseWriter, req *http.Request) {
	k := chi.URLParam(req, "key")
	h, _ := strconv.ParseUint(k, 16, 64)
	v, found := s.m[h]
	if found {
		http.Redirect(rw, req, v, http.StatusMovedPermanently)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}
