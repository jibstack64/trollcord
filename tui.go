package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// emoji stuff
const (
	SUCCESS   = "✅"
	ERROR     = "❌"
	EXTRASPIN = 4
)

// colours
var (
	PromptColour  = color.New(color.FgBlue)
	OptionColour  = color.New(color.FgCyan)
	BackColour    = color.New(color.FgHiBlack)
	SuccessColour = color.New(color.FgHiGreen)
	ErrorColour   = color.New(color.FgRed)
	FatalColour   = color.New(color.FgRed, color.Underline)
	MessageColour = color.New(color.FgHiMagenta)
)

func title() string {
	return "\n ▄▄█████▓ ██▀███   ▒█████   ██▓     ██▓     ▄████▄   ▒█████   ██▀███  ▓█████▄ \n ▓  ██▒ ▓▒▓██ ▒ ██▒▒██▒  ██▒▓██▒    ▓██▒    ▒██▀ ▀█  ▒██▒  ██▒▓██ ▒ ██▒▒██▀ ██▌\n ▒ ▓██░ ▒░▓██ ░▄█ ▒▒██░  ██▒▒██░    ▒██░    ▒▓█    ▄ ▒██░  ██▒▓██ ░▄█ ▒░██   █▌\n ░ ▓██▓ ░ ▒██▀▀█▄  ▒██   ██░▒██░    ▒██░    ▒▓▓▄ ▄██▒▒██   ██░▒██▀▀█▄  ░▓█▄   ▌\n   ▒██▒ ░ ░██▓ ▒██▒░ ████▓▒░░██████▒░██████▒▒ ▓███▀ ░░ ████▓▒░░██▓ ▒██▒░▒████▓ \n   ▒ ░░   ░ ▒▓ ░▒▓░░ ▒░▒░▒░ ░ ▒░▓  ░░ ▒░▓  ░░ ░▒ ▒  ░░ ▒░▒░▒░ ░ ▒▓ ░▒▓░ ▒▒▓  ▒ \n     ░      ░▒ ░ ▒░  ░ ▒ ▒░ ░ ░ ▒  ░░ ░ ▒  ░  ░  ▒     ░ ▒ ▒░   ░▒ ░ ▒░ ░ ▒  ▒ \n   ░        ░░   ░ ░ ░ ░ ▒    ░ ░     ░ ░   ░        ░ ░ ░ ▒    ░░   ░  ░ ░  ░ \n             ░         ░ ░      ░  ░    ░  ░░ ░          ░ ░     ░        ░    \n                                            ░                           ░      \n		    https://github.com/jibstack64/trollcord                   v0.0.3"
}

// prompts the user with 'prompt' and returns their response.
// if 'retry' is true, then the user will be prompted until a response is given.
func getInput(prompt string, retry bool, retryString *string) string {
	PromptColour.Println("\n" + prompt)
	for {
		BackColour.Print("> ")
		var input string
		inputReader := bufio.NewReader(os.Stdin)
		input, _ = inputReader.ReadString('\n')
		input = strings.ReplaceAll(input, "\n", "")
		if retry && input == "" {
			if retryString != nil {
				ErrorColour.Println(*retryString)
			}
			continue
		} else {
			return input
		}
	}
}

// clears 'n' lines.
func clearLine(n int) {
	for c := 0; c < n; c++ {
		fmt.Print("\033[1A\033[K")
	}
}

// creates a loading string.
// when finished becomes true, it stops.
func loading(text string, fn func(finished *bool, err *error)) error {
	fmt.Print("\n")
	chars := []string{
		"\\", "/", "-", "\\", "/", "-",
	}
	c := 0
	extra := 0
	// track progress
	var finished bool
	var err error
	go fn(&finished, &err)
	for {
		fmt.Printf("%s %s\n", MessageColour.Sprint(chars[c]), text)
		time.Sleep(time.Second / 4)
		clearLine(1) // clear line
		c++
		if c == len(chars) {
			c = 0
		}
		if finished {
			if extra <= EXTRASPIN {
				extra++
			} else {
				break
			}
		}
	}
	// cool!
	if err != nil {
		fmt.Printf("%s %s\n", ERROR, ErrorColour.Sprint(text))
	} else {
		fmt.Printf("%s %s\n", SUCCESS, SuccessColour.Sprint(text))
	}
	return err
}

// prompts the user with 'prompt' and awaits their input.
// returns a boolean, yes being true and no being false.
func yesOrNo(prompt string) bool {
	inputString := fmt.Sprintf("%s%s%s%s%s",
		BackColour.Sprint("("), SuccessColour.Sprint("yes"), BackColour.Sprint("/"),
		ErrorColour.Sprint("no"), BackColour.Sprint(")"))
	PromptColour.Printf("\n%s %s\n", prompt, inputString)
	for {
		BackColour.Print("> ")
		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(input) // lower
		if input == "yes" {
			return true
		} else if input == "no" {
			return false
		} else {
			continue
		}
	}
}

// creates a progress bar that is updated according to 'length' and 'done'.
func progressBar(prompt string, fn func(length *int, done *int, err *error)) error {
	fmt.Printf("\n%s\n", MessageColour.Sprint(prompt))
	progress := "[%s]"
	check := "#"
	back := "~"
	// track progress
	var length int
	var done int
	var err error
	go fn(&length, &done, &err)
	for {
		var checks string
		donevl := 0
		if done > 0 {
			donevl = done / length * 20
		}
		for i := 0; i < donevl; i++ {
			checks += check
		}
		for i := 0; i < 20-donevl; i++ {
			checks += back
		}
		fmt.Printf("%s\n", fmt.Sprintf(progress, checks))
		time.Sleep(time.Second / 4)
		clearLine(1) // clear line
		if done <= length || err != nil {
			break
		}
	}
	// cool!
	if err != nil {
		fmt.Printf("[%s]\n", ErrorColour.Sprint("~~~~~~~~~"+ERROR+"~~~~~~~~~"))
	} else {
		fmt.Printf("[%s]\n", SuccessColour.Sprint("~~~~~~~~~"+SUCCESS+"~~~~~~~~~"))
	}
	return err
}

// prompts the user with the given 'prompt' string.
// lists all options in 'options' as 'index+1: option'.
// returns selected option's index.
func fromSelection(prompt string, options []string) int {
	PromptColour.Println("\n" + prompt)
	for p, v := range options {
		if v == "" {
			continue
		}
		fmt.Printf("%s%s %s\n", OptionColour.Sprint(p+1), BackColour.Sprint("."), v)
	}
	for {
		BackColour.Print("> ")
		var input string
		fmt.Scanln(&input)
		index, err := strconv.Atoi(input)
		index -= 1
		if err != nil {
			ErrorColour.Printf("'%s' is not an integer.\n", input)
			continue
		}
		for p, v := range options {
			if index == p && v != "" {
				return p
			}
		}
		if input != "" {
			ErrorColour.Printf("%s is not a valid option.\n", input)
		}
	}
}

// prints 'text' in an underlined, red error text.
func fatal(text string) {
	FatalColour.Println(text)
}

// prints 'text' in a red error text.
func errorPr(text string) {
	ErrorColour.Println(text)
}

// prints 'text' in a light magenta tone.
func message(text string) {
	MessageColour.Println(text)
}

// prints 'text' in a light green tone.
func success(text string) {
	SuccessColour.Println(text)
}
