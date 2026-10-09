<script lang="ts">
	import { posts } from '$lib/posts';
	import { formatMonth } from '$lib/format';
	import * as Card from '$lib/components/ui/card';
</script>

<svelte:head>
	<title>Blog</title>
</svelte:head>

<h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Blog</h1>

{#if posts.length}
	<div class="mt-6 grid gap-3 sm:grid-cols-2">
		{#each posts as post}
			<a href={`/blog/${post.slug}/`} class="group">
				<Card.Root class="h-full gap-2 p-5 transition-colors hover:bg-muted/40">
					{#if post.date}
						<span class="text-[11px] uppercase tracking-wide text-muted-foreground"
							>{formatMonth(post.date)}</span
						>
					{/if}
					<span class="text-[15px] font-semibold leading-snug">{post.title}</span>
					{#if post.summary}
						<span class="text-sm text-muted-foreground">{post.summary}</span>
					{/if}
					<span class="mt-1 text-xs font-medium text-primary">Read →</span>
				</Card.Root>
			</a>
		{/each}
	</div>
{:else}
	<div
		class="mt-8 rounded-lg border border-dashed border-border p-10 text-center text-sm text-muted-foreground"
	>
		No posts yet — check back soon.
	</div>
{/if}
