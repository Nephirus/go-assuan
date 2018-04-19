package pinentry

import (
	"os/exec"
	"strconv"
	"time"

	assuan "github.com/foxcpp/go-assuan/client"
)

type Client struct {
	Session *assuan.Session

	current    Settings
	qualityBar bool
}

func Launch() (*Client, error) {
	cmd := exec.Command("pinentry")

	c := new(Client)
	var err error
	c.Session, err = assuan.InitCmd(cmd)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func LaunchCustom(path string) (Client, error) {
	cmd := exec.Command(path)

	c := Client{}
	var err error
	c.Session, err = assuan.InitCmd(cmd)
	if err != nil {
		return Client{}, err
	}
	return c, nil
}

func (c *Client) Shutdown() {
	c.Session.Close()
}

func (c *Client) Reset() {
	c.Session.Reset()
}

func (c *Client) SetDesc(text string) {
	c.Session.SimpleCmd("SETDESC", text)
	c.current.Desc = text
}

func (c *Client) SetPrompt(text string) {
	c.Session.SimpleCmd("SETPROMPT", text)
	c.current.Prompt = text
}

func (c *Client) SetError(text string) {
	c.Session.SimpleCmd("SETERROR", text)
	c.current.Error = text
}

func (c Client) SetOkBtn(text string) {
	c.Session.SimpleCmd("SETOK", text)
	c.current.OkBtn = text
}

func (c Client) SetNotOkBtn(text string) {
	c.Session.SimpleCmd("SETNOTOK", text)
	c.current.NotOkBtn = text
}

func (c Client) SetCancelBtn(text string) {
	c.Session.SimpleCmd("SETCANCEL", text)
	c.current.CancelBtn = text
}

func (c Client) SetTitle(text string) {
	c.Session.SimpleCmd("SETTITLE", text)
	c.current.Title = text
}

func (c Client) SetTimeout(timeout time.Duration) {
	c.Session.SimpleCmd("SETTIMEOUT", strconv.Itoa(int(timeout.Seconds())))
	c.current.Timeout = timeout
}

func (c Client) SetRepeatPrompt(text string) {
	c.Session.SimpleCmd("SETREPEAT", text)
	c.current.RepeatPrompt = text
}

func (c Client) SetRepeatError(text string) {
	c.Session.SimpleCmd("SETREPEATERROR", text)
	c.current.RepeatError = text
}

func (c Client) Current() Settings {
	return c.current
}

func (c Client) Apply(s Settings) {
	c.SetDesc(s.Desc)
	c.SetPrompt(s.Prompt)
	c.SetError(s.Error)
	c.SetOkBtn(s.OkBtn)
	c.SetNotOkBtn(s.NotOkBtn)
	c.SetCancelBtn(s.CancelBtn)
	c.SetTitle(s.Title)
	c.SetTimeout(s.Timeout)
	c.SetRepeatPrompt(s.RepeatPrompt)
	c.SetRepeatError(s.RepeatError)
}

// GetPIN shows window with password textbox, Cancel and Ok buttons.
// Error is returned if Cancel is pressed.
func (c Client) GetPIN() (string, error) {
	dat, err := c.Session.SimpleCmd("GETPIN", "")
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// Confirm shows window with Cancel and Ok buttons but without password
// textbox, error is returned if Cancel is pressed (as usual).
func (c Client) Confirm() error {
	_, err := c.Session.SimpleCmd("CONFIRM", "")
	return err
}

// Message just shows window with only OK button.
func (c Client) Message() {
	c.Session.SimpleCmd("MESSAGE", "")
}
