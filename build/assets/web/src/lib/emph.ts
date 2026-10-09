// Parse **bold** markers into ordered spans, mirroring internal/render.ParseSpans.
export interface Span {
	text: string;
	bold: boolean;
}

export function parseSpans(input: string): Span[] {
	const spans: Span[] = [];
	let s = input;
	let bold = false;
	while (s.length > 0) {
		const idx = s.indexOf('**');
		if (idx < 0) {
			if (s) spans.push({ text: s, bold });
			break;
		}
		if (idx > 0) spans.push({ text: s.slice(0, idx), bold });
		bold = !bold;
		s = s.slice(idx + 2);
	}
	return spans;
}
