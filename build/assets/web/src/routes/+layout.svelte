<script lang="ts">
	import '../app.css';
	import type { Snippet } from 'svelte';
	import resume from '$lib/data/resume.json';
	import type { Resume } from '$lib/types';
	import Nav from '$lib/components/Nav.svelte';

	const r = resume as Resume;
	const brand = `${r.contact.firstname} ${r.contact.lastname}`.trim();
	const year = new Date().getFullYear();

	interface Props {
		children: Snippet;
	}
	let { children }: Props = $props();
</script>

<div class="flex min-h-screen flex-col">
	<Nav {brand} />
	<main class="container flex-1 py-8 sm:py-10">
		{@render children()}
	</main>
	<footer class="border-t border-border">
		<div class="container py-6 text-center text-xs text-muted-foreground">
			© {year} {brand} · Built with SvelteKit
		</div>
	</footer>
</div>
