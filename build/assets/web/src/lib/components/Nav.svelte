<script lang="ts">
	import { page } from '$app/state';
	import ThemeToggle from './ThemeToggle.svelte';
	import { cn } from '#lib/utils';

	interface Props {
		brand: string;
	}
	let { brand }: Props = $props();

	const links = [
		{ href: '/', label: 'Home' },
		{ href: '/timeline', label: 'Timeline' },
		{ href: '/blog', label: 'Blog' }
	];

	function isActive(href: string, path: string): boolean {
		if (href === '/') return path === '/';
		return path === href || path.startsWith(href + '/');
	}
</script>

<header
	class="sticky top-0 z-30 border-b border-border bg-background/80 backdrop-blur supports-[backdrop-filter]:bg-background/60"
>
	<div class="container flex h-14 items-center justify-between">
		<a href="/" class="text-[15px] font-semibold tracking-tight">{brand}</a>
		<nav class="flex items-center gap-1 sm:gap-2">
			{#each links as l}
				<a
					href={l.href}
					class={cn(
						'rounded-md px-2.5 py-1.5 text-sm transition-colors',
						isActive(l.href, page.url.pathname)
							? 'font-medium text-foreground'
							: 'text-muted-foreground hover:text-foreground'
					)}>{l.label}</a
				>
			{/each}
			<div class="ml-1"><ThemeToggle /></div>
		</nav>
	</div>
</header>
