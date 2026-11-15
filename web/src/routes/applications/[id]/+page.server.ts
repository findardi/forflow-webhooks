import { error, fail, redirect } from '@sveltejs/kit';
import { createApplication } from '#lib/server/actions.ts';
import { api, ApiError, type Application, type Endpoint } from '#lib/server/api.ts';

export const load = async ({ fetch, params, parent }) => {
	const [{ applications }, endpoints] = await Promise.all([
		parent(),
		api<Endpoint[]>(fetch, `/applications/${params.id}/endpoints`).catch((e) => {
			if (e instanceof ApiError && e.status === 400) error(404, 'Application not found.');
			throw e;
		})
	]);
	const application = applications.find((a) => a.id === params.id);
	if (!application) error(404, 'Application not found.');

	return { application, endpoints };
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

		try {
			await api<Endpoint>(fetch, `/applications/${params.id}/endpoints`, { url });
		} catch (e) {
			if (e instanceof ApiError && messages[e.status]) {
				return fail(e.status, { url, message: messages[e.status] });
			}
			throw e;
		}
		return { added: url };
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
