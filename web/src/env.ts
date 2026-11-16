import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	API_URL: {
		description: 'Base URL of the Forflow Go API.',
		schema: (value) => value || 'http://localhost:8080'
	}
});
