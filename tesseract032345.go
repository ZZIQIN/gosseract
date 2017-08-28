package gosseract

import "fmt"
import "os/exec"
import "bytes"

type tesseract2345 struct {
	version        string
	resultFilePath string
	commandPath    string
}

func (t tesseract2345) Version() string {
	return t.version
}
func (t tesseract2345) Execute(params []string) (res string, e error) {

	// command args
	var args []string
	// Register source file
	args = append(args, params[0],"stdout","-l",params[1])

	// Register digest file
	if len(params) > 2 {
		args = append(args, params[2])
	}

	// prepare command
	fmt.Print(args)
	cmd := exec.Command(TESSERACT, args...)
	// execute
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	if e = cmd.Run(); e != nil {
		e = fmt.Errorf(stderr.String())
		return
	}
	// read result
	res = out.String()
	return
}

