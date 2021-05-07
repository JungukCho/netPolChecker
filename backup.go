// Store all information
func getEgressInfo(netPol *networkingv1.NetworkPolicy) map[string]string {
	egressPortMap := map[string]string{}
	// store all ingress Port information
	for _, egressRule := range netPol.Spec.Egress {
		for _, port := range egressRule.Ports {
			protocol := fmt.Sprintf("%s", *port.Protocol)
			portNum := port.Port.String()
			egressPortMap[protocol] = portNum
		}
	}
	return egressPortMap
}

//
func isConflict(ingressPorts, egressPorts map[string]string) bool {
	for key, val := range ingressPorts {
		netPolBVal, ok := egressPorts[key]

		if !ok {
			return true
		}

		if netPolBVal != val {
			return true
		}
	}
	return false
}

// Store all information
func getIngressInfo(netPol *networkingv1.NetworkPolicy) map[string]string {
	ingressPortMap := map[string]string{}
	// store all ingress Port information
	for _, ingressRule := range netPol.Spec.Ingress {
		for _, port := range ingressRule.Ports {
			protocol := fmt.Sprintf("%s", *port.Protocol)
			portNum := port.Port.String()
			ingressPortMap[protocol] = portNum
		}
	}
	return ingressPortMap
}

// func Diff() bool {
// 	netPolA := map[string]string{
// 		"25":  "TCP",
// 		"443": "TCP",
// 		"53":  "TCP"}

// 	netPolB := map[string]string{
// 		"443": "TCP",
// 		"53":  "TCP"}
// }
