package search

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/acheong08/cybertron/pkg/models/bert"
	"github.com/acheong08/cybertron/pkg/tasks"
	"github.com/acheong08/cybertron/pkg/tasks/textencoding"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type encoder struct {
	textencoding.Interface
	*errgroup.Group
}

func (e *encoder) extractVectors(contentLines []string, embeddingMaxLines int) ([][]float64, error) {
	var texts []string
	firstLine, lastLine := 0, embeddingMaxLines

	for {
		if lastLine >= len(contentLines) {
			lastLine = len(contentLines)
		}

		texts = append(texts, strings.Join(contentLines[firstLine:lastLine], "\n"))

		if lastLine == len(contentLines) {
			break
		}

		firstLine += embeddingMaxLines / 3
		lastLine += embeddingMaxLines / 3
	}

	return e.encodeMulti(texts)
}

func new() *encoder {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	var err error
	home_dir, err := os.UserHomeDir()
	if err != nil {
		home_dir = "."
	}
	// Create ~/.models directory if it doesn't exist
	if _, err := os.Stat(home_dir + "/.models"); os.IsNotExist(err) {
		os.Mkdir(home_dir+"/.models", 0755)
	}

	modelsDir := home_dir + "/.models"
	modelName := "sentence-transformers/all-MiniLM-L6-v2"

	m, err := tasks.Load[textencoding.Interface](&tasks.Config{ModelsDir: modelsDir, ModelName: modelName})
	if err != nil {
		panic(err)
	}

	group := &errgroup.Group{}
	group.SetLimit(10)

	return &encoder{
		Interface: m,
		Group:     group,
	}
}

func (e *encoder) encodeMulti(texts []string) ([][]float64, error) {
	var resultMutex sync.Mutex
	var wg sync.WaitGroup
	results := make([][]float64, len(texts))

	for i, text := range texts {
		i, text := i, text

		e.Group.Go(func() error {
			result, err := e.Interface.Encode(context.Background(), text, int(bert.MeanPooling))
			if err != nil {
				return nil
			}

			resultMutex.Lock()
			defer resultMutex.Unlock()
			results[i] = result.Vector.Data().F64()

			return nil
		})
	}

	wg.Wait()
	return results, nil
}
