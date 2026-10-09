// Format an ISO-ish date string (YYYY-MM-DD) for display. Falls back to the
// raw string when it can't be parsed.
export function formatDate(input: string): string {
	if (!input) return '';
	const d = new Date(input);
	if (isNaN(d.getTime())) return input;
	return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}

export function formatMonth(input: string): string {
	if (!input) return '';
	const d = new Date(input);
	if (isNaN(d.getTime())) return input;
	return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
}
