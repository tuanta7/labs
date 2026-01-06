package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/option"
)

func main() {
	action := "submit"

	ctx := context.Background()
	client, err := recaptcha.NewClient(ctx, option.WithAPIKey(os.Getenv("GCP_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating reCAPTCHA client %v", err)
	}
	defer client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Token string `json:"token"`
		}

		json.NewDecoder(r.Body).Decode(&data)

		// fmt.Println(data.Token)

		event := &recaptchapb.Event{
			Token:          data.Token,
			SiteKey:        os.Getenv("RECAPTCHA_SITE_KEY"),
			ExpectedAction: action,
		}

		assessment := &recaptchapb.Assessment{
			Event: event,
		}

		request := &recaptchapb.CreateAssessmentRequest{
			Assessment: assessment,
			Parent:     fmt.Sprintf("projects/%s", os.Getenv("PROJECT_ID")),
		}

		response, err := client.CreateAssessment(ctx, request)
		if err != nil {
			fmt.Printf("%v", err.Error())
		}

		jsonResp, _ := json.Marshal(response)
		w.Write(jsonResp)
	}).Methods("POST")

	srv := &http.Server{
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{"http://localhost:5500"}),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
		Addr: "localhost:8000",
	}

	fmt.Println("Starting server at localhost:8000")
	srv.ListenAndServe()
}
