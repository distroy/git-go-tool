/*
 * Copyright (C) distroy
 */

package resultobj

import (
	"fmt"

	"github.com/distroy/git-go-tool/core/jsoncore"
)

func (p *Result) Unmarshal(b []byte) error {
	if err := jsoncore.Unmarshal(b, p); err != nil {
		return err
	}

	dataRaw, err := jsoncore.Marshal(p.Data)
	if err != nil {
		return err
	}

	switch p.Type {
	default:
		return fmt.Errorf("invalid type: %s", p.Type)

	case TypeGoCognitive:
		p.Data = &GoComplexityData{}
	case TypeGoCoverage:
		p.Data = &GoCoverageData{}
	case TypeGoFormat:
		p.Data = &GoFormatData{}
	}

	if err := jsoncore.Unmarshal(dataRaw, p.Data); err != nil {
		return err
	}

	return nil
}
