package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	RunnerPath 	string 	= nil
	Preview 	bool 	= true
)

func main() {
	args := flag.Args()
	currentDir, _ := os.Getwd() + "/"

	fmt.Println(args)
	fmt.Println(currentDir)

	for args_i := 0; args_i < len(args); args_i++ {
		
		switch(args[args_i]) {
		case "--io-view-none":
			Preview = false;
		default:
			if args_i + 1 < len(args) {
				switch(args[args_i]) {
				case "--runner":
					args_i++
					if RunnerPath == nil {
						RunnerPath = currentDir + args[args_i]
						if err := checkExistFile(RunnerPath); err != nil {
							// Todo: Exception Manage "Error runnerPash isnot exist" (Exit 1)
						}
					} else {
						// Todo: Exception Manage "Error multiple-runner-path" (Exit 1)
					}
				}
			} else {
				switch (args[args_i]) {
				case "--runner": // Todo: Exception Manage "Error --runner don't have path of runner file" (Exit 1)
					break
				}
			}
		}
		
	}

    if (RunnerPath == nil) { 
        // Todo: Exception Manage "Error runnerPath isn't defined" (Exit 1) 
    }
	return 0
}

func checkExistFile(path string) -> error {

}
