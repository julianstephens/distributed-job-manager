package worker

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/logger"
)

const (
	DefaultSandboxTTL                    = 30 * time.Minute
	DefaultInactivityExpirationThreshold = 20 * time.Minute
	DefaultCleanupFrequency              = 5 * time.Minute
)

type Sandbox struct {
	ID         int
	ExpiresAt  time.Time
	LastUsedAt time.Time
}

type SandboxPool struct {
	mu                            sync.Mutex
	Count                         int
	Reserved                      map[string]*Sandbox
	AvailableBoxes                map[int]bool
	SandboxTTL                    time.Duration
	CleanupFrequency              time.Duration
	InactivityExpirationThreshold time.Duration
}

var ErrSandboxBusy error = errors.New("sandbox busy")

func NewSandboxPool(count int) *SandboxPool {
	s := make(map[int]bool)
	for i := range count {
		cmd := exec.Command("isolate", "--init", fmt.Sprintf("-b %v", i))
		op, err := cmd.Output()
		if err != nil {
			logger.Errorf("unable to create sandbox %v, output: %v, err: %v", i, string(op), err)
		} else {
			s[i] = true
		}
	}

	return &SandboxPool{
		Count:                         count,
		AvailableBoxes:                s,
		Reserved:                      make(map[string]*Sandbox),
		SandboxTTL:                    DefaultSandboxTTL,
		CleanupFrequency:              DefaultCleanupFrequency,
		InactivityExpirationThreshold: DefaultInactivityExpirationThreshold,
	}
}

func (s *SandboxPool) AvailableCount() int {
	return len(s.AvailableBoxes)
}

func (s *SandboxPool) GetSandbox(userID string) (*Sandbox, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.Reserved[userID]
	return b, ok
}

func (s *SandboxPool) UpdateLastUsed(userID string) {
	box := s.Reserved[userID]
	if box != nil {
		box.LastUsedAt = time.Now()
	}
}

func (s *SandboxPool) Reserve(userID string) (*Sandbox, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.AvailableBoxes) == 0 {
		return nil, ErrSandboxBusy
	}

	var selected int
	for k := range s.AvailableBoxes {
		selected = k
		break
	}

	box := &Sandbox{
		ID:         selected,
		ExpiresAt:  time.Now().Add(s.SandboxTTL),
		LastUsedAt: time.Now(),
	}

	s.Reserved[userID] = box
	delete(s.AvailableBoxes, selected)

	return box, nil
}

func (s *SandboxPool) Release(userID string) {
	s.mu.Lock()
	s.AvailableBoxes[s.Reserved[userID].ID] = true
	delete(s.Reserved, userID)
	s.mu.Unlock()
}

func (s *SandboxPool) init(boxID int) error {
	return exec.Command("isolate", "--init", fmt.Sprintf("-b %v", boxID)).Run()
}

func (s *SandboxPool) delete(boxID int) error {
	return exec.Command("isolate", "--cleanup", fmt.Sprintf("-b %v", boxID)).Run()
}

func (s *SandboxPool) ScheduleCleanup() {
	ticker := time.NewTicker(s.CleanupFrequency)
	go func() {
		for {
			<-ticker.C
			s.Cleanup()
		}
	}()
}
func (s *SandboxPool) Cleanup() {
	logger.Infof("cleaning sandboxes...")
	s.mu.Lock()

	releaseQueue := []string{}

	for userID, box := range s.Reserved {
		if box == nil {
			continue
		}

		expired := box.ExpiresAt.Before(time.Now())
		inactive := box.LastUsedAt.Add(s.InactivityExpirationThreshold).Before(time.Now())

		if expired || inactive {
			err := s.delete(box.ID)
			if err != nil {
				logger.Errorf("failed to cleanup sandbox %v %v", box.ID, err)
				continue
			}

			err = s.init(box.ID)
			if err != nil {
				logger.Errorf("failed to init sandbox %v %v", box.ID, err)
				continue
			}

			releaseQueue = append(releaseQueue, userID)
		}
	}
	s.mu.Unlock()

	for _, user := range releaseQueue {
		s.Release(user)
	}

	logger.Infof("%d sandboxes cleaned", len(releaseQueue))
}
