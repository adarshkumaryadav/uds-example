# uds-example
Unix domain socket

---

```markdown
# 🧩 Unix Domain Socket Communication in Kubernetes using Golang

This project demonstrates how **two containers inside a single Kubernetes Pod** can communicate using a **Unix Domain Socket (UDS)** — a concept heavily used in real-world systems like **Kubernetes CSI drivers**, **Docker**, and **container sidecars**.

---

## 📘 What Are Unix Domain Sockets (UDS)?

### 🧠 Basic Definition

A **Unix Domain Socket (UDS)** is an inter-process communication (IPC) mechanism that allows **two processes on the same machine** to communicate by reading and writing to a **special file** in the filesystem.

### 📁 Key Properties

| Property         | Value                             |
|------------------|------------------------------------|
| Communication    | Local only (same machine)          |
| Address Format   | File path (e.g., `/tmp/app.sock`)  |
| Performance      | Faster than TCP (no network stack) |
| Security         | Controlled via file permissions    |
| Lifetime         | Valid as long as the file exists   |

---

## ⚙️ UDS vs TCP vs hostPath

| Feature              | UDS (`emptyDir`)           | TCP Socket                  | hostPath (for UDS)       |
|----------------------|----------------------------|-----------------------------|---------------------------|
| Same Pod             | ✅ Perfect fit              | ✅                          | ✅                        |
| Cross Pod            | ❌ Not possible             | ✅ Over Service IP          | 🟡 Risky (manual path)    |
| Performance          | ⚡ Fast (low latency)        | ⚠️ Slower (network overhead) | ⚠️ Fast but risky         |
| Security             | ✅ File permission           | ⚠️ Exposed over network      | ❌ Can expose host system |
| Port management      | ❌ Not needed               | ✅ Required                  | ❌ Not needed             |

---

## 🚀 Why Use UDS in Kubernetes?

- **Used extensively in CSI drivers** for communication between:
  - CSI Plugin <--> Sidecars (e.g., node-driver-registrar)
  - Kubelet <--> CSI Plugin
- Ideal for **intra-Pod** communication using `emptyDir` shared volume
- No need to expose ports or manage network rules

---

## 📦 Project Structure

```

uds-example/
├── server/
│   ├── main.go         # Server creates UDS socket and listens
│   └── Dockerfile
├── client/
│   ├── main.go         # Client connects to UDS socket and sends message
│   └── Dockerfile
└── pod.yaml            # Kubernetes Pod with 2 containers & shared emptyDir

````

---

## 🧠 How It Works

1. **Server container** creates a Unix socket file at `/csi/csi.sock`
2. **Client container** connects to that socket and sends a message
3. Both containers share a **common volume (`emptyDir`)** where the socket is created
4. Socket exists only for the **lifetime of the Pod**

---

## 🔨 Build & Deploy

### Step 1: Build and Push Docker Images

```bash
# Build & push server
cd server
docker build -t <your-dockerhub>/uds-server:latest .
docker push <your-dockerhub>/uds-server:latest

# Build & push client
cd ../client
docker build -t <your-dockerhub>/uds-client:latest .
docker push <your-dockerhub>/uds-client:latest
````

### Step 2: Update `pod.yaml` with your DockerHub image names

```yaml
image: <your-dockerhub>/uds-server:latest
image: <your-dockerhub>/uds-client:latest
```

### Step 3: Apply Pod

```bash
kubectl apply -f pod.yaml
```

### Step 4: View Logs

```bash
kubectl logs uds-demo -c uds-server
kubectl logs uds-demo -c uds-client
```

---

## 📥 Expected Output

### Server Logs:

```
🔌 Server listening on /csi/csi.sock
📩 Server received: Hello from client
```

### Client Logs:

```
📥 Client received: Hello from server via UDS
```

---

## 🧰 Why emptyDir?

We use `emptyDir` for the following reasons:

| Reason                   | Explanation                                      |
| ------------------------ | ------------------------------------------------ |
| Temporary socket sharing | UDS only needed while Pod is alive               |
| Pod-local communication  | Containers in the same Pod share the same volume |
| Safer than hostPath      | Doesn’t expose host filesystem                   |
| Cleaned up automatically | Deleted with Pod                                 |

---

## 🧪 Real-World Usage: Kubernetes CSI Driver

In the CSI architecture:

* `node-driver-registrar` (sidecar) communicates with CSI plugin via **UDS**
* Both containers are in the **same Pod**
* Shared volume (usually mounted at `/csi`) contains the socket file (`csi.sock`)
* `kubelet` uses hostPath to talk to plugin at `/var/lib/kubelet/plugins/.../csi.sock`

This project mimics this exact behavior in a minimal setup.

---

## 💬 Top Interview Questions (with Answers)

### 1. **What is a Unix Domain Socket?**

> A file-based IPC mechanism used for fast and secure communication between processes on the same machine.

### 2. **Why use UDS over TCP in Kubernetes?**

> UDS is faster, doesn't need ports, and is more secure for intra-Pod communication.

### 3. **Can you use UDS across Pods?**

> No, unless you expose the socket using `hostPath`, which is risky and not recommended.

### 4. **What’s the role of `emptyDir` in UDS communication?**

> It provides a shared ephemeral volume for containers in the same Pod to read/write the Unix socket.

### 5. **Where is UDS used in Kubernetes?**

> Between sidecar containers and the CSI Plugin in storage architecture; also between Docker CLI and daemon via `/var/run/docker.sock`.

### 6. **Why not use hostPath for sharing UDS?**

> HostPath gives access to the node filesystem, which is unsafe and reduces portability.

---

## 📘 References

* [Go net package](https://pkg.go.dev/net)
* [Kubernetes CSI Docs](https://kubernetes-csi.github.io/docs/)
* [Kubernetes Volumes - emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir)

---

## 👨‍💻 Author

Built with ❤️ by **Adarsh Kumar Yadav**
→ Follow for more system-level Golang & Kubernetes learning!
📸 Instagram: [@computer\_science\_engineers](https://instagram.com/computer_science_engineers)

---

```

---