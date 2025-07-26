package file

import (
	"path"

	"github.com/ARTMUC/magic-video/internal/domain/composition"
)

func ImageKey(videoComposition *composition.VideoComposition, image *composition.Image) string {
	return path.Join(videoComposition.UUID.String(), image.UUID.String())
}
