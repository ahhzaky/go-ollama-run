# Ollama CLI Chat Client 🦙💬

CLI client to interact with Ollama AI models (Llama3.2-vision) with conversation logging feature.
![Demo](https://img.shields.io/badge/Demo-Coming_Soon-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/username/ollama-cli-chat)](https://goreportcard.com/report/github.com/username/ollama-cli-chat)

## Fitur Utama ✨

- 💬 Interactive chat with AI model Llama3.2-vision
- 📝 Automatic conversation log to file `conversation.txt`
- ⚡ GPU & CPU performance optimization
- 🎨 Loading indicator animation
- ⏱️ Response time measurement
- 🔧 Flexible model parameter configuration

## Prasyarat 🔧

- [Ollama](https://ollama.ai/) terinstall and running in `localhost:11434`
- Model Llama3.2-vision downloaded:
  ```bash
  ollama pull llama3.2-vision
  ```
- Go 1.20+

## Installation 🚀

```bash
git clone https://github.com/username/ollama-cli-chat.git
cd ollama-cli-chat
go build -o ollama-chat
```

## Use

```bash
./ollama-chat

Chat with Ollama AI (type 'exit' to quit)

Enter your message: [Ketik pesan Anda]
```

## Configuration ⚙️

The default parameters can be modified in the source code:

```go
req := Request{
    Model:       "llama3.2-vision",
    NumGPU:      50,   // Layer GPU
    NumThreads:  8,    // Thread CPU
    Temperature: 0.7,  // Creativity (0.0-1.0)
    TopP:        0.9,  // Filtering token
}
```

## Conversation Log 📂

Conversation history is stored in conversation.txt with the format:

```bash
User: [Pesan pengguna]
Ollama: [Respon AI]
```

## License 📄

MIT License © 2025
