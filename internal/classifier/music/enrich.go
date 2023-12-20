package music

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func PreEnrich(input model.TorrentContent) (model.TorrentContent, error) {
	if hasAudio := input.Torrent.HasFileType(model.FileTypeAudio); hasAudio.Valid && !hasAudio.Bool {
		return model.TorrentContent{}, classifier.ErrNoMatch
	}

	if input.ContentType.Valid && !input.ContentType.ContentType.IsAudio() {
		return model.TorrentContent{}, classifier.ErrNoMatch
	}

	return input, nil
}
