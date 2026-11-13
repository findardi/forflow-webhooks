import { error, fail } from '@sveltejs/kit';
import { createApplication } from '#lib/server/actions.ts';
import { api, ApiError, type Endpoint } from '#lib/server/api.ts';

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
	}
};
