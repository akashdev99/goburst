package cutefmt

import "fmt"

//Pallete from https://github.com/zakaria-chahboun/cute/blob/main/cute.go

type CuteColor string

const (
	ResetColor   CuteColor = "\033[0m"
	DefaultColor CuteColor = "\033[39m"

	Black  CuteColor = "\033[30m"
	Red    CuteColor = "\033[31m"
	Green  CuteColor = "\033[32m"
	Yellow CuteColor = "\033[33m"
	Blue   CuteColor = "\033[34m"
	Purple CuteColor = "\033[35m"
	Cyan   CuteColor = "\033[36m"
	White  CuteColor = "\033[37m"

	BrightBlack  CuteColor = "\033[90m"
	BrightRed    CuteColor = "\033[91m"
	BrightGreen  CuteColor = "\033[92m"
	BrightYellow CuteColor = "\033[93m"
	BrightBlue   CuteColor = "\033[94m"
	BrightPurple CuteColor = "\033[95m"
	BrightCyan   CuteColor = "\033[96m"
	BrightWhite  CuteColor = "\033[97m"
)

func Errorf(a ...any) (n int, err error) {
	fmt.Printf("%v", BrightRed)
	fmt.Printf("❌ ")
	fmt.Println(a...)
	return fmt.Printf("%v", ResetColor)
}

func Successf(a ...any) (n int, err error) {
	fmt.Printf("%v", BrightGreen)
	fmt.Printf("✅ ")
	fmt.Println(a...)
	return fmt.Printf("%v", ResetColor)
}
