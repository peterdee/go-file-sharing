package scheduledtasks

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/go-co-op/gocron/v2"
	"github.com/julyskies/gohelpers"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func MarkAsDeleted() {
	scheduler, schedulerError := gocron.NewScheduler()
	if schedulerError == nil {
		scheduler.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(gocron.NewAtTime(1, 0, 0)), // 01:00 AM
			),
			gocron.NewTask(func() {
				timestamp := gohelpers.MakeTimestampSeconds() - 60*60*24*14 // 2 weeks

				var metricsRecords []database.MetricsModel
				cursorError := database.MetricsService.FindAll(
					context.Background(),
					map[string]any{
						"isDeleted":      false,
						"lastDownloaded": map[string]any{"$lt": timestamp},
						"lastViewed":     map[string]any{"$lt": timestamp},
					},
					&metricsRecords,
				)
				if cursorError != nil {
					log.Fatal(cursorError)
				}

				uploadsDirectoryName := utilities.GetEnv(
					constants.ENV_NAMES.UplaodsDirectoryName,
					constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
				)
				uids := make([]string, len(metricsRecords))
				for index, metrics := range metricsRecords {
					os.Remove(filepath.Join(uploadsDirectoryName, metrics.Uid))
					uids[index] = metrics.Uid
				}
				cache.FileService.DelMany(context.Background(), uids...)

				queryError := database.FileService.UpdateMany(
					context.Background(),
					map[string]any{"uid": map[string]any{"$in": uids}},
					map[string]any{
						"deletedAt": timestamp,
						"isDeleted": true,
						"updatedAt": timestamp,
					},
				)
				if queryError != nil {
					log.Fatal(queryError)
				}
				queryError = database.MetricsService.UpdateMany(
					context.Background(),
					map[string]any{"uid": map[string]any{"$in": uids}},
					map[string]any{
						"deletedAt": timestamp,
						"isDeleted": true,
						"updatedAt": timestamp,
					},
				)
				if queryError != nil {
					log.Fatal(queryError)
				}
			}),
		)
		scheduler.Start()
	}
}
