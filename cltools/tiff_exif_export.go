package cltools

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tacusci/clover/img"
	"github.com/tacusci/logging"
)

func RunTee(ts bool, sdir string, odir string, itype string, showExportOutput bool, overwrite bool, recursive bool) {
	if len(sdir) == 0 || len(odir) == 0 || len(itype) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Clover - Running TIFF EXIF export tool...\n")

	var st time.Time
	if ts {
		st = time.Now()
	}

	err := createDirectoryIfNotExists(odir)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	supportedInputTypes := []string{".nef"}
	supportedOutputTypes := []string{".jpg", ".png"}

	inputTypePrefixToMatch, inputType, err := parseInputOutputTypes(itype, "", supportedInputTypes, supportedOutputTypes)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	doneSearchingChan := make(chan bool, 32)
	imagesToExportExifChan := make(chan img.TiffImage, 32)

	if isDir, err := isDirectory(sdir); isDir {
		//file searching wait group
		var fswg sync.WaitGroup
		//images to export EXIF wait group
		var ieewg sync.WaitGroup
		//add a wait for the initial single call of 'findImagesInDir'
		fswg.Add(1)
		go findImagesInDir(&fswg, &imagesToExportExifChan, &doneSearchingChan, sdir, inputTypePrefixToMatch, inputType, recursive)
		ieewg.Add(1)
		go exportRawImageEXIF(&ieewg, &imagesToExportExifChan, &doneSearchingChan, itype, showExportOutput, overwrite, recursive, sdir, odir)
		//main thread doesn't wait after firing these goroutines, so force it to
		//wait until the file searching thread has finished
		fswg.Wait()
		//then tell the image exif export goroutine that there's no more images coming to export exif's
		doneSearchingChan <- true
		//wait on the image exif export goroutine until it's finished with all images it's already been working on
		ieewg.Wait()
		//both worker goroutines have finished, main thread continues
	} else {
		if err != nil {
			logging.ErrorAndExit(err.Error())
		}
	}

	if ts {
		logging.Info(fmt.Sprintf("Time taken: %d ms", time.Since(st).Nanoseconds()/1000000))
	}
}

func exportRawImageEXIF(wg *sync.WaitGroup, iteec *chan img.TiffImage, dsc *chan bool, itype string, showExportOutput bool, overwrite bool, recursive bool, sdir string, odir string) {
	for {
		if !<-*dsc {
			ri := <-*iteec
			wg.Add(1)
			if ri != nil {

			}
			wg.Done()
		} else {
			wg.Done()
		}
	}
}
