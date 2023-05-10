package main

//func UploadCustomersAsync() error {
//
//	// declare variables
//	var timer time.Time
//	var global time.Time
//	var loopNumbers int
//	var results = make(chan error)
//	var url = fmt.Sprintf("%s/api/customers/customer_import", settings.DestinationHost)
//
//	// print start time
//	fmt.Printf("Start time: %s\n", time.Now().Format(time.ANSIC))
//
//	// search file(s) to load
//	files, _ := filepath.Glob(fmt.Sprintf("%s_?????%s", strings.TrimSuffix(FileCustomers, filepath.Ext(FileCustomers)), filepath.Ext(FileCustomers)))
//	if len(files) == 0 {
//		files, _ = filepath.Glob(FileCustomers)
//	}
//	if len(files) == 0 {
//		return fmt.Errorf("files to loading not found")
//	}
//
//	// start main loop
//	global = time.Now()
//	for _, file := range files {
//
//		// load list of customers from file
//		fmt.Printf("Loading list of customers from file '%s'\n", file)
//		var customers []Customer
//		err := ObjectRead(&customers, file)
//		if err != nil {
//			return err
//		}
//		fmt.Printf("Customers loaded: %d\n", len(customers))
//
//		// calculate loops number
//		loopNumbers = len(customers) / settings.PacketSize
//		if len(customers)%settings.PacketSize != 0 {
//			loopNumbers++
//		}
//
//		// start upload loop
//		for i := 0; i < loopNumbers; i++ {
//
//			//
//			timer = time.Now()
//
//			//
//			end := settings.PacketSize
//			if i+1 == loopNumbers {
//				end = len(customers) % settings.PacketSize
//			}
//
//			// call pool of request
//			for j := 0; j < end; j++ {
//				var customer, _ = json.MarshalIndent(customers[i*settings.PacketSize+j], "", "\t")
//				go ExecRequest2(url, string(customer), results)
//			}
//
//			// await result
//			for j := 0; j < end; j++ {
//				err = <-results
//				if err != nil {
//					fmt.Printf("Error: %v", err)
//				}
//			}
//
//			fmt.Printf("\rCyrcle %09d. Uploaded customers from %09d to %09d. Time: %05d ms, total: %05.2f min, average: %05.2f objects/second",
//				i+1, i*settings.PacketSize+1, i*settings.PacketSize+settings.PacketSize, time.Since(timer).Milliseconds(),
//				time.Since(global).Minutes(), float64(i*settings.PacketSize+settings.PacketSize)/time.Since(global).Seconds())
//		}
//		fmt.Printf("\n")
//	}
//
//	fmt.Printf("\nFinish time: %s\n", time.Now().Format(time.ANSIC))
//	return nil
//}
