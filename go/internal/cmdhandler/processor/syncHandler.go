package processor

import "github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
import syncer "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/sync"

func (c *Worker) doSync(syncer *syncer.DomJudgeSyncer) {
	logging.Info("start cmdhandler sync")
	if err := syncer.Sync(); err != nil {
		logging.Error("sync failed", err)
	} else {
		logging.Info("finished cmdhandler sync")
	}
}
