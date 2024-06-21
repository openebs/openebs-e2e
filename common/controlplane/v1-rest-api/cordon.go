package v1_rest_api

import "fmt"

func (cp CPv1RestApi) CordonNode(nodeName string, cordonLabel string) error {
	// #TODO implement REST api for node cordon
	panic(fmt.Errorf("not implemented REST api for node cordon"))
}

func (cp CPv1RestApi) GetCordonNodeLabels(nodeName string) ([]string, error) {
	// #TODO implement REST api for node cordon labels
	panic(fmt.Errorf("not implemented REST api for getting node cordon labels"))
}

func (cp CPv1RestApi) UnCordonNode(nodeName string, cordonLabel string) error {
	// #TODO implement REST api for node uncordon
	panic(fmt.Errorf("not implemented REST api for node uncordon"))
}
