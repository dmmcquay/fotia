package fotia

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type failure struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewFailure(msg string) *failure {
	return &failure{
		Success: false,
		Error:   msg,
	}
}

func Sleep(w http.ResponseWriter, r *http.Request) {
	defer Time()()
	i := rand.Intn(10)
	fmt.Fprintf(w, "sleeping for %d", i)
	time.Sleep(time.Second * time.Duration(i))
}

func Up(w http.ResponseWriter, r *http.Request) {
	defer Time()()
	searchreq := r.URL.Path[len("/up/"):]
	if len(searchreq) == 0 {
		Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure("url could not be parsed"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	if searchreq[len(searchreq)-1] != '/' {
		http.Redirect(w, r, "/up/"+searchreq+"/", http.StatusMovedPermanently)
		return
	}
	searchReqParsed := strings.Split(searchreq, "/")
	s, err := strconv.Atoi(searchReqParsed[0])
	if err != nil {
		Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure(fmt.Sprintf("could not convert %v to an int", searchReqParsed[0])))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	ECount.Add(float64(s))
}

func Down(w http.ResponseWriter, r *http.Request) {
	defer Time()()
	searchreq := r.URL.Path[len("/down/"):]
	if len(searchreq) == 0 {
		Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure("url could not be parsed"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	if searchreq[len(searchreq)-1] != '/' {
		http.Redirect(w, r, "/down/"+searchreq+"/", http.StatusMovedPermanently)
		return
	}
	searchReqParsed := strings.Split(searchreq, "/")
	s, err := strconv.Atoi(searchReqParsed[0])
	if err != nil {
		Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure(fmt.Sprintf("could not convert %v to an int", searchReqParsed[0])))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	ECount.Sub(float64(s))
}
