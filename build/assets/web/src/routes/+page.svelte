<script lang="ts">
	import resume from '#lib/data/resume.json';
	import type { Resume, ProjectLink } from '#lib/types';
	import { posts } from '#lib/posts';
	import { formatMonth } from '#lib/format';
	import Emph from '#lib/components/Emph.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import DownloadMenu from '#lib/components/DownloadMenu.svelte';
	import * as Card from '#lib/components/ui/card';

	const r = resume as Resume;
	const c = r.contact;
	const fullName = `${c.firstname} ${c.lastname}`.trim();

	type Social = { icon: string; label: string; href: string; external: boolean };
	const socials: Social[] = [];
	if (c.email) socials.push({ icon: 'mail', label: 'Email', href: `mailto:${c.email}`, external: false });
	if (c.socials.github)
		socials.push({ icon: 'github', label: 'GitHub', href: `https://github.com/${c.socials.github}`, external: true });
	if (c.socials.gitlab)
		socials.push({ icon: 'gitlab', label: 'GitLab', href: `https://gitlab.com/${c.socials.gitlab}`, external: true });
	if (c.socials.linkedin)
		socials.push({
			icon: 'linkedin',
			label: 'LinkedIn',
			href: `https://www.linkedin.com/in/${c.socials.linkedin}`,
			external: true
		});

	function projectIcon(p: ProjectLink): string {
		switch (p.icon) {
			case 'neovim':
				return 'code';
			case 'terminal':
				return 'terminal';
			default:
				return 'github';
		}
	}

	const latestPosts = posts.slice(0, 4);
</script>

<svelte:head>
	<title>{fullName} — {c.title}</title>
	<meta name="description" content={r.summary?.slice(0, 160)} />
</svelte:head>

<!-- Hero -->
<section class="hero-wash -mx-6 px-6 py-8 sm:rounded-xl sm:px-8">
	<div class="flex flex-col items-start gap-5 sm:flex-row sm:items-center">
		{#if c.photo}
			<img
				src={`/${c.photo}`}
				alt={fullName}
				class="h-24 w-24 rounded-full border border-border object-cover shadow-sm"
			/>
		{/if}
		<div>
			<h1 class="text-3xl font-bold tracking-tight sm:text-4xl">{fullName}</h1>
			{#if c.title}<p class="mt-1 text-lg text-muted-foreground">{c.title}</p>{/if}
			{#if c.location}
				<p class="mt-1 flex items-center gap-1.5 text-sm text-muted-foreground">
					<Icon name="map-pin" size={14} />
					{c.location}
				</p>
			{/if}
		</div>
	</div>

	<!-- Links bar -->
	<div class="mt-6 flex flex-wrap items-center gap-2">
		{#each socials as s}
			<a
				href={s.href}
				target={s.external ? '_blank' : undefined}
				rel="noopener"
				class="inline-flex items-center gap-1.5 rounded-full border border-border bg-background/70 px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-muted"
			>
				<Icon name={s.icon} size={16} />
				{s.label}
			</a>
		{/each}
		{#if r.downloads?.length}
			<DownloadMenu downloads={r.downloads} />
		{/if}
	</div>
</section>

<!-- Professional Summary -->
{#if r.summary}
	<section class="mt-10">
		<h2 class="mb-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
			Professional Summary
		</h2>
		<p class="text-[15px] leading-relaxed text-foreground/90"><Emph text={r.summary} /></p>
	</section>
{/if}

<!-- Personal Projects -->
{#if r.personalProjects?.links?.length}
	<section class="mt-10">
		<h2 class="mb-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
			Personal Projects
		</h2>
		{#if r.personalProjects.intro}
			<p class="mb-4 text-[15px] leading-relaxed text-foreground/90">
				<Emph text={r.personalProjects.intro} />
			</p>
		{/if}
		<div class="grid gap-3 sm:grid-cols-2">
			{#each r.personalProjects.links as p}
				<a href={p.url} target="_blank" rel="noopener" class="group">
					<Card.Root class="h-full gap-2 p-4 transition-colors hover:bg-muted/40">
						<div class="flex items-center gap-2 font-semibold">
							<span class="accent-gradient inline-flex size-7 items-center justify-center rounded-md">
								<Icon name={projectIcon(p)} size={16} />
							</span>
							{p.name}
							<Icon name="external" size={14} class="ml-auto text-muted-foreground" />
						</div>
						{#if p.description}
							<p class="text-sm text-muted-foreground"><Emph text={p.description} /></p>
						{/if}
					</Card.Root>
				</a>
			{/each}
		</div>
	</section>
{/if}

<!-- Blogs -->
{#if latestPosts.length}
	<section class="mt-10">
		<div class="mb-3 flex items-center justify-between">
			<h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Blogs</h2>
			{#if posts.length > latestPosts.length}
				<a href="/blog" class="text-sm font-medium text-primary hover:underline">View all →</a>
			{/if}
		</div>
		<div class="grid gap-3 sm:grid-cols-2">
			{#each latestPosts as post}
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
	</section>
{/if}
