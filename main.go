package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

func main() {
	// Replace with your values
	subscriptionKey := "xxx"
	region := "switzerlandnorth"
	// text := `You are my sunshine, my only sunshine.
	// You make me happy when skies are gray.
	// You'll never know, dear, how much I love you.
	// Please don't take my sunshine away.`
	textEsp := `Eres mi sol, mi único sol.
	Me haces feliz cuando el cielo está gris.
	Nunca sabrás, querida, cuánto te amo.
	Por favor, no me quites mi sol.`

	// Create SSML eng
// 	ssmlEng := fmt.Sprintf(`
// <speak version='1.0' xml:lang='en-US'>
//   <voice xml:lang='en-US' xml:gender='Female' name='en-US-JessaNeural'>
//     %s
//   </voice>
// </speak>`, text)

	// Create SSML esp
	ssmlEsp := fmt.Sprintf(`
<speak version='1.0' xml:lang='es-ES'>
  <voice xml:lang='es-ES' xml:gender='Female' name='es-ES-HelenaNeural'>
    %s
  </voice>
</speak>`, textEsp)

	// Set up request
	client := resty.New()
	endpoint := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", region)

	resp, err := client.R().
		SetHeader("Ocp-Apim-Subscription-Key", subscriptionKey).
		SetHeader("Content-Type", "application/ssml+xml").
		SetHeader("X-Microsoft-OutputFormat", "audio-48khz-192kbitrate-mono-mp3").
		SetHeader("User-Agent", "go-azure-tts").
		SetBody(ssmlEsp).
		SetDoNotParseResponse(true). // important for binary response
		Post(endpoint)

	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != 200 {
		body, _ := io.ReadAll(resp.RawBody())
		log.Fatalf("Azure TTS failed: %d - %s", resp.StatusCode(), string(body))
	}

	// Save to file
	outFile := "output-es.mp3"
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.RawBody())
	if err != nil {
		log.Fatalf("Failed to write audio: %v", err)
	}

	fmt.Printf("Audio saved to %s\n", outFile)
}
