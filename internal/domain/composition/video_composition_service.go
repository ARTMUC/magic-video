//go:build !goverter

package composition

import (
	"errors"
	"fmt"

	"github.com/ARTMUC/magic-video/internal/contracts"
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

type videoCompositionService struct {
	videoCompositionRepository VideoCompositionRepository
}

func (v *videoCompositionService) GetByID(id uuid.UUID) (contracts.VideoComposition, bool, error) {
	videoComposition, err := v.videoCompositionRepository.FindOne(base.ReadOptions{
		[]base.Scope{VideoCompositionScopes{}.WithID(id)}, nil,
	})
	if err != nil {
		if errors.Is(err, base.ErrRecordNotFound) {
			return contracts.VideoComposition{}, false, nil
		}

		return contracts.VideoComposition{}, false, fmt.Errorf("failed to get video compositions: %w", err)
	}

	converter := VideoCompositionConverterImpl{}
	return converter.VideoCompositionDomainToContract(*videoComposition), true, nil
}
