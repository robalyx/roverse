import { env } from "cloudflare:workers";
import { Elysia } from "elysia";
import { CloudflareAdapter } from "elysia/adapter/cloudflare-worker";

const AUTH_HEADER = "x-proxy-secret";

export default new Elysia({ adapter: CloudflareAdapter })
	.onBeforeHandle({ as: "global" }, ({ request }) => {
		const secret = env.PROXY_SECRET_KEY;
		if (!secret) {
			console.error("PROXY_SECRET_KEY is not configured");
			return new Response(
				JSON.stringify({ message: "Internal Server Error" }),
				{ status: 500, headers: { "Content-Type": "application/json" } },
			);
		}

		if (request.headers.get(AUTH_HEADER) !== secret) {
			return new Response(JSON.stringify({ message: "Unauthorized" }), {
				status: 401,
				headers: { "Content-Type": "application/json" },
			});
		}
	})
	.all("*", async ({ request }) => {
		const url = new URL(request.url);

		// Extract Roblox subdomain from path
		const pathSegments = url.pathname.split("/");
		const subdomain = pathSegments[1];
		if (!subdomain) {
			return new Response(JSON.stringify({ message: "Missing subdomain" }), {
				status: 400,
				headers: { "Content-Type": "application/json" },
			});
		}
		const remainingPath = `/${pathSegments.slice(2).join("/")}`;
		const targetURL = `https://${subdomain}.roblox.com${remainingPath}${url.search}`;

		// Forward all headers except the proxy secret
		const headers = new Headers();
		for (const [key, value] of request.headers) {
			if (key !== AUTH_HEADER) {
				headers.set(key, value);
			}
		}

		// Assume JSON for POST/PUT when no Content-Type is provided
		if (
			(request.method === "POST" || request.method === "PUT") &&
			!headers.has("content-type")
		) {
			headers.set("content-type", "application/json");
		}

		// Proxy the request to the target Roblox API
		let response: Response;
		try {
			response = await fetch(targetURL, {
				method: request.method,
				headers,
				body: request.body,
				redirect: "follow",
			});
		} catch (error) {
			console.error(`Error requesting ${targetURL}:`, error);
			return new Response(JSON.stringify({ message: "Bad Gateway" }), {
				status: 502,
				headers: { "Content-Type": "application/json" },
			});
		}

		if (response.status >= 400) {
			console.error(
				`Error requesting ${targetURL} (status: ${response.status})`,
			);
		}

		return new Response(response.body, {
			status: response.status,
			headers: response.headers,
		});
	})
	.compile();
