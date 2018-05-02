package main

import (
	"os"
	"flag"
	"github.com/tacusci/clover/cltools"
	"fmt"
)

func OutputUsage() {
	println("Usage: " + os.Args[0] + " </TOOLFLAG>")
	fmt.Printf("\t/sdc (StorageDeviceChecker) - Tool for checking size of storage devices.\n")
	fmt.Printf("\t/rtc (RawToCompressed) - Tool for batch compressing raw images.\n")
}

func OutputUsageAndClose() {
	OutputUsage()
	os.Exit(1)
}

func main() {

	if len(os.Args) == 1 {
		OutputUsageAndClose()
	}

	RunTool(os.Args[1])
}

func RunTool(toolFlag string) {
	if toolFlag == "/sdc" {
		locationPath := flag.String("location", "", "Location to write data to.")
		sizeToWrite := flag.Int("size", 0, "Size of total data to write.")
		skipFileIntegrityCheck := flag.Bool("skip-integrity-check", false, "Skip verifying output file integrity.")
		dontDeleteFiles := flag.Bool("no-delete", false, "Don't delete outputted files.")

		//kind of hack to force flag parser to find tool argument flags correctly
		os.Args = os.Args[1:]

		flag.Parse()

		cltools.RunSdc(*locationPath, *sizeToWrite, *skipFileIntegrityCheck, *dontDeleteFiles)

	} else if toolFlag == "/rtc" {
		sourceDirectory := flag.String("source-directory", "", "Location containing raw images to convert.")
		inputType := flag.String("input-type", "", "Extension of image type to convert.")

		//kind of hack to force flag parser to find tool argument flags correctly
		os.Args = os.Args[1:]

		flag.Parse()

		cltools.RunRtc(*sourceDirectory, *inputType, ".jpg")
	}  else {
		OutputUsageAndClose()
	}
}