package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// GetCursorPosition return the real position of the cursor.
func GetCursorPosition(g *gocui.Gui, v *gocui.View) (int, int, error) {
	return getCursorPosition(g, v)
}

func getCursorPosition(g *gocui.Gui, v *gocui.View) (int, int, error) {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()

	// px and py are absolute postions of the view.
	px, py := ox+cx, oy+cy
	if px < 0 || py < 0 {
		return 0, 0, fmt.Errorf("invalid point")
	}

	return px, py, nil
}

// MoveCursorUp move the cursor of view up.
func MoveCursorUp(g *gocui.Gui, v *gocui.View, dY int) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()

	_, py, _ := getCursorPosition(g, v)
	if py == dY {
		log.WithField("pkg", "view").Debugf("block the cursor, the line '%d'", py)
		if !(oy > 0) {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			return err
		}

		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}

		return nil
	}

	log.WithField("pkg", "view").Debug("move the cursor up")
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

// MoveCursorDown move the cursor of view down.
func MoveCursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	nextLine, _ := getNextViewLine(v)
	if nextLine == "" {
		log.WithField("pkg", "view").Debugf("block the cursor down, line '%d'", cy)
		return nil
	}

	log.WithField("pkg", "view").Debug("move the cursor down")
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

// MoveCursorTop move the cursor top.
func MoveCursorTop(g *gocui.Gui, v *gocui.View, dY int) error {
	_, oy := v.Origin()
	cx, cy := v.Cursor()

	for {
		if cy == 0 {
			break
		}

		if oy+cy == dY {
			break
		}

		cy--
	}

	log.WithField("pkg", "view").Debugf("move the cursor top, line '%d'", cy)
	return v.SetCursor(cx, cy)
}

// MoveCursorBottom move the cursor bottom.
func MoveCursorBottom(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	cx, cy := v.Cursor()

	for {
		if cy == maxY-1 {
			break
		}

		if nextLine, err := v.Line(cy+1); nextLine == "" || err != nil {
			break
		}

		cy++
	}

	log.WithField("pkg", "view").Debugf("move the cursor bottom, line '%d'", cy)
	return v.SetCursor(cx, cy)
}

func getNextViewLine(v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 1); err != nil {
		l = ""
	}

	return l, err
}
