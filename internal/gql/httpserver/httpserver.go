package httpserver

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	qbt "github.com/autobrr/go-qbittorrent"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Schema graphql.ExecutableSchema
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	gql := handler.NewDefaultServer(p.Schema)
	pg := playground.Handler("GraphQL playground", "/graphql")
	return Result{
		Option: &builder{
			gqlHandler: func(c *gin.Context) {
				gql.ServeHTTP(c.Writer, c.Request)
			},
			playgroundHandler: func(c *gin.Context) {
				pg.ServeHTTP(c.Writer, c.Request)
			},
		},
	}
}

type builder struct {
	gqlHandler        gin.HandlerFunc
	playgroundHandler gin.HandlerFunc
}

func (builder) Key() string {
	return "graphql"
}

func (b builder) Apply(e *gin.Engine, cfg httpserver.Config) error {
	e.POST("/graphql", b.gqlHandler)
	e.GET("/graphql", b.playgroundHandler)
	e.GET("/add/torrent/:info_hash/:category", func(c *gin.Context) {
		hash := c.Params.ByName("info_hash")
		category := c.Params.ByName("category")

		c_qbt := qbt.NewClient(qbt.Config{
			Host:     cfg.HostClientTorrent,
			Username: cfg.UserClientTorrent,
			Password: cfg.PasswordClientTorrent,
		})

		options := map[string]string{
			"category": category,
		}
		c_qbt.AddTorrentFromUrl("magnet:?xt=urn:btih:"+hash, options)
	})
	return nil
}
