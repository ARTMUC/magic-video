package file

import (
	"path"

	"github.com/ARTMUC/magic-video/internal/domain/composition"
)

func ImageKey(videoComposition *composition.VideoComposition, image *composition.Image) string {
	return path.Join(videoComposition.ID.String(), image.ID.String())
}
