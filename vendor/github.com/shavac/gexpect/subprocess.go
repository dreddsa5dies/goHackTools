package gexpect

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/shavac/gexpect/pty"
)

var (
	err error
)

type SubProcess struct {
	Term            *pty.Terminal
	cmd             *exec.Cmd
	DelayBeforeSend time.Duration
	CheckInterval   time.Duration
	Before          []byte
	After           []byte
	Match           []byte
	echo            bool
}

func (sp *SubProcess) Start() (err error) {
	return sp.Term.Start(sp.cmd)
}

func (sp *SubProcess) Close() (err error) {
	sp.Terminate()
	return sp.Term.Close()
}

func (sp *SubProcess) WaitTimeout(d time.Duration) (err error) {
	if sp.echo {
		go func() {
			io.Copy(os.Stdout, sp)
		}()
	}
	execerr := make(chan error, 1)
	if d > 0 {
		go func() {
			time.Sleep(d)
			execerr <- TIMEOUT
		}()
	}
	go func() {
		execerr <- sp.cmd.Wait()
	}()
	return <-execerr
}

func (sp *SubProcess) Wait() error {
	return sp.WaitTimeout(0)
}

func (sp *SubProcess) Terminate() error {
	return sp.cmd.Process.Kill()
}

func (sp *SubProcess) Expect(expreg ...*regexp.Regexp) (matchIndex int, err error) {
	return sp.ExpectTimeout(0, expreg...)
}

func (sp *SubProcess) ExpectTimeout(timeout time.Duration, expreg ...*regexp.Regexp) (matchIndex int, err error) {
	buf := make([]byte, 2048)
	c := make(chan byte, 1)
	checkpoint := make(chan int, 1)
	rerr := make(chan error, 1)
	go func() {
		for {
			if _, err := io.ReadAtLeast(sp, buf, 1); err != nil {
				rerr <- err
				continue
			}
			for _, b := range buf {
				c <- b
			}
			if sp.echo {
				os.Stdout.Write(buf)
				os.Stdout.Sync()
			}
		}
	}()
	if timeout > 0 {
		go func() {
			time.Sleep(timeout)
			rerr <- TIMEOUT
		}()
	}
	go func() {
		for i := 1; ; i++ {
			time.Sleep(sp.CheckInterval)
			checkpoint <- i
		}
	}()
	for {
		select {
		case c1 := <-c:
			buf = append(buf, c1)
		case e := <-rerr:
			sp.Before = append(sp.Before, buf...)
			return -1, e
		case <-checkpoint:
			for idx, re := range expreg {
				if loc := re.FindIndex(buf); loc != nil {
					sp.Match = buf[loc[0]:loc[1]]
					sp.Before = append(sp.Before, buf[0:loc[0]]...)
					buf = make([]byte, 2048)
					return idx, nil
				}
			} // no match
		}
	}
	return -1, nil
}

func (sp *SubProcess) Read(b []byte) (n int, err error) {
	return sp.Term.Read(b)
}

func (sp *SubProcess) Write(b []byte) (n int, err error) {
	time.Sleep(sp.DelayBeforeSend)
	return sp.Term.Write(b)
}

func (sp *SubProcess) Writeln(b []byte) (n int, err error) {
	bn := append(b, []byte("\r\n")...)
	return sp.Write(bn)
}

func (sp *SubProcess) Send(response string) (err error) {
	_, err = sp.Write([]byte(response))
	return
}

func (sp *SubProcess) SendLine(response string) (err error) {
	return sp.Send(response + "\r\n")
}

func (sp *SubProcess) Interact() (err error) {
	return sp.InteractTimeout(0)
}

func (sp *SubProcess) InteractTimeout(d time.Duration) (err error) {
	sp.Write(sp.After)
	sp.After = []byte{}
	oldState, _ := pty.Tcgetattr(os.Stdin)
	pty.SetRaw(os.Stdin)
	defer pty.Tcsetattr(os.Stdin, oldState)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGWINCH, syscall.SIGTSTP)
	go func() {
		for sig := range s {
			switch sig {
			case syscall.SIGINT:
				sp.Term.SendIntr()
			case syscall.SIGWINCH:
				if x, y, err := pty.GetWinSize(os.Stdout); err == nil {
					sp.Term.SetWinSize(x, y)
				}
			default:
				continue
			}
		}
	}()
	execerr := make(chan error, 1)
	if d > 0 {
		go func() {
			time.Sleep(d)
			execerr <- TIMEOUT
		}()
	}
	go func() {
		execerr <- sp.cmd.Wait()
	}()
	in := make(chan byte, 1)
	stdin := bufio.NewReader(os.Stdin)
	go func() error {
		var b byte
		for {
			if b, err = stdin.ReadByte(); err != nil {
				if err == io.EOF {
					sp.Term.SendEOF()
					continue
				} else {
					return err
				}
			}
			in <- b
		}
	}()
	go func() {
		io.Copy(os.Stdout, sp)
		return
	}()
	for {

		select {
		case err := <-execerr:
			return err
		case b := <-in:
			_, err = sp.Write([]byte{b})
		}
	}
	return
}

func (sp *SubProcess) Echo() {
	sp.echo = true
}

func (sp *SubProcess) NoEcho() {
	sp.echo = false
}

func NewSubProcess(name string, arg ...string) (sp *SubProcess, err error) {
	sp = new(SubProcess)
	if sp.Term, err = pty.NewTerminal(); err != nil {
		return
	}
	sp.cmd = exec.Command(name, arg...)
	sp.DelayBeforeSend = 50 * time.Microsecond
	sp.CheckInterval = time.Microsecond
	if x, y, err := pty.GetWinSize(os.Stdout); err == nil {
		sp.Term.SetWinSize(x, y)

	}
	return
}
