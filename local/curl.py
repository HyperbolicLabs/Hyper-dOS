import requests

url = "http://localhost:8000/v1/chat/completions"
headers = {
    "Content-Type": "application/json"
}
data = {
    "messages": [{
      "role": "user",
      "content": "Tell me about the Mycelial network"
    }],
    # "model": "meta-llama/Llama-3.2-3B-Instruct",
    # "model": "mistralai/Mistral-7B-Instruct-v0.2",
    "model": "TinyLlama/TinyLlama-1.1B-Chat-v1.0",
    "max_tokens": 512,
    # "temperature": 0.7,
    # "top_p": 0.9
}

response = requests.post(url, headers=headers, json=data)
print(response.json())
