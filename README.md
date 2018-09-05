# Clova CEK SDK Go

Go SDK for the Clova Extensions


## About Clova Extension Kit

See the official documentation for more information.

- https://clova-developers.line.me/guide/


## Installation ##

```sh
$ go get github.com/line/clova-cek-sdk-go/cek
```

## Usage

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/line/clova-cek-sdk-go/cek"
)

func main() {
	ext := cek.NewExtension("com.example.my_extension")
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		messagee, err := ext.ParseRequest(r)
		if err != nil {
			log.Printf("invalid request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var response *cek.ResponseMessage
		switch request := messagee.Request.(type) {
		case *cek.IntentRequest:
			switch request.Intent.Name {
			case "Clova.GuideIntent":
				response = cek.NewResponseBuilder().
					OutputSpeech(
						cek.NewOutputSpeechBuilder().
							AddSpeechText("話しかけてください", cek.SpeechInfoLangJA).
							Build()).
					Build()
			}
		case *cek.LaunchRequest:
			response = cek.NewResponseBuilder().
				OutputSpeech(
					cek.NewOutputSpeechBuilder().
						AddSpeechText("起動しました", cek.SpeechInfoLangJA).
						Build()).
				Build()
		}
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```


## LICENSE

See [LICENSE.txt](LICENSE.txt)
