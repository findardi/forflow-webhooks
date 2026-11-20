// UTC on purpose: server and browser render the same string, no hydration mismatch.
const long = new Intl.DateTimeFormat('en-GB', {
	day: 'numeric',
	month: 'short',
	year: 'numeric',
	hour: '2-digit',
	minute: '2-digit',
	timeZone: 'UTC'
});
const short = new Intl.DateTimeFormat('en-GB', {
	day: 'numeric',
	month: 'short',
	hour: '2-digit',
	minute: '2-digit',
	timeZone: 'UTC'
});

// Retries can land seconds apart, so attempt times carry seconds.
const exact = new Intl.DateTimeFormat('en-GB', {
	day: 'numeric',
	month: 'short',
	hour: '2-digit',
	minute: '2-digit',
	second: '2-digit',
	timeZone: 'UTC'
});

export const dateLong = (iso: string) => `${long.format(new Date(iso))} UTC`;
export const dateShort = (iso: string) => short.format(new Date(iso));
export const dateExact = (iso: string) => exact.format(new Date(iso));
