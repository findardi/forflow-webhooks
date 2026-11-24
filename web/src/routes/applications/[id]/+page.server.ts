import { error, fail, redirect } from '@sveltejs/kit';
import { createApplication } from '#lib/server/actions.ts';
import {
	api,
	ApiError,
	type Application,
	type CreatedEndpoint,
	type DeliveredMessage,
	type Endpoint,
	type EndpointSecret,
	type Message
} from '#lib/server/api.ts';

export const load = async ({ fetch, params, parent, url }) => {
	const [{ applications }, endpoints] = await Promise.all([
		parent(),
		api<Endpoint[]>(fetch, `/applications/${params.id}/endpoints`).catch((e) => {
			if (e instanceof ApiError && e.status === 400) error(404, 'Application not found.');
			throw e;
		})
	]);
	const application = applications.find((a) => a.id === params.id);
	if (!application) error(404, 'Application not found.');

	const tab = url.searchParams.get('tab') === 'messages' ? 'messages' : 'endpoints';
	if (tab === 'endpoints') {
		return { application, endpoints, tab, endpointId: null, messages: [] as DeliveredMessage[] };
	}

	// No endpoint picked (or one since deleted): every message of the application,
	// including those that went to deleted endpoints.
	const wanted = url.searchParams.get('endpoint');
	const endpointId = endpoints.some((e) => e.id === wanted) ? wanted : null;
	const messages = await api<DeliveredMessage[]>(
		fetch,
		endpointId
			? `/applications/${params.id}/endpoints/${endpointId}/messages`
			: `/applications/${params.id}/messages`
	);

	// Newest first: the one just sent is the one you look for.
	return { application, endpoints, tab, endpointId, messages: messages.reverse() };
};

const messages: Record<number, string> = {
	400: 'Enter a full http or https URL, like https://example.com/webhooks.',
	404: 'This application no longer exists.',
	409: 'This URL is already an endpoint of this application.'
};

export const actions = {
	create: createApplication,
	add: async ({ request, fetch, params }) => {
		const url = String((await request.formData()).get('url') ?? '').trim();
		if (!url) return fail(400, { url, message: messages[400] });

		let created: CreatedEndpoint;
		try {
			created = await api<CreatedEndpoint>(fetch, `/applications/${params.id}/endpoints`, { url });
		} catch (e) {
			if (e instanceof ApiError && messages[e.status]) {
				return fail(e.status, { url, message: messages[e.status] });
			}
			throw e;
		}
		return {
			reveal: {
				kind: 'created' as const,
				endpointId: created.result.id,
				url: created.result.url,
				secret: created.secret
			}
		};
	},
	rotateSecret: async ({ request, fetch, params }) => {
		const data = await request.formData();
		const id = String(data.get('endpoint_id') ?? '');
		const url = String(data.get('url') ?? '');

		let rotated: EndpointSecret;
		try {
			rotated = await api<EndpointSecret>(
				fetch,
				`/applications/${params.id}/endpoints/${id}/secret/rotate`,
				undefined,
				'POST'
			);
		} catch (e) {
			if (e instanceof ApiError && (e.status === 400 || e.status === 404)) {
				return fail(404, { rotate: id, message: 'This endpoint no longer exists.' });
			}
			throw e;
		}
		return {
			reveal: { kind: 'rotated' as const, endpointId: id, url, secret: rotated.secret }
		};
	},
	send: async ({ request, fetch, params }) => {
		const raw = String((await request.formData()).get('payload') ?? '').trim();
		let parsed: unknown;
		try {
			parsed = JSON.parse(raw);
		} catch {
			return fail(400, { payload: raw, message: 'Payload must be valid JSON.' });
		}
		if (parsed === null) {
			return fail(400, {
				payload: raw,
				message: 'Payload cannot be null. Send an object or value.'
			});
		}

		let sent: Message;
		try {
			sent = await api<Message>(fetch, `/applications/${params.id}/messages`, `{"payload":${raw}}`);
		} catch (e) {
			if (e instanceof ApiError && e.status === 404) {
				return fail(404, { payload: raw, message: 'This application no longer exists.' });
			}
			if (e instanceof ApiError && e.status === 400) {
				return fail(400, {
					payload: raw,
					message: 'The API rejected this payload. Keep it valid JSON under 4 KB.'
				});
			}
			throw e;
		}
		return { sent: sent.id };
	},
	deleteEndpoint: async ({ request, fetch, params }) => {
		const id = String((await request.formData()).get('endpoint_id') ?? '');
		try {
			await api(fetch, `/applications/${params.id}/endpoints/${id}`, undefined, 'DELETE');
		} catch (e) {
			// 404: already gone. The reload after this action shows the truth either way.
			if (!(e instanceof ApiError && e.status === 404)) throw e;
		}
		return { deletedEndpoint: id };
	},
	deleteApplication: async ({ request, fetch, params }) => {
		const typed = String((await request.formData()).get('confirm') ?? '').trim();
		// Re-check the typed name here too: the disabled button is only a client-side hint.
		const apps = await api<Application[]>(fetch, '/applications');
		const app = apps.find((a) => a.id === params.id);
		if (app && typed !== app.name) {
			return fail(400, {
				confirm: typed,
				message: 'Type the exact application name to delete it.'
			});
		}
		if (app) {
			try {
				await api(fetch, `/applications/${params.id}`, undefined, 'DELETE');
			} catch (e) {
				if (!(e instanceof ApiError && e.status === 404)) throw e;
			}
		}
		redirect(303, '/applications');
	}
};
