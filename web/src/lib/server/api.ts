import { API_URL } from '$app/env/private';
import { error } from '@sveltejs/kit';

export type Application = {
	id: string;
	name: string;
	created_at: string;
	deleted_at: string | null;
};

export type Endpoint = {
	id: string;
	application_id: string;
	url: string;
	created_at: string;
};

export type Message = {
	id: string;
	application_id: string;
	// Parsed by res.json(), so integers past 2^53 display rounded.
	payload: unknown;
	created_at: string;
};

export type DeliveryStatus = 'pending' | 'succeeded' | 'failed' | 'canceled';

/**
 * One row per delivery. Per endpoint: status only. Per application: also the delivery ID and
 * endpoint URL, all null when the message reached no endpoint.
 */
export type DeliveredMessage = Message & {
	delivery_id?: string | null;
	delivery_status: DeliveryStatus | null;
	endpoint_url?: string | null;
};

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
	}
}

const baseUrl = () => API_URL;

/**
 * GET when body is undefined, POST JSON otherwise; pass method for DELETE. Throws ApiError on non-2xx.
 * A string body is sent as already-serialized JSON, so user-typed JSON reaches the API byte for byte.
 */
export async function api<T>(
	fetch: typeof globalThis.fetch,
	path: string,
	body?: unknown,
	method = body === undefined ? 'GET' : 'POST'
) {
	let res: Response;
	try {
		res = await fetch(
			baseUrl() + path,
			body === undefined
				? { method }
				: {
						method,
						headers: { 'content-type': 'application/json' },
						body: typeof body === 'string' ? body : JSON.stringify(body)
					}
		);
	} catch {
		error(503, `Cannot reach the Forflow API at ${baseUrl()}. Is the server running?`);
	}
	if (!res.ok) throw new ApiError(res.status, (await res.text()).trim());
	if (res.status === 204) return undefined as T;
	return (await res.json()) as T;
}
