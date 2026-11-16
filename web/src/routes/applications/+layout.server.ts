import { api, type Application } from '#lib/server/api.ts';

export const load = async ({ fetch }) => ({
	applications: await api<Application[]>(fetch, '/applications')
});
