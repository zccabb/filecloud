package bolt

import (
	"filecloud/settings"
	"fmt"
	"github.com/asdine/storm/v3"
	bolt "go.etcd.io/bbolt"
)

type settingsBackend struct {
	db *storm.DB
}

func view_kv(s settingsBackend) {
	err := s.db.Bolt.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		fmt.Println("ALL")
		_ = tx.ForEach(func(k []byte, b *bolt.Bucket) error {
			fmt.Printf("key=%s\n", k)
			return nil
		})
		fmt.Println("USER")
		_ = tx.Bucket([]byte("User")).ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		fmt.Println("Config")
		_ = tx.Bucket([]byte("config")).ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		fmt.Println("Link")
		_ = tx.Bucket([]byte("Link")).ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		fmt.Println("storm_db")
		_ = tx.Bucket([]byte("__storm_db")).ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
	fmt.Println(err)
}
func (s settingsBackend) Get() (*settings.Settings, error) {
	set := &settings.Settings{}
	//view_kv(s)
	return set, get(s.db, "settings", set)
}

func (s settingsBackend) Save(set *settings.Settings) error {
	return save(s.db, "settings", set)
}

func (s settingsBackend) GetServer() (*settings.Server, error) {
	server := &settings.Server{}
	return server, get(s.db, "server", server)
}

func (s settingsBackend) SaveServer(server *settings.Server) error {
	return save(s.db, "server", server)
}
