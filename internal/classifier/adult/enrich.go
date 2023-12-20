package adult

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func PreEnrich(input model.TorrentContent) (model.TorrentContent, error) {
	if !strings.Contains(strings.ToLower(input.Torrent.Name), "xxx") {
		return model.TorrentContent{}, classifier.ErrNoMatch
	}

	titleLower := strings.ToLower(input.Torrent.Name)
	titleLower = strings.Replace(titleLower, "com_", "", -1)
	titleLower = strings.Replace(titleLower, "www.torrenting.com", "", -1)
	titleLower = strings.Replace(titleLower, "www.torrenting.org", "", -1)

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
