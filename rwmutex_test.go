package detectlock

import (
	"sync"
	"testing"
)

func BenchmarkRWMutex_RLock(b *testing.B) {
	count := 0
	lockers := []*RWMutex{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	b.Run("EnableDebug", func(b *testing.B) {
		EnableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.RLock()
					defer locker.RUnlock()
				}
				count++
			}()
		}
	})

	b.Run("DisableDebug", func(b *testing.B) {
		DisableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.RLock()
					defer locker.RUnlock()
				}
				count++
			}()
		}
	})

	b.Run("sync.RWMutex", func(b *testing.B) {
		lockers := []*sync.RWMutex{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.RLock()
					defer locker.RUnlock()
				}
				count++
			}()
		}
	})
}

func BenchmarkRWMutex_TryRLock(b *testing.B) {
	count := 0
	lockers := []*RWMutex{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	b.Run("EnableDebug", func(b *testing.B) {
		EnableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryRLock()
					defer locker.RUnlock()
				}
				count++
			}()
		}
	})

	b.Run("DisableDebug", func(b *testing.B) {
		DisableDebug()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryRLock()
					defer locker.RUnlock()
				}
				count++
			}()
		}
	})

	b.Run("sync.RWMutex", func(b *testing.B) {
		lockers := []*sync.RWMutex{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryRLock()
					defer locker.RUnlock()
				}
				count++
			}()
		}
	})
}

func BenchmarkRWMutex_Lock(b *testing.B) {
	count := 0
	lockers := []*RWMutex{
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
				count++
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
				count++
			}()
		}
	})

	b.Run("sync.RWMutex", func(b *testing.B) {
		lockers := []*sync.RWMutex{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.Lock()
					defer locker.Unlock()
				}
				count++
			}()
		}
	})
}

func BenchmarkRWMutex_TryLock(b *testing.B) {
	count := 0
	lockers := []*RWMutex{
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
				count++
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
				count++
			}()
		}
	})

	b.Run("sync.RWMutex", func(b *testing.B) {
		lockers := []*sync.RWMutex{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func() {
				for _, locker := range lockers {
					locker.TryLock()
					defer locker.Unlock()
				}
				count++
			}()
		}
	})
}
