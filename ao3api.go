package ao3api

import (
	"fmt"

	"github.com/capoverflow/ao3api/internals/fanfic"
	"github.com/capoverflow/ao3api/models"
)

func Fanfic(params models.FanficParams) (works models.Fanfic, err error) {
	fmt.Println(params)

	works, err = fanfic.GetInfo(params)
	if err != nil {
		return models.Fanfic{}, nil
	}

	return works, nil
}
