package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"regexp"
	"sort"
	"strings"

	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/micro/internal/handler"
	"github.com/serenize/snaker"
)

var (
	re = regexp.MustCompile("^[a-zA-Z0-9]+$")
	// Default address to bind to
	Address = ":8082"
	// The namespace to serve
	// Example:
	// Namespace + /[Service]/foo/bar
	// Host: Namespace.Service Endpoint: /foo/bar
	Namespace = "go.micro.web"
	// Base path sent to web service.
	// This is stripped from the request path
	// Allows the web service to define absolute paths
	BasePathHeader = "X-Micro-Web-Base-Path"
)

type server struct {
	*mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}

	s.Router.ServeHTTP(w, r)
}

func (s *server) proxy() http.Handler {
	sel := selector.NewSelector(
		selector.Registry((*cmd.DefaultOptions().Registry)),
	)

	director := func(r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			return
		}
		if !re.MatchString(parts[1]) {
			return
		}
		next, err := sel.Select(Namespace + "." + parts[1])
		if err != nil {
			return
		}
		r.URL.Scheme = "http"
		s, err := next()
		if err != nil {
			return
		}

		r.Header.Set(BasePathHeader, "/"+parts[1])
		r.URL.Host = fmt.Sprintf("%s:%d", s.Address, s.Port)
		r.URL.Path = "/" + strings.Join(parts[2:], "/")
	}
	return &httputil.ReverseProxy{
		Director: director,
	}
}

func format(v *registry.Value) string {
	if v == nil || len(v.Values) == 0 {
		return "{}"
	}
	var f []string
	for _, k := range v.Values {
		f = append(f, formatEndpoint(k, 0))
	}
	return fmt.Sprintf("{\n%s}", strings.Join(f, ""))
}

func formatEndpoint(v *registry.Value, r int) string {
	// default format is tabbed plus the value plus new line
	fparts := []string{"", "%s %s", "\n"}
	for i := 0; i < r+1; i++ {
		fparts[0] += "\t"
	}
	// its just a primitive of sorts so return
	if len(v.Values) == 0 {
		return fmt.Sprintf(strings.Join(fparts, ""), snaker.CamelToSnake(v.Name), v.Type)
	}

	// this thing has more things, it's complex
	fparts[1] += " {"

	vals := []interface{}{snaker.CamelToSnake(v.Name), v.Type}

	for _, val := range v.Values {
		fparts = append(fparts, "%s")
		vals = append(vals, formatEndpoint(val, r+1))
	}

	// at the end
	l := len(fparts) - 1
	for i := 0; i < r+1; i++ {
		fparts[l] += "\t"
	}
	fparts = append(fparts, "}\n")

	return fmt.Sprintf(strings.Join(fparts, ""), vals...)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	var webServices []string
	for _, s := range services {
		if strings.Index(s.Name, Namespace) == 0 {
			webServices = append(webServices, strings.Replace(s.Name, Namespace+".", "", 1))
		}
	}

	sort.Strings(webServices)

	type templateData struct {
		HasWebServices bool
		WebServices    []string
	}

	data := templateData{len(webServices) > 0, webServices}
	render(w, r, indexTemplate, data)
}

func registryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	svc := r.Form.Get("service")

	if len(svc) > 0 {
		s, err := (*cmd.DefaultOptions().Registry).GetService(svc)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if len(s) == 0 {
			http.Error(w, "Not found", 404)
			return
		}

		if r.Header.Get("Content-Type") == "application/json" {
			b, err := json.Marshal(map[string]interface{}{
				"services": s,
			})
			if err != nil {
				http.Error(w, "Error occurred:"+err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}

		render(w, r, serviceTemplate, s)
		return
	}

	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	sort.Sort(sortedServices{services})

	if r.Header.Get("Content-Type") == "application/json" {
		b, err := json.Marshal(map[string]interface{}{
			"services": services,
		})
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	render(w, r, registryTemplate, services)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, queryTemplate, nil)
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
	}).Parse(layoutTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	t, err = t.Parse(tmpl)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}

func run() {
	r := mux.NewRouter()
	s := &server{r}
	s.HandleFunc("/registry", registryHandler)
	s.HandleFunc("/rpc", handler.RPC)
	s.HandleFunc("/query", queryHandler)
	s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(s.proxy())
	s.HandleFunc("/", indexHandler)

	log.Infof("Listening on %s", Address)

	if err := http.ListenAndServe(Address, s); err != nil {
		log.Fatal(err)
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "web",
			Usage: "Run the micro web app",
			Action: func(c *cli.Context) {
				run()
			},
		},
	}
}
