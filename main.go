package main

import (
    "bytes"
    "fmt"
    "github.com/fatih/color"
    "gopkg.in/abiosoft/ishell.v2"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "strconv"
    "unicode/utf8"
)

const maxLen = 1000

func main() {
    var token string
    var folderId string

    // create new shell.
    // by default, new shell includes 'exit', 'help' and 'clear' commands.
    shell := ishell.New()

    // display welcome info.
    yellow := color.New(color.FgYellow).SprintFunc()
    shell.Println(yellow("Welcome to the go-speechKit cli"))

    // register a function for "init" command.
    shell.AddCmd(&ishell.Cmd{
        Name: "init",
        Help: "initialize go-speechKit",
        Func: func(c *ishell.Context) {
            c.Print("Please enter access token: ")
            cmdString := c.ReadLine()

            var err error
            token, err = GenerateKey(cmdString)
            if err != nil {
                c.Println("Error. Try again")
                return
            }

            c.Print("Please enter folderId: ")
            folderId = c.ReadLine()

            c.Println("go-speechKit is initialized. Type «process» to continue")
        },
    })

    // register a function for "process" command.
    shell.AddCmd(&ishell.Cmd{
        Name: "process",
        Help: "processing",
        Func: func(c *ishell.Context) {
            c.Print("Please enter path to text file: ")
            inputFilePath := c.ReadLine()

            textFile, err := ioutil.ReadFile(inputFilePath)
            if err != nil {
                c.Printf("Can't find file %s\n", inputFilePath)
                return
            }

            if string(textFile) == "" {
                c.Printf("File must contain string\n")
                return
            }

            c.Print("Please enter path to output file: ")
            outputMp3 := c.ReadLine()
            if string(outputMp3) == "" {
                c.Printf("String is empty\n")
            }

            c.Println("Wait a minute...")
            process(string(textFile), string(outputMp3), token, folderId)

            c.Println("Done!")
        },
    })

    // run shell
    shell.Run()
}

// process is a function that calls YandeX SpeechKit API to generate ogg file
// and converts via ffmpeg to the mp3
// First it splits the text before sending to the API
// For each text-blocks ogg file is generated
// Then merge every chunk and convert to final mp3
func process(textFile string, outputMp3 string, token string, folderId string) {
    textChunks := SplitText(textFile)
    outputTxt, err := os.Create("output.txt")
    defer outputTxt.Close()

    var chunks []string
    for index, text := range textChunks {
        oggChunk := fmt.Sprintf("audio%s.ogg", strconv.Itoa(index))
        err = SpeechKitProcess(text, oggChunk, token, folderId)
        if err != nil {
            log.Fatalf("SpeechkitProcess error: %s", err)
        }

        chunks = append(chunks, oggChunk)
        outputTxt.WriteString(fmt.Sprintf("file '%s'\n", oggChunk))
    }

    log.Println("[go-speechKit]: Converting...")
    cmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", "output.txt", "-vn", "-ar", "44100", "-ac", "2", "-ab", "192k", "-f", "mp3", outputMp3)

    // Debug exec commands
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err = cmd.Run()

    // remote output.txt file
    os.Remove("output.txt")

    // remove chunks
    for _, chunk := range chunks {
        os.Remove(chunk)
    }
}

func SplitText(longString string) []string {
    splits := []string{}

    var l, r int
    for l, r = 0, maxLen; r < len(longString); l, r = r, r+maxLen {
        for !utf8.RuneStart(longString[r]) {
            r--
        }
        splits = append(splits, longString[l:r])
    }
    splits = append(splits, longString[l:])
    return splits
}