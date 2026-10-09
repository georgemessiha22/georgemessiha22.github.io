<script lang="ts">
	import { parseSpans } from '$lib/emph';

	interface Props {
		text?: string;
		class?: string;
	}
	let { text = '', class: klass = '' }: Props = $props();
	let spans = $derived(parseSpans(text));
</script>

<span class={klass}
	>{#each spans as span}{#if span.bold}<strong class="font-semibold text-foreground"
				>{span.text}</strong
			>{:else}{span.text}{/if}{/each}</span
>
