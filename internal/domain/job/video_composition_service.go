package composition

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/domain/order"
	"github.com/ARTMUC/magic-video/internal/pkg/file"
)

type VideoCompositionService interface {
	Create() error
	Enqueue(order *order.Order) error
}

type videoCompositionService struct {
	orderLineRepository           order.OrderLineRepository
	videoCompositionJobRepository VideoCompositionJobRepository
	transactionProvider           base.TransactionProvider
	fileManager                   file.ReaderWriter
}

func NewVideoCompositionService(
	orderLineRepository order.OrderLineRepository,
	videoCompositionJobRepository VideoCompositionJobRepository,
	transactionProvider base.TransactionProvider,
	fileManager file.ReaderWriter,
) VideoCompositionService {
	return &videoCompositionService{
		orderLineRepository:           orderLineRepository,
		videoCompositionJobRepository: videoCompositionJobRepository,
		transactionProvider:           transactionProvider,
		fileManager:                   fileManager,
	}
}

func (v *videoCompositionService) Enqueue(order *order.Order) error {
	orderLines, err := v.orderLineRepository.FindMany(base.ReadOptions{
		Scopes:  []base.Scope{OrderLineScopes{}.WithOrderID(order.ID)},
		Preload: []string{OrderLinePreloadVideoComposition},
	})
	if err != nil {
		return fmt.Errorf("failed to find order lines for orderID %d in db: %w", order.ID, err)
	}

	err = v.transactionProvider.Transaction(func(tx *base.Tx) error {
		for _, orderLine := range orderLines {
			videocompositionJob := &VideoCompositionJob{
				BaseModel:          base.BaseModel{},
				VideoCompositionID: orderLine.VideoCompositionID,
				VideoComposition:   orderLine.VideoComposition,
				OrderLineID:        orderLine.ID,
				OrderLine:          &orderLine,
				Status:             VideoCompositionJobStatusCreated,
			}
			err = v.videoCompositionJobRepository.Create(base.WriteOptions{Tx: tx}, videocompositionJob)
			if err != nil {
				return fmt.Errorf("failed to create video composition job in db: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (v *videoCompositionService) Create() error {
	jobs, err := v.videoCompositionJobRepository.FindMany(base.ReadOptions{
		Scopes: []base.Scope{
			VideoCompositionJobScopes{}.WithStatus(VideoCompositionJobStatusCreated)},
		Preload: []string{
			VideoCompositionJobPreloadOrderLine,
			VideoCompositionJobPreloadVideoComposition,
			VideoCompositionJobPreloadVideoCompositionImages,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get unprocessed jobs: %w", err)
	}

	for _, job := range jobs {
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
		if err := v.videoCompositionJobRepository.Update(base.WriteOptions{}, job); err != nil {
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
