// This command isn't currently intended to be usable, it's more of a testbed for trying things out, but may become user-friendly in future.

package reclassifycmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	porndb "git.sr.ht/~dragnel/go-tpdb"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/adult/tpdb"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/gosuri/uiprogress"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Params struct {
	fx.In
	Search search.Search
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

var count int

func New(p Params) (Result, error) {
	cmd := &cli.Command{
		Name: "reclassify",

		Action: func(ctx *cli.Context) error {

			p := porndb.NewClient("HZfErPdglUruXnazXpAPWCPbzmhMaJrDtmsqwVLhddb52e89")
			psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", 5432, "bitmagnet_admin", "bitmagnet", "bitmagnet")
			db, _ := gorm.Open(postgres.Open(psqlconn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})

			count = 0
			var torrents []model.TorrentContent

			db.Where("content_id is null and lower(title) LIKE '%xxx%'").Find(&torrents)

			bar := uiprogress.AddBar(len(torrents)).AppendCompleted().PrependElapsed()
			bar.PrependFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf("Found %d", count)
			})

			fmt.Printf("Found %d torrents to classify\n", len(torrents))
			uiprogress.Start()
			bar.Incr()
			i := 0
			for _, torrent := range torrents {

				name := clean_name(strings.ToLower(torrent.Title))
				scene, err := p.Parse(name)
				if err == nil {
					contentScene, _ := tpdb.SceneToXxxModel(scene)

					torrent.Content = contentScene
					torrent.Title = contentScene.Title
					torrent.ContentType.Valid = true
					torrent.SearchString = contentScene.SearchString

					db.Save(&torrent)
					count += 1
				}
				bar.Incr()
				i++
				if i%5000 == 0 {
					time.Sleep(30 * time.Second)
				}
				// fmt.Println(name)
			}

			return nil
		},
	}
	return Result{Command: cmd}, nil
}

func clean_name(s string) string {
	trash := []string{"galaxxxy", "prt", "rarbg", "xxx", "gush", "vr18", "vsex", "xox", "n1c", "narcos", "gapfill", "bty", "ipt", "rbg", "rq", "1k", "4k",
		"mp4", "x264", "x265", "hevc", "mpeg", "wmv", "dvdrip", "dvd", "xvid", "rip",
		"ktr", "kleenex", "xvx", "iak", "sd", "hd", "xc", "p2p", "french", "wrb", "hr", "oro", ",spankhash", "gagball", "sexors", "-xleech",
		"web", "imageset", "fugli", "ghostfreakxx", "eks265", "sexytv", "yapg", "tbp", "gagball", "hushhush", "team", "-oly", "tg",
		"www.viciosaszt.com", "www.torrenting.org", "www.torrenting.com", "torrenting.com", "www.torrentday.com", "www.xbay.me", "www.iptv.memorial",
		"720p", "1080p", "480p", "1920p", "2160p", "4096", "3840p",
		"[", "]", ",", "(", ")",
	}

	for _, t := range trash {
		s = strings.Replace(s, t, "", -1)
	}
	s = strings.Replace(s, ".", " ", -1)
	s = strings.Replace(s, "-", " ", -1)
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	return s
}
