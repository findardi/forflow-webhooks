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

export const dateLong = (iso: string) => `${long.format(new Date(iso))} UTC`;
export const dateShort = (iso: string) => short.format(new Date(iso));
