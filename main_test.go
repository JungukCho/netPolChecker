package main

import (
	"reflect"
	"testing"
)

type expectedValue struct {
	want   bool
	labels map[string]string
}

func TestIngressNsSelectors(t *testing.T) {
	netPolAYaml := "./service-ingress.yaml"
	netpolObj := createNetPolObjFromYaml(netPolAYaml)

	ingressNsSelectors := IngressNsSelectors(netpolObj.Spec.Ingress)
	t.Log(ingressNsSelectors)

	want := true
	matchLables := map[string]string{
		"ingress-ns": "ingress-ns-nginx",
	}

	got := reflect.DeepEqual(ingressNsSelectors, matchLables)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestEgressNsSelectors(t *testing.T) {
	netPolAYaml := "./ok-nginx-egress-3.yaml"
	netpolObj := createNetPolObjFromYaml(netPolAYaml)

	egressNsSelectors := EgressNsSelectors(netpolObj.Spec.Egress)
	t.Log(egressNsSelectors)

	want := true
	matchLables := map[string]string{
		"egress-ns": "egress-ns-nginx",
	}

	got := reflect.DeepEqual(egressNsSelectors, matchLables)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}
