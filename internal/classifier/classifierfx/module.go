package classifierfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/adult"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/adult/adultfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/consumer"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/producer"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/publisher"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/music"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/music/musicfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/videofx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"classifier",
		fx.Provide(
			classifier.New,
			consumer.New,
			producer.New,
			publisher.New,
			video.New,
			adult.New,
			music.New,
		),
		videofx.New(),
		adultfx.New(),
		musicfx.New(),
	)
}
