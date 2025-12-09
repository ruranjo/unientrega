# Docker Setup Guide

This project has two Docker configurations:

## 🔧 Development Mode (Hot Reload)

**Use this when actively developing** - code changes are automatically detected and the app restarts.

### Start Development Environment
```bash
cd docker
docker compose -f docker-compose.dev.yml up --build
```

### Features
- ✅ **Hot reload** - Changes to `.go` files automatically rebuild and restart the app
- ✅ **Volume mounting** - Your local code is mounted into the container
- ✅ **Fast iteration** - No need to rebuild Docker image for code changes
- ✅ **Uses Air** - Automatic rebuild tool for Go

### How It Works
- Your source code is mounted as a volume
- [Air](https://github.com/air-verse/air) watches for file changes
- When you save a `.go` file, Air rebuilds and restarts the app
- Changes appear in ~1-2 seconds

---

## 🚀 Production Mode (Optimized Build)

**Use this for production deployment** - creates an optimized, minimal Docker image.

### Start Production Environment
```bash
cd docker
docker compose up --build -d
```

### Features
- ✅ **Multi-stage build** - Small final image (~20MB)
- ✅ **Optimized binary** - Compiled Go binary, no source code
- ✅ **Secure** - Minimal attack surface
- ✅ **Fast startup** - Pre-compiled binary

### When to Rebuild
You need to rebuild the image when you change code:
```bash
cd docker
docker compose up --build -d
```

---

## 📋 Common Commands

### Development
```bash
# Start dev environment
docker compose -f docker-compose.dev.yml up

# Stop dev environment
docker compose -f docker-compose.dev.yml down

# View logs
docker compose -f docker-compose.dev.yml logs -f api
```

### Production
```bash
# Start production
docker compose up -d

# Stop production
docker compose down

# Rebuild and restart
docker compose up --build -d

# View logs
docker compose logs -f api
```

### Both
```bash
# View running containers
docker ps

# View container logs
docker logs unientrega-api

# Execute command in container
docker exec -it unientrega-api sh
```

---

## 🎯 Which One Should I Use?

| Scenario | Use |
|----------|-----|
| Writing code, testing features | **Development** (`docker-compose.dev.yml`) |
| Deploying to server | **Production** (`docker-compose.yml`) |
| CI/CD pipeline | **Production** (`docker-compose.yml`) |
| Quick testing of changes | **Development** (`docker-compose.dev.yml`) |

---

## 🔍 Troubleshooting

### Changes not detected in dev mode?
- Make sure you're using `docker-compose.dev.yml`
- Check that Air is running: `docker logs unientrega-api-dev`
- Verify your changes are in a `.go` file

### Production image too large?
- The multi-stage build should keep it small (~20MB)
- Check with: `docker images | grep docker-api`

### Port already in use?
- Stop other instances: `docker compose down`
- Or change the port in `docker-compose.yml`: `"8081:8080"`
