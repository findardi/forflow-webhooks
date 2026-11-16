import { fail, redirect, type RequestEvent } from '@sveltejs/kit';
import { api, ApiError, type Application } from './api.ts';

/** Shared by every /applications page so the list panel's form posts to the current page. */
export async function createApplication({ request, fetch }: RequestEvent) {
	const name = String((await request.formData()).get('name') ?? '').trim();
	if (!name) return fail(400, { name, message: 'Give the application a name.' });

	let app: Application;
	try {
		app = await api<Application>(fetch, '/applications', { name });
	} catch (e) {
		if (e instanceof ApiError && e.status === 400) {
			return fail(400, { name, message: 'Name is required, up to 100 characters.' });
		}
		throw e;
	}
	redirect(303, `/applications/${app.id}`);
}
