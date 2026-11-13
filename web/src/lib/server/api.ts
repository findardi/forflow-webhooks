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

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
	}
}

const baseUrl = () => API_URL;

/** GET when body is undefined, POST JSON otherwise. Throws ApiError on non-2xx. */
export async function api<T>(fetch: typeof globalThis.fetch, path: string, body?: unknown) {
	let res: Response;
	try {
		res = await fetch(
			baseUrl() + path,
			body === undefined
				? undefined
				: {
						method: 'POST',
						headers: { 'content-type': 'application/json' },
						body: JSON.stringify(body)
					}
		);
	} catch {
		error(503, `Cannot reach the Forflow API at ${baseUrl()}. Is the server running?`);
	}
	if (!res.ok) throw new ApiError(res.status, (await res.text()).trim());
	return (await res.json()) as T;
}
