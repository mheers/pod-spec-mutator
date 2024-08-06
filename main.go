package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// Configuration variables
	podNameRegex   = os.Getenv("POD_NAME_REGEX")
	namespaceRegex = os.Getenv("NAMESPACE_REGEX")
	patchJSON      = os.Getenv("PATCH_JSON")

	currentNamespace = os.Getenv("POD_NAMESPACE")
	currentPodName   = os.Getenv("POD_NAME")
)

func mutate(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// Verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "invalid Content-Type, expected `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *admissionv1.AdmissionResponse
	ar := admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		admissionResponse = createAdmissionResponse(&ar)
	}

	admissionReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
	}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the request and response
	fmt.Println("Request: ", string(body))
	fmt.Println("Response: ", string(resp))
}

func createAdmissionResponse(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// Check if this pod is the controller itself
	if req.Namespace == currentNamespace && pod.Name == currentPodName {
		return &admissionv1.AdmissionResponse{
			Allowed: true,
		}
	}

	// Check if the pod name matches the regex
	if podNameRegex != "" {
		match, _ := regexp.MatchString(podNameRegex, pod.Name)
		if !match {
			return &admissionv1.AdmissionResponse{
				Allowed: true,
			}
		}
	}

	// Apply the patch
	patchedPod, err := applyPatch(&pod, []byte(patchJSON))
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// Create patch
	patchBytes, err := createPatch(pod, *patchedPod)
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch // JSONPatch from RFC 6902
			return &pt
		}(),
	}
}

func applyPatch(pod *corev1.Pod, patchJSONSpec []byte) (*corev1.Pod, error) {
	patchedPod := pod.DeepCopy()

	originalSpec, err := json.Marshal(patchedPod.Spec)
	if err != nil {
		return nil, fmt.Errorf("error marshaling original pod spec: %v", err)
	}

	// Perform a strategic merge patch
	mergedSpec, err := strategicpatch.StrategicMergePatch(originalSpec, patchJSONSpec, corev1.PodSpec{})
	if err != nil {
		return nil, fmt.Errorf("error merging patch: %v", err)
	}

	var newSpec corev1.PodSpec
	if err := json.Unmarshal(mergedSpec, &newSpec); err != nil {
		return nil, fmt.Errorf("error unmarshaling merged spec: %v", err)
	}

	patchedPod.Spec = newSpec

	return patchedPod, nil
}

func createPatch(originalPod, modifiedPod corev1.Pod) ([]byte, error) {
	patches, err := comparePods(originalPod, modifiedPod)
	if err != nil {
		return nil, fmt.Errorf("error comparing pod specs: %v", err)
	}

	return json.Marshal(patches)
}

func printInfo() {
	fmt.Println("POD_NAME_REGEX: ", podNameRegex)
	fmt.Println("NAMESPACE_REGEX: ", namespaceRegex)
	fmt.Println("PATCH_JSON: ", patchJSON)
	fmt.Println("Current Namespace: ", currentNamespace)
	fmt.Println("Current Pod Name: ", currentPodName)
}

func main() {
	printInfo()
	http.HandleFunc("/mutate", mutate)
	fmt.Println("Starting webhook server on port 8443")
	err := http.ListenAndServeTLS(":8443", "/tmp/k8s-webhook-server/serving-certs/tls.crt", "/tmp/k8s-webhook-server/serving-certs/tls.key", nil)
	if err != nil {
		panic(err)
	}
}
