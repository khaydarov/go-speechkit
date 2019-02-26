package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "github.com/pkg/errors"
   "io"
   "io/ioutil"
   "log"
   "net/http"
   "net/url"
   "os"
   "strings"
)

const speechkitPostApiUrl = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"
const requestTokenUrl = "https://iam.api.cloud.yandex.net/iam/v1/tokens"

// Cloud Token structure
type YandexCloudToken struct {
   IamToken string `json:"iamToken"`
}

// SpeechKitProcess sends to the Yandex SpeechKit API request to convert text to speech
func SpeechKitProcess(text string, file string, iamToken string, folderId string) (error) {
  form := url.Values{}
  form.Add("text", text)
  form.Add("lang", "ru-RU")
  form.Add("emotion", "neutral")
  form.Add("speed", "1.0")
  form.Add("voice", "zahar")
  form.Add("folderId", folderId)

  // Create a new request using http
  req, err := http.NewRequest("POST", speechkitPostApiUrl, strings.NewReader(form.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

  //add authorization header to the req
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", iamToken))

  // Send req using http Client
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      return errors.Wrap(err, "")
  }

  if resp.StatusCode == http.StatusOK {
      log.Printf("[YandexCloud Speechkit]: Status OK\n")
      defer resp.Body.Close()

      // Create the file
      oggFile, _ := os.Create(file)
      defer oggFile.Close()

      _, err = io.Copy(oggFile, resp.Body)
      return nil
  } else {
      return errors.Wrap(err, "[YandexCloud Speechkit]: Status is not OK")
  }
}

// GenerateKey is a Yandex cloud package method that generates iamtoken
func GenerateKey(token string) (string, error) {
    accessToken := token
    jsonStr := []byte(fmt.Sprintf(`{"yandexPassportOauthToken":"%s"}`, accessToken))

    req, err := http.NewRequest("POST", requestTokenUrl, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    // Send req using http Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", errors.Wrap(err, "Can't send request")
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", errors.Wrap(err, "Can't read response body")
    }

    var cloudToken YandexCloudToken
    err = json.Unmarshal(body, &cloudToken)

    if err != nil {
        return "", errors.Wrap(err, "Can't parse response body")
    }

    if cloudToken.IamToken == "" {
        return "", errors.Errorf("Wrong access token")
    }

    return cloudToken.IamToken, nil
}