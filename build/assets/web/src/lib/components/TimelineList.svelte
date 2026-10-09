<script lang="ts">
	import type { Entry } from '$lib/types';
	import Emph from './Emph.svelte';
	import Icon from './Icon.svelte';

	interface Props {
		entries: Entry[];
	}
	let { entries }: Props = $props();

	function range(e: Entry): string {
		if (e.start && e.end) return `${e.start} – ${e.end}`;
		return e.start || e.end || '';
	}
</script>

<div class="relative ml-1.5 border-l-2 border-border pl-6">
	{#each entries as e}
		<div class="relative pb-7 last:pb-0">
			<span
				class="accent-gradient absolute -left-[31px] top-1 h-3.5 w-3.5 rounded-full ring-4 ring-background"
			></span>
			<div class="flex flex-wrap items-baseline justify-between gap-x-3 gap-y-1">
				<h3 class="text-[15px] font-semibold">
					{#if e.url}
						<a href={e.url} target="_blank" rel="noopener" class="hover:text-primary hover:underline"
							>{e.role}<Icon name="external" size={12} class="ml-1 inline align-baseline" /></a
						>
					{:else}
						{e.role}
					{/if}
				</h3>
				{#if range(e)}
					<span class="whitespace-nowrap text-xs text-muted-foreground">{range(e)}</span>
				{/if}
			</div>

			{#if e.org || e.location || e.mode}
				<p class="mt-0.5 text-sm text-muted-foreground">
					{#if e.orgUrl}
						<a href={e.orgUrl} target="_blank" rel="noopener" class="font-medium hover:text-primary"
							>{e.org}</a
						>
					{:else if e.org}
						<span class="font-medium text-foreground/80">{e.org}</span>
					{/if}
					{#if e.location}<span> · {e.location}</span>{/if}
					{#if e.mode}
						<span
							class="ml-1.5 inline-block rounded border border-border px-1.5 py-0.5 text-[11px] align-middle"
							>{e.mode}</span
						>
					{/if}
				</p>
			{/if}

			{#if e.note}
				<p class="mt-1 text-sm italic text-muted-foreground">{e.note}</p>
			{/if}

			{#if e.intro}
				<p class="mt-2 text-sm leading-relaxed text-muted-foreground"><Emph text={e.intro} /></p>
			{/if}

			{#if e.bullets?.length}
				<ul class="mt-2 list-disc space-y-1 pl-5 text-sm leading-relaxed text-foreground/90">
					{#each e.bullets as b}
						<li><Emph text={b} /></li>
					{/each}
				</ul>
			{/if}
		</div>
	{/each}
</div>
