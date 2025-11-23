# CI/CD Documentation

This project uses both GitHub Actions and GitLab CI for continuous integration and deployment.

## GitHub Actions

### Workflows

The GitHub Actions workflow (`.github/workflows/go.yml`) includes:

1. **Build and Test Job**
   - Runs on every push and pull request
   - Tests the code with race detection
   - Generates coverage reports
   - Builds the binary artifact

2. **Build and Push Docker Image**
   - Builds multi-arch Docker images (amd64, arm64)
   - Pushes to GitHub Container Registry (ghcr.io)
   - Only runs on main branch and tags (not on PRs)
   - Uses Docker layer caching for faster builds

3. **Security Scan**
   - Runs Trivy vulnerability scanner
   - Uploads results to GitHub Security tab

### Container Registry

Images are published to: `ghcr.io/<owner>/<repo>/moonbeam`

**Tags:**
- `latest` - Latest commit on main branch
- `<branch-name>-<sha>` - Branch-specific tags
- `v<version>` - Semantic version tags (e.g., v1.0.0)

### Pulling the Image

```bash
# Login to GHCR (if private)
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull the image
docker pull ghcr.io/<owner>/<repo>/moonbeam:latest
```

### Permissions

The workflow requires:
- `contents: read` - To checkout code
- `packages: write` - To push images to GHCR
- `security-events: write` - To upload security scan results

These are automatically granted via `GITHUB_TOKEN` for public repositories.

## GitLab CI

### Stages

1. **Test Stage**
   - Format checking
   - Vet analysis
   - Test execution with race detection
   - Coverage report generation

2. **Build Stage**
   - Compiles the Go binary
   - Stores artifacts

3. **Docker Build Stage**
   - Builds Docker image
   - Pushes to GitLab Container Registry
   - Tags with commit SHA and `latest`
   - Only runs on main branch and tags

4. **Deploy Stage**
   - Manual deployment step
   - Can be triggered from GitLab UI

### Container Registry

Images are published to: `<registry>/<project>/moonbeam`

**Tags:**
- `latest` - Latest commit on main branch
- `<commit-sha>` - Specific commit SHA

### Variables Required

- `CI_REGISTRY` - GitLab registry URL
- `CI_REGISTRY_USER` - Registry username
- `CI_REGISTRY_PASSWORD` - Registry password
- `CI_REGISTRY_IMAGE` - Full image path

These are automatically provided by GitLab CI/CD.

## Local Testing

### Test the Docker build locally

```bash
cd moonbeam
docker build -t moonbeam:test .
docker run -p 8080:8080 moonbeam:test
```

### Run tests locally

```bash
cd moonbeam
go test -v -race ./...
go test -cover ./...
```

## Versioning

For GitHub Actions, use semantic versioning tags:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

This will trigger the workflow and create versioned images:
- `ghcr.io/<owner>/<repo>/moonbeam:v1.0.0`
- `ghcr.io/<owner>/<repo>/moonbeam:v1.0`
- `ghcr.io/<owner>/<repo>/moonbeam:v1`

