package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jelinden/stock-portfolio/app/config"
	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/middleware"
	"github.com/jelinden/stock-portfolio/app/routes"
	"github.com/jelinden/stock-portfolio/app/util"
	"github.com/julienschmidt/httprouter"
)

var fromEmail, emailSendingPasswd string

type HTTPError struct {
	code    int
	message string
}

func (e *HTTPError) Error() string {
	return e.message
}

func Init() {
	var configFile string
	flag.StringVar(&configFile, "c", "app/config/config-localhost.json", "config")
	flag.Parse()
	log.Println("Loading configuration from file " + configFile)
	config.SetConfig(configFile)
}

func main() {
	fromEmail = os.Getenv("FROMEMAIL")
	emailSendingPasswd = os.Getenv("EMAILSENDINGPASSWD")
	adminUser := os.Getenv("ADMINUSER")
	if fromEmail == "" || emailSendingPasswd == "" {
		log.Fatal("FROMEMAIL or EMAILSENDINGPASSWD was not set")
	}
	config.Config.FromEmail = fromEmail
	config.Config.EmailSendingPasswd = emailSendingPasswd
	config.Config.AdminUser = adminUser
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init()
	db.Init()
	db.InitBolt()

	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	fsAssets := util.JustFilesFilesystem{Fs: http.Dir("build/")}

	router.Handler("GET", "/js/*.js", http.StripPrefix("/js", util.GH(http.FileServer(fsAssets))))
	router.Handler("GET", "/css/*.css", http.StripPrefix("/css", util.GH(http.FileServer(fsAssets))))
	router.Handler("GET", "/favicon.ico", util.GH(http.FileServer(fsAssets)))
	router.Handle("GET", "/", index)
	router.Handle("GET", "/login", index)
	router.Handle("GET", "/signup", index)
	router.Handle("GET", "/portfolio/:id", index)
	router.Handle("POST", "/login", routes.Login)
	router.Handle("GET", "/logout", routes.Logout)
	router.Handle("POST", "/signup", routes.Signup)
	router.Handle("GET", "/verify", index)
	router.Handle("GET", "/verify/:verifystring", routes.Verify)

	router.Handle("GET", "/api/user", middleware.HttpLogger(middleware.Auth(util.MakeGzipHandler(routes.UserJSON))))
	router.Handle("GET", "/api/allusers", middleware.AdminAuth(util.MakeGzipHandler(routes.AllUsers)))
	router.Handle("POST", "/api/portfolio/create", middleware.Auth(util.MakeGzipHandler(routes.AddPortfolio)))
	router.Handle("GET", "/api/portfolios", middleware.Auth(util.MakeGzipHandler(routes.GetPortfolios)))
	router.Handle("GET", "/api/portfolio/get/:id", middleware.Auth(util.MakeGzipHandler(routes.GetPortfolio)))
	router.Handle("POST", "/api/portfolio/add", middleware.Auth(util.MakeGzipHandler(routes.AddStock)))
	router.Handle("GET", "/api/portfolio/remove/:portfolioid/:symbol", middleware.Auth(util.MakeGzipHandler(routes.RemoveStock)))

	gracefullShutdown()
	log.Println("starting server at port 3300")
	log.Fatal(http.ListenAndServe(":3300", router))
}

var fsPublic = util.JustFilesFilesystem{Fs: http.Dir("public/")}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	f, err := fsPublic.Open("index.html")
	if err != nil {
		log.Println(err)
	}

	result, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, string(result))
}

func gracefullShutdown() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		db.After()
		db.CloseBolt()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
