<script lang="ts">
	import { onMount } from 'svelte';
	import { Sun, Moon } from '@lucide/svelte';
	import { Button } from '#lib/components/ui/button';

	let dark = $state(false);

	onMount(() => {
		dark = document.documentElement.classList.contains('dark');
	});

	function toggle() {
		dark = !dark;
		document.documentElement.classList.toggle('dark', dark);
		try {
			localStorage.setItem('theme', dark ? 'dark' : 'light');
		} catch (e) {
			/* ignore */
		}
	}
</script>

<Button variant="outline" size="icon" onclick={toggle} aria-label="Toggle dark mode">
	{#if dark}
		<Sun />
	{:else}
		<Moon />
	{/if}
</Button>
