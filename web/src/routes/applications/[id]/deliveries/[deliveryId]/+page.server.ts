import { error } from '@sveltejs/kit';
import { api, ApiError, type Attempt, type Delivery } from '#lib/server/api.ts';

export const load = async ({ fetch, params, parent }) => {
	const base = `/applications/${params.id}/deliveries/${params.deliveryId}`;
	const [{ applications }, delivery, attempts] = await Promise.all([
		parent(),
		api<Delivery>(fetch, base),
		api<Attempt[]>(fetch, `${base}/attempts`)
	]).catch((e) => {
		// 400 is a malformed ID in the URL; to the reader it is the same missing delivery.
		if (e instanceof ApiError && (e.status === 400 || e.status === 404)) {
			error(404, 'Delivery not found.');
		}
		throw e;
	});

	const application = applications.find((a) => a.id === params.id);
	if (!application) error(404, 'Application not found.');

	// Newest first, like messages: the latest try is the one you came for.
	return { application, delivery, attempts: attempts.reverse() };
};
