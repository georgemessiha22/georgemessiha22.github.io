<script lang="ts">
	import { tick } from 'svelte';
	import { getPostComponent } from '#lib/posts';
	import { formatDate } from '#lib/format';
	import Icon from '#lib/components/Icon.svelte';

	let { data } = $props();
	let Post = $derived(getPostComponent(data.slug));

	let container = $state<HTMLElement | null>(null);

	// Shared wrap state: the toggle on any code block wraps ALL code blocks.
	let wrapAll = false;
	const wrapButtons: HTMLButtonElement[] = [];

	function applyWrap(on: boolean) {
		wrapAll = on;
		const root = container;
		if (!root) return;
		for (const pre of Array.from(root.querySelectorAll('pre.code-block'))) {
			pre.classList.toggle('wrap', on);
		}
		for (const b of wrapButtons) {
			b.textContent = on ? 'No wrap' : 'Wrap';
			b.setAttribute('aria-pressed', String(on));
		}
	}

	// Progressive enhancement: add line numbers + a word-wrap toggle to every
	// code block. Runs on mount and whenever the post (slug) changes.
	$effect(() => {
		data.slug; // track so we re-enhance on navigation
		enhanceCode();
	});

	async function enhanceCode() {
		await tick();
		const root = container;
		if (!root) return;
		wrapButtons.length = 0;
		for (const pre of Array.from(root.querySelectorAll('pre'))) {
			if ((pre as HTMLElement).dataset.enhanced) continue;
			(pre as HTMLElement).dataset.enhanced = '1';
			pre.classList.add('code-block');

			// Line numbers: render each line as a two-cell row (narrow number cell +
			// wide code cell) so wrapped code stays in the code column and never
			// mixes with the numbers. Numbers are CSS-generated (not copied).
			const code = pre.querySelector('code');
			if (code) {
				const lines = (code.textContent ?? '').replace(/\n$/, '').split('\n');
				code.textContent = '';
				for (const line of lines) {
					const row = document.createElement('span');
					row.className = 'code-line';
					const num = document.createElement('span');
					num.className = 'code-ln';
					num.setAttribute('aria-hidden', 'true');
					const txt = document.createElement('span');
					txt.className = 'code-txt';
					txt.textContent = line.length ? line : '\u200b';
					row.appendChild(num);
					row.appendChild(txt);
					code.appendChild(row);
				}
				pre.style.setProperty('--ln-width', `${String(lines.length).length + 1}ch`);
				pre.classList.add('with-line-numbers');
			}

			// Word-wrap toggle in the top-right of the block — toggles ALL blocks.
			const btn = document.createElement('button');
			btn.type = 'button';
			btn.className = 'code-wrap-toggle';
			btn.textContent = wrapAll ? 'No wrap' : 'Wrap';
			btn.setAttribute('aria-pressed', String(wrapAll));
			btn.addEventListener('click', () => applyWrap(!wrapAll));
			pre.insertBefore(btn, pre.firstChild);
			wrapButtons.push(btn);
		}
		// Keep newly-added blocks in sync with the current wrap state.
		applyWrap(wrapAll);
	}
</script>

<svelte:head>
	<title>{data.meta.title}</title>
	{#if data.meta.summary}<meta name="description" content={data.meta.summary} />{/if}
</svelte:head>

<a
	href="/blog/"
	class="inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
>
	<Icon name="arrow-left" size={15} /> Back to blog
</a>

<article class="mt-5">
	{#if data.meta.date}
		<time class="mb-2 block text-sm text-muted-foreground">{formatDate(data.meta.date)}</time>
	{/if}
	<div class="prose-post" bind:this={container}>
		{#if Post}
			<Post />
		{/if}
	</div>
</article>
