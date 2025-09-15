package colors

import (
	"github.com/fatih/color"
)

var (
	Success   = color.New(color.FgGreen, color.Bold)
	Error     = color.New(color.FgRed, color.Bold)
	Info      = color.New(color.FgCyan, color.Bold)
	Prompt    = color.New(color.FgBlue, color.Bold)
	Generated = color.New(color.FgMagenta, color.Bold)
	Title     = color.New(color.FgWhite, color.Bold)
	Subtle    = color.New(color.FgHiBlack)
)

func SuccessMsg(msg string) string {
	return Success.Sprint(msg)
}

func ErrorMsg(msg string) string {
	return Error.Sprint(msg)
}

func InfoMsg(msg string) string {
	return Info.Sprint(msg)
}

func PromptMsg(msg string) string {
	return Prompt.Sprint(msg)
}

func GeneratedMsg(msg string) string {
	return Generated.Sprint(msg)
}

func TitleMsg(msg string) string {
	return Title.Sprint(msg)
}

func SubtleMsg(msg string) string {
	return Subtle.Sprint(msg)
}
