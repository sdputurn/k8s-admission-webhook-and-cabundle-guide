package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func init() {
	_ = corev1.AddToScheme(scheme)
	_ = admissionv1.AddToScheme(scheme)
}

func main() {
	http.HandleFunc("/validate", handleValidation)
	server := &http.Server{
		Addr: ":443",
	}
	fmt.Println("Starting webhook server...")
	err := server.ListenAndServeTLS("certs/tls.crt", "certs/tls.key")
	if err != nil {
		fmt.Printf("Failed to start webhook server: %v\n", err)
	}
}

func handleValidation(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	if err := json.NewDecoder(r.Body).Decode(&admissionReview); err != nil {
		http.Error(w, fmt.Sprintf("could not decode request: %v", err), http.StatusBadRequest)
		return
	}

	req := admissionReview.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		http.Error(w, fmt.Sprintf("could not unmarshal pod object: %v", err), http.StatusBadRequest)
		return
	}

	// Log the Pod name and namespace
	log.Printf("Validating Pod: %s in Namespace: %s", pod.Name, pod.Namespace)

	allowed := validateAnnotations(&pod)
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: allowed,
	}

	if allowed {
		admissionResponse.Result = &metav1.Status{
			Message: "webhook annotation present. allowed",
		}
		log.Printf("Pod %s in Namespace %s: allowed", pod.Name, pod.Namespace)
	} else {
		admissionResponse.Result = &metav1.Status{
			Message: "webhook annotation missing. denied",
			Code:    http.StatusForbidden,
		}
		log.Printf("Pod %s in Namespace %s: denied (missing annotation)", pod.Name, pod.Namespace)
	}

	response := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: admissionResponse,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

func validateAnnotations(pod *corev1.Pod) bool {
	annotations := pod.ObjectMeta.Annotations
	if _, exists := annotations["webhook.dbd.in/v1"]; exists {
		return true
	}
	return false
}
