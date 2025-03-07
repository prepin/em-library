package services

import (
	"context"
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"fmt"
	"time"

	"resty.dev/v3"
)

type RESTSongInfoService struct {
	logger config.Logger
	config config.ServicesConfig
}

func NewRESTSongInfoService(cfg config.ServicesConfig, logger config.Logger) *RESTSongInfoService {
	return &RESTSongInfoService{
		config: cfg,
		logger: logger,
	}
}

type SongDetailResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (s *RESTSongInfoService) GetInfo(ctx context.Context, band, song string) (*entities.SongDetail, error) {
	c := resty.New()
	defer c.Close()

	responseData := SongDetailResponse{}
	var resp *resty.Response
	var err error

	resultCh := make(chan struct {
		resp *resty.Response
		err  error
	}, 1)

	go func() {
		r, e := c.R().
			SetQueryParam("group", band).
			SetQueryParam("song", song).
			SetHeader("Accept", "application/json").
			SetResult(&responseData).
			Get(s.config.InfoServiceURL + "/info")

		resultCh <- struct {
			resp *resty.Response
			err  error
		}{resp: r, err: e}
	}()

	select {
	case result := <-resultCh:
		resp = result.resp
		err = result.err
	case <-time.After(time.Duration(s.config.Timeout) * time.Millisecond):
		return nil, errs.ErrServiceProblem{Err: fmt.Errorf("SongDetailService timeout")}
	}

	if err != nil {
		return nil, errs.ErrServiceProblem{Err: err}
	}

	// если 400-ка или 500-ка
	if resp.IsError() {
		return nil, errs.ErrServiceProblem{Err: fmt.Errorf("SongDetailService fail HTTP status:%d Detail:%+v", resp.StatusCode(), resp)}
	}

	t, err := time.Parse("02.01.2006", responseData.ReleaseDate)
	if err != nil {
		return nil, errs.ErrServiceProblem{Err: fmt.Errorf("failed parsing SongDetailService response %w", err)}
	}

	songDetail := entities.SongDetail{
		ReleaseDate: t,
		Link:        responseData.Link,
		Lyrics:      responseData.Text,
	}

	return &songDetail, nil
}
