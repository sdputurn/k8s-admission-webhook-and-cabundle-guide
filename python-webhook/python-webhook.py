from flask import Flask, request, jsonify
import base64
import json

app = Flask(__name__)

@app.route('/validate', methods=['POST'])
def validate():
    request_uid = ""  # Initialize request_uid to avoid UnboundLocalError
    
    try:
        admission_review = request.get_json()
        request_uid = admission_review['request']['uid']
        
        # Safely access the 'raw' field
        print(admission_review, type(admission_review))
        raw_data = admission_review['request']['object']
        if not raw_data:
            raise ValueError("'raw' field is missing in the request object")
        
        # pod = json.loads(base64.b64decode(raw_data).decode('utf-8'))
        pod=raw_data
        
        pod_name = pod['metadata']['name']
        pod_namespace = pod['metadata']['namespace']
        
        # Logging the pod name and namespace
        print(f"Validating Pod: {pod_name} in Namespace: {pod_namespace}")
        
        annotations = pod['metadata'].get('annotations', {})
        allowed = "webhook.dbd.in/v1" in annotations
        status_message = "webhook annotation present. allowed" if allowed else "webhook annotation missing. denied"

        response = {
            "apiVersion": "admission.k8s.io/v1",
            "kind": "AdmissionReview",
            "response": {
                "uid": request_uid,
                "allowed": allowed,
                "status": {
                    "message": status_message
                }
            }
        }
        
        print(f"Pod {pod_name} in Namespace {pod_namespace}: {'allowed' if allowed else 'denied'}")
        
        return jsonify(response)
    except Exception as e:
        return jsonify({
            "apiVersion": "admission.k8s.io/v1",
            "kind": "AdmissionReview",
            "response": {
                "uid": request_uid,
                "allowed": False,
                "status": {
                    "message": f"Error processing request: {str(e)}",
                    "code": 500
                }
            }
        })


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=443, ssl_context=('certs/tls.crt', 'certs/tls.key'))
