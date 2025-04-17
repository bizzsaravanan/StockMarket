package service

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	services map[string]interface{} = make(map[string]interface{})
)

func RegisterService(serviceName string, service interface{}) {
	services[serviceName] = service
}

func Connect() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully waits for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	contextPath := viper.GetString("ContextPath")

	if viper.GetString("Host") == "" || viper.GetString("Port") == "" || contextPath == "" {
		log.Fatal("Missing required configuration values.")
	}

	r.HandleFunc(contextPath+"/api/{service}/{method}", handle)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/browser/")))

	srv := &http.Server{
		Addr:    viper.GetString("Host") + ":" + viper.GetString("Port"),
		Handler: r,
	}

	log.Printf("Starting server on %s\n", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v\n", err)
	}

	log.Println("shutting down")
	os.Exit(0)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust this to specify allowed origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	rr := make(map[string]interface{})
	rr["MetaData"] = r.Header
	rr["RequestTime"] = time.Now()
	rr["ipAddress"] = getIpAddress(r)

	params := mux.Vars(r)
	body, _ := ioutil.ReadAll(r.Body)
	rr["RequestBody"] = body
	serviceName := params["service"]
	methodName := params["method"]
	requestName := serviceName + ":" + methodName
	rr["ServiceName"] = serviceName
	rr["MethodName"] = methodName
	rr["RequestName"] = requestName
	resp, _ = handleRpc(rr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Error:HandleRpc:", err)
	}
}

func getIpAddress(r *http.Request) string {
	ipaddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !(r.Header.Get("X-Forwarded-For") == "") {
		ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
		ipaddress = ips[0]
	}
	return ipaddress
}

func handleRpc(rr map[string]interface{}) (resp interface{}, ctx context.Context) {
	var err error
	var c interface{}
	// if c, err = middleware(rr); err != nil {
	// 	log.Println("Error:HandleRpc:", err)
	// }
	if c != nil {
		ctx, _ = c.(context.Context)
	} else {
		ctx = context.Background()
	}
	mfp := make(map[string]interface{})
	if ctx != nil {
		ctx = context.WithValue(ctx, "MFP", mfp)
	}
	if err == nil {
		log.Printf("Starts ServiceName: %s, MethodName: %s, In %d millisecond", rr["ServiceName"], rr["MethodName"], time.Now().String(), 0)
		body, _ := rr["RequestBody"].([]byte)
		service, _ := rr["ServiceName"].(string)
		method, _ := rr["MethodName"].(string)

		if resp, err = invoke(service, method, string(body), ctx); err != nil {
			log.Println("Error: invoke :", err)
		}
		requestTime, _ := rr["RequestTime"].(time.Time)
		log.Printf("Ends ServiceName: %s, MethodName: %s, In %d millisecond", rr["ServiceName"], rr["MethodName"], time.Now().String(), time.Now().Sub(requestTime).Seconds()*1000)
	}
	if err != nil {
		resp = "Internal Error"
	}
	return resp, ctx
}

func invoke(serviceName, methodName, payload string, ctx context.Context) (interface{}, error) {
	if services[serviceName] == nil {
		log.Printf("Error: invoke service not found: %s\n", serviceName)
		return nil, errors.New("Internal Error")
	}

	sName := reflect.New(reflect.TypeOf(services[serviceName]).Elem()).Interface()
	method := reflect.ValueOf(sName).MethodByName(methodName)
	if !method.IsValid() {
		log.Printf("Error: invoke method not found: %s\n", methodName)
		return nil, errors.New("Internal Error")
	}

	mtype := method.Type()
	args := mtype.NumIn()

	var tp reflect.Type
	if args == 2 {
		tp = mtype.In(1)
	} else {
		tp = mtype.In(0)
	}

	par := reflect.New(tp).Interface()
	err := json.Unmarshal([]byte(payload), par)
	if err != nil {
		log.Printf("Error: invoke payload unmarshalling failed: %v\n", err)
		return nil, errors.New("Invalid Request")
	}

	v := reflect.ValueOf(par).Elem() // Get the value of the pointer
	if v.IsZero() {
		log.Println("Error: unmarshalled parameter is zero value")
		return nil, errors.New("Invalid Request")
	}

	if ctx == nil {
		log.Println("Error: context is nil")
		return nil, errors.New("Invalid Request")
	}

	var vals []reflect.Value
	if args == 2 {
		h := reflect.ValueOf(ctx)
		vals = method.Call([]reflect.Value{h, v})
	} else {
		vals = method.Call([]reflect.Value{v})
	}

	if len(vals) < 2 {
		log.Println("Error: method did not return expected values")
		return nil, errors.New("Internal Error")
	}

	respErr, _ := vals[1].Interface().(error)
	return vals[0].Interface(), respErr
}
