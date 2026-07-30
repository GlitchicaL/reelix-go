package transcoder

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func StartCleanupManager(mgr *TranscodeManager, interval, timeout time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Printf("[transcoder] Cleanup manager started (interval=%v, timeout=%v)", interval, timeout)

		for range ticker.C {
			mgr.cleanupCycle(timeout)
		}
	}()
}

func (m *TranscodeManager) cleanupCycle(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiredCount := 0
	zombieCount := 0

	for key, session := range m.sessions {
		isZombie := !isProcessAlive(session.Cmd)
		isExpired := time.Since(session.LastSeen) > timeout

		if isZombie {
			log.Printf("[%s] Zombie detected, cleaning up", session.ID)
			m.cleanupSession(session, key)
			zombieCount++
			continue
		}

		if isExpired && session.State != "expired" {
			log.Printf("[%s] Idle %v, expiring", session.ID, timeout.Round(time.Minute))
			session.State = "expired"
			session.Cmd.Process.Kill()
			os.RemoveAll(session.Directory)
			delete(m.sessions, key)
			expiredCount++
			log.Printf("[%s] Session deleted from cache", session.ID)
		}
	}

	log.Printf("[cleanup] Cycle complete: %d expired, %d zombies", expiredCount, zombieCount)
}

func (m *TranscodeManager) cleanupSession(session *Session, key string) {
	if session.Cmd.Process != nil {
		session.Cmd.Process.Kill()
		log.Printf("[%s] FFmpeg process killed", session.ID)
	}

	os.RemoveAll(session.Directory)
	log.Printf("[%s] Directory removed: %s", session.ID, session.Directory)

	delete(m.sessions, key)
	log.Printf("[%s] Session deleted from cache", session.ID)
}

func isProcessAlive(cmd *exec.Cmd) bool {
	if cmd.Process == nil {
		return false
	}

	err := cmd.Process.Signal(syscall.Signal(0))
	return err == nil
}
