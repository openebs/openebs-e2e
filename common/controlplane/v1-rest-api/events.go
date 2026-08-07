package v1_rest_api

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
)

func (cp CPv1RestApi) GetEvents(flags ...string) ([]common.EventRecord, error) {
	panic(fmt.Errorf("not implemented REST api"))
}

func (cp CPv1RestApi) GetEventsWithStderr(flags ...string) ([]common.EventRecord, string, error) {
	panic(fmt.Errorf("not implemented REST api"))
}

func (cp CPv1RestApi) GetRawEventsOutput(format string, flags ...string) ([]byte, error) {
	panic(fmt.Errorf("not implemented REST api"))
}

func (cp CPv1RestApi) CountEvents(flags ...string) int {
	panic(fmt.Errorf("not implemented REST api"))
}
