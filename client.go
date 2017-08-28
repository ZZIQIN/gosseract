package gosseract

import "fmt"
import "image"
import "os"
import "image/png"

// Client is an client to use gosseract functions
type Client struct {
	tesseract tesseractCmd
	source    path
	digest    path
	language  string
	Error     error
}
type path struct {
	value string
}

func (p *path) Ready() bool {
	return (p.value != "")
}
func (p *path) Get() string {
	return p.value
}

// NewClient provide reference to new Client
func NewClient() (c *Client, e error) {
	tess, e := getTesseractCmd()
	if e != nil {
		return
	}
	c = &Client{tesseract: tess}
	return
}

// Src accepts path to target source file
func (c *Client) Src(srcPath string) *Client {
	c.source = path{srcPath}
	return c
}

// Digest accepts path to target digest file
func (c *Client) Digest(digestPath string) *Client {
	c.digest = path{digestPath}
	return c
}

// Image accepts image object of target
func (c *Client) Image(img image.Image) *Client {
	imageFilePath, e := generateTmpFile()
	if e != nil {
		c.Error = e
		return c
	}
	f, e := os.Create(imageFilePath)
	// TODO: DRY
	if e != nil {
		c.Error = e
		return c
	}
	defer f.Close()
	png.Encode(f, img)
	// TODO: delete created file after
	c.source = path{f.Name()}
	return c
}

// Delete remove the temp file
func (c* Client) Delete() error{
	return os.Remove(c.source.Get())
}

// Out executes tesseract and gives results
func (c *Client) Out() (out string, e error) {
	if e = c.ready(); e != nil {
		return
	}
	// TODO: validation to call execute
	out, e = c.execute()
	return
}

// Must executes tesseract directly by parameter map
func (c *Client) Must(params map[string]string) (out string, e error) {
	if e = c.accept(params); e != nil {
		return
	}
	return c.Out()
}
func (c *Client) accept(params map[string]string) (e error) {
	src, ok := params["src"];
	if ok {
		c.source = path{src}
	}else if c.source.Get() == "" {
		return fmt.Errorf("Missing parameter `src`.")
	}

	if language, ok := params["language"]; ok {
		c.language = language
	}else{
		c.language = "eng"
	}

	if digest, ok := params["digest"]; ok {
		c.digest = path{digest}
	}
	return
}
func (c *Client) ready() (e error) {
	if !c.source.Ready() {
		return fmt.Errorf("Source is not set")
	}
	return
}
func (c *Client) execute() (res string, e error) {
	args := []string{
		c.source.Get(),
	}
	if c.language != "" {
		args = append(args, c.language)
	}
	if c.digest.Ready() {
		args = append(args, c.digest.Get())
	}
	return c.tesseract.Execute(args)
}
