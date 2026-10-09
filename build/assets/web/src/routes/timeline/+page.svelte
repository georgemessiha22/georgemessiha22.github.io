<script lang="ts">
	import resume from '$lib/data/resume.json';
	import type { Resume, Entry, VariantName } from '$lib/types';
	import TimelineList from '$lib/components/TimelineList.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';

	const r = resume as Resume;

	let variant = $state<VariantName>('summary');
	let v = $derived(r.variants[variant]);

	let certEntries = $derived<Entry[]>(
		v.certificates.map((c) => ({
			role: c.title,
			org: c.org,
			orgUrl: c.orgUrl,
			url: c.certUrl,
			location: c.location,
			start: c.start,
			end: c.end,
			note: c.group
		}))
	);
</script>

<svelte:head>
	<title>Timeline — {r.contact.firstname} {r.contact.lastname}</title>
</svelte:head>

<div class="mb-6 flex flex-wrap items-end justify-between gap-4">
	<div>
		<h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Timeline</h1>
		<p class="mt-1 text-sm text-muted-foreground">
			Career history, education, skills, and more.
		</p>
	</div>
	<!-- Summary / Detailed toggle -->
	<div class="inline-flex rounded-md border border-border p-0.5 text-sm">
		{#each ['summary', 'detailed'] as const as opt}
			<button
				type="button"
				onclick={() => (variant = opt)}
				class={cn(
					'rounded px-3 py-1.5 capitalize transition-colors',
					variant === opt
						? 'accent-gradient font-medium'
						: 'text-muted-foreground hover:text-foreground'
				)}>{opt}</button
			>
		{/each}
	</div>
</div>

<div class="space-y-10">
	{#if v.experience.length}
		<section>
			<h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
				Experience
			</h2>
			<TimelineList entries={v.experience} />
		</section>
	{/if}

	{#if v.education.length}
		<section>
			<h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
				Education
			</h2>
			<TimelineList entries={v.education} />
		</section>
	{/if}

	{#if r.skills.length}
		<section>
			<h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
				Skills
			</h2>
			<div class="space-y-3">
				{#each r.skills as group}
					<div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:gap-4">
						<span class="min-w-[150px] text-sm font-semibold">{group.category}</span>
						<div class="flex flex-wrap gap-1.5">
							{#each group.items as item}
								<Badge variant="secondary">{item}</Badge>
							{/each}
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if r.languages.length}
		<section>
			<h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
				Languages
			</h2>
			<div class="flex flex-wrap gap-1.5">
				{#each r.languages as lang}
					<Badge variant="secondary">{lang.name}{#if lang.level} — {lang.level}{/if}</Badge>
				{/each}
			</div>
		</section>
	{/if}

	{#if certEntries.length}
		<section>
			<h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
				Certificates & Awards
			</h2>
			<TimelineList entries={certEntries} />
		</section>
	{/if}

	{#if v.activities.length}
		<section>
			<h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
				Activities
			</h2>
			<TimelineList entries={v.activities} />
		</section>
	{/if}
</div>
