package detectlock

import (
	"sync"
	"testing"
)

func BenchmarkMutex_Lock(b *testing.B) {
	counter := 0
	lockers := []*Mutex{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	b.Run("EnableDebug", func(b *testing.B) {
		EnableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.Lock()
					defer locker.Unlock()
				}
				counter++
			}()
		}
	})

	b.Run("DisableDebug", func(b *testing.B) {
		DisableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.Lock()
					defer locker.Unlock()
				}
				counter++
			}()
		}
	})

	b.Run("sync.Mutex", func(b *testing.B) {
		lockers := []*sync.Mutex{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.Lock()
					defer locker.Unlock()
				}
				counter++
			}()
		}
	})
}

func BenchmarkMutex_TryLock(b *testing.B) {
	counter := 0
	lockers := []*Mutex{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	b.Run("EnableDebug", func(b *testing.B) {
		EnableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryLock()
					defer locker.Unlock()
				}
				counter++
			}()
		}
	})

	b.Run("DisableDebug", func(b *testing.B) {
		DisableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryLock()
					defer locker.Unlock()
				}
				counter++
			}()
		}
	})

	b.Run("sync.Mutex", func(b *testing.B) {
		lockers := []*sync.Mutex{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryLock()
					defer locker.Unlock()
				}
				counter++
			}()
		}
	})
}

func TestMutex_Lock(t *testing.T) {
	EnableDebug()
	t.Run("lock once success", func(t *testing.T) {
		l := &Mutex{}
		l.Lock()
	})

	t.Run("lock twice fail", func(t *testing.T) {
		l := &Mutex{}
		l.Lock()
		if l.TryLock() {
			t.Fail()
		}
	})
}

func TestMutex_TryLock(t *testing.T) {
	EnableDebug()
	t.Run("try-lock once success", func(t *testing.T) {
		l := &Mutex{}
		if !l.TryLock() {
			t.Fail()
		}
	})

	t.Run("try-lock twice fail", func(t *testing.T) {
		l := &Mutex{}
		l.TryLock()
		if l.TryLock() {
			t.Fail()
		}
	})
}

func TestMutex_Unlock(t *testing.T) {
	EnableDebug()
	t.Run("unlock once success", func(t *testing.T) {
		l := &Mutex{}
		counter := 0
		l.Lock()
		counter++
		l.Unlock()
	})
}
