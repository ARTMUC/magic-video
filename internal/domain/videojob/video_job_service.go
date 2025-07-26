package job

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/ARTMUC/magic-video/internal/contracts"
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/domain/order"
	"github.com/ARTMUC/magic-video/internal/pkg/file"
	"github.com/google/uuid"
)

type VideoCompositionService interface {
	Create() error
	Enqueue(orderID uuid.UUID) error
}

type VideoCompositionsGetter interface {
	GetByID(id uuid.UUID) (contracts.VideoComposition, bool, error)
}

type videoCompositionService struct {
	orderLineRepository     order.OrderLineRepository
	videoJobRepository      VideoCompositionJobRepository
	transactionProvider     base.TransactionProvider
	fileManager             file.ReaderWriter
	videoCompositionsGetter VideoCompositionsGetter
}

func NewVideoCompositionService(
	orderLineRepository order.OrderLineRepository,
	videoCompositionJobRepository VideoCompositionJobRepository,
	transactionProvider base.TransactionProvider,
	fileManager file.ReaderWriter,
	videoCompositionsGetter VideoCompositionsGetter,
) VideoCompositionService {
	return &videoCompositionService{
		orderLineRepository:     orderLineRepository,
		videoJobRepository:      videoCompositionJobRepository,
		transactionProvider:     transactionProvider,
		fileManager:             fileManager,
		videoCompositionsGetter: videoCompositionsGetter,
	}
}

func (v *videoCompositionService) Enqueue(orderID uuid.UUID) error {
	videocompositionJob := &VideoCompositionJob{
		BaseModel: base.BaseModel{},
		OrderID:   orderID,
		Status:    VideoCompositionJobStatusCreated,
	}
	err := v.videoJobRepository.Create(base.WriteOptions{}, videocompositionJob)
	if err != nil {
		return fmt.Errorf("failed to create video composition job in db: %w", err)
	}

	return nil
}

func (v *videoCompositionService) Create() error {
	jobs, err := v.videoJobRepository.FindMany(base.ReadOptions{
		Scopes: []base.Scope{
			VideoCompositionJobScopes{}.WithStatus(VideoCompositionJobStatusCreated)},
		Preload: []string{},
	})
	if err != nil {
		return fmt.Errorf("failed to get unprocessed jobs: %w", err)
	}

	for _, job := range jobs {
		videoCompositions, ok, err := v.videoCompositionsGetter.GetByID(orderLine.videoCompositionID)
		if !ok {
			return fmt.Errorf("no composition found by id: %s")
		}
		if err != nil {
			job.Error = sql.Null[string]{
				V:     err.Error(),
				Valid: true,
			}
			job.Status = VideoCompositionJobStatusFailed
			updateErr := v.videoJobRepository.Update(base.WriteOptions{}, &job)
			if updateErr != nil {
				return fmt.Errorf("failed to update job status: %w", updateErr)
			}

			return fmt.Errorf("failed to get video compositions: %w", err)
		}

		switch job.VideoComposition.VideoTemplate {
		case "SuperHero1":
			// in future we put code here
		default:
			return fmt.Errorf("unsupported video template: %s, VideoCompositionJob: %d", job.VideoComposition.VideoTemplate, job.ID)

		}

		compositionImages := make([]image.Image, len(job.VideoComposition.Images))
		for _, img := range job.VideoComposition.Images {
			f, err := v.fileManager.Read(file.ImageKey(job.VideoComposition, &img))
			if err != nil {
				return fmt.Errorf("failed to read image file: %d : %w", img.ID, err)
			}
			compositionImage, err := processImageBytes(f)
			if err != nil {
				return fmt.Errorf("failed to process image bytes: %d : %w", img.ID, err)
			}
			compositionImages[img.ID] = compositionImage
		}

		// Remove background
		photoNoBackground, err := removeBackground(photo)
		if err != nil {
			return fmt.Errorf("failed to remove background: %w", err)
		}

		// Generate AI videos
		aiVideos, err := generateAIVideos(photoNoBackground)
		if err != nil {
			return fmt.Errorf("failed to generate AI videos: %w", err)
		}

		// Compose final video
		finalVideo, err := composeVideo(aiVideos)
		if err != nil {
			return fmt.Errorf("failed to compose video: %w", err)
		}

		// Store video
		videoPath, err := storeVideo(finalVideo)
		if err != nil {
			return fmt.Errorf("failed to store video: %w", err)
		}

		// Send email with download link
		if err := notifyCustomer(job.OrderLine.Order.CustomerEmail, videoPath); err != nil {
			return fmt.Errorf("failed to notify customer: %w", err)
		}

		// Update job status
		job.Status = VideoCompositionJobStatusCompleted
		if err := v.videoJobRepository.Update(base.WriteOptions{}, job); err != nil {
			return fmt.Errorf("failed to update job status: %w", err)
		}
	}

	return nil

}

// @TODO make helper struct
func processImageBytes(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}
