package atuhdsw52ed

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/connpool"
	"go.uber.org/zap"
)

var (
	ErrorOutOfRange = errors.New("input or output is out of range")
	regGetInput     = regexp.MustCompile("x(.)AVx1")
)

//AudioVideoInputs .
func (vs *AtlonaVideoSwitcher5x1) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	vs.log.Info("Getting the current inputs")
	inputs := make(map[string]string)

	err := vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetReadDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		cmd := []byte("Status\r\n")

		n, err := conn.Write(cmd)

		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		var match [][]string

		for len(match) == 0 {
			buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)

			if err != nil {
				return fmt.Errorf("unable to read from connection: %w", err)
			}
			match = regGetInput.FindAllStringSubmatch(string(buf), -1)
		}

		inputs[""] = strings.TrimPrefix(match[0][1], "0")
		return nil

	})
	if err != nil {
		return inputs, err
	}
	vs.log.Info("Got inputs", zap.Any("inputs", inputs))
	return inputs, nil

}

//SetAudioVideoInput
func (vs *AtlonaVideoSwitcher5x1) SetAudioVideoInput(ctx context.Context, output, input string) error {
	output = "1"
	vs.log.Info("Setting audio video input", zap.String("output", output), zap.String("input", input))
	cmd := []byte(fmt.Sprintf("x%sAVx%s\r\n", input, output))

	intInput, nerr := strconv.Atoi(input)
	if nerr != nil {
		return fmt.Errorf("error occured when converting input to int: %w", nerr)
	}

	if intInput == 0 || intInput > 5 {
		return fmt.Errorf("Invalid Input. The input requested must be between 1-5. The input you requested was %v", intInput)
	}

	return vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}
		n, err := conn.Write(cmd)

		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write to connection: wrote %v%v bytes", n, len(cmd))
		}

		buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %w", err)
		}
		if strings.Contains(string(buf), "failed") {
			return ErrorOutOfRange
		}

		vs.log.Info("Successfully set audio video input", zap.String("output", output), zap.String("input", input))
		return nil
	})

}
