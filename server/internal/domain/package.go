package domain

import "cod-server/internal/utils"

type Package struct {
	ID    string   `json:"id"`
	Cards utils.Map[string, *Card] `json:"cards"`
}


