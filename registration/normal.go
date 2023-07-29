package registration

import (
	"runtime"
	"sync"
	"time"

	"github.com/deroproject/derohe/transaction"
	"github.com/deroproject/derohe/walletapi"
)

type NormalReg struct {
	Running bool
	OnFound func(tx *transaction.Transaction)

	wallet    *walletapi.Wallet_Disk
	hashRate  map[int]uint64
	hashCount map[int]uint64
	mutex     sync.RWMutex
}

func NewNormalReg() *NormalReg {
	return &NormalReg{
		hashRate:  make(map[int]uint64),
		hashCount: make(map[int]uint64),
	}
}

func (s *NormalReg) Start(workers int, wallet *walletapi.Wallet_Disk) {
	s.wallet = wallet
	s.Running = true
	for i := 0; i < workers; i++ {
		go s.run(i)
	}
}

func (s *NormalReg) Stop() {
	s.Running = false
}

func (s *NormalReg) HashRate() uint64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	sum := uint64(0)
	for _, v := range s.hashRate {
		sum += v
	}

	return sum
}

func (s *NormalReg) HashCount() uint64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	sum := uint64(0)
	for _, v := range s.hashCount {
		sum += v
	}

	return sum
}

func (s *NormalReg) run(wIndex int) {
	start := time.Now()
	count := uint64(0)
	hashRate := uint64(0)

	for {
		if !s.Running {
			s.mutex.Lock()
			s.hashCount = make(map[int]uint64)
			s.hashRate = make(map[int]uint64)
			s.mutex.Unlock()

			return
		}

		for i := 0; i < runtime.GOMAXPROCS(0); i++ {
			tx := s.wallet.GetRegistrationTX()
			count++
			hashRate++

			if time.Now().Add(-1 * time.Second).After(start) {
				s.mutex.Lock()
				s.hashRate[wIndex] = hashRate
				s.hashCount[wIndex] = count
				s.mutex.Unlock()

				start = time.Now()
				hashRate = 0
			}

			hash := tx.GetHash()
			if hash[0] == 0 && hash[1] == 0 && hash[2] == 0 {
				if tx.IsRegistrationValid() {
					s.Stop()
					s.OnFound(tx)
					break
				}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
