package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/lukesampson/figlet/figletlib"
	"golang.org/x/term"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	width, _, err = term.GetSize(0)
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ErrorCheck(err error) {
	if err != nil {
		msg := strings.Trim(err.Error(), "\n")
		fmt.Println(Pretty(msg))
	}
}

func NCenter(width int, s string) *bytes.Buffer {
	const half, space = 2, "\u0020"
	var b bytes.Buffer
	n := (width - utf8.RuneCountInString(s)) / half
	fmt.Fprintf(&b, "%s%s", strings.Repeat(space, n), s)
	return &b
}

func AppendFile(text string, name string) {

	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	ErrorCheck(err)

	defer f.Close()

	_, err = f.WriteString(text + "\n")
	ErrorCheck(err)

}

func Pretty(info string) string {
	pretty := ""
	pretty += color.HiMagentaString("[")
	pretty += color.WhiteString("+")
	pretty += color.HiMagentaString("] ")
	pretty += info
	return pretty
}

func Valid(info string) string {
	pretty := ""
	pretty += color.HiGreenString("[")
	pretty += color.WhiteString("+")
	pretty += color.HiGreenString("] ")
	pretty += info
	return pretty
}

func Invalid(info string) string {
	pretty := ""
	pretty += color.HiRedString("[")
	pretty += color.WhiteString("+")
	pretty += color.HiRedString("] ")
	pretty += info
	return pretty
}

func Clear() {
	screen.Clear()
}

func Border() {
	i := 0
	res1 := ""
	for i < width {
		res1 += "─"
		i += 1
	}
	fmt.Println(res1)
}

func Logo() {
	cwd, _ := os.Getwd()
	fontsdir := filepath.Join(cwd, "data")
	ErrorCheck(err)
	f, err := figletlib.GetFontByName(fontsdir, "4max")
	ErrorCheck(err)
	color.Set(color.FgHiMagenta)
	figletlib.PrintMsg("Kyanite", f, width, f.Settings(), "center")
	color.Set(color.FgHiWhite)
	fmt.Println()
	fmt.Println()
	Border()
}

type Profile struct {
	Username string `json:"username"`
	Discrim  string `json:"discriminator"`
}

func Check(token string, wg *sync.WaitGroup) {
	defer wg.Done()
	var data Profile

	if len(token) < 60 {
		fmt.Println(Invalid("- Invalid Token: " + token))
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	ErrorCheck(err)

	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	ErrorCheck(err)

	Works := strings.HasPrefix(resp.Status, "2")

	if Works {
		stonedeagleify, err := io.ReadAll(resp.Body)
		ErrorCheck(err)

		if err := json.Unmarshal(stonedeagleify, &data); err != nil {
			fmt.Println("failed to unmarshal:", err)
		} else {
			AppendFile(token, "tokens/valid_nousernames.txt")
			AppendFile(token+" - "+data.Username+"#"+data.Discrim, "tokens/valid_username.txt")
		}
		fmt.Println(Valid("- Valid Token: " + token + " - " + data.Username + "#" + data.Discrim))
	} else {
		AppendFile(token, "tokens/invalid.txt")
		fmt.Println(Invalid("- Invalid Token: " + token))
	}

}

func Installation() {
	errr := os.Mkdir("data", os.ModePerm)
	err := os.Mkdir("tokens", os.ModePerm)
	ErrorCheck(err)
	ErrorCheck(errr)
	client := &http.Client{}
	ErrorCheck(err)
	req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/Kyxnite/stuff/main/GO/Fonts/4max.flf", nil)
	resp, err := client.Do(req)
	stonedeagleify, err := io.ReadAll(resp.Body)
	AppendFile(string(stonedeagleify), "data/4max.flf")
	ErrorCheck(err)
}

func main() {
	Clear()
	Installation()
	fmt.Println(NCenter(width, ("- The checking is not as slow as it seems in console")))
	fmt.Println(NCenter(width, ("- The tokens are checked WAY before they are printed")))
	time.Sleep(3 * time.Second)
	var wg sync.WaitGroup
	Clear()
	Logo()
	lines, err := readLines("tokens.txt")
	ErrorCheck(err)
	for sex, token := range lines {
		wg.Add(sex)
		go Check(token, &wg)
	}
	wg.Wait()
}
