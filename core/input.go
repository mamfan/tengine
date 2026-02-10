package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/term"
)

func StartInput() (<-chan KeyboardEvent, func(), error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return nil, nil, fmt.Errorf("stdin is not a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, nil, fmt.Errorf("enable raw mode: %w", err)
	}

	events := make(chan KeyboardEvent, 32)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			b, err := reader.ReadByte()
			if err != nil {
				close(events)
				return
			}
			if ev, ok := parseKeyByte(b); ok {
				events <- ev
			}
		}
	}()

	restore := func() {
		_ = term.Restore(fd, oldState)
	}

	return events, restore, nil
}

func parseKeyByte(b byte) (KeyboardEvent, bool) {
	if b == 3 {
		return KeyboardEvent{
			Key:       int('C'),
			Action:    KeyDown,
			Modifiers: ModCtrl,
		}, true
	}

	if b >= 0x20 && b <= 0x7E {
		r := unicode.ToUpper(rune(b))
		return KeyboardEvent{
			Key:       int(r),
			Action:    KeyDown,
			Modifiers: 0,
		}, true
	}

	return KeyboardEvent{}, false
}
