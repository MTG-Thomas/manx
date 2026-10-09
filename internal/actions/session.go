// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MTG-Thomas/manx/internal/lifecycle"
	"golang.org/x/sys/unix"
)

// All views call this action. No session state influences GateRunner's safety gate.
func runSession(args []string) error {
	operation := "status"
	if len(args) > 0 {
		operation = args[0]
		args = args[1:]
	}
	flags := flag.NewFlagSet("session", flag.ContinueOnError)
	endpoint := flags.String("endpoint", os.Getenv("MANX_LIFECYCLE_ENDPOINT"), "public HTTPS enrollment origin")
	directory := flags.String("state-dir", filepath.Join(toolkitOut(), "lifecycle"), "private journal directory; operator-owned persistent scratch for reboot persistence")
	diagnostic := flags.String("diagnostic", "hardware_inventory", "typed diagnostic category")
	if err := flags.Parse(args); err != nil {
		return errors.New("invalid session flags")
	}
	if err := os.MkdirAll(*directory, 0700); err != nil {
		return errors.New("session directory unavailable")
	}
	lock, err := os.OpenFile(filepath.Join(*directory, ".lock"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return errors.New("session lock unavailable")
	}
	defer lock.Close()
	if err := unix.Flock(int(lock.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		return errors.New("session client already running")
	}
	defer unix.Flock(int(lock.Fd()), unix.LOCK_UN)
	client, err := lifecycle.Open(*directory, *endpoint)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	switch operation {
	case "register":
		err = client.Register(ctx)
	case "sync":
		err = client.Sync(ctx)
	case "record":
		err = client.Record(*diagnostic)
	case "status":
	default:
		return errors.New("session operation must be register, status, sync or record")
	}
	if err != nil {
		return err
	}
	value, err := json.Marshal(client.Status())
	if err != nil {
		return err
	}
	fmt.Println(string(value))
	return nil
}
