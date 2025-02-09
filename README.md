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
- [Just](https://just.systems/man/en/chapter_1.html)
  - Installation packages [here](https://just.systems/man/en/packages.html)
- [Cloudflare Account](https://dash.cloudflare.com/)
- [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler/)
  - Run `npm install -g wrangler`

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
   ```bash
   # Copy the environment template
   cp .env.template .env
   
   # Edit .env with your settings
   # Set PROXY_DOMAIN to your workers.dev domain or custom domain
   # Set PROXY_SECRET_KEY to your desired secret key
   ```

3. **Deploy**:
   ```bash
   just deploy
   ```

## 📖 Usage Guide

All requests to the proxy must include the `X-Proxy-Secret` header with your configured secret key. This authentication mechanism ensures that only authorized clients can access your proxy, preventing unauthorized usage and potential abuse of your worker's resources.

### Converting Roblox URLs to Worker Requests

To use the proxy, convert any Roblox API URL by replacing the official domain with your configured domain. The format is:

```bash
Roblox URL:    https://{subdomain}.roblox.com/{path}
Worker URL:    https://{subdomain}.your-domain.com/{path}
```

### Examples

Using curl:

```bash
# Get user details  
curl -X GET \
  -H "X-Proxy-Secret: your-secret-key" \
  "https://users.your-domain.com/v1/users/1"

# Get groups with query parameters
curl -X GET \
  -H "X-Proxy-Secret: your-secret-key" \
  "https://groups.your-domain.com/v1/groups/search?keyword=test&prioritizeExactMatch=false&limit=10"

# Get games with universe IDs
curl -X GET \
  -H "X-Proxy-Secret: your-secret-key" \
  "https://games.your-domain.com/v1/games?universeIds=1,2,3"
```

The proxy will keep all your original headers (except the secret key) and forward them to the Roblox API.

## 🛠️ Development

### Commands

```bash
# Generate config from template
just generate-config

# Start development server
just dev

# Build WebAssembly binary
just build

# Deploy to Cloudflare
just deploy
```

### Testing Dev Server

Before testing, you may want to modify the `PROXY_SECRET_KEY` in your `.dev.vars` file. By default, it's set to "development".

> [!NOTE]
> Make sure your system's hosts file allows subdomain resolution for localhost. Most modern operating systems support this by default.

When running the development server, you can access different Roblox API endpoints by using subdomains with localhost. For example:

```bash
# Test users endpoint
curl -H "X-Proxy-Secret: development" \
  "http://users.localhost:8787/v1/users/1"

# Test games endpoint
curl -H "X-Proxy-Secret: development" \
  "http://games.localhost:8787/v1/games?universeIds=1,2,3"

# Test groups endpoint
curl -H "X-Proxy-Secret: development" \
  "http://groups.localhost:8787/v1/groups/search?keyword=test"
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

It's important to protect your worker from unauthorized usage and potential costly bills.

Some good practices would be to use a **custom domain** instead of workers.dev, implement **Cloudflare's Web Application Firewall (WAF)** rules, regularly monitor your worker's metrics, and **periodic rotation of your secret keys** which minimizes the impact of potential key leaks.

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

The proxy supports common Roblox API endpoints including users, games, groups, friends, avatar, presence, and thumbnails. To add support for additional subdomains:

1. Open `wrangler.template.toml`
2. Add a new route entry following the existing pattern:
   ```toml
   { pattern = "your-subdomain.${PROXY_DOMAIN}", custom_domain = true }
   ```
3. Deploy your changes

If you find any endpoints that aren't working correctly, please [open an issue](https://github.com/robalyx/roverse/issues).

</details>

<details>
<summary>Why use TinyGo instead of regular Go?</summary>

Although the built WebAssembly binary doesn't exceed Cloudflare's free plan limit of 3MB, TinyGo is still required due to **memory constraints** in the Workers runtime environment. Regular Go's WebAssembly output includes a larger runtime and garbage collector that would increase the memory usage of Cloudflare Workers. TinyGo produces a more efficient WebAssembly binary with a smaller runtime footprint that works within these memory constraints.

</details>

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by the robalyx team
</p>
