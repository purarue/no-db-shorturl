package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// number of times to retry generating a random value for a specific length
// before incrementing the hash length to generate
const randomRetryAmount = 10

// configuration information
type config struct {
	port       int
	dataFolder string
	secretKey  string
}

// https://stackoverflow.com/a/22892986
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generates a hash value that doesn't already exist in ./data
func generateHashValue() string {
	hashLength := 4
	incrementHashCount := randomRetryAmount
	var hashStr string
	for {
		hashStr = randSeq(hashLength)
		if _, err := os.Stat(hashStr); err != nil {
			break
		}
		incrementHashCount -= 1
		if incrementHashCount == 0 {
			incrementHashCount = randomRetryAmount
			hashLength += 1
		}
	}
	return hashStr
}

func parseFlags() *config {
	// flag definitions
	port := flag.Int("port", 8040, "port to serve shorturl on")
	dataFolderP := flag.String("data-folder", "./data", "directory to store data in")
	secretKeyP := flag.String("secret-key", "", "secret key to authenticate POST requests")
	// parse flags
	flag.Parse()
	// make sure path is valid
	dataFolder := strings.TrimSpace(*dataFolderP)
	secretKey := strings.TrimSpace(*secretKeyP)
	_, err := os.Stat(dataFolder)
	if os.IsNotExist(err) {
		log.Printf("%s does not exist, creating...\n", dataFolder)
		mode := int(0777)
		os.Mkdir(dataFolder, os.FileMode(mode))
	} else if err != nil {
		panic(err)
	}
	if len(secretKey) == 0 {
		secretKey = os.Getenv("SHORTURL_KEY")
		if len(secretKey) == 0 {
			log.Println("Warning: no -secret-key flag or SHORTURL_KEY environment variable provided, anyone is able to create short URLs")
		}
	}
	return &config{
		port:       *port,
		dataFolder: dataFolder,
		secretKey:  secretKey,
	}
}

type postInfo struct {
	SecretKey string `json:"key"`
	Url       string `json:"url"`
	Hash      string `json:"hash"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	config := parseFlags()
	err := os.Chdir(config.dataFolder)
	if err != nil {
		panic(err)
	}
	// global handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			switch r.Method {
			case "POST":
				// Create a new shortened URL
				decoder := json.NewDecoder(r.Body)
				var postData postInfo
				err := decoder.Decode(&postData)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error: decoding POST body into JSON: %v\n", err)
					return
				}
				// fmt.Printf("%v\n", postData)
				passedKey := postData.SecretKey
				passedUrl := strings.TrimSpace(postData.Url)
				filePath := strings.TrimSpace(postData.Hash)
				// Error Checking
				if passedKey != config.secretKey {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "Error: incorrect secret key\n")
					return
				}
				if passedUrl == "" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Error: You didn't provide a url to redirect to\n")
					return
				}
				if filePath == "" {
					filePath = generateHashValue()
				}
				// Create URL
				err = os.WriteFile(filePath, []byte(passedUrl), 0644)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error: Couldn't create shorturl file: %v\n", err)
					return
				} else {
					fmt.Fprintf(w, "%s\n", filePath)
					return
				}
			default:
				fmt.Fprintf(w, "%s", `The base endpoint not using a POST request does nothing.
To add a shortened URL, make a POST request to this endpoint. Example:

curl --header "Content-Type: application/json" --request POST --data '{"key":"your_secret_key","url":"https://sean.fish"}' http://localhost:8040

or to specify the path to create the shorturl on:

curl --header "Content-Type: application/json" --request POST --data '{"key":"your_secret_key","url":"https://sean.fish","hash":"short"}' http://localhost:8040

For more info see https://github.com/purarue/no-db-shorturl
`)
				return
			}
		} else {
			trimmedUrl := strings.Trim(r.URL.Path, "/")
			if strings.Contains(trimmedUrl, "/") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Error: URL shouldn't contain '/'\n")
				return
			}
			// redirect user or 404
			if _, err := os.Stat(trimmedUrl); err == nil {
				// shorturl file exists, serve it
				contents, _ := os.ReadFile(trimmedUrl)
				http.Redirect(w, r, strings.TrimSpace(string(contents)), 302)
				return
			} else {
				// shorturl doesn't exist
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Error: shorturl %s doesn't exist\n", r.URL.Path)
				return
			}
		}
	})
	log.Printf("shorturl serving on port %d at %s", config.port, config.dataFolder)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.port), nil))
}
