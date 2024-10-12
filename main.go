package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	pb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

var (
	text   string
	output string

	langCode string
	voice    string
)

func init() {
	flag.StringVar(&text, "text", "Hello World", "text line for say something")
	flag.StringVar(&output, "output", "./speech-voice.wav", "output path & filename")

	flag.StringVar(&langCode, "lang", "ja-JP", "language code")
	flag.StringVar(&voice, "voice", "jp-JP-Standard-A", "voice name")

	flag.Parse()
}

func main() {
	fmt.Printf("text line: %s\n", text)
	fmt.Printf("output path & filename: %s\n", output)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v\n", err)
		return
	}
	defer client.Close()

	res, err := client.SynthesizeSpeech(ctx, &pb.SynthesizeSpeechRequest{
		Input: &pb.SynthesisInput{
			InputSource: &pb.SynthesisInput_Text{Text: text},
		},
		Voice: &pb.VoiceSelectionParams{
			LanguageCode: langCode,
			Name:         voice,
			SsmlGender:   pb.SsmlVoiceGender_MALE,
		},
		AudioConfig: &pb.AudioConfig{
			AudioEncoding: pb.AudioEncoding_LINEAR16,
		},
	})
	if err != nil {
		log.Printf("Failed to synthesize speech: %v\n", err)
		return
	}

	// save audio content to file
	err = saveAudioToFile(res.AudioContent, output)
	if err != nil {
		log.Printf("Failed to save audio content to file: %v\n", err)
		return
	}

}

func saveAudioToFile(audioContent []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(audioContent)
	if err != nil {
		return err
	}

	return nil
}
