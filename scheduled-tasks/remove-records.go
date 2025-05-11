package scheduledtasks

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/go-co-op/gocron/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/v2/bson"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func RemoveRecords() {
	scheduler, schedulerError := gocron.NewScheduler()
	if schedulerError == nil {
		scheduler.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(gocron.NewAtTime(1, 0, 0)), // 01:00 AM
			),
			gocron.NewTask(func() {
				timestamp := gohelpers.MakeTimestampSeconds() - 60*60*24*14 // 2 weeks

				cursor, cursorError := database.MetricsCollection.Find(
					context.Background(),
					bson.M{
						"lastDownloaded": bson.M{"$lt": timestamp},
						"lastViewed":     bson.M{"$lt": timestamp},
					},
				)
				if cursorError != nil {
					log.Fatal(cursorError)
				}
				var records []database.Metrics
				if cursorError = cursor.All(context.Background(), &records); cursorError != nil {
					log.Fatal(cursorError)
				}

				uploadsDirectoryName := utilities.GetEnv(
					constants.ENV_NAMES.UplaodsDirectoryName,
					constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
				)
				uids := make([]string, len(records))
				for index, metrics := range records {
					os.Remove(filepath.Join(uploadsDirectoryName, metrics.UID))
					uids[index] = metrics.UID
				}
				_, queryError := database.FilesCollection.DeleteMany(
					context.Background(),
					bson.M{"uid": bson.M{"$in": uids}},
				)
				if queryError != nil {
					log.Fatal(queryError)
				}
				_, queryError = database.MetricsCollection.DeleteMany(
					context.Background(),
					bson.M{"uid": bson.M{"$in": uids}},
				)
				if queryError != nil {
					log.Fatal(queryError)
				}

				cursor.Close(context.Background())
			}),
		)
		scheduler.Start()
	}
}
