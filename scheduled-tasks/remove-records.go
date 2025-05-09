package scheduledtasks

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/v2/bson"

	"file-sharing/database"
)

func RemoveRecords() {
	scheduler, schedulerError := gocron.NewScheduler()
	if schedulerError == nil {
		fmt.Println("scheduler")
		scheduler.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(gocron.NewAtTime(14, 20, 0)),
			),
			gocron.NewTask(func() {
				fmt.Println("tick")
				seconds := gohelpers.MakeTimestampSeconds() - 60
				res, queryError := database.FilesCollection.DeleteMany(
					context.Background(),
					bson.M{"createdAt": bson.M{"$lt": seconds}},
				)
				fmt.Println("res", res, "err", queryError)
				res, queryError = database.MetricsCollection.DeleteMany(
					context.Background(),
					bson.M{"createdAt": bson.M{"$lt": seconds}},
				)
				fmt.Println("res", res, "err", queryError)
			}),
		)
		scheduler.Start()
	}
}
