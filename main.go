package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/thoas/go-funk"
	networkingv1 "k8s.io/api/networking/v1"
)

const (
	TCP  string = "TCP"
	UDP  string = "UDP"
	SCTP string = "SCTP"
)

func hasEgressPolicy(netPol *networkingv1.NetworkPolicy) bool {
	for _, val := range netPol.Spec.PolicyTypes {
		if val == networkingv1.PolicyTypeEgress {
			return true
		}
	}
	return false
}

func hasIngressPolicy(netPol *networkingv1.NetworkPolicy) bool {
	for _, val := range netPol.Spec.PolicyTypes {
		if val == networkingv1.PolicyTypeIngress {
			return true
		}
	}
	return false
}

func NetPolPodSelector(netPol *networkingv1.NetworkPolicy) map[string]string {
	return netPol.Spec.PodSelector.MatchLabels
}

func IngressNsSelectors(ingressRules []networkingv1.NetworkPolicyIngressRule) map[string]string {
	ingressNsSelectors := map[string]string{}
	for _, ingressRule := range ingressRules {
		for _, networkPolicyPeer := range ingressRule.From {
			if networkPolicyPeer.NamespaceSelector != nil {
				for key, val := range networkPolicyPeer.PodSelector.MatchLabels {
					ingressNsSelectors[key] = val
				}
			}
		}
	}
	return ingressNsSelectors
}

func IngressPodSelectors(ingressRules []networkingv1.NetworkPolicyIngressRule) map[string]string {
	ingressPodSelectors := map[string]string{}
	for _, ingressRule := range ingressRules {
		for _, networkPolicyPeer := range ingressRule.From {
			if networkPolicyPeer.PodSelector != nil {
				for key, val := range networkPolicyPeer.PodSelector.MatchLabels {
					ingressPodSelectors[key] = val
				}
			}
		}
	}
	return ingressPodSelectors
}

func IngressPortInfo(ingressRules []networkingv1.NetworkPolicyIngressRule) map[string][]string {
	// map : key : protocol, value : slices of portNums
	ingressPortMap := map[string][]string{}
	// store all ingress Port information
	for _, ingressRule := range ingressRules {
		for _, port := range ingressRule.Ports {
			protocol := fmt.Sprintf("%s", *port.Protocol)
			portNum := port.Port.String()
			ingressPortMap[protocol] = append(ingressPortMap[protocol], portNum)
		}
	}
	return ingressPortMap
}

func EgressPodSelectors(egressRules []networkingv1.NetworkPolicyEgressRule) map[string]string {
	egressPodSelectors := map[string]string{}
	for _, egressRule := range egressRules {
		for _, networkPolicyPeer := range egressRule.To {
			if networkPolicyPeer.PodSelector != nil {
				for key, val := range networkPolicyPeer.PodSelector.MatchLabels {
					egressPodSelectors[key] = val
				}
			}
		}
	}
	return egressPodSelectors
}

func EgressNsSelectors(egressRules []networkingv1.NetworkPolicyEgressRule) map[string]string {
	egressNsSelectors := map[string]string{}
	for _, egressRule := range egressRules {
		for _, networkPolicyPeer := range egressRule.To {
			if networkPolicyPeer.NamespaceSelector != nil {
				for key, val := range networkPolicyPeer.NamespaceSelector.MatchLabels {
					egressNsSelectors[key] = val
				}
			}
		}
	}
	return egressNsSelectors
}

// Store all information
func EgressPortInfo(egressRules []networkingv1.NetworkPolicyEgressRule) map[string][]string {
	// map : key : protocol, value : slices of portNums
	egressPortMap := map[string][]string{}
	// store all ingress Port information
	for _, egressRule := range egressRules {
		for _, port := range egressRule.Ports {
			protocol := fmt.Sprintf("%s", *port.Protocol)
			portNum := port.Port.String()
			egressPortMap[protocol] = append(egressPortMap[protocol], portNum)
		}
	}
	return egressPortMap
}

func createNetPolObjFromYaml(file string) *networkingv1.NetworkPolicy {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}

	var netPol networkingv1.NetworkPolicy
	err = yaml.Unmarshal(bytes, &netPol)
	if err != nil {
		panic(err.Error())
	}

	//	fmt.Println(netPol)
	return &netPol
}

func isSubset(protocol string, ingressPorts, egressPorts map[string][]string) bool {
	return funk.Subset(ingressPorts[protocol], egressPorts[protocol]) // true
}

// check all protocols
func isSubsets(ingressPorts, egressPorts map[string][]string) bool {
	for protocol, _ := range ingressPorts {
		isSubset := funk.Subset(ingressPorts[protocol], egressPorts[protocol])
		if !isSubset {
			return false
		}
	}

	return true
}

func HasCommonPodSelectors(ingressPodSelectors, PodSelectorInNetPolB map[string]string) (bool, string, string) {
	for matchKey, matchVal := range ingressPodSelectors {
		val, ok := PodSelectorInNetPolB[matchKey]
		if ok {
			if matchVal == val {
				return true, matchKey, matchVal
			}
		}
	}
	return false, "", ""
}

func IsBreak(netPolAYaml, netPolBYaml string) bool {
	netPolA := createNetPolObjFromYaml(netPolAYaml)
	netPolB := createNetPolObjFromYaml(netPolBYaml)

	// Check to see whether common podSelector between PodSelector from "From" in netPolA and PodSelector from "NetworkPolicy" in netPolB
	ingressPodSelectors := IngressPodSelectors(netPolA.Spec.Ingress)
	PodSelectorInNetPolB := NetPolPodSelector(netPolB)
	hasCommonPodSelectors, key, label := HasCommonPodSelectors(ingressPodSelectors, PodSelectorInNetPolB)
	if !hasCommonPodSelectors {
		fmt.Println("No Common #1")
		return false
	}

	// Two netPols have common labels in both sides, we need to check ports information
	fmt.Println("key : ", key, " val : ", label)

	// Check to see whether common podSelector between PodSelector from "To" in netPolB and PodSelector from "NetworkPolicy" in netPolA
	egressPodSelectors := EgressPodSelectors(netPolB.Spec.Egress)
	PodSelectorInNetPolA := NetPolPodSelector(netPolA)
	hasCommonPodSelectors, key, label = HasCommonPodSelectors(egressPodSelectors, PodSelectorInNetPolA)
	if !hasCommonPodSelectors {
		fmt.Println("No Common #2")
		return false
	}

	// Two netPols have common labels in both sides, we need to check ports information
	fmt.Println("key : ", key, " val : ", label)
	// check whether podSelector's labels from "From" field in "ingress" exist in egress PodSelector in "networkPolicy"
	ingressMap := IngressPortInfo(netPolA.Spec.Ingress)
	fmt.Println("Ingress port Map ", ingressMap)

	egressMap := EgressPortInfo(netPolB.Spec.Egress)
	fmt.Println("Egress port Map ", egressMap)

	// ingressMap from netPolA - map[pod:a], map[TCP:[80]]
	// egressMap from netPolB  - map[pod:a], map[TCP:[80 53] UDP:[53]]

	// (TODO): Need to check more protocol
	isSubset := isSubsets(ingressMap, egressMap)

	fmt.Println("Is Subset ", isSubset)
	if isSubset {
		fmt.Println("Ok Egress has all port information of ingress")
		return false
	}

	return true
}
func main() {
	netPolAYaml := "./service-ingress.yaml"

	// false
	netPolBYaml := "./ok-nginx-egress.yaml"
	fmt.Println(netPolBYaml)
	fmt.Println("# ", netPolBYaml, " Final result: ", IsBreak(netPolAYaml, netPolBYaml))

	// false
	netPolBYaml = "./ok-nginx-egress-1.yaml"
	fmt.Println(netPolBYaml)
	fmt.Println("# ", netPolBYaml, " Final result: ", IsBreak(netPolAYaml, netPolBYaml))

	// false
	netPolBYaml = "./ok-nginx-egress-2.yaml"
	fmt.Println(netPolBYaml)
	fmt.Println("# ", netPolBYaml, " Final result: ", IsBreak(netPolAYaml, netPolBYaml))

	// true
	netPolBYaml = "./break-nginx-egress.yaml"
	fmt.Println(netPolBYaml)
	fmt.Println("# ", netPolBYaml, " Final result: ", IsBreak(netPolAYaml, netPolBYaml))

	// // problem case since netPolA does not have 53 port
	// fmt.Println(isSubsets(TCP, egressMap, ingressMap))
}
