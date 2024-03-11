package cards

import (
	"github.com/Clinet/clinet_storage"
	"github.com/JoshuaDoes/logger"
)

var (
	Log     *logger.Logger
	Storage *storage.Storage
)

func Init() error {
	Log = logger.NewLogger("cards", 2)
	Storage = &storage.Storage{}
	if err := Storage.LoadFrom("cards"); err != nil {
		Log.Error(err)
		return err
	}
	return nil
}