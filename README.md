# myLocal Signup Service 🪧

**Signup Service** is the micro service that powers user sign‑ups for the **myLocal** network.
It combines a **Next.js** UI (App Router) with a **Go** API, packaged and shipped as two containers but runnable together via Docker Compose or Kubernetes.

* **`/ui`** – Next · React Server Components · Tailwind CSS · Turbopack  
* **`/api`** – Go · Fiber · GORM
* **Docker‑first** – one `docker-compose.yml` spins up Postgres + Redis + API + UI for local dev  
* **Kubernetes‑ready** – manifests under `k8s/`; GitHub Actions builds, tags, patches, and deploys

---

## 🗂  Repository layout

```text
signup/
├── api/                 # Go source with Fiber and GORM
├── ui/                  # Next.js App Router
├── docker-compose.yml   # dev stack
└── k8s/                 # prod manifests
```

---

## 🚀  Run locally with Docker Compose

### 1 Clone

```bash
git clone https://github.com/mylo-ing/signup
cd signup
```

### 2 Launch stack

```bash
docker compose up --build
```

| Service  | URL / Port              | Notes                                   |
|----------|-------------------------|-----------------------------------------|
| UI       | <http://localhost:3000> | Fast refresh; proxies `/api`            |
| API      | <http://localhost:3517> | Hot‑reload via **Air**                  |
| Postgres | localhost:5432          | `postgres / password_for_dev_only`      |
| Redis    | localhost:6379          | single instance                         |

> Environment variables live in **`docker-compose.yml`**; tweak as needed.

### 3 Tear down

```bash
docker compose down -v    # -v removes the Postgres volume
```

---

## 🤖  Handy Make targets

```bash
make test      # Go unit tests
make lint      # revive + golangci-lint
make ui        # yarn dev outside Docker
```

---

## 📦  Production CI/CD (DigitalOcean Kubernetes)

1. GitHub Actions builds & pushes images  
   ```
   docker push jtjemsland/signup-api:<sha>
   docker push jtjemsland/signup-ui:<sha>
   ```
2. Workflow patches `k8s/api-deployment.yaml` & `k8s/ui-deployment.yaml` with the same `<sha>`.
3. `kubectl apply -f k8s/` triggers rolling upgrades in DOKS.

---

## 🔧  Tech stack

| Layer        | Library / Tool                         |
|--------------|----------------------------------------|
| UI           | Next, Tailwind CSS                     |
| API          | Go, Fiber, GORM                        |
| Database     | Postgres + PostGIS                     |
| Cache        | Redis                                  |
| Auth         | JWT (guest & user secrets)             |
| Email        | AWS SES                                |
| Dev tooling  | Turbopack, Air hot‑reload, Docker      |

---

## License

This project is licensed under the **GNU AFFERO GENERAL PUBLIC LICENSE v3.0 (AGPL‑3.0)**. See the [`LICENSE`](https://www.gnu.org/licenses/agpl-3.0.html) file for full details.

