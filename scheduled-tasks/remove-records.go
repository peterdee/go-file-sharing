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
				gocron.NewAtTimes(gocron.NewAtTime(1, 0, 0)),
			),
			gocron.NewTask(func() {
				seconds := gohelpers.MakeTimestampSeconds() - 60*60*24*14 // 2 weeks
				cursor, cursorError := database.FilesCollection.Find(
					context.Background(),
					bson.M{"createdAt": bson.M{"$lt": seconds}},
				)
				if cursorError != nil {
					log.Fatal(cursorError)
				}
				var filesRecords []database.Files
				if cursorError = cursor.All(context.Background(), &filesRecords); cursorError != nil {
					log.Fatal(cursorError)
				}
				uids := make([]string, len(filesRecords))
				uploadsDirectoryName := utilities.GetEnv(
					constants.ENV_NAMES.UplaodsDirectoryName,
					constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
				)
				for _, file := range filesRecords {
					uids = append(uids, file.UID)
					os.Remove(filepath.Join(uploadsDirectoryName, file.UID))
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
