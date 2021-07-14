package testserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)


func ForwardServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", ReadCheckReq)
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
	fmt.Println("forward server starts")
}

func ReadCheckReq(w http.ResponseWriter, r *http.Request){
	fmt.Println("forwarding received")
	reqData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err", err)
	}

	w.Write([]byte("forward received"))
	fmt.Printf("forward data: %s",string(reqData))


}

// TestClient sends request to slackbot server
func TestClient(visitURL string, form url.Values) (*http.Response, error){
	//visit, err := url.Parse(visitURL)
	//if err != nil {
	//	return nil, err
	//}
	//req := http.Request{
	//	Method: "POST",
	//	URL: visit,
	//	PostForm: form,
	//}

	//client := &http.Client{}
	//res, err := client.Do(&req)

	resp, err := http.PostForm(visitURL, form)
	if err != nil {
		fmt.Println(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("test client: %s\n", string(bytes))

	return resp, err
}

