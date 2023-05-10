package main

import (
	"fmt"
	_ "go/types"
	"path/filepath"
	"strings"
	"time"
)

type ObjectType string

const (
	ObjectTypeCards     ObjectType = "DiscountCards"
	ObjectTypeAccounts  ObjectType = "Accounts"
	ObjectTypeCustomers ObjectType = "Customers"
)

func UploadObjects(url string, filename string) error {

	// declare variables
	var timer time.Time
	var filetimer time.Time
	var global time.Time
	var loopNumbers int
	var result = make(chan *ErrorInfo)
	var errors []ErrorInfo
	fmt.Printf("Use url for loading: %s\n", url)
	fmt.Printf("Use packet size: %d\n", settings.PacketSize)

	// print start time
	fmt.Printf("Time start: %s\n", time.Now().Format(time.ANSIC))
	global = time.Now()

	// search file(s) to load
	files, _ := filepath.Glob(fmt.Sprintf("%s_?????%s", strings.TrimSuffix(filename, filepath.Ext(filename)), filepath.Ext(filename)))
	if len(files) == 0 {
		files, _ = filepath.Glob(filename)
	}
	if len(files) == 0 {
		return fmt.Errorf("files to loading not found")
	}
	fmt.Printf("Found files to loading: %03d\n", len(files))

	// uploading loop
	for _, file := range files {

		// start file timer
		filetimer = time.Now()

		// load list of objects from file
		timer = time.Now()
		fmt.Printf("\nLoading list of objects from file '%s'\n", file)
		var objs []any
		err := ObjectRead(&objs, file)
		if err != nil {
			return err
		}
		fmt.Printf("Loaded %07d objects in %03.02f seconds\n", len(objs), time.Since(timer).Seconds())

		// calculate loop numbers
		loopNumbers = len(objs) / settings.PacketSize
		if len(objs)%settings.PacketSize != 0 {
			loopNumbers++
		}

		// run upload loop
		for i := 0; i < loopNumbers; i++ {

			// start loop timer
			timer = time.Now()

			// executing parallel request
			end := settings.PacketSize
			if i+1 == loopNumbers {
				end = len(objs) % settings.PacketSize
			}
			for j := 0; j < end; j++ {
				go ExecRequest(url, objs[i*settings.PacketSize+j], result)
			}

			// waiting for result
			for j := 0; j < end; j++ {
				errR := <-result
				if errR != nil {
					if errR.Error != "" && errR.Object != "" {
						errors = append(errors, *errR)
					}
				}
			}

			// print results
			fmt.Printf("\rLoop %06d/%06d. Uploaded objects %06d - %06d. Time: %05.4f s, total: %05.4f m, aver: %03.0f objs/sec",
				i+1, loopNumbers, i*settings.PacketSize+1, i*settings.PacketSize+end, time.Since(timer).Seconds(),
				time.Since(filetimer).Minutes(), float64(i*settings.PacketSize+settings.PacketSize)/time.Since(filetimer).Seconds())

		}
		fmt.Printf("\n")

		//
		fmt.Printf("Objects from file %s uploaded in %04.02f minutes, %05.02f objects/second\n",
			file, time.Since(filetimer).Minutes(), float64(len(objs))/time.Since(filetimer).Seconds())

	}

	fmt.Printf("Upload completed in %04.02f minutes\n", time.Since(global).Minutes())
	if len(errors) > 0 {
		fmt.Printf("Got uploading error number: %d\n", len(errors))
		err := WriteObject(errors, "errors.log")
		if err != nil {
			fmt.Printf("%v", err)
		}
	}

	fmt.Printf("Time finish: %s\n", time.Now().Format(time.ANSIC))
	return nil
}
