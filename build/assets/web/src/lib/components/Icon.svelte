<script lang="ts">
	interface Props {
		name: string;
		size?: number;
		class?: string;
	}
	let { name, size = 16, class: klass = '' }: Props = $props();

	// Minimal inline icon set (stroke-based, currentColor).
	const paths: Record<string, string> = {
		mail: '<path d="M2 5h20v14H2z" fill="none"/><path d="M4 6h16a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1z"/><path d="m4 7 8 6 8-6"/>',
		download: '<path d="M12 3v12"/><path d="m7 10 5 5 5-5"/><path d="M5 21h14"/>',
		sun: '<circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4"/>',
		moon: '<path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z"/>',
		'chevron-down': '<path d="m6 9 6 6 6-6"/>',
		'arrow-right': '<path d="M5 12h14"/><path d="m12 5 7 7-7 7"/>',
		'arrow-left': '<path d="M19 12H5"/><path d="m12 19-7-7 7-7"/>',
		external: '<path d="M15 3h6v6"/><path d="M10 14 21 3"/><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>',
		'map-pin': '<path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0z"/><circle cx="12" cy="10" r="3"/>',
		github: '<path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.9a3.4 3.4 0 0 0-1-2.6c3-.3 6-1.5 6-6.5a5 5 0 0 0-1.4-3.5 4.6 4.6 0 0 0-.1-3.5s-1.1-.3-3.5 1.3a12 12 0 0 0-6 0C6.7.7 5.6 1 5.6 1a4.6 4.6 0 0 0-.1 3.5A5 5 0 0 0 4 8c0 5 3 6.2 6 6.5a3.4 3.4 0 0 0-1 2.6V21"/>',
		gitlab: '<path d="m12 21 3.4-10.5H8.6L12 21z"/><path d="M12 21 3 10.5h5.6L12 21z"/><path d="M3 10.5 4.5 5l2 5.5H3z"/><path d="M12 21 21 10.5h-5.6L12 21z"/><path d="M21 10.5 19.5 5l-2 5.5H21z"/>',
		linkedin: '<path d="M16 8a6 6 0 0 1 6 6v6h-4v-6a2 2 0 0 0-4 0v6h-4v-9h4v1.5"/><rect x="2" y="9" width="4" height="11"/><circle cx="4" cy="4" r="2"/>',
		terminal: '<path d="m4 17 6-6-6-6"/><path d="M12 19h8"/>',
		code: '<path d="m16 18 6-6-6-6"/><path d="m8 6-6 6 6 6"/>',
		globe: '<circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15 15 0 0 1 0 20 15 15 0 0 1 0-20z"/>'
	};

	let inner = $derived(paths[name] ?? paths.globe);
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	width={size}
	height={size}
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="2"
	stroke-linecap="round"
	stroke-linejoin="round"
	class={klass}
	aria-hidden="true"
>
	{@html inner}
</svg>
