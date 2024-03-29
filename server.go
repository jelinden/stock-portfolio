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
var indexPage string

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
	f, err := fsPublic.Open("index.html")
	if err != nil {
		log.Println(err)
	}
	result, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		log.Println(err)
	}
	indexPage = string(result)
	db.Init()
}

func main() {
	fromEmail = os.Getenv("FROMEMAIL")
	adminUser := os.Getenv("ADMINUSER")
	apiToken := os.Getenv("IEXAPITOKEN")
	if fromEmail == "" || apiToken == "" {
		log.Fatal("FROMEMAIL or IEXAPITOKEN was not set")
	}
	config.Config.FromEmail = fromEmail
	config.Config.AdminUser = adminUser
	config.Config.Token = apiToken
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init()
	db.InitBolt()

	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	fsAssets := util.JustFilesFilesystem{Fs: http.Dir("build/")}

	//router.Handler("GET", "/build/*name", http.StripPrefix("/build", util.GH(http.FileServer(fsAssets))))
	router.Handler("GET", "/js/*.js", http.StripPrefix("/js", util.GH(http.FileServer(fsAssets))))
	router.Handler("GET", "/css/*.css", http.StripPrefix("/css", util.GH(http.FileServer(fsAssets))))
	router.Handler("GET", "/img/*.png", http.StripPrefix("/img", util.GH(http.FileServer(fsAssets))))
	router.Handler("GET", "/favicon.ico", util.GH(http.FileServer(fsAssets)))
	router.Handle("GET", "/", middleware.HttpLogger(middleware.RequestCounter(index)))
	router.Handle("GET", "/login", middleware.HttpLogger(middleware.RequestCounter(index)))
	router.Handle("GET", "/signup", middleware.HttpLogger(middleware.RequestCounter(index)))
	router.Handle("GET", "/portfolio/:id", middleware.HttpLogger(middleware.RequestCounter(index)))
	router.Handle("GET", "/transactions/:id", middleware.HttpLogger(middleware.RequestCounter(index)))
	router.Handle("POST", "/login", middleware.HttpLogger(middleware.RequestCounter(routes.Login)))
	router.Handle("GET", "/logout", middleware.HttpLogger(middleware.RequestCounter(routes.Logout)))
	router.Handle("POST", "/signup", middleware.HttpLogger(middleware.RequestCounter(routes.Signup)))
	router.Handle("GET", "/verify", middleware.HttpLogger(middleware.RequestCounter(index)))
	router.Handle("GET", "/verify/:verifystring", middleware.HttpLogger(middleware.RequestCounter(routes.Verify)))
	router.Handle("GET", "/health", middleware.AdminAuth(middleware.HttpLogger(middleware.RequestCounter(index))))

	router.Handle("GET", "/api/user", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.UserJSON)))))
	router.Handle("GET", "/api/allusers", middleware.HttpLogger(middleware.RequestCounter(middleware.AdminAuth(util.MakeGzipHandler(routes.AllUsers)))))
	router.Handle("POST", "/api/portfolio/create", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.AddPortfolio)))))
	router.Handle("GET", "/api/portfolios", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.GetPortfolios)))))
	router.Handle("GET", "/api/portfolio/get/:id", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.GetPortfolio)))))
	router.Handle("POST", "/api/portfolio/add", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.AddStock)))))
	router.Handle("GET", "/api/portfolio/remove/:portfolioid/:symbol/:transactionid", middleware.RequestCounter(middleware.HttpLogger(middleware.Auth(util.MakeGzipHandler(routes.RemoveStock)))))
	router.Handle("GET", "/api/dividends", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.GetDividend)))))
	router.Handle("GET", "/api/transactions/:id", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.GetTransactions)))))
	router.Handle("GET", "/api/history/:id", middleware.HttpLogger(middleware.RequestCounter(middleware.Auth(util.MakeGzipHandler(routes.GetHistory)))))
	router.Handle("GET", "/api/health", middleware.HttpLogger(middleware.RequestCounter(middleware.AdminAuth(util.MakeGzipHandler(routes.Health)))))

	gracefullShutdown()
	log.Println("starting server at port 3300")
	log.Fatal(http.ListenAndServe(":3300", router))
}

var fsPublic = util.JustFilesFilesystem{Fs: http.Dir("build/")}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var lastModifiedTimeStamp = time.Now().Add(6 * time.Hour).Format(http.TimeFormat)
	var noBrowserCache = time.Now().Add(-6 * time.Hour).Format(http.TimeFormat)
	w.Header().Add("Cache-Control", "no-store, private, no-cache, must-revalidate")
	w.Header().Add("Expires", noBrowserCache)
	w.Header().Add("Last-Modified", lastModifiedTimeStamp)
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, indexPage)
}

func gracefullShutdown() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 3 second to finish processing")
		db.After()
		db.CloseBolt()
		middleware.CloseLogFile()
		time.Sleep(3 * time.Second)
		os.Exit(0)
	}()
}
