<h1 align="center">
  roverse
  <br>
  <a href="https://github.com/robalyx/roverse/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/robalyx/roverse?style=flat-square&color=4a92e1">
  </a>
  <a href="https://github.com/robalyx/roverse/issues">
    <img src="https://img.shields.io/github/issues/robalyx/roverse?style=flat-square&color=4a92e1">
  </a>
</h1>

<p align="center">
  <em>A secure and efficient Cloudflare Worker proxy for Roblox API endpoints.</em>
</p>

## 📚 Table of Contents

- [🔧 How It Works](#-how-it-works)
- [📋 Requirements](#-requirements)
- [🚀 Getting Started](#-getting-started)
- [📖 Usage Guide](#-usage-guide)
- [🛠️ Development](#️-development)
- [⚠️ Pitfalls](#️-pitfalls)
- [❓ FAQ](#-faq)
- [📄 License](#-license)

## 🔧 How It Works

Roverse uses [Cloudflare Workers](https://developers.cloudflare.com/workers) to create a secure proxy layer between your application and Roblox's API endpoints. When you make a request to your worker, it forwards that request to the corresponding Roblox API endpoint while keeping all necessary headers and authentication.

## 📋 Requirements

- [Go](https://go.dev/dl/) 1.23.5 or later
- [TinyGo](https://tinygo.org/getting-started/install/) 0.29.0 or later
- [Node.js](https://nodejs.org/en/download/)
- [Cloudflare Account](https://dash.cloudflare.com/)
- [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler/)
  - just run `npm install -g wrangler`

## 🚀 Getting Started

1. **Clone and Setup**:
   ```bash
   # Clone the repository
   git clone https://github.com/robalyx/roverse.git
   cd roverse
   
   # Install dependencies
   go mod tidy
   ```

2. **Configure Environment**:
   - Set your worker name in `wrangler.toml`
   - Configure your secret key:
     ```bash
     wrangler secret put PROXY_SECRET_KEY
     ```

3. **Deploy**:
   ```bash
   make deploy
   ```

## 📖 Usage Guide

All requests to the proxy must include the `X-Proxy-Secret` header with your configured secret key. This authentication mechanism ensures that only authorized clients can access your proxy, preventing unauthorized usage and potential abuse of your worker's resources.

### Converting Roblox URLs to Worker Requests

To use the proxy, convert any Roblox API URL to a worker request by taking the subdomain and path. The format is:

```bash
Roblox URL:    https://{subdomain}.roblox.com/{path}
Worker URL:    https://your-worker.workers.dev/{subdomain}/{path}
```

### Examples

Using curl:

```bash
# Get user details  
curl -X GET \
  -H "X-Proxy-Secret: your-secret-key" \
  "https://your-worker.workers.dev/users/v1/users/1"

# Get groups with query parameters
curl -X GET \
  -H "X-Proxy-Secret: your-secret-key" \
  "https://your-worker.workers.dev/groups/v1/groups/search?keyword=test&prioritizeExactMatch=false&limit=10"

# Get games with universe IDs
curl -X GET \
  -H "X-Proxy-Secret: your-secret-key" \
  "https://your-worker.workers.dev/games/v1/games?universeIds=1,2,3"
```

The proxy will keep all your original headers (except the secret key) and forward them to the Roblox API.

## 🛠️ Development

### Commands

```bash
# Start development server
make dev

# Build WebAssembly binary
make build

# Deploy to Cloudflare
make deploy
```

### Testing Dev Server

Before testing, you may want to modify the `PROXY_SECRET_KEY` in your `.dev.vars` file. By default, it's set to "development".

You can test the dev server using curl:

```bash
# Test the proxy with the users endpoint
curl -H "X-Proxy-Secret: development" \
  http://localhost:8787/users/v1/users/1
```

## ⚠️ Pitfalls

<details>
<summary>Using workers.dev Domains</summary>

Using the default `workers.dev` domain can expose your worker to unwanted traffic. There are bots that scan for new SSL certificates and monitor these domains, looking for workers to abuse. These bots can quickly find and target your worker even before you start using it.

We **strongly recommend** using a custom domain instead of the default `workers.dev` domain. Custom domains are much less likely to be targeted by automated scanning, as they require more effort to discover and aren't immediately identifiable as Cloudflare Workers.

This is especially important if you're on the **paid plan**, as unauthorized requests will still count towards your quota even if they're blocked by your authentication. You may check the other pitfalls for more information.

</details>

<details>
<summary>Triggering Cloudflare's Abuse Protection</summary>

Cloudflare's abuse protection system may trigger if your worker **receives too many requests per second**, especially on the free plan. This may also happen if too much traffic originates from a single IP address or a small range of IPs.

If you need to handle higher request volumes, consider **upgrading to the paid Workers plan** which allows for thousands of requests per second. We recommend implementing your own rate limiting and request distribution strategies to stay within these boundaries and ensure reliable service.

There is no reason for Cloudflare to block your worker as long as you're not abusing the service. You may learn more about the limits [here](https://developers.cloudflare.com/workers/platform/limits).

</details>

<details>
<summary>Protecting Against Unauthorized Usage</summary>

I'm sure you wouldn't want to wake up to a 100k dollar bill in your bank account, so to protect your worker from unauthorized usage, you can link a custom domain and implement Cloudflare's Web Application Firewall (WAF) rules:

1. Link a custom domain to your worker in the Cloudflare Dashboard under your worker's settings at Domains & Routes.

2. Navigate to your domain settings, then Security > WAF > Custom Rules to create firewall rules specific to your hostname.

3. If you're expecting requests from a specific IP only, you can create a rule with an expression like:
   ```bash
   (ip.src ne YOUR_IP_ADDRESS and http.host wildcard "your-subdomain.example.com")
   ```
   Replace `YOUR_IP_ADDRESS` and `your-subdomain.example.com` with your actual values.

This setup helps ensure your worker's request quota isn't consumed by unauthorized traffic. Please do test to ensure that your setup is working as expected.

</details>

## ❓ FAQ

<details>
<summary>Why use a proxy for Roblox APIs?</summary>

A proxy provides additional security, rate limiting control, and also helps prevent exposure of your original IP address when making API requests.

</details>

<details>
<summary>How secure is the secret key authentication?</summary>

The secret key is stored securely in Cloudflare Workers' environment variables. It's never exposed in logs or error messages, and all requests without the correct key are immediately rejected.

</details>

<details>
<summary>What endpoints are supported?</summary>

The proxy supports all Roblox API endpoints. If you find any endpoints that aren't working correctly, please [open an issue](https://github.com/robalyx/roverse/issues) and we'll investigate it.

</details>

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by the robalyx team
</p>
