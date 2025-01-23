// Copyright (c) 2024 Nvidia Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package quirks

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	maxLockTries    = 200
	HotPlugDelaySec = 12
)

var HotplugFifoName = "/tmp/.nvidiaHotPlugFIFO"
var HotPlugQuirksDisabled = false
var HotPlugDelay = HotPlugDelaySec * time.Second
var fifoMutex = &sync.Mutex{}

func ExecHotPlugQuirks(logger *logrus.Entry) error {
	if HotPlugQuirksDisabled {
		return nil
	}
	// avoid concurrency within
	fifoMutex.Lock()
	defer fifoMutex.Unlock()

	f, err := os.OpenFile(HotplugFifoName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fifoLock := &unix.Flock_t{
		Type:   unix.F_WRLCK,
		Whence: unix.SEEK_SET,
		Start:  0,
		Len:    1024,
	}
	for ix := 0; ix < maxLockTries; ix++ {
		err = unix.FcntlFlock(f.Fd(), unix.F_SETLK, fifoLock)
		if err != nil {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if err != nil {
		// insert random delay and return
		logger.Infof("Failed to acquire lock, using random delay err:%v", err)
		time.Sleep(getRandomDelay())
		return nil
	}

	fifoUnlock := &unix.Flock_t{
		Type:   unix.F_UNLCK,
		Whence: unix.SEEK_SET,
		Start:  0,
		Len:    1024,
	}
	delay := getFifoDelay(logger)
	// unlock before sleeping to let other processes read the fifo
	unix.FcntlFlock(f.Fd(), unix.F_SETLK, fifoUnlock)
	logger.Infof("QUIRK: Sleeping for:%d milliseconds",
		delay.Milliseconds())
	time.Sleep(delay)
	return nil
}

func getFifoDelay(logger *logrus.Entry) time.Duration {
	var err error
	var fifoTS time.Time
	var delay time.Duration

	delay = time.Millisecond // default to a non-zero value
	newTS := time.Now().Add(HotPlugDelay)
	fts := &fifoTS

	hpTS, err := os.ReadFile(HotplugFifoName)
	if err != nil {
		logger.Errorf("Failed to read %s %v", HotplugFifoName, err)
	} else {
		err = fts.UnmarshalBinary(hpTS)
		if err != nil {
			logger.Errorf("Failed to unmarshal ts %v", err)
		} else {
			now := time.Now()
			if fifoTS.After(now) {
				delay = fifoTS.Sub(now)
				newTS = fifoTS.Add(HotPlugDelay)
			}
		}
	}
	newTSData, err := newTS.MarshalBinary()
	if err != nil {
		logger.Errorf("Failed to marshal ts %v", err)
	} else {
		err = os.WriteFile(HotplugFifoName, newTSData, 0644)
		if err != nil {
			logger.Errorf("Failed to write %s %v", HotplugFifoName, err)
		}
	}

	return delay
}

func getRandomDelay() time.Duration {
	n := HotPlugDelaySec + rand.Intn(10*HotPlugDelaySec)
	return time.Duration(n) * time.Second
}
