package adult

import (
	"context"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/adult/tpdb"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	TpdbClient tpdb.Client
}

type Result struct {
	fx.Out
	Resolver classifier.SubResolver `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: adultResolver{
			config:     classifier.SubResolverConfig{Key: "adult", Priority: 2},
			tpdbClient: p.TpdbClient,
		},
	}
}

type adultResolver struct {
	config     classifier.SubResolverConfig
	tpdbClient tpdb.Client
}

func (r adultResolver) Config() classifier.SubResolverConfig {
	return r.config
}

func (r adultResolver) PreEnrich(content model.TorrentContent) (model.TorrentContent, error) {
	//return content, nil
	return PreEnrich(content)
}

func (r adultResolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {

	if r.tpdbClient != nil {
		titleLower := strings.ToLower(content.Torrent.Name)
		titleLower = clean_name(titleLower)
		contentAdult, err := r.tpdbClient.SearchScene(ctx, titleLower)
		if err == nil {
			content.Title = contentAdult.Title
			content.ContentType.Valid = true
			content.Content = contentAdult
			content.SearchString = contentAdult.SearchString
			return content, nil
		}
		return model.TorrentContent{}, classifier.ErrNoMatch
	}
	content.ContentType.Valid = true
	content.ContentType.ContentType = model.ContentTypeXxx
	return content, nil

	return model.TorrentContent{}, classifier.ErrNoMatch
}
