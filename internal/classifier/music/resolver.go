package music

import (
	"context"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/music/discogs"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	DiscogsClient discogs.Client
}

type Result struct {
	fx.Out
	Resolver classifier.SubResolver `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: musicResolver{
			config:        classifier.SubResolverConfig{Key: "music", Priority: 3},
			discogsClient: p.DiscogsClient,
		},
	}
}

type musicResolver struct {
	config        classifier.SubResolverConfig
	discogsClient discogs.Client
}

func (r musicResolver) Config() classifier.SubResolverConfig {
	return r.config
}

func (r musicResolver) PreEnrich(content model.TorrentContent) (model.TorrentContent, error) {
	return PreEnrich(content)
}

func (r musicResolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	titleLower := strings.ToLower(content.Torrent.Name)

	if strings.Contains(titleLower, "discography") ||
		strings.Contains(titleLower, "discographie") ||
		strings.Contains(titleLower, "discografia") ||
		strings.Contains(titleLower, "anthology") {

		artist, err := FindArtistDiscography(titleLower)
		if err != nil {
			return content, err
		}

		resultArtist, err := r.discogsClient.SearchArtist(ctx, discogs.SearchMusicParams{
			Artist:               artist,
			Album:                "",
			Track:                "",
			LevenshteinThreshold: 0,
		})

		content.ContentType.Valid = true
		content.ContentType.ContentType = model.ContentTypeMusic
		content.Content = resultArtist
		return content, err
	} else {
		for _, ext := range model.FileTypeAudio.Extensions() {
			if strings.Contains(titleLower, ext) {
				content.ContentType.Valid = true
				content.ContentType.ContentType = model.ContentTypeMusic
				return content, nil
			}
		}
	}
	return content, nil
}
