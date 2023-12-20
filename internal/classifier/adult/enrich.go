package adult

import (
	"regexp"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func PreEnrich(input model.TorrentContent) (model.TorrentContent, error) {
	if !strings.Contains(strings.ToLower(input.Torrent.Name), "xxx") {
		return model.TorrentContent{}, classifier.ErrNoMatch
	}

	titleLower := strings.ToLower(input.Torrent.Name)
	titleLower = clean_name(titleLower)

	output := input
	output.Title = titleLower

	if !output.VideoResolution.Valid {
		output.VideoResolution = model.InferVideoResolution(output.Title)
	}
	if !output.VideoSource.Valid {
		output.VideoSource = model.InferVideoSource(output.Title)
	}
	if !output.VideoModifier.Valid {
		output.VideoModifier = model.InferVideoModifier(output.Title)
	}
	if !output.Video3d.Valid {
		output.Video3d = model.InferVideo3d(output.Title)
	}
	if !output.VideoCodec.Valid || !output.ReleaseGroup.Valid {
		vc, rg := model.InferVideoCodecAndReleaseGroup(output.Title)
		if !output.VideoCodec.Valid {
			output.VideoCodec = vc
		}
		if !output.ReleaseGroup.Valid {
			output.ReleaseGroup = rg
		}
	}
	return output, nil
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
