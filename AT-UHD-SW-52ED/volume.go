package atuhdsw52ed

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/connpool"
	"go.uber.org/zap"
)

var (
	ErrorVolumeOR = errors.New("Volume is out of range")
	ErrorMute     = errors.New("Mute State is invalid")
	regGetVolume  = regexp.MustCompile("VOUT1 (...)")
	regGetMute    = regexp.MustCompile("VOUTMute1 (...)")
)

//Volumes .
func (vs *AtlonaVideoSwitcher5x1) Volumes(ctx context.Context) (map[string]int, error) {

	vs.log.Info("Getting Volume")
	volumeOut := make(map[string]int)
	var volumeStr string

	err := vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		cmd := []byte("VOUT1 sta\r\n")
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

			match = regGetVolume.FindAllStringSubmatch(string(buf), -1)

		}
		volumeStr = strings.TrimPrefix(match[0][1], "0")
		if volumeStr[len(volumeStr)-1:] == string('\r') {
			volumeStr = strings.TrimSuffix(volumeStr, string('\r'))
		}
		return nil
	})
	if err != nil {
		return volumeOut, err
	}

	volumeLevel, err := strconv.Atoi(volumeStr)
	if err != nil {
		return volumeOut, fmt.Errorf("failed to convert volume to int: %s", err.Error())
	}
	if volumeLevel < -35 {
		volumeOut[""] = 0
	} else {
		volume := ((volumeLevel + 35) * 2)
		if volume%2 != 0 {
			volume = volume + 1
		}
		volumeOut[""] = volume
	}

	vs.log.Info("Got Volume", zap.Any("Vol", volumeOut))
	return volumeOut, nil
}

//SetVolume .
func (vs *AtlonaVideoSwitcher5x1) SetVolume(ctx context.Context, level int) error {

	if level == 0 {
		level = -80
	} else {
		convertedVolume := -35 + math.Round(float64(level/2))
		level = int(convertedVolume)
	}

	levelStr := strconv.Itoa(level)

	vs.log.Info("Setting Volume", zap.String("to", levelStr))
	cmd := []byte(fmt.Sprintf("VOUT1 %s\r\n", levelStr))

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
			return fmt.Errorf("Unable to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("uanble to write to cvonnection: wrote %v%v bytes", n, len(cmd))
		}

		buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)
		if err != nil {
			return fmt.Errorf("failed to read from connection %w", err)
		}

		if strings.Contains(string(buf), "ERROR") {
			return ErrorVolumeOR
		}

		vs.log.Info("Successfully set volume", zap.String("to", levelStr))

		return nil
	})

}

//Mutes
func (vs *AtlonaVideoSwitcher5x1) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	vs.log.Info("Getting Muted State")
	mute := make(map[string]bool)

	var muteState string

	err := vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		cmd := []byte("VOUTMute1 sta\r\n")

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

			match = regGetMute.FindAllStringSubmatch(string(buf), -1)
		}
		muteState = strings.TrimPrefix(match[0][1], "0")
		if muteState[len(muteState)-1:] == string('\r') {
			muteState = strings.TrimSuffix(muteState, string('\r'))
		}

		return nil
	})

	if err != nil {
		return mute, err
	}

	if muteState == "on" {
		mute[""] = true
	} else {
		mute[""] = false
	}

	vs.log.Info("Got Mute State", zap.Any("State", mute))
	return mute, nil

}

//SetMute .
func (vs *AtlonaVideoSwitcher5x1) SetMute(ctx context.Context, output string, muted bool) error {

	var state string
	if muted {
		state = "on"
	} else {
		state = "off"
	}
	vs.log.Info("Muting System", zap.String("state", state))

	cmd := []byte(fmt.Sprintf("VOUTMute1 %s\r\n", state))

	return vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("uanble to write to connection deadline: %w", err)
		}

		n, err := conn.Write(cmd)

		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write to connection: wrote %v%v bytes", n, len(cmd))

		}

		buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)
		if err != nil {
			return fmt.Errorf("failed to read form connection %w", err)
		}

		if strings.Contains(string(buf), "failed") {
			return ErrorMute
		}

		vs.log.Info("Successfully set MUTE state", zap.String("state", state))
		return nil
	})
}
