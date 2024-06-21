package v1_rest_api

import "fmt"

func (cp CPv1RestApi) DrainNode(nodeName string, drainLabel string, drainTimeOut int) error {
	// #TODO implement REST api for node drain
	panic(fmt.Errorf("not implemented REST api for node drain"))
}

func (cp CPv1RestApi) GetDrainNodeLabels(nodeName string) ([]string, []string, error) {
	// #TODO implement REST api for node drain labels
	panic(fmt.Errorf("not implemented REST api for getting node drain labels"))
}
