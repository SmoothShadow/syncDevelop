//go:build linux

package clipboard

import (
	"os/exec"
)

func readText() (string, error) {
	// 优先使用xclip，其次xsel
	if _, err := exec.LookPath("xclip"); err == nil {
		out, err := exec.Command("xclip", "-selection", "clipboard", "-o").Output()
		if err == nil {
			return string(out), nil
		}
	}
	if _, err := exec.LookPath("xsel"); err == nil {
		out, err := exec.Command("xsel", "--clipboard", "--output").Output()
		if err == nil {
			return string(out), nil
		}
	}
	return "", nil
}

func writeText(text string) error {
	if _, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go func() {
			defer stdin.Close()
			stdin.Write([]byte(text))
		}()
		return cmd.Run()
	}
	if _, err := exec.LookPath("xsel"); err == nil {
		cmd := exec.Command("xsel", "--clipboard", "--input")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go func() {
			defer stdin.Close()
			stdin.Write([]byte(text))
		}()
		return cmd.Run()
	}
	return nil
}

func readImage() ([]byte, error) {
	if _, err := exec.LookPath("xclip"); err == nil {
		out, err := exec.Command("xclip", "-selection", "clipboard", "-o", "-t", "image/png").Output()
		if err == nil && len(out) > 0 {
			return out, nil
		}
	}
	return nil, nil
}

func writeImage(data []byte) error {
	if _, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go func() {
			defer stdin.Close()
			stdin.Write(data)
		}()
		return cmd.Run()
	}
	return nil
}
